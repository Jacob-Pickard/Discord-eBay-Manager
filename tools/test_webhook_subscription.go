package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env for testing
	godotenv.Load("../.env")

	accessToken := os.Getenv("EBAY_ACCESS_TOKEN")
	verifyToken := os.Getenv("WEBHOOK_VERIFY_TOKEN")
	environment := os.Getenv("EBAY_ENVIRONMENT")

	if accessToken == "" {
		fmt.Println("Error: EBAY_ACCESS_TOKEN not set")
		return
	}

	baseURL := "https://api.sandbox.ebay.com"
	if environment == "PRODUCTION" {
		baseURL = "https://api.ebay.com"
	}

	url := baseURL + "/commerce/notification/v1/destination"

	// Subscribe to marketplace offer notifications
	topics := []string{
		"MARKETPLACE_OFFER",
	}

	webhookURL := "https://jacob.it.com/webhook/ebay/notification"

	// Build subscription payload
	payload := map[string]interface{}{
		"name":   "Discord_Bot_Notifications",
		"status": "ENABLED",
		"deliveryConfig": map[string]interface{}{
			"endpoint":    webhookURL,
			"verifyToken": verifyToken,
		},
		"topics": buildTopics(topics),
	}

	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal payload: %v\n", err)
		return
	}

	fmt.Println("=== Creating Subscription ===")
	fmt.Printf("URL: %s\n", url)
	fmt.Printf("Environment: %s\n", environment)
	fmt.Printf("Webhook URL: %s\n", webhookURL)
	fmt.Printf("Verify Token: %s\n", verifyToken)
	fmt.Printf("Verify Token Length: %d characters\n", len(verifyToken))
	fmt.Printf("\nPayload:\n%s\n\n", string(jsonData))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	// Dump request
	reqDump, _ := httputil.DumpRequestOut(req, true)
	fmt.Println("=== Request ===")
	fmt.Println(string(reqDump))
	fmt.Println()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Dump response
	respDump, _ := httputil.DumpResponse(resp, true)
	fmt.Println("=== Response ===")
	fmt.Println(string(respDump))
	fmt.Println()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		fmt.Printf("\n❌ Error: Status %d\n", resp.StatusCode)
		fmt.Printf("Response: %s\n", string(body))

		// Try to parse error
		var errorResp struct {
			Errors []struct {
				ErrorID  int    `json:"errorId"`
				Domain   string `json:"domain"`
				Category string `json:"category"`
				Message  string `json:"message"`
			} `json:"errors"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil {
			fmt.Println("\n=== Parsed Errors ===")
			for i, e := range errorResp.Errors {
				fmt.Printf("%d. Error ID: %d\n", i+1, e.ErrorID)
				fmt.Printf("   Domain: %s\n", e.Domain)
				fmt.Printf("   Category: %s\n", e.Category)
				fmt.Printf("   Message: %s\n", e.Message)
			}
		}
		return
	}

	fmt.Printf("\n✅ Success! Status: %d\n", resp.StatusCode)
	fmt.Printf("Response: %s\n", string(body))
}

func buildTopics(topicNames []string) []map[string]string {
	topics := make([]map[string]string, len(topicNames))
	for i, name := range topicNames {
		topics[i] = map[string]string{"topicName": name}
	}
	return topics
}
