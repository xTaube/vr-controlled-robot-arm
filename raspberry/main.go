package main

import (
	"log"
	"os"

	"github.com/xTaube/vr-controlled-robot-arm/robot"
	"github.com/xTaube/vr-controlled-robot-arm/server"
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

	err = server.RunWebSocketServer(os.Getenv("PORT"), robot)
	if err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
