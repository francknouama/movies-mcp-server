package models

import (
	"time"

	messages "github.com/cucumber/messages/go/v21"
)

// Feature represents a Gherkin feature file with additional metadata
type Feature struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	FilePath    string     `json:"file_path"`
	Language    string     `json:"language"`
	Tags        []string   `json:"tags"`
	Scenarios   []Scenario `json:"scenarios"`
	Background  *Step      `json:"background,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Scenario represents a Gherkin scenario
type Scenario struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Steps       []Step    `json:"steps"`
	Examples    []Example `json:"examples,omitempty"`
	Location    Location  `json:"location"`
}

// Step represents a Gherkin step
type Step struct {
	ID       string   `json:"id"`
	Keyword  string   `json:"keyword"` // Given, When, Then, And, But
	Text     string   `json:"text"`
	Argument string   `json:"argument,omitempty"` // DocString or DataTable
	Location Location `json:"location"`
}

// Example represents scenario examples (for scenario outlines)
type Example struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Tags     []string   `json:"tags"`
	Headers  []string   `json:"headers"`
	Rows     [][]string `json:"rows"`
	Location Location   `json:"location"`
}

// Location represents the file location of a Gherkin element
type Location struct {
	Line   uint32 `json:"line"`
	Column uint32 `json:"column,omitempty"`
}

// TestResult represents the result of running Godog tests
type TestResult struct {
	ID        string          `json:"id"`
	SuiteName string          `json:"suite_name"`
	StartTime time.Time       `json:"start_time"`
	EndTime   time.Time       `json:"end_time"`
	Duration  time.Duration   `json:"duration"`
	Status    TestStatus      `json:"status"`
	Features  []FeatureResult `json:"features"`
	Summary   TestSummary     `json:"summary"`
	Metadata  TestMetadata    `json:"metadata"`
}

// FeatureResult represents the result of running a single feature
type FeatureResult struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	FilePath  string           `json:"file_path"`
	Status    TestStatus       `json:"status"`
	Duration  time.Duration    `json:"duration"`
	Scenarios []ScenarioResult `json:"scenarios"`
	Error     string           `json:"error,omitempty"`
}

// ScenarioResult represents the result of running a single scenario
type ScenarioResult struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Status   TestStatus    `json:"status"`
	Duration time.Duration `json:"duration"`
	Steps    []StepResult  `json:"steps"`
	Error    string        `json:"error,omitempty"`
}

// StepResult represents the result of running a single step
type StepResult struct {
	ID       string        `json:"id"`
	Keyword  string        `json:"keyword"`
	Text     string        `json:"text"`
	Status   TestStatus    `json:"status"`
	Duration time.Duration `json:"duration"`
	Error    string        `json:"error,omitempty"`
	Output   string        `json:"output,omitempty"`
}

// TestSummary provides summary statistics for test results
type TestSummary struct {
	TotalFeatures   int `json:"total_features"`
	PassedFeatures  int `json:"passed_features"`
	FailedFeatures  int `json:"failed_features"`
	SkippedFeatures int `json:"skipped_features"`

	TotalScenarios   int `json:"total_scenarios"`
	PassedScenarios  int `json:"passed_scenarios"`
	FailedScenarios  int `json:"failed_scenarios"`
	SkippedScenarios int `json:"skipped_scenarios"`

	TotalSteps     int `json:"total_steps"`
	PassedSteps    int `json:"passed_steps"`
	FailedSteps    int `json:"failed_steps"`
	SkippedSteps   int `json:"skipped_steps"`
	UndefinedSteps int `json:"undefined_steps"`
	PendingSteps   int `json:"pending_steps"`
}

// TestMetadata contains additional information about test execution
type TestMetadata struct {
	GodogVersion string            `json:"godog_version"`
	Platform     string            `json:"platform"`
	GoVersion    string            `json:"go_version"`
	Arguments    []string          `json:"arguments"`
	Environment  map[string]string `json:"environment"`
}

// TestStatus represents the status of test execution
type TestStatus string

const (
	StatusPassed    TestStatus = "passed"
	StatusFailed    TestStatus = "failed"
	StatusSkipped   TestStatus = "skipped"
	StatusUndefined TestStatus = "undefined"
	StatusPending   TestStatus = "pending"
)

// StepDefinition represents a Go step definition
type StepDefinition struct {
	ID          string   `json:"id"`
	Expression  string   `json:"expression"`  // Regex pattern
	Handler     string   `json:"handler"`     // Function name
	File        string   `json:"file"`        // Source file
	Line        int      `json:"line"`        // Line number
	Description string   `json:"description"` // Optional description
	Tags        []string `json:"tags"`        // Optional tags
}

// GodogRunOptions represents options for running Godog tests
type GodogRunOptions struct {
	FeaturePaths  []string          `json:"feature_paths"`
	Tags          string            `json:"tags,omitempty"`
	Format        string            `json:"format,omitempty"`
	Output        string            `json:"output,omitempty"`
	Strict        bool              `json:"strict"`
	NoColors      bool              `json:"no_colors"`
	StopOnFailure bool              `json:"stop_on_failure"`
	Randomize     int64             `json:"randomize,omitempty"`
	Concurrency   int               `json:"concurrency,omitempty"`
	Environment   map[string]string `json:"environment,omitempty"`
}

// ConvertGherkinFeature converts a gherkin feature to our Feature model
func ConvertGherkinFeature(gf *messages.Feature, filePath string) *Feature {
	feature := &Feature{
		ID:          generateID(),
		Name:        gf.Name,
		Description: gf.Description,
		FilePath:    filePath,
		Language:    gf.Language,
		Tags:        extractTags(gf.Tags),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Convert scenarios
	for _, scenario := range gf.Children {
		if scenario.Scenario != nil {
			feature.Scenarios = append(feature.Scenarios, convertGherkinScenario(scenario.Scenario))
		}
	}

	return feature
}

// Helper functions
func generateID() string {
	return time.Now().Format("20060102150405") + "-" + generateRandomString(8)
}

func generateRandomString(length int) string {
	// Simple implementation - in production, use crypto/rand
	return "abcd1234"
}

func extractTags(tags []*messages.Tag) []string {
	var result []string
	for _, tag := range tags {
		result = append(result, tag.Name)
	}
	return result
}

func convertGherkinScenario(gs *messages.Scenario) Scenario {
	scenario := Scenario{
		ID:          generateID(),
		Name:        gs.Name,
		Description: gs.Description,
		Tags:        extractTags(gs.Tags),
		Location: Location{
			Line:   uint32(gs.Location.Line),
			Column: uint32(gs.Location.Column),
		},
	}

	// Convert steps
	for _, step := range gs.Steps {
		scenario.Steps = append(scenario.Steps, convertGherkinStep(step))
	}

	return scenario
}

func convertGherkinStep(gstep *messages.Step) Step {
	return Step{
		ID:      generateID(),
		Keyword: gstep.Keyword,
		Text:    gstep.Text,
		Location: Location{
			Line:   uint32(gstep.Location.Line),
			Column: uint32(gstep.Location.Column),
		},
	}
}
