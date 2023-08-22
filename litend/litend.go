package litend

import (
	"errors"
	"io"
	"math"
	"reflect"
	"sync"
)

func Uint16(b []byte) uint16 {
	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
	return uint16(b[0]) | uint16(b[1])<<8
}

func PutUint16(b []byte, v uint16) {
	_ = b[1] // early bounds check to guarantee safety of writes below
	b[0] = byte(v)
	b[1] = byte(v >> 8)
}

func AppendUint16(b []byte, v uint16) []byte {
	return append(b,
		byte(v),
		byte(v>>8),
	)
}

func Uint32(b []byte) uint32 {
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func PutUint32(b []byte, v uint32) {
	_ = b[3] // early bounds check to guarantee safety of writes below
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
}

func AppendUint32(b []byte, v uint32) []byte {
	return append(b,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
	)
}

func Uint64(b []byte) uint64 {
	_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func PutUint64(b []byte, v uint64) {
	_ = b[7] // early bounds check to guarantee safety of writes below
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
}

func AppendUint64(b []byte, v uint64) []byte {
	return append(b,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
		byte(v>>32),
		byte(v>>40),
		byte(v>>48),
		byte(v>>56),
	)
}

func Read(r io.Reader, data any) error {
	// Fast path for basic types and slices.
	if n := intSizeOf(data); n != 0 {
		bs := make([]byte, n)
		if _, err := io.ReadFull(r, bs); err != nil {
			return err
		}
		switch data := data.(type) {
		case *bool:
			*data = bs[0] != 0
		case *int8:
			*data = int8(bs[0])
		case *uint8:
			*data = bs[0]
		case *int16:
			*data = int16(Uint16(bs))
		case *uint16:
			*data = Uint16(bs)
		case *int32:
			*data = int32(Uint32(bs))
		case *uint32:
			*data = Uint32(bs)
		case *int64:
			*data = int64(Uint64(bs))
		case *uint64:
			*data = Uint64(bs)
		case *float32:
			*data = math.Float32frombits(Uint32(bs))
		case *float64:
			*data = math.Float64frombits(Uint64(bs))
		case []bool:
			for i, x := range bs { // Easier to loop over the input for 8-bit values.
				data[i] = x != 0
			}
		case []int8:
			for i, x := range bs {
				data[i] = int8(x)
			}
		case []uint8:
			copy(data, bs)
		case []int16:
			for i := range data {
				data[i] = int16(Uint16(bs[2*i:]))
			}
		case []uint16:
			for i := range data {
				data[i] = Uint16(bs[2*i:])
			}
		case []int32:
			for i := range data {
				data[i] = int32(Uint32(bs[4*i:]))
			}
		case []uint32:
			for i := range data {
				data[i] = Uint32(bs[4*i:])
			}
		case []int64:
			for i := range data {
				data[i] = int64(Uint64(bs[8*i:]))
			}
		case []uint64:
			for i := range data {
				data[i] = Uint64(bs[8*i:])
			}
		case []float32:
			for i := range data {
				data[i] = math.Float32frombits(Uint32(bs[4*i:]))
			}
		case []float64:
			for i := range data {
				data[i] = math.Float64frombits(Uint64(bs[8*i:]))
			}
		default:
			n = 0 // fast path doesn't apply
		}
		if n != 0 {
			return nil
		}
	}

	// Fallback to reflect-based decoding.
	v := reflect.ValueOf(data)
	size := -1
	switch v.Kind() {
	case reflect.Pointer:
		v = v.Elem()
		size = SizeOf(v)
	case reflect.Slice:
		size = SizeOf(v)
	}
	if size < 0 {
		return errors.New("binary.Read: invalid type " + reflect.TypeOf(data).String())
	}
	d := &decoder{buf: make([]byte, size)}
	if _, err := io.ReadFull(r, d.buf); err != nil {
		return err
	}
	d.value(v)
	return nil
}

func Write(w io.Writer, data any) error {
	// Fast path for basic types and slices.
	if n := intSizeOf(data); n != 0 {
		bs := make([]byte, n)
		switch v := data.(type) {
		case *bool:
			if *v {
				bs[0] = 1
			} else {
				bs[0] = 0
			}
		case bool:
			if v {
				bs[0] = 1
			} else {
				bs[0] = 0
			}
		case []bool:
			for i, x := range v {
				if x {
					bs[i] = 1
				} else {
					bs[i] = 0
				}
			}
		case *int8:
			bs[0] = byte(*v)
		case int8:
			bs[0] = byte(v)
		case []int8:
			for i, x := range v {
				bs[i] = byte(x)
			}
		case *uint8:
			bs[0] = *v
		case uint8:
			bs[0] = v
		case []uint8:
			bs = v
		case *int16:
			PutUint16(bs, uint16(*v))
		case int16:
			PutUint16(bs, uint16(v))
		case []int16:
			for i, x := range v {
				PutUint16(bs[2*i:], uint16(x))
			}
		case *uint16:
			PutUint16(bs, *v)
		case uint16:
			PutUint16(bs, v)
		case []uint16:
			for i, x := range v {
				PutUint16(bs[2*i:], x)
			}
		case *int32:
			PutUint32(bs, uint32(*v))
		case int32:
			PutUint32(bs, uint32(v))
		case []int32:
			for i, x := range v {
				PutUint32(bs[4*i:], uint32(x))
			}
		case *uint32:
			PutUint32(bs, *v)
		case uint32:
			PutUint32(bs, v)
		case []uint32:
			for i, x := range v {
				PutUint32(bs[4*i:], x)
			}
		case *int64:
			PutUint64(bs, uint64(*v))
		case int64:
			PutUint64(bs, uint64(v))
		case []int64:
			for i, x := range v {
				PutUint64(bs[8*i:], uint64(x))
			}
		case *uint64:
			PutUint64(bs, *v)
		case uint64:
			PutUint64(bs, v)
		case []uint64:
			for i, x := range v {
				PutUint64(bs[8*i:], x)
			}
		case *float32:
			PutUint32(bs, math.Float32bits(*v))
		case float32:
			PutUint32(bs, math.Float32bits(v))
		case []float32:
			for i, x := range v {
				PutUint32(bs[4*i:], math.Float32bits(x))
			}
		case *float64:
			PutUint64(bs, math.Float64bits(*v))
		case float64:
			PutUint64(bs, math.Float64bits(v))
		case []float64:
			for i, x := range v {
				PutUint64(bs[8*i:], math.Float64bits(x))
			}
		}
		_, err := w.Write(bs)
		return err
	}

	// Fallback to reflect-based encoding.
	v := reflect.Indirect(reflect.ValueOf(data))
	size := SizeOf(v)
	if size < 0 {
		return errors.New("binary.Write: some values are not fixed-sized in type " + reflect.TypeOf(data).String())
	}
	buf := make([]byte, size)
	e := &encoder{buf: buf}
	e.value(v)
	_, err := w.Write(buf)
	return err
}

func Size(v any) int {
	return SizeOf(reflect.Indirect(reflect.ValueOf(v)))
}

var structSize sync.Map // map[reflect.Type]int

func SizeOf(v reflect.Value) int {
	switch v.Kind() {
	case reflect.Slice:
		if s := sizeof(v.Type().Elem()); s >= 0 {
			return s * v.Len()
		}

	case reflect.Struct:
		t := v.Type()
		if size, ok := structSize.Load(t); ok {
			return size.(int)
		}
		size := sizeof(t)
		structSize.Store(t, size)
		return size

	default:
		if v.IsValid() {
			return sizeof(v.Type())
		}
	}

	return -1
}

// sizeof returns the size >= 0 of variables for the given type or -1 if the type is not acceptable.
func sizeof(t reflect.Type) int {
	switch t.Kind() {
	case reflect.Array:
		if s := sizeof(t.Elem()); s >= 0 {
			return s * t.Len()
		}

	case reflect.Struct:
		sum := 0
		for i, n := 0, t.NumField(); i < n; i++ {
			s := sizeof(t.Field(i).Type)
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum

	case reflect.Bool,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return int(t.Size())
	}

	return -1
}

type coder struct {
	buf    []byte
	offset int
}

type (
	decoder coder
	encoder coder
)

func (d *decoder) bool() bool {
	x := d.buf[d.offset]
	d.offset++
	return x != 0
}

func (e *encoder) bool(x bool) {
	if x {
		e.buf[e.offset] = 1
	} else {
		e.buf[e.offset] = 0
	}
	e.offset++
}

func (d *decoder) uint8() uint8 {
	x := d.buf[d.offset]
	d.offset++
	return x
}

func (e *encoder) uint8(x uint8) {
	e.buf[e.offset] = x
	e.offset++
}

func (d *decoder) uint16() uint16 {
	x := Uint16(d.buf[d.offset : d.offset+2])
	d.offset += 2
	return x
}

func (e *encoder) uint16(x uint16) {
	PutUint16(e.buf[e.offset:e.offset+2], x)
	e.offset += 2
}

func (d *decoder) uint32() uint32 {
	x := Uint32(d.buf[d.offset : d.offset+4])
	d.offset += 4
	return x
}

func (e *encoder) uint32(x uint32) {
	PutUint32(e.buf[e.offset:e.offset+4], x)
	e.offset += 4
}

func (d *decoder) uint64() uint64 {
	x := Uint64(d.buf[d.offset : d.offset+8])
	d.offset += 8
	return x
}

func (e *encoder) uint64(x uint64) {
	PutUint64(e.buf[e.offset:e.offset+8], x)
	e.offset += 8
}

func (d *decoder) int8() int8 { return int8(d.uint8()) }

func (e *encoder) int8(x int8) { e.uint8(uint8(x)) }

func (d *decoder) int16() int16 { return int16(d.uint16()) }

func (e *encoder) int16(x int16) { e.uint16(uint16(x)) }

func (d *decoder) int32() int32 { return int32(d.uint32()) }

func (e *encoder) int32(x int32) { e.uint32(uint32(x)) }

func (d *decoder) int64() int64 { return int64(d.uint64()) }

func (e *encoder) int64(x int64) { e.uint64(uint64(x)) }

func (d *decoder) value(v reflect.Value) {
	switch v.Kind() {
	case reflect.Array:
		l := v.Len()
		for i := 0; i < l; i++ {
			d.value(v.Index(i))
		}

	case reflect.Struct:
		t := v.Type()
		l := v.NumField()
		for i := 0; i < l; i++ {
			// Note: Calling v.CanSet() below is an optimization.
			// It would be sufficient to check the field name,
			// but creating the StructField info for each field is
			// costly (run "go test -bench=ReadStruct" and compare
			// results when making changes to this code).
			if v := v.Field(i); v.CanSet() || t.Field(i).Name != "_" {
				d.value(v)
			} else {
				d.skip(v)
			}
		}

	case reflect.Slice:
		l := v.Len()
		for i := 0; i < l; i++ {
			d.value(v.Index(i))
		}

	case reflect.Bool:
		v.SetBool(d.bool())

	case reflect.Int8:
		v.SetInt(int64(d.int8()))
	case reflect.Int16:
		v.SetInt(int64(d.int16()))
	case reflect.Int32:
		v.SetInt(int64(d.int32()))
	case reflect.Int64:
		v.SetInt(d.int64())

	case reflect.Uint8:
		v.SetUint(uint64(d.uint8()))
	case reflect.Uint16:
		v.SetUint(uint64(d.uint16()))
	case reflect.Uint32:
		v.SetUint(uint64(d.uint32()))
	case reflect.Uint64:
		v.SetUint(d.uint64())

	case reflect.Float32:
		v.SetFloat(float64(math.Float32frombits(d.uint32())))
	case reflect.Float64:
		v.SetFloat(math.Float64frombits(d.uint64()))

	case reflect.Complex64:
		v.SetComplex(complex(
			float64(math.Float32frombits(d.uint32())),
			float64(math.Float32frombits(d.uint32())),
		))
	case reflect.Complex128:
		v.SetComplex(complex(
			math.Float64frombits(d.uint64()),
			math.Float64frombits(d.uint64()),
		))
	}
}

func (e *encoder) value(v reflect.Value) {
	switch v.Kind() {
	case reflect.Array:
		l := v.Len()
		for i := 0; i < l; i++ {
			e.value(v.Index(i))
		}

	case reflect.Struct:
		t := v.Type()
		l := v.NumField()
		for i := 0; i < l; i++ {
			// see comment for corresponding code in decoder.value()
			if v := v.Field(i); v.CanSet() || t.Field(i).Name != "_" {
				e.value(v)
			} else {
				e.skip(v)
			}
		}

	case reflect.Slice:
		l := v.Len()
		for i := 0; i < l; i++ {
			e.value(v.Index(i))
		}

	case reflect.Bool:
		e.bool(v.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v.Type().Kind() {
		case reflect.Int8:
			e.int8(int8(v.Int()))
		case reflect.Int16:
			e.int16(int16(v.Int()))
		case reflect.Int32:
			e.int32(int32(v.Int()))
		case reflect.Int64:
			e.int64(v.Int())
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch v.Type().Kind() {
		case reflect.Uint8:
			e.uint8(uint8(v.Uint()))
		case reflect.Uint16:
			e.uint16(uint16(v.Uint()))
		case reflect.Uint32:
			e.uint32(uint32(v.Uint()))
		case reflect.Uint64:
			e.uint64(v.Uint())
		}

	case reflect.Float32, reflect.Float64:
		switch v.Type().Kind() {
		case reflect.Float32:
			e.uint32(math.Float32bits(float32(v.Float())))
		case reflect.Float64:
			e.uint64(math.Float64bits(v.Float()))
		}

	case reflect.Complex64, reflect.Complex128:
		switch v.Type().Kind() {
		case reflect.Complex64:
			x := v.Complex()
			e.uint32(math.Float32bits(float32(real(x))))
			e.uint32(math.Float32bits(float32(imag(x))))
		case reflect.Complex128:
			x := v.Complex()
			e.uint64(math.Float64bits(real(x)))
			e.uint64(math.Float64bits(imag(x)))
		}
	}
}

func (d *decoder) skip(v reflect.Value) {
	d.offset += SizeOf(v)
}

func (e *encoder) skip(v reflect.Value) {
	n := SizeOf(v)
	zero := e.buf[e.offset : e.offset+n]
	for i := range zero {
		zero[i] = 0
	}
	e.offset += n
}

// intSizeOf returns the size of the data required to represent the data when encoded.
// It returns zero if the type cannot be implemented by the fast path in Read or Write.
func intSizeOf(data any) int {
	switch data := data.(type) {
	case bool, int8, uint8, *bool, *int8, *uint8:
		return 1
	case []bool:
		return len(data)
	case []int8:
		return len(data)
	case []uint8:
		return len(data)
	case int16, uint16, *int16, *uint16:
		return 2
	case []int16:
		return 2 * len(data)
	case []uint16:
		return 2 * len(data)
	case int32, uint32, *int32, *uint32:
		return 4
	case []int32:
		return 4 * len(data)
	case []uint32:
		return 4 * len(data)
	case int64, uint64, *int64, *uint64:
		return 8
	case []int64:
		return 8 * len(data)
	case []uint64:
		return 8 * len(data)
	case float32, *float32:
		return 4
	case float64, *float64:
		return 8
	case []float32:
		return 4 * len(data)
	case []float64:
		return 8 * len(data)
	}
	return 0
}
