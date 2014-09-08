REDIS Controller for Revel.
========

Simple redis db controller for revel application written on top of [redigo](https://github.com/garyburd/redigo) client.

#### CONFIGURATION
few variables are exposed for correct redis connection and authentication.
each value expressed is the default configuration.
```go
redis.host=""
// Redis host ip

redis.port=6379
// Redis tcp port

redis.password=""
// Redis password

redis.trace=true
// enable redis trace log on stderr usefull on deploy

redis.check=true
// Check idle connection after their usage with a simple PING command, on failure
// close the application.

redis.idle=10
// Max idle connections in Poll, see: http://godoc.org/github.com/garyburd/redigo/redis#Pool.

redis.timeout=240
// Expressed in seconds, after this time the idle connection is closed.
```
#### Example production configuration
```go
redis.host=$REDIS_HOST
redis.port=$REDIS_PORT
redis.password=$REDIS_PASSWORD
redis.trace=false
redis.check=false
redis.idle=3
redis.timeout=60
```
#### USAGE
just import the reredigo app package and add the reredigo controller to the application
controller.
For the redigo package documentation go to: http://godoc.org/github.com/garyburd/redigo/redis.
For Redis documentation: http://redis.io/documentation.
```go
import (
	"github.com/alex179ohm/reredigo/app"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
	reredis.RedisController
}

func (c App) Index() revel.Result {
	connRedis = c.Pool.Get()
	defer connRedis.Close()

	connRedis.Do("SET", "foo", "bar")
	bar, err := connRedis.Do("GET", "foo")
	if err != nil {
		revel.ERROR.Println(err)
	}

	c.RenderArgs["foo"] = bar
	c.Render()
}
```
