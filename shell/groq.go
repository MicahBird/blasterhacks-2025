package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// --- Structs for Request Body ---

type GroqRequest struct {
	Model          string         `json:"model"`
	Messages       []Message      `json:"messages"`
	ResponseFormat ResponseFormat `json:"response_format"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

// --- Structs for Response Body ---

// Main structure for the Groq API response
type GroqResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
	// Add other fields like system_fingerprint, x_groq if needed
}

type Choice struct {
	Index        int             `json:"index"`
	Message      ResponseMessage `json:"message"`
	FinishReason string          `json:"finish_reason"`
	// Add logprobs if needed
}

// Message structure within the response choice
type ResponseMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"` // This will contain the JSON string like {"category": "..."}
}

type Usage struct {
	QueueTime        float64 `json:"queue_time"`
	PromptTokens     int     `json:"prompt_tokens"`
	PromptTime       float64 `json:"prompt_time"`
	CompletionTokens int     `json:"completion_tokens"`
	CompletionTime   float64 `json:"completion_time"`
	TotalTokens      int     `json:"total_tokens"`
	TotalTime        float64 `json:"total_time"`
}

// --- Struct for the Inner JSON content ---

type CategoryResult struct {
	Category string `json:"category"`
}

// getGroqCatagory sends a command to the Groq API and extracts the category.
func getGroqCatagory(command string) string {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		log.Println("Error: GROQ_API_KEY environment variable not set.")
		return "" // Or handle error appropriately
	}

	groqURL := "https://api.groq.com/openai/v1/chat/completions"

	// Define the system message content
	systemContent := `You are a machine that will translate a command given by a user to category in a JSON output. The only possible categories are "pythondev", "javadev", "jsdev", "cppdev", "rustdev", "csharpdev", "phpdev", "golangdev", "swiftdev", "kotlindev", "sysdev", "webdev", and "securitydev". The output must be a JSON object with a single key "category". For example, if a user requests "gcc", you should output: {"category": "cppdev"}`

	// Construct the request body payload
	payload := GroqRequest{
		Model: "llama-3.3-70b-versatile", // Using the model specified in curl example - Update if needed
		Messages: []Message{
			{
				Role:    "system",
				Content: systemContent,
			},
			{
				Role:    "user",
				Content: command, // Use the input command here
			},
		},
		ResponseFormat: ResponseFormat{
			Type: "json_object",
		},
	}

	// Marshal the payload into JSON
	requestBodyBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling request payload: %v\n", err)
		return ""
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", groqURL, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		log.Printf("Error creating HTTP request: %v\n", err)
		return ""
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Create an HTTP client and execute the request
	client := &http.Client{Timeout: 30 * time.Second} // Add a timeout
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request to Groq API: %v\n", err)
		return ""
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Read the response body
	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		return ""
	}

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("Groq API returned non-OK status: %s\nResponse: %s\n", resp.Status, string(responseBodyBytes))
		return ""
	}

	// Unmarshal the main Groq API response
	var groqResponse GroqResponse
	err = json.Unmarshal(responseBodyBytes, &groqResponse)
	if err != nil {
		log.Printf("Error unmarshaling Groq response: %v\nResponse Body: %s\n", err, string(responseBodyBytes))
		return ""
	}

	// Basic validation: Check if choices exist
	if len(groqResponse.Choices) == 0 {
		log.Println("Error: Groq response contained no choices.")
		return ""
	}

	// Extract the content string which should be JSON
	contentJSON := groqResponse.Choices[0].Message.Content

	// Unmarshal the inner JSON content to get the category
	var categoryResult CategoryResult
	err = json.Unmarshal([]byte(contentJSON), &categoryResult)
	if err != nil {
		// Log the error and the content string that failed to parse
		log.Printf("Error unmarshaling inner category JSON: %v\nInner JSON String: %s\n", err, contentJSON)
		return ""
	}

	// Return the extracted category
	return categoryResult.Category
}
