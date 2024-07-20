// api/proto/data_service.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: data_service.proto

package data_service

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type UpdateLatestDataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbols []string `protobuf:"bytes,1,rep,name=symbols,proto3" json:"symbols,omitempty"`
}

func (x *UpdateLatestDataRequest) Reset() {
	*x = UpdateLatestDataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateLatestDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateLatestDataRequest) ProtoMessage() {}

func (x *UpdateLatestDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_data_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateLatestDataRequest.ProtoReflect.Descriptor instead.
func (*UpdateLatestDataRequest) Descriptor() ([]byte, []int) {
	return file_data_service_proto_rawDescGZIP(), []int{0}
}

func (x *UpdateLatestDataRequest) GetSymbols() []string {
	if x != nil {
		return x.Symbols
	}
	return nil
}

type UpdateLatestDataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *UpdateLatestDataResponse) Reset() {
	*x = UpdateLatestDataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateLatestDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateLatestDataResponse) ProtoMessage() {}

func (x *UpdateLatestDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_data_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateLatestDataResponse.ProtoReflect.Descriptor instead.
func (*UpdateLatestDataResponse) Descriptor() ([]byte, []int) {
	return file_data_service_proto_rawDescGZIP(), []int{1}
}

func (x *UpdateLatestDataResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *UpdateLatestDataResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type StockRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbol    string `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	StartDate string `protobuf:"bytes,2,opt,name=start_date,json=startDate,proto3" json:"start_date,omitempty"`
	EndDate   string `protobuf:"bytes,3,opt,name=end_date,json=endDate,proto3" json:"end_date,omitempty"`
	Interval  string `protobuf:"bytes,4,opt,name=interval,proto3" json:"interval,omitempty"` // 1d, 1wk, 1mo
}

func (x *StockRequest) Reset() {
	*x = StockRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StockRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StockRequest) ProtoMessage() {}

func (x *StockRequest) ProtoReflect() protoreflect.Message {
	mi := &file_data_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StockRequest.ProtoReflect.Descriptor instead.
func (*StockRequest) Descriptor() ([]byte, []int) {
	return file_data_service_proto_rawDescGZIP(), []int{2}
}

func (x *StockRequest) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *StockRequest) GetStartDate() string {
	if x != nil {
		return x.StartDate
	}
	return ""
}

func (x *StockRequest) GetEndDate() string {
	if x != nil {
		return x.EndDate
	}
	return ""
}

func (x *StockRequest) GetInterval() string {
	if x != nil {
		return x.Interval
	}
	return ""
}

type StockResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbol     string            `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	DataPoints []*StockDataPoint `protobuf:"bytes,2,rep,name=data_points,json=dataPoints,proto3" json:"data_points,omitempty"`
}

func (x *StockResponse) Reset() {
	*x = StockResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StockResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StockResponse) ProtoMessage() {}

func (x *StockResponse) ProtoReflect() protoreflect.Message {
	mi := &file_data_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StockResponse.ProtoReflect.Descriptor instead.
func (*StockResponse) Descriptor() ([]byte, []int) {
	return file_data_service_proto_rawDescGZIP(), []int{3}
}

func (x *StockResponse) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *StockResponse) GetDataPoints() []*StockDataPoint {
	if x != nil {
		return x.DataPoints
	}
	return nil
}

type StockDataPoint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp     int64   `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Open          float64 `protobuf:"fixed64,2,opt,name=open,proto3" json:"open,omitempty"`
	High          float64 `protobuf:"fixed64,3,opt,name=high,proto3" json:"high,omitempty"`
	Low           float64 `protobuf:"fixed64,4,opt,name=low,proto3" json:"low,omitempty"`
	Close         float64 `protobuf:"fixed64,5,opt,name=close,proto3" json:"close,omitempty"`
	AdjustedClose float64 `protobuf:"fixed64,6,opt,name=adjusted_close,json=adjustedClose,proto3" json:"adjusted_close,omitempty"`
	Volume        int64   `protobuf:"varint,7,opt,name=volume,proto3" json:"volume,omitempty"`
}

func (x *StockDataPoint) Reset() {
	*x = StockDataPoint{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StockDataPoint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StockDataPoint) ProtoMessage() {}

func (x *StockDataPoint) ProtoReflect() protoreflect.Message {
	mi := &file_data_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StockDataPoint.ProtoReflect.Descriptor instead.
func (*StockDataPoint) Descriptor() ([]byte, []int) {
	return file_data_service_proto_rawDescGZIP(), []int{4}
}

func (x *StockDataPoint) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *StockDataPoint) GetOpen() float64 {
	if x != nil {
		return x.Open
	}
	return 0
}

func (x *StockDataPoint) GetHigh() float64 {
	if x != nil {
		return x.High
	}
	return 0
}

func (x *StockDataPoint) GetLow() float64 {
	if x != nil {
		return x.Low
	}
	return 0
}

func (x *StockDataPoint) GetClose() float64 {
	if x != nil {
		return x.Close
	}
	return 0
}

func (x *StockDataPoint) GetAdjustedClose() float64 {
	if x != nil {
		return x.AdjustedClose
	}
	return 0
}

func (x *StockDataPoint) GetVolume() int64 {
	if x != nil {
		return x.Volume
	}
	return 0
}

type BatchStockRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbols   []string `protobuf:"bytes,1,rep,name=symbols,proto3" json:"symbols,omitempty"`
	StartDate string   `protobuf:"bytes,2,opt,name=start_date,json=startDate,proto3" json:"start_date,omitempty"`
	EndDate   string   `protobuf:"bytes,3,opt,name=end_date,json=endDate,proto3" json:"end_date,omitempty"`
	Interval  string   `protobuf:"bytes,4,opt,name=interval,proto3" json:"interval,omitempty"`
}

func (x *BatchStockRequest) Reset() {
	*x = BatchStockRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BatchStockRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchStockRequest) ProtoMessage() {}

func (x *BatchStockRequest) ProtoReflect() protoreflect.Message {
	mi := &file_data_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchStockRequest.ProtoReflect.Descriptor instead.
func (*BatchStockRequest) Descriptor() ([]byte, []int) {
	return file_data_service_proto_rawDescGZIP(), []int{5}
}

func (x *BatchStockRequest) GetSymbols() []string {
	if x != nil {
		return x.Symbols
	}
	return nil
}

func (x *BatchStockRequest) GetStartDate() string {
	if x != nil {
		return x.StartDate
	}
	return ""
}

func (x *BatchStockRequest) GetEndDate() string {
	if x != nil {
		return x.EndDate
	}
	return ""
}

func (x *BatchStockRequest) GetInterval() string {
	if x != nil {
		return x.Interval
	}
	return ""
}

type BatchStockResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StockData map[string]*StockResponse `protobuf:"bytes,1,rep,name=stock_data,json=stockData,proto3" json:"stock_data,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Errors    map[string]string         `protobuf:"bytes,2,rep,name=errors,proto3" json:"errors,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *BatchStockResponse) Reset() {
	*x = BatchStockResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BatchStockResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchStockResponse) ProtoMessage() {}

func (x *BatchStockResponse) ProtoReflect() protoreflect.Message {
	mi := &file_data_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchStockResponse.ProtoReflect.Descriptor instead.
func (*BatchStockResponse) Descriptor() ([]byte, []int) {
	return file_data_service_proto_rawDescGZIP(), []int{6}
}

func (x *BatchStockResponse) GetStockData() map[string]*StockResponse {
	if x != nil {
		return x.StockData
	}
	return nil
}

func (x *BatchStockResponse) GetErrors() map[string]string {
	if x != nil {
		return x.Errors
	}
	return nil
}

var File_data_service_proto protoreflect.FileDescriptor

var file_data_service_proto_rawDesc = []byte{
	0x0a, 0x12, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x22, 0x33, 0x0a, 0x17, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4c, 0x61, 0x74, 0x65, 0x73,
	0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x73,
	0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x73, 0x22, 0x4e, 0x0a, 0x18, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x7c, 0x0a, 0x0c, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x12, 0x1d,
	0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x44, 0x61, 0x74, 0x65, 0x12, 0x19, 0x0a,
	0x08, 0x65, 0x6e, 0x64, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x65, 0x6e, 0x64, 0x44, 0x61, 0x74, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x76, 0x61, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x76, 0x61, 0x6c, 0x22, 0x65, 0x0a, 0x0d, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x12, 0x3c, 0x0a,
	0x0b, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x44, 0x61, 0x74, 0x61, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52,
	0x0a, 0x64, 0x61, 0x74, 0x61, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x22, 0xbd, 0x01, 0x0a, 0x0e,
	0x53, 0x74, 0x6f, 0x63, 0x6b, 0x44, 0x61, 0x74, 0x61, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x1c,
	0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x12, 0x0a, 0x04,
	0x6f, 0x70, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x6f, 0x70, 0x65, 0x6e,
	0x12, 0x12, 0x0a, 0x04, 0x68, 0x69, 0x67, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04,
	0x68, 0x69, 0x67, 0x68, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x6f, 0x77, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x03, 0x6c, 0x6f, 0x77, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x0e,
	0x61, 0x64, 0x6a, 0x75, 0x73, 0x74, 0x65, 0x64, 0x5f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x0d, 0x61, 0x64, 0x6a, 0x75, 0x73, 0x74, 0x65, 0x64, 0x43, 0x6c,
	0x6f, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x06, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x22, 0x83, 0x01, 0x0a, 0x11,
	0x42, 0x61, 0x74, 0x63, 0x68, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x07, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x44, 0x61, 0x74, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x6e,
	0x64, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x6e,
	0x64, 0x44, 0x61, 0x74, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61,
	0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61,
	0x6c, 0x22, 0xbd, 0x02, 0x0a, 0x12, 0x42, 0x61, 0x74, 0x63, 0x68, 0x53, 0x74, 0x6f, 0x63, 0x6b,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4d, 0x0a, 0x0a, 0x73, 0x74, 0x6f, 0x63,
	0x6b, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x64,
	0x61, 0x74, 0x61, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x42, 0x61, 0x74, 0x63, 0x68,
	0x53, 0x74, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x53, 0x74,
	0x6f, 0x63, 0x6b, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x09, 0x73, 0x74,
	0x6f, 0x63, 0x6b, 0x44, 0x61, 0x74, 0x61, 0x12, 0x43, 0x0a, 0x06, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x53, 0x74, 0x6f, 0x63, 0x6b,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x1a, 0x58, 0x0a, 0x0e,
	0x53, 0x74, 0x6f, 0x63, 0x6b, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x30, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x74,
	0x6f, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x39, 0x0a, 0x0b, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x32, 0x91, 0x02, 0x0a, 0x0b, 0x44, 0x61, 0x74, 0x61, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x47, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x44, 0x61, 0x74,
	0x61, 0x12, 0x19, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x53, 0x74, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x64,
	0x61, 0x74, 0x61, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x74, 0x6f, 0x63, 0x6b,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x56, 0x0a, 0x11, 0x47, 0x65,
	0x74, 0x42, 0x61, 0x74, 0x63, 0x68, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x1e, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x42, 0x61,
	0x74, 0x63, 0x68, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1f, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x42, 0x61,
	0x74, 0x63, 0x68, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x61, 0x0a, 0x10, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4c, 0x61, 0x74, 0x65,
	0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12, 0x24, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4c, 0x61, 0x74, 0x65, 0x73,
	0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x64,
	0x61, 0x74, 0x61, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x32, 0x5a, 0x30, 0x6d, 0x6f, 0x6d, 0x65, 0x6e, 0x74, 0x75,
	0x6d, 0x2d, 0x74, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x2d, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f,
	0x72, 0x6d, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x61, 0x74,
	0x61, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_data_service_proto_rawDescOnce sync.Once
	file_data_service_proto_rawDescData = file_data_service_proto_rawDesc
)

func file_data_service_proto_rawDescGZIP() []byte {
	file_data_service_proto_rawDescOnce.Do(func() {
		file_data_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_data_service_proto_rawDescData)
	})
	return file_data_service_proto_rawDescData
}

var file_data_service_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_data_service_proto_goTypes = []any{
	(*UpdateLatestDataRequest)(nil),  // 0: dataservice.UpdateLatestDataRequest
	(*UpdateLatestDataResponse)(nil), // 1: dataservice.UpdateLatestDataResponse
	(*StockRequest)(nil),             // 2: dataservice.StockRequest
	(*StockResponse)(nil),            // 3: dataservice.StockResponse
	(*StockDataPoint)(nil),           // 4: dataservice.StockDataPoint
	(*BatchStockRequest)(nil),        // 5: dataservice.BatchStockRequest
	(*BatchStockResponse)(nil),       // 6: dataservice.BatchStockResponse
	nil,                              // 7: dataservice.BatchStockResponse.StockDataEntry
	nil,                              // 8: dataservice.BatchStockResponse.ErrorsEntry
}
var file_data_service_proto_depIdxs = []int32{
	4, // 0: dataservice.StockResponse.data_points:type_name -> dataservice.StockDataPoint
	7, // 1: dataservice.BatchStockResponse.stock_data:type_name -> dataservice.BatchStockResponse.StockDataEntry
	8, // 2: dataservice.BatchStockResponse.errors:type_name -> dataservice.BatchStockResponse.ErrorsEntry
	3, // 3: dataservice.BatchStockResponse.StockDataEntry.value:type_name -> dataservice.StockResponse
	2, // 4: dataservice.DataService.GetStockData:input_type -> dataservice.StockRequest
	5, // 5: dataservice.DataService.GetBatchStockData:input_type -> dataservice.BatchStockRequest
	0, // 6: dataservice.DataService.UpdateLatestData:input_type -> dataservice.UpdateLatestDataRequest
	3, // 7: dataservice.DataService.GetStockData:output_type -> dataservice.StockResponse
	6, // 8: dataservice.DataService.GetBatchStockData:output_type -> dataservice.BatchStockResponse
	1, // 9: dataservice.DataService.UpdateLatestData:output_type -> dataservice.UpdateLatestDataResponse
	7, // [7:10] is the sub-list for method output_type
	4, // [4:7] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_data_service_proto_init() }
func file_data_service_proto_init() {
	if File_data_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_data_service_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateLatestDataRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_data_service_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateLatestDataResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_data_service_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*StockRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_data_service_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*StockResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_data_service_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*StockDataPoint); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_data_service_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*BatchStockRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_data_service_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*BatchStockResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_data_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_data_service_proto_goTypes,
		DependencyIndexes: file_data_service_proto_depIdxs,
		MessageInfos:      file_data_service_proto_msgTypes,
	}.Build()
	File_data_service_proto = out.File
	file_data_service_proto_rawDesc = nil
	file_data_service_proto_goTypes = nil
	file_data_service_proto_depIdxs = nil
}
