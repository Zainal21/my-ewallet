package helpers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strings"
)

func GenerateSignature(secretKey string, payload string) string {
	header := `{"typ":"JWT","alg":"HS256"}`
	// Encode Header to Base64Url String
	headerBytes := []byte(header)
	base64UrlHeader := base64.URLEncoding.EncodeToString(headerBytes)

	// Replace characters
	base64UrlHeader = strings.Replace(base64UrlHeader, "+", "-", -1)
	base64UrlHeader = strings.Replace(base64UrlHeader, "/", "_", -1)
	base64UrlHeader = strings.Replace(base64UrlHeader, "=", "", -1)

	// Payload
	stringPayload := payload
	// payload encoded
	base64UrlPayload := base64.RawURLEncoding.EncodeToString([]byte(stringPayload))

	// Create the signature
	signingString := base64UrlHeader + "." + base64UrlPayload
	signature := GenerateHMACSHA256([]byte(secretKey), []byte(signingString))
	base64UrlSignature := base64.RawURLEncoding.EncodeToString(signature)

	// JWT token
	return base64UrlHeader + "." + base64UrlPayload + "." + base64UrlSignature
}

func GenerateHMACSHA256(key []byte, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}
