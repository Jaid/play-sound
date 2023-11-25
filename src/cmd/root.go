package cmd

import (
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/spf13/afero"

	"github.com/jaid/play-sound/src/embed"
	"github.com/jaid/play-sound/src/play"
	"github.com/spf13/cobra"
)

var (
	verbose  bool
	fileType string
	volume   uint8
	samples  uint32
	channels uint8
	bits     uint8
)

func setupCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "play-sound",
		Short: "Play sound files",
		Run:   runCommand,
	}
	command.PersistentFlags().BoolVar(&verbose, "verbose", false, "Verbose output")
	command.PersistentFlags().StringVar(&fileType, "file-type", "", "File type")
	command.PersistentFlags().Uint8Var(&volume, "volume", 100, `Volume from 0 to 100`)
	command.PersistentFlags().Uint32Var(&samples, "samples", 0, `Sample rate`)
	command.PersistentFlags().Uint8Var(&channels, "channels", 0, `Channel count`)
	command.PersistentFlags().Uint8Var(&bits, "bits", 0, `Bit depth`)
	return command
}

func Main() {
	command := setupCommand()
	err := command.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func getFileTypeFromExtension(extension string) play.FileType {
	extensionNormalized := strings.ToLower(extension)
	switch extensionNormalized {
	case "wav":
		return play.FileTypeWav
	case "mp3":
		return play.FileTypeMp3
	case "flac":
		return play.FileTypeFlac
	case "opus":
		return play.FileTypeOpus
	default:
		return play.FileTypeWav
	}
}

func runCommand(cobraCommand *cobra.Command, args []string) {
	if verbose {
		log.Println("Verbose: ", verbose)
		log.Println("Type: ", fileType)
		log.Println("Volume: ", volume)
	}
	var fileAfero afero.File
	var options play.PlayOptions
	options.Volume = volume
	if len(args) == 0 {
		var err error
		fs := afero.FromIOFS{
			FS: embed.Fs,
		}
		fileAfero, err = fs.Open("files/sound.wav")
		if err != nil {
			log.Fatal(err)
		}
		options.FileType = play.FileTypeWav
		options.BitDepth = 16
		options.ChannelCount = 1
		options.SampleRate = 48000
		options.BufferSize = time.Millisecond * 10
	} else {
		file := args[0]
		fs := afero.NewOsFs()
		var err error
		fileAfero, err = fs.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(fileType)
		if lo.IsEmpty(fileType) {
			extension := filepath.Ext(fileAfero.Name())
			log.Println(extension)
			options.FileType = getFileTypeFromExtension(strings.TrimPrefix(extension, "."))
		} else {
			options.FileType = getFileTypeFromExtension(fileType)
		}
	}
	if lo.IsNotEmpty(samples) {
		options.SampleRate = samples
	}
	if lo.IsNotEmpty(channels) {
		options.ChannelCount = channels
	}
	if lo.IsNotEmpty(bits) {
		options.BitDepth = bits
	}
	log.Println(fileAfero.Name())
	play.Main(fileAfero, options)
	err := fileAfero.Close()
	if err != nil {
		log.Fatal(err)
	}
}
