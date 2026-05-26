package secrets

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
)

// GenerateEd25519Keys generates Ed25519 key pair and saves it to PEM files
func GenerateEd25519Keys(pubKeyPath, privKeyPath string) error {
	// 1) Generate Ed25519 key pair
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return err
	}

	// 2) Save private key to PEM file
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return err
	}
	privPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	})
	if err = os.WriteFile(privKeyPath, privPem, 0600); err != nil {
		return err
	}

	// 3) Save public key to PEM file
	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return err
	}
	pubPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})
	if err = os.WriteFile(pubKeyPath, pubPem, 0644); err != nil {
		return err
	}

	fmt.Printf("Saved Ed25519 keys to '%s' and '%s'\n", pubKeyPath, privKeyPath)
	return nil
}

// LoadEd25519PrivateKey loads Ed25519 private key from PEM file
func LoadEd25519PrivateKey(path string) (ed25519.PrivateKey, error) {
	privPem, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	pemBlock, _ := pem.Decode(privPem)

	if pemBlock == nil || pemBlock.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode private key PEM block")
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	privateKey, ok := parsedKey.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("invalid ed25519 private key")
	}

	return privateKey, nil
}

// LoadEd25519PublicKey loads Ed25519 public key from PEM file
func LoadEd25519PublicKey(path string) (ed25519.PublicKey, error) {
	pubPem, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	pemBlock, _ := pem.Decode(pubPem)

	if pemBlock == nil || pemBlock.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode public key PEM block")
	}

	parsedKey, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey, ok := parsedKey.(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid ed25519 public key")
	}

	return publicKey, nil
}

// Ed25519SignMessage generates a message signature using Ed25519
func Ed25519SignMessage(message string, privKeyPath string) (string, error) {
	privKey, err := LoadEd25519PrivateKey(privKeyPath)
	if err != nil {
		return "", err
	}

	signature := hex.EncodeToString(ed25519.Sign(privKey, []byte(message)))
	return signature, nil
}
