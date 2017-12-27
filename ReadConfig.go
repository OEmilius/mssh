package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

//var cfg Config

//type Config struct {
//	Hostst []string
//}

var login string
var password string
var cmd string

func ProcessFlags() {
	flag.StringVar(&login, "login", "root", "login for all hosts")
	flag.StringVar(&password, "password", "tour", "password for all hosts")
	flag.StringVar(&cmd, "cmd", "uname", "command to execute")
	flag.Parse()
}

func ReadConfig(fname string) Config {
	file, err := os.Open(fname)
	if err != nil {
		//dbg.Println("opening config file error", err)
		panic(err)
	}
	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		//dbg.Println("error reading config file", err)
		fmt.Println("error reading config file", err)
		panic(err)
	}
	//fmt.Println("config=", config)
	return config
}

func init() {
	ProcessFlags()
	cfg = ReadConfig("config.json")
}
