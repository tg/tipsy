package tipsy

import (
	"bytes"
	"testing"
	"testing/quick"
)

type encoderFunc func([]byte, []byte) []byte
type decoderFunc func([]byte, []byte) ([]byte, error)

func Test_encodeDecode(t *testing.T) {
	cases := [][]byte{
		{},
		{0},
		{1},
		{1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3},
		{1, 1, 1, 1, 7, 1, 1, 2, 2, 2, 2, 2, 7, 2, 3, 3, 3, 3, 7, 3, 3},
		{1, 1, 1, 1, 15, 1, 1, 2, 2, 2, 2, 2, 15, 2, 3, 3, 3, 3, 15, 3, 3},
		{1, 1, 1, 1, 30, 1, 1, 2, 2, 2, 2, 2, 30, 2, 3, 3, 3, 3, 30, 3, 3},
		{1, 1, 1, 1, 60, 1, 1, 2, 2, 2, 2, 2, 60, 2, 3, 3, 3, 3, 60, 3, 3},
		{1, 1, 1, 1, 120, 1, 1, 2, 2, 2, 2, 2, 120, 2, 3, 3, 3, 3, 120, 3, 3},
		{1, 1, 1, 1, 200, 1, 1, 2, 2, 2, 2, 2, 200, 2, 3, 3, 3, 3, 200, 3, 3},
		{4},
		{6},
		{17},
		{20},
		{60},
		{56, 60},
		{253},
		{255},
		{0, 0, 0, 0},
		{255, 255, 255, 255},
		{0x9b, 0x44, 0x4b},
		{0, 0, 1, 0},
		{1, 1, 1, 1},
		{0, 1, 0, 2},
		{4, 0},
		{0, 4},
		{0, 0, 4},
		{0, 0, 0, 4},
		{1, 2, 3, 4},
		{1, 2, 3, 2},
		{1, 2, 3, 2, 4},
		{0, 0, 0, 0, 0, 8},
		{0, 0, 0, 0, 0, 20},
		{0, 0, 0, 0, 0, 40},
		{255, 0, 255, 0},
		{255, 0, 255, 0, 0},
		{7},
		{7, 1},
		{2, 4, 5, 9, 4, 4, 2, 3, 6, 4, 3, 4, 3, 4, 4, 2, 3, 2, 9, 7, 1},
		{2, 4, 5, 9, 4, 4, 2, 3, 6, 4, 3, 4, 3, 4, 4, 2, 3, 2, 9, 7},
		{0, 4, 1, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 1, 2, 0, 0, 1, 0, 0, 0, 2, 0, 0, 3, 0, 1, 0, 1, 3, 0, 0, 0, 0, 1, 2, 0, 0, 0, 0, 2, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 2, 0, 0, 1, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 3, 1, 0, 3, 0, 3, 2, 0, 0, 1, 1, 2, 1, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 5, 0, 1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 1, 0, 0, 1, 0, 3, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 4, 0, 0, 1, 0, 0, 2, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 5, 3, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 7, 0, 8, 0, 0, 2, 0, 0, 0, 4, 1, 2, 0, 6, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 5, 2, 0, 0, 5, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 3, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 1, 0, 0, 1, 1, 0, 6, 0, 0, 0, 0, 5, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 1, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 2, 0, 1, 0, 0, 0, 0, 1, 1, 1, 2, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 6, 0, 0, 0, 0, 5, 0, 1, 0, 0, 1, 0, 0, 2, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 2, 0, 0, 1, 0, 0, 0, 0, 2, 0, 1, 6, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 1},

		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},

		make([]byte, 7),
		make([]byte, 14),
		make([]byte, 21),
		make([]byte, 28),
		make([]byte, 30),
		make([]byte, 100),
		make([]byte, 1000),
	}

	for cn, src := range cases {
		dst := Encode(nil, src)
		t.Logf("enc: %v => %08b\n", src, dst)
		decoded, err := Decode(nil, dst)
		if err != nil {
			t.Errorf("case %d: %s", cn, err)
		}
		if !bytes.Equal(src, decoded) {
			t.Errorf("case %d: mismatch: %v != %v", cn, src, decoded)
		}
	}
}

// test random src slices
func Test_encodeDecode_quick(t *testing.T) {
	err := quick.Check(
		func(src []byte) bool {
			encoded := Encode(nil, src)
			decoded, err := Decode(nil, encoded)
			if err != nil {
				return false
			}
			return bytes.Equal(src, decoded)
		},
		nil,
	)
	if err != nil {
		panic(err)
	}
}

func TestDecoder_reuseBuffer(t *testing.T) {
	buf := make([]byte, 100)
	for n := range buf {
		buf[n] = 255
	}

	src := []byte{100, 100, 100, 100, 100, 100, 100, 0, 0, 0, 0, 0}
	encoded := Encode(nil, src)
	decoded, err := Decode(buf, encoded)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(src, decoded) {
		t.Errorf("mismatch: %v != %v", src, decoded)
	}
}

func TestDecoder_srcTooLarge(t *testing.T) {
	src := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	encoded := Encode(nil, src)

	encoded = append(encoded, 1<<7)
	_, err := Decode(nil, encoded)

	if err != ErrCorrupt {
		t.Fatal("expected ErrCorrupt, got:", err)
	}
}

// test size of the prefix (containing size of src array)
func TestEncoder_emptySize(t *testing.T) {
	cases := []struct {
		srcSize, dstSize int
	}{
		{0, 0},
		{1, 1},
		{4, 1},
		{32, 1},
		{33, 1},

		{63, 1},
		{64, 1},
		{65, 2},

		{127, 2},
		{128, 1},
		{129, 2},

		{256, 1},

		{1 << 13, 1},
		{(1<<13 - 1), 2},

		{1 << 16, 1},
		{(1<<16 + 1), 3},
	}

	for _, tc := range cases {
		encoded := Encode(nil, make([]byte, tc.srcSize))
		if len(encoded) != tc.dstSize {
			t.Errorf("%d bytes encoded to %d bytes, expected %d",
				tc.srcSize, len(encoded), tc.dstSize)
		}
		if len(encoded) > 0 && encoded[0] == 0 {
			t.Errorf("[%d] first byte of encoded array is zero", tc.srcSize)
		}

		decoded, err := Decode(nil, encoded)
		if err != nil {
			t.Error(tc.srcSize, err)
		}
		if len(decoded) != tc.srcSize {
			t.Errorf("decoded array has %d bytes, expected %d", len(decoded), tc.srcSize)
		}
		for n := range decoded {
			if decoded[n] != 0 {
				t.Errorf("[%d] decoded array contains non-zero byte", tc.srcSize)
			}
		}
	}
}
