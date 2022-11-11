package awscred

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
	awsProfile    = "eksa-packages"
	awsConfigFile = "/root/aws-ecr-creds"
	awsPath       = "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
	awsHome       = "/root"
)

//go:embed config/aws-cred-secret.yaml
var awsEcrCredSecretTemplate string

//go:embed config/cred-provider-config.yaml
var credProviderConfigTemplate string

type AwsCred struct {
}

func NewAwsCred() *AwsCred {
	return &AwsCred{}
}

func (a *AwsCred) GenerateAwsConfigPackages() ([]byte, error) {
	eksaAccessKeyId, eksaSecretKey, eksaRegion := os.Getenv(config.EksaAccessKeyIdEnv),
		os.Getenv(config.EksaSecretAccessKeyEnv),
		os.Getenv(config.EksaRegionEnv)

	if eksaAccessKeyId == "" || eksaSecretKey == "" || eksaRegion == "" {
		return nil, fmt.Errorf("missing credentials")
	}

	awsConfig := fmt.Sprintf(
		"[profile %s]\n"+
			"aws_access_key_id=%s\n"+
			"aws_secret_access_key=%s\n"+
			"region=%s",
		awsProfile, eksaAccessKeyId, eksaSecretKey, eksaRegion)

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

func (a *AwsCred) GenerateCredProviderConfig() ([]byte, error) {
	credProviderConfig := fmt.Sprintf("apiVersion: kubelet.config.k8s.io/v1alpha1\n"+
		"kind: CredentialProviderConfig\n"+
		"providers:\n"+
		"  - name: ecr-credential-provider\n"+
		"    matchImages: \n"+
		"      - \"*.dkr.ecr.*.amazonaws.com\"\n"+
		"    defaultCacheDuration: \"12h\"\n"+
		"    apiVersion: credentialprovider.kubelet.k8s.io/v1alpha1\n"+
		"    env:\n"+
		"      - name: AWS_PROFILE\n"+
		"        value: %s\n"+
		"      - name: AWS_CONFIG_FILE\n"+
		"        value: %s\n"+
		"      - name: PATH\n"+
		"        value: %s\n"+
		"      - name: HOME\n"+
		"        value: %s",
		awsProfile, awsConfigFile, awsPath, awsHome)

	credProviderConfigBytes := base64.StdEncoding.EncodeToString([]byte(credProviderConfig))
	data := map[string]string{
		"namespace":               constants.EksaSystemNamespace,
		"credProviderConfigBytes": credProviderConfigBytes,
	}

	credProviderConfigSecret, err := templater.Execute(credProviderConfigTemplate, data)
	if err != nil {
		return nil, fmt.Errorf("generating aws-config secret: %v", err)
	}
	return credProviderConfigSecret, nil
}
