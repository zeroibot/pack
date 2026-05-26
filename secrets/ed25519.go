package secrets

import (
	"crypto/ed25519"
	"crypto/x509"
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
