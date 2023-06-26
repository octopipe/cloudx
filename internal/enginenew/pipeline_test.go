package engine

// Basic imports
import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PipelineTestSuite struct {
	suite.Suite
	VariableThatShouldStartAtFive int
}

func (suite *PipelineTestSuite) SetupTest() {
	suite.VariableThatShouldStartAtFive = 5
}

func (suite *PipelineTestSuite) TestExample() {
	assert.Equal(suite.T(), 5, suite.VariableThatShouldStartAtFive)
}

func TestPipelineTestSuite(t *testing.T) {
	suite.Run(t, new(PipelineTestSuite))
}
