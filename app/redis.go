package redis

import (
    "github.com/garyburd/redigo/redis";
    "github.com/revel/revel";
    "strings";
    "time";
)

var (
    pool *redis.Pool
)

func newPool(proto, server, password string, idle, active int) *redis.Pool {
    if len(proto) == 0 { proto = "tcp" }
    return &redis.Pool{
        MaxIdle: idle,
        MaxActive: active,
        IdleTimeout: 240 * time.Second,
        Dial: func () (redis.Conn, error) {
            c, err := redis.Dial(proto, server)
            if err != nil {
                return nil, err
            }
            if len(password) > 0 {
                if _, err := c.Do("AUTH", password); err != nil {
                    c.Close()
                    return nil, err
                }
            } else {
                if _, err := c.Do("PING"); err != nil {
                    c.Close()
                    return nil, err
                }
            }
            return c, err
        },
        TestOnBorrow: func(c redis.Conn, t time.Time) error {
            if _, err := c.Do("PING"); err != nil {
                return err
            }
            return nil
        },
    }
}

func Init() {
    // configurations
    var (
        found bool
        host string
        port string
        password string
    )

    //TODO: include a regexp match and REDIS_URL env variable

    // set default configirations

    //TODO: add unix url
    if host, found = revel.Config.String("redis.host"); !found {
        host = ""
        revel.WARN.Printf("Redis: redis.host not found")
    }
    if port, found = revel.Config.String("redis.port"); !found {
        port = "6379"
        revel.WARN.Printf("Redis: redis.port not found")
    }
    if password, found = revel.Config.String("redis.password"); !found {
        password = ""
        revel.WARN.Printf("Redis: redis.password not found")
    }

    url := []string{host, port}

    pool = newPool("tcp", strings.Join(url, ":"), password, 3, 0)
}

type RedisController struct {
    *revel.Controller
    Redis *redis.Pool
}

func (c *RedisController) Begin() revel.Result {
    c.Redis = pool
    return nil
}

func init() {
    revel.OnAppStart(Init)
    revel.InterceptMethod((*RedisController).Begin, revel.BEFORE)
}
