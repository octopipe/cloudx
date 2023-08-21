package main

import (
	"context"

	"github.com/go-logr/zapr"
	"github.com/joho/godotenv"
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/annotation"
	"github.com/octopipe/cloudx/internal/controller/repository"
	"github.com/octopipe/cloudx/pkg/twice/cache"
	"github.com/octopipe/cloudx/pkg/twice/reconciler"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(commonv1alpha1.AddToScheme(scheme))
}

func main() {
	logger, _ := zap.NewProduction()
	_ = godotenv.Load()

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     "0",
		Port:                   9444,
		HealthProbeBindAddress: ":8002",
		LeaderElection:         false,
		LeaderElectionID:       "dec90f56.circlerr.io",
	})
	if err != nil {
		panic(err)
	}

	clusterCache := cache.NewLocalCache()
	reconciler := reconciler.NewReconciler(zapr.NewLogger(logger), mgr.GetConfig(), clusterCache)
	k8sClient := kubernetes.NewForConfigOrDie(mgr.GetConfig())

	repositoryController := repository.NewController(
		logger,
		mgr.GetClient(),
		mgr.GetScheme(),
		k8sClient,
	)

	if err := repositoryController.SetupWithManager(mgr); err != nil {
		panic(err)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		panic(err)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		panic(err)
	}

	logger.Info("preloading cluster cache...")
	err = reconciler.Preload(context.Background(), func(un *unstructured.Unstructured) bool {
		return un.GetAnnotations()[annotation.ManagedByAnnotation] == "cloudx"
	}, true)

	if err != nil {
		logger.Error("failed to preload", zap.Error(err))
		panic(err)
	}

	logger.Info("start repository controller")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		panic(err)
	}
}
