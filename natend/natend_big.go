//go:build armbe || arm64be || m68k || mips || mips64 || mips64p32 || ppc || ppc64 || s390 || s390x || shbe || sparc || sparc64

package natend

import (
	"io"

	"github.com/go-perf/encoding/bigend"
)

func Uint16(b []byte) uint16 {
	return bigend.Uint16(b)
}

func PutUint16(b []byte, v uint16) {
	bigend.PutUint16(b, v)
}

func AppendUint16(b []byte, v uint16) []byte {
	return bigend.AppendUint16(b, v)
}

func Uint32(b []byte) uint32 {
	return bigend.Uint32(b)
}

func PutUint32(b []byte, v uint32) {
	bigend.PutUint32(b, v)
}

func AppendUint32(b []byte, v uint32) []byte {
	return bigend.AppendUint32(b, v)
}

func Uint64(b []byte) uint64 {
	return bigend.Uint64(b)
}

func PutUint64(b []byte, v uint64) {
	bigend.PutUint64(b, v)
}

func AppendUint64(b []byte, v uint64) []byte {
	return bigend.AppendUint64(b, v)
}

func Read(r io.Reader, data any) error {
	return bigend.Read(r, data)
}

func Write(w io.Writer, data any) error {
	return bigend.Write(w, data)
}

func Size(v any) int {
	return bigend.Size(v)
}
