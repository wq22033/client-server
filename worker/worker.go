package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"os"

	"uk.ac.bris.cs/gameoflife/stubs"
)

var worldChannel = make(chan stubs.Params)
var broker *rpc.Client

func getOutboundIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr).IP.String()
	return localAddr
}

func buildWorld(width int, height int) [][]byte {
	world := make([][]byte, height)
	for i := range world {
		world[i] = make([]byte, width)
	}
	return world
}

func golWorker(world [][]byte, turn, startY, endY, startX, endX int) [][]byte {
	// Create a new 2D slice to store the updated world.
	newWorld := buildWorld(endX, endY-startY)
	// Iterate over each cell in the world.
	for y := startY; y < endY; y++ {
		for x := startX; x < endX; x++ {
			// Count the number of alive neighbors around the current cell.
			aliveNeighbors := 0
			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					// Skip the current cell itself.
					if i == 0 && j == 0 {
						continue
					}
					// Calculate the coordinates of the neighboring cell.
					neighborX := x + i
					neighborY := y + j

					// Wrap around the world if necessary.
					if neighborX < 0 {
						neighborX = endX - 1
					} else if neighborX >= endX {
						neighborX = 0
					}
					if neighborY < 0 {
						neighborY = endX - 1
					} else if neighborY >= endX {
						neighborY = 0
					}

					// Check if the neighboring cell is alive.
					if world[neighborY][neighborX] == 255 {
						aliveNeighbors++
					}
				}
			}

			// Update the newWorld slice based on the Game of Life rules.
			if world[y][x] == 255 {
				if aliveNeighbors < 2 || aliveNeighbors > 3 {
					newWorld[y-startY][x] = 0

				} else {
					newWorld[y-startY][x] = 255

				}
			} else {
				if aliveNeighbors == 3 {
					newWorld[y-startY][x] = 255
				} else {
					newWorld[y-startY][x] = 0

				}
			}
		}
	}
	return newWorld
}

func handleTermination() {
	fmt.Println("Received termination signal. Cleaning up and exiting worker...")
	os.Exit(0)
}

type GOLWorker struct{}

func (w *GOLWorker) ProcessTurns(req stubs.WorkerParams, res *stubs.WorkerResponse) (err error) {
	res.World = golWorker(req.World, req.Turns, req.StartY, req.EndY, req.StartX, req.EndX)
	return
}

func (w *GOLWorker) Terminate(req bool, res *bool) error {
	*res = true
	handleTermination()
	return nil
}

func main() {
	pAddr := flag.String("port", "8050", "Port to listen on")
	pIP := flag.String("ip", getOutboundIP(), "IP address of worker")
	brokerAddr := flag.String("broker", "127.0.0.1:8030", "Address of broker instance")
	flag.Parse()
	broker, _ = rpc.Dial("tcp", *brokerAddr)
	status := new(stubs.StatusReport)
	address := *pIP + ":" + *pAddr
	fmt.Println(address)

	rpc.Register(&GOLWorker{})
	listener, err := net.Listen("tcp", ":"+*pAddr)
	if err != nil {
		panic(err)
	}
	broker.Call(stubs.Subscribe, stubs.Subscription{Address: address}, status)
	defer listener.Close()
	rpc.Accept(listener)
	flag.Parse()
	// pAddr := flag.String("port", "8030", "Port to listen on")
	// flag.Parse()
	// rpc.Register(&GOLWorker{})
	// listener, _ := net.Listen("tcp", ":"+*pAddr)
	// defer listener.Close()
	// rpc.Accept(listener)

}
