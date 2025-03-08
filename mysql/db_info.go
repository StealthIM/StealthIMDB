package mysql

import (
	"StealthIMDB/config"
	"database/sql"
	_ "embed" // Embed
)

const dbCnt = 7

var conns = make([]ConnObj, dbCnt+1)
var dbs = [dbCnt + 1]string{
	"",
	"stimser_users",
	"stimser_msg",
	"stimser_file",
	"stimser_logging",
	"stimser_groups",
	"stimser_masterdb",
	"stimser_session",
}

// GetConn 获取数据库连接
func GetConn(dbindex int32) func() *sql.DB {
	return func() *sql.DB {
		if !conns[dbindex].online {
			return nil
		}
		if conns[dbindex].conn == nil {
			return nil
		}
		return conns[dbindex].conn
	}
}

type dbConnObj struct {
	Users    func() *sql.DB
	Msg      func() *sql.DB
	File     func() *sql.DB
	Logging  func() *sql.DB
	Groups   func() *sql.DB
	MasterDB func() *sql.DB
	Session  func() *sql.DB
}

// DBConn 数据库链接表
var DBConn = dbConnObj{
	Users:    GetConn(1),
	Msg:      GetConn(2),
	File:     GetConn(3),
	Logging:  GetConn(4),
	Groups:   GetConn(5),
	MasterDB: GetConn(6),
	Session:  GetConn(7),
}

//go:embed sql/users.sql
var sUsers string

//go:embed sql/msg.sql
var sMsg string

//go:embed sql/file.sql
var sFile string

//go:embed sql/logging.sql
var sLogging string

//go:embed sql/groups.sql
var sGroups string

//go:embed sql/masterdb.sql
var sMasterdb string

//go:embed sql/session.sql
var sSession string

var initArgs = []string{"", sUsers, sMsg, sFile, sLogging, sGroups, sMasterdb, sSession}

var cfgArgs [dbCnt + 1]config.MysqlNodeConfig

func setCfgNode(cfgID int, nodecfg config.MysqlNodeConfig) {
	if nodecfg.Host != "" {
		cfgArgs[cfgID].Host = nodecfg.Host
	} else {
		cfgArgs[cfgID].Host = cfg.Mysql.Host
	}
	if nodecfg.User != "" {
		cfgArgs[cfgID].User = nodecfg.User
	} else {
		cfgArgs[cfgID].User = cfg.Mysql.User
	}
	if nodecfg.Password != "" {
		cfgArgs[cfgID].Password = nodecfg.Password
	} else {
		cfgArgs[cfgID].Password = cfg.Mysql.Password
	}
	if nodecfg.Port != 0 {
		cfgArgs[cfgID].Port = nodecfg.Port
	} else {
		cfgArgs[cfgID].Port = cfg.Mysql.Port
	}
	if nodecfg.MaxConn != 0 {
		cfgArgs[cfgID].MaxConn = nodecfg.MaxConn
	} else {
		cfgArgs[cfgID].MaxConn = cfg.Mysql.MaxConn
	}
	if nodecfg.MinConn != 0 {
		cfgArgs[cfgID].MinConn = nodecfg.MinConn
	} else {
		cfgArgs[cfgID].MinConn = cfg.Mysql.MinConn
	}
}

// Setcfg 设置数据库配置
func Setcfg() {
	setCfgNode(1, cfg.Mysql.DBusers)
	setCfgNode(2, cfg.Mysql.DBmsg)
	setCfgNode(3, cfg.Mysql.DBfile)
	setCfgNode(4, cfg.Mysql.DBlogging)
	setCfgNode(5, cfg.Mysql.DBgroups)
	setCfgNode(6, cfg.Mysql.DBmasterdb)
	setCfgNode(7, cfg.Mysql.DBsession)
}
