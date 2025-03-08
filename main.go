package main

import (
	"StealthIMDB/config"
	"StealthIMDB/grpc"
	"StealthIMDB/mysql"
	"StealthIMDB/redis"
	"log"
)

func main() {
	cfg := config.ReadConf()
	log.Printf("Start server [%v]\n", config.Version)
	log.Printf("+ GRPC\n")
	log.Printf("    Host: %s\n", cfg.GRPCProxy.Host)
	log.Printf("    Port: %d\n", cfg.GRPCProxy.Port)
	log.Printf("+ Mysql\n")
	log.Printf("    Host: %s\n", cfg.Mysql.Host)
	log.Printf("    Port: %d\n", cfg.Mysql.Port)
	log.Printf("+ Redis\n")
	log.Printf("    Host: %s\n", cfg.Redis.Host)
	log.Printf("    Port: %d\n", cfg.Redis.Port)
	// log.Printf("+ RocketMQ\n")
	// log.Printf("    Host: %s\n", cfg.RocketMQ.Host)
	// log.Printf("    Port: %d\n", cfg.RocketMQ.Port)
	go mysql.Connect(cfg)
	go redis.Connect(cfg)
	grpc.Start(cfg)
}
