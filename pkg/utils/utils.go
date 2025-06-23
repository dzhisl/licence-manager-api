package utils

import (
	"crypto/rand"
	"math/big"
	"strings"

	"github.com/dzhisl/license-api/pkg/config"
)

func GenLicense() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	prefix := config.AppConfig.LicensePrefix
	length := config.AppConfig.LicenseLen
	if length <= 0 {
		length = 16 // fallback default
	}

	var sb strings.Builder
	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		sb.WriteByte(charset[n.Int64()])
	}

	if prefix != "" {
		return prefix + "-" + sb.String()
	}
	return sb.String()
}
