package lambda

import (
	"github.com/b-b3rn4rd/acm-approver-lambda/pkg/certificate"
	"github.com/sirupsen/logrus"
)

// Input input parameters
type Input struct {
	DomainName              string   `json:"domain-name"`
	SubjectAlternativeNames []string `json:"subject-alternative-names"`
}

// Lambda interface
type Lambda interface {
	Handler(Input) error
}

// ApproverLambda lambda struct
type ApproverLambda struct {
	logger *logrus.Logger
	cert   certificate.Certificate
}

// New creates a new lambda struct
func New(cert certificate.Certificate, logger *logrus.Logger) *ApproverLambda {
	return &ApproverLambda{
		cert:   cert,
		logger: logger,
	}
}

// Handler lambda request handler
func (a *ApproverLambda) Handler(input Input) error {
	return a.cert.Request(input.DomainName, input.SubjectAlternativeNames)
}
