// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kava/evmutil/v1beta1/query.proto

package types

import (
	context "context"
	fmt "fmt"
	query "github.com/cosmos/cosmos-sdk/types/query"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// QueryParamsRequest defines the request type for querying x/evmutil parameters.
type QueryParamsRequest struct {
}

func (m *QueryParamsRequest) Reset()         { *m = QueryParamsRequest{} }
func (m *QueryParamsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryParamsRequest) ProtoMessage()    {}
func (*QueryParamsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a8d0512331709e7, []int{0}
}
func (m *QueryParamsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryParamsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryParamsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryParamsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryParamsRequest.Merge(m, src)
}
func (m *QueryParamsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryParamsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryParamsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryParamsRequest proto.InternalMessageInfo

// QueryParamsResponse defines the response type for querying x/evmutil parameters.
type QueryParamsResponse struct {
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}

func (m *QueryParamsResponse) Reset()         { *m = QueryParamsResponse{} }
func (m *QueryParamsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryParamsResponse) ProtoMessage()    {}
func (*QueryParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a8d0512331709e7, []int{1}
}
func (m *QueryParamsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryParamsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryParamsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryParamsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryParamsResponse.Merge(m, src)
}
func (m *QueryParamsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryParamsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryParamsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryParamsResponse proto.InternalMessageInfo

func (m *QueryParamsResponse) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

// QueryDeployedCosmosCoinContractsRequest defines the request type for Query/DeployedCosmosCoinContracts method.
type QueryDeployedCosmosCoinContractsRequest struct {
	// optional query param to only return specific denoms in the list
	// denoms that do not have deployed contracts will be omitted from the result
	// must request fewer than 100 denoms at a time.
	CosmosDenoms []string `protobuf:"bytes,1,rep,name=cosmos_denoms,json=cosmosDenoms,proto3" json:"cosmos_denoms,omitempty"`
	// pagination defines an optional pagination for the request.
	Pagination *query.PageRequest `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryDeployedCosmosCoinContractsRequest) Reset() {
	*m = QueryDeployedCosmosCoinContractsRequest{}
}
func (m *QueryDeployedCosmosCoinContractsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryDeployedCosmosCoinContractsRequest) ProtoMessage()    {}
func (*QueryDeployedCosmosCoinContractsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a8d0512331709e7, []int{2}
}
func (m *QueryDeployedCosmosCoinContractsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryDeployedCosmosCoinContractsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryDeployedCosmosCoinContractsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryDeployedCosmosCoinContractsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryDeployedCosmosCoinContractsRequest.Merge(m, src)
}
func (m *QueryDeployedCosmosCoinContractsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryDeployedCosmosCoinContractsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryDeployedCosmosCoinContractsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryDeployedCosmosCoinContractsRequest proto.InternalMessageInfo

// QueryDeployedCosmosCoinContractsResponse defines the response type for the Query/DeployedCosmosCoinContracts method.
type QueryDeployedCosmosCoinContractsResponse struct {
	// deployed_cosmos_coin_contracts is a list of cosmos-sdk coin denom and its deployed contract address
	DeployedCosmosCoinContracts []DeployedCosmosCoinContract `protobuf:"bytes,1,rep,name=deployed_cosmos_coin_contracts,json=deployedCosmosCoinContracts,proto3" json:"deployed_cosmos_coin_contracts"`
	// pagination defines the pagination in the response.
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryDeployedCosmosCoinContractsResponse) Reset() {
	*m = QueryDeployedCosmosCoinContractsResponse{}
}
func (m *QueryDeployedCosmosCoinContractsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryDeployedCosmosCoinContractsResponse) ProtoMessage()    {}
func (*QueryDeployedCosmosCoinContractsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a8d0512331709e7, []int{3}
}
func (m *QueryDeployedCosmosCoinContractsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryDeployedCosmosCoinContractsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryDeployedCosmosCoinContractsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryDeployedCosmosCoinContractsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryDeployedCosmosCoinContractsResponse.Merge(m, src)
}
func (m *QueryDeployedCosmosCoinContractsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryDeployedCosmosCoinContractsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryDeployedCosmosCoinContractsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryDeployedCosmosCoinContractsResponse proto.InternalMessageInfo

func (m *QueryDeployedCosmosCoinContractsResponse) GetDeployedCosmosCoinContracts() []DeployedCosmosCoinContract {
	if m != nil {
		return m.DeployedCosmosCoinContracts
	}
	return nil
}

func (m *QueryDeployedCosmosCoinContractsResponse) GetPagination() *query.PageResponse {
	if m != nil {
		return m.Pagination
	}
	return nil
}

// DeployedCosmosCoinContract defines a deployed token contract to the evm representing a native cosmos-sdk coin
type DeployedCosmosCoinContract struct {
	CosmosDenom string              `protobuf:"bytes,1,opt,name=cosmos_denom,json=cosmosDenom,proto3" json:"cosmos_denom,omitempty"`
	Address     *InternalEVMAddress `protobuf:"bytes,2,opt,name=address,proto3,customtype=InternalEVMAddress" json:"address,omitempty"`
}

func (m *DeployedCosmosCoinContract) Reset()         { *m = DeployedCosmosCoinContract{} }
func (m *DeployedCosmosCoinContract) String() string { return proto.CompactTextString(m) }
func (*DeployedCosmosCoinContract) ProtoMessage()    {}
func (*DeployedCosmosCoinContract) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a8d0512331709e7, []int{4}
}
func (m *DeployedCosmosCoinContract) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DeployedCosmosCoinContract) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DeployedCosmosCoinContract.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DeployedCosmosCoinContract) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeployedCosmosCoinContract.Merge(m, src)
}
func (m *DeployedCosmosCoinContract) XXX_Size() int {
	return m.Size()
}
func (m *DeployedCosmosCoinContract) XXX_DiscardUnknown() {
	xxx_messageInfo_DeployedCosmosCoinContract.DiscardUnknown(m)
}

var xxx_messageInfo_DeployedCosmosCoinContract proto.InternalMessageInfo

func (m *DeployedCosmosCoinContract) GetCosmosDenom() string {
	if m != nil {
		return m.CosmosDenom
	}
	return ""
}

func init() {
	proto.RegisterType((*QueryParamsRequest)(nil), "kava.evmutil.v1beta1.QueryParamsRequest")
	proto.RegisterType((*QueryParamsResponse)(nil), "kava.evmutil.v1beta1.QueryParamsResponse")
	proto.RegisterType((*QueryDeployedCosmosCoinContractsRequest)(nil), "kava.evmutil.v1beta1.QueryDeployedCosmosCoinContractsRequest")
	proto.RegisterType((*QueryDeployedCosmosCoinContractsResponse)(nil), "kava.evmutil.v1beta1.QueryDeployedCosmosCoinContractsResponse")
	proto.RegisterType((*DeployedCosmosCoinContract)(nil), "kava.evmutil.v1beta1.DeployedCosmosCoinContract")
}

func init() { proto.RegisterFile("kava/evmutil/v1beta1/query.proto", fileDescriptor_4a8d0512331709e7) }

var fileDescriptor_4a8d0512331709e7 = []byte{
	// 542 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x54, 0x4f, 0x6f, 0xd3, 0x30,
	0x14, 0x8f, 0x0b, 0x14, 0xea, 0x8e, 0x8b, 0xa9, 0xd0, 0xd4, 0x55, 0xe9, 0x08, 0x88, 0x75, 0x48,
	0x38, 0x5b, 0x41, 0x1c, 0x26, 0x40, 0xa2, 0x1d, 0x20, 0x0e, 0x48, 0x2c, 0x07, 0x0e, 0x5c, 0x2a,
	0x27, 0xb1, 0x42, 0x44, 0x6a, 0xa7, 0xb1, 0x5b, 0x51, 0x71, 0x83, 0x0b, 0x47, 0x24, 0xbe, 0x40,
	0x3f, 0xce, 0x8e, 0x93, 0xb8, 0xa0, 0x1d, 0x26, 0xd4, 0x72, 0x40, 0x9c, 0xf8, 0x08, 0xa8, 0xb6,
	0xbb, 0x15, 0x91, 0xb6, 0x88, 0x9b, 0xf5, 0xfc, 0x7b, 0xfe, 0xfd, 0x79, 0x2f, 0x81, 0x9b, 0x6f,
	0xc8, 0x80, 0xb8, 0x74, 0xd0, 0xed, 0xcb, 0x38, 0x71, 0x07, 0xbb, 0x3e, 0x95, 0x64, 0xd7, 0xed,
	0xf5, 0x69, 0x36, 0xc4, 0x69, 0xc6, 0x25, 0x47, 0x95, 0x29, 0x02, 0x1b, 0x04, 0x36, 0x88, 0xea,
	0xad, 0x80, 0x8b, 0x2e, 0x17, 0xae, 0x4f, 0x04, 0xd5, 0xf0, 0xd3, 0xe6, 0x94, 0x44, 0x31, 0x23,
	0x32, 0xe6, 0x4c, 0xbf, 0x50, 0xad, 0x44, 0x3c, 0xe2, 0xea, 0xe8, 0x4e, 0x4f, 0xa6, 0x5a, 0x8b,
	0x38, 0x8f, 0x12, 0xea, 0x92, 0x34, 0x76, 0x09, 0x63, 0x5c, 0xaa, 0x16, 0x61, 0x6e, 0x9d, 0x5c,
	0x5d, 0x11, 0x65, 0x54, 0xc4, 0x06, 0xe3, 0x54, 0x20, 0x3a, 0x98, 0x32, 0xbf, 0x20, 0x19, 0xe9,
	0x0a, 0x8f, 0xf6, 0xfa, 0x54, 0x48, 0xe7, 0x00, 0x5e, 0xf9, 0xa3, 0x2a, 0x52, 0xce, 0x04, 0x45,
	0x7b, 0xb0, 0x98, 0xaa, 0xca, 0x3a, 0xd8, 0x04, 0x8d, 0x72, 0xb3, 0x86, 0xf3, 0x7c, 0x61, 0xdd,
	0xd5, 0x3a, 0x7f, 0x78, 0x52, 0xb7, 0x3c, 0xd3, 0xe1, 0x8c, 0x00, 0xdc, 0x52, 0x6f, 0xee, 0xd3,
	0x34, 0xe1, 0x43, 0x1a, 0xb6, 0x95, 0xf9, 0x36, 0x8f, 0x59, 0x9b, 0x33, 0x99, 0x91, 0x40, 0xce,
	0xe8, 0xd1, 0x75, 0x78, 0x59, 0x47, 0xd3, 0x09, 0x29, 0xe3, 0x8a, 0xee, 0x5c, 0xa3, 0xe4, 0xad,
	0xe9, 0xe2, 0xbe, 0xaa, 0xa1, 0x27, 0x10, 0x9e, 0xa5, 0xb4, 0x5e, 0x50, 0x82, 0x6e, 0x62, 0x0d,
	0xc1, 0xd3, 0x48, 0xb1, 0x9e, 0xc0, 0x99, 0xaa, 0x88, 0x1a, 0x02, 0x6f, 0xae, 0x73, 0xef, 0xd2,
	0xc7, 0x51, 0xdd, 0xfa, 0x31, 0xaa, 0x5b, 0xce, 0x2f, 0x00, 0x1b, 0xab, 0x25, 0x9a, 0x2c, 0xde,
	0x41, 0x3b, 0x34, 0xb0, 0x8e, 0x11, 0x1b, 0xf0, 0x98, 0x75, 0x82, 0x19, 0x52, 0x89, 0x2e, 0x37,
	0x77, 0xf2, 0x33, 0x5a, 0x4c, 0x61, 0x72, 0xdb, 0x08, 0x17, 0x8b, 0x40, 0x4f, 0x73, 0xbc, 0x6f,
	0xad, 0xf4, 0xae, 0x95, 0xcf, 0x9b, 0x77, 0x7a, 0xb0, 0xba, 0x58, 0x09, 0xba, 0x06, 0xd7, 0xe6,
	0xe7, 0xa0, 0xa6, 0x5e, 0xf2, 0xca, 0x73, 0x63, 0x40, 0x3b, 0xf0, 0x22, 0x09, 0xc3, 0x8c, 0x0a,
	0xa1, 0x64, 0x94, 0x5a, 0x57, 0x8f, 0x4f, 0xea, 0xe8, 0x19, 0x93, 0x34, 0x63, 0x24, 0x79, 0xfc,
	0xf2, 0xf9, 0x23, 0x7d, 0xeb, 0xcd, 0x60, 0xcd, 0x9f, 0x05, 0x78, 0x41, 0xa5, 0x8c, 0x3e, 0x00,
	0x58, 0xd4, 0xbb, 0x82, 0x1a, 0xf9, 0x29, 0xfd, 0xbd, 0x9a, 0xd5, 0xed, 0x7f, 0x40, 0x6a, 0xa3,
	0xce, 0x8d, 0xf7, 0x5f, 0xbe, 0x7f, 0x2e, 0xd8, 0xa8, 0xe6, 0xe6, 0x7e, 0x08, 0x7a, 0x31, 0xd1,
	0x31, 0x80, 0x1b, 0x4b, 0x06, 0x8e, 0x1e, 0x2c, 0x21, 0x5c, 0xbd, 0xcb, 0xd5, 0x87, 0xff, 0xdb,
	0x6e, 0x4c, 0xdc, 0x57, 0x26, 0xee, 0xa1, 0xbb, 0xf9, 0x26, 0x96, 0xef, 0x60, 0xab, 0x7d, 0x38,
	0xb6, 0xc1, 0xd1, 0xd8, 0x06, 0xdf, 0xc6, 0x36, 0xf8, 0x34, 0xb1, 0xad, 0xa3, 0x89, 0x6d, 0x7d,
	0x9d, 0xd8, 0xd6, 0xab, 0xed, 0x28, 0x96, 0xaf, 0xfb, 0x3e, 0x0e, 0x78, 0x57, 0xbd, 0x7c, 0x3b,
	0x21, 0xbe, 0xd0, 0x1c, 0x6f, 0x4f, 0x59, 0xe4, 0x30, 0xa5, 0xc2, 0x2f, 0xaa, 0x5f, 0xc5, 0x9d,
	0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x1f, 0xfa, 0x86, 0x41, 0xe8, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// Params queries all parameters of the evmutil module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	// DeployedCosmosCoinContracts queries a list cosmos coin denom and their deployed erc20 address
	DeployedCosmosCoinContracts(ctx context.Context, in *QueryDeployedCosmosCoinContractsRequest, opts ...grpc.CallOption) (*QueryDeployedCosmosCoinContractsResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, "/kava.evmutil.v1beta1.Query/Params", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) DeployedCosmosCoinContracts(ctx context.Context, in *QueryDeployedCosmosCoinContractsRequest, opts ...grpc.CallOption) (*QueryDeployedCosmosCoinContractsResponse, error) {
	out := new(QueryDeployedCosmosCoinContractsResponse)
	err := c.cc.Invoke(ctx, "/kava.evmutil.v1beta1.Query/DeployedCosmosCoinContracts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Params queries all parameters of the evmutil module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	// DeployedCosmosCoinContracts queries a list cosmos coin denom and their deployed erc20 address
	DeployedCosmosCoinContracts(context.Context, *QueryDeployedCosmosCoinContractsRequest) (*QueryDeployedCosmosCoinContractsResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Params(ctx context.Context, req *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (*UnimplementedQueryServer) DeployedCosmosCoinContracts(ctx context.Context, req *QueryDeployedCosmosCoinContractsRequest) (*QueryDeployedCosmosCoinContractsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeployedCosmosCoinContracts not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kava.evmutil.v1beta1.Query/Params",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_DeployedCosmosCoinContracts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryDeployedCosmosCoinContractsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).DeployedCosmosCoinContracts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kava.evmutil.v1beta1.Query/DeployedCosmosCoinContracts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).DeployedCosmosCoinContracts(ctx, req.(*QueryDeployedCosmosCoinContractsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "kava.evmutil.v1beta1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "DeployedCosmosCoinContracts",
			Handler:    _Query_DeployedCosmosCoinContracts_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kava/evmutil/v1beta1/query.proto",
}

func (m *QueryParamsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryParamsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryParamsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryParamsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryParamsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryParamsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *QueryDeployedCosmosCoinContractsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryDeployedCosmosCoinContractsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryDeployedCosmosCoinContractsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.CosmosDenoms) > 0 {
		for iNdEx := len(m.CosmosDenoms) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.CosmosDenoms[iNdEx])
			copy(dAtA[i:], m.CosmosDenoms[iNdEx])
			i = encodeVarintQuery(dAtA, i, uint64(len(m.CosmosDenoms[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *QueryDeployedCosmosCoinContractsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryDeployedCosmosCoinContractsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryDeployedCosmosCoinContractsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.DeployedCosmosCoinContracts) > 0 {
		for iNdEx := len(m.DeployedCosmosCoinContracts) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.DeployedCosmosCoinContracts[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQuery(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *DeployedCosmosCoinContract) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DeployedCosmosCoinContract) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DeployedCosmosCoinContract) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Address != nil {
		{
			size := m.Address.Size()
			i -= size
			if _, err := m.Address.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.CosmosDenom) > 0 {
		i -= len(m.CosmosDenom)
		copy(dAtA[i:], m.CosmosDenom)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.CosmosDenom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryParamsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryParamsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryDeployedCosmosCoinContractsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.CosmosDenoms) > 0 {
		for _, s := range m.CosmosDenoms {
			l = len(s)
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryDeployedCosmosCoinContractsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.DeployedCosmosCoinContracts) > 0 {
		for _, e := range m.DeployedCosmosCoinContracts {
			l = e.Size()
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *DeployedCosmosCoinContract) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.CosmosDenom)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	if m.Address != nil {
		l = m.Address.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryParamsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryParamsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryParamsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryParamsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryParamsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryParamsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryDeployedCosmosCoinContractsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryDeployedCosmosCoinContractsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryDeployedCosmosCoinContractsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CosmosDenoms", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CosmosDenoms = append(m.CosmosDenoms, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageRequest{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryDeployedCosmosCoinContractsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryDeployedCosmosCoinContractsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryDeployedCosmosCoinContractsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DeployedCosmosCoinContracts", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DeployedCosmosCoinContracts = append(m.DeployedCosmosCoinContracts, DeployedCosmosCoinContract{})
			if err := m.DeployedCosmosCoinContracts[len(m.DeployedCosmosCoinContracts)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageResponse{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *DeployedCosmosCoinContract) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: DeployedCosmosCoinContract: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DeployedCosmosCoinContract: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CosmosDenom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CosmosDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v InternalEVMAddress
			m.Address = &v
			if err := m.Address.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)