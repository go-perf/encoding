package bench

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
	"testing"

	"github.com/go-perf/encoding/bigend"
	"github.com/go-perf/encoding/litend"
)

func BenchmarkReadSlice1000Int32s(b *testing.B) {
	b.Run("stdlib", func(b *testing.B) {
		bsr := &byteSliceReader{}
		slice := make([]int32, 1000)
		buf := make([]byte, len(slice)*4)
		b.SetBytes(int64(len(buf)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bsr.remain = buf
			binary.Read(bsr, binary.BigEndian, slice)
		}
	})
	b.Run("litend", func(b *testing.B) {
		bsr := &byteSliceReader{}
		slice := make([]int32, 1000)
		buf := make([]byte, len(slice)*4)
		b.SetBytes(int64(len(buf)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bsr.remain = buf
			litend.Read(bsr, slice)
		}
	})
	b.Run("bigend", func(b *testing.B) {
		bsr := &byteSliceReader{}
		slice := make([]int32, 1000)
		buf := make([]byte, len(slice)*4)
		b.SetBytes(int64(len(buf)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bsr.remain = buf
			bigend.Read(bsr, slice)
		}
	})
}

// func BenchmarkReadStruct(b *testing.B) {
// b.Run("stdlib", func(b *testing.B) {
// })
// b.Run("litend", func(b *testing.B) {

// })
// b.Run("bigend", func(b *testing.B) {

// })
// 	bsr := &byteSliceReader{}
// 	var buf bytes.Buffer
// 	binary.Write(&buf, binary.BigEndian, &s)
// 	b.SetBytes(int64(dataSize(reflect.ValueOf(s))))
// 	t := s
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		bsr.remain = buf.Bytes()
// 		binary.Read(bsr, binary.BigEndian, &t)
// 	}
// 	b.StopTimer()
// 	if b.N > 0 && !reflect.DeepEqual(s, t) {
// 		b.Fatalf("struct doesn't match:\ngot  %v;\nwant %v", t, s)
// 	}
// }

func BenchmarkWriteStruct(b *testing.B) {
	b.Run("stdlib", func(b *testing.B) {
		b.SetBytes(int64(binary.Size(&s)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			binary.Write(io.Discard, binary.BigEndian, &s)
		}
	})
	b.Run("litend", func(b *testing.B) {
		b.SetBytes(int64(litend.Size(&s)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			litend.Write(io.Discard, &s)
		}
	})
	b.Run("bigend", func(b *testing.B) {
		b.SetBytes(int64(bigend.Size(&s)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bigend.Write(io.Discard, &s)
		}
	})
}

func BenchmarkReadInts(b *testing.B) {
	b.Run("stdlib", func(b *testing.B) {
		var ls Struct
		bsr := &byteSliceReader{}
		var r io.Reader = bsr
		b.SetBytes(2 * (1 + 2 + 4 + 8))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bsr.remain = big
			binary.Read(r, binary.BigEndian, &ls.Int8)
			binary.Read(r, binary.BigEndian, &ls.Int16)
			binary.Read(r, binary.BigEndian, &ls.Int32)
			binary.Read(r, binary.BigEndian, &ls.Int64)
			binary.Read(r, binary.BigEndian, &ls.Uint8)
			binary.Read(r, binary.BigEndian, &ls.Uint16)
			binary.Read(r, binary.BigEndian, &ls.Uint32)
			binary.Read(r, binary.BigEndian, &ls.Uint64)
		}
		b.StopTimer()
		want := s
		want.Float32 = 0
		want.Float64 = 0
		want.Complex64 = 0
		want.Complex128 = 0
		want.Array = [4]uint8{0, 0, 0, 0}
		want.Bool = false
		want.BoolArray = [4]bool{false, false, false, false}
		if b.N > 0 && !reflect.DeepEqual(ls, want) {
			b.Fatalf("struct doesn't match:\ngot  %v;\nwant %v", ls, want)
		}
	})
	b.Run("litend", func(b *testing.B) {
		b.Skip()
		var ls Struct
		bsr := &byteSliceReader{}
		var r io.Reader = bsr
		b.SetBytes(2 * (1 + 2 + 4 + 8))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bsr.remain = big
			litend.Read(r, &ls.Int8)
			litend.Read(r, &ls.Int16)
			litend.Read(r, &ls.Int32)
			litend.Read(r, &ls.Int64)
			litend.Read(r, &ls.Uint8)
			litend.Read(r, &ls.Uint16)
			litend.Read(r, &ls.Uint32)
			litend.Read(r, &ls.Uint64)
		}
		b.StopTimer()
		want := s
		want.Float32 = 0
		want.Float64 = 0
		want.Complex64 = 0
		want.Complex128 = 0
		want.Array = [4]uint8{0, 0, 0, 0}
		want.Bool = false
		want.BoolArray = [4]bool{false, false, false, false}
		if b.N > 0 && !reflect.DeepEqual(ls, want) {
			b.Fatalf("struct doesn't match:\ngot  %v;\nwant %v", ls, want)
		}
	})
	b.Run("bigend", func(b *testing.B) {
		var ls Struct
		bsr := &byteSliceReader{}
		var r io.Reader = bsr
		b.SetBytes(2 * (1 + 2 + 4 + 8))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bsr.remain = big
			bigend.Read(r, &ls.Int8)
			bigend.Read(r, &ls.Int16)
			bigend.Read(r, &ls.Int32)
			bigend.Read(r, &ls.Int64)
			bigend.Read(r, &ls.Uint8)
			bigend.Read(r, &ls.Uint16)
			bigend.Read(r, &ls.Uint32)
			bigend.Read(r, &ls.Uint64)
		}
		b.StopTimer()
		want := s
		want.Float32 = 0
		want.Float64 = 0
		want.Complex64 = 0
		want.Complex128 = 0
		want.Array = [4]uint8{0, 0, 0, 0}
		want.Bool = false
		want.BoolArray = [4]bool{false, false, false, false}
		if b.N > 0 && !reflect.DeepEqual(ls, want) {
			b.Fatalf("struct doesn't match:\ngot  %v;\nwant %v", ls, want)
		}
	})
}

func BenchmarkWriteInts(b *testing.B) {
	b.Run("stdlib", func(b *testing.B) {
		buf := new(bytes.Buffer)
		var w io.Writer = buf
		b.SetBytes(2 * (1 + 2 + 4 + 8))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			binary.Write(w, binary.BigEndian, s.Int8)
			binary.Write(w, binary.BigEndian, s.Int16)
			binary.Write(w, binary.BigEndian, s.Int32)
			binary.Write(w, binary.BigEndian, s.Int64)
			binary.Write(w, binary.BigEndian, s.Uint8)
			binary.Write(w, binary.BigEndian, s.Uint16)
			binary.Write(w, binary.BigEndian, s.Uint32)
			binary.Write(w, binary.BigEndian, s.Uint64)
		}
		b.StopTimer()
		if b.N > 0 && !bytes.Equal(buf.Bytes(), big[:30]) {
			b.Fatalf("first half doesn't match: %x %x", buf.Bytes(), big[:30])
		}
	})
	b.Run("litend", func(b *testing.B) {
		b.Skip()
		buf := new(bytes.Buffer)
		var w io.Writer = buf
		b.SetBytes(2 * (1 + 2 + 4 + 8))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			litend.Write(w, s.Int8)
			litend.Write(w, s.Int16)
			litend.Write(w, s.Int32)
			litend.Write(w, s.Int64)
			litend.Write(w, s.Uint8)
			litend.Write(w, s.Uint16)
			litend.Write(w, s.Uint32)
			litend.Write(w, s.Uint64)
		}
		b.StopTimer()
		if b.N > 0 && !bytes.Equal(buf.Bytes(), big[:30]) {
			b.Fatalf("first half doesn't match: %x %x", buf.Bytes(), big[:30])
		}
	})
	b.Run("bigend", func(b *testing.B) {
		buf := new(bytes.Buffer)
		var w io.Writer = buf
		b.SetBytes(2 * (1 + 2 + 4 + 8))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			bigend.Write(w, s.Int8)
			bigend.Write(w, s.Int16)
			bigend.Write(w, s.Int32)
			bigend.Write(w, s.Int64)
			bigend.Write(w, s.Uint8)
			bigend.Write(w, s.Uint16)
			bigend.Write(w, s.Uint32)
			bigend.Write(w, s.Uint64)
		}
		b.StopTimer()
		if b.N > 0 && !bytes.Equal(buf.Bytes(), big[:30]) {
			b.Fatalf("first half doesn't match: %x %x", buf.Bytes(), big[:30])
		}
	})
}

func BenchmarkWriteSlice1000Int32s(b *testing.B) {
	b.Run("stdlib", func(b *testing.B) {
		slice := make([]int32, 1000)
		buf := new(bytes.Buffer)
		var w io.Writer = buf
		b.SetBytes(4 * 1000)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			binary.Write(w, binary.BigEndian, slice)
		}
		b.StopTimer()
	})
	b.Run("litend", func(b *testing.B) {
		slice := make([]int32, 1000)
		buf := new(bytes.Buffer)
		var w io.Writer = buf
		b.SetBytes(4 * 1000)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			litend.Write(w, slice)
		}
		b.StopTimer()
	})
	b.Run("bigend", func(b *testing.B) {
		slice := make([]int32, 1000)
		buf := new(bytes.Buffer)
		var w io.Writer = buf
		b.SetBytes(4 * 1000)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			bigend.Write(w, slice)
		}
		b.StopTimer()
	})
}

func BenchmarkPutUint16(b *testing.B) {
	b.Run("stdlib", func(b *testing.B) {
		b.SetBytes(2)
		for i := 0; i < b.N; i++ {
			binary.BigEndian.PutUint16(putbuf[:2], uint16(i))
		}
	})
	b.Run("litend", func(b *testing.B) {
		b.SetBytes(2)
		for i := 0; i < b.N; i++ {
			litend.PutUint16(putbuf[:2], uint16(i))
		}
	})
	b.Run("bigend", func(b *testing.B) {
		b.SetBytes(2)
		for i := 0; i < b.N; i++ {
			bigend.PutUint16(putbuf[:2], uint16(i))
		}
	})
}

func BenchmarkAppendUint16(b *testing.B) {
	b.Run("stdlib", func(b *testing.B) {
		b.SetBytes(2)
		for i := 0; i < b.N; i++ {
			putbuf = binary.BigEndian.AppendUint16(putbuf[:0], uint16(i))
		}
	})
	b.Run("litend", func(b *testing.B) {
		b.SetBytes(2)
		for i := 0; i < b.N; i++ {
			putbuf = litend.AppendUint16(putbuf[:0], uint16(i))
		}
	})
	b.Run("bigend", func(b *testing.B) {
		b.SetBytes(2)
		for i := 0; i < b.N; i++ {
			putbuf = bigend.AppendUint16(putbuf[:0], uint16(i))
		}
	})
}

func BenchmarkPutUint32(b *testing.B) {
	b.Run("stdlib", func(b *testing.B) {
		b.SetBytes(4)
		for i := 0; i < b.N; i++ {
			binary.BigEndian.PutUint32(putbuf[:4], uint32(i))
		}
	})
	b.Run("litend", func(b *testing.B) {
		b.SetBytes(4)
		for i := 0; i < b.N; i++ {
			litend.PutUint32(putbuf[:4], uint32(i))
		}
	})
	b.Run("bigend", func(b *testing.B) {
		b.SetBytes(4)
		for i := 0; i < b.N; i++ {
			bigend.PutUint32(putbuf[:4], uint32(i))
		}
	})
}

func BenchmarkAppendUint32(b *testing.B) {
	b.Run("stdlib", func(b *testing.B) {
		b.SetBytes(4)
		for i := 0; i < b.N; i++ {
			putbuf = binary.BigEndian.AppendUint32(putbuf[:0], uint32(i))
		}
	})
	b.Run("litend", func(b *testing.B) {
		b.SetBytes(4)
		for i := 0; i < b.N; i++ {
			putbuf = litend.AppendUint32(putbuf[:0], uint32(i))
		}
	})
	b.Run("bigend", func(b *testing.B) {
		b.SetBytes(4)
		for i := 0; i < b.N; i++ {
			putbuf = bigend.AppendUint32(putbuf[:0], uint32(i))
		}
	})
}

func BenchmarkPutUint64(b *testing.B) {
	b.Run("stdlib", func(b *testing.B) {
		b.SetBytes(8)
		for i := 0; i < b.N; i++ {
			binary.BigEndian.PutUint64(putbuf[:8], uint64(i))
		}
	})
	b.Run("litend", func(b *testing.B) {
		b.SetBytes(8)
		for i := 0; i < b.N; i++ {
			litend.PutUint64(putbuf[:8], uint64(i))
		}
	})
	b.Run("bigend", func(b *testing.B) {
		b.SetBytes(8)
		for i := 0; i < b.N; i++ {
			bigend.PutUint64(putbuf[:8], uint64(i))
		}
	})
}

func BenchmarkAppendUint64(b *testing.B) {
	b.Run("stdlib", func(b *testing.B) {
		b.SetBytes(8)
		for i := 0; i < b.N; i++ {
			putbuf = binary.BigEndian.AppendUint64(putbuf[:0], uint64(i))
		}
	})
	b.Run("litend", func(b *testing.B) {
		b.SetBytes(8)
		for i := 0; i < b.N; i++ {
			putbuf = litend.AppendUint64(putbuf[:0], uint64(i))
		}
	})
	b.Run("bigend", func(b *testing.B) {
		b.SetBytes(8)
		for i := 0; i < b.N; i++ {
			putbuf = bigend.AppendUint64(putbuf[:0], uint64(i))
		}
	})
}

func BenchmarkLittleEndianPutUint16(b *testing.B) {
	b.Run("stdlib", func(b *testing.B) {
		b.SetBytes(2)
		for i := 0; i < b.N; i++ {
			binary.LittleEndian.PutUint16(putbuf[:2], uint16(i))
		}
	})
	b.Run("litend", func(b *testing.B) {
		b.SetBytes(2)
		for i := 0; i < b.N; i++ {
			litend.PutUint16(putbuf[:2], uint16(i))
		}
	})
	b.Run("bigend", func(b *testing.B) {
		b.SetBytes(2)
		for i := 0; i < b.N; i++ {
			bigend.PutUint16(putbuf[:2], uint16(i))
		}
	})
}

func BenchmarkLittleEndianAppendUint16(b *testing.B) {
	b.Run("stdlib", func(b *testing.B) {
		b.SetBytes(2)
		for i := 0; i < b.N; i++ {
			putbuf = binary.LittleEndian.AppendUint16(putbuf[:0], uint16(i))
		}
	})
	b.Run("litend", func(b *testing.B) {
		b.SetBytes(2)
		for i := 0; i < b.N; i++ {
			putbuf = litend.AppendUint16(putbuf[:0], uint16(i))
		}
	})
	b.Run("bigend", func(b *testing.B) {
		b.SetBytes(2)
		for i := 0; i < b.N; i++ {
			putbuf = bigend.AppendUint16(putbuf[:0], uint16(i))
		}
	})
}

func BenchmarkLittleEndianPutUint32(b *testing.B) {
	b.Run("stdlib", func(b *testing.B) {
		b.SetBytes(4)
		for i := 0; i < b.N; i++ {
			binary.LittleEndian.PutUint32(putbuf[:4], uint32(i))
		}
	})
	b.Run("litend", func(b *testing.B) {
		b.SetBytes(4)
		for i := 0; i < b.N; i++ {
			litend.PutUint32(putbuf[:4], uint32(i))
		}
	})
	b.Run("bigend", func(b *testing.B) {
		b.SetBytes(4)
		for i := 0; i < b.N; i++ {
			bigend.PutUint32(putbuf[:4], uint32(i))
		}
	})
}

func BenchmarkLittleEndianAppendUint32(b *testing.B) {
	b.Run("stdlib", func(b *testing.B) {
		b.SetBytes(4)
		for i := 0; i < b.N; i++ {
			putbuf = binary.LittleEndian.AppendUint32(putbuf[:0], uint32(i))
		}
	})
	b.Run("litend", func(b *testing.B) {
		b.SetBytes(4)
		for i := 0; i < b.N; i++ {
			putbuf = litend.AppendUint32(putbuf[:0], uint32(i))
		}
	})
	b.Run("bigend", func(b *testing.B) {
		b.SetBytes(4)
		for i := 0; i < b.N; i++ {
			putbuf = bigend.AppendUint32(putbuf[:0], uint32(i))
		}
	})
}

func BenchmarkLittleEndianPutUint64(b *testing.B) {
	b.Run("stdlib", func(b *testing.B) {
		b.SetBytes(8)
		for i := 0; i < b.N; i++ {
			binary.LittleEndian.PutUint64(putbuf[:8], uint64(i))
		}
	})
	b.Run("litend", func(b *testing.B) {
		b.SetBytes(8)
		for i := 0; i < b.N; i++ {
			litend.PutUint64(putbuf[:8], uint64(i))
		}
	})
	b.Run("bigend", func(b *testing.B) {
		b.SetBytes(8)
		for i := 0; i < b.N; i++ {
			bigend.PutUint64(putbuf[:8], uint64(i))
		}
	})
}

func BenchmarkLittleEndianAppendUint64(b *testing.B) {
	b.Run("stdlib", func(b *testing.B) {
		b.SetBytes(8)
		for i := 0; i < b.N; i++ {
			putbuf = binary.LittleEndian.AppendUint64(putbuf[:0], uint64(i))
		}
	})
	b.Run("litend", func(b *testing.B) {
		b.SetBytes(8)
		for i := 0; i < b.N; i++ {
			putbuf = litend.AppendUint64(putbuf[:0], uint64(i))
		}
	})
	b.Run("bigend", func(b *testing.B) {
		b.SetBytes(8)
		for i := 0; i < b.N; i++ {
			putbuf = bigend.AppendUint64(putbuf[:0], uint64(i))
		}
	})
}

func BenchmarkReadFloats(b *testing.B) {
	b.Run("stdlib", func(b *testing.B) {
		var ls Struct
		bsr := &byteSliceReader{}
		var r io.Reader = bsr
		b.SetBytes(4 + 8)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bsr.remain = big[30:]
			binary.Read(r, binary.BigEndian, &ls.Float32)
			binary.Read(r, binary.BigEndian, &ls.Float64)
		}
		b.StopTimer()
		want := s
		want.Int8 = 0
		want.Int16 = 0
		want.Int32 = 0
		want.Int64 = 0
		want.Uint8 = 0
		want.Uint16 = 0
		want.Uint32 = 0
		want.Uint64 = 0
		want.Complex64 = 0
		want.Complex128 = 0
		want.Array = [4]uint8{0, 0, 0, 0}
		want.Bool = false
		want.BoolArray = [4]bool{false, false, false, false}
		if b.N > 0 && !reflect.DeepEqual(ls, want) {
			b.Fatalf("struct doesn't match:\ngot  %v;\nwant %v", ls, want)
		}
	})
	b.Run("litend", func(b *testing.B) {
		b.Skip()
		var ls Struct
		bsr := &byteSliceReader{}
		var r io.Reader = bsr
		b.SetBytes(4 + 8)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bsr.remain = big[30:]
			litend.Read(r, &ls.Float32)
			litend.Read(r, &ls.Float64)
		}
		b.StopTimer()
		want := s
		want.Int8 = 0
		want.Int16 = 0
		want.Int32 = 0
		want.Int64 = 0
		want.Uint8 = 0
		want.Uint16 = 0
		want.Uint32 = 0
		want.Uint64 = 0
		want.Complex64 = 0
		want.Complex128 = 0
		want.Array = [4]uint8{0, 0, 0, 0}
		want.Bool = false
		want.BoolArray = [4]bool{false, false, false, false}
		if b.N > 0 && !reflect.DeepEqual(ls, want) {
			b.Fatalf("struct doesn't match:\ngot  %v;\nwant %v", ls, want)
		}
	})
	b.Run("bigend", func(b *testing.B) {
		var ls Struct
		bsr := &byteSliceReader{}
		var r io.Reader = bsr
		b.SetBytes(4 + 8)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bsr.remain = big[30:]
			bigend.Read(r, &ls.Float32)
			bigend.Read(r, &ls.Float64)
		}
		b.StopTimer()
		want := s
		want.Int8 = 0
		want.Int16 = 0
		want.Int32 = 0
		want.Int64 = 0
		want.Uint8 = 0
		want.Uint16 = 0
		want.Uint32 = 0
		want.Uint64 = 0
		want.Complex64 = 0
		want.Complex128 = 0
		want.Array = [4]uint8{0, 0, 0, 0}
		want.Bool = false
		want.BoolArray = [4]bool{false, false, false, false}
		if b.N > 0 && !reflect.DeepEqual(ls, want) {
			b.Fatalf("struct doesn't match:\ngot  %v;\nwant %v", ls, want)
		}
	})
}

func BenchmarkWriteFloats(b *testing.B) {
	b.Run("stdlib", func(b *testing.B) {
		buf := new(bytes.Buffer)
		var w io.Writer = buf
		b.SetBytes(4 + 8)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			binary.Write(w, binary.BigEndian, s.Float32)
			binary.Write(w, binary.BigEndian, s.Float64)
		}
		b.StopTimer()
		if b.N > 0 && !bytes.Equal(buf.Bytes(), big[30:30+4+8]) {
			b.Fatalf("first half doesn't match: %x %x", buf.Bytes(), big[30:30+4+8])
		}
	})
	b.Run("litend", func(b *testing.B) {
		b.Skip()
		buf := new(bytes.Buffer)
		var w io.Writer = buf
		b.SetBytes(4 + 8)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			litend.Write(w, s.Float32)
			litend.Write(w, s.Float64)
		}
		b.StopTimer()
		if b.N > 0 && !bytes.Equal(buf.Bytes(), big[30:30+4+8]) {
			b.Fatalf("first half doesn't match: %x %x", buf.Bytes(), big[30:30+4+8])
		}
	})
	b.Run("bigend", func(b *testing.B) {
		buf := new(bytes.Buffer)
		var w io.Writer = buf
		b.SetBytes(4 + 8)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			bigend.Write(w, s.Float32)
			bigend.Write(w, s.Float64)
		}
		b.StopTimer()
		if b.N > 0 && !bytes.Equal(buf.Bytes(), big[30:30+4+8]) {
			b.Fatalf("first half doesn't match: %x %x", buf.Bytes(), big[30:30+4+8])
		}
	})
}

func BenchmarkReadSlice1000Float32s(b *testing.B) {
	b.Run("stdlib", func(b *testing.B) {
		bsr := &byteSliceReader{}
		slice := make([]float32, 1000)
		buf := make([]byte, len(slice)*4)
		b.SetBytes(int64(len(buf)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bsr.remain = buf
			binary.Read(bsr, binary.BigEndian, slice)
		}
	})
	b.Run("litend", func(b *testing.B) {
		bsr := &byteSliceReader{}
		slice := make([]float32, 1000)
		buf := make([]byte, len(slice)*4)
		b.SetBytes(int64(len(buf)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bsr.remain = buf
			litend.Read(bsr, slice)
		}
	})
	b.Run("bigend", func(b *testing.B) {
		bsr := &byteSliceReader{}
		slice := make([]float32, 1000)
		buf := make([]byte, len(slice)*4)
		b.SetBytes(int64(len(buf)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bsr.remain = buf
			bigend.Read(bsr, slice)
		}
	})
}

func BenchmarkWriteSlice1000Float32s(b *testing.B) {
	b.Run("stdlib", func(b *testing.B) {
		slice := make([]float32, 1000)
		buf := new(bytes.Buffer)
		var w io.Writer = buf
		b.SetBytes(4 * 1000)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			binary.Write(w, binary.BigEndian, slice)
		}
		b.StopTimer()
	})
	b.Run("litend", func(b *testing.B) {
		slice := make([]float32, 1000)
		buf := new(bytes.Buffer)
		var w io.Writer = buf
		b.SetBytes(4 * 1000)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			litend.Write(w, slice)
		}
		b.StopTimer()
	})
	b.Run("bigend", func(b *testing.B) {
		slice := make([]float32, 1000)
		buf := new(bytes.Buffer)
		var w io.Writer = buf
		b.SetBytes(4 * 1000)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			bigend.Write(w, slice)
		}
		b.StopTimer()
	})
}

func BenchmarkReadSlice1000Uint8s(b *testing.B) {
	b.Run("stdlib", func(b *testing.B) {
		bsr := &byteSliceReader{}
		slice := make([]uint8, 1000)
		buf := make([]byte, len(slice))
		b.SetBytes(int64(len(buf)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bsr.remain = buf
			binary.Read(bsr, binary.BigEndian, slice)
		}
	})
	b.Run("litend", func(b *testing.B) {
		bsr := &byteSliceReader{}
		slice := make([]uint8, 1000)
		buf := make([]byte, len(slice))
		b.SetBytes(int64(len(buf)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bsr.remain = buf
			litend.Read(bsr, slice)
		}
	})
	b.Run("bigend", func(b *testing.B) {
		bsr := &byteSliceReader{}
		slice := make([]uint8, 1000)
		buf := make([]byte, len(slice))
		b.SetBytes(int64(len(buf)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bsr.remain = buf
			bigend.Read(bsr, slice)
		}
	})
}

func BenchmarkWriteSlice1000Uint8s(b *testing.B) {
	b.Run("stdlib", func(b *testing.B) {
		slice := make([]uint8, 1000)
		buf := new(bytes.Buffer)
		var w io.Writer = buf
		b.SetBytes(1000)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			binary.Write(w, binary.BigEndian, slice)
		}
	})
	b.Run("litend", func(b *testing.B) {
		slice := make([]uint8, 1000)
		buf := new(bytes.Buffer)
		var w io.Writer = buf
		b.SetBytes(1000)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			litend.Write(w, slice)
		}
	})
	b.Run("bigend", func(b *testing.B) {
		slice := make([]uint8, 1000)
		buf := new(bytes.Buffer)
		var w io.Writer = buf
		b.SetBytes(1000)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			bigend.Write(w, slice)
		}
	})
}
