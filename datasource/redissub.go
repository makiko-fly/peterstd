package datasource

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"gitlab.wallstcn.com/spider/peterstd/logging"
	"gitlab.wallstcn.com/spider/peterstd/util"
	"time"
)

type RedisSubConfig struct {
	RedisConfig              `yaml:"redis"`
	Channel                  string `yaml:"channel"`
	HealthCheckIntervalInSec int64  `yaml:"health_check_interval_in_sec"`
}

type RedisSubscriber struct {
	PubSubConn               redis.PubSubConn
	Handler                  RedisChannelMsgHandler
	Channel                  string
	HealthCheckIntervalInSec int64
}

type RedisChannelMsgHandler interface {
	HandleChannelMsg(channel string, data []byte) error
}

func (c RedisSubConfig) NewSubscriber() *RedisSubscriber {
	if c.HealthCheckIntervalInSec < 30 {
		c.HealthCheckIntervalInSec = 30 // should be at least 30 seconds
	}
	conn, err := redis.Dial("tcp",
		fmt.Sprintf("%s:%d", c.Host, c.Port),
		redis.DialPassword(c.Auth),
		redis.DialDatabase(c.DB),
		redis.DialReadTimeout(10*time.Second+time.Duration(c.HealthCheckIntervalInSec)*time.Second),
		redis.DialWriteTimeout(10*time.Second),
	)
	if err != nil {
		panic(err)
	}
	psc := redis.PubSubConn{Conn: conn}
	// Can't Ping before it has subscribed to some channel.
	//if err := psc.Ping(""); err != nil {
	//	panic(fmt.Sprintf("fails to ping redis server, err: %v", err))
	//}
	return &RedisSubscriber{PubSubConn: psc, Channel: c.Channel, HealthCheckIntervalInSec: c.HealthCheckIntervalInSec}
}

func (s *RedisSubscriber) ServeForever(errorChan chan<- error) {
	// this PubSub connection will never be closed by redis-server due to timeout, according to following docs:
	// https://redis.io/topics/clients, section "Client timeouts", quote: "Note that the timeout only applies to
	// normal clients and it does not apply to Pub/Sub clients, since a Pub/Sub connection is a push style connection
	// so a client that is idle is the norm."
	if err := s.PubSubConn.Subscribe(s.Channel); err != nil {
		errorChan <- util.AppendLineNumToErr(err)
		return
	}
	go func() {
		for {
			switch n := s.PubSubConn.Receive().(type) {
			case error:
				s.PubSubConn.Close()
				errorChan <- util.AppendLineNumToErr(n)
				return
			case redis.Message:
				if err := s.Handler.HandleChannelMsg(n.Channel, n.Data); err != nil {
					// no need to send to global error channel, just log it
					logging.GetLogger().Errorf("fails to handle channel msg, err: %v", err)
				}
			case redis.Subscription:
			default:
				// pass
			}
		}
	}()

	ticker := time.NewTicker(time.Duration(s.HealthCheckIntervalInSec) * time.Second)
	for {
		select {
		case <-ticker.C:
			if err := s.PubSubConn.Ping(""); err != nil {
				//errorChan <- util.AppendLineNumToErr(err)
				//return
				// only log it, no other operation
				logging.GetLogger().Errorf("fails to ping from PubSubConn, err: %v", err)
			}
			fmt.Println("RedisSubscriber ping called.")
		}
	}
}

// shouldn't be more than one handler
func (s *RedisSubscriber) AddHandler(h RedisChannelMsgHandler) {
	s.Handler = h
}
