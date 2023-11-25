package stubs

var Publish = "Broker.Publish"
var Subscribe = "Broker.Subscribe"
var Terminate = "Broker.TerminateWorker"
var Pause = "Broker.Pause"
var PublishCurrentWorld = "Broker.CurrentWorld"

var ProcessTurns = "GOLWorker.ProcessTurns"
var TerminateHandler = "GOLWorker.Terminate"

type Params struct {
	World       [][]byte
	Turns       int
	ImageWidth  int
	ImageHeight int
}

type WorkerParams struct {
	World  [][]byte
	Turns  int
	StartY int
	EndY   int
	StartX int
	EndX   int
}

type WorkerResponse struct {
	World [][]byte
}

type StatusReport struct {
	Message string
}

type Subscription struct {
	Address string
	Params  Params
}

type Ticker struct{}

type CurrentWorld struct {
	World [][]byte
	Turn  int
}
