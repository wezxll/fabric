// Code generated by protoc-gen-go.
// source: openchain.proto
// DO NOT EDIT!

package protos

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "google/protobuf"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Transaction_Type int32

const (
	Transaction_UNDEFINED          Transaction_Type = 0
	Transaction_CHAINLET_NEW       Transaction_Type = 1
	Transaction_CHAINLET_UPDATE    Transaction_Type = 2
	Transaction_CHAINLET_EXECUTE   Transaction_Type = 3
	Transaction_CHAINLET_TERMINATE Transaction_Type = 4
)

var Transaction_Type_name = map[int32]string{
	0: "UNDEFINED",
	1: "CHAINLET_NEW",
	2: "CHAINLET_UPDATE",
	3: "CHAINLET_EXECUTE",
	4: "CHAINLET_TERMINATE",
}
var Transaction_Type_value = map[string]int32{
	"UNDEFINED":          0,
	"CHAINLET_NEW":       1,
	"CHAINLET_UPDATE":    2,
	"CHAINLET_EXECUTE":   3,
	"CHAINLET_TERMINATE": 4,
}

func (x Transaction_Type) String() string {
	return proto.EnumName(Transaction_Type_name, int32(x))
}

type OpenchainMessage_Type int32

const (
	OpenchainMessage_UNDEFINED              OpenchainMessage_Type = 0
	OpenchainMessage_DISC_HELLO             OpenchainMessage_Type = 1
	OpenchainMessage_DISC_DISCONNECT        OpenchainMessage_Type = 2
	OpenchainMessage_DISC_GET_PEERS         OpenchainMessage_Type = 3
	OpenchainMessage_DISC_PEERS             OpenchainMessage_Type = 4
	OpenchainMessage_DISC_NEWMSG            OpenchainMessage_Type = 5
	OpenchainMessage_CHAIN_STATUS           OpenchainMessage_Type = 7
	OpenchainMessage_CHAIN_GET_TRANSACTIONS OpenchainMessage_Type = 8
	OpenchainMessage_CHAIN_TRANSACTIONS     OpenchainMessage_Type = 9
	OpenchainMessage_CHAIN_GET_BLOCKS       OpenchainMessage_Type = 10
	OpenchainMessage_CHAIN_BLOCKS           OpenchainMessage_Type = 11
	OpenchainMessage_CHAIN_NEW_BLOCK        OpenchainMessage_Type = 12
	OpenchainMessage_REQUEST                OpenchainMessage_Type = 14
	OpenchainMessage_CONSENSUS              OpenchainMessage_Type = 15
)

var OpenchainMessage_Type_name = map[int32]string{
	0:  "UNDEFINED",
	1:  "DISC_HELLO",
	2:  "DISC_DISCONNECT",
	3:  "DISC_GET_PEERS",
	4:  "DISC_PEERS",
	5:  "DISC_NEWMSG",
	7:  "CHAIN_STATUS",
	8:  "CHAIN_GET_TRANSACTIONS",
	9:  "CHAIN_TRANSACTIONS",
	10: "CHAIN_GET_BLOCKS",
	11: "CHAIN_BLOCKS",
	12: "CHAIN_NEW_BLOCK",
	14: "REQUEST",
	15: "CONSENSUS",
}
var OpenchainMessage_Type_value = map[string]int32{
	"UNDEFINED":              0,
	"DISC_HELLO":             1,
	"DISC_DISCONNECT":        2,
	"DISC_GET_PEERS":         3,
	"DISC_PEERS":             4,
	"DISC_NEWMSG":            5,
	"CHAIN_STATUS":           7,
	"CHAIN_GET_TRANSACTIONS": 8,
	"CHAIN_TRANSACTIONS":     9,
	"CHAIN_GET_BLOCKS":       10,
	"CHAIN_BLOCKS":           11,
	"CHAIN_NEW_BLOCK":        12,
	"REQUEST":                14,
	"CONSENSUS":              15,
}

func (x OpenchainMessage_Type) String() string {
	return proto.EnumName(OpenchainMessage_Type_name, int32(x))
}

// Transaction defines a function call to a contract.
// `args` is an array of type string so that the chaincode writer can choose
// whatever format they wish for the arguments for their chaincode.
// For example, they may wish to use JSON, XML, or a custom format.
// TODO: Defined remaining fields.
type Transaction struct {
	Type       Transaction_Type `protobuf:"varint,1,opt,name=type,enum=protos.Transaction_Type" json:"type,omitempty"`
	ChainletID *ChainletID      `protobuf:"bytes,2,opt,name=chainletID" json:"chainletID,omitempty"`
	Function   string           `protobuf:"bytes,3,opt,name=function" json:"function,omitempty"`
	Args       []string         `protobuf:"bytes,4,rep,name=args" json:"args,omitempty"`
	Payload    []byte           `protobuf:"bytes,5,opt,name=payload,proto3" json:"payload,omitempty"`
	Uuid       string           `protobuf:"bytes,6,opt,name=uuid" json:"uuid,omitempty"`
}

func (m *Transaction) Reset()         { *m = Transaction{} }
func (m *Transaction) String() string { return proto.CompactTextString(m) }
func (*Transaction) ProtoMessage()    {}

func (m *Transaction) GetChainletID() *ChainletID {
	if m != nil {
		return m.ChainletID
	}
	return nil
}

// TransactionBlock carries a batch of transactions.
type TransactionBlock struct {
	Transactions []*Transaction `protobuf:"bytes,1,rep,name=transactions" json:"transactions,omitempty"`
}

func (m *TransactionBlock) Reset()         { *m = TransactionBlock{} }
func (m *TransactionBlock) String() string { return proto.CompactTextString(m) }
func (*TransactionBlock) ProtoMessage()    {}

func (m *TransactionBlock) GetTransactions() []*Transaction {
	if m != nil {
		return m.Transactions
	}
	return nil
}

// TODO: Explain when this message type is used.
// TODO: Explain fields.
// TODO: Rename field names according to style guide:
// https://developers.google.com/protocol-buffers/docs/style#message-and-field-names
type Block struct {
	ProposerID        string                     `protobuf:"bytes,1,opt,name=proposerID" json:"proposerID,omitempty"`
	Timestamp         *google_protobuf.Timestamp `protobuf:"bytes,2,opt,name=Timestamp" json:"Timestamp,omitempty"`
	Transactions      []*Transaction             `protobuf:"bytes,3,rep,name=transactions" json:"transactions,omitempty"`
	StateHash         []byte                     `protobuf:"bytes,4,opt,name=stateHash,proto3" json:"stateHash,omitempty"`
	PreviousBlockHash []byte                     `protobuf:"bytes,5,opt,name=previousBlockHash,proto3" json:"previousBlockHash,omitempty"`
}

func (m *Block) Reset()         { *m = Block{} }
func (m *Block) String() string { return proto.CompactTextString(m) }
func (*Block) ProtoMessage()    {}

func (m *Block) GetTimestamp() *google_protobuf.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *Block) GetTransactions() []*Transaction {
	if m != nil {
		return m.Transactions
	}
	return nil
}

type OpenchainMessage struct {
	Type      OpenchainMessage_Type      `protobuf:"varint,1,opt,name=type,enum=protos.OpenchainMessage_Type" json:"type,omitempty"`
	Timestamp *google_protobuf.Timestamp `protobuf:"bytes,2,opt,name=timestamp" json:"timestamp,omitempty"`
	Payload   []byte                     `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *OpenchainMessage) Reset()         { *m = OpenchainMessage{} }
func (m *OpenchainMessage) String() string { return proto.CompactTextString(m) }
func (*OpenchainMessage) ProtoMessage()    {}

func (m *OpenchainMessage) GetTimestamp() *google_protobuf.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func init() {
	proto.RegisterEnum("protos.Transaction_Type", Transaction_Type_name, Transaction_Type_value)
	proto.RegisterEnum("protos.OpenchainMessage_Type", OpenchainMessage_Type_name, OpenchainMessage_Type_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Client API for Peer service

type PeerClient interface {
	// Accepts a stream of OpenchainMessage during chat session, while receiving
	// other OpenchainMessage (e.g. from other peers).
	Chat(ctx context.Context, opts ...grpc.CallOption) (Peer_ChatClient, error)
}

type peerClient struct {
	cc *grpc.ClientConn
}

func NewPeerClient(cc *grpc.ClientConn) PeerClient {
	return &peerClient{cc}
}

func (c *peerClient) Chat(ctx context.Context, opts ...grpc.CallOption) (Peer_ChatClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Peer_serviceDesc.Streams[0], c.cc, "/protos.Peer/Chat", opts...)
	if err != nil {
		return nil, err
	}
	x := &peerChatClient{stream}
	return x, nil
}

type Peer_ChatClient interface {
	Send(*OpenchainMessage) error
	Recv() (*OpenchainMessage, error)
	grpc.ClientStream
}

type peerChatClient struct {
	grpc.ClientStream
}

func (x *peerChatClient) Send(m *OpenchainMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *peerChatClient) Recv() (*OpenchainMessage, error) {
	m := new(OpenchainMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Peer service

type PeerServer interface {
	// Accepts a stream of OpenchainMessage during chat session, while receiving
	// other OpenchainMessage (e.g. from other peers).
	Chat(Peer_ChatServer) error
}

func RegisterPeerServer(s *grpc.Server, srv PeerServer) {
	s.RegisterService(&_Peer_serviceDesc, srv)
}

func _Peer_Chat_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PeerServer).Chat(&peerChatServer{stream})
}

type Peer_ChatServer interface {
	Send(*OpenchainMessage) error
	Recv() (*OpenchainMessage, error)
	grpc.ServerStream
}

type peerChatServer struct {
	grpc.ServerStream
}

func (x *peerChatServer) Send(m *OpenchainMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *peerChatServer) Recv() (*OpenchainMessage, error) {
	m := new(OpenchainMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Peer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.Peer",
	HandlerType: (*PeerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Chat",
			Handler:       _Peer_Chat_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
}
