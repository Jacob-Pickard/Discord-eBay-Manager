package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SubscriptionManager handles eBay notification subscription management
type SubscriptionManager struct {
	accessToken string
	environment string
}

// NewSubscriptionManager creates a new subscription manager
func NewSubscriptionManager(accessToken, environment string) *SubscriptionManager {
	return &SubscriptionManager{
		accessToken: accessToken,
		environment: environment,
	}
}

// GetBaseURL returns the appropriate eBay API URL
func (sm *SubscriptionManager) GetBaseURL() string {
	if sm.environment == "PRODUCTION" {
		return "https://api.ebay.com"
	}
	return "https://api.sandbox.ebay.com"
}

// CreateSubscription creates a new webhook subscription with eBay
func (sm *SubscriptionManager) CreateSubscription(webhookURL, verifyToken string, topics []string) error {
	url := sm.GetBaseURL() + "/commerce/notification/v1/destination"

	// Build subscription payload
	payload := map[string]interface{}{
		"name":   "Discord_Bot_Notifications",
		"status": "ENABLED",
		"deliveryConfig": map[string]interface{}{
			"endpoint":     webhookURL,
			"verifyToken":  verifyToken,
		},
		"topics": sm.buildTopics(topics),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+sm.accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	
	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to create subscription (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// ListSubscriptions retrieves all active subscriptions
func (sm *SubscriptionManager) ListSubscriptions() ([]Subscription, error) {
	url := sm.GetBaseURL() + "/commerce/notification/v1/destination"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+sm.accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("failed to list subscriptions (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Destinations []Subscription `json:"destinations"`
	}
	
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Destinations, nil
}

// DeleteSubscription removes a webhook subscription
func (sm *SubscriptionManager) DeleteSubscription(destinationID string) error {
	url := sm.GetBaseURL() + "/commerce/notification/v1/destination/" + destinationID

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+sm.accessToken)

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
func (sm *SubscriptionManager) buildTopics(topicNames []string) []map[string]string {
	topics := make([]map[string]string, len(topicNames))
	for i, name := range topicNames {
		topics[i] = map[string]string{"topicName": name}
	}
	return topics
}

// Subscription represents an eBay notification subscription
type Subscription struct {
	DestinationID  string                 `json:"destinationId"`
	Name           string                 `json:"name"`
	Status         string                 `json:"status"`
	DeliveryConfig map[string]interface{} `json:"deliveryConfig"`
	Topics         []map[string]string    `json:"topics"`
}
