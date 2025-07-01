package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/cucumber/godog"
	"github.com/francknouama/movies-mcp-server/godog-server/step_definitions"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"examples/features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(sc *godog.ScenarioContext) {
	// Create test context for this scenario
	ctx := step_definitions.NewTestContext()

	// Register step definitions for all feature areas
	step_definitions.RegisterMCPProtocolSteps(sc, ctx)
	step_definitions.RegisterMovieSteps(sc, ctx)
	step_definitions.RegisterActorSteps(sc, ctx)

	// Setup hooks for scenario lifecycle
	sc.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		// This runs before each scenario
		fmt.Printf("Starting scenario: %s\n", sc.Name)
		return ctx, nil
	})

	sc.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		// This runs after each scenario
		testCtx := step_definitions.NewTestContext()

		if err != nil {
			fmt.Printf("Scenario failed: %s - %v\n", sc.Name, err)
		} else {
			fmt.Printf("Scenario passed: %s\n", sc.Name)
		}

		// Clean up test data
		if cleanupErr := testCtx.CleanDatabase(); cleanupErr != nil {
			fmt.Printf("Warning: failed to clean database after scenario: %v\n", cleanupErr)
		}

		// Stop MCP server if running
		if stopErr := testCtx.StopMCPServer(); stopErr != nil {
			fmt.Printf("Warning: failed to stop MCP server: %v\n", stopErr)
		}

		return ctx, nil
	})
}