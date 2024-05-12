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
	ServerDomain   string
	ServerPort     int
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
		ServerDomain:   cfg.Section("api").Key("domain").String(),
		ServerPort:     cfg.Section("api").Key("port").MustInt(),
		FEUrl:          cfg.Section("fe").Key("url").String(),
	}
}
