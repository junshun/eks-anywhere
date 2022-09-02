package awsecrcred_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/eks-anywhere/internal/test"
	"github.com/aws/eks-anywhere/pkg/awsecrcred"
	"github.com/aws/eks-anywhere/pkg/config"
)

const (
	wantSecretContent = "testdata/want-aws-ecr-cred-secret.yaml"
	testAwsAccessKey  = "test_access_key"
	testAwsSecretKey  = "test_secret_key"
	testAwsRegion     = "test_region"
)

func TestGenerateAwsConfigSuccess(t *testing.T) {
	awsEcrCred := awsecrcred.NewAwsEcrCred()
	err := setAwsConfigEnvVars(t)
	if err != nil {
		t.Fatalf("setAwsConfigEnvVars()\n error = %v\n wantErr = nil", err)
	}
	gotFileContent, err := awsEcrCred.GenerateAwsConfig()
	if err != nil {
		t.Fatalf("awsecrcred.GenerateAwsConfig()\n error = %v\n wantErr = nil", err)
	}

	test.AssertContentToFile(t, string(gotFileContent), wantSecretContent)
}

func TestGenerateAwsConfigFail(t *testing.T) {
	wantErr := fmt.Errorf("missing credentials")
	awsEcrCred := awsecrcred.NewAwsEcrCred()

	_, err := awsEcrCred.GenerateAwsConfig()
	if !reflect.DeepEqual(err, wantErr) {
		t.Fatalf("error = %v\n wantErr = %v", err, wantErr)
	}
}

func setAwsConfigEnvVars(t *testing.T) error {
	t.Setenv(config.EksaAccessKeyIdEnv, testAwsAccessKey)
	t.Setenv(config.EksaSecretAcessKeyEnv, testAwsSecretKey)
	t.Setenv(config.EksaRegionEnv, testAwsRegion)

	return nil
}
