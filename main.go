package main

import (
	"log"
	"os"

	"github.com/jaid/play-sound/lib/playWav"
)

func main() {
	var file string
	if len(os.Args) < 2 {
		file = "private/sound.wav"
	} else {
		file = os.Args[1]
	}
	fileInstance, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	playWav.Play(fileInstance)
}
