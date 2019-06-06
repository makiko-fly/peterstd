package server

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gomodule/redigo/redis"
)

// Server
type Server interface {
	ServeForever(errChan chan<- error)
}

// Server工厂
type ServerFactory interface {
	Build() Server
}

var _globalErrChan = make(chan error, 1)

// RunServer 运行所有server
func RunServer(servers ...Server) {
	for _, _server := range servers {
		go _server.ServeForever(_globalErrChan)
	}
}

func RunUniqueServer(server Server, rp *redis.Pool, name string) error {
	conn := rp.Get()
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("SETNX", name, 1))
	if err != nil {
		return err
	}

	if ok {
		_, err := conn.Do("EXPIRE", name, 60)
		if err != nil {
			return err
		}
		go server.ServeForever(_globalErrChan)
	}
	return nil
}

// Wait 监控channel直到检测到error或者接收到interrupt信号
func Wait() {
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case err := <-_globalErrChan:
			if err != nil {
				panic(err)
			}
		case s := <-interrupt:
			fmt.Printf("\nReceive a signal: %v, exit peacefully.\n", s)
			os.Exit(0)
		}
	}
}
