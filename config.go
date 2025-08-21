package main

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
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
	CloudApiKey               = strings.TrimSpace(os.Getenv("aqara_api_key"))
	CloudAPIBase              = strings.TrimSpace(os.Getenv("aqara_base_url")) // http://host:port/echo/mcp
	DeviceID                  = ""
	AppID                     = ""
	AppSecret                 = ""
	Verbose                   = false
	CloudMode                 = strings.ToLower(os.Getenv("cloud_mode")) == "true"
	CloudCtxFromEnv *CloudCtx = nil
)

const (
	Version                         = "0.0.3"
	RequestSignatureHeaderAccessKey = "X-Access-Key"
	RequestSignatureHeaderSignature = "X-Signature"
	RequestSignatureHeaderTimestamp = "X-Timestamp"
	RequestSignatureHeaderNonce     = "X-Nonce"
	DefaultAPITimeout               = 10 * time.Second
	DefaultAPPTimeout               = 15 * time.Second
)

func init() {
	if !CloudMode {
		DeviceID = genDeviceID()
		AppID = genAppID()
		AppSecret = genSecret()
	}
	Verbose = os.Getenv("verbose") == "true"
}

func genSecret() string {
	url := CloudAPIBase + "/secret"
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

func getCloudCtx(ctx context.Context) CloudCtx {
	if CloudCtxFromEnv != nil {
		return *CloudCtxFromEnv
	}
	cloudCtxAny := ctx.Value(CloudCtx{})
	if cloudCtxAny == nil {
		return CloudCtx{
			ApiBase: CloudAPIBase,
			ApiKey:  CloudApiKey,
		}
	}
	return cloudCtxAny.(CloudCtx)
}
