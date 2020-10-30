package web

import (
	"context"
	"github.com/forkpoons/cleanserver/core"
	"github.com/gaarx/gaarx"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Service struct {
	log func() *logrus.Entry
	app *gaarx.App
	ctx context.Context
	srv *http.Server
}

func Create(ctx context.Context) *Service {
	return &Service{
		ctx: ctx,
	}
}

func (m *Service) Start(app *gaarx.App) error {
	m.app = app
	m.log = func() *logrus.Entry {
		return app.GetLog().WithField("service", "Web service")
	}

	r := gin.Default()
	r.LoadHTMLFiles("/var/www/templates/index.html")
	srv := &http.Server{
		Addr:    ":80",
		Handler: r,
	}
	r.GET("/", index)
	r.Static("/style", "/var/www/templates/style")
	r.Static("/img", "/var/www/templates/img")
	r.Static("/scripts", "/var/www/templates/scripts")
	r.POST("/send.go", m.messageHandler())
	m.srv = srv
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func index(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"status": "success",
	})
}

func (m *Service) messageHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{
			"status": "success",
		})
		var qwe interface{}
		_ = c.ShouldBindJSON(&qwe)

		m.app.Event(core.MessageEvent).Dispatch(qwe)
	}
}

func (m *Service) Stop() {
}

func (m *Service) GetName() string {
	return "Web service"
}
