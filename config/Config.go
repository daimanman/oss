package config

import (
	"strings"
)

const (
	ContextPath        string = ""
	AppPort            string = ":9090"
	MaxFileSize        int64  = 128
	SeaWeedMasterHost  string = "192.168.227.129:9333"
	SeaWeedMasterHost1 string = "192.168.1.191:9333"
)

func GetHandlerPath(basePath string) string {
	if len(basePath) > 0 {
		return ContextPath + strings.Trim(basePath, " ")
	}
	return ""
}
