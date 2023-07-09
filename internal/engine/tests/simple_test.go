package tests

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/connectioninterface"
	engine "github.com/octopipe/cloudx/internal/engine"
	rpcclientmocks "github.com/octopipe/cloudx/mocks/rpcclient"
	mocks "github.com/octopipe/cloudx/mocks/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
	"k8s.io/apimachinery/pkg/types"
)

type SimpleTestSuite struct {
	suite.Suite
	VariableThatShouldStartAtFive int
}

func (suite *SimpleTestSuite) SetupTest() {
	suite.VariableThatShouldStartAtFive = 5
}

func (suite *SimpleTestSuite) TestSimpleCase() {
	assert.Equal(suite.T(), 5, suite.VariableThatShouldStartAtFive)

	logger := zaptest.NewLogger(suite.T())

	absPath, _ := filepath.Abs("./data/simple-infra.json")
	simpleInfraJSON, _ := ioutil.ReadFile(absPath)

	currentSharedInfra := &commonv1alpha1.SharedInfra{}
	err := json.Unmarshal(simpleInfraJSON, currentSharedInfra)

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

	// connectionInterface := &commonv1alpha1.ConnectionInterface{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name:      "ecr-test",
	// 		Namespace: "default",
	// 	},
	// 	Spec: commonv1alpha1.ConnectionInterfaceSpec{
	// 		Outputs: []commonv1alpha1.ConnectionInterfaceSpecItem{
	// 			{Key: "repository_url", Value: "repository.org"},
	// 		},
	// 	},
	// }

	fakeRpcClient := new(rpcclientmocks.Client)

	fakeRpcClient.On(
		"Call",
		"GetConnectionInterface",
		connectioninterface.RPCGetConnectionInterfaceArgs{Ref: types.NamespacedName{Namespace: "default", Name: "ecr-test"}},
		commonv1alpha1.ConnectionInterface{},
	).Return(nil)

	newEngine := engine.NewEngine(logger, fakeRpcClient, terraformProvider)

	assert.NoError(suite.T(), err)

	chann := make(chan commonv1alpha1.ExecutionStatus)

	status := newEngine.Apply(*currentSharedInfra, chann)
	assert.Empty(suite.T(), status.Error)
	assert.Equal(suite.T(), engine.ExecutionSuccessStatus, status.Status)
}

func TestSimpleTestSuite(t *testing.T) {
	suite.Run(t, new(SimpleTestSuite))
}
