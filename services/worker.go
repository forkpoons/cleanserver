package services

import (
	"context"
	"github.com/gaarx/gaarx"
	"github.com/sirupsen/logrus"
)

func Create(ctx context.Context) *WorkerService {
	return &WorkerService{
		ctx: ctx,
	}
}

func (m *WorkerService) Start(app *gaarx.App) error {
	m.app = app
	m.log = func() *logrus.Entry {
		return app.GetLog().WithField("service", "Worker")
	}
	return nil
}

func (m *WorkerService) Stop() {
}

func (m *WorkerService) GetName() string {
	return "Worker"
}
