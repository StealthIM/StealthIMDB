package mysql

import (
	"StealthIMDB/errorcode"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func initr(connID int) error {
	if !conns[connID].online {
		return &SQLDBErrorObj{Code: int(errorcode.ServerInternalNetworkError), Message: fmt.Sprintf("MySQL [%s] is not online", dbs[connID])}
	}
	tx, err := conns[connID].conn.Begin()
	if err != nil {
		return err
	}
	sqlStatements := strings.Split(initArgs[connID], ";")
	for i, statement := range sqlStatements {
		statement = strings.TrimSpace(statement)
		if statement != "" {
			sqlStatements[i] = statement + ";"
		}
	}
	if len(sqlStatements) > 0 {
		sqlStatements = sqlStatements[:len(sqlStatements)-1]
	}
	for _, statement := range sqlStatements {
		_, err = tx.Exec(statement)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func initConn(connID int) error {
	log.Printf("[MySQL]Init MySQL [%s]\n", dbs[connID])
	dbinfo := cfgArgs[connID].User + ":" + cfgArgs[connID].Password + "@tcp(" + cfgArgs[connID].Host + ":" + strconv.Itoa(cfgArgs[connID].Port) + ")/" + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dbinfo)
	if err != nil {
		log.Printf("[MySQL]Error init MySQL [%s]: %+v\n", dbs[connID], err)
		conns[0].online = false
		time.Sleep(3000)
		return err
	}
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)
	conns[0].conn = db
	conns[0].online = true
	return nil
}

// Init 初始化数据库
func Init() error {
	// for {
	// 	err := conn(0)
	// 	if err == nil {
	// 		break
	// 	} else {
	// 		time.Sleep(3 * time.Second)
	// 	}
	// }
	// for {
	// 	err := initr(0)
	// 	if err == nil {
	// 		break
	// 	} else {
	// 		log.Printf("[MySQL]MySQL [] init error: %v\n", err)
	// 		time.Sleep(3 * time.Second)
	// 	}
	// }
	for i := 1; i <= dbCnt; i++ {
		initArgs[0] = "CREATE DATABASE IF NOT EXISTS `" + dbs[i] + "`;"
		for {
			err := initConn(i)
			if err == nil {
				break
			} else {
				panic(err)
			}
		}
		for {
			err := initr(0)
			if err == nil {
				break
			} else {
				log.Printf("[MySQL]MySQL [%s] init error: %v\n", dbs[i], err)
				panic(err)
			}
		}
		conns[0].conn.Close()
		conns[0].online = false
	}
	for i := 1; i <= dbCnt; i++ {
		for {
			err := conn(i)
			if err == nil {
				break
			} else {
				time.Sleep(3 * time.Second)
			}
		}
		for {
			err := initr(i)
			if err == nil {
				break
			} else {
				log.Printf("[MySQL]MySQL [%s] connect error: %v\n", dbs[i], err)
				time.Sleep(3 * time.Second)
			}
		}
	}
	log.Printf("[MySQL]MySQL inited!\n")
	return nil
}
