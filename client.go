package confredis

import (
	"fmt"
	"sync"

	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

type Redis struct {
	Protocol string
	Endpoint string `comment:"alias: tcp://:password@host:port/db"`
	Host     string `env:"host"`
	Port     int    `env:"port"`
	Password string `env:"password"`
	DB       int    `env:"db"`

	MaxIdle   int
	MaxActive int

	IdleTimeout     time.Duration `comment:"seconds"`
	MaxConnLifetime time.Duration `comment:"seconds"`
	ConnectTimeout  time.Duration `comment:"seconds"`
	ReadTimeout     time.Duration `comment:"seconds"`

	pool *redis.Pool
}

func (r *Redis) SetDefaults() {
	if r.Port == 0 {
		r.Port = 6379
	}

	if r.Protocol == "" {
		r.Protocol = "tcp"
	}
	if r.MaxConnLifetime == 0 {
		r.MaxConnLifetime = time.Duration(300) * time.Second
	}

	if r.IdleTimeout == 0 {
		r.IdleTimeout = time.Duration(60) * time.Second
	}

	if r.MaxActive == 0 {
		r.MaxActive = 5
	}

	if r.MaxIdle == 0 {
		r.MaxIdle = 5
	}

	if r.ConnectTimeout == 0 {
		r.ConnectTimeout = time.Duration(5) * time.Second
	}

	if r.ReadTimeout == 0 {
		r.ReadTimeout = time.Duration(5) * time.Second
	}

}

func (r *Redis) initial() *redis.Pool {

	dialFunc := func() (redis.Conn, error) {
		c, err := redis.Dial("tcp",
			fmt.Sprintf("%s:%d", r.Host, r.Port),
			redis.DialReadTimeout(time.Duration(r.ReadTimeout)),
			redis.DialConnectTimeout(time.Duration(r.ConnectTimeout)),
		)
		if err != nil {
			logrus.Errorf("redis dial error: %s", err)
			return nil, err
		}

		if r.Password != "" {
			if _, err := c.Do("AUTH", r.Password); err != nil {
				logrus.Errorf("redis auth error: %s", err)
				return nil, err
			}
		}

		if r.DB != 0 {
			if _, err := c.Do("SELECT", r.DB); err != nil {
				logrus.Errorf("redis select db error: %s", err)
				return nil, err
			}
		}
		return c, nil
	}

	return &redis.Pool{
		Dial:            dialFunc,
		MaxIdle:         r.MaxIdle,
		MaxActive:       r.MaxActive,
		IdleTimeout:     time.Duration(r.IdleTimeout) * time.Second,
		MaxConnLifetime: time.Duration(r.MaxConnLifetime) * time.Second,
	}
}

var (
	mutex sync.Mutex
)

func (r *Redis) Init() {

	/*
		单例模式: 线程安全 , 锁
		// https://zhuanlan.zhihu.com/p/33102022
	*/

	mutex.Lock()
	defer mutex.Unlock()

	r.SetDefaults()
	if r.pool == nil {
		r.pool = r.initial()
	}
}

func (r *Redis) Get() redis.Conn {
	if r.pool != nil {
		return r.pool.Get()
	}
	return nil
}

func (r *Redis) PING() error {

	pong, err := r.Do("PING")
	if err != nil {
		return err
	}

	logrus.Debugf("redis conn ping sucess: %s", pong)
	return nil
}

func (r *Redis) Do(command string, args ...interface{}) (interface{}, error) {

	c := r.Get()
	defer func() {
		err := c.Close()
		if err != nil {
			logrus.Errorf("redis conn close err: %s", err)
		}
	}()

	if len(args) == 0 {
		return c.Do(command)
	}

	return c.Do(command, args...)

}
