package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"golang.org/x/crypto/ripemd160"
)

// Checksum generates a double sha256 checksum for given data
func checksum(data []byte) []byte {
	sha := sha256.Sum256(data)
	sha = sha256.Sum256(sha[:])

	return sha[:AddressChecksumLen]
}

func HashPubKey(key []byte) []byte {
	pubKeyHash := sha256.Sum256(key)
	hasher := ripemd160.New()
	hasher.Write(pubKeyHash[:])
	return hasher.Sum(nil)
}

func newKeyPair() (*ecdsa.PrivateKey, []byte, error) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return private, pubKey, nil
}

func IsValidAddress(address []byte) bool {
	if len(address) < AddressChecksumLen+2 {
		return false
	}
	addrInfo := GetAddressInfo(address)
	if addrInfo.Version > Version {
		return false
	}
	payload := make([]byte, 0, len(addrInfo.PublicKeyHash)+1)
	payload = append(payload, addrInfo.Version)
	payload = append(payload, addrInfo.PublicKeyHash...)
	checkSum := checksum(payload)

	return bytes.Compare(checkSum, addrInfo.ChechSum) == 0
}
