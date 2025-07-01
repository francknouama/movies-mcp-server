package godog

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/francknouama/movies-mcp-server/shared-mcp/pkg/errors"
	"github.com/francknouama/movies-mcp-server/shared-mcp/pkg/logging"

	gherkin "github.com/cucumber/gherkin/go/v26"
	messages "github.com/cucumber/messages/go/v21"
	"github.com/francknouama/movies-mcp-server/godog-server/internal/config"
	"github.com/francknouama/movies-mcp-server/godog-server/internal/models"
)

// Runner handles Godog test execution and feature management
type Runner struct {
	config *config.Config
	logger *logging.Logger
}

// NewRunner creates a new Godog runner instance
func NewRunner(cfg *config.Config, logger *logging.Logger) *Runner {
	return &Runner{
		config: cfg,
		logger: logger,
	}
}

// CheckAvailability verifies that Godog is available and working
func (r *Runner) CheckAvailability() error {
	cmd := exec.Command(r.config.GodogBinary, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.NewGodogNotFound("Godog binary not found or not executable: " + err.Error())
	}

	r.logger.WithField("version", strings.TrimSpace(string(output))).Info("Godog version detected")
	return nil
}

// ValidateFeature parses and validates a Gherkin feature file
func (r *Runner) ValidateFeature(filePath string) (any, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, errors.NewFeatureParseError("Feature file not found: " + filePath)
	}

	// Read the feature file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.NewFeatureParseError("Failed to read feature file: " + err.Error())
	}

	// Parse with Gherkin
	gherkinDoc, err := gherkin.ParseGherkinDocument(strings.NewReader(string(content)), (&messages.Incrementing{}).NewId)
	if err != nil {
		return nil, errors.NewFeatureParseError("Gherkin parsing failed: " + err.Error())
	}

	if gherkinDoc.Feature == nil {
		return nil, errors.NewFeatureParseError("No feature found in file: " + filePath)
	}

	// Convert to our model
	feature := models.ConvertGherkinFeature(gherkinDoc.Feature, filePath)

	return map[string]any{
		"valid":   true,
		"feature": feature,
		"message": "Feature file is valid",
	}, nil
}

// ListFeatures discovers and lists all feature files
func (r *Runner) ListFeatures(directory string, includeContent bool) (any, error) {
	searchDir := r.config.FeaturesDir
	if directory != "" {
		searchDir = directory
	}

	var features []models.Feature
	var featureList []map[string]any

	err := filepath.Walk(searchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".feature") {
			if includeContent {
				// Parse each feature file with full content
				result, parseErr := r.ValidateFeature(path)
				if parseErr != nil {
					r.logger.WithField("file", path).WithField("error", parseErr).Warn("Failed to parse feature file")
					return nil // Continue walking, don't fail the entire operation
				}

				if resultMap, ok := result.(map[string]any); ok {
					if feature, ok := resultMap["feature"].(models.Feature); ok {
						features = append(features, feature)
					}
				}
			} else {
				// Just collect basic file information
				relPath, _ := filepath.Rel(searchDir, path)
				featureInfo := map[string]any{
					"file_path":     path,
					"relative_path": relPath,
					"name":          strings.TrimSuffix(filepath.Base(path), ".feature"),
					"size":          info.Size(),
					"modified_at":   info.ModTime(),
				}
				featureList = append(featureList, featureInfo)
			}
		}

		return nil
	})

	if err != nil {
		return nil, errors.NewFeatureParseError("Failed to scan features directory: " + err.Error())
	}

	if includeContent {
		return map[string]any{
			"features":     features,
			"count":        len(features),
			"directory":    searchDir,
			"with_content": true,
		}, nil
	} else {
		return map[string]any{
			"feature_files": featureList,
			"count":         len(featureList),
			"directory":     searchDir,
			"with_content":  false,
		}, nil
	}
}

// GetFeatureContent retrieves the content of a specific feature file
func (r *Runner) GetFeatureContent(filePath string, includeParsed bool) (any, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, errors.NewFeatureParseError("Feature file not found: " + filePath)
	}

	// Read the raw content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.NewFeatureParseError("Failed to read feature file: " + err.Error())
	}

	result := map[string]any{
		"file_path":   filePath,
		"raw_content": string(content),
		"size":        len(content),
	}

	// Add file metadata
	if info, err := os.Stat(filePath); err == nil {
		result["modified_at"] = info.ModTime()
		result["name"] = strings.TrimSuffix(filepath.Base(filePath), ".feature")
	}

	if includeParsed {
		// Parse the Gherkin content
		validateResult, parseErr := r.ValidateFeature(filePath)
		if parseErr != nil {
			result["parse_error"] = parseErr.Error()
			result["valid"] = false
		} else {
			if resultMap, ok := validateResult.(map[string]any); ok {
				result["parsed_feature"] = resultMap["feature"]
				result["valid"] = resultMap["valid"]
				result["parse_message"] = resultMap["message"]
			}
		}
	}

	return result, nil
}

// RunFeature executes a specific feature file
func (r *Runner) RunFeature(filePath string, options *models.GodogRunOptions) (*models.TestResult, error) {
	if options == nil {
		options = &models.GodogRunOptions{
			FeaturePaths: []string{filePath},
			Format:       "cucumber",
			Strict:       true,
		}
	}

	return r.runGodog(options)
}

// RunSuite executes the complete test suite
func (r *Runner) RunSuite(options *models.GodogRunOptions) (*models.TestResult, error) {
	if options == nil {
		options = &models.GodogRunOptions{
			FeaturePaths: []string{r.config.FeaturesDir},
			Format:       "cucumber",
			Strict:       true,
		}
	}

	return r.runGodog(options)
}

// GetLatestReport retrieves the most recent test results
func (r *Runner) GetLatestReport() (any, error) {
	// Check if reports directory exists
	if _, err := os.Stat(r.config.ReportsDir); os.IsNotExist(err) {
		return map[string]any{
			"message":     "No reports found",
			"reports_dir": r.config.ReportsDir,
		}, nil
	}

	// Find the most recent report file
	files, err := filepath.Glob(filepath.Join(r.config.ReportsDir, "*.json"))
	if err != nil {
		return nil, errors.NewReportGenerationError("Failed to scan reports directory: " + err.Error())
	}

	if len(files) == 0 {
		return map[string]any{
			"message":     "No JSON reports found",
			"reports_dir": r.config.ReportsDir,
		}, nil
	}

	// Find the most recent file
	var latestFile string
	var latestTime time.Time

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		if info.ModTime().After(latestTime) {
			latestTime = info.ModTime()
			latestFile = file
		}
	}

	if latestFile == "" {
		return map[string]any{
			"message": "No accessible reports found",
		}, nil
	}

	// Read and return the latest report
	content, err := os.ReadFile(latestFile)
	if err != nil {
		return nil, errors.NewReportGenerationError("Failed to read report file: " + err.Error())
	}

	var report any
	if err := json.Unmarshal(content, &report); err != nil {
		return nil, errors.NewReportGenerationError("Failed to parse report JSON: " + err.Error())
	}

	return map[string]any{
		"report":      report,
		"file":        latestFile,
		"modified_at": latestTime,
	}, nil
}

// runGodog executes Godog with the specified options
func (r *Runner) runGodog(options *models.GodogRunOptions) (*models.TestResult, error) {
	startTime := time.Now()

	// Prepare command arguments
	args := []string{}

	// Add feature paths
	args = append(args, options.FeaturePaths...)

	// Add format
	if options.Format != "" {
		args = append(args, "--format", options.Format)
	}

	// Add output file for JSON reports
	reportFile := filepath.Join(r.config.ReportsDir, fmt.Sprintf("report_%s.json", startTime.Format("20060102_150405")))
	if options.Format == "cucumber" || options.Format == "" {
		// Ensure reports directory exists
		os.MkdirAll(r.config.ReportsDir, 0755)
		args = append(args, "--output", reportFile)
	}

	// Add tags
	if options.Tags != "" {
		args = append(args, "--tags", options.Tags)
	}

	// Add other options
	if options.Strict {
		args = append(args, "--strict")
	}

	if options.NoColors {
		args = append(args, "--no-colors")
	}

	if options.StopOnFailure {
		args = append(args, "--stop-on-failure")
	}

	if options.Concurrency > 0 {
		args = append(args, "--concurrency", fmt.Sprintf("%d", options.Concurrency))
	}

	// Execute Godog
	cmd := exec.Command(r.config.GodogBinary, args...)

	// Set environment variables
	if options.Environment != nil {
		env := os.Environ()
		for key, value := range options.Environment {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
		cmd.Env = env
	}

	output, err := cmd.CombinedOutput()
	endTime := time.Now()

	// Create test result
	testResult := &models.TestResult{
		ID:        generateTestID(),
		SuiteName: "Godog Test Suite",
		StartTime: startTime,
		EndTime:   endTime,
		Duration:  endTime.Sub(startTime),
		Status:    models.StatusPassed,
		Metadata: models.TestMetadata{
			Arguments:   args,
			Environment: options.Environment,
		},
	}

	// Determine status based on exit code
	if err != nil {
		testResult.Status = models.StatusFailed
		r.logger.WithField("error", err).WithField("output", string(output)).Error("Godog execution failed")
	}

	// Try to parse JSON report if it exists
	if options.Format == "cucumber" || options.Format == "" {
		if _, readErr := os.ReadFile(reportFile); readErr == nil {
			// TODO: Parse Cucumber JSON format and populate FeatureResults
			r.logger.WithField("report_file", reportFile).Info("Generated JSON report")
		}
	}

	return testResult, nil
}

// generateTestID creates a unique test run identifier
func generateTestID() string {
	return fmt.Sprintf("test_%s_%d", time.Now().Format("20060102_150405"), time.Now().UnixNano()%1000)
}
