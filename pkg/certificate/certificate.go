package certificate

import (
	"fmt"
	"strings"

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

// Approve approve given certificate
func (a *AcmCertificate) Approve(certificateArn string, ttl int64) error {
	a.logger.WithField("certificateArn", certificateArn).Debug("received request to approve certificate")
	res, err := a.acmvsvc.DescribeCertificate(&acm.DescribeCertificateInput{
		CertificateArn: aws.String(certificateArn),
	})

	if err != nil {
		return err
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
