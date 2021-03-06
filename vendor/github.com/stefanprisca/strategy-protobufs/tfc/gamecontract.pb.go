// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gamecontract.proto

package tfc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type BuildType int32

const (
	BuildType_ROAD   BuildType = 0
	BuildType_SETTLE BuildType = 1
)

var BuildType_name = map[int32]string{
	0: "ROAD",
	1: "SETTLE",
}
var BuildType_value = map[string]int32{
	"ROAD":   0,
	"SETTLE": 1,
}

func (x BuildType) String() string {
	return proto.EnumName(BuildType_name, int32(x))
}
func (BuildType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_gamecontract_17bdb9e11a6eb477, []int{0}
}

type GameTrxType int32

const (
	GameTrxType_JOIN      GameTrxType = 0
	GameTrxType_NEXT      GameTrxType = 1
	GameTrxType_ROLL      GameTrxType = 2
	GameTrxType_TRADE     GameTrxType = 3
	GameTrxType_DEV       GameTrxType = 4
	GameTrxType_BATTLE    GameTrxType = 5
	GameTrxType_REGISTERL GameTrxType = 6
)

var GameTrxType_name = map[int32]string{
	0: "JOIN",
	1: "NEXT",
	2: "ROLL",
	3: "TRADE",
	4: "DEV",
	5: "BATTLE",
	6: "REGISTERL",
}
var GameTrxType_value = map[string]int32{
	"JOIN":      0,
	"NEXT":      1,
	"ROLL":      2,
	"TRADE":     3,
	"DEV":       4,
	"BATTLE":    5,
	"REGISTERL": 6,
}

func (x GameTrxType) String() string {
	return proto.EnumName(GameTrxType_name, int32(x))
}
func (GameTrxType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_gamecontract_17bdb9e11a6eb477, []int{1}
}

type JoinTrxPayload struct {
	Player               Player   `protobuf:"varint,1,opt,name=player,proto3,enum=tfc.Player" json:"player,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JoinTrxPayload) Reset()         { *m = JoinTrxPayload{} }
func (m *JoinTrxPayload) String() string { return proto.CompactTextString(m) }
func (*JoinTrxPayload) ProtoMessage()    {}
func (*JoinTrxPayload) Descriptor() ([]byte, []int) {
	return fileDescriptor_gamecontract_17bdb9e11a6eb477, []int{0}
}
func (m *JoinTrxPayload) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JoinTrxPayload.Unmarshal(m, b)
}
func (m *JoinTrxPayload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JoinTrxPayload.Marshal(b, m, deterministic)
}
func (dst *JoinTrxPayload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JoinTrxPayload.Merge(dst, src)
}
func (m *JoinTrxPayload) XXX_Size() int {
	return xxx_messageInfo_JoinTrxPayload.Size(m)
}
func (m *JoinTrxPayload) XXX_DiscardUnknown() {
	xxx_messageInfo_JoinTrxPayload.DiscardUnknown(m)
}

var xxx_messageInfo_JoinTrxPayload proto.InternalMessageInfo

func (m *JoinTrxPayload) GetPlayer() Player {
	if m != nil {
		return m.Player
	}
	return Player_RED
}

type TradeTrxPayload struct {
	Source               Player               `protobuf:"varint,1,opt,name=source,proto3,enum=tfc.Player" json:"source,omitempty"`
	Dest                 Player               `protobuf:"varint,2,opt,name=dest,proto3,enum=tfc.Player" json:"dest,omitempty"`
	Resource             Resource             `protobuf:"varint,3,opt,name=resource,proto3,enum=tfc.Resource" json:"resource,omitempty"`
	Amount               int32                `protobuf:"varint,4,opt,name=amount,proto3" json:"amount,omitempty"`
	LastUpdated          *timestamp.Timestamp `protobuf:"bytes,9,opt,name=last_updated,json=lastUpdated,proto3" json:"last_updated,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *TradeTrxPayload) Reset()         { *m = TradeTrxPayload{} }
func (m *TradeTrxPayload) String() string { return proto.CompactTextString(m) }
func (*TradeTrxPayload) ProtoMessage()    {}
func (*TradeTrxPayload) Descriptor() ([]byte, []int) {
	return fileDescriptor_gamecontract_17bdb9e11a6eb477, []int{1}
}
func (m *TradeTrxPayload) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TradeTrxPayload.Unmarshal(m, b)
}
func (m *TradeTrxPayload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TradeTrxPayload.Marshal(b, m, deterministic)
}
func (dst *TradeTrxPayload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TradeTrxPayload.Merge(dst, src)
}
func (m *TradeTrxPayload) XXX_Size() int {
	return xxx_messageInfo_TradeTrxPayload.Size(m)
}
func (m *TradeTrxPayload) XXX_DiscardUnknown() {
	xxx_messageInfo_TradeTrxPayload.DiscardUnknown(m)
}

var xxx_messageInfo_TradeTrxPayload proto.InternalMessageInfo

func (m *TradeTrxPayload) GetSource() Player {
	if m != nil {
		return m.Source
	}
	return Player_RED
}

func (m *TradeTrxPayload) GetDest() Player {
	if m != nil {
		return m.Dest
	}
	return Player_RED
}

func (m *TradeTrxPayload) GetResource() Resource {
	if m != nil {
		return m.Resource
	}
	return Resource_HILL
}

func (m *TradeTrxPayload) GetAmount() int32 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *TradeTrxPayload) GetLastUpdated() *timestamp.Timestamp {
	if m != nil {
		return m.LastUpdated
	}
	return nil
}

type BuildRoadPayload struct {
	EdgeID               uint32               `protobuf:"varint,1,opt,name=edgeID,proto3" json:"edgeID,omitempty"`
	Player               Player               `protobuf:"varint,2,opt,name=player,proto3,enum=tfc.Player" json:"player,omitempty"`
	LastUpdated          *timestamp.Timestamp `protobuf:"bytes,9,opt,name=last_updated,json=lastUpdated,proto3" json:"last_updated,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *BuildRoadPayload) Reset()         { *m = BuildRoadPayload{} }
func (m *BuildRoadPayload) String() string { return proto.CompactTextString(m) }
func (*BuildRoadPayload) ProtoMessage()    {}
func (*BuildRoadPayload) Descriptor() ([]byte, []int) {
	return fileDescriptor_gamecontract_17bdb9e11a6eb477, []int{2}
}
func (m *BuildRoadPayload) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BuildRoadPayload.Unmarshal(m, b)
}
func (m *BuildRoadPayload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BuildRoadPayload.Marshal(b, m, deterministic)
}
func (dst *BuildRoadPayload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BuildRoadPayload.Merge(dst, src)
}
func (m *BuildRoadPayload) XXX_Size() int {
	return xxx_messageInfo_BuildRoadPayload.Size(m)
}
func (m *BuildRoadPayload) XXX_DiscardUnknown() {
	xxx_messageInfo_BuildRoadPayload.DiscardUnknown(m)
}

var xxx_messageInfo_BuildRoadPayload proto.InternalMessageInfo

func (m *BuildRoadPayload) GetEdgeID() uint32 {
	if m != nil {
		return m.EdgeID
	}
	return 0
}

func (m *BuildRoadPayload) GetPlayer() Player {
	if m != nil {
		return m.Player
	}
	return Player_RED
}

func (m *BuildRoadPayload) GetLastUpdated() *timestamp.Timestamp {
	if m != nil {
		return m.LastUpdated
	}
	return nil
}

type BuildSettlePayload struct {
	SettleID             uint32               `protobuf:"varint,1,opt,name=settleID,proto3" json:"settleID,omitempty"`
	Player               Player               `protobuf:"varint,2,opt,name=player,proto3,enum=tfc.Player" json:"player,omitempty"`
	LastUpdated          *timestamp.Timestamp `protobuf:"bytes,9,opt,name=last_updated,json=lastUpdated,proto3" json:"last_updated,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *BuildSettlePayload) Reset()         { *m = BuildSettlePayload{} }
func (m *BuildSettlePayload) String() string { return proto.CompactTextString(m) }
func (*BuildSettlePayload) ProtoMessage()    {}
func (*BuildSettlePayload) Descriptor() ([]byte, []int) {
	return fileDescriptor_gamecontract_17bdb9e11a6eb477, []int{3}
}
func (m *BuildSettlePayload) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BuildSettlePayload.Unmarshal(m, b)
}
func (m *BuildSettlePayload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BuildSettlePayload.Marshal(b, m, deterministic)
}
func (dst *BuildSettlePayload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BuildSettlePayload.Merge(dst, src)
}
func (m *BuildSettlePayload) XXX_Size() int {
	return xxx_messageInfo_BuildSettlePayload.Size(m)
}
func (m *BuildSettlePayload) XXX_DiscardUnknown() {
	xxx_messageInfo_BuildSettlePayload.DiscardUnknown(m)
}

var xxx_messageInfo_BuildSettlePayload proto.InternalMessageInfo

func (m *BuildSettlePayload) GetSettleID() uint32 {
	if m != nil {
		return m.SettleID
	}
	return 0
}

func (m *BuildSettlePayload) GetPlayer() Player {
	if m != nil {
		return m.Player
	}
	return Player_RED
}

func (m *BuildSettlePayload) GetLastUpdated() *timestamp.Timestamp {
	if m != nil {
		return m.LastUpdated
	}
	return nil
}

type BuildTrxPayload struct {
	Type                 BuildType            `protobuf:"varint,1,opt,name=type,proto3,enum=tfc.BuildType" json:"type,omitempty"`
	BuildRoadPayload     *BuildRoadPayload    `protobuf:"bytes,2,opt,name=buildRoadPayload,proto3" json:"buildRoadPayload,omitempty"`
	BuildSettlePayload   *BuildSettlePayload  `protobuf:"bytes,3,opt,name=buildSettlePayload,proto3" json:"buildSettlePayload,omitempty"`
	LastUpdated          *timestamp.Timestamp `protobuf:"bytes,9,opt,name=last_updated,json=lastUpdated,proto3" json:"last_updated,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *BuildTrxPayload) Reset()         { *m = BuildTrxPayload{} }
func (m *BuildTrxPayload) String() string { return proto.CompactTextString(m) }
func (*BuildTrxPayload) ProtoMessage()    {}
func (*BuildTrxPayload) Descriptor() ([]byte, []int) {
	return fileDescriptor_gamecontract_17bdb9e11a6eb477, []int{4}
}
func (m *BuildTrxPayload) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BuildTrxPayload.Unmarshal(m, b)
}
func (m *BuildTrxPayload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BuildTrxPayload.Marshal(b, m, deterministic)
}
func (dst *BuildTrxPayload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BuildTrxPayload.Merge(dst, src)
}
func (m *BuildTrxPayload) XXX_Size() int {
	return xxx_messageInfo_BuildTrxPayload.Size(m)
}
func (m *BuildTrxPayload) XXX_DiscardUnknown() {
	xxx_messageInfo_BuildTrxPayload.DiscardUnknown(m)
}

var xxx_messageInfo_BuildTrxPayload proto.InternalMessageInfo

func (m *BuildTrxPayload) GetType() BuildType {
	if m != nil {
		return m.Type
	}
	return BuildType_ROAD
}

func (m *BuildTrxPayload) GetBuildRoadPayload() *BuildRoadPayload {
	if m != nil {
		return m.BuildRoadPayload
	}
	return nil
}

func (m *BuildTrxPayload) GetBuildSettlePayload() *BuildSettlePayload {
	if m != nil {
		return m.BuildSettlePayload
	}
	return nil
}

func (m *BuildTrxPayload) GetLastUpdated() *timestamp.Timestamp {
	if m != nil {
		return m.LastUpdated
	}
	return nil
}

type GameContractTrxArgs struct {
	Type                 GameTrxType          `protobuf:"varint,1,opt,name=type,proto3,enum=tfc.GameTrxType" json:"type,omitempty"`
	TradeTrxPayload      *TradeTrxPayload     `protobuf:"bytes,2,opt,name=tradeTrxPayload,proto3" json:"tradeTrxPayload,omitempty"`
	BuildTrxPayload      *BuildTrxPayload     `protobuf:"bytes,3,opt,name=buildTrxPayload,proto3" json:"buildTrxPayload,omitempty"`
	BattleTrxPayload     *BattleTrxPayload    `protobuf:"bytes,4,opt,name=battleTrxPayload,proto3" json:"battleTrxPayload,omitempty"`
	JoinTrxPayload       *JoinTrxPayload      `protobuf:"bytes,5,opt,name=joinTrxPayload,proto3" json:"joinTrxPayload,omitempty"`
	LastUpdated          *timestamp.Timestamp `protobuf:"bytes,9,opt,name=last_updated,json=lastUpdated,proto3" json:"last_updated,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *GameContractTrxArgs) Reset()         { *m = GameContractTrxArgs{} }
func (m *GameContractTrxArgs) String() string { return proto.CompactTextString(m) }
func (*GameContractTrxArgs) ProtoMessage()    {}
func (*GameContractTrxArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_gamecontract_17bdb9e11a6eb477, []int{5}
}
func (m *GameContractTrxArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GameContractTrxArgs.Unmarshal(m, b)
}
func (m *GameContractTrxArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GameContractTrxArgs.Marshal(b, m, deterministic)
}
func (dst *GameContractTrxArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GameContractTrxArgs.Merge(dst, src)
}
func (m *GameContractTrxArgs) XXX_Size() int {
	return xxx_messageInfo_GameContractTrxArgs.Size(m)
}
func (m *GameContractTrxArgs) XXX_DiscardUnknown() {
	xxx_messageInfo_GameContractTrxArgs.DiscardUnknown(m)
}

var xxx_messageInfo_GameContractTrxArgs proto.InternalMessageInfo

func (m *GameContractTrxArgs) GetType() GameTrxType {
	if m != nil {
		return m.Type
	}
	return GameTrxType_JOIN
}

func (m *GameContractTrxArgs) GetTradeTrxPayload() *TradeTrxPayload {
	if m != nil {
		return m.TradeTrxPayload
	}
	return nil
}

func (m *GameContractTrxArgs) GetBuildTrxPayload() *BuildTrxPayload {
	if m != nil {
		return m.BuildTrxPayload
	}
	return nil
}

func (m *GameContractTrxArgs) GetBattleTrxPayload() *BattleTrxPayload {
	if m != nil {
		return m.BattleTrxPayload
	}
	return nil
}

func (m *GameContractTrxArgs) GetJoinTrxPayload() *JoinTrxPayload {
	if m != nil {
		return m.JoinTrxPayload
	}
	return nil
}

func (m *GameContractTrxArgs) GetLastUpdated() *timestamp.Timestamp {
	if m != nil {
		return m.LastUpdated
	}
	return nil
}

type GameContractInitArgs struct {
	Uuid                 []byte               `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	LastUpdated          *timestamp.Timestamp `protobuf:"bytes,9,opt,name=last_updated,json=lastUpdated,proto3" json:"last_updated,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *GameContractInitArgs) Reset()         { *m = GameContractInitArgs{} }
func (m *GameContractInitArgs) String() string { return proto.CompactTextString(m) }
func (*GameContractInitArgs) ProtoMessage()    {}
func (*GameContractInitArgs) Descriptor() ([]byte, []int) {
	return fileDescriptor_gamecontract_17bdb9e11a6eb477, []int{6}
}
func (m *GameContractInitArgs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GameContractInitArgs.Unmarshal(m, b)
}
func (m *GameContractInitArgs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GameContractInitArgs.Marshal(b, m, deterministic)
}
func (dst *GameContractInitArgs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GameContractInitArgs.Merge(dst, src)
}
func (m *GameContractInitArgs) XXX_Size() int {
	return xxx_messageInfo_GameContractInitArgs.Size(m)
}
func (m *GameContractInitArgs) XXX_DiscardUnknown() {
	xxx_messageInfo_GameContractInitArgs.DiscardUnknown(m)
}

var xxx_messageInfo_GameContractInitArgs proto.InternalMessageInfo

func (m *GameContractInitArgs) GetUuid() []byte {
	if m != nil {
		return m.Uuid
	}
	return nil
}

func (m *GameContractInitArgs) GetLastUpdated() *timestamp.Timestamp {
	if m != nil {
		return m.LastUpdated
	}
	return nil
}

func init() {
	proto.RegisterType((*JoinTrxPayload)(nil), "tfc.JoinTrxPayload")
	proto.RegisterType((*TradeTrxPayload)(nil), "tfc.TradeTrxPayload")
	proto.RegisterType((*BuildRoadPayload)(nil), "tfc.BuildRoadPayload")
	proto.RegisterType((*BuildSettlePayload)(nil), "tfc.BuildSettlePayload")
	proto.RegisterType((*BuildTrxPayload)(nil), "tfc.BuildTrxPayload")
	proto.RegisterType((*GameContractTrxArgs)(nil), "tfc.GameContractTrxArgs")
	proto.RegisterType((*GameContractInitArgs)(nil), "tfc.GameContractInitArgs")
	proto.RegisterEnum("tfc.BuildType", BuildType_name, BuildType_value)
	proto.RegisterEnum("tfc.GameTrxType", GameTrxType_name, GameTrxType_value)
}

func init() { proto.RegisterFile("gamecontract.proto", fileDescriptor_gamecontract_17bdb9e11a6eb477) }

var fileDescriptor_gamecontract_17bdb9e11a6eb477 = []byte{
	// 583 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xbc, 0x94, 0x4f, 0x8f, 0xd2, 0x40,
	0x18, 0xc6, 0x2d, 0x14, 0x5c, 0xde, 0x2e, 0xd0, 0xcc, 0xa2, 0x12, 0x2e, 0xab, 0xe8, 0x41, 0xf7,
	0xd0, 0x4d, 0x30, 0x9e, 0x8c, 0x26, 0xac, 0x34, 0x84, 0x0d, 0xd9, 0xdd, 0x0c, 0xd5, 0x18, 0x2f,
	0x66, 0x4a, 0x67, 0x49, 0x4d, 0xa1, 0xa4, 0x4c, 0x13, 0xb9, 0xf9, 0x09, 0x3c, 0xf9, 0xe5, 0xfc,
	0x2a, 0x9e, 0x9c, 0xbe, 0x6d, 0xd9, 0x69, 0x59, 0x4f, 0x24, 0xde, 0x3a, 0xf3, 0x3c, 0xef, 0x3b,
	0x4f, 0x7f, 0xf3, 0x07, 0xc8, 0x82, 0x2d, 0xf9, 0x3c, 0x5c, 0x89, 0x88, 0xcd, 0x85, 0xb5, 0x8e,
	0x42, 0x11, 0x92, 0xaa, 0xb8, 0x9d, 0xf7, 0x5a, 0x89, 0xe0, 0x31, 0xc1, 0xd2, 0xc9, 0x5e, 0xc7,
	0x65, 0x42, 0x04, 0x25, 0x6b, 0xef, 0x74, 0x11, 0x86, 0x8b, 0x80, 0x9f, 0xe3, 0xc8, 0x8d, 0x6f,
	0xcf, 0x85, 0xbf, 0xe4, 0x1b, 0xc1, 0x96, 0xeb, 0xd4, 0xd0, 0x7f, 0x03, 0xad, 0xcb, 0xd0, 0x5f,
	0x39, 0xd1, 0xf7, 0x1b, 0xb6, 0x0d, 0x42, 0xe6, 0x91, 0xe7, 0x50, 0x5f, 0x07, 0x6c, 0xcb, 0xa3,
	0xae, 0xf6, 0x54, 0x7b, 0xd9, 0x1a, 0x18, 0x96, 0x5c, 0xce, 0xba, 0xc1, 0x29, 0x9a, 0x49, 0xfd,
	0xdf, 0x1a, 0xb4, 0x9d, 0x88, 0x79, 0xbc, 0x58, 0xb8, 0x09, 0xe3, 0x68, 0xce, 0xef, 0x2d, 0x4c,
	0x25, 0x72, 0x0a, 0xba, 0x27, 0x03, 0x74, 0x2b, 0xfb, 0x16, 0x14, 0xc8, 0x2b, 0x38, 0x8a, 0x78,
	0xd6, 0xa7, 0x8a, 0xa6, 0x26, 0x9a, 0x68, 0x36, 0x49, 0x77, 0x32, 0x79, 0x0c, 0x75, 0xb6, 0x0c,
	0xe3, 0x95, 0xe8, 0xea, 0xd2, 0x58, 0xa3, 0xd9, 0x88, 0xbc, 0x83, 0xe3, 0x80, 0x6d, 0xc4, 0xd7,
	0x78, 0x2d, 0xf9, 0x70, 0xaf, 0xdb, 0x90, 0xaa, 0x31, 0xe8, 0x59, 0x29, 0x0b, 0x2b, 0x67, 0x61,
	0x39, 0x39, 0x0b, 0x6a, 0x24, 0xfe, 0x8f, 0xa9, 0xbd, 0xff, 0x53, 0x03, 0xf3, 0x22, 0xf6, 0x03,
	0x8f, 0xca, 0xbf, 0xca, 0x7f, 0x4e, 0xae, 0xc5, 0xbd, 0x05, 0x9f, 0x8c, 0xf0, 0xe7, 0x9a, 0x34,
	0x1b, 0x29, 0xb4, 0x2a, 0xff, 0xa4, 0x75, 0x68, 0xa0, 0x5f, 0x1a, 0x10, 0x0c, 0x34, 0xe3, 0xc9,
	0x16, 0xe7, 0x91, 0x7a, 0x70, 0xb4, 0xc1, 0x89, 0x5d, 0xa8, 0xdd, 0xf8, 0xbf, 0xc4, 0xfa, 0x51,
	0x81, 0x36, 0xc6, 0x52, 0xce, 0x40, 0x1f, 0x74, 0xb1, 0x5d, 0xe7, 0x27, 0xa0, 0x85, 0xab, 0xa6,
	0x1e, 0x39, 0x4b, 0x51, 0x23, 0x43, 0x30, 0xdd, 0x12, 0x5e, 0x4c, 0x69, 0x0c, 0x1e, 0xdd, 0xf9,
	0x15, 0x91, 0xee, 0xd9, 0xc9, 0x18, 0x88, 0xbb, 0x07, 0x04, 0x8f, 0x8b, 0x31, 0x78, 0x72, 0xd7,
	0xa4, 0x20, 0xd3, 0x7b, 0x4a, 0x0e, 0x45, 0xf0, 0xa7, 0x02, 0x27, 0x63, 0x79, 0x0f, 0x3f, 0x64,
	0xb7, 0x4e, 0x92, 0x18, 0x46, 0x8b, 0x0d, 0x79, 0x51, 0xc0, 0x60, 0x62, 0xa2, 0xc4, 0x27, 0x75,
	0x05, 0xc4, 0x7b, 0x68, 0x8b, 0xe2, 0x1d, 0xca, 0x38, 0x74, 0xb0, 0xa0, 0x74, 0xbf, 0x68, 0xd9,
	0x9c, 0xd4, 0xbb, 0x45, 0xfe, 0x19, 0x82, 0x8e, 0xc2, 0x5d, 0xa9, 0x2f, 0x99, 0x71, 0x23, 0xf0,
	0xd1, 0x50, 0x1a, 0xe8, 0xea, 0x46, 0x94, 0x44, 0xba, 0x67, 0x27, 0x6f, 0xa1, 0xf5, 0xad, 0xf0,
	0x7c, 0x74, 0x6b, 0xd8, 0xe0, 0x04, 0x1b, 0x14, 0x5f, 0x16, 0x5a, 0xb2, 0x1e, 0x0a, 0xdf, 0x87,
	0x8e, 0xca, 0x7e, 0xb2, 0xf2, 0x05, 0xc2, 0x27, 0xa0, 0xc7, 0xb1, 0xef, 0x21, 0xfc, 0x63, 0x8a,
	0xdf, 0x07, 0x2e, 0x75, 0xf6, 0x0c, 0x1a, 0xbb, 0x53, 0x4c, 0x8e, 0x40, 0xa7, 0xd7, 0xc3, 0x91,
	0xf9, 0x80, 0x00, 0xd4, 0x67, 0xb6, 0xe3, 0x4c, 0x6d, 0x53, 0x3b, 0xfb, 0x02, 0x86, 0xb2, 0xc3,
	0x89, 0xe9, 0xf2, 0x7a, 0x72, 0x25, 0x4d, 0xf2, 0xeb, 0xca, 0xfe, 0xec, 0x98, 0x5a, 0x5a, 0x38,
	0x9d, 0x9a, 0x15, 0xd2, 0x80, 0x9a, 0x43, 0x87, 0x23, 0xdb, 0xac, 0x92, 0x87, 0x50, 0x1d, 0xd9,
	0x9f, 0x4c, 0x3d, 0x69, 0x76, 0x31, 0xc4, 0x66, 0x35, 0xd2, 0x84, 0x06, 0xb5, 0xc7, 0x93, 0x99,
	0x63, 0xd3, 0xa9, 0x59, 0x77, 0xeb, 0x98, 0xef, 0xf5, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x5e,
	0xbe, 0x25, 0x35, 0x0d, 0x06, 0x00, 0x00,
}
