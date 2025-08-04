package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func sample_main() {
	// Configuration
	baseURL := "http://localhost:8080/api/v1"
	apiKey := "your-api-key-here"

	// 1. Create scholarship scheme
	schemePayload := `{
		"scheme_id": "summer-scholarship-2023",
		"title": "Summer Scholarship Program",
		"description": "Financial aid for underprivileged students",
		"eligibility_criteria": {
			"age": {"operator": "lt", "value": 25},
			"income": {"operator": "lt", "value": 300000}
		},
		"scheme_amount": 50000
	}`

	// Push the scheme to ONEST
	schemeResp := makeAPIRequest(baseURL, "POST", "/organizations/1/schemes", schemePayload, apiKey)
	schemeID := schemeResp["data"].(map[string]interface{})["scheme_id"].(string)

	// 2. Get applications for the scheme
	applicationsResp := makeAPIRequest(baseURL, "GET", "/schemes/"+schemeID+"/applications", "", apiKey)

	// 3. Update application status (if any applications exist)
	if apps, ok := applicationsResp["data"].([]interface{}); ok && len(apps) > 0 {
		appID := fmt.Sprintf("%.0f", apps[0].(map[string]interface{})["id"].(float64))
		makeAPIRequest(baseURL, "PUT", "/applications/"+appID+"/status", `{"status":"APPROVED"}`, apiKey)
	}
}

func makeAPIRequest(baseURL, method, endpoint, payload, apiKey string) map[string]interface{} {
	req, _ := http.NewRequest(method, baseURL+endpoint, bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result
}
