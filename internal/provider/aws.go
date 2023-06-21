package provider

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
)

type AWSCreds struct {
	AccessKeyId  string
	AccessKey    string
	SessionToken string
}

func GetCreds(ctx context.Context, awsProviderConfig commonv1alpha1.AWSProviderConfig) (AWSCreds, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(opts *config.LoadOptions) error {
		// opts.Region = "us-east-1"
		return nil
	})
	if err != nil {
		return AWSCreds{}, err
	}

	if awsProviderConfig.Source != "CurrentAccount" {

		stsclient := sts.NewFromConfig(cfg)
		cfg, err = config.LoadDefaultConfig(
			// ctx, config.WithRegion("{aws-region}"),
			ctx,
			config.WithCredentialsProvider(aws.NewCredentialsCache(
				stscreds.NewAssumeRoleProvider(
					stsclient,
					awsProviderConfig.Source,
				)),
			),
		)

		if err != nil {
			return AWSCreds{}, err
		}
	}

	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return AWSCreds{}, err
	}

	return AWSCreds{
		AccessKeyId:  creds.AccessKeyID,
		AccessKey:    creds.SecretAccessKey,
		SessionToken: creds.SessionToken,
	}, nil
}
