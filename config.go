package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Global variables
var (
	Token           = ""
	Region          = ""
	DeviceID        = ""
	AppID           = ""
	AppSecret       = ""
	MCPCloudAPIBase = ""
)

const (
	Version                         = "0.0.1"
	RequestSignatureHeaderAccessKey = "X-Access-Key"
	RequestSignatureHeaderSignature = "X-Signature"
	RequestSignatureHeaderTimestamp = "X-Timestamp"
	RequestSignatureHeaderNonce     = "X-Nonce"
	DefaultAPITimeout               = 5 * time.Second
	DefaultAPPTimeout               = 8 * time.Second
)

func init() {
	Token = os.Getenv("token")
	Region = strings.ToUpper(os.Getenv("region"))
	DeviceID = genDeviceID()
	MCPCloudAPIBase = getCloudeAPIBase()
	AppID = genAppID(DeviceID)
	AppSecret = genSecret(AppID)
}

func genSecret(did string) string {
	url := MCPCloudAPIBase + "/mcp/secret"
	result, err := httpGet[map[string]string](url, map[string]string{"key": did})
	if err == nil {
		if v, ok := (*result)["secret_key"]; ok {
			return v
		}
	}
	return ""
}

func getCloudeAPIBase() string {
	var apiHost string

	switch Region {
	case "CN":
		apiHost = "https://ai-echo.aqara.cn"
	case "KR":
		apiHost = "https://ai-echo-kr.aqara.com"
	case "SG":
		apiHost = "https://ai-echo-sg.aqara.com"
	case "US":
		apiHost = "https://ai-echo-us.aqara.com"
	case "EU":
		apiHost = "https://ai-echo-ger.aqara.com"
	case "RU":
		apiHost = "https://ai-echo-ru.aqara.com"

	case "TEST":
		apiHost = "https://ai-echo-test.aqara.cn"
	case "DEV":
		apiHost = "http://localhost:5000"
	default:
		apiHost = "https://ai-echo-us.aqara.com"
	}

	return fmt.Sprintf("%s/echo", apiHost)
}

// genDeviceID generates a unique device identifier.
func genDeviceID() string {
	// Get MAC address.
	var macAddr string
	if interfaces, err := net.Interfaces(); err == nil {
		for _, inter := range interfaces {
			if inter.Flags&net.FlagLoopback == 0 && len(inter.HardwareAddr) > 0 {
				macAddr = inter.HardwareAddr.String()
				break
			}
		}
	}

	// Get hostname.
	hostname, _ := os.Hostname()

	prefix := "mcp0."
	if macAddr == "" {
		macAddr = uuid.NewString()
		prefix = "mcp1."
	}

	// Add OS information.
	osInfo := runtime.GOOS + "_" + runtime.GOARCH

	// Combine these characteristic values.
	// Using a separator for better readability of the raw string before hashing.
	rawID := fmt.Sprintf("%s|%s|%s", macAddr, hostname, osInfo)

	// Generate unique ID using SHA1 hash.
	h := sha1.Sum([]byte(rawID))
	return prefix + hex.EncodeToString(h[:])
}

func genAppID(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
