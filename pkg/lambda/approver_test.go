package lambda_test

import (
	"testing"

	"github.com/b-b3rn4rd/acm-approver-lambda/mocks"
	"github.com/b-b3rn4rd/acm-approver-lambda/pkg/lambda"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestApproverLambda(t *testing.T) {
	t.Run(" lambda calls Approve method with input", func(t *testing.T) {
		logger, _ := test.NewNullLogger()
		s := &mocks.Certificate{}

		s.On("Approve", "1", int64(300)).Return(nil)

		l := lambda.New(s, logger)

		err := l.Handler(lambda.Input{
			CertificateArn: "1",
			TTL:            int64(300)})

		assert.Nil(t, err)
	})

}
