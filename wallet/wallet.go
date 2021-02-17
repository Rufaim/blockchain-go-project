package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"

	pb "github.com/Rufaim/blockchain/message"

	"github.com/Rufaim/blockchain/base58"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func (w *Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)
	payload := make([]byte, 0, len(pubKeyHash)+1+AddressChecksumLen)
	payload = append(payload, Version)
	payload = append(payload, pubKeyHash...)

	checksum := checksum(payload)
	payload = append(payload, checksum...)
	return base58.Base58Encode(payload)
}

func (w *Wallet) ToProto() *pb.Wallet {
	return &pb.Wallet{
		PrivateKey: &pb.Wallet_PrivateKey{
			Curve: &pb.Wallet_PrivateKey_CurveParams{
				P:       w.PrivateKey.Curve.Params().P.Bytes(),
				N:       w.PrivateKey.Curve.Params().N.Bytes(),
				B:       w.PrivateKey.Curve.Params().B.Bytes(),
				Gx:      w.PrivateKey.Curve.Params().Gx.Bytes(),
				Gy:      w.PrivateKey.Curve.Params().Gy.Bytes(),
				BitSize: int32(w.PrivateKey.Curve.Params().BitSize),
			},
			X: w.PrivateKey.X.Bytes(),
			Y: w.PrivateKey.Y.Bytes(),
			D: w.PrivateKey.D.Bytes(),
		},
		PublicKey: w.PublicKey[:],
	}
}

func NewFromProto(pbw *pb.Wallet) *Wallet {
	curve := elliptic.P256()
	curve.Params().P = new(big.Int).SetBytes(pbw.PrivateKey.Curve.P)
	curve.Params().N = new(big.Int).SetBytes(pbw.PrivateKey.Curve.N)
	curve.Params().B = new(big.Int).SetBytes(pbw.PrivateKey.Curve.B)
	curve.Params().Gx = new(big.Int).SetBytes(pbw.PrivateKey.Curve.Gx)
	curve.Params().Gy = new(big.Int).SetBytes(pbw.PrivateKey.Curve.Gy)
	curve.Params().BitSize = int(pbw.PrivateKey.Curve.BitSize)
	return &Wallet{
		PrivateKey: ecdsa.PrivateKey{
			PublicKey: ecdsa.PublicKey{
				Curve: curve,
				X:     new(big.Int).SetBytes(pbw.PrivateKey.X),
				Y:     new(big.Int).SetBytes(pbw.PrivateKey.Y),
			},
			D: new(big.Int).SetBytes(pbw.PrivateKey.D),
		},
		PublicKey: pbw.PublicKey[:],
	}
}

func NewWallet() (*Wallet, error) {
	private, public, err := newKeyPair()
	if err != nil {
		return nil, err
	}
	return &Wallet{
		PrivateKey: *private,
		PublicKey:  public,
	}, nil
}
