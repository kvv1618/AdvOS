// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.15.8
// source: service.proto

package service

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Empty struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Empty) Reset() {
	*x = Empty{}
	mi := &file_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{0}
}

type JobDetailsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	FilePath      string                 `protobuf:"bytes,1,opt,name=filePath,proto3" json:"filePath,omitempty"`
	StartSeg      int32                  `protobuf:"varint,2,opt,name=startSeg,proto3" json:"startSeg,omitempty"`
	SegLen        int32                  `protobuf:"varint,3,opt,name=segLen,proto3" json:"segLen,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *JobDetailsResponse) Reset() {
	*x = JobDetailsResponse{}
	mi := &file_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *JobDetailsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobDetailsResponse) ProtoMessage() {}

func (x *JobDetailsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobDetailsResponse.ProtoReflect.Descriptor instead.
func (*JobDetailsResponse) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{1}
}

func (x *JobDetailsResponse) GetFilePath() string {
	if x != nil {
		return x.FilePath
	}
	return ""
}

func (x *JobDetailsResponse) GetStartSeg() int32 {
	if x != nil {
		return x.StartSeg
	}
	return 0
}

func (x *JobDetailsResponse) GetSegLen() int32 {
	if x != nil {
		return x.SegLen
	}
	return 0
}

type JobDataResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Data          []byte                 `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *JobDataResponse) Reset() {
	*x = JobDataResponse{}
	mi := &file_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *JobDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobDataResponse) ProtoMessage() {}

func (x *JobDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobDataResponse.ProtoReflect.Descriptor instead.
func (*JobDataResponse) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{2}
}

func (x *JobDataResponse) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type PartialResults struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	FilePath      string                 `protobuf:"bytes,1,opt,name=filePath,proto3" json:"filePath,omitempty"`
	StartSeg      int32                  `protobuf:"varint,2,opt,name=startSeg,proto3" json:"startSeg,omitempty"`
	SegLen        int32                  `protobuf:"varint,3,opt,name=segLen,proto3" json:"segLen,omitempty"`
	NumPrimes     int32                  `protobuf:"varint,4,opt,name=numPrimes,proto3" json:"numPrimes,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PartialResults) Reset() {
	*x = PartialResults{}
	mi := &file_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PartialResults) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PartialResults) ProtoMessage() {}

func (x *PartialResults) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PartialResults.ProtoReflect.Descriptor instead.
func (*PartialResults) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{3}
}

func (x *PartialResults) GetFilePath() string {
	if x != nil {
		return x.FilePath
	}
	return ""
}

func (x *PartialResults) GetStartSeg() int32 {
	if x != nil {
		return x.StartSeg
	}
	return 0
}

func (x *PartialResults) GetSegLen() int32 {
	if x != nil {
		return x.SegLen
	}
	return 0
}

func (x *PartialResults) GetNumPrimes() int32 {
	if x != nil {
		return x.NumPrimes
	}
	return 0
}

type StopWorkersRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StopWorkersRequest) Reset() {
	*x = StopWorkersRequest{}
	mi := &file_service_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StopWorkersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopWorkersRequest) ProtoMessage() {}

func (x *StopWorkersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopWorkersRequest.ProtoReflect.Descriptor instead.
func (*StopWorkersRequest) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{4}
}

func (x *StopWorkersRequest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_service_proto protoreflect.FileDescriptor

const file_service_proto_rawDesc = "" +
	"\n" +
	"\rservice.proto\x12\aservice\"\a\n" +
	"\x05empty\"d\n" +
	"\x12jobDetailsResponse\x12\x1a\n" +
	"\bfilePath\x18\x01 \x01(\tR\bfilePath\x12\x1a\n" +
	"\bstartSeg\x18\x02 \x01(\x05R\bstartSeg\x12\x16\n" +
	"\x06segLen\x18\x03 \x01(\x05R\x06segLen\"%\n" +
	"\x0fjobDataResponse\x12\x12\n" +
	"\x04data\x18\x01 \x01(\fR\x04data\"~\n" +
	"\x0epartialResults\x12\x1a\n" +
	"\bfilePath\x18\x01 \x01(\tR\bfilePath\x12\x1a\n" +
	"\bstartSeg\x18\x02 \x01(\x05R\bstartSeg\x12\x16\n" +
	"\x06segLen\x18\x03 \x01(\x05R\x06segLen\x12\x1c\n" +
	"\tnumPrimes\x18\x04 \x01(\x05R\tnumPrimes\".\n" +
	"\x12stopWorkersRequest\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage2G\n" +
	"\n" +
	"jobService\x129\n" +
	"\n" +
	"jobDetails\x12\x0e.service.empty\x1a\x1b.service.jobDetailsResponse2T\n" +
	"\x0ejobDataService\x12B\n" +
	"\ajobData\x12\x1b.service.jobDetailsResponse\x1a\x18.service.jobDataResponse0\x012T\n" +
	"\x16condenseResultsService\x12:\n" +
	"\x0fcondenseResults\x12\x17.service.partialResults\x1a\x0e.service.empty2P\n" +
	"\x12stopWorkersService\x12:\n" +
	"\vstopWorkers\x12\x1b.service.stopWorkersRequest\x1a\x0e.service.emptyB\x10Z\x0eprotoc/serviceb\x06proto3"

var (
	file_service_proto_rawDescOnce sync.Once
	file_service_proto_rawDescData []byte
)

func file_service_proto_rawDescGZIP() []byte {
	file_service_proto_rawDescOnce.Do(func() {
		file_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_service_proto_rawDesc), len(file_service_proto_rawDesc)))
	})
	return file_service_proto_rawDescData
}

var file_service_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_service_proto_goTypes = []any{
	(*Empty)(nil),              // 0: service.empty
	(*JobDetailsResponse)(nil), // 1: service.jobDetailsResponse
	(*JobDataResponse)(nil),    // 2: service.jobDataResponse
	(*PartialResults)(nil),     // 3: service.partialResults
	(*StopWorkersRequest)(nil), // 4: service.stopWorkersRequest
}
var file_service_proto_depIdxs = []int32{
	0, // 0: service.jobService.jobDetails:input_type -> service.empty
	1, // 1: service.jobDataService.jobData:input_type -> service.jobDetailsResponse
	3, // 2: service.condenseResultsService.condenseResults:input_type -> service.partialResults
	4, // 3: service.stopWorkersService.stopWorkers:input_type -> service.stopWorkersRequest
	1, // 4: service.jobService.jobDetails:output_type -> service.jobDetailsResponse
	2, // 5: service.jobDataService.jobData:output_type -> service.jobDataResponse
	0, // 6: service.condenseResultsService.condenseResults:output_type -> service.empty
	0, // 7: service.stopWorkersService.stopWorkers:output_type -> service.empty
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_service_proto_init() }
func file_service_proto_init() {
	if File_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_service_proto_rawDesc), len(file_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   4,
		},
		GoTypes:           file_service_proto_goTypes,
		DependencyIndexes: file_service_proto_depIdxs,
		MessageInfos:      file_service_proto_msgTypes,
	}.Build()
	File_service_proto = out.File
	file_service_proto_goTypes = nil
	file_service_proto_depIdxs = nil
}
