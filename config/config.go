package config

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	DBDriverName   string
	DBName         string
	DBUserName     string
	DBUserPassword string
	DBHost         string
	DBPort         string
	AppEnv         string
	ServerDomain   string
	ServerPort     int
	ServerGrpcPort string
	CertFile       string
	KeyFile        string
	FEUrl          string
}

var Config ConfigList

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	Config = ConfigList{
		DBDriverName:   cfg.Section("db").Key("driver").String(),
		DBName:         cfg.Section("db").Key("db_name").String(),
		DBUserName:     cfg.Section("db").Key("user").String(),
		DBUserPassword: cfg.Section("db").Key("password").String(),
		DBHost:         cfg.Section("db").Key("host").String(),
		DBPort:         cfg.Section("db").Key("port").String(),
		AppEnv:         cfg.Section("api").Key("app_env").String(),
		ServerDomain:   cfg.Section("api").Key("domain").String(),
		ServerPort:     cfg.Section("api").Key("port").MustInt(),
		ServerGrpcPort: cfg.Section("api").Key("grpc_port").String(),
		CertFile:       cfg.Section("api").Key("certFile").String(),
		KeyFile:        cfg.Section("api").Key("keyFile").String(),
		FEUrl:          cfg.Section("fe").Key("url").String(),
	}
}
