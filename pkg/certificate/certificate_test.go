package certificate_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/route53"
	mocks2 "github.com/b-b3rn4rd/acm-approver-lambda/mocks"
	"github.com/b-b3rn4rd/acm-approver-lambda/pkg/certificate"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestApprove(t *testing.T) {
	r53 := &mocks2.Route53API{}
	acmapi := &mocks2.ACMAPI{}

	listHostedZonesRes := func(input *route53.ListHostedZonesInput) *route53.ListHostedZonesOutput {
		return &route53.ListHostedZonesOutput{HostedZones: []*route53.HostedZone{
			{
				Id:   aws.String("333"),
				Name: aws.String(".net."),
			},
			{
				Id:   aws.String("111"),
				Name: aws.String("example.com."),
			},
			{
				Id:   aws.String("222"),
				Name: aws.String("example.net."),
			},
		}}
	}

	expectedInput := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String("222"),
		ChangeBatch: &route53.ChangeBatch{Changes: []*route53.Change{
			{
				Action: aws.String(route53.ChangeActionUpsert),
				ResourceRecordSet: &route53.ResourceRecordSet{
					Name: aws.String("www.example.net"),
					Type: aws.String("TEST"),
					TTL:  aws.Int64(300),
					ResourceRecords: []*route53.ResourceRecord{
						{
							Value: aws.String("secret"),
						},
					},
				},
			},
		}},
	}

	changeResourceRecordSetsRes := func(input *route53.ChangeResourceRecordSetsInput) *route53.ChangeResourceRecordSetsOutput {
		return &route53.ChangeResourceRecordSetsOutput{}
	}

	describeCertificateRes := func(input *acm.DescribeCertificateInput) *acm.DescribeCertificateOutput {
		return &acm.DescribeCertificateOutput{Certificate: &acm.CertificateDetail{
			DomainName: aws.String("www.example.net"),
			DomainValidationOptions: []*acm.DomainValidation{{
				ResourceRecord: &acm.ResourceRecord{
					Name:  aws.String("www.example.net"),
					Type:  aws.String("TEST"),
					Value: aws.String("secret")}}},
		}}
	}

	acmapi.On("DescribeCertificate", mock.AnythingOfType("*acm.DescribeCertificateInput")).Return(describeCertificateRes, nil)
	r53.On("ListHostedZones", mock.AnythingOfType("*route53.ListHostedZonesInput")).Return(listHostedZonesRes, nil)
	r53.On("ChangeResourceRecordSets", mock.MatchedBy(func(input *route53.ChangeResourceRecordSetsInput) bool {
		return assert.Equal(t, *expectedInput, *input)

	})).Return(changeResourceRecordSetsRes, nil)

	t.Run("Testing approval process", func(t *testing.T) {
		logger, _ := test.NewNullLogger()
		c := certificate.New(acmapi, r53, logger)
		c.Approve("abc", 300)
	})

}
