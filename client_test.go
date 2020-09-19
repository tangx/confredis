package confredis

import (
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestMain(t *testing.T) {

	redi := Redis{
		Host:     "127.0.0.1",
		Port:     6379,
		Password: "Pass12345",
		DB:       9,
	}

	logrus.SetLevel(9)

	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			redi.Init()
		}()
		wg.Wait()
	}

	// time.Sleep(1 * time.Second)
	_ = redi.PING()

}
