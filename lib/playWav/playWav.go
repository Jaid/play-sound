package playWav

import (
	"log"
	"os"
	"time"

	"github.com/ebitengine/oto/v3"
	wavInfo "github.com/go-audio/wav"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

func Play(fileReader *os.File) {
	wavInfoDecoder := wavInfo.NewDecoder(fileReader)
	wavInfoDecoder.ReadInfo()
	_, err := fileReader.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}
	sampleRate := int(wavInfoDecoder.SampleRate)
	channelCount := int(wavInfoDecoder.NumChans)
	// bitDepth := int(wavInfoDecoder.BitDepth)
	decoder, err := wav.DecodeWithSampleRate(sampleRate, fileReader)
	if err != nil {
		log.Fatal(err)
	}
	otoOptions := &oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: channelCount,
		Format:       oto.FormatSignedInt16LE,
		BufferSize:   8192,
	}
	otoContext, readyChan, err := oto.NewContext(otoOptions)
	if err != nil {
		log.Fatal(err)
	}
	<-readyChan
	otoPlayer := otoContext.NewPlayer(decoder)
	otoPlayer.SetVolume(1)
	otoPlayer.Play()
	for otoPlayer.IsPlaying() {
		time.Sleep(time.Millisecond)
	}
	err = otoPlayer.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = fileReader.Close()
	if err != nil {
		log.Fatal(err)
	}
}
