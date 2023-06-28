package engine

import (
	"encoding/json"
	"testing"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	mocks "github.com/octopipe/cloudx/mocks/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var demo1JSON = `
{
	"apiVersion": "commons.cloudx.io/v1alpha1",
	"kind": "SharedInfra",
	"metadata": {
			"annotations": {
					"kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"commons.cloudx.io/v1alpha1\",\"kind\":\"SharedInfra\",\"metadata\":{\"annotations\":{},\"labels\":{\"revision\":\"0.0.1\"},\"name\":\"demo-1\",\"namespace\":\"default\"},\"spec\":{\"author\":\"Maycon Pacheco\",\"description\":\"SharedInfra for example 2\",\"plugins\":[{\"depends\":[],\"inputs\":[{\"key\":\"name\",\"value\":\"demo-1-sns\"},{\"key\":\"region\",\"value\":\"us-east-1\"}],\"name\":\"demo-1-sns\",\"outputs\":[],\"ref\":\"mayconjrpacheco/plugin:sns-1\",\"type\":\"terraform\"},{\"depends\":[],\"inputs\":[{\"key\":\"lambda_name\",\"value\":\"demo-1-function\"},{\"key\":\"region\",\"value\":\"us-east-1\"}],\"name\":\"demo-1-lambda-role\",\"outputs\":[],\"ref\":\"mayconjrpacheco/plugin:lambda-role-1\",\"type\":\"terraform\"},{\"depends\":[\"demo-1-lambda-role\"],\"inputs\":[{\"key\":\"name\",\"value\":\"demo-2-function\"},{\"key\":\"region\",\"value\":\"us-east-1\"},{\"key\":\"role_arn\",\"value\":\"{{ this.demo-1-lambda-role.arn }}\"},{\"key\":\"image_uri\",\"value\":\"{{ connection-interfaces.ecr-test.repository_url }}\"}],\"name\":\"demo-2-lambda\",\"outputs\":[],\"ref\":\"mayconjrpacheco/plugin:lambda-1\",\"type\":\"terraform\"},{\"depends\":[\"demo-2-lambda\",\"demo-1-sns\"],\"inputs\":[{\"key\":\"lambda_name\",\"value\":\"demo-1-function\"},{\"key\":\"sns_arn\",\"value\":\"{{ demo-1-sns.outputs.arn }}\"},{\"key\":\"region\",\"value\":\"us-east-1\"}],\"name\":\"demo-1-lambda-sns-trigger\",\"outputs\":[],\"ref\":\"mayconjrpacheco/plugin:sns-lambda-trigger-1\",\"type\":\"terraform\"}],\"providerConfigRef\":{\"name\":\"aws-config\",\"namespace\":\"default\"}}}\n"
			},
			"creationTimestamp": "2023-06-27T13:49:58Z",
			"generation": 1,
			"labels": {
					"revision": "0.0.1"
			},
			"name": "demo-1",
			"namespace": "default",
			"resourceVersion": "663968",
			"uid": "8c385578-e030-4240-bbab-2da0cd9ada27"
	},
	"spec": {
			"author": "Maycon Pacheco",
			"description": "SharedInfra for example 2",
			"plugins": [
					{
							"depends": [],
							"inputs": [
									{
											"key": "name",
											"value": "demo-1-sns"
									},
									{
											"key": "region",
											"value": "us-east-1"
									}
							],
							"name": "demo-1-sns",
							"outputs": [],
							"ref": "mayconjrpacheco/plugin:sns-1",
							"type": "terraform"
					},
					{
							"depends": [],
							"inputs": [
									{
											"key": "lambda_name",
											"value": "demo-1-function"
									},
									{
											"key": "region",
											"value": "us-east-1"
									}
							],
							"name": "demo-1-lambda-role",
							"outputs": [],
							"ref": "mayconjrpacheco/plugin:lambda-role-1",
							"type": "terraform"
					},
					{
							"depends": [
									"demo-1-lambda-role"
							],
							"inputs": [
									{
											"key": "name",
											"value": "demo-2-function"
									},
									{
											"key": "region",
											"value": "us-east-1"
									},
									{
											"key": "role_arn",
											"value": "{{ this.demo-1-lambda-role.arn }}"
									},
									{
											"key": "image_uri",
											"value": "{{ connection-interface.ecr-test.repository_url }}:latest"
									}
							],
							"name": "demo-2-lambda",
							"outputs": [],
							"ref": "mayconjrpacheco/plugin:lambda-1",
							"type": "terraform"
					},
					{
							"depends": [
									"demo-2-lambda",
									"demo-1-sns"
							],
							"inputs": [
									{
											"key": "lambda_name",
											"value": "demo-1-function"
									},
									{
											"key": "sns_arn",
											"value": "{{ this.demo-1-sns.arn }}"
									},
									{
											"key": "region",
											"value": "us-east-1"
									}
							],
							"name": "demo-1-lambda-sns-trigger",
							"outputs": [],
							"ref": "mayconjrpacheco/plugin:sns-lambda-trigger-1",
							"type": "terraform"
					}
			],
			"providerConfigRef": {
					"name": "aws-config",
					"namespace": "default"
			}
	},
	"status": {
			"executions": [
					{}
			]
	}
}
`

type EngineTestSuite struct {
	suite.Suite
	VariableThatShouldStartAtFive int
}

func (suite *EngineTestSuite) SetupTest() {
	suite.VariableThatShouldStartAtFive = 5
}

func (suite *EngineTestSuite) TestExample() {
	assert.Equal(suite.T(), 5, suite.VariableThatShouldStartAtFive)

	logger := zaptest.NewLogger(suite.T())

	currentSharedInfra := &commonv1alpha1.SharedInfra{}
	err := json.Unmarshal([]byte(demo1JSON), currentSharedInfra)

	terraformProvider := new(mocks.TerraformProvider)

	terraformProvider.On("Apply",
		"mayconjrpacheco/plugin:sns-1",
		currentSharedInfra.Spec.Plugins[0].Inputs,
		"",
		"",
	).Return(map[string]any{
		"arn": "arn:aws:sns-arn",
	}, "", "", nil)

	terraformProvider.On("Apply",
		"mayconjrpacheco/plugin:lambda-role-1",
		currentSharedInfra.Spec.Plugins[1].Inputs,
		"",
		"",
	).Return(map[string]any{
		"arn": "arn:aws:role-arn",
	}, "", "", nil)

	terraformProvider.On("Apply",
		"mayconjrpacheco/plugin:lambda-1",
		[]commonv1alpha1.SharedInfraPluginInput{
			currentSharedInfra.Spec.Plugins[2].Inputs[0],
			currentSharedInfra.Spec.Plugins[2].Inputs[1],
			{Key: "role_arn", Value: "arn:aws:role-arn"},
			{Key: "image_uri", Value: "repository.org:latest"},
		},
		"",
		"",
	).Return(map[string]any{
		"arn": "role-arn",
	}, "", "", nil)

	terraformProvider.On("Apply",
		"mayconjrpacheco/plugin:sns-lambda-trigger-1",
		[]commonv1alpha1.SharedInfraPluginInput{
			currentSharedInfra.Spec.Plugins[3].Inputs[0],
			{Key: "sns_arn", Value: "arn:aws:sns-arn"},
			currentSharedInfra.Spec.Plugins[3].Inputs[2],
		},
		"",
		"",
	).Return(map[string]any{
		"arn": "role-arn",
	}, "", "", nil)

	connectionInterface := &commonv1alpha1.ConnectionInterface{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ecr-test",
			Namespace: "default",
		},
		Spec: commonv1alpha1.ConnectionInterfaceSpec{
			Outputs: []commonv1alpha1.ConnectionInterfaceSpecItem{
				{Key: "repository_url", Value: "repository.org"},
			},
		},
	}
	scheme := runtime.NewScheme()

	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(commonv1alpha1.AddToScheme(scheme))

	objs := []runtime.Object{connectionInterface}
	fakeClient := fake.NewClientBuilder()
	fakeClient.WithScheme(scheme)
	fakeClient.WithRuntimeObjects(objs...)

	engine := NewEngine(logger, fakeClient.Build(), terraformProvider)

	assert.NoError(suite.T(), err)

	status := engine.Apply(commonv1alpha1.Execution{}, *currentSharedInfra)
	assert.Empty(suite.T(), status.Error)
}

func TestEngineTestSuite(t *testing.T) {
	suite.Run(t, new(EngineTestSuite))
}
