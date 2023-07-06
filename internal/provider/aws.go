package provider

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type AWSCreds struct {
	AccessKeyId  string
	AccessKey    string
	SessionToken string
}

type Provider struct {
	k8sClient client.Client
}

func NewProvider(k8sClient client.Client) Provider {
	return Provider{
		k8sClient: k8sClient,
	}
}

func (p Provider) GetCreds(ctx context.Context, providerConfig commonv1alpha1.ProviderConfig) (AWSCreds, error) {
	var stsclient *sts.Client
	var cfg aws.Config
	var err error

	if providerConfig.Spec.Source != "SECRET" {
		var secret v1.Secret

		err := p.k8sClient.Get(ctx, types.NamespacedName{
			Name:      providerConfig.Spec.SecretRef.Name,
			Namespace: providerConfig.Spec.SecretRef.Namespace,
		}, &secret)
		if err != nil {
			return AWSCreds{}, err
		}

		cfg, err = config.LoadDefaultConfig(context.TODO(), func(opts *config.LoadOptions) error {
			opts.Region = providerConfig.Spec.AWSConfig.Region
			opts.Credentials = credentials.NewStaticCredentialsProvider(
				string(secret.Data["aws_access_key_id"]),
				string(secret.Data["aws_secret_access_key"]),
				"",
			)
			return nil
		})

		if err != nil {
			return AWSCreds{}, err
		}

	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO(), func(opts *config.LoadOptions) error {
			opts.Region = providerConfig.Spec.AWSConfig.Region
			return nil
		})
		if err != nil {
			return AWSCreds{}, err
		}

	}

	if providerConfig.Spec.AWSConfig.Role != "" {
		stsclient = sts.NewFromConfig(cfg)
		credsCache := stscreds.NewAssumeRoleProvider(stsclient, providerConfig.Spec.AWSConfig.Role)
		cfg, err = config.LoadDefaultConfig(
			ctx,
			config.WithRegion(providerConfig.Spec.AWSConfig.Region),
			config.WithCredentialsProvider(aws.NewCredentialsCache(credsCache)),
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
