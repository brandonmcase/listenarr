package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	baseURL = flag.String("url", "http://localhost:8686", "Base URL of the API")
	apiKey  = flag.String("apikey", "", "API key for authentication")
	verbose = flag.Bool("v", false, "Verbose output")
)

type TestResult struct {
	Name       string
	Method     string
	Endpoint   string
	Status     string
	StatusCode int
	Expected   int
	Message    string
	Error      error
}

func main() {
	flag.Parse()

	// Try to get API key from environment or config
	if *apiKey == "" {
		*apiKey = os.Getenv("LISTENARR_API_KEY")
	}

	results := []TestResult{}

	fmt.Println("==========================================")
	fmt.Println("Listenarr API Endpoint Verification")
	fmt.Println("==========================================")
	fmt.Printf("Base URL: %s\n", *baseURL)
	if *apiKey != "" {
		fmt.Printf("API Key: %s...\n", (*apiKey)[:min(8, len(*apiKey))])
	} else {
		fmt.Println("API Key: Not provided (some tests will be skipped)")
	}
	fmt.Println()

	// Test 1: Health Check (no auth required)
	fmt.Println("=== Public Endpoints ===")
	result := testEndpoint("GET", "/api/health", 200, "", nil)
	results = append(results, result)
	printResult(result)

	// Protected endpoints
	fmt.Println()
	fmt.Println("=== Protected Endpoints ===")

	if *apiKey == "" {
		fmt.Println("Skipping protected endpoints - no API key provided")
	} else {
		// Test 2: Get Library
		result = testEndpoint("GET", "/api/v1/library", 200, *apiKey, nil)
		results = append(results, result)
		printResult(result)

		// Test 3: Add to Library
		data := map[string]string{
			"title":  "Test Book",
			"author": "Test Author",
		}
		result = testEndpoint("POST", "/api/v1/library", 200, *apiKey, data)
		results = append(results, result)
		printResult(result)

		// Test 4: Get Downloads
		result = testEndpoint("GET", "/api/v1/downloads", 200, *apiKey, nil)
		results = append(results, result)
		printResult(result)

		// Test 5: Start Download
		downloadData := map[string]string{
			"libraryItemId": "test",
			"releaseId":     "test",
		}
		result = testEndpoint("POST", "/api/v1/downloads", 200, *apiKey, downloadData)
		results = append(results, result)
		printResult(result)

		// Test 6: Get Processing Queue
		result = testEndpoint("GET", "/api/v1/processing", 200, *apiKey, nil)
		results = append(results, result)
		printResult(result)

		// Test 7: Search Audiobooks
		result = testEndpoint("GET", "/api/v1/search?q=test", 200, *apiKey, nil)
		results = append(results, result)
		printResult(result)

		// Test 8: Remove from Library
		result = testEndpoint("DELETE", "/api/v1/library/test-id", 200, *apiKey, nil)
		results = append(results, result)
		printResult(result)

		// Test 9: Test authentication failure
		fmt.Println()
		fmt.Println("=== Authentication Tests ===")
		result = testEndpoint("GET", "/api/v1/library", 401, "invalid-key", nil)
		results = append(results, result)
		printResult(result)
	}

	// Test invalid endpoint
	fmt.Println()
	fmt.Println("=== Error Handling ===")
	result = testEndpoint("GET", "/api/v1/nonexistent", 404, *apiKey, nil)
	results = append(results, result)
	printResult(result)

	// Summary
	fmt.Println()
	fmt.Println("==========================================")
	fmt.Println("Test Summary")
	fmt.Println("==========================================")

	passed := 0
	failed := 0
	skipped := 0

	for _, r := range results {
		if r.Status == "PASS" {
			passed++
		} else if r.Status == "FAIL" {
			failed++
		} else {
			skipped++
		}
	}

	fmt.Printf("Passed: %d\n", passed)
	fmt.Printf("Failed: %d\n", failed)
	fmt.Printf("Skipped: %d\n", skipped)
	fmt.Println()

	if failed == 0 {
		fmt.Println("✅ All tests passed!")
		os.Exit(0)
	} else {
		fmt.Println("❌ Some tests failed!")
		os.Exit(1)
	}
}

func testEndpoint(method, endpoint string, expectedStatus int, authKey string, data interface{}) TestResult {
	url := *baseURL + endpoint
	client := &http.Client{Timeout: 5 * time.Second}

	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return TestResult{
				Name:       fmt.Sprintf("%s %s", method, endpoint),
				Method:     method,
				Endpoint:   endpoint,
				Status:     "FAIL",
				StatusCode: 0,
				Expected:   expectedStatus,
				Message:    "Failed to marshal request data",
				Error:      err,
			}
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return TestResult{
			Name:       fmt.Sprintf("%s %s", method, endpoint),
			Method:     method,
			Endpoint:   endpoint,
			Status:     "FAIL",
			StatusCode: 0,
			Expected:   expectedStatus,
			Message:    "Failed to create request",
			Error:      err,
		}
	}

	if authKey != "" {
		req.Header.Set("X-API-Key", authKey)
	}

	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return TestResult{
			Name:       fmt.Sprintf("%s %s", method, endpoint),
			Method:     method,
			Endpoint:   endpoint,
			Status:     "FAIL",
			StatusCode: 0,
			Expected:   expectedStatus,
			Message:    "Request failed",
			Error:      err,
		}
	}
	defer resp.Body.Close()

	status := "FAIL"
	message := fmt.Sprintf("Expected %d, got %d", expectedStatus, resp.StatusCode)

	if resp.StatusCode == expectedStatus {
		status = "PASS"
		message = "OK"
	}

	return TestResult{
		Name:       fmt.Sprintf("%s %s", method, endpoint),
		Method:     method,
		Endpoint:   endpoint,
		Status:     status,
		StatusCode: resp.StatusCode,
		Expected:   expectedStatus,
		Message:    message,
	}
}

func printResult(result TestResult) {
	var icon string
	var color string

	switch result.Status {
	case "PASS":
		icon = "✓"
		color = "\033[0;32m" // Green
	case "FAIL":
		icon = "✗"
		color = "\033[0;31m" // Red
	default:
		icon = "⊘"
		color = "\033[1;33m" // Yellow
	}

	fmt.Printf("%s%s%s %s - %s", color, icon, "\033[0m", result.Name, result.Message)
	if result.Error != nil && *verbose {
		fmt.Printf(" (Error: %v)", result.Error)
	}
	fmt.Println()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
