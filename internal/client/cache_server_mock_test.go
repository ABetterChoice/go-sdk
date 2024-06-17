// Package client ...
package client

import (
	"context"
	"testing"

	"github.com/abetterchoice/protoc_cache_server"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
)

func TestBatchGetExperimentBucketInfo(t *testing.T) {
	mockClient := NewMockClient(gomock.NewController(t))
	mockClient.EXPECT().BatchGetExperimentBucketInfo(gomock.Any(), gomock.Any()).Return(&protoc_cache_server.BatchGetExperimentBucketResp{
		Code:        protoc_cache_server.Code_CODE_SUCCESS,
		Message:     "mock success",
		BucketIndex: nil,
	}, nil)
	tests := []struct {
		name    string
		req     *protoc_cache_server.BatchGetExperimentBucketReq
		want    *protoc_cache_server.BatchGetExperimentBucketResp
		wantErr bool
	}{
		{
			name: "normal",
			want: &protoc_cache_server.BatchGetExperimentBucketResp{
				Code:        protoc_cache_server.Code_CODE_SUCCESS,
				Message:     "mock success",
				BucketIndex: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mockClient.BatchGetExperimentBucketInfo(context.TODO(), tt.req)
			if !proto.Equal(got, tt.want) {
				t.Errorf("BatchGetExperimentBucketInfo() = %v, want %v", got, tt.want)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("BatchGetExperimentBucketInfo(); err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBatchGetGroupBucketInfo(t *testing.T) {
	mockClient := NewMockClient(gomock.NewController(t))
	mockClient.EXPECT().BatchGetGroupBucketInfo(gomock.Any(), gomock.Any()).Return(&protoc_cache_server.BatchGetGroupBucketResp{
		Code:        protoc_cache_server.Code_CODE_SUCCESS,
		Message:     "mock success",
		BucketIndex: nil,
	}, nil)
	tests := []struct {
		name    string
		req     *protoc_cache_server.BatchGetGroupBucketReq
		want    *protoc_cache_server.BatchGetGroupBucketResp
		wantErr bool
	}{
		{
			name: "normal",
			want: &protoc_cache_server.BatchGetGroupBucketResp{
				Code:        protoc_cache_server.Code_CODE_SUCCESS,
				Message:     "mock success",
				BucketIndex: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := mockClient.BatchGetGroupBucketInfo(context.TODO(), tt.req)
			if !proto.Equal(got, tt.want) {
				t.Errorf("BatchGetGroupBucketInfo() = %v, want %v", got, tt.want)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("BatchGetGroupBucketInfo(); err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetTabConfigData(t *testing.T) {
	mockClient := NewMockClient(gomock.NewController(t))
	mockClient.EXPECT().GetTabConfigData(gomock.Any(), gomock.Any()).Return(&protoc_cache_server.GetTabConfigResp{
		Code:    protoc_cache_server.Code_CODE_SUCCESS,
		Message: "mock success",
	}, nil)
	tests := []struct {
		name    string
		req     *protoc_cache_server.GetTabConfigReq
		want    *protoc_cache_server.GetTabConfigResp
		wantErr bool
	}{
		{
			name: "normal",
			want: &protoc_cache_server.GetTabConfigResp{
				Code:    protoc_cache_server.Code_CODE_SUCCESS,
				Message: "mock success",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := mockClient.GetTabConfigData(context.TODO(), tt.req)
			if !proto.Equal(got, tt.want) {
				t.Errorf("GetTabConfigData() = %v, want %v", got, tt.want)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTabConfigData(); err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
