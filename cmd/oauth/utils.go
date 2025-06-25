package oauth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

	"hirelin-auth/internal/logger"
)

func GenerateStateHash(data StateType) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	base64Data := base64.RawURLEncoding.EncodeToString(jsonData)
	secret := logger.GetEnv("AUTH_SECRET")

	hashedData := hmac.New(sha256.New, []byte(secret))
	hashedData.Write([]byte(base64Data))
	signature := hashedData.Sum(nil)

	return base64Data + "." + hex.EncodeToString(signature), nil
}

func GetDataFromStateHash(hashedData string) (StateType, error) {
	parts := strings.Split(hashedData, ".")
	if len(parts) != 2 {
		return StateType{}, errors.New("invalid state format")
	}

	data, sig := parts[0], parts[1]
	secret := logger.GetEnv("AUTH_SECRET")

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	expectedSig := hex.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return StateType{}, errors.New("invalid state signature")
	}

	jsonBytes, err := base64.RawURLEncoding.DecodeString(data)
	if err != nil {
		return StateType{}, err
	}

	var payload StateType
	if err := json.Unmarshal(jsonBytes, &payload); err != nil {
		return StateType{}, err
	}

	return payload, nil
}
