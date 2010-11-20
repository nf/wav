package wav

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
)

type File struct {
	SampleRate uint32
	SignificantBits uint16
	Channels uint16
	b bytes.Buffer
	done bool
}

func (f *File) Write(data []byte) (int, os.Error) {
	if f.b.Len() == 0 {
		if err := writeFmt(f); err != nil {
			return 0, err
		}
	}
	if err := writeChunk(&f.b, "data", data); err != nil {
		return 0, err
	}
	return 8+len(data), nil
}

func (f *File) WriteTo(w io.Writer) (n int64, err os.Error) {
	defer func() {
		if e, ok := recover().(os.Error); ok {
			err = e
		}
	}()
	write(w, "RIFF")
	n += 4
	write(w, uint32(f.b.Len()))
	n += 4
	write(w, "WAVE")
	n += 4
	write(w, f.b.Bytes())
	n += int64(f.b.Len())
	return
}

func writeFmt(f *File) (err os.Error) {
	defer func() {
		if e, ok := recover().(os.Error); ok {
			err = e
		}
	}()
	var b bytes.Buffer
	write(&b, uint16(1)) // uncompressed/PCM
	write(&b, f.Channels)
	write(&b, f.SampleRate)
	write(&b, f.SampleRate*uint32(f.SignificantBits)) // bytes per second
	write(&b, f.SignificantBits)
	return writeChunk(&f.b, "fmt ", b.Bytes())
}

func writeChunk(w io.Writer, id string, data []byte) (err os.Error) {
	if len(id) != 4 {
		return os.NewError("invalid chunk id")
	}
	defer func() {
		if e, ok := recover().(os.Error); ok {
			err = e
		}
	}()
	write(w, id)
	write(w, uint32(len(data)))
	write(w, data)
	return
}

func write(w io.Writer, data interface{}) {
	if b, ok := data.([]byte); ok {
		for c := 0; c < len(b); {
			n, err := w.Write(b[c:])
			if err != nil {
				panic(err)
			}
			c += n
		}
		return
	}
	if err := binary.Write(w, binary.LittleEndian, data); err != nil {
		panic(err)
	}
}
