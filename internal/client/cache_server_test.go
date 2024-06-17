// Package client ...
package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/abetterchoice/go-sdk/env"
	protoctabcacheserver "github.com/abetterchoice/protoc_cache_server"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
)

var (
	invalidHTTPStatus = "invalid http status"
)

func TestNewTABCacheClient(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want Client
	}{
		{
			name: "normal",
			args: args{opts: nil},
			want: &tabCacheClient{
				httpClient: &http.Client{
					Timeout: 10 * time.Second,
				},
				addr: env.GetAddr(env.TypePrd),
			},
		},
		{
			name: "with http client",
			args: args{opts: []Option{WithHTTPClient(&http.Client{Timeout: 1 * time.Second})}},
			want: &tabCacheClient{
				httpClient: &http.Client{
					Timeout: 1 * time.Second,
				},
				addr: env.GetAddr(env.TypePrd),
			},
		},
		{
			name: "with env type",
			args: args{opts: []Option{WithHTTPClient(&http.Client{Timeout: 1 * time.Second}), WithEnvType(env.TypeTest)}},
			want: &tabCacheClient{
				httpClient: &http.Client{
					Timeout: 1 * time.Second,
				},
				addr: env.GetAddr(env.TypeTest),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTABCacheClient(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTABCacheClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegisterCacheClient(t *testing.T) {
	type args struct {
		client Client
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "normal",
			args: args{client: NewTABCacheClient()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterCacheClient(tt.args.client)
		})
	}
}

func Test_BatchGetExperimentBucketInfo(t *testing.T) {
	batchGetExperimentBucketInfoURI = mockBatchGetExperimentBucketInfo(t).URL
	type fields struct {
		httpClient *http.Client
		addr       string
	}
	type args struct {
		ctx context.Context
		req *protoctabcacheserver.BatchGetExperimentBucketReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *protoctabcacheserver.BatchGetExperimentBucketResp
		wantErr bool
	}{
		{
			name: "bucketVersionIndex is required",
			fields: fields{
				httpClient: &http.Client{
					Timeout: 10 * time.Second,
				},
				addr: "",
			},
			args: args{
				ctx: context.TODO(),
				req: &protoctabcacheserver.BatchGetExperimentBucketReq{},
			},
			want: &protoctabcacheserver.BatchGetExperimentBucketResp{
				Code:        protoctabcacheserver.Code_CODE_SUCCESS,
				Message:     "empty resp",
				BucketIndex: make(map[int64]*protoctabcacheserver.BucketInfo),
			},
			wantErr: false,
		},
		{
			name: "normal",
			fields: fields{
				httpClient: &http.Client{
					Timeout: 10 * time.Second,
				},
				addr: "",
			},
			args: args{
				ctx: context.TODO(),
				req: &protoctabcacheserver.BatchGetExperimentBucketReq{
					ProjectId:  "mock 123",
					SdkVersion: env.SDKVersion,
					BucketVersionIndex: map[int64]string{
						1: "",
						2: "",
					},
				},
			},
			want: &protoctabcacheserver.BatchGetExperimentBucketResp{
				Code:    protoctabcacheserver.Code_CODE_SUCCESS,
				Message: "mock success",
				BucketIndex: map[int64]*protoctabcacheserver.BucketInfo{
					1: &protoctabcacheserver.BucketInfo{
						BucketType: protoctabcacheserver.BucketType_BUCKET_TYPE_RANGE,
						TrafficRange: &protoctabcacheserver.TrafficRange{
							Left:  100,
							Right: 10000,
						},
						ModifyType: protoctabcacheserver.ModifyType_MODIFY_UPDATE,
					},
					2: &protoctabcacheserver.BucketInfo{
						BucketType: protoctabcacheserver.BucketType_BUCKET_TYPE_RANGE,
						TrafficRange: &protoctabcacheserver.TrafficRange{
							Left:  0,
							Right: 99,
						},
						ModifyType: protoctabcacheserver.ModifyType_MODIFY_UPDATE,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &tabCacheClient{
				httpClient: tt.fields.httpClient,
				addr:       tt.fields.addr,
			}
			got, err := c.BatchGetExperimentBucketInfo(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("BatchGetExperimentBucketInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !proto.Equal(got, tt.want) {
				t.Errorf("BatchGetExperimentBucketInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_BatchGetExpBucketInfoFailure(t *testing.T) {
	batchGetExperimentBucketInfoURI = mockGetTabConfigFailure(t).URL
	type fields struct {
		httpClient *http.Client
		addr       string
	}
	type args struct {
		ctx context.Context
		req *protoctabcacheserver.BatchGetExperimentBucketReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *protoctabcacheserver.BatchGetExperimentBucketResp
		wantErr bool
	}{
		{
			name: invalidHTTPStatus,
			fields: fields{
				httpClient: &http.Client{
					Timeout: 10 * time.Second,
				},
				addr: "",
			},
			args: args{
				ctx: context.TODO(),
				req: &protoctabcacheserver.BatchGetExperimentBucketReq{
					ProjectId:  "mock 123",
					SdkVersion: env.SDKVersion,
					BucketVersionIndex: map[int64]string{
						1: "",
						2: "",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &tabCacheClient{
				httpClient: tt.fields.httpClient,
				addr:       tt.fields.addr,
			}
			got, err := c.BatchGetExperimentBucketInfo(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("BatchGetExperimentBucketInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !proto.Equal(got, tt.want) {
				t.Errorf("BatchGetExperimentBucketInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_BatchGetGroupBucketInfo(t *testing.T) {
	batchGetGroupBucketInfoURI = mockBatchGetGroupBucketInfo(t).URL
	type fields struct {
		httpClient *http.Client
		addr       string
	}
	type args struct {
		ctx context.Context
		req *protoctabcacheserver.BatchGetGroupBucketReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *protoctabcacheserver.BatchGetGroupBucketResp
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				httpClient: &http.Client{
					Timeout: 10 * time.Second,
				},
				addr: "",
			},
			args: args{
				ctx: context.TODO(),
				req: &protoctabcacheserver.BatchGetGroupBucketReq{
					ProjectId:  "mock 123",
					SdkVersion: env.SDKVersion,
					BucketVersionIndex: map[int64]string{
						1: "",
						2: "",
					},
				},
			},
			want: &protoctabcacheserver.BatchGetGroupBucketResp{
				Code:    protoctabcacheserver.Code_CODE_SUCCESS,
				Message: "mock success",
				BucketIndex: map[int64]*protoctabcacheserver.BucketInfo{
					1: &protoctabcacheserver.BucketInfo{
						BucketType: protoctabcacheserver.BucketType_BUCKET_TYPE_RANGE,
						TrafficRange: &protoctabcacheserver.TrafficRange{
							Left:  100,
							Right: 10000,
						},
						ModifyType: protoctabcacheserver.ModifyType_MODIFY_UPDATE,
					},
					2: &protoctabcacheserver.BucketInfo{
						BucketType: protoctabcacheserver.BucketType_BUCKET_TYPE_RANGE,
						TrafficRange: &protoctabcacheserver.TrafficRange{
							Left:  0,
							Right: 99,
						},
						ModifyType: protoctabcacheserver.ModifyType_MODIFY_UPDATE,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "bucketVersionIndex is required",
			fields: fields{
				httpClient: &http.Client{
					Timeout: 10 * time.Second,
				},
				addr: "",
			},
			args: args{
				ctx: context.TODO(),
				req: &protoctabcacheserver.BatchGetGroupBucketReq{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &tabCacheClient{
				httpClient: tt.fields.httpClient,
				addr:       tt.fields.addr,
			}
			got, err := c.BatchGetGroupBucketInfo(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("BatchGetGroupBucketInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !proto.Equal(got, tt.want) {
				t.Errorf("BatchGetGroupBucketInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_BatchGetGroupBucketFailure(t *testing.T) {
	batchGetGroupBucketInfoURI = mockGetTabConfigFailure(t).URL
	type fields struct {
		httpClient *http.Client
		addr       string
	}
	type args struct {
		ctx context.Context
		req *protoctabcacheserver.BatchGetGroupBucketReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *protoctabcacheserver.BatchGetGroupBucketResp
		wantErr bool
	}{
		{
			name: invalidHTTPStatus,
			fields: fields{
				httpClient: &http.Client{
					Timeout: 10 * time.Second,
				},
				addr: "",
			},
			args: args{
				ctx: context.TODO(),
				req: &protoctabcacheserver.BatchGetGroupBucketReq{
					ProjectId:  "mock 123",
					SdkVersion: env.SDKVersion,
					BucketVersionIndex: map[int64]string{
						1: "",
						2: "",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &tabCacheClient{
				httpClient: tt.fields.httpClient,
				addr:       tt.fields.addr,
			}
			got, err := c.BatchGetGroupBucketInfo(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("BatchGetExperimentBucketInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !proto.Equal(got, tt.want) {
				t.Errorf("BatchGetExperimentBucketInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetTabConfigData(t *testing.T) {
	getTabConfigURI = mockGetTabConfig(t).URL
	type fields struct {
		httpClient *http.Client
		addr       string
	}
	type args struct {
		ctx context.Context
		req *protoctabcacheserver.GetTabConfigReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *protoctabcacheserver.GetTabConfigResp
		wantErr bool
	}{
		{
			name: "request is required",
			fields: fields{
				httpClient: nil,
				addr:       "",
			},
			args: args{
				ctx: nil,
				req: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "normal",
			fields: fields{
				httpClient: &http.Client{
					Timeout: 10 * time.Second,
				},
				addr: "",
			},
			args: args{
				ctx: context.TODO(),
				req: &protoctabcacheserver.GetTabConfigReq{},
			},
			want: &protoctabcacheserver.GetTabConfigResp{
				Code:    protoctabcacheserver.Code_CODE_SUCCESS,
				Message: "mock success",
				TabConfigManager: &protoctabcacheserver.TabConfigManager{
					ProjectId:  "mock 123",
					UpdateType: protoctabcacheserver.UpdateType_UPDATE_TYPE_COMPLETE,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &tabCacheClient{
				httpClient: tt.fields.httpClient,
				addr:       tt.fields.addr,
			}
			got, err := c.GetTabConfigData(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTabConfigData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !proto.Equal(got, tt.want) {
				t.Errorf("GetTabConfigData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetTabConfigDataFailure(t *testing.T) {
	getTabConfigURI = mockGetTabConfigFailure(t).URL
	type fields struct {
		httpClient *http.Client
		addr       string
	}
	type args struct {
		ctx context.Context
		req *protoctabcacheserver.GetTabConfigReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *protoctabcacheserver.GetTabConfigResp
		wantErr bool
	}{
		{
			name: invalidHTTPStatus,
			fields: fields{
				httpClient: &http.Client{
					Timeout: 10 * time.Second,
				},
				addr: "",
			},
			args: args{
				ctx: context.TODO(),
				req: &protoctabcacheserver.GetTabConfigReq{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &tabCacheClient{
				httpClient: tt.fields.httpClient,
				addr:       tt.fields.addr,
			}
			got, err := c.GetTabConfigData(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTabConfigData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !proto.Equal(got, tt.want) {
				t.Errorf("GetTabConfigData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockGetTabConfig(t gomock.TestReporter) *httptest.Server {
	h := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		resp := &protoctabcacheserver.GetTabConfigResp{
			Code:    protoctabcacheserver.Code_CODE_SUCCESS,
			Message: "mock success",
			TabConfigManager: &protoctabcacheserver.TabConfigManager{
				ProjectId:  "mock 123",
				Version:    "",
				UpdateType: protoctabcacheserver.UpdateType_UPDATE_TYPE_COMPLETE,
				TabConfig:  nil,
			},
		}
		body, err := proto.Marshal(resp)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Write(body)
	})
	ts := httptest.NewServer(h)
	return ts
}

func mockGetTabConfigFailure(t gomock.TestReporter) *httptest.Server {
	h := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	})
	ts := httptest.NewServer(h)
	return ts
}

func mockBatchGetExperimentBucketInfo(t gomock.TestReporter) *httptest.Server {
	h := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		resp := &protoctabcacheserver.BatchGetExperimentBucketResp{
			Code:    protoctabcacheserver.Code_CODE_SUCCESS,
			Message: "mock success",
			BucketIndex: map[int64]*protoctabcacheserver.BucketInfo{
				1: &protoctabcacheserver.BucketInfo{
					BucketType: protoctabcacheserver.BucketType_BUCKET_TYPE_RANGE,
					TrafficRange: &protoctabcacheserver.TrafficRange{
						Left:  100,
						Right: 10000,
					},
					ModifyType: protoctabcacheserver.ModifyType_MODIFY_UPDATE,
				},
				2: &protoctabcacheserver.BucketInfo{
					BucketType: protoctabcacheserver.BucketType_BUCKET_TYPE_RANGE,
					TrafficRange: &protoctabcacheserver.TrafficRange{
						Left:  0,
						Right: 99,
					},
					ModifyType: protoctabcacheserver.ModifyType_MODIFY_UPDATE,
				},
			},
		}
		body, err := proto.Marshal(resp)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Write(body)
	})
	ts := httptest.NewServer(h)
	return ts
}

func mockBatchGetGroupBucketInfo(t gomock.TestReporter) *httptest.Server {
	h := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		resp := &protoctabcacheserver.BatchGetGroupBucketResp{
			Code:    protoctabcacheserver.Code_CODE_SUCCESS,
			Message: "mock success",
			BucketIndex: map[int64]*protoctabcacheserver.BucketInfo{
				1: &protoctabcacheserver.BucketInfo{
					BucketType: protoctabcacheserver.BucketType_BUCKET_TYPE_RANGE,
					TrafficRange: &protoctabcacheserver.TrafficRange{
						Left:  100,
						Right: 10000,
					},
					ModifyType: protoctabcacheserver.ModifyType_MODIFY_UPDATE,
				},
				2: &protoctabcacheserver.BucketInfo{
					BucketType: protoctabcacheserver.BucketType_BUCKET_TYPE_RANGE,
					TrafficRange: &protoctabcacheserver.TrafficRange{
						Left:  0,
						Right: 99,
					},
					ModifyType: protoctabcacheserver.ModifyType_MODIFY_UPDATE,
				},
			},
		}
		body, err := proto.Marshal(resp)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Write(body)
	})
	ts := httptest.NewServer(h)
	return ts
}
