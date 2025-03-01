package api

import (
	"bytes"
	"citizenship-tracker-cli/pkg/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Auth(uci, password string) (*model.AuthResponse, error) {
	headers := map[string]string{
		"Content-Type":     "application/x-amz-json-1.1",
		"Accept":           "*/*",
		"Sec-Fetch-Site":   "cross-site",
		"Accept-Language":  "en-US,en;q=0.9",
		"Cache-Control":    "max-age=0",
		"Sec-Fetch-Mode":   "cors",
		"Accept-Encoding":  "gzip, deflate, br",
		"Origin":           "https://tracker-suivi.apps.cic.gc.ca",
		"User-Agent":       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.3 Safari/605.1.15",
		"Referer":          "https://tracker-suivi.apps.cic.gc.ca/",
		"Content-Length":   "169",
		"Sec-Fetch-Dest":   "empty",
		"X-Amz-User-Agent": "aws-amplify/5.0.4 js",
		"Priority":         "u=3, i",
		"X-Amz-Target":     "AWSCognitoIdentityProviderService.InitiateAuth",
	}

	payload := fmt.Sprintf(`{"AuthFlow":"USER_PASSWORD_AUTH","ClientId":"mtnf1qn9p739g2v8aij2anpju","AuthParameters":{"USERNAME":"%s","PASSWORD":"%s"},"ClientMetadata":{}}`, uci, password)

	body, err := Post("https://cognito-idp.ca-central-1.amazonaws.com/", headers, payload)
	if err != nil {
		// fmt.Println("Error:", err)
		return nil, err
	}

	var authResponse model.AuthResponse
	err = json.Unmarshal(body, &authResponse)
	if err != nil {
		return nil, err
	}

	return &authResponse, nil
}

func Post(url string, headers map[string]string, payload string) ([]byte, error) {
	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return nil, err
	}

	// Set the headers for the request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Print the request information
	// fmt.Println("Request URL:", req.URL)
	// fmt.Println("Request Method:", req.Method)
	// fmt.Println("Request Headers:")
	// for key, values := range req.Header {
	// 	for _, value := range values {
	// 		fmt.Printf("  %s: %s\n", key, value)
	// 	}
	// }
	// fmt.Println("Request Body:", payload)

	// Send the request using the HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status %d", resp.StatusCode)
	}

	// Read and return the response body
	return ioutil.ReadAll(resp.Body)
}

func GetStatus(authToken string, applicationNumber string) (*model.StatusResponse, error) {
	url := "https://api.tracker-suivi.apps.cic.gc.ca/user"
	payload := fmt.Sprintf(`{"method":"get-application-details","applicationNumber":"%s"}`, applicationNumber)

	headers := map[string]string{
		"Content-Type":    "application/json",
		"Accept":          "application/json",
		"Authorization":   authToken,
		"Sec-Fetch-Site":  "same-site",
		"Accept-Language": "en-US,en;q=0.9",
		"Accept-Encoding": "gzip, deflate, br",
		"Sec-Fetch-Mode":  "cors",
		"Origin":          "https://tracker-suivi.apps.cic.gc.ca",
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.3 Safari/605.1.15",
		"Content-Length":  "69",
		"Referer":         "https://tracker-suivi.apps.cic.gc.ca/",
		"Sec-Fetch-Dest":  "empty",
		"Priority":        "u=3, i",
	}

	responseBody, err := Post(url, headers, payload)
	if err != nil {
		// fmt.Println("Error:", err)
		return nil, err
	}

	var statusResponse *model.StatusResponse
	err = json.Unmarshal(responseBody, &statusResponse)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("Response: %s\n", string(responseBody))

	return statusResponse, nil
}
