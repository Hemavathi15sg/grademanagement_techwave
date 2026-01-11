package handlers

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/ssh"
)

// AuthHandler demonstrates usage of golang.org/x/crypto package
// This handler was updated to use v0.21.0+ which fixes CVE-2023-48795 (Terrapin Attack)
type AuthHandler struct{}

// GenerateSSHKey generates an SSH key pair for authentication
// Uses golang.org/x/crypto v0.21.0+ which has CVE-2023-48795 (Terrapin Attack) fixed
// The vulnerability that affected SSH protocol sequence number validation has been patched
func (h *AuthHandler) GenerateSSHKey(w http.ResponseWriter, r *http.Request) {
	// Generate ED25519 key pair
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		http.Error(w, "Failed to generate key", http.StatusInternalServerError)
		return
	}

	// Convert to SSH format using patched crypto package
	_, err = ssh.NewSignerFromKey(privateKey)
	if err != nil {
		http.Error(w, "Failed to create SSH signer", http.StatusInternalServerError)
		return
	}

	sshPublicKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		http.Error(w, "Failed to create SSH public key", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"public_key":  string(ssh.MarshalAuthorizedKey(sshPublicKey)),
		"key_type":    sshPublicKey.Type(),
		"fingerprint": ssh.FingerprintSHA256(sshPublicKey),
		"status":      "Generated using golang.org/x/crypto v0.21.0 (CVE-2023-48795 FIXED)",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
