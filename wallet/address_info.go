package wallet

import "github.com/Rufaim/blockchain/base58"

type AddressInfo struct {
	Version       byte
	PublicKeyHash []byte
	ChechSum      []byte
}

func GetAddressInfo(address []byte) *AddressInfo {
	addrDecoded := base58.Base58Decode(address)
	return &AddressInfo{
		Version:       addrDecoded[0],
		PublicKeyHash: addrDecoded[1 : len(addrDecoded)-AddressChecksumLen],
		ChechSum:      addrDecoded[len(addrDecoded)-AddressChecksumLen:],
	}
}

func (w *Wallet) GetAddressInfo() *AddressInfo {
	pubKeyHash := HashPubKey(w.PublicKey)
	payload := make([]byte, 0, len(pubKeyHash)+1)
	payload = append(payload, Version)
	payload = append(payload, pubKeyHash...)

	return &AddressInfo{
		Version:       Version,
		PublicKeyHash: pubKeyHash,
		ChechSum:      checksum(payload),
	}
}
