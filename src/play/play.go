package play

import (
	"io"
	"log"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/gopxl/beep"
	"github.com/jaid/play-sound/lib/flac"
	"github.com/jaid/play-sound/lib/wav"
	"github.com/samber/lo"
)

type FileType uint8

type PlayOptions struct {
	Volume       uint8
	FileType     FileType
	ChannelCount uint8
	SampleRate   uint32
	BitDepth     uint8
	BufferSize   time.Duration
}

const (
	FileTypeWav FileType = iota
	FileTypeMp3
	FileTypeFlac
	FileTypeOpus
)

func Main(reader io.ReadSeeker, options PlayOptions) {
	var detectedSampleRate uint32
	var detectedChannelCount uint8
	var decodedReader io.Reader
	var detectedBitDepth uint8
	if options.FileType == FileTypeWav {
		detectedSampleRate, detectedChannelCount, detectedBitDepth = wav.Detect(reader)
		_, err := reader.Seek(0, 0)
		if err != nil {
			log.Fatal(err)
		}
		if options.SampleRate > 0 {
			decodedReader = wav.DecodeWithSampleRate(reader, int(options.SampleRate))
		} else {
			decodedReader = wav.Decode(reader)
		}
	}
	if options.FileType == FileTypeFlac {
		var beepReader beep.StreamSeeker
		beepReader, detectedSampleRate, detectedChannelCount, detectedBitDepth = flac.Get(reader)
		decodedReader = &flac.StreamReader{
			BeepReader: beepReader,
		}
	}
	if options.SampleRate > 0 {
		detectedSampleRate = uint32(options.SampleRate)
	}
	if options.ChannelCount > 0 {
		detectedChannelCount = options.ChannelCount
	}
	sampleRate, _ := lo.Coalesce(options.SampleRate, detectedSampleRate)
	channelCount, _ := lo.Coalesce(options.ChannelCount, detectedChannelCount)
	bitDepth, _ := lo.Coalesce(options.BitDepth, detectedBitDepth)
	log.Println(sampleRate)
	log.Println(channelCount)
	log.Println(bitDepth)
	otoOptions := &oto.NewContextOptions{
		SampleRate:   int(sampleRate),
		ChannelCount: int(channelCount),
		Format:       oto.FormatSignedInt16LE,
		BufferSize:   options.BufferSize,
	}
	otoContext, readyChan, err := oto.NewContext(otoOptions)
	if err != nil {
		log.Fatal(err)
	}
	<-readyChan
	otoPlayer := otoContext.NewPlayer(decodedReader)
	if options.Volume != 100 {
		otoPlayer.SetVolume(float64(options.Volume) / 100)
	}
	otoPlayer.Play()
	for otoPlayer.IsPlaying() {
		time.Sleep(time.Millisecond)
	}
	err = otoPlayer.Close()
	if err != nil {
		panic(err)
	}
}
