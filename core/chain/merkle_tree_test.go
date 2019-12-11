package chain

import (
	"crypto/sha256"
	"testing"

	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
)

func checksumFromStr(str string) Checksum256 {
	h := sha256.New()
	_, _ = h.Write([]byte(str))
	return h.Sum(nil)
}

// From github.com/eosspark/eos-go

func makeCanonicalLeft(val common.DigestType) common.DigestType {
	val.Hash[0] &= 0xFFFFFFFFFFFFFF7F
	return val
}
func makeCanonicalRight(val common.DigestType) common.DigestType {
	val.Hash[0] |= 0x0000000000000080
	return val
}

func TestMerkle(t *testing.T) {
	t1 := checksumFromStr("123")
	t2 := checksumFromStr("321")
	t.Logf("t1  : %s %s", t1.String(), t2.String())

	tt1 := crypto.NewSha256Byte(t1)
	tt2 := crypto.NewSha256Byte(t2)

	t.Logf("tt1 : %s %s", tt1.String(), tt2.String())

	if makeCanonicalLeft(*tt1).String() != mkCanonicalLeft(t1).String() {
		t.Logf("tt  : %s", makeCanonicalLeft(*tt1).String())
		t.Logf("tt  : %s", mkCanonicalLeft(t1).String())
		t.Fatalf("makeCanonicalLeft output error")
	}

	if types.Merkle([]common.DigestType{*tt1, *tt2}).String() != Merkle([]Checksum256{t1, t2}).String() {
		t.Fatalf("merkle output error")
	}
}
