package redisConn
import (
	"github.com/garyburd/redigo/redis"
	"../config"
	"sync"
)
var redisConn redis.Conn
var once sync.Once
var pool = NewPool()

/*func GetRedis() (redis.Conn, error) {

	var err error
	once.Do(func(){
		var conf = config.GetConfig()
		redisConn, err = redis.Dial("tcp", conf.String("redis::host") + ":" + conf.String("redis::port"))
		_, err = redisConn.Do("select", conf.String("redis::index"))
	})
	return redisConn, err
}*/

func  GetString(key string) (value string, err error)  {

	c := pool.Get()
	value, err = redis.String(c.Do("get", key))

	return value, err
}

func  SetString(key string, value string)  {

	c := pool.Get()
	c.Do("set", key, value)

	return
}

func NewPool() *redis.Pool {
	var conf = config.GetConfig()

	return &redis.Pool{
		MaxIdle: 80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.String("redis::host") + ":" + conf.String("redis::port"))
			if err != nil {
				panic(err.Error())
			}
			c.Do("SELECT", conf.String("redis::index"))
			return c, err
		},
	}
}