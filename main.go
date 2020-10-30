package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/forkpoons/cleanserver/core"
	_ "github.com/forkpoons/cleanserver/core"
	"github.com/forkpoons/cleanserver/services/notify"
	"github.com/forkpoons/cleanserver/services/web"
	"github.com/gaarx/gaarx"
	"github.com/sirupsen/logrus"
	"github.com/zergu1ar/logrus-filename"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	flag.Parse()
	var stop = make(chan os.Signal)
	var done = make(chan bool, 1)
	ctx, finish := context.WithCancel(context.Background())
	var application = &gaarx.App{}
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	_ = application.InitializeLogger(gaarx.FileLog, "log.log", &logrus.TextFormatter{DisableColors: true})
	filenameHook := filename.NewHook()
	filenameHook.Field = "line"
	application.GetLog().AddHook(filenameHook)
	application.Initialize(
		gaarx.WithContext(ctx),
		gaarx.WithStorage(core.TelegramUsersScope),
		gaarx.WithServices(
			web.Create(ctx),
			notify.Create(ctx),
		),
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
