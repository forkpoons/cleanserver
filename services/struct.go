package services

import "context"

type WorkerService struct {
	log func() *logrus.Entry
	app *gaarx.App
	ctx context.Context
}
