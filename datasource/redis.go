package datasource

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

type RedisConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Auth string `yaml:"auth"`
	DB   int    `yaml:"db"`
}

//func (c RedisConfig) NewRedisClient() *goredis.Client {
//	address := fmt.Sprintf("%s:%d", c.Host, c.Port)
//	cli := NewRedisClient(address, c.Auth, c.DB)
//	if err := cli.Ping().Err(); err != nil {
//		panic(err)
//	}
//	return cli
//}

//func (c RedisConfig) NewRedisClient() *redis.Conn {
//
//}

func (c RedisConfig) NewRedisPool() *redis.Pool {
	address := fmt.Sprintf("%s:%d", c.Host, c.Port)
	return NewRedisPool(address, c.Auth, c.DB)
}

//func NewRedisClient(address, password string, db int) *goredis.Client {
//	return goredis.NewClient(&goredis.Options{
//		Addr:     address,
//		Password: password,
//		DB:       db,
//	})
//}

func NewRedisPool(address, password string, db int) *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address,
				redis.DialPassword(password),
				redis.DialDatabase(db),
			)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:   100,
		MaxActive: 1000,
		Wait:      true,
	}
}
