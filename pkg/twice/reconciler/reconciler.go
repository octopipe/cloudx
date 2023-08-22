package reconciler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"time"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/go-logr/logr"
	"github.com/octopipe/cloudx/pkg/twice/cache"
	"github.com/octopipe/cloudx/pkg/twice/resource"
	"golang.org/x/sync/errgroup"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	clientCache "k8s.io/client-go/tools/cache"
	watchutil "k8s.io/client-go/tools/watch"
	"k8s.io/client-go/util/retry"
)

const (
	LastAppliedConfigurationAnnotation = "twice.io/last-applied-configuration"
)

const (
	PlanImmutableAction = "IMMUTABLE"
	PlanUpdateAction    = "UPDATE"
	PlanCreateAction    = "CREATE"
	PlanDeleteAction    = "DELETE"
)

var (
	ignoredResources = map[string]bool{
		"events": true,
	}
)

type Reconciler interface {
	Preload(ctx context.Context, isManaged isManagedFunc, liveUpdate bool) error
	Plan(ctx context.Context, manifests []string, namespace string, isManaged isManagedFunc) ([]PlanResult, error)
	Apply(ctx context.Context, planResults []PlanResult, namespace string, metadata map[string]string) (ApplyResult, error)
}

type PlanResult struct {
	resource.Resource
	action         string
	srcManifest    string
	targetManifest string
	diffString     []string
	err            error
}

type ApplyResult struct{}

type isManagedFunc func(un *unstructured.Unstructured) bool

type reconciler struct {
	logger logr.Logger
	config *rest.Config
	cache  cache.Cache

	dynamicClient   *dynamic.DynamicClient
	discoveryClient *discovery.DiscoveryClient
}

func NewReconciler(logger logr.Logger, config *rest.Config, cache cache.Cache) Reconciler {
	dynamicClient := dynamic.NewForConfigOrDie(config)
	discoveryClient := discovery.NewDiscoveryClientForConfigOrDie(config)

	return reconciler{
		logger:          logger,
		config:          config,
		cache:           cache,
		dynamicClient:   dynamicClient,
		discoveryClient: discoveryClient,
	}
}

func isSupportedVerb(verbs []string) bool {
	foundList := false
	foundWatch := false
	for _, verb := range verbs {
		if verb == "list" {
			foundList = true
			continue
		}

		if verb == "watch" {
			foundWatch = true
			continue
		}
	}

	return foundList && foundWatch
}

func (r reconciler) syncCache(ctx context.Context, resourceList *v1.APIResourceList, apiResource v1.APIResource, liveUpdate bool, isManaged isManagedFunc) func() error {
	return func() error {
		gvk := schema.FromAPIVersionAndKind(resourceList.GroupVersion, apiResource.Kind)
		gvr := gvk.GroupVersion().WithResource(apiResource.Name)

		dynamicInterface := r.dynamicClient.Resource(gvr)
		uns, err := dynamicInterface.List(ctx, v1.ListOptions{})
		if err != nil {
			return err
		}

		for _, un := range uns.Items {
			if !isManaged(&un) {
				continue
			}

			isManaged := isManaged(&un)
			newRes := resource.NewResourceByUnstructured(un, un.GetNamespace(), apiResource.Name, isManaged)
			r.cache.Set(newRes.GetResourceIdentifier(), newRes)
		}

		if liveUpdate {
			go r.watch(ctx, uns.GetResourceVersion(), apiResource.Name, dynamicInterface, isManaged)
		}

		return nil
	}
}

func (r reconciler) watch(ctx context.Context, resourceVersion string, apiResourceName string, dynamicInterface dynamic.NamespaceableResourceInterface, isManaged isManagedFunc) {
	wait.PollImmediateUntil(time.Second*3, func() (bool, error) {
		w, err := watchutil.NewRetryWatcher(resourceVersion, &clientCache.ListWatch{
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				res, err := dynamicInterface.Watch(ctx, options)
				if k8sErrors.IsNotFound(err) {
					fmt.Println("RES NOT FOUND")
				}

				return res, err
			},
		})
		if err != nil {
			return false, err
		}

		defer w.Stop()

		for {
			select {
			case <-ctx.Done():
				return true, nil
			case <-w.Done():
				return false, errors.New("was done on init")
			case event, ok := <-w.ResultChan():
				if !ok {
					return false, errors.New("was closed on init")
				}

				obj, ok := event.Object.(*unstructured.Unstructured)
				if !ok {
					return false, errors.New("was closed")
				}

				res := resource.NewResourceByUnstructured(*obj, obj.GetNamespace(), apiResourceName, isManaged(obj))
				key := res.GetResourceIdentifier()
				if event.Type == watch.Deleted && r.cache.Has(key) {
					r.cache.Delete(key)
				} else {
					r.cache.Set(key, res)
				}
			}
		}
	}, ctx.Done())
}

func (r reconciler) Preload(ctx context.Context, isManaged isManagedFunc, liveUpdate bool) error {
	apiResouceList, err := r.discoveryClient.ServerPreferredResources()
	if err != nil {
		return err
	}

	g, _ := errgroup.WithContext(ctx)
	for _, resourceList := range apiResouceList {
		for _, apiResource := range resourceList.APIResources {
			if _, ok := ignoredResources[apiResource.Name]; ok || !isSupportedVerb(apiResource.Verbs) {
				continue
			}

			g.Go(r.syncCache(ctx, resourceList, apiResource, liveUpdate, isManaged))
		}
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (r reconciler) splitManifest(manifest []byte) ([][]byte, error) {
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(manifest), 4096)
	manifests := [][]byte{}

	for {
		newManifest := runtime.RawExtension{}
		err := decoder.Decode(&newManifest)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		newManifest.Raw = bytes.TrimSpace(newManifest.Raw)
		if len(newManifest.Raw) == 0 || bytes.Equal(newManifest.Raw, []byte("null")) {
			continue
		}

		manifests = append(manifests, newManifest.Raw)
	}

	return manifests, nil
}

func (r reconciler) deserializer(manifest []byte) (*unstructured.Unstructured, error) {
	un := &unstructured.Unstructured{}
	if err := json.Unmarshal([]byte(manifest), un); err != nil {
		return nil, err
	}

	return un, nil
}

func (r reconciler) getLastAppliedConfiguration(un *unstructured.Unstructured) string {
	annotations := un.GetAnnotations()

	kubectlLastAppliedConfigurationAnnotation := "kubectl.kubernetes.io/last-applied-configuration"
	kubectlLastAppliedConfiguration, ok := annotations[kubectlLastAppliedConfigurationAnnotation]
	if ok {
		return kubectlLastAppliedConfiguration
	}

	twiceLastAppliedConfig, ok := annotations[LastAppliedConfigurationAnnotation]
	if ok {
		return twiceLastAppliedConfig
	}

	return ""
}

func removeNulls(m map[string]interface{}) {
	val := reflect.ValueOf(m)
	for _, e := range val.MapKeys() {
		v := val.MapIndex(e)
		if v.IsNil() {
			delete(m, e.String())
			continue
		}
		switch t := v.Interface().(type) {
		// If key is a JSON object (Go Map), use recursion to go deeper
		case map[string]interface{}:
			removeNulls(t)
		}
	}
}

func (r reconciler) getMergePatch(original []byte, modified []byte) ([]byte, error) {
	patch, err := jsonpatch.CreateMergePatch(original, modified)
	if err != nil {
		return nil, err
	}

	p := map[string]interface{}{}
	err = json.Unmarshal(patch, &p)
	if err != nil {
		return nil, err
	}

	removeNulls(p)

	l, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (r reconciler) getResourceName(un *unstructured.Unstructured) (string, error) {
	apiResourceList, err := r.discoveryClient.ServerResourcesForGroupVersion(un.GroupVersionKind().GroupVersion().String())
	if err != nil {
		return "", err
	}

	for _, apiResource := range apiResourceList.APIResources {
		if !isSupportedVerb(apiResource.Verbs) {
			continue
		}

		if apiResource.Kind == un.GetKind() {
			return apiResource.Name, nil
		}
	}

	return "", errors.New("server resource not supported")
}

func (r reconciler) Plan(ctx context.Context, manifests []string, namespace string, isManaged isManagedFunc) ([]PlanResult, error) {
	// TODO: Convert manifests to objects ()
	// TODO: get all managed resources from cache
	// TODO: make diff from cache with managed resources
	// TODO: identify action by diff's (create, update or delete)
	// TODO: return plan result

	allManifests := [][]byte{}
	result := []PlanResult{}

	for _, m := range manifests {
		splitedManifest, err := r.splitManifest([]byte(m))
		if err != nil {
			return nil, err
		}

		allManifests = append(allManifests, splitedManifest...)
	}

	for _, m := range allManifests {
		un, err := r.deserializer(m)
		if err != nil {
			return nil, err
		}

		resourceName, err := r.getResourceName(un)
		if err != nil {
			return nil, err
		}

		res := resource.NewResourceByUnstructured(*un, namespace, resourceName, true)

		if !r.cache.Has(res.GetResourceIdentifier()) {
			result = append(result, PlanResult{
				Resource:       res,
				action:         PlanCreateAction,
				srcManifest:    string(m),
				targetManifest: string(m),
				diffString:     []string{},
			})
			continue
		}

		if err != nil {
			return nil, err
		}

		currentResource := r.cache.Get(res.GetResourceIdentifier())
		lastAppliedConfiguration := r.getLastAppliedConfiguration(currentResource.Object)
		patch, err := r.getMergePatch([]byte(lastAppliedConfiguration), m)
		if err != nil {
			return nil, err
		}

		currentAction := PlanImmutableAction
		target := []byte(lastAppliedConfiguration)
		if string(patch) != "{}" {
			currentAction = PlanUpdateAction

			target, err = jsonpatch.MergePatch(target, patch)
			if err != nil {
				return nil, err
			}
		}

		result = append(result, PlanResult{
			Resource:       res,
			action:         currentAction,
			srcManifest:    string(m),
			targetManifest: string(target),
			diffString:     []string{},
		})
	}

	resultsForDeletion := r.getPlanResultsForDeletion(isManaged, result)
	result = append(result, resultsForDeletion...)

	return result, nil
}

func (r *reconciler) getPlanResultsForDeletion(isManaged isManagedFunc, currentResults []PlanResult) []PlanResult {
	result := []PlanResult{}
	cachedResources := r.cache.List(func(res resource.Resource) bool {
		return res.Object != nil && isManaged(res.Object)
	})
	for _, cachedKey := range cachedResources {
		forDeletion := true
		for _, planResult := range currentResults {
			if cachedKey == planResult.GetResourceIdentifier() {
				forDeletion = false
			}
		}

		if forDeletion {
			cachedItem := r.cache.Get(cachedKey)

			isControlled := false
			// Verifying if cached resource has a controller to prevent accidentally deleting
			for _, owner := range cachedItem.Object.GetOwnerReferences() {
				isController := owner.Controller
				if isController != nil && *isController {
					isControlled = true
					break
				}
			}

			if !isControlled {
				result = append(result, PlanResult{
					Resource:       cachedItem,
					action:         PlanDeleteAction,
					srcManifest:    r.getLastAppliedConfiguration(cachedItem.Object),
					targetManifest: "",
				})
			}
		}
	}

	return result
}

func applyMetadataToObject(planResult PlanResult, un *unstructured.Unstructured, metadata map[string]string) *unstructured.Unstructured {

	currentAnnotations := un.GetAnnotations()
	if currentAnnotations == nil {
		currentAnnotations = map[string]string{}
	}

	for k, v := range metadata {
		currentAnnotations[k] = v
	}

	currentAnnotations[LastAppliedConfigurationAnnotation] = string(planResult.targetManifest)
	un.SetAnnotations(currentAnnotations)

	return un
}

func (r reconciler) Apply(ctx context.Context, planResults []PlanResult, namespace string, metadata map[string]string) (ApplyResult, error) {
	// TODO: create, update or delete objects by plan result (add snapshots and
	// metadata)
	// TODO: update cache after apply
	// TODO: return apply result

	if len(planResults) <= 0 {
		return ApplyResult{}, nil
	}

	for _, res := range planResults {

		dynamicInterface := r.dynamicClient.Resource(schema.GroupVersionResource{
			Group:    res.Group,
			Version:  res.Version,
			Resource: res.ResourceName,
		}).Namespace(namespace)

		if res.action == PlanCreateAction {
			modifiedObject := applyMetadataToObject(res, res.Object, metadata)
			_, err := dynamicInterface.Create(ctx, modifiedObject, v1.CreateOptions{})

			if err != nil {
				return ApplyResult{}, err
			}

			r.cache.Set(res.GetResourceIdentifier(), res.Resource)
		}

		if res.action == PlanUpdateAction {
			err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {

				target := &unstructured.Unstructured{}
				err := json.Unmarshal([]byte(res.targetManifest), target)
				if err != nil {
					return err
				}

				curr, err := dynamicInterface.Get(ctx, target.GetName(), v1.GetOptions{})
				if err != nil {
					return err
				}

				target.SetResourceVersion(curr.GetResourceVersion())

				modifiedObject := applyMetadataToObject(res, target, metadata)
				_, err = dynamicInterface.Update(ctx, modifiedObject, v1.UpdateOptions{})
				if err != nil {
					return err
				}

				return nil
			})
			if err != nil {
				return ApplyResult{}, err
			}

			r.cache.Set(res.GetResourceIdentifier(), res.Resource)
		}

		if res.action == PlanDeleteAction {
			err := dynamicInterface.Delete(ctx, res.Name, v1.DeleteOptions{})
			if err != nil {
				return ApplyResult{}, err
			}

			r.cache.Delete(res.GetResourceIdentifier())
		}

	}

	return ApplyResult{}, nil
}
