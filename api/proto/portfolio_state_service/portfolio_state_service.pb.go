// api/proto/portfolio_state_service.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: portfolio_state_service.proto

package portfolio_state_service

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

type PortfolioState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Date        string      `protobuf:"bytes,1,opt,name=date,proto3" json:"date,omitempty"`
	Positions   []*Position `protobuf:"bytes,2,rep,name=positions,proto3" json:"positions,omitempty"`
	CashBalance float64     `protobuf:"fixed64,3,opt,name=cash_balance,json=cashBalance,proto3" json:"cash_balance,omitempty"`
	TotalValue  float64     `protobuf:"fixed64,4,opt,name=total_value,json=totalValue,proto3" json:"total_value,omitempty"`
}

func (x *PortfolioState) Reset() {
	*x = PortfolioState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_portfolio_state_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PortfolioState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PortfolioState) ProtoMessage() {}

func (x *PortfolioState) ProtoReflect() protoreflect.Message {
	mi := &file_portfolio_state_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PortfolioState.ProtoReflect.Descriptor instead.
func (*PortfolioState) Descriptor() ([]byte, []int) {
	return file_portfolio_state_service_proto_rawDescGZIP(), []int{0}
}

func (x *PortfolioState) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

func (x *PortfolioState) GetPositions() []*Position {
	if x != nil {
		return x.Positions
	}
	return nil
}

func (x *PortfolioState) GetCashBalance() float64 {
	if x != nil {
		return x.CashBalance
	}
	return 0
}

func (x *PortfolioState) GetTotalValue() float64 {
	if x != nil {
		return x.TotalValue
	}
	return 0
}

type Position struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbol       string  `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Quantity     int32   `protobuf:"varint,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
	AveragePrice float64 `protobuf:"fixed64,3,opt,name=average_price,json=averagePrice,proto3" json:"average_price,omitempty"`
	CurrentPrice float64 `protobuf:"fixed64,4,opt,name=current_price,json=currentPrice,proto3" json:"current_price,omitempty"`
	MarketValue  float64 `protobuf:"fixed64,5,opt,name=market_value,json=marketValue,proto3" json:"market_value,omitempty"`
}

func (x *Position) Reset() {
	*x = Position{}
	if protoimpl.UnsafeEnabled {
		mi := &file_portfolio_state_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Position) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Position) ProtoMessage() {}

func (x *Position) ProtoReflect() protoreflect.Message {
	mi := &file_portfolio_state_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Position.ProtoReflect.Descriptor instead.
func (*Position) Descriptor() ([]byte, []int) {
	return file_portfolio_state_service_proto_rawDescGZIP(), []int{1}
}

func (x *Position) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *Position) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *Position) GetAveragePrice() float64 {
	if x != nil {
		return x.AveragePrice
	}
	return 0
}

func (x *Position) GetCurrentPrice() float64 {
	if x != nil {
		return x.CurrentPrice
	}
	return 0
}

func (x *Position) GetMarketValue() float64 {
	if x != nil {
		return x.MarketValue
	}
	return 0
}

type SaveResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *SaveResponse) Reset() {
	*x = SaveResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_portfolio_state_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveResponse) ProtoMessage() {}

func (x *SaveResponse) ProtoReflect() protoreflect.Message {
	mi := &file_portfolio_state_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveResponse.ProtoReflect.Descriptor instead.
func (*SaveResponse) Descriptor() ([]byte, []int) {
	return file_portfolio_state_service_proto_rawDescGZIP(), []int{2}
}

func (x *SaveResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *SaveResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type LoadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Date string `protobuf:"bytes,1,opt,name=date,proto3" json:"date,omitempty"`
}

func (x *LoadRequest) Reset() {
	*x = LoadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_portfolio_state_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoadRequest) ProtoMessage() {}

func (x *LoadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_portfolio_state_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoadRequest.ProtoReflect.Descriptor instead.
func (*LoadRequest) Descriptor() ([]byte, []int) {
	return file_portfolio_state_service_proto_rawDescGZIP(), []int{3}
}

func (x *LoadRequest) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

type HistoryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StartDate string `protobuf:"bytes,1,opt,name=start_date,json=startDate,proto3" json:"start_date,omitempty"`
	EndDate   string `protobuf:"bytes,2,opt,name=end_date,json=endDate,proto3" json:"end_date,omitempty"`
}

func (x *HistoryRequest) Reset() {
	*x = HistoryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_portfolio_state_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HistoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HistoryRequest) ProtoMessage() {}

func (x *HistoryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_portfolio_state_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HistoryRequest.ProtoReflect.Descriptor instead.
func (*HistoryRequest) Descriptor() ([]byte, []int) {
	return file_portfolio_state_service_proto_rawDescGZIP(), []int{4}
}

func (x *HistoryRequest) GetStartDate() string {
	if x != nil {
		return x.StartDate
	}
	return ""
}

func (x *HistoryRequest) GetEndDate() string {
	if x != nil {
		return x.EndDate
	}
	return ""
}

type PortfolioHistory struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	States []*PortfolioState `protobuf:"bytes,1,rep,name=states,proto3" json:"states,omitempty"`
}

func (x *PortfolioHistory) Reset() {
	*x = PortfolioHistory{}
	if protoimpl.UnsafeEnabled {
		mi := &file_portfolio_state_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PortfolioHistory) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PortfolioHistory) ProtoMessage() {}

func (x *PortfolioHistory) ProtoReflect() protoreflect.Message {
	mi := &file_portfolio_state_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PortfolioHistory.ProtoReflect.Descriptor instead.
func (*PortfolioHistory) Descriptor() ([]byte, []int) {
	return file_portfolio_state_service_proto_rawDescGZIP(), []int{5}
}

func (x *PortfolioHistory) GetStates() []*PortfolioState {
	if x != nil {
		return x.States
	}
	return nil
}

var File_portfolio_state_service_proto protoreflect.FileDescriptor

var file_portfolio_state_service_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x70, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x5f, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x15, 0x70, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22, 0xa7, 0x01, 0x0a, 0x0e, 0x50, 0x6f, 0x72, 0x74, 0x66,
	0x6f, 0x6c, 0x69, 0x6f, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x65, 0x12, 0x3d, 0x0a,
	0x09, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x1f, 0x2e, 0x70, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x09, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x21, 0x0a, 0x0c,
	0x63, 0x61, 0x73, 0x68, 0x5f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x0b, 0x63, 0x61, 0x73, 0x68, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x12,
	0x1f, 0x0a, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x22, 0xab, 0x01, 0x0a, 0x08, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73,
	0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x12, 0x23, 0x0a, 0x0d, 0x61, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x5f, 0x70, 0x72, 0x69,
	0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0c, 0x61, 0x76, 0x65, 0x72, 0x61, 0x67,
	0x65, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x74, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0c, 0x63,
	0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x6d,
	0x61, 0x72, 0x6b, 0x65, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x0b, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x42,
	0x0a, 0x0c, 0x53, 0x61, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x22, 0x21, 0x0a, 0x0b, 0x4c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x65, 0x22, 0x4a, 0x0a, 0x0e, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x44, 0x61, 0x74, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x64, 0x61,
	0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x44, 0x61, 0x74,
	0x65, 0x22, 0x51, 0x0a, 0x10, 0x50, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x48, 0x69,
	0x73, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x3d, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x65, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x70, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69,
	0x6f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x50, 0x6f,
	0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x73, 0x32, 0xc7, 0x02, 0x0a, 0x15, 0x50, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c,
	0x69, 0x6f, 0x53, 0x74, 0x61, 0x74, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x62,
	0x0a, 0x12, 0x53, 0x61, 0x76, 0x65, 0x50, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x12, 0x25, 0x2e, 0x70, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x50, 0x6f, 0x72,
	0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x53, 0x74, 0x61, 0x74, 0x65, 0x1a, 0x23, 0x2e, 0x70, 0x6f,
	0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x53, 0x61, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x61, 0x0a, 0x12, 0x4c, 0x6f, 0x61, 0x64, 0x50, 0x6f, 0x72, 0x74, 0x66, 0x6f,
	0x6c, 0x69, 0x6f, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x22, 0x2e, 0x70, 0x6f, 0x72, 0x74, 0x66,
	0x6f, 0x6c, 0x69, 0x6f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x4c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x70,
	0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x50, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x22, 0x00, 0x12, 0x67, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x50, 0x6f, 0x72, 0x74,
	0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x25, 0x2e, 0x70,
	0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e, 0x70, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x50, 0x6f, 0x72, 0x74,
	0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x22, 0x00, 0x42, 0x3d,
	0x5a, 0x3b, 0x6d, 0x6f, 0x6d, 0x65, 0x6e, 0x74, 0x75, 0x6d, 0x2d, 0x74, 0x72, 0x61, 0x64, 0x69,
	0x6e, 0x67, 0x2d, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x5f,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_portfolio_state_service_proto_rawDescOnce sync.Once
	file_portfolio_state_service_proto_rawDescData = file_portfolio_state_service_proto_rawDesc
)

func file_portfolio_state_service_proto_rawDescGZIP() []byte {
	file_portfolio_state_service_proto_rawDescOnce.Do(func() {
		file_portfolio_state_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_portfolio_state_service_proto_rawDescData)
	})
	return file_portfolio_state_service_proto_rawDescData
}

var file_portfolio_state_service_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_portfolio_state_service_proto_goTypes = []any{
	(*PortfolioState)(nil),   // 0: portfoliostateservice.PortfolioState
	(*Position)(nil),         // 1: portfoliostateservice.Position
	(*SaveResponse)(nil),     // 2: portfoliostateservice.SaveResponse
	(*LoadRequest)(nil),      // 3: portfoliostateservice.LoadRequest
	(*HistoryRequest)(nil),   // 4: portfoliostateservice.HistoryRequest
	(*PortfolioHistory)(nil), // 5: portfoliostateservice.PortfolioHistory
}
var file_portfolio_state_service_proto_depIdxs = []int32{
	1, // 0: portfoliostateservice.PortfolioState.positions:type_name -> portfoliostateservice.Position
	0, // 1: portfoliostateservice.PortfolioHistory.states:type_name -> portfoliostateservice.PortfolioState
	0, // 2: portfoliostateservice.PortfolioStateService.SavePortfolioState:input_type -> portfoliostateservice.PortfolioState
	3, // 3: portfoliostateservice.PortfolioStateService.LoadPortfolioState:input_type -> portfoliostateservice.LoadRequest
	4, // 4: portfoliostateservice.PortfolioStateService.GetPortfolioHistory:input_type -> portfoliostateservice.HistoryRequest
	2, // 5: portfoliostateservice.PortfolioStateService.SavePortfolioState:output_type -> portfoliostateservice.SaveResponse
	0, // 6: portfoliostateservice.PortfolioStateService.LoadPortfolioState:output_type -> portfoliostateservice.PortfolioState
	5, // 7: portfoliostateservice.PortfolioStateService.GetPortfolioHistory:output_type -> portfoliostateservice.PortfolioHistory
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_portfolio_state_service_proto_init() }
func file_portfolio_state_service_proto_init() {
	if File_portfolio_state_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_portfolio_state_service_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*PortfolioState); i {
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
		file_portfolio_state_service_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Position); i {
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
		file_portfolio_state_service_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*SaveResponse); i {
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
		file_portfolio_state_service_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*LoadRequest); i {
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
		file_portfolio_state_service_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*HistoryRequest); i {
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
		file_portfolio_state_service_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*PortfolioHistory); i {
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
			RawDescriptor: file_portfolio_state_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_portfolio_state_service_proto_goTypes,
		DependencyIndexes: file_portfolio_state_service_proto_depIdxs,
		MessageInfos:      file_portfolio_state_service_proto_msgTypes,
	}.Build()
	File_portfolio_state_service_proto = out.File
	file_portfolio_state_service_proto_rawDesc = nil
	file_portfolio_state_service_proto_goTypes = nil
	file_portfolio_state_service_proto_depIdxs = nil
}
