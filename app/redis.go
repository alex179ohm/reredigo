package reredis

import (
	"github.com/agtorre/gocolorize"
	"github.com/garyburd/redigo/redis"
	"github.com/revel/revel"
	"log"
	"os"
	"strings"
	"time"
)

var (
	pool *redis.Pool
)

func newPool(proto, server, password string,
	idle, active, timeout int,
	trace, check bool) *redis.Pool {
	if len(proto) == 0 {
		proto = "tcp"
	}

	redis_log := gocolorize.NewColor("magenta+b")
	r := redis_log.Paint
	redisLog := log.New(os.Stdout, r("REDIS "), log.Ldate|log.Ltime|log.Lshortfile)
	return &redis.Pool{
		MaxIdle:     idle,
		MaxActive:   active,
		IdleTimeout: timeout * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(proto, server)
			if trace {
				c = redis.NewLoggingConn(c, redisLog, "")
			}
			if err != nil {
				return nil, err
			}
			if len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if check {
				if _, err := c.Do("PING"); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if check {
				if _, err := c.Do("PING"); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func Init() {
	// configurations
	var (
		found bool
		//        redisUrl string
		host     string
		port     string
		password string
		idle     int
		timeout  int
		//		active int
		trace bool
		check bool
	)

	//TODO: include a regexp match and REDIS_URL env variable

	// set default configurations

	/*
	   var redis_url = 0;
	   //TODO: add unixUrl parser and check
	   if redisUrl, found = revel.Config.String("redis.url"); !found {
	       redisUrl = ""
	       redis_url = 0;
	       revel.INFO.Printf("Redis: redis")
	   } else {
	       redis_url = 1;
	   }*/
	//    if redis_url != 1 {
	if host, found = revel.Config.String("redis.host"); !found {
		host = ""
		revel.INFO.Println("Redis: redis.host not found default is localhost")
	}
	if port, found = revel.Config.String("redis.port"); !found {
		port = "6379"
		revel.INFO.Println("Redis: redis.port not found default is 6379")
	}
	if password, found = revel.Config.String("redis.password"); !found {
		password = ""
		revel.INFO.Println("Redis: redis.password not found default is blank string")
	}
	//    }

	if idle, found = revel.Config.Int("redis.idle"); !found {
		idle = 10
		revel.INFO.Println("REDIS: redis.idle not found, default is 10")
	}
	if timeout, found = revel.Config("redis.timeout"); !found {
		timeout = 240
		revel.INFO.Println("REDIS: redis.timeout not found, default is 240")
	}
	if trace, found = revel.Config.Bool("redis.trace"); !found {
		trace = true
		revel.INFO.Println("REDIS: redis.trace not found default is true")
	}
	if check, found = revel.Config.Bool("redis.check"); !found {
		check = true
		revel.INFO.Println("REDIS: redis.check not found default is true")
	}

	url := []string{host, port}

	pool = newPool("tcp", strings.Join(url, ":"), password, idle, 0, timeout, trace, check)
}

type RedisController struct {
	*revel.Controller
	Pool *redis.Pool
}

func (c *RedisController) Begin() revel.Result {
	c.Pool = pool
	return nil
}

func init() {
	revel.OnAppStart(Init)
	revel.InterceptMethod((*RedisController).Begin, revel.BEFORE)
}
