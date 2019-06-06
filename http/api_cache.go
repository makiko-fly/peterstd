package http

import (
	"errors"
	"math/rand"
	"time"

	"github.com/gomodule/redigo/redis"
	json "github.com/json-iterator/go"
	"github.com/labstack/echo"
	"gitlab.wallstcn.com/spider/peterstd"
)

type (
	apiCache struct {
		rp *redis.Pool
	}
	APICache = *apiCache
)

func NewAPICache(rp *redis.Pool) APICache {
	return &apiCache{
		rp: rp,
	}
}

func (a *apiCache) WithCache(expiration time.Duration) peterstd.HTTPHandlerWrapper {
	return func(handler peterstd.HTTPHandler) peterstd.HTTPHandler {
		return func(ctx echo.Context) (interface{}, error) {
			// 1. 在缓存中查找url
			// 2. 如果没有，执行handler，并缓存结果
			// 3. 如果有，直接返回结果

			var ttl = int32(expiration/time.Second) + rand.Int31()%5

			url := ctx.Request().URL
			key := "flash:api:cache:" + url.Path
			if url.RawQuery != "" {
				key += "?" + url.RawQuery
			}

			conn := a.rp.Get()
			defer conn.Close()

			for i := 0; i < 15; i++ {
				content := lock(conn, key, 3)
				if content == LockSuccess {
					res, err := handler(ctx)
					if err != nil {
						return nil, err
					}
					body, err := json.Marshal(res)
					if err != nil {
						return nil, err
					}

					conn.Send("MULTI")
					conn.Send("SET", key, body)
					conn.Send("EXPIRE", key, ttl)
					if _, err := conn.Do("EXEC"); err != nil {
						peterstd.WithError(err).Error("Error in set data to redis")
						return nil, err
					}
					return res, err
				} else if content == LockExists {
					time.Sleep(time.Millisecond * 200)
					continue
				} else {
					return json.RawMessage(content), nil
				}
			}
			return nil, errors.New("timed out")
		}
	}
}

const (
	LockSuccess = "lock_success"
	LockExists  = "lock_exists"
)

func lock(conn redis.Conn, key string, ttl int32) string {
	content, _ := redis.String(conn.Do("EVAL", `
local success = redis.call("SETNX", KEYS[1], ARGV[1])
local content = "lock_failed"
if success == 1 then
    redis.call("EXPIRE", KEYS[1], ARGV[2])
    content = "lock_success"
else
    content = redis.call("GET", KEYS[1])
end
return content 
`, 1, key, LockExists, ttl))
	return content
}
