package step_definitions

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cucumber/godog"
)

// RegisterMCPProtocolSteps registers step definitions for MCP protocol testing
func RegisterMCPProtocolSteps(sc *godog.ScenarioContext, ctx *TestContext) {
	// Background steps
	sc.Step(`^the MCP server is running$`, ctx.theMCPServerIsRunning)
	sc.Step(`^I have a valid MCP client connection$`, ctx.iHaveAValidMCPClientConnection)
	sc.Step(`^the MCP connection is initialized$`, ctx.theMCPConnectionIsInitialized)
	
	// Protocol steps
	sc.Step(`^I send an initialize request with:$`, ctx.iSendAnInitializeRequestWith)
	sc.Step(`^I send a tools/list request$`, ctx.iSendAToolsListRequest)
	sc.Step(`^I send a resources/list request$`, ctx.iSendAResourcesListRequest)
	sc.Step(`^I send a request with invalid method "([^"]*)"$`, ctx.iSendARequestWithInvalidMethod)
	sc.Step(`^I send an initialize request with protocol version "([^"]*)"$`, ctx.iSendAnInitializeRequestWithProtocolVersion)
	
	// Response validation steps
	sc.Step(`^the response should be successful$`, ctx.theResponseShouldBeSuccessful)
	sc.Step(`^the response should contain server capabilities$`, ctx.theResponseShouldContainServerCapabilities)
	sc.Step(`^the protocol version should be "([^"]*)"$`, ctx.theProtocolVersionShouldBe)
	sc.Step(`^the response should contain the following tools:$`, ctx.theResponseShouldContainTheFollowingTools)
	sc.Step(`^the response should contain the following resources:$`, ctx.theResponseShouldContainTheFollowingResources)
	sc.Step(`^the response should contain an error$`, ctx.theResponseShouldContainAnError)
	sc.Step(`^the error code should be (-?\d+)$`, ctx.theErrorCodeShouldBe)
	sc.Step(`^the error message should contain "([^"]*)"$`, ctx.theErrorMessageShouldContain)
	sc.Step(`^the error should indicate unsupported protocol version$`, ctx.theErrorShouldIndicateUnsupportedProtocolVersion)
}

func (ctx *TestContext) theMCPServerIsRunning() error {
	// Start the MCP server if not already running
	if ctx.mcpServerCmd == nil {
		err := ctx.StartMCPServer()
		if err != nil {
			return fmt.Errorf("failed to start MCP server: %w", err)
		}
		
		// Wait for server to be ready
		
		err = ctx.WaitForServer()
		if err != nil {
			return fmt.Errorf("MCP server not ready: %w", err)
		}
	}
	return nil
}

func (ctx *TestContext) iHaveAValidMCPClientConnection() error {
	// This step assumes the connection is established
	// In a real implementation, this might set up client credentials
	return nil
}

func (ctx *TestContext) theMCPConnectionIsInitialized() error {
	// Send initialize request
	request := &MCPRequest{
		JSONRPC: "2.0",
		ID:      "init-1",
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "godog-test-client",
				"version": "1.0.0",
			},
		},
	}
	
	err := ctx.SendMCPRequest(request)
	if err != nil {
		return fmt.Errorf("failed to send initialize request: %w", err)
	}
	
	// Verify successful initialization
	response, err := ctx.GetLastMCPResponse()
	if err != nil {
		return fmt.Errorf("failed to parse initialize response: %w", err)
	}
	
	if response.Error != nil {
		return fmt.Errorf("initialize request failed: %s", response.Error.Message)
	}
	
	return nil
}

func (ctx *TestContext) iSendAnInitializeRequestWith(docString *godog.DocString) error {
	var requestData map[string]interface{}
	err := json.Unmarshal([]byte(docString.Content), &requestData)
	if err != nil {
		return fmt.Errorf("failed to parse request JSON: %w", err)
	}
	
	request := &MCPRequest{
		JSONRPC: "2.0",
		ID:      "init-test",
		Method:  requestData["method"].(string),
		Params:  requestData["params"],
	}
	
	return ctx.SendMCPRequest(request)
}

func (ctx *TestContext) iSendAToolsListRequest() error {
	request := &MCPRequest{
		JSONRPC: "2.0",
		ID:      "tools-list",
		Method:  "tools/list",
	}
	
	return ctx.SendMCPRequest(request)
}

func (ctx *TestContext) iSendAResourcesListRequest() error {
	request := &MCPRequest{
		JSONRPC: "2.0",
		ID:      "resources-list",
		Method:  "resources/list",
	}
	
	return ctx.SendMCPRequest(request)
}

func (ctx *TestContext) iSendARequestWithInvalidMethod(method string) error {
	request := &MCPRequest{
		JSONRPC: "2.0",
		ID:      "invalid-method",
		Method:  method,
	}
	
	return ctx.SendMCPRequest(request)
}

func (ctx *TestContext) iSendAnInitializeRequestWithProtocolVersion(version string) error {
	request := &MCPRequest{
		JSONRPC: "2.0",
		ID:      "init-version-test",
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": version,
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "test-client",
				"version": "1.0.0",
			},
		},
	}
	
	return ctx.SendMCPRequest(request)
}

func (ctx *TestContext) theResponseShouldBeSuccessful() error {
	response, err := ctx.GetLastMCPResponse()
	if err != nil {
		return fmt.Errorf("failed to get response: %w", err)
	}
	
	if response.Error != nil {
		return fmt.Errorf("response contains error: %s", response.Error.Message)
	}
	
	if response.Result == nil {
		return fmt.Errorf("response missing result field")
	}
	
	return nil
}

func (ctx *TestContext) theResponseShouldContainServerCapabilities() error {
	response, err := ctx.GetLastMCPResponse()
	if err != nil {
		return err
	}
	
	result, ok := response.Result.(map[string]interface{})
	if !ok {
		return fmt.Errorf("response result is not an object")
	}
	
	_, hasCapabilities := result["capabilities"]
	if !hasCapabilities {
		return fmt.Errorf("response missing capabilities field")
	}
	
	return nil
}

func (ctx *TestContext) theProtocolVersionShouldBe(expectedVersion string) error {
	response, err := ctx.GetLastMCPResponse()
	if err != nil {
		return err
	}
	
	result, ok := response.Result.(map[string]interface{})
	if !ok {
		return fmt.Errorf("response result is not an object")
	}
	
	version, hasVersion := result["protocolVersion"]
	if !hasVersion {
		return fmt.Errorf("response missing protocolVersion field")
	}
	
	if version != expectedVersion {
		return fmt.Errorf("expected protocol version %s, got %v", expectedVersion, version)
	}
	
	return nil
}

func (ctx *TestContext) theResponseShouldContainTheFollowingTools(table *godog.Table) error {
	response, err := ctx.GetLastMCPResponse()
	if err != nil {
		return err
	}
	
	result, ok := response.Result.(map[string]interface{})
	if !ok {
		return fmt.Errorf("response result is not an object")
	}
	
	toolsInterface, hasTools := result["tools"]
	if !hasTools {
		return fmt.Errorf("response missing tools field")
	}
	
	toolsArray, ok := toolsInterface.([]interface{})
	if !ok {
		return fmt.Errorf("tools field is not an array")
	}
	
	// Convert tools to map for easier lookup
	toolsMap := make(map[string]string)
	for _, toolInterface := range toolsArray {
		tool, ok := toolInterface.(map[string]interface{})
		if !ok {
			continue
		}
		if name, hasName := tool["name"].(string); hasName {
			if desc, hasDesc := tool["description"].(string); hasDesc {
				toolsMap[name] = desc
			}
		}
	}
	
	// Check each expected tool
	for i := 1; i < len(table.Rows); i++ { // Skip header row
		row := table.Rows[i]
		if len(row.Cells) < 2 {
			continue
		}
		
		expectedName := row.Cells[0].Value
		expectedDesc := row.Cells[1].Value
		
		actualDesc, found := toolsMap[expectedName]
		if !found {
			return fmt.Errorf("tool %s not found in response", expectedName)
		}
		
		if !strings.Contains(actualDesc, expectedDesc) {
			return fmt.Errorf("tool %s description mismatch: expected %s, got %s", expectedName, expectedDesc, actualDesc)
		}
	}
	
	return nil
}

func (ctx *TestContext) theResponseShouldContainTheFollowingResources(table *godog.Table) error {
	response, err := ctx.GetLastMCPResponse()
	if err != nil {
		return err
	}
	
	result, ok := response.Result.(map[string]interface{})
	if !ok {
		return fmt.Errorf("response result is not an object")
	}
	
	resourcesInterface, hasResources := result["resources"]
	if !hasResources {
		return fmt.Errorf("response missing resources field")
	}
	
	resourcesArray, ok := resourcesInterface.([]interface{})
	if !ok {
		return fmt.Errorf("resources field is not an array")
	}
	
	// Convert resources to map for easier lookup
	resourcesMap := make(map[string]map[string]string)
	for _, resourceInterface := range resourcesArray {
		resource, ok := resourceInterface.(map[string]interface{})
		if !ok {
			continue
		}
		if uri, hasURI := resource["uri"].(string); hasURI {
			resourceInfo := make(map[string]string)
			if name, hasName := resource["name"].(string); hasName {
				resourceInfo["name"] = name
			}
			if desc, hasDesc := resource["description"].(string); hasDesc {
				resourceInfo["description"] = desc
			}
			resourcesMap[uri] = resourceInfo
		}
	}
	
	// Check each expected resource
	for i := 1; i < len(table.Rows); i++ { // Skip header row
		row := table.Rows[i]
		if len(row.Cells) < 3 {
			continue
		}
		
		expectedURI := row.Cells[0].Value
		expectedName := row.Cells[1].Value
		expectedDesc := row.Cells[2].Value
		
		resource, found := resourcesMap[expectedURI]
		if !found {
			return fmt.Errorf("resource %s not found in response", expectedURI)
		}
		
		if actualName, hasName := resource["name"]; hasName {
			if actualName != expectedName {
				return fmt.Errorf("resource %s name mismatch: expected %s, got %s", expectedURI, expectedName, actualName)
			}
		}
		
		if actualDesc, hasDesc := resource["description"]; hasDesc {
			if !strings.Contains(actualDesc, expectedDesc) {
				return fmt.Errorf("resource %s description mismatch: expected %s, got %s", expectedURI, expectedDesc, actualDesc)
			}
		}
	}
	
	return nil
}

func (ctx *TestContext) theResponseShouldContainAnError() error {
	response, err := ctx.GetLastMCPResponse()
	if err != nil {
		return err
	}
	
	if response.Error == nil {
		return fmt.Errorf("response should contain an error but doesn't")
	}
	
	return nil
}

func (ctx *TestContext) theErrorCodeShouldBe(expectedCode int) error {
	response, err := ctx.GetLastMCPResponse()
	if err != nil {
		return err
	}
	
	if response.Error == nil {
		return fmt.Errorf("response should contain an error")
	}
	
	if response.Error.Code != expectedCode {
		return fmt.Errorf("expected error code %d, got %d", expectedCode, response.Error.Code)
	}
	
	return nil
}

func (ctx *TestContext) theErrorMessageShouldContain(expectedMessage string) error {
	response, err := ctx.GetLastMCPResponse()
	if err != nil {
		return err
	}
	
	if response.Error == nil {
		return fmt.Errorf("response should contain an error")
	}
	
	if !strings.Contains(response.Error.Message, expectedMessage) {
		return fmt.Errorf("error message should contain %s, got: %s", expectedMessage, response.Error.Message)
	}
	
	return nil
}

func (ctx *TestContext) theErrorShouldIndicateUnsupportedProtocolVersion() error {
	response, err := ctx.GetLastMCPResponse()
	if err != nil {
		return err
	}
	
	if response.Error == nil {
		return fmt.Errorf("response should contain an error")
	}
	
	message := strings.ToLower(response.Error.Message)
	if !strings.Contains(message, "protocol") && !strings.Contains(message, "version") {
		return fmt.Errorf("error message should indicate protocol version issue, got: %s", response.Error.Message)
	}
	
	return nil
}