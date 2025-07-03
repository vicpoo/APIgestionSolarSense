package infrastructure

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// decodeTokenWithoutVerification decodifica manualmente el token JWT de Google sin verificación
// (solo para desarrollo/pruebas, en producción deberías usar la librería oficial)
func decodeTokenWithoutVerification(idToken string) (map[string]interface{}, error) {
	// Dividir el token JWT en sus partes
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	// Decodificar la parte de los claims (payload)
	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode token claims: %v", err)
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse token claims: %v", err)
	}

	// Extraer datos del usuario
	userData := map[string]interface{}{
		"uid":         getClaim(claims, "sub", generateUID()),
		"email":       getClaim(claims, "email", ""),
		"displayName": getClaim(claims, "name", ""),
		"photoURL":    getClaim(claims, "picture", ""),
	}

	return userData, nil
}

// generateUID genera un ID único para usuarios
func generateUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// Si falla crypto/rand, usamos timestamp como respaldo
		return fmt.Sprintf("email-%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("email-%x-%x-%x", b[0:4], b[4:6], b[6:8])
}

// getClaim extrae un claim del token o devuelve un valor por defecto
func getClaim(claims map[string]interface{}, key, defaultValue string) string {
	if val, ok := claims[key].(string); ok {
		return val
	}
	return defaultValue
}