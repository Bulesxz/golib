package redislock

import (
	"fmt"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
)

func TestLock(t *testing.T) {
	rd := Redispool.Get()
	defer rd.Close()

	go func() {
		Alock := RedisLock{lockKey: "xxxxx"}
		err := Alock.Lock(&rd, 5) //5 秒后自动删除Alock

		time.Sleep(7 * time.Second) //等待7秒
		fmt.Println("111", err)
		Alock.Unlock(&rd) //想删除的是Alock锁，但是Alock 已经被自动删除 ,Block由于value 不一样，所以也不会删除
	}()

	time.Sleep(6 * time.Second) //此时Alock 已经被删除
	Block := RedisLock{lockKey: "xxxxx"}
	err := Block.Lock(&rd, 5) //此时 会获取新的lock Block
	fmt.Println("222", err)

	time.Sleep(2 * time.Second)
	Clock := RedisLock{lockKey: "xxxxx"}
	err = Clock.Lock(&rd, 5) //想获取新的lock Clock，但由于 Block还存在，返回错误
	fmt.Println("333", err)

	time.Sleep(10 * time.Second)

}

var Redispool *redis.Pool

func init() {
	Redispool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			tcp := fmt.Sprintf("%s:%d", "127.0.0.1", 32768)
			c, err := redis.Dial("tcp", tcp)
			if err != nil {
				return nil, err
			}
			//fmt.Println("connect redis success!")
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

}
