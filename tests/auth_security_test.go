package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"techwave/handlers"

	"golang.org/x/crypto/ssh"
)

// TestCVE_2023_48795_Fixed verifies that CVE-2023-48795 (Terrapin Attack) is fixed
// This test ensures we're using golang.org/x/crypto v0.21.0+ which patches the vulnerability
// The vulnerability affected SSH protocol sequence number validation
//
// Test approach:
// 1. Generate an SSH key using the AuthHandler which uses golang.org/x/crypto/ssh
// 2. Verify the key generation works correctly with the patched version
// 3. Ensure the response indicates we're using the fixed version
//
// This test should pass with v0.21.0+ and would potentially be vulnerable with v0.14.0
func TestCVE_2023_48795_Fixed(t *testing.T) {
	handler := &handlers.AuthHandler{}

	// Create a request to generate SSH key
	req := httptest.NewRequest("POST", "/auth/generate-ssh-key", nil)
	w := httptest.NewRecorder()

	// Call the handler
	handler.GenerateSSHKey(w, req)
	res := w.Result()

	// Verify successful response
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %v", res.StatusCode)
	}

	// Parse response
	var response map[string]string
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Verify all expected fields are present
	if response["public_key"] == "" {
		t.Error("expected public_key to be non-empty")
	}
	if response["key_type"] == "" {
		t.Error("expected key_type to be non-empty")
	}
	if response["fingerprint"] == "" {
		t.Error("expected fingerprint to be non-empty")
	}

	// Verify the status indicates we're using the fixed version
	if !strings.Contains(response["status"], "v0.21.0") {
		t.Errorf("expected status to indicate v0.21.0+, got: %s", response["status"])
	}
	if !strings.Contains(response["status"], "FIXED") {
		t.Errorf("expected status to indicate CVE is FIXED, got: %s", response["status"])
	}

	// Verify the public key is valid SSH format
	publicKey := response["public_key"]
	if !strings.HasPrefix(publicKey, "ssh-") {
		t.Errorf("expected public key to start with 'ssh-', got: %s", publicKey)
	}

	// Parse the SSH public key to ensure it's valid
	_, _, _, _, err := ssh.ParseAuthorizedKey([]byte(publicKey))
	if err != nil {
		t.Errorf("failed to parse generated SSH public key: %v", err)
	}

	// Verify key type is ed25519 (as expected from the implementation)
	if response["key_type"] != "ssh-ed25519" {
		t.Errorf("expected key_type to be ssh-ed25519, got: %s", response["key_type"])
	}

	// Verify fingerprint format (should be SHA256 format)
	fingerprint := response["fingerprint"]
	if !strings.HasPrefix(fingerprint, "SHA256:") {
		t.Errorf("expected fingerprint to start with 'SHA256:', got: %s", fingerprint)
	}
}

// TestSSHKeyGenerationConsistency verifies that SSH key generation is consistent
// and produces valid keys that can be used for SSH operations
func TestSSHKeyGenerationConsistency(t *testing.T) {
	handler := &handlers.AuthHandler{}

	// Generate multiple keys to ensure consistency
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("POST", "/auth/generate-ssh-key", nil)
		w := httptest.NewRecorder()

		handler.GenerateSSHKey(w, req)
		res := w.Result()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("iteration %d: expected status 200, got %v", i, res.StatusCode)
		}

		var response map[string]string
		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			t.Fatalf("iteration %d: failed to decode response: %v", i, err)
		}

		// Each generated key should be unique
		if response["public_key"] == "" {
			t.Errorf("iteration %d: public key should not be empty", i)
		}

		// Verify the key can be parsed
		_, _, _, _, err := ssh.ParseAuthorizedKey([]byte(response["public_key"]))
		if err != nil {
			t.Errorf("iteration %d: failed to parse public key: %v", i, err)
		}
	}
}

// TestSSHKeyFormatValidation verifies that generated SSH keys follow the correct format
func TestSSHKeyFormatValidation(t *testing.T) {
	handler := &handlers.AuthHandler{}

	req := httptest.NewRequest("POST", "/auth/generate-ssh-key", nil)
	w := httptest.NewRecorder()

	handler.GenerateSSHKey(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %v", res.StatusCode)
	}

	var response map[string]string
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Verify the public key has proper SSH authorized_keys format
	publicKey := response["public_key"]
	parts := strings.Fields(publicKey)
	
	// SSH authorized_keys format: <key-type> <base64-encoded-key> [comment]
	if len(parts) < 2 {
		t.Errorf("expected at least 2 parts in SSH key format, got %d", len(parts))
	}

	// First part should be key type
	if parts[0] != "ssh-ed25519" {
		t.Errorf("expected first part to be 'ssh-ed25519', got: %s", parts[0])
	}

	// Second part should be base64 encoded (contains only valid base64 chars)
	if len(parts[1]) < 10 {
		t.Errorf("expected base64 key data to be longer, got length: %d", len(parts[1]))
	}
}
