package verifier

import (
	"github.com/eoscanada/eos-go"
	"github.com/eosspark/eos-go/crypto"
)

// ToSha256 From eos-go Checksum256 to geos Sha256
func ToSha256(sum eos.Checksum256) crypto.Sha256 {
	return *crypto.NewSha256Byte([]byte(sum))
}
