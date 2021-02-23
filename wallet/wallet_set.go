package wallet

import (
	"io/ioutil"
	"os"

	pb "github.com/Rufaim/blockchain/message"
	"google.golang.org/protobuf/proto"
)

type WalletSet map[string]*Wallet

func (ws *WalletSet) LoadFromFile(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return err
	}

	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	pws := &pb.WalletSet{}
	if err := proto.Unmarshal(fileContent, pws); err != nil {
		return err
	}

	for s, pw := range pws.Set {
		(*ws)[s] = NewFromProto(pw)
	}

	return nil
}

func (ws *WalletSet) SaveToFile(filename string) error {
	pws := &pb.WalletSet{}
	pws.Set = make(map[string]*pb.Wallet)
	for s, w := range *ws {
		pws.Set[s] = w.ToProto()
	}

	encoded, err := proto.Marshal(pws)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, encoded, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (ws *WalletSet) CreateWallet() (string, error) {
	wallet, err := NewWallet()
	if err != nil {
		return "", err
	}
	address := string(wallet.GetAddress())
	(*ws)[address] = wallet
	return address, nil
}

func (ws *WalletSet) GetAllAddresses() []string {
	adds := make([]string, 0, len(*ws))
	for addr, _ := range *ws {
		adds = append(adds, addr)
	}
	return adds
}

func (ws *WalletSet) GetWalletByAddress(address string) (*Wallet, error) {
	w, ok := (*ws)[address]
	if !ok {
		return nil, ErrorWalletDoesNotExist
	}
	return w, nil
}

func NewWalletSet() *WalletSet {
	ws := WalletSet(make(map[string]*Wallet))
	return &ws
}
