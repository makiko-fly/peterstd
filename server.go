package peterstd

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gomodule/redigo/redis"
)

type Server interface {
	ServeForever(errChan chan<- error)
}

type ServerFactory interface {
	Build() (Server, error)
}

func WaitError(errorChan <-chan error) {
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case err := <-errorChan:
			if err != nil {
				Panic(err)
				panic(err)
			}
		case s := <-interrupt:
			fmt.Printf("\nReceive a signal: %v, exit peacefully.\n", s)
			os.Exit(0)
		}
	}
}

func RunAllServer(servers ...Server) {
	Info("========== Init success, run all servers ==========")
	errorChan := make(chan error, 42)
	for _, server := range servers {
		go server.ServeForever(errorChan)
	}
	WaitError(errorChan)
}

func RunServerWithLock(server Server, errChan chan<- error, rp *redis.Pool, name string) error {
	conn := rp.Get()
	ok, err := redis.Bool(conn.Do("SETNX", name, 1))
	if err != nil {
		return err
	}

	if ok {
		_, err := conn.Do("EXPIRE", name, 60)
		if err != nil {
			return err
		}
		go server.ServeForever(errChan)
	}
	conn.Close()
	return nil
}
