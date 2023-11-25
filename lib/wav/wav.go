package wav

import (
	"io"
	"log"

	wavInfo "github.com/go-audio/wav"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

func Detect(inputReader io.ReadSeeker) (uint32, uint8, uint8) {
	wavInfoDecoder := wavInfo.NewDecoder(inputReader)
	wavInfoDecoder.ReadInfo()
	sampleRate := wavInfoDecoder.SampleRate
	channelCount := uint8(wavInfoDecoder.NumChans)
	bitDepth := uint8(wavInfoDecoder.BitDepth)
	return sampleRate, channelCount, bitDepth
}

func Decode(inputReader io.ReadSeeker) io.Reader {
	decodedReader, err := wav.DecodeWithoutResampling(inputReader)
	if err != nil {
		log.Fatal(err)
	}
	return decodedReader
}

func DecodeWithSampleRate(inputReader io.ReadSeeker, sampleRate int) io.Reader {
	decodedReader, err := wav.DecodeWithSampleRate(sampleRate, inputReader)
	if err != nil {
		log.Fatal(err)
	}
	return decodedReader
}
