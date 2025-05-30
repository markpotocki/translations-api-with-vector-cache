// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: embed.proto

package embed

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

// The request message containing the text to generate an embedding for.
type EmbeddingRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Text          string                 `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"` // The text to embed.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EmbeddingRequest) Reset() {
	*x = EmbeddingRequest{}
	mi := &file_embed_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EmbeddingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmbeddingRequest) ProtoMessage() {}

func (x *EmbeddingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_embed_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmbeddingRequest.ProtoReflect.Descriptor instead.
func (*EmbeddingRequest) Descriptor() ([]byte, []int) {
	return file_embed_proto_rawDescGZIP(), []int{0}
}

func (x *EmbeddingRequest) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

// The response message containing the embedding as a list of floats.
type EmbeddingResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Embedding     []float32              `protobuf:"fixed32,1,rep,packed,name=embedding,proto3" json:"embedding,omitempty"` // The embedding vector.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EmbeddingResponse) Reset() {
	*x = EmbeddingResponse{}
	mi := &file_embed_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EmbeddingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmbeddingResponse) ProtoMessage() {}

func (x *EmbeddingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_embed_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmbeddingResponse.ProtoReflect.Descriptor instead.
func (*EmbeddingResponse) Descriptor() ([]byte, []int) {
	return file_embed_proto_rawDescGZIP(), []int{1}
}

func (x *EmbeddingResponse) GetEmbedding() []float32 {
	if x != nil {
		return x.Embedding
	}
	return nil
}

var File_embed_proto protoreflect.FileDescriptor

const file_embed_proto_rawDesc = "" +
	"\n" +
	"\vembed.proto\x12\x05embed\"&\n" +
	"\x10EmbeddingRequest\x12\x12\n" +
	"\x04text\x18\x01 \x01(\tR\x04text\"1\n" +
	"\x11EmbeddingResponse\x12\x1c\n" +
	"\tembedding\x18\x01 \x03(\x02R\tembedding2R\n" +
	"\bEmbedder\x12F\n" +
	"\x11GenerateEmbedding\x12\x17.embed.EmbeddingRequest\x1a\x18.embed.EmbeddingResponseB\x1cZ\x1aembeddingapi/service;embedb\x06proto3"

var (
	file_embed_proto_rawDescOnce sync.Once
	file_embed_proto_rawDescData []byte
)

func file_embed_proto_rawDescGZIP() []byte {
	file_embed_proto_rawDescOnce.Do(func() {
		file_embed_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_embed_proto_rawDesc), len(file_embed_proto_rawDesc)))
	})
	return file_embed_proto_rawDescData
}

var file_embed_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_embed_proto_goTypes = []any{
	(*EmbeddingRequest)(nil),  // 0: embed.EmbeddingRequest
	(*EmbeddingResponse)(nil), // 1: embed.EmbeddingResponse
}
var file_embed_proto_depIdxs = []int32{
	0, // 0: embed.Embedder.GenerateEmbedding:input_type -> embed.EmbeddingRequest
	1, // 1: embed.Embedder.GenerateEmbedding:output_type -> embed.EmbeddingResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_embed_proto_init() }
func file_embed_proto_init() {
	if File_embed_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_embed_proto_rawDesc), len(file_embed_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_embed_proto_goTypes,
		DependencyIndexes: file_embed_proto_depIdxs,
		MessageInfos:      file_embed_proto_msgTypes,
	}.Build()
	File_embed_proto = out.File
	file_embed_proto_goTypes = nil
	file_embed_proto_depIdxs = nil
}
