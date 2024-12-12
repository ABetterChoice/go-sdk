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

// BatchGetTagValue TODO
func (e emptyDMPClient) BatchGetTagValue(ctx context.Context,
	req *protoc_dmp_proxy_server.BatchGetTagValueReq) (*protoc_dmp_proxy_server.BatchGetTagValueResp, error) {
	resp := &protoc_dmp_proxy_server.BatchGetTagValueResp{
		RetCode:         protoc_dmp_proxy_server.RetCode_RET_CODE_SUCCESS,
		Message:         "mock success",
		UnitId:          req.UnitId,
		UnitType:        req.UnitType,
		DmpPlatformCode: req.DmpPlatformCode,
		TagResult:       map[string]string{},
	}
	for _, tag := range req.TagList {
		if tag == "tagTest123" { // tagTest123 结果 123
			resp.TagResult[tag] = "123"
			continue
		}
		if tag == "dmpCodeTestHit" { // dmpCodeTestHit 结果命中
			resp.TagResult[tag] = "1"
			continue
		}
		resp.TagResult[tag] = "0"
	}
	return resp, nil
}

var (
	// MockEmptyDMPClient TODO
	// Deprecated: for test
	MockEmptyDMPClient = &emptyDMPClient{}
)
