package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
)

/* Go Cryptography Examples
This project is a simple Go program that demonstrates two fundamental cryptography techniques: symmetric encryption
using AES and asymmetric encryption using RSA. It is designed to provide a clear, practical example of how to
implement these methods for securing data.*/

// --- Symmetric Cryptography (AES) ---

// encryptSymmetric encrypts plaintext using AES-GCM.
// It returns the ciphertext and the key used for encryption.
func encryptSymmetric(plaintext []byte) ([]byte, []byte, error) {
	// Generate a new random 256-bit (32-byte) key for AES.
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, nil, err
	}

	// Create a new AES cipher block from the key.
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	// Create a new GCM (Galois/Counter Mode) cipher, which is a modern
	// and secure mode of operation.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	// Create a nonce (number used once). GCM requires a nonce for each encryption.
	// Its size is determined by the GCM implementation.
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	// Encrypt the data. The nonce is prepended to the ciphertext.
	// This is a standard practice as the nonce is not secret and is required for decryption.
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, key, nil
}

// decryptSymmetric decrypts ciphertext using AES-GCM.
// It requires the same key that was used for encryption.
func decryptSymmetric(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// The nonce is prepended to the ciphertext, so we extract it.
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the data. If the key or nonce is incorrect, or if the ciphertext
	// has been tampered with, this will return an error.
	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// --- Asymmetric Cryptography (RSA) ---

// generateRSAKeys generates a new RSA public/private key pair.
func generateRSAKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	// Generate a new private key with a key size of 2048 bits.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

// encryptAsymmetric encrypts a message using an RSA public key.
func encryptAsymmetric(plaintext []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	// OAEP (Optimal Asymmetric Encryption Padding) is a padding scheme that
	// provides enhanced security. SHA-256 is used as the hash function.
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, plaintext, nil)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// decryptAsymmetric decrypts a message using an RSA private key.
func decryptAsymmetric(ciphertext []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// --- Main Function to Demonstrate Usage ---

func main() {
	originalMessage := "This is a secret message that needs to be protected."
	fmt.Printf("Original Message: %s\n", originalMessage)

	fmt.Println("\n---  Symmetric Cryptography (AES) ---")
	// Encrypt the message
	symmetricCiphertext, symmetricKey, err := encryptSymmetric([]byte(originalMessage))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Symmetric encryption failed: %v\n", err)
		return
	}
	fmt.Printf("AES Key (first 8 bytes): %x...\n", symmetricKey[:8])
	fmt.Printf("Symmetric Ciphertext (first 16 bytes): %x...\n", symmetricCiphertext[:16])

	// Decrypt the message
	symmetricPlaintext, err := decryptSymmetric(symmetricCiphertext, symmetricKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Symmetric decryption failed: %v\n", err)
		return
	}
	fmt.Printf("Decrypted Symmetric Plaintext: %s\n", symmetricPlaintext)

	fmt.Println("\n--- Asymmetric Cryptography (RSA) ---")
	// Generate a new key pair
	privateKey, publicKey, err := generateRSAKeys()
	if err != nil {
		fmt.Fprintf(os.Stderr, "RSA key generation failed: %v\n", err)
		return
	}

	// For demonstration, let's print the public key in PEM format
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal public key: %v\n", err)
		return
	}
	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKeyBytes,
	})
	fmt.Printf("Generated Public Key (PEM format):\n%s\n", pubKeyPEM)

	// Encrypt the message with the public key
	asymmetricCiphertext, err := encryptAsymmetric([]byte(originalMessage), publicKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Asymmetric encryption failed: %v\n", err)
		return
	}
	fmt.Printf("Asymmetric Ciphertext (first 16 bytes): %x...\n", asymmetricCiphertext[:16])

	// Decrypt the message with the private key
	asymmetricPlaintext, err := decryptAsymmetric(asymmetricCiphertext, privateKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Asymmetric decryption failed: %v\n", err)
		return
	}
	fmt.Printf("Decrypted Asymmetric Plaintext: %s\n", asymmetricPlaintext)
}
