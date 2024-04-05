// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v3.12.4
// source: task_service.proto

package task_servicepb

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

type TaskID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *TaskID) Reset() {
	*x = TaskID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_task_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskID) ProtoMessage() {}

func (x *TaskID) ProtoReflect() protoreflect.Message {
	mi := &file_task_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskID.ProtoReflect.Descriptor instead.
func (*TaskID) Descriptor() ([]byte, []int) {
	return file_task_service_proto_rawDescGZIP(), []int{0}
}

func (x *TaskID) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type TaskContent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title           string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Description     string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Status          string `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
	CreatorUsername string `protobuf:"bytes,5,opt,name=creator_username,json=creatorUsername,proto3" json:"creator_username,omitempty"`
}

func (x *TaskContent) Reset() {
	*x = TaskContent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_task_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskContent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskContent) ProtoMessage() {}

func (x *TaskContent) ProtoReflect() protoreflect.Message {
	mi := &file_task_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskContent.ProtoReflect.Descriptor instead.
func (*TaskContent) Descriptor() ([]byte, []int) {
	return file_task_service_proto_rawDescGZIP(), []int{1}
}

func (x *TaskContent) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *TaskContent) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *TaskContent) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *TaskContent) GetCreatorUsername() string {
	if x != nil {
		return x.CreatorUsername
	}
	return ""
}

type Task struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   *TaskID      `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Task *TaskContent `protobuf:"bytes,2,opt,name=task,proto3" json:"task,omitempty"`
}

func (x *Task) Reset() {
	*x = Task{}
	if protoimpl.UnsafeEnabled {
		mi := &file_task_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Task) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Task) ProtoMessage() {}

func (x *Task) ProtoReflect() protoreflect.Message {
	mi := &file_task_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Task.ProtoReflect.Descriptor instead.
func (*Task) Descriptor() ([]byte, []int) {
	return file_task_service_proto_rawDescGZIP(), []int{2}
}

func (x *Task) GetId() *TaskID {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *Task) GetTask() *TaskContent {
	if x != nil {
		return x.Task
	}
	return nil
}

type TaskList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tasks    []*Task `protobuf:"bytes,1,rep,name=tasks,proto3" json:"tasks,omitempty"`
	PageSize int32   `protobuf:"varint,3,opt,name=pageSize,proto3" json:"pageSize,omitempty"`
}

func (x *TaskList) Reset() {
	*x = TaskList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_task_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskList) ProtoMessage() {}

func (x *TaskList) ProtoReflect() protoreflect.Message {
	mi := &file_task_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskList.ProtoReflect.Descriptor instead.
func (*TaskList) Descriptor() ([]byte, []int) {
	return file_task_service_proto_rawDescGZIP(), []int{3}
}

func (x *TaskList) GetTasks() []*Task {
	if x != nil {
		return x.Tasks
	}
	return nil
}

func (x *TaskList) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

type RequestByID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                *TaskID `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	RequestorUsername string  `protobuf:"bytes,2,opt,name=requestor_username,json=requestorUsername,proto3" json:"requestor_username,omitempty"`
}

func (x *RequestByID) Reset() {
	*x = RequestByID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_task_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestByID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestByID) ProtoMessage() {}

func (x *RequestByID) ProtoReflect() protoreflect.Message {
	mi := &file_task_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestByID.ProtoReflect.Descriptor instead.
func (*RequestByID) Descriptor() ([]byte, []int) {
	return file_task_service_proto_rawDescGZIP(), []int{4}
}

func (x *RequestByID) GetId() *TaskID {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *RequestByID) GetRequestorUsername() string {
	if x != nil {
		return x.RequestorUsername
	}
	return ""
}

type TaskPageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestorUsername string `protobuf:"bytes,1,opt,name=requestor_username,json=requestorUsername,proto3" json:"requestor_username,omitempty"`
	Offset            int32  `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	PageSize          int32  `protobuf:"varint,3,opt,name=pageSize,proto3" json:"pageSize,omitempty"`
}

func (x *TaskPageRequest) Reset() {
	*x = TaskPageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_task_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskPageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskPageRequest) ProtoMessage() {}

func (x *TaskPageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_task_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskPageRequest.ProtoReflect.Descriptor instead.
func (*TaskPageRequest) Descriptor() ([]byte, []int) {
	return file_task_service_proto_rawDescGZIP(), []int{5}
}

func (x *TaskPageRequest) GetRequestorUsername() string {
	if x != nil {
		return x.RequestorUsername
	}
	return ""
}

func (x *TaskPageRequest) GetOffset() int32 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *TaskPageRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

var File_task_service_proto protoreflect.FileDescriptor

var file_task_service_proto_rawDesc = []byte{
	0x0a, 0x12, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x22, 0x18, 0x0a, 0x06, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x44, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x88, 0x01, 0x0a,
	0x0b, 0x54, 0x61, 0x73, 0x6b, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x29, 0x0a, 0x10,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x55,
	0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x5b, 0x0a, 0x04, 0x54, 0x61, 0x73, 0x6b, 0x12,
	0x24, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x74, 0x61,
	0x73, 0x6b, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x49,
	0x44, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2d, 0x0a, 0x04, 0x74, 0x61, 0x73, 0x6b, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x52, 0x04,
	0x74, 0x61, 0x73, 0x6b, 0x22, 0x50, 0x0a, 0x08, 0x54, 0x61, 0x73, 0x6b, 0x4c, 0x69, 0x73, 0x74,
	0x12, 0x28, 0x0a, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x12, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x54,
	0x61, 0x73, 0x6b, 0x52, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61,
	0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61,
	0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x22, 0x62, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x42, 0x79, 0x49, 0x44, 0x12, 0x24, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x14, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x44, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2d, 0x0a, 0x12, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x6f, 0x72, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x74, 0x0a, 0x0f, 0x54, 0x61,
	0x73, 0x6b, 0x50, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2d, 0x0a,
	0x12, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x6f, 0x72, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6f, 0x66,
	0x66, 0x73, 0x65, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65,
	0x32, 0xd1, 0x02, 0x0a, 0x0b, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x3f, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x19,
	0x2e, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x54, 0x61,
	0x73, 0x6b, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x1a, 0x14, 0x2e, 0x74, 0x61, 0x73, 0x6b,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x44, 0x22,
	0x00, 0x12, 0x38, 0x0a, 0x0a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x12,
	0x12, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x54,
	0x61, 0x73, 0x6b, 0x1a, 0x14, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x44, 0x22, 0x00, 0x12, 0x3f, 0x0a, 0x0a, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x19, 0x2e, 0x74, 0x61, 0x73, 0x6b,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x42, 0x79, 0x49, 0x44, 0x1a, 0x14, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x44, 0x22, 0x00, 0x12, 0x3e, 0x0a, 0x0b,
	0x47, 0x65, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x42, 0x79, 0x49, 0x64, 0x12, 0x19, 0x2e, 0x74, 0x61,
	0x73, 0x6b, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x42, 0x79, 0x49, 0x44, 0x1a, 0x12, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x22, 0x00, 0x12, 0x46, 0x0a, 0x0b,
	0x47, 0x65, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x1d, 0x2e, 0x74, 0x61,
	0x73, 0x6b, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x50,
	0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x74, 0x61, 0x73,
	0x6b, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x4c, 0x69,
	0x73, 0x74, 0x22, 0x00, 0x42, 0x1e, 0x5a, 0x1c, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2f, 0x3b, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_task_service_proto_rawDescOnce sync.Once
	file_task_service_proto_rawDescData = file_task_service_proto_rawDesc
)

func file_task_service_proto_rawDescGZIP() []byte {
	file_task_service_proto_rawDescOnce.Do(func() {
		file_task_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_task_service_proto_rawDescData)
	})
	return file_task_service_proto_rawDescData
}

var file_task_service_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_task_service_proto_goTypes = []interface{}{
	(*TaskID)(nil),          // 0: task_service.TaskID
	(*TaskContent)(nil),     // 1: task_service.TaskContent
	(*Task)(nil),            // 2: task_service.Task
	(*TaskList)(nil),        // 3: task_service.TaskList
	(*RequestByID)(nil),     // 4: task_service.RequestByID
	(*TaskPageRequest)(nil), // 5: task_service.TaskPageRequest
}
var file_task_service_proto_depIdxs = []int32{
	0, // 0: task_service.Task.id:type_name -> task_service.TaskID
	1, // 1: task_service.Task.task:type_name -> task_service.TaskContent
	2, // 2: task_service.TaskList.tasks:type_name -> task_service.Task
	0, // 3: task_service.RequestByID.id:type_name -> task_service.TaskID
	1, // 4: task_service.TaskService.CreateTask:input_type -> task_service.TaskContent
	2, // 5: task_service.TaskService.UpdateTask:input_type -> task_service.Task
	4, // 6: task_service.TaskService.DeleteTask:input_type -> task_service.RequestByID
	4, // 7: task_service.TaskService.GetTaskById:input_type -> task_service.RequestByID
	5, // 8: task_service.TaskService.GetTaskList:input_type -> task_service.TaskPageRequest
	0, // 9: task_service.TaskService.CreateTask:output_type -> task_service.TaskID
	0, // 10: task_service.TaskService.UpdateTask:output_type -> task_service.TaskID
	0, // 11: task_service.TaskService.DeleteTask:output_type -> task_service.TaskID
	2, // 12: task_service.TaskService.GetTaskById:output_type -> task_service.Task
	3, // 13: task_service.TaskService.GetTaskList:output_type -> task_service.TaskList
	9, // [9:14] is the sub-list for method output_type
	4, // [4:9] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_task_service_proto_init() }
func file_task_service_proto_init() {
	if File_task_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_task_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskID); i {
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
		file_task_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskContent); i {
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
		file_task_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Task); i {
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
		file_task_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskList); i {
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
		file_task_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestByID); i {
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
		file_task_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskPageRequest); i {
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
			RawDescriptor: file_task_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_task_service_proto_goTypes,
		DependencyIndexes: file_task_service_proto_depIdxs,
		MessageInfos:      file_task_service_proto_msgTypes,
	}.Build()
	File_task_service_proto = out.File
	file_task_service_proto_rawDesc = nil
	file_task_service_proto_goTypes = nil
	file_task_service_proto_depIdxs = nil
}