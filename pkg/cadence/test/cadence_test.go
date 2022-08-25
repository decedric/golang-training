package cadence_test

import (
	"testing"

	c "owt/fibonacci/pkg/cadence"

	"github.com/stretchr/testify/suite"
	"go.uber.org/cadence/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UnitTestSuite) Test_StartFibonacciWorkflow() {
	s.env.ExecuteWorkflow(c.StartFibonacciWorkflow, "100", "testId")
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	s.True(c.GetResult("testId").Text(10) == "354224848179261915075")
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
