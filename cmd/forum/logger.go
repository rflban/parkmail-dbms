package main

import (
	"context"
	"github.com/rflban/parkmail-dbms/internal/pkg/forum/constants"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

const helloMessage = " ________ ________  ________  ___  ___  _____ ______      \n|\\  _____\\\\   __  \\|\\   __  \\|\\  \\|\\  \\|\\   _ \\  _   \\    \n\\ \\  \\__/\\ \\  \\|\\  \\ \\  \\|\\  \\ \\  \\\\\\  \\ \\  \\\\\\__\\ \\  \\   \n \\ \\   __\\\\ \\  \\\\\\  \\ \\   _  _\\ \\  \\\\\\  \\ \\  \\\\|__| \\  \\  \n  \\ \\  \\_| \\ \\  \\\\\\  \\ \\  \\\\  \\\\ \\  \\\\\\  \\ \\  \\    \\ \\  \\ \n   \\ \\__\\   \\ \\_______\\ \\__\\\\ _\\\\ \\_______\\ \\__\\    \\ \\__\\\n    \\|__|    \\|_______|\\|__|\\|__|\\|_______|\\|__|     \\|__|\n                                                          "

func BindLoggers(ctx context.Context) context.Context {
	logrus.SetOutput(ioutil.Discard)

	withSetupLog := context.WithValue(
		ctx,
		constants.SetupLogKey,
		logrus.WithField("type", "setup"),
	)
	withAccessLog := context.WithValue(
		withSetupLog,
		constants.AccessLogKey,
		logrus.WithField("type", "access"),
	)
	withDeliveryLog := context.WithValue(
		withAccessLog,
		constants.DeliveryLogKey,
		logrus.WithField("type", "delivery"),
	)
	withUseCaseLog := context.WithValue(
		withDeliveryLog,
		constants.UseCaseLogKey,
		logrus.WithField("type", "usecase"),
	)
	withRepoLog := context.WithValue(
		withUseCaseLog,
		constants.RepoLogKey,
		logrus.WithField("type", "repo"),
	)

	return withRepoLog
}
