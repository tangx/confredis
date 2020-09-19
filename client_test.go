package confredis

import (
	"testing"
)

func TestMain(t *testing.T) {

	redi := Redis{
		Host:     "127.0.0.1",
		Port:     6379,
		Password: "Pass123452",
		DB:       99,
	}

	redi.Init()

	_ = redi.PING()

}
