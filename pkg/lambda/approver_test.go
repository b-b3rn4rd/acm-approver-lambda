package lambda_test

import (
	"testing"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/b-b3rn4rd/acm-approver-lambda/mocks"
	"github.com/b-b3rn4rd/acm-approver-lambda/pkg/lambda"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestApproverLambdaWithAlternativeNames(t *testing.T) {
	t.Run(" lambda calls Request method with input", func(t *testing.T) {
		logger, _ := test.NewNullLogger()
		s := &mocks.Certificate{}

		s.On("Request", "www.example.net", []string{"test.example.net"}).Return(nil)

		l := lambda.New(s, logger)

		_, _, err := l.Handler(nil, cfn.Event{
			ResourceType: string(cfn.RequestCreate),
			ResourceProperties: map[string]interface{}{
				"DomainName":              "www.example.net",
				"SubjectAlternativeNames": []interface{}{"test.example.net"},
			},
		})

		assert.Nil(t, err)
	})

}

func TestApproverLambdaWithOutAlternativeNames(t *testing.T) {
	t.Run(" lambda calls Request method with input", func(t *testing.T) {
		logger, _ := test.NewNullLogger()
		s := &mocks.Certificate{}

		s.On("Request", "www.example.net", []string{}).Return(nil)

		l := lambda.New(s, logger)

		_, _, err := l.Handler(nil, cfn.Event{
			ResourceType: string(cfn.RequestCreate),
			ResourceProperties: map[string]interface{}{
				"DomainName": "www.example.net",
			},
		})

		assert.Nil(t, err)
	})

}

func TestApproverLambdaDelete(t *testing.T) {
	t.Run(" lambda calls Delete method with input", func(t *testing.T) {
		logger, _ := test.NewNullLogger()
		s := &mocks.Certificate{}

		s.On("Delete", "abc").Return(nil)

		l := lambda.New(s, logger)

		_, _, err := l.Handler(nil, cfn.Event{
			ResourceType:       string(cfn.RequestDelete),
			PhysicalResourceID: "abc",
			ResourceProperties: map[string]interface{}{
				"DomainName": "www.example.net",
			},
		})

		assert.Nil(t, err)
	})

}
