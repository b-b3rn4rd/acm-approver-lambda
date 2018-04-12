package lambda

import (
	"context"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/b-b3rn4rd/acm-approver-lambda/pkg/certificate"
	"github.com/sirupsen/logrus"
)

// Lambda interface
type Lambda interface {
	Handler(ctx context.Context, event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error)
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
func (a *ApproverLambda) Handler(ctx context.Context, event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {
	subjectAlternativeNames := []string{}

	if subjectAlternativeNamesInput, ok := event.ResourceProperties["SubjectAlternativeNames"].([]interface{}); ok {
		for _, subjectAlternativeName := range subjectAlternativeNamesInput {
			subjectAlternativeNames = append(subjectAlternativeNames, subjectAlternativeName.(string))
		}
	}
	switch event.RequestType {
	case cfn.RequestDelete:
		if event.PhysicalResourceID != "" {
			err = a.cert.Delete(event.PhysicalResourceID)
		}
		physicalResourceID = event.PhysicalResourceID
	case cfn.RequestCreate:
		physicalResourceID, err = a.cert.Request(
			event.ResourceProperties["DomainName"].(string),
			subjectAlternativeNames,
		)
	default:
		physicalResourceID = event.PhysicalResourceID
	}

	return
}
