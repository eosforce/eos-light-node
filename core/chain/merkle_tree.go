package chain

import (
	"crypto/sha256"
)

func mkCanonicalLeft(val Checksum256) Checksum256 {
	val[0] &= 0x7F
	return val
}

func mkCanonicalRight(val Checksum256) Checksum256 {
	val[0] |= 0x80
	return val
}

func mkCanonicalPairHash(l Checksum256, r Checksum256) Checksum256 {
	h := sha256.New()
	_, _ = h.Write([]byte(mkCanonicalLeft(l)))
	_, _ = h.Write([]byte(mkCanonicalRight(r)))
	return h.Sum(nil)
}

// Merkle calc merkle hash from checksums
func Merkle(ids []Checksum256) Checksum256 {
	if 0 == len(ids) {
		return Checksum256{}
	}

	for len(ids) > 1 {
		if len(ids)%2 > 0 {
			ids = append(ids, ids[len(ids)-1])
		}

		for i := 0; i < len(ids)/2; i++ {
			ids[i] = mkCanonicalPairHash(ids[2*i], ids[(2*i)+1])
		}

		ids = ids[:len(ids)/2]
	}

	return ids[0]
}
