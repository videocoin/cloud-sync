// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: syncer/v1/syncer_service.proto

package v1

import (
	context "context"
	encoding_binary "encoding/binary"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/types"
	golang_proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type SyncRequest struct {
	Data                 []byte   `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Path                 string   `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	ContentType          string   `protobuf:"bytes,3,opt,name=content_type,json=contentType,proto3" json:"content_type,omitempty"`
	Duration             float64  `protobuf:"fixed64,4,opt,name=duration,proto3" json:"duration,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SyncRequest) Reset()         { *m = SyncRequest{} }
func (m *SyncRequest) String() string { return proto.CompactTextString(m) }
func (*SyncRequest) ProtoMessage()    {}
func (*SyncRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e555e157b8570a45, []int{0}
}
func (m *SyncRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SyncRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SyncRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SyncRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncRequest.Merge(m, src)
}
func (m *SyncRequest) XXX_Size() int {
	return m.Size()
}
func (m *SyncRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SyncRequest proto.InternalMessageInfo

func (m *SyncRequest) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *SyncRequest) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *SyncRequest) GetContentType() string {
	if m != nil {
		return m.ContentType
	}
	return ""
}

func (m *SyncRequest) GetDuration() float64 {
	if m != nil {
		return m.Duration
	}
	return 0
}

func (*SyncRequest) XXX_MessageName() string {
	return "cloud.api.syncer.v1.SyncRequest"
}
func init() {
	proto.RegisterType((*SyncRequest)(nil), "cloud.api.syncer.v1.SyncRequest")
	golang_proto.RegisterType((*SyncRequest)(nil), "cloud.api.syncer.v1.SyncRequest")
}

func init() { proto.RegisterFile("syncer/v1/syncer_service.proto", fileDescriptor_e555e157b8570a45) }
func init() {
	golang_proto.RegisterFile("syncer/v1/syncer_service.proto", fileDescriptor_e555e157b8570a45)
}

var fileDescriptor_e555e157b8570a45 = []byte{
	// 246 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x8f, 0xcf, 0x4a, 0xc3, 0x40,
	0x10, 0xc6, 0x99, 0x5a, 0x44, 0xb7, 0x15, 0x21, 0x5e, 0x42, 0x84, 0x25, 0x7a, 0xca, 0xc5, 0x84,
	0xe8, 0x1b, 0xf8, 0x08, 0x89, 0x27, 0x2f, 0x25, 0xd9, 0x8c, 0x69, 0xa0, 0xcd, 0xac, 0xc9, 0x6c,
	0x60, 0xdf, 0xce, 0x63, 0x8f, 0x3e, 0x82, 0xa4, 0x2f, 0x22, 0xdd, 0x15, 0x7b, 0xfb, 0x7d, 0x7f,
	0x06, 0xe6, 0x13, 0x72, 0xb4, 0xbd, 0xc2, 0x21, 0x9b, 0xf2, 0xcc, 0xd3, 0x66, 0xc4, 0x61, 0xea,
	0x14, 0xa6, 0x7a, 0x20, 0xa6, 0xe0, 0x4e, 0xed, 0xc8, 0x34, 0x69, 0xa5, 0xbb, 0xd4, 0xe7, 0xe9,
	0x94, 0x47, 0xf7, 0x2d, 0x51, 0xbb, 0xc3, 0xcc, 0x55, 0x6a, 0xf3, 0x91, 0xe1, 0x5e, 0xb3, 0xf5,
	0x17, 0xd1, 0x53, 0xdb, 0xf1, 0xd6, 0xd4, 0xa9, 0xa2, 0x7d, 0xd6, 0x52, 0x4b, 0xe7, 0xd6, 0x49,
	0x39, 0xe1, 0xc8, 0xd7, 0x1f, 0x59, 0xac, 0x4a, 0xdb, 0xab, 0x02, 0x3f, 0x0d, 0x8e, 0x1c, 0x04,
	0x62, 0xd9, 0x54, 0x5c, 0x85, 0x10, 0x43, 0xb2, 0x2e, 0x1c, 0x9f, 0x3c, 0x5d, 0xf1, 0x36, 0x5c,
	0xc4, 0x90, 0x5c, 0x17, 0x8e, 0x83, 0x07, 0xb1, 0x56, 0xd4, 0x33, 0xf6, 0xbc, 0x61, 0xab, 0x31,
	0xbc, 0x70, 0xd9, 0xea, 0xcf, 0x7b, 0xb3, 0x1a, 0x83, 0x48, 0x5c, 0x35, 0x66, 0xa8, 0xb8, 0xa3,
	0x3e, 0x5c, 0xc6, 0x90, 0x40, 0xf1, 0xaf, 0x9f, 0x6f, 0xc5, 0x4d, 0xe9, 0xe6, 0x94, 0x7e, 0xed,
	0x6b, 0x78, 0x98, 0x25, 0x7c, 0xcf, 0x12, 0x7e, 0x66, 0x09, 0x5f, 0x47, 0x09, 0x87, 0xa3, 0x84,
	0xf7, 0xc5, 0x94, 0xd7, 0x97, 0xee, 0xcf, 0x97, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb0, 0xbe,
	0x44, 0xf6, 0x2a, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SyncerServiceClient is the client API for SyncerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SyncerServiceClient interface {
}

type syncerServiceClient struct {
	cc *grpc.ClientConn
}

func NewSyncerServiceClient(cc *grpc.ClientConn) SyncerServiceClient {
	return &syncerServiceClient{cc}
}

// SyncerServiceServer is the server API for SyncerService service.
type SyncerServiceServer interface {
}

// UnimplementedSyncerServiceServer can be embedded to have forward compatible implementations.
type UnimplementedSyncerServiceServer struct {
}

func RegisterSyncerServiceServer(s *grpc.Server, srv SyncerServiceServer) {
	s.RegisterService(&_SyncerService_serviceDesc, srv)
}

var _SyncerService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "cloud.api.syncer.v1.SyncerService",
	HandlerType: (*SyncerServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "syncer/v1/syncer_service.proto",
}

func (m *SyncRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SyncRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SyncRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Duration != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Duration))))
		i--
		dAtA[i] = 0x21
	}
	if len(m.ContentType) > 0 {
		i -= len(m.ContentType)
		copy(dAtA[i:], m.ContentType)
		i = encodeVarintSyncerService(dAtA, i, uint64(len(m.ContentType)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Path) > 0 {
		i -= len(m.Path)
		copy(dAtA[i:], m.Path)
		i = encodeVarintSyncerService(dAtA, i, uint64(len(m.Path)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Data) > 0 {
		i -= len(m.Data)
		copy(dAtA[i:], m.Data)
		i = encodeVarintSyncerService(dAtA, i, uint64(len(m.Data)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintSyncerService(dAtA []byte, offset int, v uint64) int {
	offset -= sovSyncerService(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *SyncRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovSyncerService(uint64(l))
	}
	l = len(m.Path)
	if l > 0 {
		n += 1 + l + sovSyncerService(uint64(l))
	}
	l = len(m.ContentType)
	if l > 0 {
		n += 1 + l + sovSyncerService(uint64(l))
	}
	if m.Duration != 0 {
		n += 9
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovSyncerService(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSyncerService(x uint64) (n int) {
	return sovSyncerService(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SyncRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSyncerService
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
			return fmt.Errorf("proto: SyncRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SyncRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSyncerService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSyncerService
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthSyncerService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append(m.Data[:0], dAtA[iNdEx:postIndex]...)
			if m.Data == nil {
				m.Data = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Path", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSyncerService
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
				return ErrInvalidLengthSyncerService
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSyncerService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Path = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ContentType", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSyncerService
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
				return ErrInvalidLengthSyncerService
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSyncerService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ContentType = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Duration", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Duration = float64(math.Float64frombits(v))
		default:
			iNdEx = preIndex
			skippy, err := skipSyncerService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthSyncerService
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthSyncerService
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipSyncerService(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSyncerService
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
					return 0, ErrIntOverflowSyncerService
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
					return 0, ErrIntOverflowSyncerService
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
				return 0, ErrInvalidLengthSyncerService
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSyncerService
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSyncerService
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSyncerService        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSyncerService          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSyncerService = fmt.Errorf("proto: unexpected end of group")
)
