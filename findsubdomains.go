package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// fetchJSONWithHeader performs an HTTP GET request with a custom header and decodes the JSON response.
func fetchJSONWithHeader(url string, headerName string, headerValue string, wrapper interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set(headerName, headerValue)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	dec := json.NewDecoder(resp.Body)
	return dec.Decode(wrapper)
}

func fetchSpyseV4Subdomains(domain string, apiKey string) ([]string, error) {
	out := make([]string, 0)
	// According to Spyse API v4 docs, limit can be up to 1000.
	// Free tier might have lower limits or only return the first page.
	// For simplicity, we'll fetch one page with a reasonable limit.
	// Pagination would require handling the 'offset' and 'total_count'.
	limit := 100 
	fetchURL := fmt.Sprintf(
		"https://api.spyse.com/v4/data/domain/subdomain?domain_name=%s&limit=%d",
		domain, limit,
	)

	type SpyseSubdomainItem struct {
		Domain string `json:"domain"`
	}
	type SpyseData struct {
		Items      []SpyseSubdomainItem `json:"items"`
		TotalCount int                  `json:"total_count"`
		Offset     int                  `json:"offset"`
		SearchID   string               `json:"search_id"`
	}
	type SpyseResponse struct {
		Data   SpyseData `json:"data"`
		Status string    `json:"status"`
	}

	var wrapper SpyseResponse
	err := fetchJSONWithHeader(fetchURL, "X-API-KEY", apiKey, &wrapper)
	if err != nil {
		return nil, fmt.Errorf("Spyse API V4 request failed: %w", err)
	}

	if wrapper.Status != "ok" && wrapper.Status != "" { // API might return empty status on success
		return nil, fmt.Errorf("Spyse API V4 returned status: %s", wrapper.Status)
	}
	
	for _, item := range wrapper.Data.Items {
		out = append(out, item.Domain)
	}

	return out, nil
}

func fetchFindSubDomains(domain string) ([]string, error) {
	apiKey := os.Getenv("SPYSE_API_TOKEN")
	if apiKey == "" {
		// Must have an API token
		return []string{}, nil
	}

	subdomains, err := fetchSpyseV4Subdomains(domain, apiKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching subdomains from Spyse V4 for domain %s: %v\n", domain, err)
		return []string{}, nil // Return empty list on error to allow other sources to proceed
	}

	return subdomains, nil
}
