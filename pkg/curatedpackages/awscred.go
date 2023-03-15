package curatedpackages

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/aws/eks-anywhere/pkg/config"
	"github.com/aws/eks-anywhere/pkg/constants"
	"github.com/aws/eks-anywhere/pkg/templater"
)

const (
	awsProfile = "default"
)

//go:embed config/aws-cred-secret.yaml
var awsEcrCredSecretTemplate string

//go:embed config/credential-package.yaml
var credPackageTemplates string

type AwsCred struct {
}

func NewAwsCred() *AwsCred {
	return &AwsCred{}
}

func (a *AwsCred) GenerateAwsConfigSecret() ([]byte, error) {
	eksaAccessKeyId, eksaSecretKey, eksaRegion := os.Getenv(config.EksaAccessKeyIdEnv),
		os.Getenv(config.EksaSecretAccessKeyEnv),
		os.Getenv(config.EksaRegionEnv)

	if eksaAccessKeyId == "" || eksaSecretKey == "" || eksaRegion == "" {
		return nil, fmt.Errorf("missing credentials")
	}

	awsConfig := fmt.Sprintf(
		"[%s]\n"+
			"aws_access_key_id=%s\n"+
			"aws_secret_access_key=%s\n"+
			"region=%s", awsProfile, eksaAccessKeyId, eksaSecretKey, eksaRegion)

	awsConfigBytes := base64.StdEncoding.EncodeToString([]byte(awsConfig))

	values := map[string]string{
		"namespace":      constants.EksaPackagesName,
		"awsConfigBytes": awsConfigBytes,
	}

	awsConfigSecret, err := templater.Execute(awsEcrCredSecretTemplate, values)
	if err != nil {
		return nil, fmt.Errorf("generating aws-config secret: %v", err)
	}
	return awsConfigSecret, nil
}

func (a *AwsCred) GenerateCredentialPackageConfig(clusterName string, sourceRegistry string) ([]byte, error) {
	values := map[string]string{
		"clusterName":    clusterName,
		"sourceRegistry": sourceRegistry,
		"profile":        awsProfile,
	}
	credPackageConfig, err := templater.Execute(credPackageTemplates, values)
	if err != nil {
		return nil, fmt.Errorf("generating credential package configuratio %v", err)
	}
	return credPackageConfig, nil
}
