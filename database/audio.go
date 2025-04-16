package database

import (
	"encoding/json"
	"fmt"
	"os"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type FFProbeOutput struct {
	Format struct {
		FormatName string `json:"format_name"`
		Duration   string `json:"duration"`
		BitRate    string `json:"bit_rate"`
	} `json:"format"`
	Streams []struct {
		CodecType  string `json:"codec_type"`
		CodecName  string `json:"codec_name"`
		SampleRate string `json:"sample_rate,omitempty"`
		Channels   int    `json:"channels,omitempty"`
	} `json:"streams"`
}

func getDetailedAudioInfo(filePath string) (map[string]interface{}, error) {
	// Run ffprobe to get file information in JSON format
	data, err := ffmpeg.Probe(filePath)
	if err != nil {
		return nil, fmt.Errorf("ffprobe failed: %w", err)
	}

	// Parse the JSON output
	var result FFProbeOutput
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, fmt.Errorf("failed to parse ffprobe output: %w", err)
	}

	// Extract relevant audio information
	info := make(map[string]interface{})
	info["format"] = result.Format.FormatName
	info["duration"] = result.Format.Duration
	info["bitrate"] = result.Format.BitRate

	// Find audio stream
	for _, stream := range result.Streams {
		if stream.CodecType == "audio" {
			info["codec"] = stream.CodecName
			info["sample_rate"] = stream.SampleRate
			info["channels"] = stream.Channels
			break
		}
	}

	return info, nil
}

func ConvertAudioFileToWav(inputFile string, outputFile string) error {
	// See if this is the right format already.
	audio_info, err := getDetailedAudioInfo(inputFile)
	fmt.Printf("[%v] Audio Info: %v\n", inputFile, audio_info)
	if err != nil {
		return fmt.Errorf("failed to get audio info: %w", err)
	}
	if audio_info["format"] == "wav" && audio_info["sameple_rate"] == "44100" && audio_info["channels"] == 2 {
		// This is already a WAV file, so we can skip the conversion.
		fmt.Printf("Audio file is already in WAV format: %v\n", inputFile)
		return nil
	}
	// Run ffmpeg to convert the audio file to WAV format
	err = ffmpeg.Input(inputFile).Output(outputFile, ffmpeg.KwArgs{
		"ar":  "44100",     // Set sample rate to 44100 Hz
		"ac":  "2",         // Set number of channels to 2 (stereo)
		"c:a": "pcm_f32le", // Set output format to WAV floating-point 32-bit little-endian
	}).OverWriteOutput().Run()

	if err != nil {
		return fmt.Errorf("audio conversion failed: %w", err)
	}
	return nil
}

func InstallWav(audioFile string, mem *MemorySlot, looperDataDir string, clearExisting bool) error {
	wavdir, wavfilename := mem.GetWavFileLocation(looperDataDir)

	// wavdir should exist.
	info, err := os.Stat(wavdir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("WAV output directory does not exist: %v", wavdir)
		}
		return fmt.Errorf("failed to check WAV output directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("WAV output directory is not a directory: %v", wavdir)
	}

	fmt.Printf("Converting: [%v] -> [%v]\n", audioFile, wavdir+"/"+wavfilename)

	// Convert to the format we want.
	err = ConvertAudioFileToWav(audioFile, wavdir+"/"+wavfilename)
	if err != nil {
		return fmt.Errorf("failed to convert audio file: %w", err)
	}

	if clearExisting {
		// Remove all files in wavdir aside from our new one.
		files, err := os.ReadDir(wavdir)
		if err != nil {
			return fmt.Errorf("failed to read WAV output directory: %w", err)
		}
		for _, file := range files {
			if file.Name() != wavfilename {
				err := os.Remove(wavdir + "/" + file.Name())
				if err != nil {
					return fmt.Errorf("failed to remove existing WAV file: %w", err)
				}
				fmt.Printf("Removed existing WAV file: %v\n", wavdir+"/"+file.Name())
			}
		}
	}

	return nil
}
