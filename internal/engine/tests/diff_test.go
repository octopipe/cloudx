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
	terraformmocks "github.com/octopipe/cloudx/mocks/terraform"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type DiffTestSuite struct {
	suite.Suite
	VariableThatShouldStartAtFive int
}

func (suite *DiffTestSuite) SetupTest() {
	suite.VariableThatShouldStartAtFive = 5
}

func (suite *DiffTestSuite) TestSimpleCase() {
	assert.Equal(suite.T(), 5, suite.VariableThatShouldStartAtFive)

	logger := zaptest.NewLogger(suite.T())

	absPath, _ := filepath.Abs("./data/simple-infra.json")
	simpleInfraJSON, _ := ioutil.ReadFile(absPath)

	currentInfra := &commonv1alpha1.Infra{}
	err := json.Unmarshal(simpleInfraJSON, currentInfra)

	terraformProvider := new(terraformmocks.TerraformProvider)

	terraformProvider.On("Destroy",
		"mayconjrpacheco/task:sns-1",
		[]commonv1alpha1.InfraTaskInput{},
		"",
		"",
	).Return(nil)

	terraformProvider.On("Apply",
		"mayconjrpacheco/task:sns-1",
		currentInfra.Spec.Tasks[0].Inputs,
		"",
		"",
	).Return(map[string]any{
		"arn": "arn:aws:sns-arn",
	}, "", "", nil)

	terraformProvider.On("Apply",
		"mayconjrpacheco/task:lambda-role-1",
		currentInfra.Spec.Tasks[1].Inputs,
		"",
		"",
	).Return(map[string]any{
		"arn": "arn:aws:role-arn",
	}, "", "", nil)

	terraformProvider.On("Apply",
		"mayconjrpacheco/task:lambda-1",
		[]commonv1alpha1.InfraTaskInput{
			currentInfra.Spec.Tasks[2].Inputs[0],
			currentInfra.Spec.Tasks[2].Inputs[1],
			{Key: "role_arn", Value: "arn:aws:role-arn"},
			{Key: "image_uri", Value: "repository.org:latest"},
		},
		"",
		"",
	).Return(map[string]any{
		"arn": "role-arn",
	}, "", "", nil)

	terraformProvider.On("Apply",
		"mayconjrpacheco/task:sns-lambda-trigger-1",
		[]commonv1alpha1.InfraTaskInput{
			currentInfra.Spec.Tasks[3].Inputs[0],
			{Key: "sns_arn", Value: "arn:aws:sns-arn"},
			currentInfra.Spec.Tasks[3].Inputs[2],
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
	fakeRpcClient := new(rpcclientmocks.Client)
	fakeRpcClient.On(
		"Call",
		"GetConnectionInterface",
		connectioninterface.RPCGetConnectionInterfaceArgs{Ref: types.NamespacedName{Namespace: "default", Name: "ecr-test"}},
		connectionInterface,
	).Return(nil)
	newEngine := engine.NewEngine(logger, fakeRpcClient, terraformProvider)

	assert.NoError(suite.T(), err)

	lastExecution := commonv1alpha1.ExecutionStatus{
		Tasks: []commonv1alpha1.TaskExecutionStatus{
			{Name: "demo-1-sns", TaskType: "terraform", Inputs: []commonv1alpha1.InfraTaskInput{}},
			{Name: "demo-1-lambda-role", TaskType: "terraform"},
			{Name: "demo-2-lambda", TaskType: "terraform"},
			{Name: "demo-1-lambda-sns-trigger", TaskType: "terraform"},
			{Name: "demo-1-after-1", TaskType: "terraform", Ref: "mayconjrpacheco/task:sns-1", Inputs: []commonv1alpha1.InfraTaskInput{}},
			{Name: "demo-1-after", Ref: "mayconjrpacheco/task:sns-1", Depends: []string{"demo-1-after-1"}, TaskType: "terraform", Inputs: []commonv1alpha1.InfraTaskInput{}},
		},
	}

	chann := make(chan commonv1alpha1.ExecutionStatus)

	currentInfra.Status.LastExecution = lastExecution

	status := newEngine.Apply(*currentInfra, chann)
	assert.Empty(suite.T(), status.Error)
	assert.Equal(suite.T(), engine.ExecutionSuccessStatus, status.Status)
}

func TestDiffTestSuite(t *testing.T) {
	suite.Run(t, new(DiffTestSuite))
}
