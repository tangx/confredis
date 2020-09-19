package confredis

import (
	"testing"
)

func TestMain(t *testing.T) {

	redi := Redis{
		Host:     "127.0.01",
		Port:     6379,
		Password: "Pass12345",
		DB:       1,
	}

	redi.Init()

	_ = redi.PING()

}
