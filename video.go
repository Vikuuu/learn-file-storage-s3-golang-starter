package main

import (
	"bytes"
	"encoding/json"
	"os/exec"
)

type videoMetadata struct {
	Streams []struct {
		Width  int    `json:"width"`
		Height int    `json:"height"`
		DAS    string `json:"display_aspect_ratio"`
	} `json:"streams"`
}

func getVideoAspectRatio(filepath string) (string, error) {
	cmd := exec.Command(
		"ffprobe",
		"-v",
		"error",
		"-print_format",
		"json",
		"-show_streams",
		filepath,
	)

	var buf bytes.Buffer
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return "", err
	}

	var output videoMetadata
	err := json.Unmarshal(buf.Bytes(), &output)
	if err != nil {
		return "", err
	}

	if output.Streams[0].DAS == "16:9" {
		return "landscape", nil
	} else if output.Streams[0].DAS == "9:16" {
		return "portrait", nil
	}
	return "other", nil
}

func processVideoForFastStart(filepath string) (string, error) {
	outputFile := filepath + ".processing"
	cmd := exec.Command(
		"ffmpeg",
		"-i",
		filepath,
		"-c",
		"copy",
		"-movflags",
		"faststart",
		"-f",
		"mp4",
		outputFile,
	)

	if err := cmd.Run(); err != nil {
		return "", err
	}
	return outputFile, nil
}
