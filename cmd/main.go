package main

import (
	"log"
	"os"

	spc "github.com/doodles526/syncplayBot/client"
	"github.com/doodles526/syncplayBot/web"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	shouldServe = kingpin.Command("serve", "Should we serve?(currently only command and default)").Default().Action(serve)

	// Server args
	listenAddr = kingpin.Flag("listen", "Address for server to listen").Short('l').TCP()

	// SP Client args
	spHost     = kingpin.Flag("syncplay-server-host", "Host for syncplay server").Default("localhost").String()
	spPort     = kingpin.Flag("syncplay-server-port", "Port for syncplay server").Default("8999").String()
	spUsername = kingpin.Flag("syncplay-username", "Username for the bot to assume on syncplay").String()
	spPassword = kingpin.Flag("syncplay-password", "Password for the bot to use on syncplay").String()
	spRoom     = kingpin.Flag("syncplay-room", "Room to enter for the bot").String()
	// Known working sp version
	spVersion = kingpin.Flag("syncplay-version", "Room to enter for the bot").Default("1.2.255").String()
)

func main() {
	kingpin.Parse()
}

func serve(p *kingpin.ParseContext) error {
	l := log.New(os.Stderr, "SyncplayBot API Server", 0)

	cArgs := &spc.Args{
		Host:     *spHost,
		Port:     *spPort,
		Username: *spUsername,
		Password: *spPassword,
		Room:     *spRoom,
		Version:  *spVersion,
	}

	wArgs := &web.Args{
		ListenAddr:         *listenAddr,
		SyncplayClientArgs: cArgs,
	}

	if err := web.ServeForBot(wArgs); err != nil {
		l.Println(err)
		return err
	}
	return nil
}
