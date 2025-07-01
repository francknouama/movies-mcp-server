package steps

import (
	"context"
	"fmt"
	"time"

	"github.com/cucumber/godog"
	
	bddContext "github.com/francknouama/movies-mcp-server/mcp-server/tests/bdd/context"
	"github.com/francknouama/movies-mcp-server/mcp-server/tests/bdd/support"
)

// DatabaseInterface defines common database operations for both implementations
type DatabaseInterface interface {
	LoadFixtures(fixtureName string) error
	CleanupAfterScenario() error
	CountRows(table string, whereClause string, args ...interface{}) (int, error)
	VerifyMovieExists(title, director string, year int) (bool, error)
	VerifyActorExists(name string, birthYear int) (bool, error)
}

// CommonStepContext provides shared step context for all BDD scenarios
type CommonStepContext struct {
	bddContext      *bddContext.BDDContext
	testDB          DatabaseInterface
	testContainerDB *support.TestContainerDatabase
	dataManager     *support.TestDataManager
	ctx             context.Context
}

// NewCommonStepContext creates a new common step context
func NewCommonStepContext() *CommonStepContext {
	return &CommonStepContext{
		bddContext:  bddContext.NewBDDContext(),
		dataManager: support.NewTestDataManager(),
		ctx:         context.Background(),
	}
}

// InitializeMCPSteps registers common MCP protocol step definitions
func InitializeMCPSteps(ctx *godog.ScenarioContext) {
	stepContext := NewCommonStepContext()

	// Setup and teardown hooks
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		return ctx, stepContext.setupScenario()
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		return ctx, stepContext.teardownScenario()
	})

	// Common step definitions
	ctx.Step(`^the MCP server is running$`, stepContext.theMCPServerIsRunning)
	ctx.Step(`^the MCP connection is initialized$`, stepContext.theMCPConnectionIsInitialized)
	ctx.Step(`^the database is clean$`, stepContext.theDatabaseIsClean)
	ctx.Step(`^the response should be successful$`, stepContext.theResponseShouldBeSuccessful)
	ctx.Step(`^the response should contain an error$`, stepContext.theResponseShouldContainAnError)
	ctx.Step(`^the response should contain an error with message "([^"]*)"$`, stepContext.theResponseShouldContainAnErrorWithMessage)
}

// setupScenario initializes the scenario context
func (c *CommonStepContext) setupScenario() error {
	var err error
	
	// Try to initialize Testcontainer database first, fallback to regular database
	c.testContainerDB, err = support.NewTestContainerDatabase(c.ctx)
	if err != nil {
		fmt.Printf("Warning: Could not initialize Testcontainer database (%v), falling back to regular database\n", err)
		
		// Fallback to regular test database
		regularDB, fallbackErr := support.NewTestDatabase()
		if fallbackErr != nil {
			return fmt.Errorf("failed to initialize any test database - Testcontainer: %v, Regular: %w", err, fallbackErr)
		}
		c.testDB = regularDB
	} else {
		// Use Testcontainer database
		c.testDB = c.testContainerDB
	}

	// Clean database before each scenario
	err = c.testDB.CleanupAfterScenario()
	if err != nil {
		return fmt.Errorf("failed to clean database before scenario: %w", err)
	}

	// Configure database environment for MCP server if using Testcontainer
	if c.testContainerDB != nil {
		connStr, err := c.testContainerDB.GetConnectionString()
		if err != nil {
			return fmt.Errorf("failed to get container connection string: %w", err)
		}
		
		// Set database environment for the MCP server
		err = c.bddContext.SetDatabaseEnvironment(connStr)
		if err != nil {
			return fmt.Errorf("failed to set database environment: %w", err)
		}
	}

	// Start MCP server
	err = c.bddContext.StartMCPServer()
	if err != nil {
		return fmt.Errorf("failed to start MCP server: %w", err)
	}

	// Wait for server to be ready
	err = c.bddContext.WaitForServer(10 * time.Second)
	if err != nil {
		return fmt.Errorf("MCP server not ready: %w", err)
	}

	return nil
}

// teardownScenario cleans up after each scenario
func (c *CommonStepContext) teardownScenario() error {
	var errors []error

	// Clean up BDD context
	if c.bddContext != nil {
		if err := c.bddContext.Cleanup(); err != nil {
			errors = append(errors, fmt.Errorf("BDD context cleanup failed: %w", err))
		}
	}

	// Clean up test database
	if c.testDB != nil {
		if err := c.testDB.CleanupAfterScenario(); err != nil {
			errors = append(errors, fmt.Errorf("database cleanup failed: %w", err))
		}
	}

	// Clean up Testcontainer if used
	if c.testContainerDB != nil {
		if err := c.testContainerDB.Cleanup(); err != nil {
			errors = append(errors, fmt.Errorf("testcontainer cleanup failed: %w", err))
		}
	}

	// Clear test data manager
	if c.dataManager != nil {
		c.dataManager.Clear()
	}

	if len(errors) > 0 {
		return fmt.Errorf("teardown errors: %v", errors)
	}
	return nil
}

// theMCPServerIsRunning step implementation
func (c *CommonStepContext) theMCPServerIsRunning() error {
	// Server is already started in setupScenario
	// This step just verifies it's running
	if c.bddContext == nil {
		return fmt.Errorf("BDD context not initialized")
	}
	return nil
}

// theMCPConnectionIsInitialized step implementation
func (c *CommonStepContext) theMCPConnectionIsInitialized() error {
	// Connection is already initialized in setupScenario
	// This step just verifies it's connected
	return c.bddContext.WaitForServer(5 * time.Second)
}

// theDatabaseIsClean step implementation
func (c *CommonStepContext) theDatabaseIsClean() error {
	// Database is already cleaned in setupScenario
	// This step can verify the database is empty
	count, err := c.testDB.CountRows("movies", "")
	if err != nil {
		return fmt.Errorf("failed to count movies: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("database not clean: found %d movies", count)
	}

	count, err = c.testDB.CountRows("actors", "")
	if err != nil {
		return fmt.Errorf("failed to count actors: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("database not clean: found %d actors", count)
	}

	return nil
}

// theResponseShouldBeSuccessful step implementation
func (c *CommonStepContext) theResponseShouldBeSuccessful() error {
	if c.bddContext.HasError() {
		return fmt.Errorf("expected successful response, but got error: %s", c.bddContext.GetErrorMessage())
	}

	response := c.bddContext.GetLastResponse()
	if response == nil {
		return fmt.Errorf("no response received")
	}

	if response.IsError {
		return fmt.Errorf("response contains error: %v", response.Content)
	}

	return nil
}

// theResponseShouldContainAnError step implementation
func (c *CommonStepContext) theResponseShouldContainAnError() error {
	if !c.bddContext.HasError() {
		return fmt.Errorf("expected error response, but got successful response")
	}
	return nil
}

// theResponseShouldContainAnErrorWithMessage step implementation
func (c *CommonStepContext) theResponseShouldContainAnErrorWithMessage(expectedMessage string) error {
	if !c.bddContext.HasError() {
		return fmt.Errorf("expected error response, but got successful response")
	}

	actualMessage := c.bddContext.GetErrorMessage()
	if actualMessage != expectedMessage {
		return fmt.Errorf("expected error message '%s', but got '%s'", expectedMessage, actualMessage)
	}

	return nil
}