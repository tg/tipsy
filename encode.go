package tipsy

import (
	"encoding/binary"
	"math/bits"
)

// Encode returns encoded form of src.
// Returned slice can be subset of dst if it has enough capacity to hold the
// output; otherwise newly allocated slice is returned. Src and dst cannot overlap.
func Encode(dst, src []byte) []byte {
	dst = dst[:0]

	if len(src) == 0 {
		return dst
	}

	// write prefix
	dst = appendPrefix(dst, len(src))

	return encodeBlocks(dst, src)
}

func encodeBlocks(dst, src []byte) []byte {
	// empty is a counter of empty blocks
	var empty uint64

	for ; len(src) >= 7; src = src[7:] {
		// rebind source block
		sb := src[:7:7]

		// get bitwise OR of all the bytes to figure out maximum number of
		// bits required to fit every byte in the block.
		sum := uint8(sb[0] | sb[1] | sb[2] | sb[3] | sb[4] | sb[5] | sb[6])

		// increase empty blocks counter and skip to the next block
		if sum == 0 {
			empty++
			continue
		}

		// if we have any empty blocks collected then it's time to flush them
		if empty > 0 {
			dst = encodeEmptyBlocks(dst, empty)
			empty = 0
		}

		// encode block depeneding on max bit size of bytes in the block
		switch bits.LeadingZeros8(sum) {
		case 7:
			// bits: X654 3210
			dst = append(dst, byte((1<<7)|sb[0]|(sb[1]<<1)|(sb[2]<<2)|(sb[3]<<3)|(sb[4]<<4)|(sb[5]<<5)|(sb[6]<<6)))
		case 6:
			// bits: xX66 5544 | 3322 1100
			v := (1 << 14) | uint16(sb[0]) | (uint16(sb[1]) << 2) | (uint16(sb[2]) << 4) | (uint16(sb[3]) << 6) | (uint16(sb[4]) << 8) | (uint16(sb[5]) << 10) | (uint16(sb[6]) << 12)
			dst = append(dst,
				byte(v>>8),
				byte(v),
			)
		case 5:
			// bits: xxX6 6655 | 5444 3332 | 2211 1000
			v := (1 << 21) | uint32(sb[0]) | (uint32(sb[1]) << 3) | (uint32(sb[2]) << 6) | (uint32(sb[3]) << 9) | (uint32(sb[4]) << 12) | (uint32(sb[5]) << 15) | (uint32(sb[6]) << 18)
			dst = append(dst,
				byte(v>>16),
				byte(v>>8),
				byte(v),
			)
		case 4:
			// bits: xxxX 6666 | 5555 4444 | 3333 2222 | 1111 0000
			v := (1 << 28) | uint32(sb[0]) | (uint32(sb[1]) << 4) | (uint32(sb[2]) << 8) | (uint32(sb[3]) << 12) | (uint32(sb[4]) << 16) | (uint32(sb[5]) << 20) | (uint32(sb[6]) << 24)
			dst = append(dst,
				byte(v>>24),
				byte(v>>16),
				byte(v>>8),
				byte(v),
			)
		case 3:
			// bits: xxxx X666 | 6655 5554 | 4444 3333 | 3222 2211 | 1110 0000
			v := (1 << 35) | uint64(sb[0]) | (uint64(sb[1]) << 5) | (uint64(sb[2]) << 10) | (uint64(sb[3]) << 15) | (uint64(sb[4]) << 20) | (uint64(sb[5]) << 25) | (uint64(sb[6]) << 30)
			dst = append(dst,
				byte(v>>32),
				byte(v>>24),
				byte(v>>16),
				byte(v>>8),
				byte(v),
			)
		case 2:
			// bits: xxxx xX66 | 6666 5555 | 5544 4444 | 3333 3322 | 2222 1111 | 1100 0000
			v := (1 << 42) | uint64(sb[0]) | (uint64(sb[1]) << 6) | (uint64(sb[2]) << 12) | (uint64(sb[3]) << 18) | (uint64(sb[4]) << 24) | (uint64(sb[5]) << 30) | (uint64(sb[6]) << 36)
			dst = append(dst,
				byte(v>>40),
				byte(v>>32),
				byte(v>>24),
				byte(v>>16),
				byte(v>>8),
				byte(v),
			)
		case 1:
			// bits: xxxx xxX6 | 6666 6655 | 5555 5444 | 4444 3333 | 3332 2222 | 2211 1111 | 1000 0000
			v := (1 << 49) | uint64(sb[0]) | (uint64(sb[1]) << 7) | (uint64(sb[2]) << 14) | (uint64(sb[3]) << 21) | (uint64(sb[4]) << 28) | (uint64(sb[5]) << 35) | (uint64(sb[6]) << 42)
			dst = append(dst,
				byte(v>>48),
				byte(v>>40),
				byte(v>>32),
				byte(v>>24),
				byte(v>>16),
				byte(v>>8),
				byte(v),
			)
		default:
			// write fixed block without any encoding
			dst = append(dst, 1)
			dst = append(dst, sb...)
		}
	}

	// write residual if non-empty
	if len(src) > 0 {
		for _, b := range src {
			if b != 0 {
				// flush empty blocks if any
				if empty > 0 {
					dst = encodeEmptyBlocks(dst, empty)
				}
				// write remaining bytes as fixed block
				// TODO: could be optimized to use regular encoding
				dst = append(dst, 1)
				dst = append(dst, src...)
				break
			}
		}
	}

	return dst
}

// appendPrefix writes size to dst;
// when size is large enough power of two then we store exponent;
// first bit set to 1 says we store exponent, otherwise absolute value.
// NOTE: size shouldn't be zero, as we don't encode zero-length inputs
// and prefix value of zero is reserved for future use (streaming
// implementation when input size is not known in advance).
func appendPrefix(dst []byte, size int) []byte {
	var prefix uint64

	if size >= 64 && (size&(size-1)) == 0 {
		prefix = (uint64(bits.TrailingZeros64(uint64(size))) << 1) | 1
	} else {
		prefix = uint64(size) << 1
	}

	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, prefix)

	return append(dst, buf[:n]...)
}

// encodeEmptyBlocks writes information about consecutive empty blocks.
// If there is only one block, we use standard encoding for 1-bit block
// (as it only takes 1 byte). Otherwise we write empty byte followed by number
// of empty blocks minus one (as a varint). We never emit varint of value zero
// as this is reserved for future use (probably a stream delimiter).
func encodeEmptyBlocks(dst []byte, n uint64) []byte {
	if n == 1 {
		return append(dst, 1<<7)
	} else if n == 0 {
		panic("n is zero")
	}

	b := make([]byte, binary.MaxVarintLen64+1)
	bn := binary.PutUvarint(b[1:], n-1)
	return append(dst, b[:bn+1]...)
}
