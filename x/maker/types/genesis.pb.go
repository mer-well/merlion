// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: merlion/maker/v1/genesis.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// GenesisState defines the maker module's genesis state.
type GenesisState struct {
	Params          Params                                 `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	CollateralRatio github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=collateral_ratio,json=collateralRatio,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"collateral_ratio" yaml:"collateral_ratio"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_60c1f3d62222051a, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

// Params defines the parameters for the maker module.
type Params struct {
	// adjusting collateral step
	CollateralRatioStep github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,1,opt,name=collateral_ratio_step,json=collateralRatioStep,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"collateral_ratio_step" yaml:"collateral_ratio_step"`
	// price band for adjusting collateral ratio
	CollateralRatioPriceBand github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=collateral_ratio_price_band,json=collateralRatioPriceBand,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"collateral_ratio_price_band" yaml:"collateral_ratio_price_band"`
	// cooldown period for adjusting collateral ratio
	CollateralRatioCooldownPeriod int64 `protobuf:"varint,3,opt,name=collateral_ratio_cooldown_period,json=collateralRatioCooldownPeriod,proto3" json:"collateral_ratio_cooldown_period,omitempty" yaml:"collateral_ration_cooldown_period"`
	// mint Mer price bias ratio
	MintPriceBias github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,4,opt,name=mint_price_bias,json=mintPriceBias,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"mint_price_bias" yaml:"mint_price_bias"`
	// burn Mer price bias ratio
	BurnPriceBias github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,5,opt,name=burn_price_bias,json=burnPriceBias,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"burn_price_bias" yaml:"burn_price_bias"`
	// recollateralization bonus ratio
	RecollateralizeBonus github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,6,opt,name=recollateralize_bonus,json=recollateralizeBonus,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"recollateralize_bonus" yaml:"recollateralize_bonus"`
	// liquidation commission fee ratio
	LiquidationCommissionFee github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,7,opt,name=liquidation_commission_fee,json=liquidationCommissionFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"liquidation_commission_fee" yaml:"liquidation_commission_fee"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_60c1f3d62222051a, []int{1}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetCollateralRatioCooldownPeriod() int64 {
	if m != nil {
		return m.CollateralRatioCooldownPeriod
	}
	return 0
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "merlion.maker.v1.GenesisState")
	proto.RegisterType((*Params)(nil), "merlion.maker.v1.Params")
}

func init() { proto.RegisterFile("merlion/maker/v1/genesis.proto", fileDescriptor_60c1f3d62222051a) }

var fileDescriptor_60c1f3d62222051a = []byte{
	// 532 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x94, 0xb1, 0x6f, 0xd3, 0x4e,
	0x14, 0xc7, 0x73, 0xbf, 0xf6, 0x97, 0xd2, 0x83, 0xaa, 0x95, 0x69, 0xc1, 0x2a, 0x60, 0x87, 0x1b,
	0x50, 0x06, 0x62, 0xab, 0x20, 0x31, 0x74, 0x74, 0x11, 0x85, 0x05, 0x45, 0x17, 0x26, 0x96, 0xe8,
	0x62, 0x1f, 0xe1, 0x54, 0xfb, 0xce, 0xdc, 0x5d, 0x0a, 0xed, 0x08, 0xff, 0x00, 0x88, 0x85, 0xb1,
	0xe2, 0xef, 0xe0, 0x0f, 0xe8, 0xd8, 0x11, 0x31, 0x44, 0x28, 0x59, 0x98, 0xfb, 0x17, 0x20, 0x9f,
	0x1d, 0x12, 0x5d, 0xd3, 0x21, 0xea, 0xe4, 0xf3, 0xfb, 0xbe, 0xf7, 0x7d, 0x9f, 0x77, 0xb6, 0x1e,
	0xf4, 0x32, 0x2a, 0x53, 0x26, 0x78, 0x98, 0x91, 0x03, 0x2a, 0xc3, 0xc3, 0x9d, 0xb0, 0x4f, 0x39,
	0x55, 0x4c, 0x05, 0xb9, 0x14, 0x5a, 0x38, 0x1b, 0x95, 0x1e, 0x18, 0x3d, 0x38, 0xdc, 0xd9, 0xde,
	0xec, 0x8b, 0xbe, 0x30, 0x62, 0x58, 0x9c, 0xca, 0x3c, 0xf4, 0x03, 0xc0, 0x1b, 0xfb, 0x65, 0x65,
	0x47, 0x13, 0x4d, 0x9d, 0x27, 0xb0, 0x9e, 0x13, 0x49, 0x32, 0xe5, 0x82, 0x06, 0x68, 0x5e, 0x7f,
	0xe4, 0x06, 0xb6, 0x53, 0xd0, 0x36, 0x7a, 0xb4, 0x7c, 0x3a, 0xf4, 0x6b, 0xb8, 0xca, 0x76, 0x34,
	0xdc, 0x88, 0x45, 0x9a, 0x12, 0x4d, 0x25, 0x49, 0xbb, 0x92, 0x68, 0x26, 0xdc, 0xff, 0x1a, 0xa0,
	0xb9, 0x1a, 0xbd, 0x28, 0xf2, 0x7e, 0x0d, 0xfd, 0x07, 0x7d, 0xa6, 0xdf, 0x0e, 0x7a, 0x41, 0x2c,
	0xb2, 0x30, 0x16, 0x2a, 0x13, 0xaa, 0x7a, 0xb4, 0x54, 0x72, 0x10, 0xea, 0xa3, 0x9c, 0xaa, 0xe0,
	0x29, 0x8d, 0xcf, 0x87, 0xfe, 0xed, 0x23, 0x92, 0xa5, 0xbb, 0xc8, 0xf6, 0x43, 0x78, 0x7d, 0x1a,
	0xc2, 0x26, 0xf2, 0x7d, 0x05, 0xd6, 0x4b, 0x1c, 0xe7, 0x23, 0x80, 0x5b, 0x76, 0x45, 0x57, 0x69,
	0x9a, 0x9b, 0x41, 0x56, 0xa3, 0x97, 0x0b, 0x63, 0xdc, 0x9d, 0x8f, 0x61, 0x4c, 0x11, 0xbe, 0x69,
	0xb1, 0x74, 0x34, 0xcd, 0x9d, 0xaf, 0x00, 0xde, 0xb9, 0x90, 0x9f, 0x4b, 0x16, 0xd3, 0x6e, 0x8f,
	0xf0, 0xa4, 0xba, 0x91, 0x57, 0x0b, 0xa3, 0xa0, 0x4b, 0x50, 0xa6, 0xd6, 0x08, 0xbb, 0x16, 0x50,
	0xbb, 0xd0, 0x22, 0xc2, 0x13, 0x67, 0x00, 0x1b, 0x17, 0x2a, 0x63, 0x21, 0xd2, 0x44, 0xbc, 0xe7,
	0xdd, 0x9c, 0x4a, 0x26, 0x12, 0x77, 0xa9, 0x01, 0x9a, 0x4b, 0xd1, 0xc3, 0xf3, 0xa1, 0xdf, 0x9c,
	0xdf, 0x8b, 0xdb, 0x25, 0x08, 0xdf, 0xb3, 0x3a, 0xee, 0x55, 0x09, 0x6d, 0xa3, 0x3b, 0x39, 0x5c,
	0xcf, 0x18, 0xd7, 0x13, 0x48, 0x46, 0x94, 0xbb, 0x6c, 0xe6, 0x7f, 0xbe, 0xf0, 0xfc, 0xb7, 0x4a,
	0x26, 0xcb, 0x0e, 0xe1, 0xb5, 0x22, 0x52, 0x0e, 0xca, 0x88, 0x2a, 0x3a, 0xf6, 0x06, 0x92, 0xcf,
	0x76, 0xfc, 0xff, 0x6a, 0x1d, 0x2d, 0x3b, 0x84, 0xd7, 0x8a, 0xc8, 0xb4, 0xe3, 0x27, 0x00, 0xb7,
	0x24, 0x9d, 0xde, 0x03, 0x3b, 0xa6, 0xdd, 0x9e, 0xe0, 0x03, 0xe5, 0xd6, 0xaf, 0xf6, 0xd7, 0xcd,
	0x35, 0x45, 0x78, 0xd3, 0x8a, 0x47, 0x45, 0xd8, 0xf9, 0x02, 0xe0, 0x76, 0xca, 0xde, 0x0d, 0x58,
	0x32, 0xf9, 0x52, 0x59, 0xc6, 0x94, 0x2a, 0x8e, 0x6f, 0x28, 0x75, 0x57, 0x0c, 0x4a, 0x67, 0x61,
	0x94, 0xfb, 0x25, 0xca, 0xe5, 0xce, 0x08, 0xbb, 0x33, 0xe2, 0xde, 0x3f, 0xed, 0x19, 0xa5, 0xbb,
	0xd7, 0xbe, 0x9d, 0xf8, 0xb5, 0x3f, 0x27, 0x3e, 0x88, 0xf6, 0x4f, 0x47, 0x1e, 0x38, 0x1b, 0x79,
	0xe0, 0xf7, 0xc8, 0x03, 0x9f, 0xc7, 0x5e, 0xed, 0x6c, 0xec, 0xd5, 0x7e, 0x8e, 0xbd, 0xda, 0xeb,
	0xd6, 0x0c, 0x4a, 0xb5, 0x66, 0x5a, 0xc7, 0x82, 0xd3, 0xc9, 0x4b, 0xf8, 0xa1, 0xda, 0x6f, 0x86,
	0xaa, 0x57, 0x37, 0x3b, 0xeb, 0xf1, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf5, 0x7b, 0xbf, 0x3c,
	0xfd, 0x04, 0x00, 0x00,
}

func (this *Params) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Params)
	if !ok {
		that2, ok := that.(Params)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.CollateralRatioStep.Equal(that1.CollateralRatioStep) {
		return false
	}
	if !this.CollateralRatioPriceBand.Equal(that1.CollateralRatioPriceBand) {
		return false
	}
	if this.CollateralRatioCooldownPeriod != that1.CollateralRatioCooldownPeriod {
		return false
	}
	if !this.MintPriceBias.Equal(that1.MintPriceBias) {
		return false
	}
	if !this.BurnPriceBias.Equal(that1.BurnPriceBias) {
		return false
	}
	if !this.RecollateralizeBonus.Equal(that1.RecollateralizeBonus) {
		return false
	}
	if !this.LiquidationCommissionFee.Equal(that1.LiquidationCommissionFee) {
		return false
	}
	return true
}
func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.CollateralRatio.Size()
		i -= size
		if _, err := m.CollateralRatio.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.LiquidationCommissionFee.Size()
		i -= size
		if _, err := m.LiquidationCommissionFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	{
		size := m.RecollateralizeBonus.Size()
		i -= size
		if _, err := m.RecollateralizeBonus.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	{
		size := m.BurnPriceBias.Size()
		i -= size
		if _, err := m.BurnPriceBias.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.MintPriceBias.Size()
		i -= size
		if _, err := m.MintPriceBias.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if m.CollateralRatioCooldownPeriod != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.CollateralRatioCooldownPeriod))
		i--
		dAtA[i] = 0x18
	}
	{
		size := m.CollateralRatioPriceBand.Size()
		i -= size
		if _, err := m.CollateralRatioPriceBand.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.CollateralRatioStep.Size()
		i -= size
		if _, err := m.CollateralRatioStep.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	l = m.CollateralRatio.Size()
	n += 1 + l + sovGenesis(uint64(l))
	return n
}

func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.CollateralRatioStep.Size()
	n += 1 + l + sovGenesis(uint64(l))
	l = m.CollateralRatioPriceBand.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if m.CollateralRatioCooldownPeriod != 0 {
		n += 1 + sovGenesis(uint64(m.CollateralRatioCooldownPeriod))
	}
	l = m.MintPriceBias.Size()
	n += 1 + l + sovGenesis(uint64(l))
	l = m.BurnPriceBias.Size()
	n += 1 + l + sovGenesis(uint64(l))
	l = m.RecollateralizeBonus.Size()
	n += 1 + l + sovGenesis(uint64(l))
	l = m.LiquidationCommissionFee.Size()
	n += 1 + l + sovGenesis(uint64(l))
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CollateralRatio", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.CollateralRatio.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CollateralRatioStep", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.CollateralRatioStep.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CollateralRatioPriceBand", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.CollateralRatioPriceBand.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CollateralRatioCooldownPeriod", wireType)
			}
			m.CollateralRatioCooldownPeriod = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CollateralRatioCooldownPeriod |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MintPriceBias", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MintPriceBias.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BurnPriceBias", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.BurnPriceBias.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RecollateralizeBonus", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.RecollateralizeBonus.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LiquidationCommissionFee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.LiquidationCommissionFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
