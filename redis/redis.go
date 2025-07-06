package redis

import (
	"StealthIMDB/config"
	"StealthIMDB/errorcode"
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const dbCnt int = 1

var ctx = context.Background()

// ConnObj 链接对象
type ConnObj struct {
	db     *redis.Client
	online bool
}

// CacheRedisError 错误对象
type CacheRedisError struct {
	Code    int
	Message string
}

func (e *CacheRedisError) Error() string {
	return fmt.Sprintf("[Redis]ErrCode: %d, Msg: %s", e.Code, e.Message)
}

var cfg config.Config
var conndb = make([]ConnObj, dbCnt)

// GetConn 获取链接
func GetConn(id int) *redis.Client {
	if id >= dbCnt {
		return nil
	}
	if conndb[id].online == false {
		return nil
	}
	if conndb[id].db == nil {
		return nil
	}
	return conndb[id].db
}

func autoReconn(connID int) {
	for {
		if !conndb[connID].online {
			conn(connID)
		}
		time.Sleep(8 * time.Second)
		_, err := conndb[connID].db.Ping(ctx).Result()
		if err != nil {
			log.Printf("[Redis]Redis connect [%d] error: %v\n", connID, err)
			conndb[connID].online = false
		}
	}
}

func initr(connID int) error {
	db := conndb[connID].db
	err := db.Set(ctx, "GatewayInfo", fmt.Sprintf("StealthIM:GTWINFO;Version:%s", config.Version), 0).Err()
	if err != nil {
		return &CacheRedisError{Code: int(errorcode.ServerInternalComponentError), Message: "Redis init error"}
	}
	return nil
}

func conn(connID int) error {
	log.Printf("[Redis]Connect to Redis [%d]\n", connID)
	db := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + strconv.Itoa(cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DBID + connID,
	})
	conndb[connID].db = db
	conndb[connID].online = true
	time.Sleep(1 * time.Second)
	_, err := db.Ping(ctx).Result()
	if err != nil {
		return &CacheRedisError{Code: int(errorcode.ServerInternalNetworkError), Message: "Redis connect error"}
	}
	return nil
}

// Connect 链接Redis
func Connect(setCfg config.Config) {
	cfg = setCfg
	for i := range dbCnt {
		for {
			err := conn(i)
			if err == nil {
				break
			} else {
				log.Printf("[Redis]Redis connect error: %v\n", err)
				time.Sleep(5 * time.Second)
			}
		}
		for {
			err := initr(i)
			if err == nil {
				break
			} else {
				log.Printf("[Redis]Redis init error: %v\n", err)
				time.Sleep(5 * time.Second)
			}
		}
	}
	go autoReconn(0)
}
