package dbinit

import (
	"fmt"
	"os"

	"github.com/gomodule/redigo/redis"
)

func GetRedis() (redis.Conn, error) {
	host := os.Getenv("hostRD")
	port := os.Getenv("portRD")
	c, err := redis.DialURL(fmt.Sprintf("redis://user:@%s:%s/0", host, port))
	if err != nil {
		return nil, err
	}
	return c, nil

}
