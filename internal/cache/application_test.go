// Package cache ...
package cache

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/RoaringBitmap/roaring"
	"github.com/abetterchoice/go-sdk/internal/client"
	"github.com/abetterchoice/go-sdk/testdata"
	protoctabcacheserver "github.com/abetterchoice/protoc_cache_server"
)

var (
	projectID     = "123"
	projectIDList = []string{"123"}
)

func TestGetApplication(t *testing.T) {
	defer Release()
	type args struct {
		projectID string
	}
	tests := []struct {
		name string
		args args
		want *Application
	}{
		{
			name: "not found",
			args: args{projectID: projectID},
			want: nil,
		},
		{
			name: "normal",
			args: args{projectID: "mock123"},
			want: &Application{
				ProjectID: "mock123",
			},
		},
		{
			name: "fake data",
			args: args{projectID: "fake123"},
			want: nil,
		},
	}
	setApplication(&Application{
		ProjectID: "mock123",
	})
	setApplication(nil)
	localApplicationCache.Store("fake123", "fake123")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetApplication(tt.args.projectID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetApplication() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitLocalCache(t *testing.T) {
	defer func() {
		client.CacheClient = nil
	}()
	client.CacheClient = testdata.MockCacheClient(t)
	type args struct {
		ctx           context.Context
		projectIDList []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				ctx:           context.Background(),
				projectIDList: projectIDList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitLocalCache(tt.args.ctx, tt.args.projectIDList); (err != nil) != tt.wantErr {
				t.Errorf("InitLocalCache() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInitLocalCacheFailure(t *testing.T) {
	defer func() {
		client.CacheClient = nil
	}()
	client.CacheClient = testdata.MockFakeCacheClient(t)
	type args struct {
		ctx           context.Context
		projectIDList []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "invalid projectID",
			args: args{
				ctx:           context.Background(),
				projectIDList: []string{"mock123"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitLocalCache(tt.args.ctx, tt.args.projectIDList); (err != nil) != tt.wantErr {
				t.Errorf("InitLocalCache() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewAndSetApplication(t *testing.T) {
	type args struct {
		ctx       context.Context
		projectID string
	}
	tests := []struct {
		name            string
		args            args
		wantApplication *Application
		wantErr         bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotApplication, err := NewAndSetApplication(tt.args.ctx, tt.args.projectID)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAndSetApplication() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotApplication, tt.wantApplication) {
				t.Errorf("NewAndSetApplication() gotApplication = %v, want %v", gotApplication, tt.wantApplication)
			}
		})
	}
}

func TestRelease(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Release()
		})
	}
}

func Test_asyncRefreshLocalCache(t *testing.T) {
	type args struct {
		projectIDList []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asyncRefreshLocalCache(tt.args.projectIDList)
		})
	}
}

func Test_continuousFetch(t *testing.T) {
	type args struct {
		projectID string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			continuousFetch(tt.args.projectID)
		})
	}
}

func Test_genExperimentVersionIndex(t *testing.T) {
	type args struct {
		application *Application
	}
	tests := []struct {
		name string
		args args
		want map[int64]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := genExperimentVersionIndex(tt.args.application); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("genExperimentVersionIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_genGroupVersionIndex(t *testing.T) {
	type args struct {
		application *Application
	}
	tests := []struct {
		name string
		args args
		want map[int64]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := genGroupVersionIndex(tt.args.application); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("genGroupVersionIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLocalCacheWithDefault(t *testing.T) {
	type args struct {
		projectID string
	}
	tests := []struct {
		name string
		args args
		want *Application
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLocalCacheWithDefault(tt.args.projectID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getLocalCacheWithDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNewApplication(t *testing.T) {
	type args struct {
		curApplication *Application
	}
	tests := []struct {
		name string
		args args
		want *Application
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNewApplication(tt.args.curApplication); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNewApplication() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNewBucketInfoIndex(t *testing.T) {
	type args struct {
		curIndex map[int64]*protoctabcacheserver.BucketInfo
	}
	tests := []struct {
		name string
		args args
		want map[int64]*protoctabcacheserver.BucketInfo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNewBucketInfoIndex(tt.args.curIndex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNewBucketInfoIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNewRoaringBitmapIndex(t *testing.T) {
	type args struct {
		curIndex map[int64]*roaring.Bitmap
	}
	tests := []struct {
		name string
		args args
		want map[int64]*roaring.Bitmap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNewRoaringBitmapIndex(tt.args.curIndex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNewRoaringBitmapIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hasTraffic(t *testing.T) {
	type args struct {
		metadata *protoctabcacheserver.DomainMetadata
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasTraffic(tt.args.metadata); got != tt.want {
				t.Errorf("hasTraffic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isFullFlowDomain(t *testing.T) {
	type args struct {
		metadata *protoctabcacheserver.DomainMetadata
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty traffic",
			args: args{metadata: &protoctabcacheserver.DomainMetadata{
				BucketSize:       10000,
				TrafficRangeList: nil,
			}},
			want: false,
		},
		{
			name: "full traffic",
			args: args{metadata: &protoctabcacheserver.DomainMetadata{
				BucketSize: 10000,
				TrafficRangeList: []*protoctabcacheserver.TrafficRange{
					{
						Left:  0,
						Right: 10000,
					},
				},
			}},
			want: true,
		},
		{
			name: "full traffic with multi traffic range",
			args: args{metadata: &protoctabcacheserver.DomainMetadata{
				BucketSize: 10000,
				TrafficRangeList: []*protoctabcacheserver.TrafficRange{
					{
						Left:  5000,
						Right: 10000,
					},
					{
						Left:  0,
						Right: 1000,
					},
					{
						Left:  1000,
						Right: 5000,
					},
				},
			}},
			want: true,
		},
		{
			name: "full traffic with multi traffic range",
			args: args{metadata: &protoctabcacheserver.DomainMetadata{
				BucketSize: 10000,
				TrafficRangeList: []*protoctabcacheserver.TrafficRange{
					{
						Left:  0,
						Right: 1000,
					},
					{
						Left:  1000,
						Right: 5000,
					},
					{
						Left:  5000,
						Right: 10000,
					},
				},
			}},
			want: true,
		},
		{
			name: "not full traffic with multi traffic range",
			args: args{metadata: &protoctabcacheserver.DomainMetadata{
				BucketSize: 10000,
				TrafficRangeList: []*protoctabcacheserver.TrafficRange{
					{
						Left:  0,
						Right: 1000,
					},
					{
						Left:  5000,
						Right: 9000,
					},
					{
						Left:  5000,
						Right: 5000,
					},
				},
			}},
			want: false,
		},
		{
			name: "not full traffic with multi traffic range",
			args: args{metadata: &protoctabcacheserver.DomainMetadata{
				BucketSize: 10000,
				TrafficRangeList: []*protoctabcacheserver.TrafficRange{
					{
						Left:  5000,
						Right: 10000,
					},
					{
						Left:  0,
						Right: 1000,
					},
				},
			}},
			want: false,
		},
		{
			name: "not full traffic",
			args: args{metadata: &protoctabcacheserver.DomainMetadata{
				BucketSize: 10000,
				TrafficRangeList: []*protoctabcacheserver.TrafficRange{
					{
						Left:  0,
						Right: 9000,
					},
				},
			}},
			want: false,
		},
		{
			name: "full traffic",
			args: args{metadata: &protoctabcacheserver.DomainMetadata{
				BucketSize: 10000,
				TrafficRangeList: []*protoctabcacheserver.TrafficRange{
					{
						Left:  -1000,
						Right: 1000000,
					},
				},
			}},
			want: true,
		},
		{
			name: "full traffic",
			args: args{metadata: &protoctabcacheserver.DomainMetadata{
				BucketSize: 10000,
				TrafficRangeList: []*protoctabcacheserver.TrafficRange{
					{
						Left:  1,
						Right: 10000,
					},
				},
			}},
			want: true,
		},
		{
			name: "full traffic",
			args: args{metadata: &protoctabcacheserver.DomainMetadata{
				BucketSize: 10000,
				TrafficRangeList: []*protoctabcacheserver.TrafficRange{
					{
						Left:  1,
						Right: 9000,
					},
				},
			}},
			want: false,
		},
		{
			name: "full traffic",
			args: args{metadata: &protoctabcacheserver.DomainMetadata{
				BucketSize: 10000,
				TrafficRangeList: []*protoctabcacheserver.TrafficRange{
					{
						Left:  1,
						Right: 9999,
					},
				},
			}},
			want: false,
		},
		{
			name: "full traffic",
			args: args{metadata: &protoctabcacheserver.DomainMetadata{
				BucketSize: 10000,
				TrafficRangeList: []*protoctabcacheserver.TrafficRange{
					{
						Left:  2,
						Right: 10000,
					},
				},
			}},
			want: false,
		},
		{
			name: "full traffic",
			args: args{metadata: &protoctabcacheserver.DomainMetadata{
				BucketSize: 11,
				TrafficRangeList: []*protoctabcacheserver.TrafficRange{
					{
						Left:  1,
						Right: 10,
					},
					{
						Left:  7,
						Right: 11,
					},
					{
						Left:  2,
						Right: 5,
					},
				},
			}},
			want: true,
		},
		{
			name: "full traffic",
			args: args{metadata: &protoctabcacheserver.DomainMetadata{
				BucketSize: 11,
				TrafficRangeList: []*protoctabcacheserver.TrafficRange{
					{
						Left:  1,
						Right: 11,
					},
					{
						Left:  2,
						Right: 5,
					},
				},
			}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isFullFlowDomain(tt.args.metadata); got != tt.want {
				t.Errorf("isFullFlowDomain() = %v, want %v, %+v", got, tt.want, tt.args.metadata.TrafficRangeList)
			}
		})
	}
}

func Test_manualFetchEvent(t *testing.T) {
	type args struct {
		projectID string
		latency   time.Duration
		err       error
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manualFetchEvent(tt.args.projectID, tt.args.latency, tt.args.err)
		})
	}
}

func Test_refreshApplication(t *testing.T) {
	type args struct {
		ctx       context.Context
		projectID string
	}
	tests := []struct {
		name    string
		args    args
		want    *Application
		want1   bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := refreshApplication(tt.args.ctx, tt.args.projectID)
			if (err != nil) != tt.wantErr {
				t.Errorf("refreshApplication() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("refreshApplication() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("refreshApplication() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_refreshInterval(t *testing.T) {
	type args struct {
		projectID string
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := refreshInterval(tt.args.projectID); got != tt.want {
				t.Errorf("refreshInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setApplication(t *testing.T) {
	type args struct {
		application *Application
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setApplication(tt.args.application)
		})
	}
}

func Test_setLayerIndex(t *testing.T) {
	type args struct {
		layerIndex map[string]*protoctabcacheserver.Layer
		layerList  []*protoctabcacheserver.Layer
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setLayerIndex(tt.args.layerIndex, tt.args.layerList)
		})
	}
}

func Test_setupDMPTagInfo(t *testing.T) {
	type args struct {
		application *Application
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setupDMPTagInfo(tt.args.application); (err != nil) != tt.wantErr {
				t.Errorf("setupDMPTagInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_setupExperimentBucketInfo(t *testing.T) {
	type args struct {
		ctx         context.Context
		application *Application
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setupExperimentBucketInfo(tt.args.ctx, tt.args.application); (err != nil) != tt.wantErr {
				t.Errorf("setupExperimentBucketInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_setupFillFlowLayerIndexDomain(t *testing.T) {
	type args struct {
		domain *protoctabcacheserver.Domain
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*protoctabcacheserver.Layer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := setupFillFlowLayerIndexDomain(tt.args.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("setupFillFlowLayerIndexDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setupFillFlowLayerIndexDomain() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setupFullFlowLayerIndex(t *testing.T) {
	type args struct {
		application *Application
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setupFullFlowLayerIndex(tt.args.application); (err != nil) != tt.wantErr {
				t.Errorf("setupFullFlowLayerIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_setupGroupBucketInfo(t *testing.T) {
	type args struct {
		ctx         context.Context
		application *Application
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setupGroupBucketInfo(tt.args.ctx, tt.args.application); (err != nil) != tt.wantErr {
				t.Errorf("setupGroupBucketInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_setupLayerIndex(t *testing.T) {
	type args struct {
		application *Application
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setupLayerIndex(tt.args.application); (err != nil) != tt.wantErr {
				t.Errorf("setupLayerIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_setupLayerIndexDomain(t *testing.T) {
	type args struct {
		domain *protoctabcacheserver.Domain
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*protoctabcacheserver.Layer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := setupLayerIndexDomain(tt.args.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("setupLayerIndexDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setupLayerIndexDomain() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setupTabConfig(t *testing.T) {
	type args struct {
		ctx         context.Context
		application *Application
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setupTabConfig(tt.args.ctx, tt.args.application); (err != nil) != tt.wantErr {
				t.Errorf("setupTabConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateTabConfig(t *testing.T) {
	type args struct {
		tabConfigData *protoctabcacheserver.GetTabConfigResp
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateTabConfig(tt.args.tabConfigData); got != tt.want {
				t.Errorf("validateTabConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
