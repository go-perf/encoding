//go:build 386 || amd64 || amd64p32 || alpha || arm || arm64 || loong64 || mipsle || mips64le || mips64p32le || nios2 || ppc64le || riscv || riscv64 || sh || wasm

package natend

import (
	"io"

	"github.com/go-perf/encoding/litend"
)

func Uint16(b []byte) uint16 {
	return litend.Uint16(b)
}

func PutUint16(b []byte, v uint16) {
	litend.PutUint16(b, v)
}

func AppendUint16(b []byte, v uint16) []byte {
	return litend.AppendUint16(b, v)
}

func Uint32(b []byte) uint32 {
	return litend.Uint32(b)
}

func PutUint32(b []byte, v uint32) {
	litend.PutUint32(b, v)
}

func AppendUint32(b []byte, v uint32) []byte {
	return litend.AppendUint32(b, v)
}

func Uint64(b []byte) uint64 {
	return litend.Uint64(b)
}

func PutUint64(b []byte, v uint64) {
	litend.PutUint64(b, v)
}

func AppendUint64(b []byte, v uint64) []byte {
	return litend.AppendUint64(b, v)
}

func Read(r io.Reader, data any) error {
	return litend.Read(r, data)
}

func Write(w io.Writer, data any) error {
	return litend.Write(w, data)
}

func Size(v any) int {
	return litend.Size(v)
}
