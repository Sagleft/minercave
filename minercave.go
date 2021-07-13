package main

import (
	"minercave/app"
	"minercave/net"
)

var cfg net.Config

func main() {
	app.Configure(&cfg)
	app.Exec(&cfg)
}
