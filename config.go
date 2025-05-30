package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
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
	DefaultAPITimeout               = 10 * time.Second
	DefaultAPPTimeout               = 15 * time.Second
)

func init() {
	Token = os.Getenv("token")
	Region = strings.ToUpper(os.Getenv("region"))
	DeviceID = genDeviceID()
	MCPCloudAPIBase = getCloudAPIBase()
	AppID = genAppID()
	AppSecret = genSecret()
}

func genSecret() string {
	url := MCPCloudAPIBase + "/mcp/secret"
	result, err := httpGet[map[string]string](url, map[string]string{"key": AppID})
	if err != nil {
		log.Printf("[ERROR] Failed to generate secret: %v", err)
		return ""
	}
	if result == nil {
		log.Printf("[WARN] No secret returned from server")
		return ""
	}
	if v, ok := (*result)["secret_key"]; ok {
		return v
	}
	log.Printf("[WARN] Secret key not found in response")
	return ""
}

func getCloudAPIBase() string {
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
	var macAddr string
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && !strings.HasPrefix(i.Name, "lo") && len(i.HardwareAddr) > 0 {
				macAddr = i.HardwareAddr.String()
				break
			}
		}
	}

	prefix := "mcp0."
	if macAddr == "" {
		macAddr = uuid.NewString()
		prefix = "mcp1."
	}

	hostname, _ := os.Hostname()
	osInfo := runtime.GOOS + "-" + runtime.GOARCH

	baseInfo := strings.Join([]string{macAddr, hostname, osInfo}, "-")
	hash := sha1.New()
	hash.Write([]byte(baseInfo))

	return prefix + hex.EncodeToString(hash.Sum(nil))
}

// genAppID generates an application identifier.
func genAppID() string {
	prefix := "mcp-"
	return prefix + md5Hash(prefix+DeviceID)
}

func md5Hash(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}
