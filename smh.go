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

type CloudCtx struct {
	ApiBase string
	ApiKey  string
}

// ---------- API Wrappers ----------

// Login authenticates a user and returns the login result and error message, if any.
func Login(username, password, region string, cloudCtx CloudCtx) (*LoginResult, string) {
	if strings.TrimSpace(username) == "" {
		return nil, "Username cannot be empty"
	}
	if strings.TrimSpace(password) == "" {
		return nil, "Password cannot be empty"
	}
	if strings.TrimSpace(region) == "" {
		return nil, "Region cannot be empty"
	}

	result, err := CallService[LoginResult]("Login", struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Region   string `json:"region"`
	}{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
		Region:   strings.ToUpper(strings.TrimSpace(region)),
	}, cloudCtx)
	return result, err
}

// DeviceControl sends a device control command.
func DeviceControl(devices []int, slots map[string]any, cloudCtx CloudCtx) string {
	if len(devices) == 0 {
		return "Device list cannot be empty"
	}
	if len(slots) == 0 {
		return "Control parameters cannot be empty"
	}

	data := map[string]any{
		"devices": devices,
		"slots":   []map[string]any{slots},
	}
	_, message := CallService[string]("DeviceControl", data, cloudCtx)
	if message != "" {
		return message
	}
	return "Device control success"
}

// DeviceQuery queries the device list by positions and types.
func DeviceQuery(positions []string, types []string, cloudCtx CloudCtx) string {
	if positions == nil {
		positions = []string{}
	}
	if types == nil {
		types = []string{}
	}

	data := map[string]any{
		"positions":    positions,
		"device_types": types,
	}
	result, message := CallService[string]("DeviceQuery", data, cloudCtx)
	if message != "" {
		return message
	}
	if result == nil {
		return "No device data available"
	}
	return *result
}

// DeviceStatusQuery fetches device status information.
func DeviceStatusQuery(positions []string, types []string, cloudCtx CloudCtx) string {
	if positions == nil {
		positions = []string{}
	}
	if types == nil {
		types = []string{}
	}

	data := map[string]any{
		"positions":    positions,
		"device_types": types,
	}
	result, message := CallService[string]("DeviceStatusQuery", data, cloudCtx)
	if message != "" {
		return message
	}
	if result == nil {
		return "No device status data available"
	}
	return *result
}

// GetScenes queries automation scenes for specified positions.
func GetScenes(positions []string, cloudCtx CloudCtx) string {
	if positions == nil {
		positions = []string{}
	}

	data := map[string]any{
		"positions": positions,
	}
	result, message := CallService[string]("GetScenes", data, cloudCtx)
	if message != "" {
		return message
	}
	if result == nil {
		return "No scenes available"
	}
	return *result
}

// RunScenes executes the specified scenes.
func RunScenes(scenes []int, cloudCtx CloudCtx) string {
	if len(scenes) == 0 {
		return "Scene list cannot be empty"
	}

	data := map[string]any{
		"scenes": scenes,
	}
	_, message := CallService[any]("RunScenes", data, cloudCtx)
	if message != "" {
		return message
	}
	return "Scene executed successfully"
}

// GetHomes retrieves the list of user homes.
func GetHomes(cloudCtx CloudCtx) ([]string, string) {
	result, err := CallService[[]string]("GetHomes", nil, cloudCtx)
	if err != "" {
		return nil, err
	}
	if result == nil {
		return nil, "No homes available"
	}
	return *result, ""
}

// SwitchHome switches the current user home.
func SwitchHome(homeName string, cloudCtx CloudCtx) (bool, string) {
	if strings.TrimSpace(homeName) == "" {
		return false, "Home name cannot be empty"
	}

	result, message := CallService[string]("SwitchHome", struct {
		HomeName string `json:"home_name"`
	}{
		HomeName: strings.TrimSpace(homeName),
	}, cloudCtx)
	if message != "" {
		return false, message
	}
	if result == nil {
		return false, "Home switch failed: no response from server"
	}
	return true, ""
}

// AutomationConfig configures a scheduled device control task.
func AutomationConfig(scheduledTime string, endpointIDs []int, controlParams map[string]any, taskName string, executionOnce bool, cloudCtx CloudCtx) string {
	if strings.TrimSpace(scheduledTime) == "" {
		return "Scheduled time cannot be empty"
	}
	if len(endpointIDs) == 0 {
		return "Device list cannot be empty"
	}
	if len(controlParams) == 0 {
		return "Control parameters cannot be empty"
	}
	if strings.TrimSpace(taskName) == "" {
		return "Task name cannot be empty"
	}

	data := map[string]any{
		"scheduled_time": strings.TrimSpace(scheduledTime),
		"devices":        endpointIDs,
		"slots":          []map[string]any{controlParams},
		"task_name":      strings.TrimSpace(taskName),
		"execution_once": executionOnce,
	}

	_, message := CallService[string]("AutomationConfig", data, cloudCtx)
	if message != "" {
		return message
	}
	return "Automation configuration successful"
}

// DeviceLogQuery queries device historical log information
func DeviceLogQuery(endpointIDs []int, startDatetime, endDatetime string, attributes []string, cloudCtx CloudCtx) string {
	log.Printf("[INFO] [DeviceLogQuery] Querying device logs for endpoints: %v, start: %s, end: %s, attributes: %v",
		endpointIDs, startDatetime, endDatetime, attributes)

	if len(endpointIDs) == 0 {
		return "Device list cannot be empty"
	}

	timeSpan := make([]string, 0)

	// Add optional parameters if provided
	if strings.TrimSpace(startDatetime) != "" {
		timeSpan = append(timeSpan, strings.TrimSpace(startDatetime))
	}
	if strings.TrimSpace(endDatetime) != "" {
		timeSpan = append(timeSpan, strings.TrimSpace(endDatetime))
	}

	data := map[string]any{
		"devices":   endpointIDs,
		"time_span": timeSpan,
	}

	if len(attributes) > 0 {
		data["attributes"] = attributes
	}

	result, message := CallService[string]("DeviceLogQuery", data, cloudCtx)
	if message != "" {
		return message
	}
	if result == nil {
		return "No device log data available"
	}
	return *result
}

// CallService calls the specific service with payload and returns parsed result and error message.
func CallService[T any](serviceName string, data any, cloudCtx CloudCtx) (*T, string) {
	if cloudCtx.ApiBase == "" || cloudCtx.ApiKey == "" {
		log.Println("[ERROR] Valid API key and base URL not configured")
		return nil, "[ERROR] Valid API key and base URL not configured"
	}

	requestURL := cloudCtx.ApiBase + "/call"
	reqData := RequestBody{
		Token:     cloudCtx.ApiKey,
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
		log.Printf("[ERROR] API call failed: %s, Status code: %d, Response: %s\n", url, resp.StatusCode, string(body))
		return nil, fmt.Sprintf("API call failed. status code: %d", resp.StatusCode)
	}

	var result = RespBody[T]{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("[ERROR] JSON parsing failed: %v, Response: %s\n", err, string(body))
		if result.Message != "" {
			return nil, result.Message
		}
		return nil, "The received data is not in a valid JSON format. Please try again later."
	}
	if result.Code == 0 {
		return &result.Result, ""
	}

	log.Printf("[WARN] Request error: (%d), Details: %s\n", result.Code, result.MsgDetails)
	if result.MsgDetails != "" {
		return nil, result.MsgDetails
	}
	return nil, result.Message
}

// httpGet executes an HTTP GET request and returns the parsed result.
func httpGet[T any](baseURL string, queryParams map[string]string) (*T, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		log.Printf("[ERROR] Failed to parse base URL '%s': %v\n", baseURL, err)
		return nil, fmt.Errorf("failed to parse: %w", err)
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
		log.Printf("[ERROR] Failed to send GET request to '%s': %v\n", finalURL, err)
		return nil, fmt.Errorf("failed to send GET: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request to '%s' returned non-OK status: %d %s", finalURL, resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read response body from '%s': %v\n", finalURL, err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result T

	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("[ERROR] JSON parsing failed: %v, Response: %s\n", err, string(body))
		return nil, fmt.Errorf("the received data is not in a valid JSON format. please try again later")
	}
	return &result, nil
}

// calculateSignature computes the signature for the request.
func calculateSignature(secret, method, path, timestamp, bodyHash string) string {
	if secret == "" {
		return ""
	}
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
		log.Printf("[ERROR] Failed to generate nonce: %v\n", err)
	}
	return hex.EncodeToString(b)
}
