package video

import (
	"os"
	"os/exec"
)


func StartFFMPEGVideoStreaming() error {
	cmd := exec.Command("ffmpeg", "-f", "v4l2", "-framerate", "30", "re", "-stream_loop", "-1", "-video_size", "1920x1080", "-input_format", "mjpeg", "-i", os.Getenv("CAMERA_DEVICE_PATH"), "-c", "copy", "-f", "rtsp", os.Getenv("SERVER_ADDRESS"))
	return cmd.Run()
}

func StopFFMPEGVideoStreaming() error {
	cmd := exec.Command("pkill", "-f ffmpeg")
	return cmd.Run()
}
