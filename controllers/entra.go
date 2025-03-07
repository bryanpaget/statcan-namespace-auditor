package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// Configuration for Entra ID (replace these with your actual values or load them from a config map)
const (
	tenantID     = "your-tenant-id"
	clientID     = "your-client-id"
	clientSecret = "your-client-secret"
	graphAPIURL  = "https://graph.microsoft.com/v1.0"
)

// Graph API response structure for user lookup
type GraphAPIResponse struct {
	Value []struct {
		Mail string `json:"mail"`
	} `json:"value"`
}

// getEntraIDClient creates an OAuth2 client for authenticating with Microsoft Graph API
func getEntraIDClient() (*http.Client, error) {
	config := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID),
		Scopes:       []string{"https://graph.microsoft.com/.default"},
	}
	token, err := config.Token(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	return oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token)), nil
}

// checkEmailInEntraID checks if a given email exists in Entra ID
func checkEmailInEntraID(email string) (bool, error) {
	client, err := getEntraIDClient()
	if err != nil {
		return false, err
	}

	// Construct the Graph API URL to search for the user by email
	url := fmt.Sprintf("%s/users?$filter=mail eq '%s'", graphAPIURL, email)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode the response
	var graphResponse GraphAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&graphResponse); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check if any users were found
	return len(graphResponse.Value) > 0, nil
}

// isStatCanEmail checks if an email belongs to the StatCan domain
func isStatCanEmail(email string) bool {
	return strings.HasSuffix(email, "@statcan.gc.ca")
}
