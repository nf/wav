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
}

func (f *File) WriteData(w io.Writer, data []byte) (err os.Error) {
	defer func() {
		if e, ok := recover().(os.Error); ok {
			err = e
		}
	}()
	var buf bytes.Buffer
	writeFmt(&buf, f)
	writeChunk(&buf, "data", data)
	write(w, []byte("RIFF"))
	write(w, uint32(buf.Len()))
	write(w, []byte("WAVE"))
	write(w, buf.Bytes())
	return
}

func writeFmt(w io.Writer, f *File) (err os.Error) {
	var b bytes.Buffer
	write(&b, uint16(1)) // uncompressed/PCM
	write(&b, f.Channels)
	write(&b, f.SampleRate)
	write(&b, f.SampleRate*uint32(f.SignificantBits)) // bytes per second
	write(&b, f.SignificantBits / 8 * f.Channels) // block align
	write(&b, f.SignificantBits)
	write(&b, uint16(0)) // extra format bytes
	return writeChunk(w, "fmt ", b.Bytes())
}

func writeChunk(w io.Writer, id string, data []byte) (err os.Error) {
	if len(id) != 4 {
		panic(os.NewError("invalid chunk id"))
	}
	write(w, []byte(id))
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
