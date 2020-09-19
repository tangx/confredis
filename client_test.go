package confredis

import (
	"sync"
	"testing"

	. "github.com/onsi/gomega"
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

	_ = redi.PING()

	t.Run("redis-set", func(t *testing.T) {
		ok, err := redi.Do("SET", "book", "golang.com")
		NewWithT(t).Expect(err).Should(BeNil())
		NewWithT(t).Expect(ok.(string)).Should(Equal("OK"))
	})
	t.Run("redis-get", func(t *testing.T) {
		v, err := redi.Do("GET", "book")
		NewWithT(t).Expect(err).Should(BeNil())
		NewWithT(t).Expect(v.([]uint8)).Should(Equal([]uint8("golang.com")))
	})
}
