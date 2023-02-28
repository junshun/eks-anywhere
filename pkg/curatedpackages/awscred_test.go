package curatedpackages_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/eks-anywhere/internal/test"
	"github.com/aws/eks-anywhere/pkg/config"
	"github.com/aws/eks-anywhere/pkg/curatedpackages"
)

const (
	wantCredSecretContent    = "testdata/want-aws-cred-secret.yaml"
	wantCredPkgConfigContent = "testdata/want-credential-package.yaml"
	testAwsAccessKey         = "test_access_key"
	testAwsSecretKey         = "test_secret_key"
	testAwsRegion            = "test_region"
)

func TestGenerateAwsConfigSecretSuccess(t *testing.T) {
	awsCred := curatedpackages.NewAwsCred()
	err := setAwsConfigEnvVars(t)
	if err != nil {
		t.Fatalf("setAwsConfigEnvVars()\n error = %v\n wantErr = nil", err)
	}
	gotFileContent, err := awsCred.GenerateAwsConfigSecret()
	if err != nil {
		t.Fatalf("awsCred.GenerateAwsConfigSecret()\n error = %v\n wantErr = nil", err)
	}

	test.AssertContentToFile(t, string(gotFileContent), wantCredSecretContent)
}

func TestGenerateAwsConfigSecretFail(t *testing.T) {
	wantErr := fmt.Errorf("missing credentials")
	awsCred := curatedpackages.NewAwsCred()

	_, err := awsCred.GenerateAwsConfigSecret()
	if !reflect.DeepEqual(err, wantErr) {
		t.Fatalf("error = %v\n wantErr = %v", err, wantErr)
	}
}

func TestGenerateCredentialPackageConfigSuccess(t *testing.T) {
	clusterName := "billy"
	sourceRegistry := "test_registry/test"
	awsCred := curatedpackages.NewAwsCred()
	gotFileContent, err := awsCred.GenerateCredentialPackageConfig(clusterName, sourceRegistry)
	if err != nil {
		t.Fatalf("awsCred.GenerateCredentialPackageConfig()\n error = %v\n wantErr = nil", err)
	}

	test.AssertContentToFile(t, string(gotFileContent), wantCredPkgConfigContent)
}

func setAwsConfigEnvVars(t *testing.T) error {
	t.Setenv(config.EksaAccessKeyIdEnv, testAwsAccessKey)
	t.Setenv(config.EksaSecretAccessKeyEnv, testAwsSecretKey)
	t.Setenv(config.EksaRegionEnv, testAwsRegion)

	return nil
}
