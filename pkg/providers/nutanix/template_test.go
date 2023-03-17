package nutanix

import (
	_ "embed"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"

	"github.com/aws/eks-anywhere/internal/test"
	anywherev1 "github.com/aws/eks-anywhere/pkg/api/v1alpha1"
	"github.com/aws/eks-anywhere/pkg/constants"
)

//go:embed testdata/eksa-cluster.yaml
var nutanixClusterConfigSpec string

//go:embed testdata/datacenterConfig.yaml
var nutanixDatacenterConfigSpec string

//go:embed testdata/machineConfig.yaml
var nutanixMachineConfigSpec string

//go:embed testdata/machineConfig_project.yaml
var nutanixMachineConfigSpecWithProject string

func fakemarshal(v interface{}) ([]byte, error) {
	return []byte{}, errors.New("marshalling failed")
}

func restoremarshal(replace func(v interface{}) ([]byte, error)) {
	jsonMarshal = replace
}

func TestNewNutanixTemplateBuilder(t *testing.T) {
	dcConf, machineConf, workerConfs := minimalNutanixConfigSpec(t)

	t.Setenv(constants.EksaNutanixUsernameKey, "admin")
	t.Setenv(constants.EksaNutanixPasswordKey, "password")
	creds := GetCredsFromEnv()
	builder := NewNutanixTemplateBuilder(&dcConf.Spec, &machineConf.Spec, &machineConf.Spec, workerConfs, creds, time.Now)
	assert.NotNil(t, builder)

	buildSpec := test.NewFullClusterSpec(t, "testdata/eksa-cluster.yaml")

	cpSpec, err := builder.GenerateCAPISpecControlPlane(buildSpec)
	assert.NoError(t, err)
	assert.NotNil(t, cpSpec)

	workloadTemplateNames := map[string]string{
		"eksa-unit-test": "eksa-unit-test",
	}
	kubeadmconfigTemplateNames := map[string]string{
		"eksa-unit-test": "eksa-unit-test",
	}
	workerSpec, err := builder.GenerateCAPISpecWorkers(buildSpec, workloadTemplateNames, kubeadmconfigTemplateNames)
	assert.NoError(t, err)
	assert.NotNil(t, workerSpec)

	secretSpec, err := builder.GenerateCAPISpecSecret(buildSpec)
	assert.NoError(t, err)
	assert.NotNil(t, secretSpec)
	expectedSecret, err := os.ReadFile("testdata/templated_secret.yaml")
	require.NoError(t, err)
	assert.Equal(t, expectedSecret, secretSpec)

	secretSpec, err = builder.GenerateEKSASpecSecret(buildSpec)
	assert.NoError(t, err)
	assert.NotNil(t, secretSpec)
	expectedSecret, err = os.ReadFile("testdata/templated_secret_eksa.yaml")
	require.NoError(t, err)
	assert.Equal(t, expectedSecret, secretSpec)
}

func TestNewNutanixTemplateBuilderGenerateSpecSecretFailure(t *testing.T) {
	storedMarshal := jsonMarshal
	jsonMarshal = fakemarshal
	defer restoremarshal(storedMarshal)

	t.Setenv(constants.NutanixUsernameKey, "admin")
	t.Setenv(constants.NutanixPasswordKey, "password")
	creds := GetCredsFromEnv()
	builder := NewNutanixTemplateBuilder(nil, nil, nil, nil, creds, time.Now)
	assert.NotNil(t, builder)

	buildSpec := test.NewFullClusterSpec(t, "testdata/eksa-cluster.yaml")

	secretSpec, err := builder.GenerateCAPISpecSecret(buildSpec)
	assert.Nil(t, secretSpec)
	assert.Error(t, err)

	secretSpec, err = builder.GenerateEKSASpecSecret(buildSpec)
	assert.Nil(t, secretSpec)
	assert.Error(t, err)
}

func TestNutanixTemplateBuilderGenerateCAPISpecForCreateWithAutoscalingConfiguration(t *testing.T) {
	dcConf, machineConf, workerConfs := minimalNutanixConfigSpec(t)

	t.Setenv(constants.NutanixUsernameKey, "admin")
	t.Setenv(constants.NutanixPasswordKey, "password")
	creds := GetCredsFromEnv()

	builder := NewNutanixTemplateBuilder(&dcConf.Spec, &machineConf.Spec, &machineConf.Spec, workerConfs, creds, time.Now)
	assert.NotNil(t, builder)

	buildSpec := test.NewFullClusterSpec(t, "testdata/eksa-cluster-autoscaler.yaml")

	workloadTemplateNames := map[string]string{
		"eksa-unit-test": "eksa-unit-test",
	}
	kubeadmconfigTemplateNames := map[string]string{
		"eksa-unit-test": "eksa-unit-test",
	}
	workerSpec, err := builder.GenerateCAPISpecWorkers(buildSpec, workloadTemplateNames, kubeadmconfigTemplateNames)
	assert.NoError(t, err)
	expectedWorkerSpec, err := os.ReadFile("testdata/expected_results_autoscaling_md.yaml")
	require.NoError(t, err)
	assert.Equal(t, string(workerSpec), string(expectedWorkerSpec))
}

func TestNewNutanixTemplateBuilderOIDCConfig(t *testing.T) {
	dcConf, machineConf, workerConfs := minimalNutanixConfigSpec(t)

	t.Setenv(constants.EksaNutanixUsernameKey, "admin")
	t.Setenv(constants.EksaNutanixPasswordKey, "password")
	creds := GetCredsFromEnv()
	builder := NewNutanixTemplateBuilder(&dcConf.Spec, &machineConf.Spec, &machineConf.Spec, workerConfs, creds, time.Now)
	assert.NotNil(t, builder)

	buildSpec := test.NewFullClusterSpec(t, "testdata/eksa-cluster-oidc.yaml")

	cpSpec, err := builder.GenerateCAPISpecControlPlane(buildSpec)
	assert.NoError(t, err)
	assert.NotNil(t, cpSpec)

	expectedControlPlaneSpec, err := os.ReadFile("testdata/expected_results_oidc.yaml")
	require.NoError(t, err)
	assert.Equal(t, expectedControlPlaneSpec, cpSpec)
}

func TestNewNutanixTemplateBuilderRegistryMirrorConfig(t *testing.T) {
	t.Setenv(constants.RegistryUsername, "username")
	t.Setenv(constants.RegistryPassword, "password")
	dcConf, machineConf, workerConfs := minimalNutanixConfigSpec(t)

	t.Setenv(constants.EksaNutanixUsernameKey, "admin")
	t.Setenv(constants.EksaNutanixPasswordKey, "password")
	creds := GetCredsFromEnv()
	builder := NewNutanixTemplateBuilder(&dcConf.Spec, &machineConf.Spec, &machineConf.Spec, workerConfs, creds, time.Now)
	assert.NotNil(t, builder)

	buildSpec := test.NewFullClusterSpec(t, "testdata/eksa-cluster-registry-mirror.yaml")

	cpSpec, err := builder.GenerateCAPISpecControlPlane(buildSpec)
	assert.NoError(t, err)
	assert.NotNil(t, cpSpec)

	expectedControlPlaneSpec, err := os.ReadFile("testdata/expected_results_registry_mirror.yaml")
	require.NoError(t, err)
	assert.Equal(t, expectedControlPlaneSpec, cpSpec)

	workloadTemplateNames := map[string]string{
		"eksa-unit-test": "eksa-unit-test",
	}
	kubeadmconfigTemplateNames := map[string]string{
		"eksa-unit-test": "eksa-unit-test",
	}
	workerSpec, err := builder.GenerateCAPISpecWorkers(buildSpec, workloadTemplateNames, kubeadmconfigTemplateNames)
	assert.NoError(t, err)
	assert.NotNil(t, workerSpec)

	expectedWorkersSpec, err := os.ReadFile("testdata/expected_results_registry_mirror_md.yaml")
	require.NoError(t, err)
	assert.Equal(t, expectedWorkersSpec, workerSpec)
}

func TestNewNutanixTemplateBuilderRegistryMirrorConfigNoRegistryCredsSet(t *testing.T) {
	dcConf, machineConf, workerConfs := minimalNutanixConfigSpec(t)

	t.Setenv(constants.EksaNutanixUsernameKey, "admin")
	t.Setenv(constants.EksaNutanixPasswordKey, "password")
	creds := GetCredsFromEnv()
	builder := NewNutanixTemplateBuilder(&dcConf.Spec, &machineConf.Spec, &machineConf.Spec, workerConfs, creds, time.Now)
	assert.NotNil(t, builder)

	buildSpec := test.NewFullClusterSpec(t, "testdata/eksa-cluster-registry-mirror.yaml")

	_, err := builder.GenerateCAPISpecControlPlane(buildSpec)
	assert.Error(t, err)

	workloadTemplateNames := map[string]string{
		"eksa-unit-test": "eksa-unit-test",
	}
	kubeadmconfigTemplateNames := map[string]string{
		"eksa-unit-test": "eksa-unit-test",
	}
	_, err = builder.GenerateCAPISpecWorkers(buildSpec, workloadTemplateNames, kubeadmconfigTemplateNames)
	assert.Error(t, err)
}

func TestNewNutanixTemplateBuilderProject(t *testing.T) {
	t.Setenv(constants.EksaNutanixUsernameKey, "admin")
	t.Setenv(constants.EksaNutanixPasswordKey, "password")
	creds := GetCredsFromEnv()

	dcConf, _, _ := minimalNutanixConfigSpec(t)
	machineConf := &anywherev1.NutanixMachineConfig{}
	err := yaml.Unmarshal([]byte(nutanixMachineConfigSpecWithProject), machineConf)
	require.NoError(t, err)

	workerConfs := map[string]anywherev1.NutanixMachineConfigSpec{
		"eksa-unit-test": machineConf.Spec,
	}
	builder := NewNutanixTemplateBuilder(&dcConf.Spec, &machineConf.Spec, &machineConf.Spec, workerConfs, creds, time.Now)
	assert.NotNil(t, builder)

	buildSpec := test.NewFullClusterSpec(t, "testdata/eksa-cluster-project.yaml")
	cpSpec, err := builder.GenerateCAPISpecControlPlane(buildSpec)
	assert.NoError(t, err)
	assert.NotNil(t, cpSpec)

	expectedControlPlaneSpec, err := os.ReadFile("testdata/expected_results_project.yaml")
	require.NoError(t, err)
	assert.Equal(t, expectedControlPlaneSpec, cpSpec)

	workloadTemplateNames := map[string]string{
		"eksa-unit-test": "eksa-unit-test",
	}
	kubeadmconfigTemplateNames := map[string]string{
		"eksa-unit-test": "eksa-unit-test",
	}
	workerSpec, err := builder.GenerateCAPISpecWorkers(buildSpec, workloadTemplateNames, kubeadmconfigTemplateNames)
	assert.NoError(t, err)
	assert.NotNil(t, workerSpec)

	expectedWorkersSpec, err := os.ReadFile("testdata/expected_results_project_md.yaml")
	require.NoError(t, err)
	assert.Equal(t, expectedWorkersSpec, workerSpec)
}

func minimalNutanixConfigSpec(t *testing.T) (*anywherev1.NutanixDatacenterConfig, *anywherev1.NutanixMachineConfig, map[string]anywherev1.NutanixMachineConfigSpec) {
	dcConf := &anywherev1.NutanixDatacenterConfig{}
	err := yaml.Unmarshal([]byte(nutanixDatacenterConfigSpec), dcConf)
	require.NoError(t, err)

	machineConf := &anywherev1.NutanixMachineConfig{}
	err = yaml.Unmarshal([]byte(nutanixMachineConfigSpec), machineConf)
	require.NoError(t, err)

	workerConfs := map[string]anywherev1.NutanixMachineConfigSpec{
		"eksa-unit-test": machineConf.Spec,
	}

	return dcConf, machineConf, workerConfs
}
