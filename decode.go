package tipsy

import (
	"encoding/binary"
	"errors"
	"math/bits"
)

var (
	ErrCorrupt = errors.New("tipsy: corrupt input")
)

func Decode(dst, src []byte) ([]byte, error) {
	if len(src) == 0 {
		return nil, nil
	}

	prefix, sn := binary.Uvarint(src)
	if sn <= 0 {
		return nil, ErrCorrupt
	}
	src = src[sn:]

	// decode destination size
	size := prefix >> 1
	if prefix&1 == 1 {
		size = 1 << size
	}

	if size == 0 {
		return nil, nil
	}

	if len(dst) >= int(size) {
		dst = dst[:size]
		for n := range dst {
			dst[n] = 0
		}
	} else {
		dst = make([]byte, size)
	}

	// dst write index
	dn := 0

	for len(src) > 0 {
		switch bits.LeadingZeros8(src[0]) {
		case 0:
			if dn+7 > len(dst) {
				return nil, ErrCorrupt
			}

			v := src[0]
			const mask = 1

			d := dst[dn : dn+7 : dn+7]
			d[0] = v & mask
			d[1] = (v >> 1) & mask
			d[2] = (v >> 2) & mask
			d[3] = (v >> 3) & mask
			d[4] = (v >> 4) & mask
			d[5] = (v >> 5) & mask
			d[6] = (v >> 6) & mask

			src = src[1:]
			dn += 7
		case 1:
			if len(src) < 2 || dn+7 > len(dst) {
				return nil, ErrCorrupt
			}
			v := (uint64(src[0]) << 8) | uint64(src[1])

			const mask = 0b11

			d := dst[dn : dn+7 : dn+7]
			d[0] = byte(v & mask)
			d[1] = byte((v >> 2) & mask)
			d[2] = byte((v >> 4) & mask)
			d[3] = byte((v >> 6) & mask)
			d[4] = byte((v >> 8) & mask)
			d[5] = byte((v >> 10) & mask)
			d[6] = byte((v >> 12) & mask)

			src = src[2:]
			dn += 7
		case 2:
			if len(src) < 3 || dn+7 > len(dst) {
				return nil, ErrCorrupt
			}
			v := (uint64(src[0]) << 16) | (uint64(src[1]) << 8) | uint64(src[2])

			const mask = 0b111

			d := dst[dn : dn+7 : dn+7]
			d[0] = byte(v & mask)
			d[1] = byte((v >> 3) & mask)
			d[2] = byte((v >> 6) & mask)
			d[3] = byte((v >> 9) & mask)
			d[4] = byte((v >> 12) & mask)
			d[5] = byte((v >> 15) & mask)
			d[6] = byte((v >> 18) & mask)

			src = src[3:]
			dn += 7
		case 3:
			if len(src) < 4 || dn+7 > len(dst) {
				return nil, ErrCorrupt
			}
			v := (uint64(src[0]) << 24) | (uint64(src[1]) << 16) | (uint64(src[2]) << 8) | uint64(src[3])

			const mask = 0b1111

			d := dst[dn : dn+7 : dn+7]
			d[0] = byte(v & mask)
			d[1] = byte((v >> 4) & mask)
			d[2] = byte((v >> 8) & mask)
			d[3] = byte((v >> 12) & mask)
			d[4] = byte((v >> 16) & mask)
			d[5] = byte((v >> 20) & mask)
			d[6] = byte((v >> 24) & mask)

			src = src[4:]
			dn += 7
		case 4:
			if len(src) < 5 || dn+7 > len(dst) {
				return nil, ErrCorrupt
			}
			v := (uint64(src[0]) << 32) | (uint64(src[1]) << 24) | (uint64(src[2]) << 16) |
				(uint64(src[3]) << 8) | uint64(src[4])

			const mask = 0b11111

			d := dst[dn : dn+7 : dn+7]
			d[0] = byte(v & mask)
			d[1] = byte((v >> 5) & mask)
			d[2] = byte((v >> 10) & mask)
			d[3] = byte((v >> 15) & mask)
			d[4] = byte((v >> 20) & mask)
			d[5] = byte((v >> 25) & mask)
			d[6] = byte((v >> 30) & mask)

			src = src[5:]
			dn += 7
		case 5:
			if len(src) < 6 || dn+7 > len(dst) {
				return nil, ErrCorrupt
			}
			v := (uint64(src[0]) << 40) | (uint64(src[1]) << 32) | (uint64(src[2]) << 24) |
				(uint64(src[3]) << 16) | (uint64(src[4]) << 8) | uint64(src[5])

			const mask = 0b111111

			d := dst[dn : dn+7 : dn+7]
			d[0] = byte(v & mask)
			d[1] = byte((v >> 6) & mask)
			d[2] = byte((v >> 12) & mask)
			d[3] = byte((v >> 18) & mask)
			d[4] = byte((v >> 24) & mask)
			d[5] = byte((v >> 30) & mask)
			d[6] = byte((v >> 36) & mask)

			src = src[6:]
			dn += 7
		case 6:
			if len(src) < 7 || dn+7 > len(dst) {
				return nil, ErrCorrupt
			}
			v := (uint64(src[0]) << 48) | (uint64(src[1]) << 40) | (uint64(src[2]) << 32) |
				(uint64(src[3]) << 24) | (uint64(src[4]) << 16) | (uint64(src[5]) << 8) | uint64(src[6])

			const mask = 0b1111111

			d := dst[dn : dn+7 : dn+7]
			d[0] = byte(v & mask)
			d[1] = byte((v >> 7) & mask)
			d[2] = byte((v >> 14) & mask)
			d[3] = byte((v >> 21) & mask)
			d[4] = byte((v >> 28) & mask)
			d[5] = byte((v >> 35) & mask)
			d[6] = byte((v >> 42) & mask)

			src = src[7:]
			dn += 7
		case 7:
			// fixed block (up to seven bytes)

			n := 8
			if len(src) < n {
				n = len(src)
			}
			dn += copy(dst[dn:], src[1:n])
			src = src[n:]
		case 8:
			// empty blocks

			v, n := binary.Uvarint(src[1:])
			if n <= 0 {
				return nil, ErrCorrupt
			}

			dn += 7 + int(v)*7
			if dn > len(dst) {
				return nil, ErrCorrupt
			}

			src = src[1+n:]
		}
	}

	return dst, nil
}
