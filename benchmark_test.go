package tipsy

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/golang/snappy"
	"github.com/tg/hyperloglog"
	"github.com/twmb/murmur3"
)

type fixedHash32 uint32

func (h fixedHash32) Sum32() uint32 {
	return uint32(h)
}

func hashInt(v int) fixedHash32 {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(v))
	h := murmur3.New32()
	h.Write(buf[:])
	return fixedHash32(h.Sum32())
}

func runBenchmarks(name string,
	b *testing.B,
	encoder encoderFunc,
	decoder decoderFunc,
	twoway bool) {

	var benchFunc func(b *testing.B, src []byte)

	if decoder == nil {
		if encoder == nil {
			benchFunc = func(b *testing.B, src []byte) {
				hll, _ := hyperloglog.NewReg(src)
				for n := 0; n < b.N; n++ {
					hll.Add(fixedHash32(murmur3.Sum32([]byte{1, 2, 3})))
					hll.Count()
				}
			}
		} else {
			benchFunc = func(b *testing.B, src []byte) {
				var dst []byte
				for n := 0; n < b.N; n++ {
					dst = encoder(dst, src)
				}
			}
		}
	} else if !twoway {
		benchFunc = func(b *testing.B, src []byte) {
			var dst []byte
			encoded := encoder(nil, src)
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				var err error
				dst, err = decoder(dst, encoded)
				if err != nil {
					b.Error(err)
				}
			}
		}

	} else {
		benchFunc = func(b *testing.B, src []byte) {
			var encoded, decoded []byte
			for n := 0; n < b.N; n++ {
				var err error
				encoded = encoder(encoded, src)
				decoded, err = decoder(decoded, encoded)
				if !bytes.Equal(decoded, src) {
					b.Errorf("not equal:\n%v\n%v", decoded, src)
				}
				if err != nil {
					panic(err)
				}
			}
		}
	}

	for _, maxn := range []int{1, 2, 3, 4, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 10e3, 100e3} {
		// var prec uint8 = 9
		hll, err := hyperloglog.New(12)
		if err != nil {
			panic(err)
		}
		for n := 0; n < maxn; n++ {
			hll.Add(hashInt(n))
		}

		src := hll.Registers()

		b.Run(fmt.Sprintf("%s_hll_%d_%d", name, len(src), maxn), func(b *testing.B) {
			b.SetBytes(int64(len(src)))
			benchFunc(b, src)
		})
	}
}

func runAllBenchmarks(b *testing.B, encoder encoderFunc, decoder decoderFunc) {
	runBenchmarks("encode", b, encoder, nil, false)
	runBenchmarks("decode", b, encoder, decoder, false)
	runBenchmarks("twoway", b, encoder, decoder, true)
}

var benchCases = func() [][]byte {
	cases := make([][]byte, 0)

	for _, maxn := range []int{1, 2, 3, 4, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 10e3, 100e3} {
		hll, err := hyperloglog.New(9)
		if err != nil {
			panic(err)
		}
		for n := 0; n < maxn; n++ {
			hll.Add(hashInt(n))
		}
		cases = append(cases, hll.Registers())
	}

	return cases
}

func BenchmarkTipsy(b *testing.B) {
	runAllBenchmarks(b, Encode, Decode)
}

func BenchmarkSnappy(b *testing.B) {
	runAllBenchmarks(b, snappy.Encode, snappy.Decode)
}

// Benchmark HLL runs a reference benchmark performing Add and Count on HLL.
func BenchmarkHLL(b *testing.B) {
	runAllBenchmarks(b, nil, nil)
}
