// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/awsiamauth/awsiamauth.go

// Package awsiamauth_test is a generated GoMock package.
package awsiamauth_test

import (
	context "context"
	reflect "reflect"

	types "github.com/aws/eks-anywhere/pkg/types"
	gomock "github.com/golang/mock/gomock"
)

// MockKubernetesClient is a mock of KubernetesClient interface.
type MockKubernetesClient struct {
	ctrl     *gomock.Controller
	recorder *MockKubernetesClientMockRecorder
}

// MockKubernetesClientMockRecorder is the mock recorder for MockKubernetesClient.
type MockKubernetesClientMockRecorder struct {
	mock *MockKubernetesClient
}

// NewMockKubernetesClient creates a new mock instance.
func NewMockKubernetesClient(ctrl *gomock.Controller) *MockKubernetesClient {
	mock := &MockKubernetesClient{ctrl: ctrl}
	mock.recorder = &MockKubernetesClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKubernetesClient) EXPECT() *MockKubernetesClientMockRecorder {
	return m.recorder
}

// ApplyKubeSpecFromBytes mocks base method.
func (m *MockKubernetesClient) ApplyKubeSpecFromBytes(ctx context.Context, cluster *types.Cluster, data []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ApplyKubeSpecFromBytes", ctx, cluster, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// ApplyKubeSpecFromBytes indicates an expected call of ApplyKubeSpecFromBytes.
func (mr *MockKubernetesClientMockRecorder) ApplyKubeSpecFromBytes(ctx, cluster, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplyKubeSpecFromBytes", reflect.TypeOf((*MockKubernetesClient)(nil).ApplyKubeSpecFromBytes), ctx, cluster, data)
}

// GetApiServerUrl mocks base method.
func (m *MockKubernetesClient) GetApiServerUrl(ctx context.Context, cluster *types.Cluster) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetApiServerUrl", ctx, cluster)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetApiServerUrl indicates an expected call of GetApiServerUrl.
func (mr *MockKubernetesClientMockRecorder) GetApiServerUrl(ctx, cluster interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetApiServerUrl", reflect.TypeOf((*MockKubernetesClient)(nil).GetApiServerUrl), ctx, cluster)
}

// GetClusterCATlsCert mocks base method.
func (m *MockKubernetesClient) GetClusterCATlsCert(ctx context.Context, clusterName string, cluster *types.Cluster, namespace string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClusterCATlsCert", ctx, clusterName, cluster, namespace)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClusterCATlsCert indicates an expected call of GetClusterCATlsCert.
func (mr *MockKubernetesClientMockRecorder) GetClusterCATlsCert(ctx, clusterName, cluster, namespace interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClusterCATlsCert", reflect.TypeOf((*MockKubernetesClient)(nil).GetClusterCATlsCert), ctx, clusterName, cluster, namespace)
}