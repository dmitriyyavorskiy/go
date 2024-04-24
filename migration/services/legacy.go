package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/oxxo-labs/go/models"
	"io/ioutil"
	"net/http"
)

const (
	username   = "a311f83ac0b3425ab8d6f454592c6d55"
	password   = "dc1c5ac593894809bed8b8240fd4acad"
	apiKey     = "904bf6e3-421d-437b-a7d8-da74d5f2e413"
	authServer = "https://ipaas-stg.oxxo.io"
	apiServer  = "https://api-stg.oxxo.io"
	authUrl    = "%s/integration/rest/oAuth/getToken?grant_type=client_credentials"
	docsUrl    = "%s/xapi-multichannel-outbound/api/v1/documents"
)

var (
	client = &http.Client{}
)

func createAuthHeader() string {

	auth := username + ":" + password

	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))

	header := "Basic " + encodedAuth

	fmt.Println(header)

	return header
}

func Auth() (models.TokenResponse, error) {
	authHeader := createAuthHeader()
	url := fmt.Sprintf(authUrl, authServer)

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte("")))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return models.TokenResponse{}, err
	}

	// Add headers
	req.Header.Add("Authorization", authHeader)

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return models.TokenResponse{}, err
	}
	defer resp.Body.Close()

	var response models.TokenResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	if err != nil {
		fmt.Println("Error decoding response body:", err)
		return models.TokenResponse{}, err
	}

	return response, nil
}

func GetDocumentList(token string, request models.DocumentListRequest) (models.DocumentListResponse, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshalling request body:", err)
		return models.DocumentListResponse{}, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf(docsUrl, apiServer), bytes.NewBuffer(requestBody))

	if err != nil {
		fmt.Println("Error creating request:", err)
		return models.DocumentListResponse{}, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-Gateway-APIKey", apiKey)
	req.Header.Add("Tracking-Id", "4bd9e9c2-4989-4f21-8398-ad9d19ef0f66")
	req.Header.Add("Channel-Id", "WEB")
	req.Header.Add("Country-Code", "MX")
	req.Header.Add("Language", "SPA")
	req.Header.Add("User-Agent", "postman/1.0.0")

	fmt.Printf("Request is: %v\n", req)

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return models.DocumentListResponse{}, err
	}
	defer resp.Body.Close()

	// Read and print the response body
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	fmt.Printf("Response body is: %v\n", resp.Body)

	var response models.DocumentListResponse
	//decoder := json.NewDecoder(resp.Body)
	//err = decoder.Decode(&response)
	//if err != nil {
	//	fmt.Println("Error decoding response body:", err)
	//	return models.DocumentListResponse{}, err
	//}

	return response, nil
}
