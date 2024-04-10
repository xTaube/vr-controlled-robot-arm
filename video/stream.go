package video

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type StreamOnError struct {}
func (err *StreamOnError) Error() string {
	return "Stream is already on"
}

type StreamOffError struct {}
func(err *StreamOffError) Error() string {
	return "Stream is already off"
}

type InputFormat string

const (
	MJPEG InputFormat = "mjpeg"
)

type Framerate int

const (
	FPS30 Framerate = 30
)

type Resoulution struct {
	Width int
	Height int
}

type VideoStream struct {
	device string
	resolution Resoulution
	framerate Framerate
	inputFormat InputFormat
	outputServerAddres string
	ffmpegProcess *exec.Cmd
}

func (vs *VideoStream) Start() (string, error) {
	if vs.ffmpegProcess != nil {
		return "", &StreamOnError{}
	}

	logFile, err := os.Create(fmt.Sprintf("stream-logs/%s.txt", time.Now()))
	if err != nil {
		return "", err
	}

	vs.ffmpegProcess = exec.Command(
		"ffmpeg",
		"-f", 
		"v4l2", 
		"-framerate", 
		fmt.Sprintf("%d", vs.framerate),
		"-re",
		"-stream_loop", 
		"-1", 
		"-video_size", 
		fmt.Sprintf("%dx%d", vs.resolution.Height, vs.resolution.Width),
		"-input_format", 
		string(vs.inputFormat),
		"-i",
		vs.device, 
		"-c", 
		"copy", 
		"-f", 
		"rtsp", 
		vs.outputServerAddres,
	)

	vs.ffmpegProcess.Stdout = logFile
	vs.ffmpegProcess.Stderr = logFile

	err = vs.ffmpegProcess.Start()
	if err != nil {
		return "", err
	}

	return strings.Replace(vs.outputServerAddres, "localhost", os.Getenv("RASPBERRY_ADDRESS"), 1), nil
}

func (vs *VideoStream) Stop() error {
	if vs.ffmpegProcess == nil {
		return &StreamOffError{}
	}
	
	err := vs.ffmpegProcess.Process.Kill()
	if err != nil {
		return err
	}
	vs.ffmpegProcess = nil
	return nil
}

func InitVideoStream(
	device string,
	resolution Resoulution, 
	framerate Framerate,
	inputFormat InputFormat,
	outputServerAddres string,
) *VideoStream {
	return &VideoStream{
		device: device, 
		resolution: resolution,
		framerate: framerate,
		inputFormat: inputFormat,
		outputServerAddres: outputServerAddres,
		ffmpegProcess: nil,
	}
}