REDIS Controller for Revel.
========

Simple redis db controller for revel application written on top of [redigo](https://github.com/garyburd/redigo) client api.

Installation
------------
```sh
go get github.com/alex179ohm/reredigo
```
#### external dependecies
[revel](https://github.com/revel/revel) #revel app  
[gocolorize](https://github.com/agtorre/gocolorize) #beauty logs  
[redigo](https://github.com/garyburd/redigo) #redis client api  

Configuration
-------------
Few variables are exposed for redis connection and authentication.
each value expressed is the default configuration.  
Default configuration is optimized for a development environment.
```go
redis.host=""
// Redis host ip (blank string means localhost)

redis.port=6379
// Redis tcp port

redis.password=""
// Redis password

redis.trace=true
// enable redis trace log on stderr usefull on development

redis.check=true
// Check idle connection after their usge, on failure
// close the application. (PING command is used to check connection)

redis.idle=10
// Max idle connections in Poll, see: [Pool](http://godoc.org/github.com/garyburd/redigo/redis#Pool).

redis.timeout=240
// Expressed in seconds, after this time the idle connection is closed.
```
#### Example configuration used on production
NOTE: This is just a configuration example, redis.idle and redis.timeout values
have to be relative at the connections (respectively number and duration)
of your specific application.
```go
redis.host=$REDIS_HOST
redis.port=$REDIS_PORT
redis.password=$REDIS_PASSWORD
redis.trace=false
redis.check=false
redis.idle=3
redis.timeout=60
```
Usage
-----
Just import the reredigo app package and add the reredigo controller to the application
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

LICENSE:
Licensed under the BSD 2-clause License (Simplied BSD License) http://opensource.org/licenses/BSD-2-Clause
