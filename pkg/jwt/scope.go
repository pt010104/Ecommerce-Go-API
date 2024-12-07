package jwt

import (
	"encoding/base64"
	"encoding/json"

	"github.com/pt010104/api-golang/internal/models"
)

// NewScope creates a new scope.
func NewScope(payload Payload) models.Scope {
	return models.Scope{
		UserID: payload.UserID,
	}
}

func SetRoleToScope(scope models.Scope, role int) models.Scope {
	scope.Role = role
	return scope
}

func CreateScopeHeader(scope models.Scope) (string, error) {
	jsonData, err := json.Marshal(scope)
	if err != nil {
		return "", err
	}

	// Encode the JSON data as Base64
	base64Data := base64.StdEncoding.EncodeToString(jsonData)

	return base64Data, nil
}

func ParseScopeHeader(scopeHeader string) (models.Scope, error) {
	// Decode the Base64 data
	jsonData, err := base64.StdEncoding.DecodeString(scopeHeader)
	if err != nil {
		return models.Scope{}, err
	}

	// Unmarshal the JSON data
	var scope models.Scope
	err = json.Unmarshal(jsonData, &scope)
	if err != nil {
		return models.Scope{}, err
	}

	return scope, nil
}

func VerifyScope(scopeHeader string) (models.Scope, error) {
	scope, err := ParseScopeHeader(scopeHeader)
	if err != nil {
		return models.Scope{}, err
	}

	return scope, nil
}
