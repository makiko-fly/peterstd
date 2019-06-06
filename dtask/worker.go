package dtask

import "github.com/RichardKnop/machinery/v1"

type Worker struct {
	server *machinery.Server
	Name   string
}

func (w *Worker) Register(name string, task interface{}) error {
	return w.server.RegisterTask(name, task)
}

func (w *Worker) ServeForever(errorChan chan<- error) {
	if cap(errorChan) == 0 {
		panic("Capacity of error channel should > 0")
	}
	w.server.NewWorker(w.Name, 0).LaunchAsync(errorChan)
}
