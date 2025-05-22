package main

import (
	"context"
	"errors"
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
		"1.0.0",
		server.WithResourceCapabilities(false, false),
		server.WithPromptCapabilities(false),
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

	// Define device control tool
	deviceControlTool := mcp.NewTool("device_control",
		mcp.WithDescription(`Control the status or attributes of smart home devices (on/off, temperature, brightness, color, color temperature, etc.)
Examples:
	control_device([1,2], {"action":str,"attribute":str,"value":str|int,"unit":str|null})
	*action, attribute, value come from the table below*
		device_type,attribute,action,value_range,default_value,unit,category_name_zh
		ALL,on_off,"on,off","[1#on,0#off]",,,没有限制
		Light,brightness,"set,up,down",[1-100],10,%,灯
		Light,color_temperature,"['set','up','down','cooler'#调冷,'warmer'#调暖]",[6500-2700],500,K,灯
		Light,color,set,[red|green|blue|yellow|cyan|magenta|orange|purple|pink|light|pale|pale|brown|gray|...],,,灯
		AirConditioner,ac_mode,set,"['cool'#制冷模式,'dry'#除湿模式,'heat'#制热模式,'fan'#通风模式,'auto'#自动模式]",,,空调
		AirConditioner,temperature,"['set','up','down','cooler'#调冷,'warmer'#调暖]",[18-30],1,℃,空调
		AirConditioner,wind_speed,"up,down",,,,空调
		AirConditioner,wind_speed,set,"['low_speed'#低风,'medium_speed'#中风,'high_speed'#高风,'auto'#自动风]",,,空调
		AirConditioner,wind_direction,set,"['left_right'#左右扫风,'up_down'#上下扫风,'stop'#停止扫风,'on'#开始扫风]",,,空调
		WindowCovering/ClotheDryingMachine,percentage,"set,up,down",[0-100],10,%,窗帘/晾衣架
		WindowCovering/ClotheDryingMachine,motion,set,"['stop'#暂停,'continue'#继续]",,,窗帘/晾衣架
		SweepingRobot,sweep,set,"['stop'#暂停打扫,'on'#开始打扫,'continue'#继续打扫,'off'#停止打扫]",,,扫地机器人
Returns:
	Response result of device control`),
		mcp.WithArray("endpoint_ids",
			mcp.Required(),
			mcp.Description("List of ENDPOINT IDs to control (type: Int)"),
		),
		mcp.WithObject("control_params",
			mcp.Required(),
			mcp.Description("Control parameters, including action (operation), attr (attribute), and value"),
		),
	)

	// Define device inquiry tool
	deviceInquiryTool := mcp.NewTool("device_inquiry",
		mcp.WithDescription(`Get the list of devices by specified location and device type (excluding status information)
Returns:
	Device list in Markdown format`),
		mcp.WithArray("positions",
			mcp.Description("List of room names; an empty list means all rooms."),
		),
		mcp.WithArray("device_types",
			mcp.Description(`List of device types. Optional values include ["Light", "WindowCovering", "AirConditioner", "Button"]. An empty list means all types.`),
		),
	)

	// Define device status inquiry tool
	deviceStatusTool := mcp.NewTool("search_device_status",
		mcp.WithDescription(`Get the current status information of devices (used for queries related to status such as color, brightness, on/off, etc. For other queries, please use the "device_inquiry" tool)
Returns:
	Device status information in Markdown format`),
		mcp.WithArray("positions",
			mcp.Description("List of room names; an empty list means all rooms."),
		),
		mcp.WithArray("device_types",
			mcp.Description(`List of device types. Optional values include ["Light", "WindowCovering", "AirConditioner", "Button"]. An empty list means all types.`),
		),
	)

	// Define device history tool
	deviceHistoryTool := mcp.NewTool("search_device_history",
		mcp.WithDescription(`Get device history information
Returns:
	Device status information in Markdown format`),
		mcp.WithArray("endpoint_ids",
			mcp.Required(),
			mcp.Description("List of ENDPOINT IDs to query (type: Int)"),
		),
		mcp.WithString("start_datetime",
			mcp.Required(),
			mcp.Description(`Start time (format: YYYY-MM-DD HH:MM:SS, e.g., 2025-05-16 12:12:12)`),
		),
		mcp.WithString("end_datetime",
			mcp.Required(),
			mcp.Description(`End time (format: YYYY-MM-DD HH:MM:SS, e.g., 2025-05-16 12:12:12)`),
		),
		mcp.WithString("attribute",
			mcp.Description("The attribute of the device to query history for (e.g., on_off, brightness)."),
		),
	)

	// Define scene inquiry tool
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

	// Home inquiry tool
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

	// Add tools and bind processing handlers
	mcpServer.AddTool(deviceControlTool, RequestWrapper(DeviceControlHandler))
	mcpServer.AddTool(deviceInquiryTool, RequestWrapper(DeviceInquiryHandler))
	mcpServer.AddTool(deviceStatusTool, RequestWrapper(SearchDeviceStatusHandler))
	mcpServer.AddTool(deviceHistoryTool, RequestWrapper(SearchDeviceHistoryHandler))

	mcpServer.AddTool(getScenesTool, RequestWrapper(GetScenesHandler))
	mcpServer.AddTool(runSceneTool, RequestWrapper(RunSceneHandler))

	mcpServer.AddTool(getHomesTool, RequestWrapper(GetHomesHandler))
	mcpServer.AddTool(switchHomeTool, RequestWrapper(SwitchHomeHandler))

	// Start the server
	if transport == "sse" {
		host, port := "", ""

		err := viper.UnmarshalKey("port", &port)
		if err != nil {
			log.Printf("failed to unmarshal port from config: %v. Using default or expecting environment override.\n", err)
		}
		err = viper.UnmarshalKey("host", &host)
		if err != nil {
			log.Printf("failed to unmarshal host from config: %v. Using default or expecting environment override.\n", err)
		}

		baseUrl := host + ":" + port

		sseServer := server.NewSSEServer(mcpServer, server.WithBaseURL("http://"+baseUrl))
		log.Printf("SSE server listening on %s\n", baseUrl)
		if err := sseServer.Start(baseUrl); err != nil {
			log.Fatalf("SSE Server error: %v", err)
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
	endpointIDsFloat, err := extractSlice[float64](request.Params.Arguments["endpoint_ids"])
	if err != nil {
		log.Printf("[ERROR] [DeviceControlHandler] 'endpoint_ids' parameter error: %v, value: %+v", err, request.Params.Arguments["endpoint_ids"])
		return nil, fmt.Errorf("invalid 'endpoint_ids' parameter: %w", err)
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
	controlParams, ok := request.Params.Arguments["control_params"].(map[string]any)
	if !ok {
		log.Printf("[ERROR] [DeviceControlHandler] 'control_params' parameter type error: %+v", request.Params.Arguments["control_params"])
		return nil, fmt.Errorf("invalid 'control_params' parameter: expected an object")
	}

	log.Printf("[INFO] [DeviceControlHandler] Calling ControlDevice, endpointIDs: %+v, controlParams: %+v", endpointIDs, controlParams)
	// Call service
	result := ControlDevice(endpointIDs, controlParams)
	log.Printf("[INFO] [DeviceControlHandler] ControlDevice result: %v", result)

	return mcp.NewToolResultText(result), nil
}

// DeviceInquiryHandler handles device inquiry requests.
func DeviceInquiryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("[INFO] [DeviceInquiryHandler] Request parameters: %+v", request.Params.Arguments)
	positions, err := extractSlice[string](request.Params.Arguments["positions"])
	if err != nil {
		log.Printf("[ERROR] [DeviceInquiryHandler] 'positions' parameter error: %v, value: %+v", err, request.Params.Arguments["positions"])
		return nil, fmt.Errorf("invalid 'positions' parameter: %w", err)
	}

	deviceTypes, err := extractSlice[string](request.Params.Arguments["device_types"])
	if err != nil {
		log.Printf("[ERROR] [DeviceInquiryHandler] 'device_types' parameter error: %v, value: %+v", err, request.Params.Arguments["device_types"])
		return nil, fmt.Errorf("invalid 'device_types' parameter: %w", err)
	}

	log.Printf("[INFO] [DeviceInquiryHandler] Calling InquireDevice, positions: %+v, deviceTypes: %+v", positions, deviceTypes)
	result := InquireDevice(positions, deviceTypes)
	log.Printf("[INFO] [DeviceInquiryHandler] InquireDevice result: %v", result)

	return mcp.NewToolResultText(result), nil
}

// SearchDeviceStatusHandler handles device status inquiry requests.
func SearchDeviceStatusHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("[INFO] [SearchDeviceStatusHandler] Request parameters: %+v", request.Params.Arguments)
	positions, err := extractSlice[string](request.Params.Arguments["positions"])
	if err != nil {
		log.Printf("[ERROR] [SearchDeviceStatusHandler] 'positions' parameter error: %v, value: %+v", err, request.Params.Arguments["positions"])
		return nil, fmt.Errorf("invalid 'positions' parameter: %w", err)
	}

	deviceTypes, err := extractSlice[string](request.Params.Arguments["device_types"])
	if err != nil {
		log.Printf("[ERROR] [SearchDeviceStatusHandler] 'device_types' parameter error: %v, value: %+v", err, request.Params.Arguments["device_types"])
		return nil, fmt.Errorf("invalid 'device_types' parameter: %w", err)
	}

	log.Printf("[INFO] [SearchDeviceStatusHandler] Calling SearchDeviceStatus, positions: %+v, deviceTypes: %+v", positions, deviceTypes)
	result := SearchDeviceStatus(positions, deviceTypes)
	log.Printf("[INFO] [SearchDeviceStatusHandler] SearchDeviceStatus result: %v", result)

	return mcp.NewToolResultText(result), nil
}

// SearchDeviceHistoryHandler handles device history query requests.
// Currently, this feature is not implemented and returns a placeholder message.
func SearchDeviceHistoryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("[INFO] [SearchDeviceHistoryHandler] Request parameters: %+v", request.Params.Arguments)
	endpointIDsFloat, err := extractSlice[float64](request.Params.Arguments["endpoint_ids"])
	if err != nil {
		log.Printf("[ERROR] [SearchDeviceHistoryHandler] 'endpoint_ids' parameter error: %v, value: %+v", err, request.Params.Arguments["endpoint_ids"])
		return nil, fmt.Errorf("invalid 'endpoint_ids' parameter: %w", err)
	}
	var endpointIDs []int
	for _, floatVal := range endpointIDsFloat {
		endpointIDs = append(endpointIDs, int(floatVal))
	}
	if len(endpointIDs) == 0 {
		log.Printf("[WARN] [SearchDeviceHistoryHandler] 'endpoint_ids' is empty: %+v", endpointIDsFloat)
		return mcp.NewToolResultText("Please provide valid device IDs."), nil
	}
	// TODO: Implement device history fetching logic.
	// Parameters like start_datetime, end_datetime, attribute would be used here.
	return mcp.NewToolResultText("This feature will be available soon."), nil
}

// GetScenesHandler handles querying available scenes.
func GetScenesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("[INFO] [GetScenesHandler] Request parameters: %+v", request.Params.Arguments)
	positions, err := extractSlice[string](request.Params.Arguments["positions"])
	if err != nil {
		log.Printf("[ERROR] [GetScenesHandler] 'positions' parameter error: %v, value: %+v", err, request.Params.Arguments["positions"])
		return nil, fmt.Errorf("invalid 'positions' parameter: %w", err)
	}
	result := GetScenes(positions)
	log.Printf("[INFO] [GetScenesHandler] GetScenes result: %v", result)
	return mcp.NewToolResultText(result), nil
}

// RunSceneHandler handles executing a scene by its ID.
func RunSceneHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("[INFO] [RunSceneHandler] Request parameters: %+v", request.Params.Arguments)
	sceneIDsFloat, err := extractSlice[float64](request.Params.Arguments["scenes"])
	if err != nil {
		log.Printf("[ERROR] [RunSceneHandler] 'scenes' parameter error: %v, value: %+v", err, request.Params.Arguments["scenes"])
		return nil, fmt.Errorf("invalid 'scenes' parameter: %w", err)
	}

	var sceneIDs []int
	for _, floatVal := range sceneIDsFloat {
		sceneIDs = append(sceneIDs, int(floatVal))
	}

	if len(sceneIDs) == 0 {
		log.Printf("[WARN] [RunSceneHandler] 'scenes' list is empty: %+v", sceneIDsFloat)
		return mcp.NewToolResultText("Please provide valid scene."), nil
	}

	result := RunScene(sceneIDs)
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
	homeName, ok := extractString(request.Params.Arguments["home_name"])
	if !ok {
		log.Printf("[ERROR] [SwitchHomeHandler] 'home_name' parameter type error: %+v", request.Params.Arguments["home_name"])
		return nil, fmt.Errorf("invalid home name: expected a string")
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

// extractSlice extracts a typed slice from an interface{}.
// It handles cases where the input is already []T or []any.
// If elements are strings, it attempts to decode unicode sequences.
// Renamed from extractTArray.
func extractSlice[T any](value any) ([]T, error) {
	if value == nil {
		// Return empty slice instead of error if nil means "no values provided" and that's acceptable.
		// Or return an error if the parameter is mandatory.
		// For now, matching original behavior of returning an error for nil.
		return nil, errors.New("parameter is nil, expected a slice")
	}

	// Check if it's already the target type []T
	if tArray, ok := value.([]T); ok {
		// If T is string, iterate and decode unicode for each element.
		var zeroT T
		if _, isString := any(zeroT).(string); isString {
			decodedTArray := make([]T, len(tArray))
			for i, v := range tArray {
				s := any(v).(string)                                 // Cast element to string
				decodedTArray[i] = any(decodeUnicodeIfString(s)).(T) // Decode and cast back to T
			}
			return decodedTArray, nil
		}
		return tArray, nil
	}

	// Check if it's []any (common for JSON unmarshalled arrays)
	anyArray, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("parameter is not a slice, but type %T", value)
	}

	// Convert []any to []T
	newValue := make([]T, 0, len(anyArray))
	var zeroT T
	_, targetTypeIsString := any(zeroT).(string)

	for i, v := range anyArray {
		var typedVal T
		var conversionOk bool

		if targetTypeIsString {
			// If target type is string, ensure element is string and decode unicode
			sVal, isElemString := v.(string)
			if !isElemString {
				return nil, fmt.Errorf("element at index %d is type %T, expected string", i, v)
			}
			// The type T is string here, so this cast is safe.
			typedVal = any(decodeUnicodeIfString(sVal)).(T)
			conversionOk = true
		} else {
			// For non-string target types, directly assert the type.
			typedVal, conversionOk = v.(T)
		}

		if !conversionOk {
			return nil, fmt.Errorf("element at index %d (value: %v, type: %T) is not assignable to target type %T", i, v, v, zeroT)
		}
		newValue = append(newValue, typedVal)
	}
	return newValue, nil
}

// decodeUnicodeIfString attempts to decode unicode escape sequences (e.g., \uXXXX) if present in the string.
// Returns the original string if no decoding is needed or if decoding fails.
func decodeUnicodeIfString(val string) string {
	// Check for common unicode escape prefixes.
	if strings.Contains(val, `\u`) || strings.Contains(val, `\U`) {
		// strconv.Unquote expects the string to be a valid Go string literal (e.g., quoted).
		// Adding quotes around the value to make it a valid literal for Unquote.
		if decoded, err := strconv.Unquote(`"` + val + `"`); err == nil {
			return decoded
		}
		// If Unquote fails, log or handle error, and return original string as fallback.
		// log.Printf("Failed to decode unicode string '%s': %v", val, err)
	}
	return val
}

// extractString extracts a string from an interface{}, decoding unicode if necessary.
// Returns the string and a boolean indicating success.
func extractString(value any) (string, bool) {
	sValue, ok := value.(string)
	if !ok {
		return "", false
	}
	return decodeUnicodeIfString(sValue), true
}
