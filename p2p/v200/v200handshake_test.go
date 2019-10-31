/*
 * @file
 * @copyright defined in aergo/LICENSE.txt
 */

package v200

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/libp2p/go-libp2p-core/crypto"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/aergoio/aergo-lib/log"
	"github.com/aergoio/aergo/config"
	"github.com/aergoio/aergo/p2p/p2pcommon"
	"github.com/aergoio/aergo/p2p/p2pkey"
	"github.com/aergoio/aergo/p2p/p2pmock"
	"github.com/aergoio/aergo/p2p/p2putil"
	"github.com/aergoio/aergo/types"
	"github.com/golang/mock/gomock"
)

var (
	myChainID, theirChainID *types.ChainID
	theirChainBytes         []byte

	samplePeerID, _   = types.IDB58Decode("16Uiu2HAmFqptXPfcdaCdwipB2fhHATgKGVFVPehDAPZsDKSU7jRm")
	dummyBlockHash, _ = hex.DecodeString("4f461d85e869ade8a0544f8313987c33a9c06534e50c4ad941498299579bd7ac")
	dummyBlockID      = types.MustParseBlockID(dummyBlockHash)

	dummyBlockHeight uint64 = 100215

	dummyGenHash = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	diffGenesis  = []byte{0xff, 0xfe, 0xfd, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	sampleVersion = "v2.0.0-test"
)

type fakeChainID struct {
	genID    types.ChainID
	versions []uint64
}

func newFC(genID types.ChainID, vers ...uint64) fakeChainID {
	genID.Version = 0
	sort.Sort(BlkNoASC(vers))
	return fakeChainID{genID: genID, versions: vers}
}
func (f fakeChainID) getChainID(no types.BlockNo) *types.ChainID {
	cp := f.genID
	for i := len(f.versions) - 1; i >= 0; i-- {
		if f.versions[i] <= no {
			cp.Version = int32(i + 1)
			break
		}
	}
	return &cp
}

type BlkNoASC []uint64

func (a BlkNoASC) Len() int           { return len(a) }
func (a BlkNoASC) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BlkNoASC) Less(i, j int) bool { return a[i] < a[j] }

func init() {
	myChainID = types.NewChainID()
	myChainID.Magic = "itSmain1"

	theirChainID = types.NewChainID()
	theirChainID.Magic = "itsdiff2"
	theirChainBytes, _ = theirChainID.Bytes()

	sampleKeyFile := "../../test/sample.key"
	baseCfg := &config.BaseConfig{AuthDir: "test"}
	p2pCfg := &config.P2PConfig{NPKey: sampleKeyFile}
	p2pkey.InitNodeInfo(baseCfg, p2pCfg, "0.0.1-test", log.NewLogger("v200.test"))
}

func TestDeepEqual(t *testing.T) {
	b1, _ := myChainID.Bytes()
	b2 := make([]byte, len(b1), len(b1)<<1)
	copy(b2, b1)

	s1 := &types.Status{ChainID: b1}
	s2 := &types.Status{ChainID: b2}

	if !reflect.DeepEqual(s1, s2) {
		t.Errorf("byte slice cant do DeepEqual! %v, %v", b1, b2)
	}

}

func TestV200StatusHS_doForOutbound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := log.NewLogger("test")
	mockActor := p2pmock.NewMockActorService(ctrl)
	mockCA := p2pmock.NewMockChainAccessor(ctrl)

	fc := newFC(*myChainID, 10000, 20000, dummyBlockHeight+100)
	localChainID := *fc.getChainID(dummyBlockHeight)
	localChainBytes, _ := localChainID.Bytes()
	oldChainID := fc.getChainID(10000)
	oldChainBytes, _ := oldChainID.Bytes()
	newChainID := fc.getChainID(600000)
	newChainBytes, _ := newChainID.Bytes()

	diffBlockNo := dummyBlockHeight + 100000
	dummyMeta := p2pcommon.NewMetaWith1Addr(samplePeerID, "dummy.aergo.io", 7846, "v2.0.0")
	dummyMeta.Version = sampleVersion
	dummyAddr := dummyMeta.ToPeerAddress()
	dummyBlock := &types.Block{Hash: dummyBlockHash, Header: &types.BlockHeader{BlockNo: dummyBlockHeight}}
	mockActor.EXPECT().GetChainAccessor().Return(mockCA).AnyTimes()
	mockCA.EXPECT().GetBestBlock().Return(dummyBlock, nil).AnyTimes()

	dummyGenHash := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	diffGenesis := []byte{0xff, 0xfe, 0xfd, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	dummyStatusMsg := &types.Status{ChainID: localChainBytes, Sender: &dummyAddr, Genesis: dummyGenHash, BestBlockHash: dummyBlockHash, BestHeight: dummyBlockHeight}
	succResult := &p2pcommon.HandshakeResult{Meta: dummyMeta, BestBlockHash: dummyBlockID, BestBlockNo: dummyBlockHeight}
	diffGenesisStatusMsg := &types.Status{ChainID: localChainBytes, Sender: &dummyAddr, Genesis: diffGenesis, BestBlockHash: dummyBlockHash, BestHeight: dummyBlockHeight}
	nilGenesisStatusMsg := &types.Status{ChainID: localChainBytes, Sender: &dummyAddr, Genesis: nil, BestBlockHash: dummyBlockHash, BestHeight: dummyBlockHeight}
	nilSenderStatusMsg := &types.Status{ChainID: localChainBytes, Sender: nil, Genesis: dummyGenHash, BestBlockHash: dummyBlockHash, BestHeight: dummyBlockHeight}
	diffStatusMsg := &types.Status{ChainID: theirChainBytes, Sender: &dummyAddr, Genesis: dummyGenHash, BestBlockHash: dummyBlockHash, BestHeight: diffBlockNo}
	olderStatusMsg := &types.Status{ChainID: oldChainBytes, Sender: &dummyAddr, Genesis: dummyGenHash, BestBlockHash: dummyBlockHash, BestHeight: 10000}
	newerStatusMsg := &types.Status{ChainID: newChainBytes, Sender: &dummyAddr, Genesis: dummyGenHash, BestBlockHash: dummyBlockHash, BestHeight: 600000}
	diffVersionStatusMsg := &types.Status{ChainID: newChainBytes, Sender: &dummyAddr, Genesis: dummyGenHash, BestBlockHash: dummyBlockHash, BestHeight: dummyBlockHeight}
	wrongBlkIDStatusMsg := &types.Status{ChainID: localChainBytes, Sender: &dummyAddr, Genesis: dummyGenHash, BestBlockHash: []byte{}, BestHeight: dummyBlockHeight}

	tests := []struct {
		name       string
		readReturn *types.Status
		readError  error
		writeError error
		want       *p2pcommon.HandshakeResult
		wantErr    bool
		wantGoAway bool
	}{
		{"TSuccess", dummyStatusMsg, nil, nil, succResult, false, false},
		{"TOldChain", olderStatusMsg, nil, nil, succResult, false, false},
		{"TNewChain", newerStatusMsg, nil, nil, succResult, false, false},
		{"TUnexpMsg", nil, nil, nil, nil, true, true},
		{"TRFail", dummyStatusMsg, fmt.Errorf("failed"), nil, nil, true, true},
		{"TRNoSender", nilSenderStatusMsg, nil, nil, nil, true, true},
		{"TWFail", dummyStatusMsg, nil, fmt.Errorf("failed"), nil, true, false},
		{"TDiffChain", diffStatusMsg, nil, nil, nil, true, true},
		{"TNilGenesis", nilGenesisStatusMsg, nil, nil, nil, true, true},
		{"TDiffGenesis", diffGenesisStatusMsg, nil, nil, nil, true, true},
		{"TDiffChainVersion", diffVersionStatusMsg, nil, nil, nil, true, true},
		{"TWrongBestHash", wrongBlkIDStatusMsg, nil, nil, nil, true, true},

		//{"TSuccess", dummyStatusMsg, nil, nil, dummyStatusMsg, false, false},
		//{"TUnexpMsg", nil, nil, nil, nil, true, true},
		//{"TRFail", dummyStatusMsg, fmt.Errorf("failed"), nil, nil, true, true},
		//{"TRNoSender", nilSenderStatusMsg, nil, nil, nil, true, true},
		//{"TWFail", dummyStatusMsg, nil, fmt.Errorf("failed"), nil, true, false},
		//{"TDiffChain", diffStatusMsg, nil, nil, nil, true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dummyReader := p2pmock.NewMockReadWriteCloser(ctrl)
			mockRW := p2pmock.NewMockMsgReadWriter(ctrl)
			mockVM := p2pmock.NewMockVersionedManager(ctrl)
			var containerMsg *p2pcommon.MessageValue
			if tt.readReturn != nil {
				containerMsg = p2pcommon.NewSimpleMsgVal(p2pcommon.StatusRequest, p2pcommon.NewMsgID())
				statusBytes, _ := p2putil.MarshalMessageBody(tt.readReturn)
				containerMsg.SetPayload(statusBytes)
			} else {
				containerMsg = p2pcommon.NewSimpleMsgVal(p2pcommon.AddressesRequest, p2pcommon.NewMsgID())
			}
			mockRW.EXPECT().ReadMsg().Return(containerMsg, tt.readError).AnyTimes()
			if tt.wantGoAway {
				mockRW.EXPECT().WriteMsg(&MsgMatcher{p2pcommon.GoAway}).Return(tt.writeError)
			}
			mockRW.EXPECT().WriteMsg(&MsgMatcher{p2pcommon.StatusRequest}).Return(tt.writeError).MaxTimes(1)
			mockVM.EXPECT().GetBestChainID().Return(myChainID).AnyTimes()
			mockVM.EXPECT().GetChainID(gomock.Any()).DoAndReturn(fc.getChainID).AnyTimes()

			h := NewV200VersionedHS(dummyMeta, mockActor, logger, mockVM, nil, samplePeerID, dummyReader, dummyGenHash)
			h.msgRW = mockRW
			got, err := h.DoForOutbound(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("PeerHandshaker.handshakeOutboundPeer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil {
				if !got.Meta.Equals(tt.want.Meta) {
					t.Errorf("PeerHandshaker.handshakeOutboundPeer() peerID = %v, want %v", got.Meta, tt.want.Meta)
				}
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerHandshaker.handshakeOutboundPeer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestV200VersionedHS_DoForInbound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// t.SkipNow()
	logger := log.NewLogger("test")
	mockActor := p2pmock.NewMockActorService(ctrl)
	mockCA := p2pmock.NewMockChainAccessor(ctrl)

	fc := newFC(*myChainID, 10000, 20000, dummyBlockHeight+100)
	localChainID := *fc.getChainID(dummyBlockHeight)
	localChainBytes, _ := localChainID.Bytes()
	oldChainID := fc.getChainID(10000)
	oldChainBytes, _ := oldChainID.Bytes()
	newChainID := fc.getChainID(600000)
	newChainBytes, _ := newChainID.Bytes()

	dummyMeta := p2pcommon.NewMetaWith1Addr(samplePeerID, "dummy.aergo.io", 7846, "v2.0.0")
	dummyMeta.Version = sampleVersion
	dummyAddr := dummyMeta.ToPeerAddress()
	dummyBlock := &types.Block{Hash: dummyBlockHash, Header: &types.BlockHeader{BlockNo: dummyBlockHeight}}
	//dummyBlkRsp := message.GetBestBlockRsp{Block: dummyBlock}
	mockActor.EXPECT().GetChainAccessor().Return(mockCA).AnyTimes()
	mockCA.EXPECT().GetBestBlock().Return(dummyBlock, nil).AnyTimes()

	dummyGenHash := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	diffGenHash := []byte{0xff, 0xfe, 0xfd, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	dummyStatusMsg := &types.Status{ChainID: localChainBytes, Sender: &dummyAddr, Genesis: dummyGenHash, BestBlockHash: dummyBlockHash, BestHeight: dummyBlockHeight}
	succResult := &p2pcommon.HandshakeResult{Meta: dummyMeta, BestBlockHash: dummyBlockID, BestBlockNo: dummyBlockHeight}
	diffGenesisStatusMsg := &types.Status{ChainID: localChainBytes, Sender: &dummyAddr, Genesis: diffGenHash, BestBlockHash: dummyBlockHash, BestHeight: dummyBlockHeight}
	nilGenesisStatusMsg := &types.Status{ChainID: localChainBytes, Sender: &dummyAddr, Genesis: nil, BestBlockHash: dummyBlockHash, BestHeight: dummyBlockHeight}
	nilSenderStatusMsg := &types.Status{ChainID: localChainBytes, Sender: nil, Genesis: dummyGenHash, BestBlockHash: dummyBlockHash, BestHeight: dummyBlockHeight}
	diffStatusMsg := &types.Status{ChainID: theirChainBytes, Sender: &dummyAddr, Genesis: dummyGenHash, BestBlockHash: dummyBlockHash, BestHeight: dummyBlockHeight}
	olderStatusMsg := &types.Status{ChainID: oldChainBytes, Sender: &dummyAddr, Genesis: dummyGenHash, BestBlockHash: dummyBlockHash, BestHeight: 10000}
	diffVersionStatusMsg := &types.Status{ChainID: newChainBytes, Sender: &dummyAddr, Genesis: dummyGenHash, BestBlockHash: dummyBlockHash, BestHeight: dummyBlockHeight}
	wrongBlkIDStatusMsg := &types.Status{ChainID: localChainBytes, Sender: &dummyAddr, Genesis: dummyGenHash, BestBlockHash: []byte{}, BestHeight: dummyBlockHeight}

	tests := []struct {
		name       string
		readReturn *types.Status
		readError  error
		writeError error
		want       *p2pcommon.HandshakeResult
		wantErr    bool
		wantGoAway bool
	}{
		{"TSuccess", dummyStatusMsg, nil, nil, succResult, false, false},
		{"TOldChain", olderStatusMsg, nil, nil, succResult, false, false},
		{"TUnexpMsg", nil, nil, nil, nil, true, true},
		{"TRFail", dummyStatusMsg, fmt.Errorf("failed"), nil, nil, true, true},
		{"TRNoSender", nilSenderStatusMsg, nil, nil, nil, true, true},
		{"TWFail", dummyStatusMsg, nil, fmt.Errorf("failed"), nil, true, false},
		{"TDiffChain", diffStatusMsg, nil, nil, nil, true, true},
		{"TNilGenesis", nilGenesisStatusMsg, nil, nil, nil, true, true},
		{"TDiffGenesis", diffGenesisStatusMsg, nil, nil, nil, true, true},
		{"TDiffChainVersion", diffVersionStatusMsg, nil, nil, nil, true, true},
		{"TWrongBestHash", wrongBlkIDStatusMsg, nil, nil, nil, true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dummyReader := p2pmock.NewMockReadWriteCloser(ctrl)
			mockRW := p2pmock.NewMockMsgReadWriter(ctrl)
			mockVM := p2pmock.NewMockVersionedManager(ctrl)

			containerMsg := &p2pcommon.MessageValue{}
			if tt.readReturn != nil {
				containerMsg = p2pcommon.NewSimpleMsgVal(p2pcommon.StatusRequest, p2pcommon.NewMsgID())
				statusBytes, _ := p2putil.MarshalMessageBody(tt.readReturn)
				containerMsg.SetPayload(statusBytes)
			} else {
				containerMsg = p2pcommon.NewSimpleMsgVal(p2pcommon.AddressesRequest, p2pcommon.NewMsgID())
			}

			mockRW.EXPECT().ReadMsg().Return(containerMsg, tt.readError).AnyTimes()
			if tt.wantGoAway {
				mockRW.EXPECT().WriteMsg(&MsgMatcher{p2pcommon.GoAway}).Return(tt.writeError)
			}
			mockRW.EXPECT().WriteMsg(gomock.Any()).Return(tt.writeError).AnyTimes()
			mockVM.EXPECT().GetBestChainID().Return(myChainID).AnyTimes()
			mockVM.EXPECT().GetChainID(gomock.Any()).DoAndReturn(fc.getChainID).AnyTimes()

			h := NewV200VersionedHS(dummyMeta, mockActor, logger, mockVM, nil, samplePeerID, dummyReader, dummyGenHash)
			h.msgRW = mockRW
			got, err := h.DoForInbound(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("PeerHandshaker.DoForInbound() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil {
				if !got.Meta.Equals(tt.want.Meta) {
					t.Errorf("PeerHandshaker.handshakeOutboundPeer() peerID = %v, want %v", got.Meta, tt.want.Meta)
				}
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerHandshaker.DoForInbound() = %v, want %v", got, tt.want)
			}
		})
	}
}

type MsgMatcher struct {
	sub p2pcommon.SubProtocol
}

func (m MsgMatcher) Matches(x interface{}) bool {
	return x.(p2pcommon.Message).Subprotocol() == m.sub
}

func (m MsgMatcher) String() string {
	return "matcher " + m.sub.String()
}

func Test_createMessage(t *testing.T) {
	type args struct {
		protocolID p2pcommon.SubProtocol
		msgBody    p2pcommon.MessageBody
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
	}{
		{"TStatus", args{protocolID: p2pcommon.StatusRequest, msgBody: &types.Status{Version: "11"}}, false},
		{"TGOAway", args{protocolID: p2pcommon.GoAway, msgBody: &types.GoAwayNotice{Message: "test"}}, false},
		{"TNil", args{protocolID: p2pcommon.StatusRequest, msgBody: nil}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createMessage(tt.args.protocolID, p2pcommon.NewMsgID(), tt.args.msgBody)
			if (got == nil) != tt.wantNil {
				t.Errorf("createMessage() = %v, want nil %v", got, tt.wantNil)
			}
			if got != nil && got.Subprotocol() != tt.args.protocolID {
				t.Errorf("status.ProtocolID = %v, want %v", got.Subprotocol(), tt.args.protocolID)
			}
		})
	}
}

func TestV200Handshaker_createLocalStatus(t *testing.T) {
	logger := log.NewLogger("handshake.test")
	dummyMeta := p2pcommon.NewMetaWith1Addr(samplePeerID, "dummy.aergo.io", 7846, "v2.0.0")
	dummyMeta.Version = sampleVersion

	sampleSize := 5
	certs := make([]*p2pcommon.AgentCertificateV1,sampleSize)
	pids := make([]types.PeerID,sampleSize)

	for i := 0 ; i<5 ; i++ {
		priv, _, _ := crypto.GenerateKeyPair(crypto.Secp256k1, 256)
		id, _ := types.IDFromPrivateKey(priv)
		pids[i] = id
 		certs[i], _ = p2putil.NewAgentCertV1(id, samplePeerID,p2putil.ConvertPKToBTCEC(priv), []string{"192.168.1.3"}, time.Hour*24 )
	}
	wrongCert := *certs[0]
	wrongCert.AgentAddress = []string{}
	type args struct {
		role types.PeerRole
		pids []types.PeerID
		cert []*p2pcommon.AgentCertificateV1
	}
	tests := []struct {
		name string

		args        args

		wantProdIDs []types.PeerID
		wantCert    []*p2pcommon.AgentCertificateV1
		wantErr     bool
	}{
		{"TBP", args{types.PeerRole_Producer, nil, nil}, nil, nil, false},
		{"TWatcher", args{types.PeerRole_Watcher, nil, nil}, nil, nil, false},
		{"TAgent", args{types.PeerRole_Agent, pids, certs}, pids, nil, false},
		{"TAgentLessCert", args{types.PeerRole_Agent, pids, certs[1:3]},  pids, certs[1:3], false},
		{"TWrongCert", args{types.PeerRole_Agent, pids, []*p2pcommon.AgentCertificateV1{&wrongCert}},  pids, certs[1:3], true},

		//{"TAgentUnknownCert", args{types.PeerRole_Agent, pids[:2], certs}, nil, nil, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dummyReader := p2pmock.NewMockReadWriteCloser(ctrl)
			mockActor := p2pmock.NewMockActorService(ctrl)
			mockVM := p2pmock.NewMockVersionedManager(ctrl)
			mockCM := p2pmock.NewMockCertificateManager(ctrl)

			inMeta := dummyMeta
			inMeta.Role = tt.args.role
			inMeta.ProducerIDs = tt.args.pids
			sampleChainID := &types.ChainID{}
			sampleBlock := &types.Block{Hash: dummyBlockHash, Header: &types.BlockHeader{}}
			mockCM.EXPECT().GetCertificates().Return(tt.args.cert).MaxTimes(1)

			h := NewV200VersionedHS(inMeta, mockActor, logger, mockVM, mockCM, samplePeerID, dummyReader, dummyGenHash)

			got, err := h.createLocalStatus(sampleChainID, sampleBlock)
			if (err != nil) != tt.wantErr {
				t.Errorf("createLocalStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				sender := got.Sender
				if sender.Role != tt.args.role {
					t.Errorf("createLocalStatus() role = %v, want %v", sender.Role, tt.args.role)
				}
				if len(sender.ProducerIDs) != len(tt.wantProdIDs) {
					t.Errorf("createLocalStatus() producers = %v, want %v", sender.ProducerIDs, tt.wantProdIDs)
				} else {
					for i, pid := range tt.wantProdIDs {
						gpid := types.PeerID(sender.ProducerIDs[i])
						if !types.IsSamePeerID(gpid, pid) {
							t.Errorf("createLocalStatus() producers = %v, wantErr %v", sender.ProducerIDs, tt.wantProdIDs)
							return
						}
					}
				}
			}
		})
	}
}

func TestV200Handshaker_checkAgent(t *testing.T) {
	type fields struct {
		cm               p2pcommon.CertificateManager
		vm               p2pcommon.VersionedManager
		selfMeta         p2pcommon.PeerMeta
		actor            p2pcommon.ActorService
		logger           *log.Logger
		peerID           types.PeerID
		msgRW            p2pcommon.MsgReadWriter
		localGenesisHash []byte
		remoteMeta       p2pcommon.PeerMeta
		remoteCerts      []*p2pcommon.AgentCertificateV1
		remoteHash       types.BlockID
		remoteNo         types.BlockNo
	}
	type args struct {
		status *types.Status
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &V200Handshaker{
				cm:               tt.fields.cm,
				vm:               tt.fields.vm,
				selfMeta:         tt.fields.selfMeta,
				actor:            tt.fields.actor,
				logger:           tt.fields.logger,
				peerID:           tt.fields.peerID,
				msgRW:            tt.fields.msgRW,
				localGenesisHash: tt.fields.localGenesisHash,
				remoteMeta:       tt.fields.remoteMeta,
				remoteCerts:      tt.fields.remoteCerts,
				remoteHash:       tt.fields.remoteHash,
				remoteNo:         tt.fields.remoteNo,
			}
			if err := h.checkAgent(tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("checkAgent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}