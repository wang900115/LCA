package encode

import "github.com/btcsuite/btcutil/base58"

func Base58Encode(input []byte) string {
	return base58.Encode(input)
}

func Base58Decode(input string) []byte {
	return base58.Decode(input)
}
