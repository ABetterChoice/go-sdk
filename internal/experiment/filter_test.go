// Package experiment ...
package experiment

import (
	"context"
	"testing"

	protoccacheserver "github.com/abetterchoice/protoc_cache_server"
	"github.com/stretchr/testify/assert"
)

func Test_experimentKeyFilter(t *testing.T) {
	type args struct {
		ctx     context.Context
		layer   *protoccacheserver.Layer
		options *Options
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "normal",
			args: args{
				ctx: context.TODO(),
				layer: &protoccacheserver.Layer{
					Metadata: &protoccacheserver.LayerMetadata{
						Key:          "filterLayer",
						DefaultGroup: nil,
						SceneIdList:  nil,
					},
					ExperimentIndex: map[int64]*protoccacheserver.Experiment{
						123: &protoccacheserver.Experiment{
							Id:  123,
							Key: "123",
						},
					},
				},
				options: &Options{
					SceneIDs:       nil,
					LayerKeys:      nil,
					ExperimentKeys: nil,
				},
			},
			want: true,
		},
		{
			name: "normal",
			args: args{
				ctx: context.TODO(),
				layer: &protoccacheserver.Layer{
					Metadata: &protoccacheserver.LayerMetadata{
						Key:          "filterLayer",
						DefaultGroup: nil,
						SceneIdList:  nil,
					},
					ExperimentIndex: map[int64]*protoccacheserver.Experiment{
						123: &protoccacheserver.Experiment{
							Id:  123,
							Key: "123",
						},
					},
				},
				options: &Options{
					SceneIDs:       nil,
					LayerKeys:      nil,
					ExperimentKeys: map[string]bool{"123": true},
				},
			},
			want: true,
		},
		{
			name: "normal",
			args: args{
				ctx: context.TODO(),
				layer: &protoccacheserver.Layer{
					Metadata: &protoccacheserver.LayerMetadata{
						Key:          "filterLayer",
						DefaultGroup: nil,
						SceneIdList:  nil,
					},
					ExperimentIndex: map[int64]*protoccacheserver.Experiment{
						123: &protoccacheserver.Experiment{
							Id:  123,
							Key: "123",
						},
					},
				},
				options: &Options{
					SceneIDs:       nil,
					LayerKeys:      nil,
					ExperimentKeys: map[string]bool{"1234": true},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, experimentKeyFilter(tt.args.ctx, tt.args.layer, tt.args.options),
				"experimentKeyFilter(%v, %v, %v)", tt.args.ctx, tt.args.layer, tt.args.options)
		})
	}
}

func Test_layerKeyFilter(t *testing.T) {
	type args struct {
		ctx     context.Context
		layer   *protoccacheserver.Layer
		options *Options
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "normal",
			args: args{
				ctx: context.TODO(),
				layer: &protoccacheserver.Layer{
					Metadata: &protoccacheserver.LayerMetadata{
						Key:          "filterLayer",
						DefaultGroup: nil,
						SceneIdList:  nil,
					},
					ExperimentIndex: map[int64]*protoccacheserver.Experiment{
						123: &protoccacheserver.Experiment{
							Id:  123,
							Key: "123",
						},
					},
				},
				options: &Options{
					SceneIDs:       nil,
					LayerKeys:      map[string]bool{"filterLayer": true},
					ExperimentKeys: map[string]bool{"123": true},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, layerKeyFilter(tt.args.ctx, tt.args.layer, tt.args.options),
				"layerKeyFilter(%v, %v, %v)", tt.args.ctx, tt.args.layer, tt.args.options)
		})
	}
}

func Test_sceneIDListFilter(t *testing.T) {
	type args struct {
		ctx     context.Context
		layer   *protoccacheserver.Layer
		options *Options
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "normal",
			args: args{
				ctx: context.TODO(),
				layer: &protoccacheserver.Layer{
					Metadata: &protoccacheserver.LayerMetadata{
						Key:          "filterLayer",
						DefaultGroup: nil,
						SceneIdList:  []int64{1},
					},
					ExperimentIndex: map[int64]*protoccacheserver.Experiment{
						123: &protoccacheserver.Experiment{
							Id:  123,
							Key: "123",
						},
					},
				},
				options: &Options{
					SceneIDs: map[int64]bool{
						1: true,
					},
					LayerKeys:      nil,
					ExperimentKeys: map[string]bool{"123": true},
				},
			},
			want: true,
		},
		{
			name: "normal",
			args: args{
				ctx: context.TODO(),
				layer: &protoccacheserver.Layer{
					Metadata: &protoccacheserver.LayerMetadata{
						Key:          "filterLayer",
						DefaultGroup: nil,
						SceneIdList:  []int64{1},
					},
					ExperimentIndex: map[int64]*protoccacheserver.Experiment{
						123: &protoccacheserver.Experiment{
							Id:  123,
							Key: "123",
						},
					},
				},
				options: &Options{
					SceneIDs: map[int64]bool{
						2: true,
					},
					LayerKeys:      nil,
					ExperimentKeys: map[string]bool{"123": true},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, sceneIDListFilter(tt.args.ctx, tt.args.layer, tt.args.options),
				"sceneIDListFilter(%v, %v, %v)", tt.args.ctx, tt.args.layer, tt.args.options)
		})
	}
}
