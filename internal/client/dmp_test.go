// Package client ...
package client

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/abetterchoice/go-sdk/env"
	protocdmpproxyserver "github.com/abetterchoice/protoc_dmp_proxy_server"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
)

func TestNewDMPClient(t *testing.T) {
	type args struct {
		opts []DMPOption
	}
	tests := []struct {
		name string
		args args
		want DMPClient
	}{
		{
			name: "normal",
			args: args{opts: []DMPOption{}},
			want: &tabDMPClient{
				httpClient: &http.Client{
					Timeout: 10 * time.Second,
				},
				addr: env.DefaultDMPAddrPrd,
			},
		},
		{
			name: "test env",
			args: args{opts: []DMPOption{WithEnvTypeOption(env.TypeTest)}},
			want: &tabDMPClient{
				httpClient: &http.Client{
					Timeout: 10 * time.Second,
				},
				addr: env.DefaultDMPAddrTest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDMPClient(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDMPClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegisterDMPClient(t *testing.T) {
	type args struct {
		client DMPClient
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "normal",
			args: args{client: NewDMPClient()},
		},
		{
			name: "nil",
			args: args{client: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterDMPClient(tt.args.client)
		})
	}
}

func mockDMPServer(t gomock.TestReporter) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestBody, err := ioutil.ReadAll(request.Body)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		batchGetDMPTagResultReq := &protocdmpproxyserver.BatchGetDMPTagResultReq{}
		err = proto.Unmarshal(requestBody, batchGetDMPTagResultReq)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		var respList = make([]*protocdmpproxyserver.GetDMPTagResultResp, 0, len(batchGetDMPTagResultReq.ReqList))
		for _, getDMPTagResultReq := range batchGetDMPTagResultReq.ReqList {
			if getDMPTagResultReq == nil {
				continue
			}
			dmpResult := make(map[string]protocdmpproxyserver.StatusCode, len(getDMPTagResultReq.TagList))
			for _, tagCode := range getDMPTagResultReq.TagList {
				dmpResult[tagCode] = protocdmpproxyserver.StatusCode_STATUS_CODE_HIT
			}
			respList = append(respList, &protocdmpproxyserver.GetDMPTagResultResp{
				RetCode:         protocdmpproxyserver.RetCode_RET_CODE_SUCCESS,
				Message:         "mock success",
				UnitType:        getDMPTagResultReq.UnitType,
				UnitId:          getDMPTagResultReq.UnitId,
				DmpPlatformCode: getDMPTagResultReq.DmpPlatformCode,
				DmpResult:       dmpResult,
			})
		}
		resultBody, err := proto.Marshal(&protocdmpproxyserver.BatchGetDMPTagResultResp{RespList: respList})
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Write(resultBody)
	}))
	return ts
}

func mockFakeDMPServer(t gomock.TestReporter, h func(writer http.ResponseWriter, request *http.Request)) *httptest.Server {
	// mock http client
	ts := httptest.NewServer(http.HandlerFunc(h))
	return ts
}

func Test_BatchGetDMPResultHTTP(t *testing.T) {
	batchGetDMPTagResultURI = mockDMPServer(t).URL
	type fields struct {
		httpClient *http.Client
		addr       string
	}
	type args struct {
		ctx context.Context
		req *protocdmpproxyserver.BatchGetDMPTagResultReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *protocdmpproxyserver.BatchGetDMPTagResultResp
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
				ctx: context.Background(),
				req: &protocdmpproxyserver.BatchGetDMPTagResultReq{
					ReqList: []*protocdmpproxyserver.GetDMPTagResultReq{
						{
							UnitType:        1,
							DmpPlatformCode: 3,
						},
					},
				},
			},
			want: &protocdmpproxyserver.BatchGetDMPTagResultResp{RespList: []*protocdmpproxyserver.GetDMPTagResultResp{
				&protocdmpproxyserver.GetDMPTagResultResp{
					RetCode:         protocdmpproxyserver.RetCode_RET_CODE_SUCCESS,
					Message:         "mock success",
					UnitType:        1,
					DmpPlatformCode: 3,
				},
			}},
			wantErr: false,
		},
		{
			name: "marshal called with nil",
			fields: fields{
				httpClient: &http.Client{
					Timeout: 10 * time.Second,
				},
				addr: "",
			},
			args: args{
				ctx: context.Background(),
				req: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &tabDMPClient{
				httpClient: tt.fields.httpClient,
				addr:       tt.fields.addr,
			}
			got, err := c.BatchGetDMPResultHTTP(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("BatchGetDMPResultHTTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !proto.Equal(got, tt.want) {
				t.Errorf("BatchGetDMPResultHTTP() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_BatchGetDMPResultHTTPFailure(t *testing.T) {
	batchGetDMPTagResultURI = mockFakeDMPServer(t, func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
	}).URL
	type fields struct {
		httpClient *http.Client
		addr       string
	}
	type args struct {
		ctx context.Context
		req *protocdmpproxyserver.BatchGetDMPTagResultReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *protocdmpproxyserver.BatchGetDMPTagResultResp
		wantErr bool
	}{
		{
			name: "invalid http status",
			fields: fields{
				httpClient: &http.Client{
					Timeout: 10 * time.Second,
				},
				addr: "",
			},
			args: args{
				ctx: context.Background(),
				req: &protocdmpproxyserver.BatchGetDMPTagResultReq{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &tabDMPClient{
				httpClient: tt.fields.httpClient,
				addr:       tt.fields.addr,
			}
			got, err := c.BatchGetDMPResultHTTP(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("BatchGetDMPResultHTTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !proto.Equal(got, tt.want) {
				t.Errorf("BatchGetDMPResultHTTP() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_BatchGetDMPResultInvalidResp(t *testing.T) {
	batchGetDMPTagResultURI = mockFakeDMPServer(t, func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("xx"))
	}).URL
	type fields struct {
		httpClient *http.Client
		addr       string
	}
	type args struct {
		ctx context.Context
		req *protocdmpproxyserver.BatchGetDMPTagResultReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *protocdmpproxyserver.BatchGetDMPTagResultResp
		wantErr bool
	}{
		{
			name: "invalid resp",
			fields: fields{
				httpClient: &http.Client{
					Timeout: 10 * time.Second,
				},
				addr: "",
			},
			args: args{
				ctx: context.Background(),
				req: &protocdmpproxyserver.BatchGetDMPTagResultReq{},
			},
			want:    &protocdmpproxyserver.BatchGetDMPTagResultResp{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &tabDMPClient{
				httpClient: tt.fields.httpClient,
				addr:       tt.fields.addr,
			}
			got, err := c.BatchGetDMPTagResult(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("BatchGetDMPResultHTTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if proto.Equal(got, tt.want) { // proto.unmarshal 对非法的 body，转换不会报错，但得到的数据是非法的数据
				// 所以这里暂时判定 如果相等就是非法
				t.Errorf("BatchGetDMPResultHTTP() got = %v, want %v", got, tt.want)
			}
		})
	}
}
