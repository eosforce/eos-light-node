package verifier

import (
	"bytes"
	"crypto/sha256"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/eosspark/eos-go/crypto"
)

// ToSha256 From eos-go Checksum256 to geos Sha256
func ToSha256(sum eos.Checksum256) crypto.Sha256 {
	return *crypto.NewSha256Byte([]byte(sum))
}

func GetBlockHeaderHash(block *eos.BlockHeader) eos.Checksum256 {
	raws, _ := eos.MarshalBinary(block)

	h := sha256.New()
	_, _ = h.Write(raws)

	return h.Sum(nil)
}

func HashCheckSumPair(c1, c2 eos.Checksum256) eos.Checksum256 {
	h := sha256.New()

	if len(c1) == 0 {
		h.Write(bytes.Repeat([]byte{0}, eos.TypeSize.Checksum256))
	} else {
		h.Write(c1)
	}

	if len(c2) == 0 {
		h.Write(bytes.Repeat([]byte{0}, eos.TypeSize.Checksum256))
	} else {
		h.Write(c2)
	}

	return h.Sum(nil)
}

func HashCheckSumPairH256(c1 eos.Checksum256, c2 crypto.Sha256) eos.Checksum256 {
	h := sha256.New()

	if len(c1) == 0 {
		h.Write(bytes.Repeat([]byte{0}, eos.TypeSize.Checksum256))
	} else {
		h.Write(c1)
	}

	h.Write(c2.Bytes())

	return h.Sum(nil)
}

func IsSamePubKey(p1, p2 ecc.PublicKey) bool {
	return p1.Curve == p2.Curve && bytes.Equal(p1.Content, p2.Content)
}
