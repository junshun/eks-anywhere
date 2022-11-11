package awscred_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/eks-anywhere/internal/test"
	"github.com/aws/eks-anywhere/pkg/awscred"
	"github.com/aws/eks-anywhere/pkg/config"
)

const (
	wantCredSecretContent         = "testdata/want-aws-cred-secret.yaml"
	wantCredProviderConfigContent = "testdata/want-cred-provider-config.yaml"
	testAwsAccessKey              = "test_access_key"
	testAwsSecretKey              = "test_secret_key"
	testAwsRegion                 = "test_region"
)

func TestGenerateAwsConfigPackagesSuccess(t *testing.T) {
	awsCred := awscred.NewAwsCred()
	err := setAwsConfigEnvVars(t)
	if err != nil {
		t.Fatalf("setAwsConfigEnvVars()\n error = %v\n wantErr = nil", err)
	}
	gotFileContent, err := awsCred.GenerateAwsConfigPackages()
	if err != nil {
		t.Fatalf("awsecrcred.GenerateAwsConfigPackages()\n error = %v\n wantErr = nil", err)
	}

	test.AssertContentToFile(t, string(gotFileContent), wantCredSecretContent)
}

func TestGenerateAwsConfigPackagesFail(t *testing.T) {
	wantErr := fmt.Errorf("missing credentials")
	awsCred := awscred.NewAwsCred()

	_, err := awsCred.GenerateAwsConfigPackages()
	if !reflect.DeepEqual(err, wantErr) {
		t.Fatalf("error = %v\n wantErr = %v", err, wantErr)
	}
}

func TestGenerateCredProviderConfigSuccess(t *testing.T) {
	awsCred := awscred.NewAwsCred()

	gotFileContent, err := awsCred.GenerateCredProviderConfig()
	if err != nil {
		t.Fatalf("awsecrcred.GenerateAwsConfigPackages()\n error = %v\n wantErr = nil", err)
	}

	test.AssertContentToFile(t, string(gotFileContent), wantCredProviderConfigContent)
}

func setAwsConfigEnvVars(t *testing.T) error {
	t.Setenv(config.EksaAccessKeyIdEnv, testAwsAccessKey)
	t.Setenv(config.EksaSecretAccessKeyEnv, testAwsSecretKey)
	t.Setenv(config.EksaRegionEnv, testAwsRegion)

	return nil
}
