package service

import (
	"errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// helper for get extension
// sample to get IssueService.GetIssue's method_signature is below
//     ref := (&pb.Issue{}).ProtoReflect()
//     ms, err := GetMethodExt(ref, "IssueService", "GetIssue", pb.E_MethodSignature)
//     "issue_id" == ms[0] // should be true
func GetMethodExt(x protoreflect.Message, serviceName, methodName string, xt protoreflect.ExtensionType) (*[]string, error) {
	svs := x.Descriptor().ParentFile().Services().ByName(protoreflect.Name(serviceName))
	if svs == nil {
		return nil, errors.New("service not found")
	}
	met := svs.Methods().ByName(protoreflect.Name(methodName))
	if met == nil {
		return nil, errors.New("method not found")
	}
	opts, ok := met.Options().(*descriptorpb.MethodOptions)
	if !ok {
		return nil, errors.New("convert Options failed")
	}
	ext, ok := proto.GetExtension(opts, xt).([]string)
	if !ok {
		return nil, errors.New("convert Ext failed")
	}
	return &ext, nil
}
