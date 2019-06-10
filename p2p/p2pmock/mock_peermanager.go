// Code generated by MockGen. DO NOT EDIT.
// Source: peermanager.go

// Package p2pmock is a generated GoMock package.
package p2pmock

import (
	message "github.com/aergoio/aergo/message"
	p2pcommon "github.com/aergoio/aergo/p2p/p2pcommon"
	types "github.com/aergoio/aergo/types"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockPeerManager is a mock of PeerManager interface
type MockPeerManager struct {
	ctrl     *gomock.Controller
	recorder *MockPeerManagerMockRecorder
}

// MockPeerManagerMockRecorder is the mock recorder for MockPeerManager
type MockPeerManagerMockRecorder struct {
	mock *MockPeerManager
}

// NewMockPeerManager creates a new mock instance
func NewMockPeerManager(ctrl *gomock.Controller) *MockPeerManager {
	mock := &MockPeerManager{ctrl: ctrl}
	mock.recorder = &MockPeerManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPeerManager) EXPECT() *MockPeerManagerMockRecorder {
	return m.recorder
}

// Start mocks base method
func (m *MockPeerManager) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start
func (mr *MockPeerManagerMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockPeerManager)(nil).Start))
}

// Stop mocks base method
func (m *MockPeerManager) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop
func (mr *MockPeerManagerMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockPeerManager)(nil).Stop))
}

// SelfMeta mocks base method
func (m *MockPeerManager) SelfMeta() p2pcommon.PeerMeta {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelfMeta")
	ret0, _ := ret[0].(p2pcommon.PeerMeta)
	return ret0
}

// SelfMeta indicates an expected call of SelfMeta
func (mr *MockPeerManagerMockRecorder) SelfMeta() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelfMeta", reflect.TypeOf((*MockPeerManager)(nil).SelfMeta))
}

// SelfNodeID mocks base method
func (m *MockPeerManager) SelfNodeID() types.PeerID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelfNodeID")
	ret0, _ := ret[0].(types.PeerID)
	return ret0
}

// SelfNodeID indicates an expected call of SelfNodeID
func (mr *MockPeerManagerMockRecorder) SelfNodeID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelfNodeID", reflect.TypeOf((*MockPeerManager)(nil).SelfNodeID))
}

// AddNewPeer mocks base method
func (m *MockPeerManager) AddNewPeer(peer p2pcommon.PeerMeta) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddNewPeer", peer)
}

// AddNewPeer indicates an expected call of AddNewPeer
func (mr *MockPeerManagerMockRecorder) AddNewPeer(peer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewPeer", reflect.TypeOf((*MockPeerManager)(nil).AddNewPeer), peer)
}

// RemovePeer mocks base method
func (m *MockPeerManager) RemovePeer(peer p2pcommon.RemotePeer) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemovePeer", peer)
}

// RemovePeer indicates an expected call of RemovePeer
func (mr *MockPeerManagerMockRecorder) RemovePeer(peer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemovePeer", reflect.TypeOf((*MockPeerManager)(nil).RemovePeer), peer)
}

// UpdatePeerRole mocks base method
func (m *MockPeerManager) UpdatePeerRole(changes []p2pcommon.AttrModifier) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdatePeerRole", changes)
}

// UpdatePeerRole indicates an expected call of UpdatePeerRole
func (mr *MockPeerManagerMockRecorder) UpdatePeerRole(changes interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePeerRole", reflect.TypeOf((*MockPeerManager)(nil).UpdatePeerRole), changes)
}

// NotifyPeerAddressReceived mocks base method
func (m *MockPeerManager) NotifyPeerAddressReceived(arg0 []p2pcommon.PeerMeta) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "NotifyPeerAddressReceived", arg0)
}

// NotifyPeerAddressReceived indicates an expected call of NotifyPeerAddressReceived
func (mr *MockPeerManagerMockRecorder) NotifyPeerAddressReceived(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NotifyPeerAddressReceived", reflect.TypeOf((*MockPeerManager)(nil).NotifyPeerAddressReceived), arg0)
}

// GetPeer mocks base method
func (m *MockPeerManager) GetPeer(ID types.PeerID) (p2pcommon.RemotePeer, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPeer", ID)
	ret0, _ := ret[0].(p2pcommon.RemotePeer)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetPeer indicates an expected call of GetPeer
func (mr *MockPeerManagerMockRecorder) GetPeer(ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPeer", reflect.TypeOf((*MockPeerManager)(nil).GetPeer), ID)
}

// GetPeers mocks base method
func (m *MockPeerManager) GetPeers() []p2pcommon.RemotePeer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPeers")
	ret0, _ := ret[0].([]p2pcommon.RemotePeer)
	return ret0
}

// GetPeers indicates an expected call of GetPeers
func (mr *MockPeerManagerMockRecorder) GetPeers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPeers", reflect.TypeOf((*MockPeerManager)(nil).GetPeers))
}

// GetPeerAddresses mocks base method
func (m *MockPeerManager) GetPeerAddresses(noHidden, showSelf bool) []*message.PeerInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPeerAddresses", noHidden, showSelf)
	ret0, _ := ret[0].([]*message.PeerInfo)
	return ret0
}

// GetPeerAddresses indicates an expected call of GetPeerAddresses
func (mr *MockPeerManagerMockRecorder) GetPeerAddresses(noHidden, showSelf interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPeerAddresses", reflect.TypeOf((*MockPeerManager)(nil).GetPeerAddresses), noHidden, showSelf)
}

// GetPeerBlockInfos mocks base method
func (m *MockPeerManager) GetPeerBlockInfos() []types.PeerBlockInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPeerBlockInfos")
	ret0, _ := ret[0].([]types.PeerBlockInfo)
	return ret0
}

// GetPeerBlockInfos indicates an expected call of GetPeerBlockInfos
func (mr *MockPeerManagerMockRecorder) GetPeerBlockInfos() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPeerBlockInfos", reflect.TypeOf((*MockPeerManager)(nil).GetPeerBlockInfos))
}
