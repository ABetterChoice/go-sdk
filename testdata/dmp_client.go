// Package testdata TODO
package testdata

import (
	"context"

	"github.com/abetterchoice/protoc_dmp_proxy_server"
)

type emptyDMPClient struct{}

// BatchGetDMPTagResult TODO
func (e emptyDMPClient) BatchGetDMPTagResult(ctx context.Context,
	req *protoc_dmp_proxy_server.BatchGetDMPTagResultReq) (*protoc_dmp_proxy_server.BatchGetDMPTagResultResp, error) {
	var respList = make([]*protoc_dmp_proxy_server.GetDMPTagResultResp, 0, len(req.ReqList))
	for _, request := range req.ReqList {
		var dmpResult = make(map[string]protoc_dmp_proxy_server.StatusCode, len(request.TagList))
		for _, tagCode := range request.TagList {
			if tagCode == "dmpCodeTestHit" { // dmpCodeTestHit 结果命中
				dmpResult[tagCode] = protoc_dmp_proxy_server.StatusCode_STATUS_CODE_HIT
				continue
			}
			dmpResult[tagCode] = protoc_dmp_proxy_server.StatusCode_STATUS_CODE_MISS
		}
		respList = append(respList, &protoc_dmp_proxy_server.GetDMPTagResultResp{
			RetCode:         protoc_dmp_proxy_server.RetCode_RET_CODE_SUCCESS,
			Message:         "mock success",
			UnitId:          request.UnitId,
			UnitType:        request.UnitType,
			DmpPlatformCode: request.DmpPlatformCode,
			DmpResult:       dmpResult,
		})
	}
	return &protoc_dmp_proxy_server.BatchGetDMPTagResultResp{RespList: respList}, nil
}

var (
	// MockEmptyDMPClient TODO
	// Deprecated: 模拟 dmp client，测试专用
	MockEmptyDMPClient = &emptyDMPClient{}
)
