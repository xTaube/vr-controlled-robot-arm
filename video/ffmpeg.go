package video

import (
	"fmt"
	"os/exec"
)

type StreamFormat string

const (
	V4L2 StreamFormat = "v4l2"
)

type InputFormat string

const (
	H264  InputFormat = "h264"
	MJPEG             = "mjpeg"
)

type FFMPEGCommandBuilder struct {
	command string
}

func (builder *FFMPEGCommandBuilder) SetStreamFormat(format StreamFormat) *FFMPEGCommandBuilder {
	builder.command += fmt.Sprintf(" -f %s", format)
	return builder
}

func (builder *FFMPEGCommandBuilder) SetInputFormat(format InputFormat) *FFMPEGCommandBuilder {
	builder.command += fmt.Sprintf(" -input_format %s", format)
	return builder
}

func (builder *FFMPEGCommandBuilder) SetVideoSize(width int, height int) *FFMPEGCommandBuilder {
	builder.command += fmt.Sprintf(" -video_size %dx%d", width, height)
	return builder
}

func (builder *FFMPEGCommandBuilder) SetFramerate(framerate int) *FFMPEGCommandBuilder {
	builder.command += fmt.Sprintf(" -framerate %d", framerate)
	return builder
}

func (builder *FFMPEGCommandBuilder) SetDevice(devicePath string) *FFMPEGCommandBuilder {
	builder.command += fmt.Sprintf(" -i %s", devicePath)
	return builder
}

func (builder *FFMPEGCommandBuilder) SetRTSPOutput(serverAddress string) *FFMPEGCommandBuilder {
	builder.command += fmt.Sprintf(" -c copy -f rtsp %s", serverAddress)
	return builder
}

func (builder *FFMPEGCommandBuilder) SetStreamLoop() *FFMPEGCommandBuilder {
	builder.command += " -stream_loop -1"
	return builder
}

func (builder *FFMPEGCommandBuilder) SetRE() *FFMPEGCommandBuilder {
	builder.command += " -re"
	return builder
}

func (builder *FFMPEGCommandBuilder) Execute() error {
	cmd := exec.Command(builder.command)
	return cmd.Run()
}

func InitFFMPEGCommandBuilder() *FFMPEGCommandBuilder {
	return &FFMPEGCommandBuilder{command: "ffmpeg"}
}

func StopFFMPEGVideoStreaming() error {
	cmd := exec.Command("pkill", "-f ffmpeg")
	return cmd.Run()
}
