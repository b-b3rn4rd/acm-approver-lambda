package certificate

import (
	"fmt"
	"strings"

	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/acm/acmiface"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Certificate interface
type Certificate interface {
	Approve(string, int64) error
	Request(string, []string) (string, error)
	Delete(certificateArn string) error
}

// AcmCertificate storage struct
type AcmCertificate struct {
	acmvsvc acmiface.ACMAPI
	r53svc  route53iface.Route53API
	logger  *logrus.Logger
}

// New create new acm approve
func New(acmvsvc acmiface.ACMAPI, r53svc route53iface.Route53API, logger *logrus.Logger) *AcmCertificate {
	return &AcmCertificate{
		acmvsvc: acmvsvc,
		r53svc:  r53svc,
		logger:  logger,
	}
}

// Delete delete certificate
func (a *AcmCertificate) Delete(certificateArn string) error {
	a.logger.WithField("certificateArn", certificateArn).Debug("Deleting certificate")

	_, err := a.acmvsvc.DeleteCertificate(&acm.DeleteCertificateInput{
		CertificateArn: aws.String(certificateArn)})

	return err
}

// Request request required certificate
func (a *AcmCertificate) Request(domainName string, subjectAlternativeNames []string) (string, error) {
	a.logger.WithField("domainName", domainName).WithField("alternativeNames", subjectAlternativeNames).Debug("Requesting certificate")

	input := &acm.RequestCertificateInput{
		DomainName:       aws.String(domainName),
		ValidationMethod: aws.String(acm.ValidationMethodDns),
	}

	if len(subjectAlternativeNames) > 0 {
		input.SubjectAlternativeNames = aws.StringSlice(subjectAlternativeNames)
	}

	res, err := a.acmvsvc.RequestCertificate(input)

	if err != nil {
		return "", err
	}

	return *res.CertificateArn, a.Approve(*res.CertificateArn, 300)
}

// Approve approve given certificate
func (a *AcmCertificate) Approve(certificateArn string, ttl int64) error {
	a.logger.WithField("certificateArn", certificateArn).Debug("received request to approve certificate")
	polls := 5

	var res *acm.DescribeCertificateOutput
	var err error

	for i := 1; i <= polls; i++ {
		polls--
		res, err = a.acmvsvc.DescribeCertificate(&acm.DescribeCertificateInput{
			CertificateArn: aws.String(certificateArn),
		})

		if err != nil {
			return err
		}

		if res.Certificate.DomainValidationOptions[0].ResourceRecord != nil {
			a.logger.WithField("certificateArn", certificateArn).Debug("certificate contains confirmation record")
			break
		}

		a.logger.WithField("certificateArn", certificateArn).Debug("certificate does not contain confirmation record, doing another poll")
		time.Sleep(5 * time.Second)
	}

	hostedZoneID, err := a.getHostedZoneID(res.Certificate.DomainName)

	if err != nil {
		return errors.Wrap(err, "error while fetching hosted zone")
	}

	a.logger.WithField("hostedZoneID", hostedZoneID).Debug("found hosted zone id for certificate")

	for _, record := range res.Certificate.DomainValidationOptions {
		_, err = a.r53svc.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
			HostedZoneId: aws.String(hostedZoneID),
			ChangeBatch: &route53.ChangeBatch{Changes: []*route53.Change{
				{
					Action: aws.String(route53.ChangeActionUpsert),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: record.ResourceRecord.Name,
						Type: record.ResourceRecord.Type,
						TTL:  aws.Int64(ttl),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: record.ResourceRecord.Value,
							},
						},
					},
				},
			}},
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (a *AcmCertificate) getHostedZoneID(domainName *string) (string, error) {
	a.logger.WithField("domainName", *domainName).Debug("searching hosted zone for certificate domain")

	res, err := a.r53svc.ListHostedZones(&route53.ListHostedZonesInput{})

	if err != nil {
		return "", err
	}

	longestMatch := 0
	hostedZoneID := ""

	longestCommonSuffix := func(a, b string) (i int) {
		for ; i < len(a) && i < len(b); i++ {
			if a[len(a)-1-i] != b[len(b)-1-i] {
				break
			}
		}
		return
	}

	for _, hostedZone := range res.HostedZones {
		hostedZoneName := strings.TrimRight(*hostedZone.Name, ".")

		if strings.HasSuffix(*domainName, hostedZoneName) {
			a.logger.WithField("hostedZoneName", hostedZoneName).Debug("found hosted zone with common suffix")
			l := longestCommonSuffix(*domainName, hostedZoneName)
			if l > longestMatch {
				a.logger.WithField("length", l).WithField("hostedZoneName", hostedZoneName).Debug("found a new best candidate")
				longestMatch = l
				hostedZoneID = *hostedZone.Id
			}
		}
	}
	if hostedZoneID == "" {
		return "", fmt.Errorf("cant find hosted zone for %s domain", *domainName)
	}

	return hostedZoneID, nil

}
