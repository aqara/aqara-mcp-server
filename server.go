package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/viper"
)

func serverStart(transport string) {
	mcpServer := server.NewMCPServer(
		"Aqara MCP Server",
		"0.0.1",
		server.WithResourceCapabilities(false, false),
		server.WithPromptCapabilities(false),
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

	// Define device control tool
	deviceControlTool := mcp.NewTool("device_control",
		mcp.WithDescription(`Control the status or attributes of smart home devices (on/off, temperature, brightness, color, color temperature, etc.)
Examples:
  control_device([1,2], {"action":"set","attribute":"brightness","value":50,"unit":"%"})
  control_device([3], {"action":"on","attribute":"on_off"})
Common operations and properties:
  Device Type | Attribute | Action | Value Range | Default | Unit | Description
  ------------|-----------|--------|-------------|---------|------|-------------
  ALL | on_off | on,off | [1=on,0=off] | - | - | Power control
  Light | brightness | set,up,down | [1-100] | 10 | % | Brightness control
  Light | color_temperature | set,up,down,cooler,warmer | [6500-2700] | 500 | K | Color temperature
  Light | color | set | red,green,blue,yellow,cyan,magenta,orange,purple,pink,white,brown,gray | - | - | Color control
  AirConditioner | ac_mode | set | cool,dry,heat,fan,auto | - | - | AC mode
  AirConditioner | temperature | set,up,down,cooler,warmer | [18-30] | 1 | °C | Temperature
  AirConditioner | wind_speed | up,down,set | low_speed,medium_speed,high_speed,auto | - | - | Fan speed
  AirConditioner | wind_direction | set | left_right,up_down,stop,on | - | - | Wind direction
  WindowCovering | percentage | set,up,down | [0-100] | 10 | % | Position control
  WindowCovering | motion | set | stop,continue | - | - | Motion control
  SweepingRobot | sweep | set | stop,on,continue,off | - | - | Sweep control
Returns:
  Device control operation result`),
		mcp.WithArray("endpoint_ids",
			mcp.Required(),
			mcp.Description("List of ENDPOINT IDs to control (type: Int)"),
		),
		mcp.WithObject("control_params",
			mcp.Required(),
			mcp.Description("Control parameters, including action (operation), attribute, and value"),
		),
	)

	// Define device query tool
	deviceQueryTool := mcp.NewTool("device_query",
		mcp.WithDescription(`Get the list of devices
Returns:
	Device list in Markdown format`),
		mcp.WithArray("positions",
			mcp.Description("List of room names; an empty list means all rooms."),
		),
		mcp.WithArray("device_types",
			mcp.Description(`List of device types. Optional values include ["Light", "WindowCovering", "AirConditioner", "Button"]. An empty list means all types.`),
		),
	)

	// Define device status query tool
	deviceStatusQueryTool := mcp.NewTool("device_status_query",
		mcp.WithDescription(`Get the current status information of devices (used for queries related to status such as color, brightness, on/off, etc. For other queries, please use the "device_query" tool)
Returns:
	Device status information in Markdown format`),
		mcp.WithArray("positions",
			mcp.Description("List of room names; an empty list means all rooms."),
		),
		mcp.WithArray("device_types",
			mcp.Description(`List of device types. Optional values include ["Light", "WindowCovering", "AirConditioner", "Button"]. An empty list means all types.`),
		),
	)

	// Define device log query tool
	deviceLogQueryTool := mcp.NewTool("device_log_query",
		mcp.WithDescription(`Get device log information
Returns:
	Device log information in Markdown format`),
		mcp.WithArray("endpoint_ids",
			mcp.Required(),
			mcp.Description("List of ENDPOINT IDs to query (type: Int)"),
		),
		mcp.WithString("start_datetime",
			mcp.Description(`Start time (format: YYYY-MM-DD HH:MM:SS, e.g., 2025-05-16 12:12:12)`),
		),
		mcp.WithString("end_datetime",
			mcp.Description(`End time (format: YYYY-MM-DD HH:MM:SS, e.g., 2025-05-16 12:12:12)`),
		),
		mcp.WithArray("attributes",
			mcp.Description("List of device attributes to query logs for (e.g., [\"on_off\", \"brightness\"])."),
		),
	)

	// Define scene query tool
	getScenesTool := mcp.NewTool("get_scenes",
		mcp.WithDescription(`Get all scene data under the user's home, or query scenes in a specified room.
Returns:
	Scene information in Markdown format`),
		mcp.WithArray("positions",
			mcp.Description("List of room names; an empty list means querying scenes for the entire home."),
		),
	)

	// Define scene execution tool
	runSceneTool := mcp.NewTool("run_scenes",
		mcp.WithDescription(`Execute the specified scene by scene ID.
Returns:
	Scene execution result message.`),
		mcp.WithArray("scenes",
			mcp.Required(),
			mcp.Description("List of scene IDs (type: Int)."),
		),
	)

	// Home query tool
	getHomesTool := mcp.NewTool("get_homes",
		mcp.WithDescription(`Get all homes under the user (useful when the user wants to query/switch homes).
Returns:
	Comma-separated list of home names; returns an empty string or specific message if no data.`),
	)

	// Switch home tool
	switchHomeTool := mcp.NewTool("switch_home",
		mcp.WithDescription(`Switch the user's current home.
Returns:
	Switch result message.`),
		mcp.WithString("home_name",
			mcp.Required(),
			mcp.Description("Target home name (should come from the available list provided by the system, e.g., via 'get_homes' tool)."),
		),
	)

	// Define automation configuration tool
	automationConfigTool := mcp.NewTool("automation_config",
		mcp.WithDescription(`Automation configuration (currently only supports scheduled delay control)
Returns:
	Automation configuration result message`),
		mcp.WithString("scheduled_time",
			mcp.Required(),
			mcp.Description("Scheduled execution time using standard Crontab format \"minute hour day month weekday\". Examples: \"30 14 * * *\" (execute daily at 14:30), \"0 9 * * 1\" (execute every Monday at 9:00)"),
		),
		mcp.WithArray("endpoint_ids",
			mcp.Required(),
			mcp.Description("List of ENDPOINT IDs to control (type: Int)"),
		),
		mcp.WithObject("control_params",
			mcp.Required(),
			mcp.Description("Device control parameters, using the same format as 'control_params' in the 'device_control' tool (including action, attribute, value, etc.)."),
		),
		mcp.WithString("task_name",
			mcp.Required(),
			mcp.Description("Name or description for this automation task (used for identification and management)."),
		),
		mcp.WithBoolean("execution_once",
			mcp.Description("Whether to execute only once. false: Execute the task periodically and repeatedly (such as daily, weekly, etc.) (type: Boolean, default true)"),
		),
	)

	// Add tools and bind processing handlers
	mcpServer.AddTool(deviceControlTool, RequestWrapper(DeviceControlHandler))
	mcpServer.AddTool(deviceQueryTool, RequestWrapper(DeviceQueryHandler))
	mcpServer.AddTool(deviceStatusQueryTool, RequestWrapper(DeviceStatusQueryHandler))
	mcpServer.AddTool(deviceLogQueryTool, RequestWrapper(DeviceLogQueryHandler))

	mcpServer.AddTool(getScenesTool, RequestWrapper(GetScenesHandler))
	mcpServer.AddTool(runSceneTool, RequestWrapper(RunSceneHandler))
	mcpServer.AddTool(automationConfigTool, RequestWrapper(AutomationConfigHandler))

	mcpServer.AddTool(getHomesTool, RequestWrapper(GetHomesHandler))
	mcpServer.AddTool(switchHomeTool, RequestWrapper(SwitchHomeHandler))

	if transport == "http" {
		host, port := "", ""

		err := viper.UnmarshalKey("port", &port)
		if err != nil {
			log.Printf("failed to unmarshal port from config: %v. Using default or expecting environment override.\n", err)
		}
		err = viper.UnmarshalKey("host", &host)
		if err != nil {
			log.Printf("failed to unmarshal host from config: %v. Using default or expecting environment override.\n", err)
		}

		startServer := server.NewStreamableHTTPServer(mcpServer)

		if startServer == nil {
			log.Fatalln("Http should not be nil")
		} else {
			baseURL := host + ":" + port
			log.Printf("HTTP server listening on %s\n", baseURL)
			// Start the server
			if err := startServer.Start(baseURL); err != nil {
				log.Fatalf("HTTP Server error: %v\n", err)
			}
		}
	} else {
		log.Println("Starting server with stdio transport.")
		if err := server.ServeStdio(mcpServer); err != nil {
			log.Fatalf("Stdio Server error: %v", err)
		}
	}
}

// DeviceControlHandler handles device control requests.
func DeviceControlHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("[INFO] [DeviceControlHandler] Request parameters: %+v", request.Params.Arguments)

	// Extract parameters from request
	endpointIDsFloat, err := extractSlice[float64](request.GetArguments()["endpoint_ids"])
	if err != nil {
		log.Printf("[ERROR] [DeviceControlHandler] 'endpoint_ids' parameter error: %v, value: %+v", err, request.GetArguments()["endpoint_ids"])
		return mcp.NewToolResultText("Invalid 'endpoint_ids' parameter: must be a list of numbers"), nil
	}

	// Convert endpoint_ids to integer array
	var endpointIDs []int
	for _, floatVal := range endpointIDsFloat {
		endpointIDs = append(endpointIDs, int(floatVal))
	}

	if len(endpointIDs) == 0 {
		log.Printf("[WARN] [DeviceControlHandler] 'endpoint_ids' is empty: %+v", endpointIDsFloat)
		return mcp.NewToolResultText("Please provide valid device IDs."), nil
	}

	// Get control parameters
	controlParams, ok := request.GetArguments()["control_params"].(map[string]any)
	if !ok {
		log.Printf("[ERROR] [DeviceControlHandler] 'control_params' parameter type error: %+v", request.GetArguments()["control_params"])
		return mcp.NewToolResultText("Invalid 'control_params' parameter: must be an object"), nil
	}

	log.Printf("[INFO] [DeviceControlHandler] Calling DeviceControl, endpointIDs: %+v, controlParams: %+v", endpointIDs, controlParams)
	// Call service
	result := DeviceControl(endpointIDs, controlParams)
	log.Printf("[INFO] [DeviceControlHandler] DeviceControl result: %v", result)

	return mcp.NewToolResultText(result), nil
}

// DeviceQueryHandler handles device query requests.
func DeviceQueryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("[INFO] [DeviceQueryHandler] Request parameters: %+v", request.Params.Arguments)
	positions, err := extractSlice[string](request.GetArguments()["positions"])
	if err != nil {
		log.Printf("[ERROR] [DeviceQueryHandler] 'positions' parameter error: %v, value: %+v", err, request.GetArguments()["positions"])
		return mcp.NewToolResultText("Invalid 'positions' parameter: must be a list of strings"), nil
	}

	deviceTypes, err := extractSlice[string](request.GetArguments()["device_types"])
	if err != nil {
		log.Printf("[ERROR] [DeviceQueryHandler] 'device_types' parameter error: %v, value: %+v", err, request.GetArguments()["device_types"])
		return mcp.NewToolResultText("Invalid 'device_types' parameter: must be a list of strings"), nil
	}

	log.Printf("[INFO] [DeviceQueryHandler] Calling DeviceQuery, positions: %+v, deviceTypes: %+v", positions, deviceTypes)
	result := DeviceQuery(positions, deviceTypes)
	log.Printf("[INFO] [DeviceQueryHandler] DeviceQuery result: %v", result)

	return mcp.NewToolResultText(result), nil
}

// DeviceStatusQueryHandler handles device status query requests.
func DeviceStatusQueryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("[INFO] [DeviceStatusQueryHandler] Request parameters: %+v", request.Params.Arguments)
	positions, err := extractSlice[string](request.GetArguments()["positions"])
	if err != nil {
		log.Printf("[ERROR] [DeviceStatusQueryHandler] 'positions' parameter error: %v, value: %+v", err, request.GetArguments()["positions"])
		return mcp.NewToolResultText("Invalid 'positions' parameter: must be a list of strings"), nil
	}

	deviceTypes, err := extractSlice[string](request.GetArguments()["device_types"])
	if err != nil {
		log.Printf("[ERROR] [DeviceStatusQueryHandler] 'device_types' parameter error: %v, value: %+v", err, request.GetArguments()["device_types"])
		return mcp.NewToolResultText("Invalid 'device_types' parameter: must be a list of strings"), nil
	}

	log.Printf("[INFO] [DeviceStatusQueryHandler] Calling DeviceStatusQuery, positions: %+v, deviceTypes: %+v", positions, deviceTypes)
	result := DeviceStatusQuery(positions, deviceTypes)
	log.Printf("[INFO] [DeviceStatusQueryHandler] DeviceStatusQuery result: %v", result)

	return mcp.NewToolResultText(result), nil
}

// DeviceLogQueryHandler handles device log query requests.
func DeviceLogQueryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("[INFO] [DeviceLogQueryHandler] Request parameters: %+v", request.Params.Arguments)

	endpointIDsFloat, err := request.RequireFloatSlice("endpoint_ids")
	if err != nil {
		log.Printf("[ERROR] [DeviceLogQueryHandler] 'endpoint_ids' parameter error: %v, value: %+v", err, request.GetArguments()["endpoint_ids"])
		return mcp.NewToolResultText("Invalid 'endpoint_ids' parameter: must be a list of numbers"), nil
	}

	var endpointIDs []int
	for _, floatVal := range endpointIDsFloat {
		endpointIDs = append(endpointIDs, int(floatVal))
	}

	if len(endpointIDs) == 0 {
		log.Printf("[WARN] [DeviceLogQueryHandler] 'endpoint_ids' is empty: %+v", endpointIDsFloat)
		return mcp.NewToolResultText("Please provide valid device IDs."), nil
	}

	// Extract optional parameters
	startDatetime, _ := request.RequireString("start_datetime")
	endDatetime, _ := request.RequireString("end_datetime")
	attributes, err := extractSlice[string](request.GetArguments()["attributes"])
	if err != nil {
		log.Printf("[ERROR] [DeviceLogQueryHandler] 'attributes' parameter error: %v, value: %+v", err, request.GetArguments()["attributes"])
		return mcp.NewToolResultText("Invalid 'attributes' parameter: must be a list of strings"), nil
	}

	log.Printf("[INFO] [DeviceLogQueryHandler] Calling DeviceLogQuery, endpointIDs: %+v, startDatetime: %s, endDatetime: %s, attributes: %v",
		endpointIDs, startDatetime, endDatetime, attributes)

	result := DeviceLogQuery(endpointIDs, startDatetime, endDatetime, attributes)
	log.Printf("[INFO] [DeviceLogQueryHandler] DeviceLogQuery result: %v", result)

	return mcp.NewToolResultText(result), nil
}

// GetScenesHandler handles querying available scenes.
func GetScenesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("[INFO] [GetScenesHandler] Request parameters: %+v", request.Params.Arguments)
	positions, err := extractSlice[string](request.GetArguments()["positions"])
	if err != nil {
		log.Printf("[ERROR] [GetScenesHandler] 'positions' parameter error: %v, value: %+v", err, request.GetArguments()["positions"])
		return mcp.NewToolResultText("Invalid 'positions' parameter: must be a list of strings"), nil
	}
	result := GetScenes(positions)
	log.Printf("[INFO] [GetScenesHandler] GetScenes result: %v", result)
	return mcp.NewToolResultText(result), nil
}

// RunSceneHandler handles executing a scene by its ID.
func RunSceneHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("[INFO] [RunSceneHandler] Request parameters: %+v", request.Params.Arguments)
	sceneIDsFloat, err := extractSlice[float64](request.GetArguments()["scenes"])
	if err != nil {
		log.Printf("[ERROR] [RunSceneHandler] 'scenes' parameter error: %v, value: %+v", err, request.GetArguments()["scenes"])
		return mcp.NewToolResultText("Invalid 'scenes' parameter: must be a list of numbers"), nil
	}

	var sceneIDs []int
	for _, floatVal := range sceneIDsFloat {
		sceneIDs = append(sceneIDs, int(floatVal))
	}

	if len(sceneIDs) == 0 {
		log.Printf("[WARN] [RunSceneHandler] 'scenes' list is empty: %+v", sceneIDsFloat)
		return mcp.NewToolResultText("Please provide valid scene IDs."), nil
	}

	result := RunScenes(sceneIDs)
	log.Printf("[INFO] [RunSceneHandler] RunScene result: %v", result)
	return mcp.NewToolResultText(result), nil
}

// GetHomesHandler retrieves all homes under the current region.
func GetHomesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("[INFO] [GetHomesHandler] Request parameters: %+v", request.Params.Arguments)
	homes, message := GetHomes()
	if message != "" {
		log.Printf("[ERROR] [GetHomesHandler] GetHomes error: %v", message)
		return mcp.NewToolResultText(message), nil
	}
	log.Printf("[INFO] [GetHomesHandler] Home list: %+v", homes)
	if len(homes) == 0 {
		return mcp.NewToolResultText("No homes found."), nil
	}
	return mcp.NewToolResultText(strings.Join(homes, ", ")), nil
}

// SwitchHomeHandler handles requests to switch current home.
func SwitchHomeHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("[INFO] [SwitchHomeHandler] Request parameters: %+v", request.Params.Arguments)
	homeName, ok := extractString(request.GetArguments()["home_name"])
	if !ok {
		log.Printf("[ERROR] [SwitchHomeHandler] 'home_name' parameter type error: %+v", request.GetArguments()["home_name"])
		return mcp.NewToolResultText("Invalid home name: expected a string"), nil
	}
	log.Printf("[INFO] [SwitchHomeHandler] Calling SwitchHome, homeName: %s", homeName)
	success, message := SwitchHome(homeName)
	if !success {
		log.Printf("[ERROR] [SwitchHomeHandler] Home switch failed: %v", message)
		// Ensure a message is always returned on failure.
		if message == "" {
			message = "Home switch failed due to an unknown error."
		}
		return mcp.NewToolResultText(message), nil
	}
	log.Printf("[INFO] [SwitchHomeHandler] Switched to home: %s", homeName)
	return mcp.NewToolResultText(fmt.Sprintf("Successfully switched to home \"%s\"", homeName)), nil
}

// AutomationConfigHandler handles requests to configure a scheduled device control task.
// Note: This handler currently only logs the request and acknowledges it.
// Actual cron job scheduling and execution are not yet implemented.
func AutomationConfigHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("[INFO] [AutomationConfigHandler] Request parameters: %+v", request.Params.Arguments)

	scheduledTime, ok := extractString(request.GetArguments()["scheduled_time"])
	if !ok || scheduledTime == "" {
		log.Printf("[ERROR] [AutomationConfigHandler] 'scheduled_time' parameter error or empty: %+v", request.GetArguments()["scheduled_time"])
		return mcp.NewToolResultText("Invalid or missing 'scheduled_time' parameter: expected a non-empty string"), nil
	}

	endpointIDsFloat, err := extractSlice[float64](request.GetArguments()["endpoint_ids"])
	if err != nil {
		log.Printf("[ERROR] [AutomationConfigHandler] 'endpoint_ids' parameter error: %v, value: %+v", err, request.GetArguments()["endpoint_ids"])
		return mcp.NewToolResultText("Invalid 'endpoint_ids' parameter: must be a list of numbers"), nil
	}

	var endpointIDs []int
	for _, floatVal := range endpointIDsFloat {
		endpointIDs = append(endpointIDs, int(floatVal))
	}

	if len(endpointIDs) == 0 {
		log.Printf("[WARN] [AutomationConfigHandler] 'endpoint_ids' is empty: %+v", endpointIDsFloat)
		return mcp.NewToolResultText("Please provide valid device IDs."), nil
	}

	controlParams, ok := request.GetArguments()["control_params"].(map[string]any)
	if !ok {
		log.Printf("[ERROR] [AutomationConfigHandler] 'control_params' parameter type error: %+v", request.GetArguments()["control_params"])
		return mcp.NewToolResultText("Invalid 'control_params' parameter: expected an object"), nil
	}

	taskName, ok := extractString(request.GetArguments()["task_name"])
	if !ok || taskName == "" {
		log.Printf("[ERROR] [AutomationConfigHandler] 'task_name' parameter error or empty: %+v", request.GetArguments()["task_name"])
		return mcp.NewToolResultText("Invalid or missing 'task_name' parameter: expected a non-empty string"), nil
	}

	// Use extractBool helper function for better type handling
	executionOnce := true // default value
	if execOnceValue, exists := request.GetArguments()["execution_once"]; exists {
		if boolVal, ok := execOnceValue.(bool); ok {
			executionOnce = boolVal
		} else {
			log.Printf("[WARN] [AutomationConfigHandler] 'execution_once' parameter type error, using default true: %+v", execOnceValue)
		}
	}

	log.Printf("[INFO] [AutomationConfigHandler] Calling AutomationConfig, scheduledTime: %s, endpointIDs: %v, controlParams: %+v, taskName: %s, executionOnce: %t",
		scheduledTime, endpointIDs, controlParams, taskName, executionOnce)

	result := AutomationConfig(scheduledTime, endpointIDs, controlParams, taskName, executionOnce)
	log.Printf("[INFO] [AutomationConfigHandler] AutomationConfig result: %v", result)

	return mcp.NewToolResultText(result), nil
}

// extractSlice extracts a typed slice from an interface{}.
// It handles cases where the input is already []T or []any.
// If elements are strings, it attempts to decode unicode sequences.
func extractSlice[T any](value any) ([]T, error) {
	if value == nil {
		return nil, nil
	}

	// Check if it's already the target type []T
	if typedSlice, ok := value.([]T); ok {
		// For string slices, decode unicode sequences
		if isStringSlice[T]() {
			result := make([]T, len(typedSlice))
			for i, item := range typedSlice {
				str := any(item).(string) // Safe cast since we know T is string
				result[i] = any(decodeUnicodeIfString(str)).(T)
			}
			return result, nil
		}
		return typedSlice, nil
	}

	isStringOfT := isStringSlice[T]()

	// Check if it's []any (common for JSON unmarshalled arrays)
	if anySlice, ok := value.([]any); ok {
		result := make([]T, 0, len(anySlice))

		// Check if T is string type once, outside the loop
		for i, item := range anySlice {
			if typedItem, ok := item.(T); ok {
				if isStringOfT {
					result = append(result, any(decodeUnicodeIfString(any(typedItem).(string))).(T))
				} else {
					result = append(result, typedItem)
				}
			} else {
				return nil, fmt.Errorf("element at index %d cannot be converted to target type", i)
			}
		}

		return result, nil
	}

	return nil, fmt.Errorf("value is not a slice")
}

// isStringSlice checks if T is string type at compile time
func isStringSlice[T any]() bool {
	var zero T
	_, ok := any(zero).(string)
	return ok
}

// decodeUnicodeIfString attempts to decode unicode escape sequences (e.g., \uXXXX) if present in the string.
// Returns the original string if no decoding is needed or if decoding fails.
func decodeUnicodeIfString(s string) string {
	// Quick check to avoid unnecessary work
	if !strings.Contains(s, "\\u") {
		return s
	}

	// Use strconv.Unquote to handle unicode escapes
	if decoded, err := strconv.Unquote(`"` + s + `"`); err == nil {
		return decoded
	}
	return s
}

// extractString extracts a string from an interface{}, decoding unicode if necessary.
// Returns the string and a boolean indicating success.
func extractString(value any) (string, bool) {
	if value == nil {
		return "", false
	}

	if str, ok := value.(string); ok {
		return decodeUnicodeIfString(str), true
	}

	return "", false
}
