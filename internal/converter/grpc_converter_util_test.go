package converter

import (
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/opendatahub-io/model-registry/internal/model/db"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestConvertArtifact_State(t *testing.T) {
	type args struct {
		source *proto.Artifact_State
	}
	pending := proto.Artifact_PENDING
	var pendingInt int8 = 1
	tests := []struct {
		name string
		args args
		want *int8
	}{
		{
			name: "valid state",
			args: args{source: &pending},
			want: &pendingInt,
		},
		{
			name: "nil state",
			args: args{source: nil},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertArtifact_State(tt.args.source); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertArtifact_State() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertToArtifact_State(t *testing.T) {
	type args struct {
		source *int8
	}
	pendingState := int8(db.PENDING)
	pendingProto := proto.Artifact_PENDING
	invalidState := int8(-1)
	tests := []struct {
		name    string
		args    args
		want    *proto.Artifact_State
		wantErr bool
	}{
		{
			name:    "valid state",
			args:    args{source: &pendingState},
			want:    &pendingProto,
			wantErr: false,
		},
		{
			name:    "invalid state",
			args:    args{source: &invalidState},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "nil state",
			args:    args{source: nil},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToArtifact_State(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToArtifact_State() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToArtifact_State() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertExecution_State(t *testing.T) {
	type args struct {
		source *proto.Execution_State
	}
	newState := proto.Execution_NEW
	newResult := int8(db.NEW)
	tests := []struct {
		name string
		args args
		want *int8
	}{
		{
			name: "valid state",
			args: args{source: &newState},
			want: &newResult,
		},
		{
			name: "nil state",
			args: args{source: nil},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertExecution_State(tt.args.source); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertExecution_State() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertToExecution_State(t *testing.T) {
	type args struct {
		source *int8
	}
	newState := int8(db.NEW)
	newProtoState := proto.Execution_NEW
	invalidState := int8(-1)
	tests := []struct {
		name    string
		args    args
		want    *proto.Execution_State
		wantErr bool
	}{
		{
			name:    "valid state",
			args:    args{source: &newState},
			want:    &newProtoState,
			wantErr: false,
		},
		{
			name:    "invalid state",
			args:    args{source: &invalidState},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "nil state",
			args:    args{source: nil},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToExecution_State(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToExecution_State() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToExecution_State() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertProtoArtifactProperties(t *testing.T) {
	type args struct {
		source *proto.Artifact
	}
	testStr := "test"
	intValue := int64(10)
	boolValue := true
	doubleValue := 98.6
	tests := []struct {
		name    string
		args    args
		want    []db.ArtifactProperty
		wantErr bool
	}{
		{
			name: "nil artifact id",
			args: args{source: &proto.Artifact{
				Id:               nil,
				Properties:       map[string]*proto.Value{"prop1": {Value: &proto.Value_StringValue{StringValue: testStr}}},
				CustomProperties: map[string]*proto.Value{"prop2": {Value: &proto.Value_IntValue{IntValue: intValue}}},
			}},
			want: []db.ArtifactProperty{
				{Name: "prop1", StringValue: &testStr},
				{Name: "prop2", IsCustomProperty: true, IntValue: &intValue},
			},
			wantErr: false,
		},
		{
			name: "valid artifact id",
			args: args{source: &proto.Artifact{
				Id:               &intValue,
				Properties:       map[string]*proto.Value{"prop1": {Value: &proto.Value_DoubleValue{DoubleValue: doubleValue}}},
				CustomProperties: map[string]*proto.Value{"prop2": {Value: &proto.Value_BoolValue{BoolValue: boolValue}}},
			}},
			want: []db.ArtifactProperty{
				{ArtifactID: intValue, Name: "prop1", DoubleValue: &doubleValue},
				{ArtifactID: intValue, Name: "prop2", IsCustomProperty: true, BoolValue: &boolValue},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertProtoArtifactProperties(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertProtoArtifactProperties() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertProtoArtifactProperties() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertProtoContextProperties(t *testing.T) {
	type args struct {
		source *proto.Context
	}
	testStr := "test"
	intValue := int64(10)
	tests := []struct {
		name    string
		args    args
		want    []db.ContextProperty
		wantErr bool
	}{
		{
			name: "nil context id",
			args: args{source: &proto.Context{
				Id:               nil,
				Properties:       map[string]*proto.Value{"prop1": {Value: &proto.Value_StringValue{StringValue: testStr}}},
				CustomProperties: map[string]*proto.Value{"prop2": {Value: &proto.Value_IntValue{IntValue: intValue}}},
			}},
			want: []db.ContextProperty{
				{Name: "prop1", StringValue: &testStr},
				{Name: "prop2", IsCustomProperty: true, IntValue: &intValue},
			},
			wantErr: false,
		},
		{
			name: "valid context id",
			args: args{source: &proto.Context{
				Id:               &intValue,
				Properties:       map[string]*proto.Value{"prop1": {Value: &proto.Value_StringValue{StringValue: testStr}}},
				CustomProperties: map[string]*proto.Value{"prop2": {Value: &proto.Value_IntValue{IntValue: intValue}}},
			}},
			want: []db.ContextProperty{
				{ContextID: intValue, Name: "prop1", StringValue: &testStr},
				{ContextID: intValue, Name: "prop2", IsCustomProperty: true, IntValue: &intValue},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertProtoContextProperties(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertProtoContextProperties() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertProtoContextProperties() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertProtoEventPath(t *testing.T) {
	type args struct {
		source *proto.Event_Path
	}
	stepIndex := 1
	stepKey := "my-key"
	tests := []struct {
		name    string
		args    args
		want    []db.EventPath
		wantErr bool
	}{
		{
			name:    "step index",
			args:    args{source: &proto.Event_Path{Steps: []*proto.Event_Path_Step{{Value: &proto.Event_Path_Step_Index{Index: int64(stepIndex)}}}}},
			want:    []db.EventPath{{IsIndexStep: true, StepIndex: &stepIndex}},
			wantErr: false,
		},
		{
			name:    "step key",
			args:    args{source: &proto.Event_Path{Steps: []*proto.Event_Path_Step{{Value: &proto.Event_Path_Step_Key{Key: stepKey}}}}},
			want:    []db.EventPath{{StepKey: &stepKey}},
			wantErr: false,
		},
		{
			name:    "nil path",
			args:    args{},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertProtoEventPath(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertProtoEventPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertProtoEventPath() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertProtoEventType(t *testing.T) {
	type args struct {
		source *proto.Event_Type
	}
	validEventType := proto.Event_DECLARED_INPUT
	tests := []struct {
		name string
		args args
		want int8
	}{
		{
			name: "valid proto event type",
			args: args{source: &validEventType},
			want: 2,
		},
		{
			name: "nil proto event type",
			args: args{},
			want: int8(db.EVENT_TYPE_UNKNOWN),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertProtoEventType(tt.args.source); got != tt.want {
				t.Errorf("ConvertProtoEventType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertProtoExecutionProperties(t *testing.T) {
	type args struct {
		source *proto.Execution
	}
	testStr := "test"
	intValue := int64(10)
	tests := []struct {
		name    string
		args    args
		want    []db.ExecutionProperty
		wantErr bool
	}{
		{
			name: "nil execution id",
			args: args{source: &proto.Execution{
				Id:               nil,
				Properties:       map[string]*proto.Value{"prop1": {Value: &proto.Value_StringValue{StringValue: testStr}}},
				CustomProperties: map[string]*proto.Value{"prop2": {Value: &proto.Value_IntValue{IntValue: intValue}}},
			}},
			want: []db.ExecutionProperty{
				{Name: "prop1", StringValue: &testStr},
				{Name: "prop2", IsCustomProperty: true, IntValue: &intValue},
			},
			wantErr: false,
		},
		{
			name: "valid execution id",
			args: args{source: &proto.Execution{
				Id:               &intValue,
				Properties:       map[string]*proto.Value{"prop1": {Value: &proto.Value_StringValue{StringValue: testStr}}},
				CustomProperties: map[string]*proto.Value{"prop2": {Value: &proto.Value_IntValue{IntValue: intValue}}},
			}},
			want: []db.ExecutionProperty{
				{ExecutionID: intValue, Name: "prop1", StringValue: &testStr},
				{ExecutionID: intValue, Name: "prop2", IsCustomProperty: true, IntValue: &intValue},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertProtoExecutionProperties(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertProtoExecutionProperties() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertProtoExecutionProperties() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertToProtoArtifactCustomProperties(t *testing.T) {
	type args struct {
		source []db.ArtifactProperty
	}
	intValue := int64(1)
	stringValue := "test"
	bytes := []byte(stringValue)
	tests := []struct {
		name    string
		args    args
		want    map[string]*proto.Value
		wantErr bool
	}{
		{
			name: "valid custom properties",
			args: args{source: []db.ArtifactProperty{
				{Name: "notCustomProp", IntValue: &intValue},
				{Name: "customProp", IsCustomProperty: true, StringValue: &stringValue},
			}},
			want:    map[string]*proto.Value{"customProp": {Value: &proto.Value_StringValue{StringValue: stringValue}}},
			wantErr: false,
		},
		{
			name: "invalid bytes property",
			args: args{source: []db.ArtifactProperty{
				{Name: "notCustomProp", IntValue: &intValue},
				{Name: "customProp", IsCustomProperty: true, ByteValue: &bytes},
			}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "no properties",
			want:    map[string]*proto.Value{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToProtoArtifactCustomProperties(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToProtoArtifactCustomProperties() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToProtoArtifactCustomProperties() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertToProtoArtifactProperties(t *testing.T) {
	type args struct {
		source []db.ArtifactProperty
	}
	doubleValue := 98.6
	boolValue := true
	bytes := []byte("test")
	tests := []struct {
		name    string
		args    args
		want    map[string]*proto.Value
		wantErr bool
	}{
		{
			name: "valid non-custom properties",
			args: args{source: []db.ArtifactProperty{
				{Name: "notCustomProp", DoubleValue: &doubleValue},
				{Name: "customProp", IsCustomProperty: true, BoolValue: &boolValue},
			}},
			want:    map[string]*proto.Value{"notCustomProp": {Value: &proto.Value_DoubleValue{DoubleValue: doubleValue}}},
			wantErr: false,
		},
		{
			name: "invalid bytes property",
			args: args{source: []db.ArtifactProperty{
				{Name: "notCustomProp", ByteValue: &bytes},
				{Name: "customProp", IsCustomProperty: true, DoubleValue: &doubleValue},
			}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "no properties",
			want:    map[string]*proto.Value{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToProtoArtifactProperties(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToProtoArtifactProperties() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToProtoArtifactProperties() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertToProtoContextCustomProperties(t *testing.T) {
	type args struct {
		source []db.ContextProperty
	}
	intValue := int64(1)
	stringValue := "test"
	bytes := []byte(stringValue)
	boolValue := true
	tests := []struct {
		name    string
		args    args
		want    map[string]*proto.Value
		wantErr bool
	}{
		{
			name: "valid custom properties",
			args: args{source: []db.ContextProperty{
				{Name: "notCustomProp", IntValue: &intValue},
				{Name: "customProp", IsCustomProperty: true, BoolValue: &boolValue},
			}},
			want:    map[string]*proto.Value{"customProp": {Value: &proto.Value_BoolValue{BoolValue: boolValue}}},
			wantErr: false,
		},
		{
			name: "invalid bytes property",
			args: args{source: []db.ContextProperty{
				{Name: "notCustomProp", IntValue: &intValue},
				{Name: "customProp", IsCustomProperty: true, ByteValue: &bytes},
			}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "no properties",
			want:    map[string]*proto.Value{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToProtoContextCustomProperties(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToProtoContextCustomProperties() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToProtoContextCustomProperties() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertToProtoContextProperties(t *testing.T) {
	type args struct {
		source []db.ContextProperty
	}
	doubleValue := 98.6
	boolValue := true
	bytes := []byte("test")
	typeUrl := "url://test"
	tests := []struct {
		name    string
		args    args
		want    map[string]*proto.Value
		wantErr bool
	}{
		{
			name: "valid non-custom properties",
			args: args{source: []db.ContextProperty{
				{Name: "notCustomProp", ProtoValue: &bytes, TypeURL: &typeUrl},
				{Name: "customProp", IsCustomProperty: true, BoolValue: &boolValue},
			}},
			want: map[string]*proto.Value{
				"notCustomProp": {Value: &proto.Value_ProtoValue{ProtoValue: &anypb.Any{TypeUrl: typeUrl, Value: bytes}}},
			},
			wantErr: false,
		},
		{
			name: "invalid bytes property",
			args: args{source: []db.ContextProperty{
				{Name: "notCustomProp", ByteValue: &bytes},
				{Name: "customProp", IsCustomProperty: true, DoubleValue: &doubleValue},
			}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "no properties",
			want:    map[string]*proto.Value{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToProtoContextProperties(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToProtoContextProperties() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToProtoContextProperties() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertToProtoEventPath(t *testing.T) {
	type args struct {
		source []db.EventPath
	}
	stepIndex := 1
	stepKey := "test"
	tests := []struct {
		name    string
		args    args
		want    *proto.Event_Path
		wantErr bool
	}{
		{
			name: "valid path index",
			args: args{source: []db.EventPath{{IsIndexStep: true, StepIndex: &stepIndex}}},
			want: &proto.Event_Path{Steps: []*proto.Event_Path_Step{
				{Value: &proto.Event_Path_Step_Index{Index: int64(stepIndex)}},
			}},
			wantErr: false,
		},
		{
			name: "valid path key",
			args: args{source: []db.EventPath{{StepKey: &stepKey}}},
			want: &proto.Event_Path{Steps: []*proto.Event_Path_Step{
				{Value: &proto.Event_Path_Step_Key{Key: stepKey}},
			}},
			wantErr: false,
		},
		{
			name: "nil path",
			want: nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToProtoEventPath(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToProtoEventPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToProtoEventPath() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertToProtoEventType(t *testing.T) {
	type args struct {
		source int8
	}
	validEventType := proto.Event_DECLARED_INPUT
	unknownEventType := proto.Event_UNKNOWN
	tests := []struct {
		name    string
		args    args
		want    *proto.Event_Type
		wantErr bool
	}{
		{
			name:    "valid event type",
			args:    args{source: int8(db.DECLARED_INPUT)},
			want:    &validEventType,
			wantErr: false,
		},
		{
			name:    "nil event type",
			args:    args{},
			want:    &unknownEventType,
			wantErr: false,
		},
		{
			name:    "invalid event type",
			args:    args{source: int8(-1)},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToProtoEventType(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToProtoEventType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToProtoEventType() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertToProtoExecutionCustomProperties(t *testing.T) {
	type args struct {
		source []db.ExecutionProperty
	}
	intValue := int64(1)
	stringValue := "test"
	bytes := []byte(stringValue)
	boolValue := true
	tests := []struct {
		name    string
		args    args
		want    map[string]*proto.Value
		wantErr bool
	}{
		{
			name: "valid custom properties",
			args: args{source: []db.ExecutionProperty{
				{Name: "notCustomProp", IntValue: &intValue},
				{Name: "customProp", IsCustomProperty: true, BoolValue: &boolValue},
			}},
			want:    map[string]*proto.Value{"customProp": {Value: &proto.Value_BoolValue{BoolValue: boolValue}}},
			wantErr: false,
		},
		{
			name: "invalid bytes property",
			args: args{source: []db.ExecutionProperty{
				{Name: "notCustomProp", IntValue: &intValue},
				{Name: "customProp", IsCustomProperty: true, ByteValue: &bytes},
			}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "no properties",
			want:    map[string]*proto.Value{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToProtoExecutionCustomProperties(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToProtoExecutionCustomProperties() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToProtoExecutionCustomProperties() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertToProtoExecutionProperties(t *testing.T) {
	type args struct {
		source []db.ExecutionProperty
	}
	doubleValue := 98.6
	boolValue := true
	bytes := []byte("test")
	typeUrl := "url://test"
	tests := []struct {
		name    string
		args    args
		want    map[string]*proto.Value
		wantErr bool
	}{
		{
			name: "valid non-custom properties",
			args: args{source: []db.ExecutionProperty{
				{Name: "notCustomProp", ProtoValue: &bytes, TypeURL: &typeUrl},
				{Name: "customProp", IsCustomProperty: true, BoolValue: &boolValue},
			}},
			want: map[string]*proto.Value{
				"notCustomProp": {Value: &proto.Value_ProtoValue{ProtoValue: &anypb.Any{TypeUrl: typeUrl, Value: bytes}}},
			},
			wantErr: false,
		},
		{
			name: "invalid bytes property",
			args: args{source: []db.ExecutionProperty{
				{Name: "notCustomProp", ByteValue: &bytes},
				{Name: "customProp", IsCustomProperty: true, DoubleValue: &doubleValue},
			}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "no properties",
			want:    map[string]*proto.Value{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToProtoExecutionProperties(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToProtoExecutionProperties() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToProtoExecutionProperties() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertTypeIDToName(t *testing.T) {
	// setup mock DB
	logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)
	dbConn, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger})
	if err != nil {
		t.Errorf("failed to create mock Gorm connection: %v", err)
	}
	_ = dbConn.Migrator().AutoMigrate(db.Type{})
	_ = dbConn.Exec("INSERT INTO Type VALUES(1,'mlmd.Dataset',NULL,1,NULL,NULL,NULL,NULL)").Error
	_ = dbConn.Exec("INSERT INTO Type VALUES(2,'mlmd.Model',NULL,1,NULL,NULL,NULL,NULL)").Error
	SetConverterDB(dbConn)

	type args struct {
		id int64
	}
	dataSet := "mlmd.Dataset"
	model := "mlmd.Model"
	tests := []struct {
		name    string
		args    args
		want    *string
		wantErr bool
	}{
		{
			name:    dataSet,
			args:    args{id: 1},
			want:    &dataSet,
			wantErr: false,
		},
		{
			name:    model,
			args:    args{id: 2},
			want:    &model,
			wantErr: false,
		},
		{
			name:    "missing id",
			args:    args{id: 3},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertTypeIDToName(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertTypeIDToName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertTypeIDToName() got = %v, want %v", got, tt.want)
			}
		})
	}
}
