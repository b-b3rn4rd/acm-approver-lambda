package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/b-b3rn4rd/acm-approver-lambda/pkg/certificate"
	approver "github.com/b-b3rn4rd/acm-approver-lambda/pkg/lambda"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.SetLevel(logrus.DebugLevel)

	err := xray.Configure(xray.Config{LogLevel: "trace"})

	if err != nil {
		logrus.WithError(err).Fatal("Error configuring xray")
	}

	acmSvc := acm.New(session.Must(session.NewSession()))
	r53Svc := route53.New(session.Must(session.NewSession()))
	l := approver.New(certificate.New(acmSvc, r53Svc, logger), logger)
	lambda.Start(l.Handler)
}
