package config

// Config 主配置
type Config struct {
	GRPCProxy GRPCProxyConfig `toml:"grpc"`
	Mysql     MysqlConfig     `toml:"mysql"`
	Redis     RedisConfig     `toml:"redis"`
	// RocketMQ  RocketMQConfig  `toml:"rocketmq"`
}

// GRPCProxyConfig grpc Server配置
type GRPCProxyConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
	Log  bool   `toml:"log"`
}

// MysqlConfig mysql配置
type MysqlConfig struct {
	Host       string          `toml:"host"`
	User       string          `toml:"user"`
	Password   string          `toml:"password"`
	Port       int             `toml:"port"`
	MaxConn    int             `toml:"maxconn"`
	MinConn    int             `toml:"minconn"`
	DBmsg      MysqlNodeConfig `toml:"db_msg"`
	DBusers    MysqlNodeConfig `toml:"db_users"`
	DBmasterdb MysqlNodeConfig `toml:"db_masterdb"`
	DBlogging  MysqlNodeConfig `toml:"db_logging"`
	DBgroups   MysqlNodeConfig `toml:"db_groups"`
	DBfile     MysqlNodeConfig `toml:"db_file"`
	DBsession  MysqlNodeConfig `toml:"db_session"`
	Prefix     string          `toml:"prefix"`
}

// MysqlNodeConfig mysql链接节点
type MysqlNodeConfig struct {
	Host     string `toml:"host"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Port     int    `toml:"port"`
	MaxConn  int    `toml:"maxconn"`
	MinConn  int    `toml:"minconn"`
}

// RedisConfig redis配置
type RedisConfig struct {
	Host      string `toml:"host"`
	Port      int    `toml:"port"`
	Password  string `toml:"password"`
	DBID      int    `toml:"dbname"`
	CacheTime int    `toml:"cachetime"`
}

// type RocketMQConfig struct {
// 	Host string `toml:"host"`
// 	Port int    `toml:"port"`
// }
