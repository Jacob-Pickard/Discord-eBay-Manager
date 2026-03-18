package ebay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Subscription represents an eBay notification subscription
type Subscription struct {
	DestinationID  string                 `json:"destinationId"`
	Name           string                 `json:"name"`
	Status         string                 `json:"status"`
	DeliveryConfig map[string]interface{} `json:"deliveryConfig"`
	Topics         []map[string]string    `json:"topics"`
}

// CreateWebhookSubscription subscribes to eBay notifications
func (c *Client) CreateWebhookSubscription(webhookURL string) error {
	if c.config.AccessToken == "" {
		return fmt.Errorf("no access token - run /ebay-authorize first")
	}

	url := c.baseURL + "/commerce/notification/v1/destination"

	// Subscribe to all important notification topics
	// Note: MARKETPLACE_ACCOUNT_DELETION removed - not required even with exemption
	topics := []string{
		"MARKETPLACE_OFFER", // Offer events (most important for your use case)
		"MARKETPLACE_ORDER", // All order events
		"ITEM_INVENTORY",    // Inventory changes
	}

	// Log verification token for debugging
	log.Printf("[DEBUG] Creating subscription with verify token: '%s' (length: %d)", c.config.WebhookVerifyToken, len(c.config.WebhookVerifyToken))

	// Build subscription payload
	payload := map[string]interface{}{
		"name":   "Discord_Bot_Notifications",
		"status": "ENABLED",
		"deliveryConfig": map[string]interface{}{
			"endpoint":    webhookURL,
			"verifyToken": c.config.WebhookVerifyToken,
		},
		"topics": buildTopics(topics),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	log.Printf("[DEBUG] Subscription payload: %s", string(jsonData))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		log.Printf("[DEBUG] eBay API error response: %s", string(body))
		log.Printf("[DEBUG] Response headers: %v", resp.Header)
		// Extract rlogid from response headers if present
		rlogid := resp.Header.Get("Rlogid")
		if rlogid == "" {
			rlogid = resp.Header.Get("X-Ebay-C-Request-Id")
		}
		if rlogid == "" {
			rlogid = resp.Header.Get("X-EBAY-C-REQUEST-ID")
		}
		if rlogid != "" {
			log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			log.Printf("rlogid: %s", rlogid)
			log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		}
		return fmt.Errorf("failed to create subscription (status %d): %s", resp.StatusCode, string(body))
	}

	log.Printf("[DEBUG] Subscription created successfully. Response: %s", string(body))
	return nil
}

// ListWebhookSubscriptions returns all active webhook subscriptions
func (c *Client) ListWebhookSubscriptions() ([]Subscription, error) {
	if c.config.AccessToken == "" {
		return nil, fmt.Errorf("no access token - run /ebay-authorize first")
	}

	url := c.baseURL + "/commerce/notification/v1/destination"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	log.Printf("[DEBUG] List subscriptions response status: %d", resp.StatusCode)
	log.Printf("[DEBUG] List subscriptions response body: %s", string(body))

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("failed to list subscriptions (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Destinations []Subscription `json:"destinations"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	log.Printf("[DEBUG] Parsed %d destinations from response", len(result.Destinations))

	return result.Destinations, nil
}

// DeleteWebhookSubscription removes a webhook subscription
func (c *Client) DeleteWebhookSubscription(destinationID string) error {
	if c.config.AccessToken == "" {
		return fmt.Errorf("no access token - run /ebay-authorize first")
	}

	url := c.baseURL + "/commerce/notification/v1/destination/" + destinationID

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete subscription (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// buildTopics creates the topics array for subscription
func buildTopics(topicNames []string) []map[string]string {
	topics := make([]map[string]string, len(topicNames))
	for i, name := range topicNames {
		topics[i] = map[string]string{"topicName": name}
	}
	return topics
}
