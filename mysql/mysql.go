package mysql

import (
	"StealthIMDB/config"
	"database/sql"
	_ "embed" // Embed
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql driver
)

// ConnObj 链接对象
type ConnObj struct {
	conn   *sql.DB
	online bool
}

// SQLDBErrorObj Mysql错误对象
type SQLDBErrorObj struct {
	Code    int
	Message string
}

func (e *SQLDBErrorObj) Error() string {
	return fmt.Sprintf("[MySQL]ErrCode: %d, Msg: %s", e.Code, e.Message)
}

var cfg config.Config

func autoReconn(connID int) {
	for {
		if !conns[connID].online {
			conn(connID)
		}
		time.Sleep(8 * time.Second)
		err := conns[connID].conn.Ping()
		if err != nil {
			log.Printf("[MySQL]MySQL [%s] connect error: %v\n", dbs[connID], err)
			conns[connID].online = false
		}
	}
}

func conn(connID int) error {
	log.Printf("[MySQL]Connect to MySQL [%s]\n", dbs[connID])
	dbinfo := cfgArgs[connID].User + ":" + cfgArgs[connID].Password + "@tcp(" + cfgArgs[connID].Host + ":" + strconv.Itoa(cfgArgs[connID].Port) + ")/" + dbs[connID] + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dbinfo)
	if err != nil {
		log.Printf("[MySQL]Error connecting to MySQL [%s]: %+v\n", dbs[connID], err)
		conns[connID].online = false
		time.Sleep(3000)
		return err
	}
	db.SetMaxIdleConns(cfgArgs[connID].MinConn)
	db.SetMaxOpenConns(cfgArgs[connID].MaxConn)
	conns[connID].conn = db
	conns[connID].online = true
	return nil
}

// Connect 初始化并链接
func Connect(setCfg config.Config) {
	cfg = setCfg
	Setcfg()
	Init()
	for i := 1; i <= dbCnt; i++ {
		go autoReconn(i)
	}
}
