package bdd

import (
	"fmt"
	"testing"
	"time"

	"github.com/cucumber/godog"

	"github.com/francknouama/movies-mcp-server/mcp-server/tests/bdd/steps"
)

func TestBDDFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	// Initialize all step definitions for Phase 3 BDD infrastructure
	steps.InitializeMCPSteps(ctx)              // Common MCP protocol steps
	steps.InitializeMovieSteps(ctx)            // Movie operation steps
	steps.InitializeActorSteps(ctx)            // Actor operation steps
	steps.InitializeMCPProtocolSteps(ctx)      // MCP protocol specific steps
	steps.InitializeAdvancedSearchSteps(ctx)   // Advanced search and integration steps
}

// Optional: Run BDD tests from command line
// go test -tags bdd ./tests/bdd/...
func TestMain(m *testing.M) {
	status := godog.TestSuite{
		Name:                "movies-mcp-server BDD",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:    "pretty",
			Paths:     []string{"features"},
			Randomize: time.Now().UTC().UnixNano(), // randomize scenario execution order
		},
	}.Run()

	if status == 2 {
		// Proper exit code when no tests are found
		return
	}

	if status != 0 {
		fmt.Printf("BDD test suite failed with status: %d\n", status)
	}
}