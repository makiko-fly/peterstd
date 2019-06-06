package peterstd

import (
	"io/ioutil"

	"gitlab.wallstcn.com/spider/peterstd/datasource"
	goredis "gopkg.in/redis.v5"
)

type (
	MySQLConfig    = datasource.MySQLConfig
	RedisConfig    = datasource.RedisConfig
	RedisSubConfig = datasource.RedisSubConfig
)

func LoadScript(rdb *goredis.Client, filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	sha1, err := rdb.ScriptLoad(string(bytes)).Result()
	if err != nil {
		panic(err)
	}
	With("script_name", filename).WithField("sha1", sha1).Info("Load script success")
	return sha1
}
