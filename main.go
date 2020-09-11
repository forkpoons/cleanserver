package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gaarx/database"
	"github.com/gaarx/gaarx"
	"github.com/sirupsen/logrus"
	"github.com/zergu1ar/logrus-filename"
	"md/core"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var configFile = flag.String("config", "config", "Project config")

func main() {
	flag.Parse()
	var stop = make(chan os.Signal)
	var done = make(chan bool, 1)
	ctx, finish := context.WithCancel(context.Background())
	var application = &gaarx.App{}
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	err := application.LoadConfig(*configFile, "toml", &Configs{})
	if err != nil {
		panic(err)
	}
	configs := application.Config().(*Configs)
	// TODO: handle errors here
	_ = application.InitializeLogger(gaarx.FileLog, configs.System.Log, &logrus.TextFormatter{DisableColors: true})
	filenameHook := filename.NewHook()
	filenameHook.Field = "line"
	application.GetLog().AddHook(filenameHook)
	application.GetLog().SetLevel(logrus.Level(configs.System.LogLevel))
	application.Initialize(
		database.WithDatabase(configs.System.DB, &core.Account{}, &core.Session{}, &core.Configuration{}),
		gaarx.WithContext(ctx),
		gaarx.WithStorage(core.ScopeSessions, core.ScopeAccounts, core.ScopeConf, core.ScopeBalance),
		//gaarx.WithServices(
		//	worker.Create(ctx),
		//),
	)
	go func() {
		sig := <-stop
		time.Sleep(2 * time.Second)
		finish()
		fmt.Printf("caught sig: %+v\n", sig)
		done <- true
	}()
	application.Work()
	<-ctx.Done()
	os.Exit(0)
}
