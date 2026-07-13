package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run scripts/generate_api_key.go <tenant_id> <name>")
		os.Exit(1)
	}

	tenantID := os.Args[1]
	name := os.Args[2]

	apiKey := generateRandomKey()
	hash := sha256.Sum256([]byte(apiKey))
	keyHash := hex.EncodeToString(hash[:])

	fmt.Printf("API Key: %s\n", apiKey)
	fmt.Printf("Key Hash: %s\n", keyHash)
	fmt.Printf("Tenant ID: %s\n", tenantID)
	fmt.Printf("Name: %s\n", name)
	fmt.Println("\nStore the API Key securely - it will not be shown again.")
	fmt.Println("\nInsert this into api_keys table:")
	fmt.Printf("INSERT INTO api_keys (key_hash, tenant_id, name) VALUES ('%s', '%s', '%s');\n", keyHash, tenantID, name)
}

func generateRandomKey() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
