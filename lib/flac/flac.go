package flac

import (
	"io"
	"log"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/flac"
)

func Get(inputReader io.ReadSeeker) (beep.StreamSeekCloser, uint32, uint8, uint8) {
	decoder, format, err := flac.Decode(inputReader)
	if err != nil {
		log.Fatal(err)
	}
	return decoder, uint32(format.SampleRate), uint8(format.NumChannels), uint8(format.Precision * 8)
}

type StreamReader struct {
	BeepReader beep.StreamSeeker
}

// from Copilot
func (sr *StreamReader) Read(p []byte) (n int, err error) {
	// Convert p to a format suitable for sr.stream.Stream
	frames := len(p) / 4
	buf := make([][2]float64, frames)
	// Stream from sr.stream into buf
	n, ok := sr.BeepReader.Stream(buf)
	if !ok {
		err = io.EOF
	}
	// Convert buf back into bytes and copy into p
	for i, frame := range buf[:n] {
		for c, sample := range frame {
			u := uint32((sample + 1.0) / 2.0 * 4294967295.0)
			for j := 0; j < 4; j++ {
				p[i*8+c*4+j] = byte(u >> (j * 8))
			}
		}
	}
	return n * 8, err
}
