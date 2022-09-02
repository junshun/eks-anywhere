package awsecrcred

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/aws/eks-anywhere/pkg/config"
	"github.com/aws/eks-anywhere/pkg/constants"
	"github.com/aws/eks-anywhere/pkg/templater"
)

//go:embed config/aws-ecr-cred-secret.yaml
var awsEcrCredSecretTemplate string

type AwsEcrCred struct {
}

func NewAwsEcrCred() *AwsEcrCred {
	return &AwsEcrCred{}
}

func (a *AwsEcrCred) GenerateAwsConfig() ([]byte, error) {
	eksaAccessKeyId, eksaSecretKey, eksaRegion := os.Getenv(config.EksaAccessKeyIdEnv),
		os.Getenv(config.EksaSecretAcessKeyEnv),
		os.Getenv(config.EksaRegionEnv)

	if eksaAccessKeyId == "" || eksaSecretKey == "" || eksaRegion == "" {
		return nil, fmt.Errorf("missing credentials")
	}

	awsConfig := fmt.Sprintf(
		"[eksa-packages]\n"+
			"aws_access_key_id=%s\n"+
			"aws_secret_access_key=%s\n"+
			"region=%s",
		eksaAccessKeyId, eksaSecretKey, eksaRegion)

	awsConfigBytes := base64.StdEncoding.EncodeToString([]byte(awsConfig))

	data := map[string]string{
		"namespace":      constants.EksaSystemNamespace,
		"awsConfigBytes": awsConfigBytes,
	}

	awsConfigSecret, err := templater.Execute(awsEcrCredSecretTemplate, data)
	if err != nil {
		return nil, fmt.Errorf("generating aws-config secret: %v", err)
	}
	return awsConfigSecret, nil
}
