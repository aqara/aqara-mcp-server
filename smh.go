package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ---------- Structs ----------

// LoginResult represents the result of a login operation.
type LoginResult struct {
	Token  string `json:"token"`
	Region string `json:"region"`
}

// HomeEntity represents home entity information.
type HomeEntity struct {
	PositionName string `json:"position_name"`
	Permission   int    `json:"permission"`
	LocationId   string `json:"location_id"`
}

// RequestBody defines the general API request payload.
type RequestBody struct {
	Token     string `json:"token"`
	Region    string `json:"region"`
	Version   string `json:"version"`
	Fn        string `json:"fn"`
	Params    any    `json:"params"`
	DeviceID  string `json:"device_id"`
	RequestID string `json:"request_id"`
}

// RespBody is a generic API response structure.
type RespBody[T any] struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Result     T      `json:"result"`
	MsgDetails string `json:"msgDetails"`
}

// ---------- API Wrappers ----------

// Login authenticates a user and returns the login result and error message, if any.
func Login(username, password, region string) (*LoginResult, string) {
	result, err := CallService[LoginResult]("Login", struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Region   string `json:"region"`
	}{
		Username: username,
		Password: password,
		Region:   region,
	})
	return result, err
}

// ControlDevice sends a device control command.
func ControlDevice(devices []int, slots map[string]any) string {
	data := map[string]any{
		"devices": devices,
		"slots":   []map[string]any{slots},
	}
	_, message := CallService[string]("DeviceControl", data)
	if message != "" {
		return message
	}
	return "Device control success"
}

// InquireDevice queries the device list by positions and types.
func InquireDevice(positions []string, types []string) string {
	data := map[string]any{
		"positions":    positions,
		"device_types": types,
	}
	result, message := CallService[string]("DeviceInquiry", data)
	if message != "" {
		return message
	}
	if result == nil {
		return "No Data"
	}
	return *result
}

// SearchDeviceStatus fetches device status information.
func SearchDeviceStatus(positions []string, types []string) string {
	data := map[string]any{
		"positions":    positions,
		"device_types": types,
	}
	result, message := CallService[string]("SearchDeviceStatus", data)
	if message != "" {
		return message
	}
	if result == nil {
		return "No Data"
	}
	return *result
}

// GetScenes queries automation scenes for specified positions.
func GetScenes(positions []string) string {
	data := map[string]any{
		"positions": positions,
	}
	result, message := CallService[string]("GetScenes", data)
	if message != "" {
		return message
	}
	if result == nil {
		return "No Scenes"
	}
	return *result
}

// RunScene executes the specified scenes.
func RunScene(scenes []int) string {
	data := map[string]any{
		"scenes": scenes,
	}
	_, message := CallService[any]("RunScene", data)
	if message != "" {
		return message
	}
	return "Scene executed successfully"
}

// GetHomes retrieves the list of user homes.
func GetHomes() ([]string, string) {
	result, err := CallService[[]string]("GetHomes", nil)
	if err != "" {
		return nil, err
	}
	if result == nil {
		return nil, "No Homes"
	}
	return *result, err
}

// SwitchRegion switches the service region.
func SwitchRegion(regionName string) (bool, string) {
	result, message := CallService[string]("SwitchRegion", struct {
		RegionName string `json:"region_name"`
	}{
		RegionName: regionName,
	})
	if message != "" {
		return false, message
	}
	if result == nil {
		return false, "Region switch failed"
	}
	return true, message
}

// SwitchHome switches the current user home.
func SwitchHome(homeName string) (bool, string) {
	result, message := CallService[string]("SwitchHome", struct {
		HomeName string `json:"home_name"`
	}{
		HomeName: homeName,
	})
	if message != "" {
		return false, message
	}
	if result == nil {
		return false, "Home switch failed"
	}
	return true, message
}

// CallService calls the specific service with payload and returns parsed result and error message.
func CallService[T any](serviceName string, data any) (*T, string) {
	requestURL := MCPCloudAPIBase + "/mcp/call"
	reqData := RequestBody{
		Token:     Token,
		Region:    Region,
		Version:   Version,
		Fn:        serviceName,
		Params:    data,
		DeviceID:  DeviceID,
		RequestID: strings.Replace(uuid.NewString(), "-", "", -1),
	}
	return Post[T](requestURL, serviceName, reqData)
}

// GetHeader returns the default headers for API requests.
func GetHeader() map[string]string {
	return map[string]string{
		"app_lang":     "",
		"lang":         "",
		"app_id":       "",
		"time_zone":    "",
		"Content-Type": "application/json",
	}
}

// Post sends a POST request and returns the decoded response or error message.
func Post[T any](url string, serviceName string, body any) (*T, string) {
	headers := GetHeader()
	response, message := httpPost[T](url, body, headers)
	if message != "" {
		return nil, message
	}
	return response, ""
}

// httpPost executes a HTTP POST with necessary signing and returns the parsed result.
func httpPost[T any](url string, data any, headers map[string]string) (*T, string) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, "Data format error (invalid JSON data). Please try again later."
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, "Failed to create HTTP request: invalid parameters or request body."
	}
	// Set request headers.
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	// Add signature headers.
	{
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		bodyHash, _ := calculateSignatureRequestBodyHash(jsonData)
		signature := calculateSignature(AppSecret, request.Method, request.URL.RequestURI(), timestamp, bodyHash)

		request.Header.Add(RequestSignatureHeaderAccessKey, AppID)
		request.Header.Add(RequestSignatureHeaderTimestamp, timestamp)
		request.Header.Add(RequestSignatureHeaderNonce, generateNonce(16))
		request.Header.Add(RequestSignatureHeaderSignature, signature)
	}

	client := &http.Client{
		Timeout: DefaultAPITimeout,
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Sprintf("An error occurred while requesting the cloud service. %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Sprintf("Failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("API call failed: %s, Status code: %d, Response: %s\n", url, resp.StatusCode, string(body))
		return nil, fmt.Sprintf("API call failed: %s, Status code: %d", url, resp.StatusCode)
	}

	var result = RespBody[T]{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("JSON parsing failed: %v, Response: %s\n", err, string(body))
		if result.Message != "" {
			return nil, result.Message
		}
		return nil, "The received data is not in a valid JSON format. Please try again later."
	}
	if result.Code == 0 {
		return &result.Result, ""
	}

	log.Printf("Request error: (%d) %v, Details: %s\n", result.Code, err, result.MsgDetails)
	if result.MsgDetails != "" {
		return nil, result.MsgDetails
	}
	return nil, result.Message
}

func httpGet[T any](baseURL string, queryParams map[string]string) (*T, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL '%s': %w", baseURL, err)
	}

	if len(queryParams) > 0 {
		params := url.Values{}
		for key, value := range queryParams {
			params.Add(key, value)
		}
		parsedURL.RawQuery = params.Encode()
	}

	finalURL := parsedURL.String()
	resp, err := http.Get(finalURL)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request to '%s': %w", finalURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request to '%s' returned non-OK status: %d %s", finalURL, resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from '%s': %w", finalURL, err)
	}

	var result T

	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("JSON parsing failed: %v, Response: %s\n", err, string(body))
		return nil, fmt.Errorf("the received data is not in a valid JSON format. please try again later")
	}
	return &result, nil
}

// calculateSignature computes the signature for the request.
func calculateSignature(secret, method, path, timestamp, bodyHash string) string {
	payload := strings.Join([]string{method, path, timestamp, bodyHash}, "\n")
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

// calculateSignatureRequestBodyHash returns the SHA256 hash of the request body.
func calculateSignatureRequestBodyHash(dataBytes []byte) (string, error) {
	h := sha256.New()
	h.Write(dataBytes)
	return hex.EncodeToString(h.Sum(nil)), nil
}

// generateNonce generates a random hexadecimal string of the specified length.
func generateNonce(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("failed to generate nonce: %v\n", err)
	}
	return hex.EncodeToString(b)
}
