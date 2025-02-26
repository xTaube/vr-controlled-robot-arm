package main

import (
	"log"
	"os"

	"github.com/xTaube/vr-controlled-robot-arm/robot"
	"github.com/xTaube/vr-controlled-robot-arm/server"
	"github.com/xTaube/vr-controlled-robot-arm/video"
	"go.bug.st/serial"
)

func main() {
	log.Println("Initializing robot arm...")
	robot, err := robot.InitRobot(
		robot.UartConfig{
			PortName: os.Getenv("UART_PORT"),
			Parity:   serial.EvenParity,
			StopBits: serial.OneStopBit,
			BaudRate: 115200,
			DataBits: 8,
		},
	)
	if err != nil {
		log.Printf("Error initializing robot arm: %s.\n", err)
		return
	}
	defer robot.ShutDown()
	log.Println("Robot arm initialized.")

	log.Println("Initializing camera 0 ...")
	video0 := video.InitVideoStream(
		os.Getenv("CAMERA_DEVICE0_PATH"),
		video.Resoulution{Width: 1280, Height: 720},
		video.FPS30,
		video.MJPEG,
		"rtsp://localhost:8554/video/feed0",
	)
	defer video0.Stop()
	log.Println("Camera 0 initialized.")

	log.Println("Initializing camera 1 ...")
	video1 := video.InitVideoStream(
		os.Getenv("CAMERA_DEVICE1_PATH"),
		video.Resoulution{Width: 1280, Height: 720},
		video.FPS30,
		video.MJPEG,
		"rtsp://localhost:8554/video/feed1",
	)
	defer video1.Stop()
	log.Println("Camera 1 initialized.")

	err = server.RunWebSocketServer(os.Getenv("PORT"), robot, video0, video1)
	if err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
