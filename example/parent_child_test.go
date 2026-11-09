package example

import (
	"context"
	"fmt"
	"os"
	"sort"
	"testing"

	sdk "github.com/abetterchoice/go-sdk"
	"github.com/abetterchoice/go-sdk/env"
	"github.com/abetterchoice/go-sdk/internal/cache"
	protoccacheserver "github.com/abetterchoice/protoc_cache_server"
)

// 拉测试环境真实下发，验证父子实验分流。
//
//	ABC_PROJECT_ID=45985 \
//	ABC_SECRET_KEY=... \
//	ABC_PARENT_LAYER=test_child_1 \
//	go test -run TestParentChildLiveConfig -v ./example
func TestParentChildLiveConfig(t *testing.T) {
	projectID := os.Getenv("ABC_PROJECT_ID")
	secretKey := os.Getenv("ABC_SECRET_KEY")
	parentLayerKey := os.Getenv("ABC_PARENT_LAYER")
	if projectID == "" || secretKey == "" || parentLayerKey == "" {
		t.Skip("set ABC_PROJECT_ID, ABC_SECRET_KEY and ABC_PARENT_LAYER")
	}

	defer sdk.Release()
	if err := sdk.Init(context.TODO(), []string{projectID},
		sdk.WithEnvType(env.TypeTest),
		sdk.WithSecretKey(secretKey),
		sdk.WithDisableReport(true),
	); err != nil {
		t.Fatalf("Init fail: %v", err)
	}

	application := cache.GetApplication(projectID)
	if application == nil {
		t.Fatalf("application not found")
	}
	parentLayer := application.LayerIndex[parentLayerKey]
	if parentLayer == nil {
		t.Fatalf("parent layer %s not delivered", parentLayerKey)
	}
	dumpLayer(t, "父层", parentLayer)

	// 找所有 LIVE_PARENT 源版本落在父层上的子层
	type childRoute struct {
		layer         *protoccacheserver.Layer
		sourceGroupID int64
	}
	var children []childRoute
	for _, layer := range application.LayerIndex {
		for _, routing := range layer.TrafficRouting {
			if routing.GetRoutingMode() !=
				protoccacheserver.ExpGroupTrafficRoutingMode_EXP_GROUP_TRAFFIC_ROUTING_MODE_LIVE_PARENT {
				continue
			}
			for _, route := range routing.GetRoutes() {
				if _, ok := parentLayer.GroupIndex[route.GetSourceExperimentId()]; ok {
					children = append(children, childRoute{layer, route.GetSourceExperimentId()})
				}
			}
		}
	}
	if len(children) == 0 {
		t.Fatalf("no LIVE_PARENT child layer references parent layer %s", parentLayerKey)
	}
	for _, c := range children {
		dumpLayer(t, "子层", c.layer)
		t.Logf("索引 LiveParentSourceLayerIndex[%d] -> %v",
			c.sourceGroupID, layerKeyOf(application.LiveParentSourceLayerIndex[c.sourceGroupID]))
	}

	const userNum = 2000
	for _, c := range children {
		childLayerKey := c.layer.GetMetadata().GetKey()
		parentDist := map[string]int{}
		childDist := map[string]int{}
		violations := 0
		for i := 0; i < userNum; i++ {
			uc := sdk.NewUserContext(fmt.Sprintf("parent-child-user-%d", i))
			parentRes, err := uc.GetExperiment(context.TODO(), projectID, parentLayerKey)
			if err != nil {
				t.Fatalf("GetExperiment parent fail: %v", err)
			}
			// 单独查子层：父层不在本次查询范围内，影子评估必须仍能判定父命中
			childRes, err := uc.GetExperiment(context.TODO(), projectID, childLayerKey)
			if err != nil {
				t.Fatalf("GetExperiment child fail: %v", err)
			}
			parentGroup := groupOf(parentRes)
			childGroup := groupOf(childRes)
			parentDist[resultLabel(parentGroup)]++
			childDist[resultLabel(childGroup)]++

			hitSource := parentGroup != nil && parentGroup.ID == c.sourceGroupID
			inChild := childGroup != nil && !childGroup.IsDefault
			if inChild != hitSource {
				violations++
				if violations <= 5 {
					t.Errorf("user %d: parent=%s child=%s，进子实验与命中父版本 %d 不一致",
						i, resultLabel(parentGroup), resultLabel(childGroup), c.sourceGroupID)
				}
			}
		}
		t.Logf("子层 %s（源版本 %d）%d 个用户：", childLayerKey, c.sourceGroupID, userNum)
		t.Logf("  父层分布 %s", formatDist(parentDist))
		t.Logf("  子层分布 %s", formatDist(childDist))
		t.Logf("  不一致 %d 个", violations)
	}
}

func groupOf(res *sdk.ExperimentResult) *sdk.Group {
	if res == nil {
		return nil
	}
	return res.Group
}

func resultLabel(g *sdk.Group) string {
	if g == nil {
		return "<nil>"
	}
	if g.IsDefault {
		return fmt.Sprintf("default(%d)", g.ID)
	}
	return fmt.Sprintf("%s(%d)", g.Key, g.ID)
}

func groupLabel(g *protoccacheserver.Group) string {
	if g == nil {
		return "<nil>"
	}
	if g.GetIsDefault() {
		return fmt.Sprintf("default(%d)", g.GetId())
	}
	return fmt.Sprintf("%s(%d)", g.GetGroupKey(), g.GetId())
}

func layerKeyOf(layer *protoccacheserver.Layer) string {
	if layer == nil {
		return "<nil>"
	}
	return layer.GetMetadata().GetKey()
}

func formatDist(dist map[string]int) string {
	keys := make([]string, 0, len(dist))
	for k := range dist {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	s := ""
	for _, k := range keys {
		s += fmt.Sprintf("%s=%d ", k, dist[k])
	}
	return s
}

func dumpLayer(t *testing.T, title string, layer *protoccacheserver.Layer) {
	md := layer.GetMetadata()
	t.Logf("==== %s %s hashType=%v bucketSize=%d defaultGroup=%v defaultGroupParams=%v",
		title, md.GetKey(), md.GetHashType(), md.GetBucketSize(), groupLabel(md.GetDefaultGroup()),
		md.GetDefaultGroup().GetParams())
	for id, g := range layer.GroupIndex {
		t.Logf("  group %d key=%s exp=%d default=%v params=%v", id, g.GetGroupKey(), g.GetExperimentId(),
			g.GetIsDefault(), g.GetParams())
	}
	for id, e := range layer.ExperimentIndex {
		ids := make([]int64, 0, len(e.GetGroupIdIndex()))
		for gid := range e.GetGroupIdIndex() {
			ids = append(ids, gid)
		}
		t.Logf("  experiment %d key=%s issueType=%v groups=%v", id, e.GetKey(), e.GetIssueType(), ids)
	}
	for _, r := range layer.TrafficRouting {
		t.Logf("  routing mode=%v source_exp_group=%d target_exp_group=%d", r.GetRoutingMode(),
			r.GetSourceExpGroupId(), r.GetTargetExpGroupId())
		for _, route := range r.GetRoutes() {
			t.Logf("    route source_experiment_id=%d target_experiment_ids=%v allocation=%v",
				route.GetSourceExperimentId(), route.GetTargetExperimentIds(), route.GetAllocationMode())
		}
	}
}
