package axelar

import "github.com/bxlkm/go-wallet-sdk/coins/cosmos"

const (
	HRP = "axelar"
)

func NewAddress(privateKeyHex string) (string, error) {
	return cosmos.NewAddress(privateKeyHex, HRP, false)
}

func ValidateAddress(address string) bool {
	return cosmos.ValidateAddress(address, HRP)
}
