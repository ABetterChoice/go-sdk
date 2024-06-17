// Package client TODO
package client

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/abetterchoice/go-sdk/env"
	"github.com/abetterchoice/go-sdk/internal"
	protocdmpproxyserver "github.com/abetterchoice/protoc_dmp_proxy_server"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

// DMPClient dmp client definition, as an abstract class, shielding the underlying specific implementation
type DMPClient interface {
	// BatchGetDMPTagResult Get the judgment results of unitID dmp tags in batches
	BatchGetDMPTagResult(ctx context.Context, req *protocdmpproxyserver.BatchGetDMPTagResultReq) (
		*protocdmpproxyserver.BatchGetDMPTagResultResp, error)
}

var (
	// DC Abbreviation of dmpClient, the default dmp client
	DC = NewDMPClient()
)

// NewDMPClient Create a new dmp client
func NewDMPClient(opts ...DMPOption) DMPClient {
	c := &tabDMPClient{
		httpClient: &http.Client{ // The specific client timeout can also be reset by ctx
			Transport: transport(env.DMPServerSocket5Addr),
			Timeout:   10 * time.Second, // The default timeout is 10s, and data is refreshed asynchronously,
			// which has little impact on the user. It can be customized through WithHTTPClient
		},
		addr: env.GetDMPAddr(env.TypePrd), // Default formal environment
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// DMPOption Create a new option for dmp client, modify the backend address, request protocol, etc.
type DMPOption func(client *tabDMPClient)

// WithEnvTypeOption Setting up the environment
func WithEnvTypeOption(envType env.Type) DMPOption {
	return func(client *tabDMPClient) {
		client.addr = env.GetDMPAddr(envType)
	}
}

// RegisterDMPClient Register the dmp client. When using dmp to determine whether the dmp tag is hit,
// the underlying dmp client is transparent to the caller.
func RegisterDMPClient(client DMPClient) {
	if client == nil { // 避免传入非法的 client，保证 defaultDMPClient 永远不会 nil
		return
	}
	DC = client
}

// dmp The proxy service is implemented by default and supports the http protocol
type tabDMPClient struct {
	httpClient *http.Client // http client，Request background cache service through http protocol
	addr       string       // http request addr=scheme+host，eg: https://openapi.abetterchoice.ai
}

var (
	// batchGetDMPTagResultURI Interface uri address
	batchGetDMPTagResultURI = "/opensource.tab.dmp_proxy_server.APIServer/BatchGetDMPTagResult"
)

// BatchGetDMPTagResult dmp default implementation, get unit dmp tag judgment result
func (c *tabDMPClient) BatchGetDMPTagResult(ctx context.Context, req *protocdmpproxyserver.BatchGetDMPTagResultReq) (
	*protocdmpproxyserver.BatchGetDMPTagResultResp, error) {
	return c.BatchGetDMPResultHTTP(ctx, req)
}

// BatchGetDMPResultHTTP Request dmp proxy service through http protocol
func (c *tabDMPClient) BatchGetDMPResultHTTP(ctx context.Context, req *protocdmpproxyserver.BatchGetDMPTagResultReq) (
	*protocdmpproxyserver.BatchGetDMPTagResultResp, error) {
	body, err := proto.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "proto marshal")
	}
	httpReq, err := http.NewRequest(http.MethodPost, c.addr+batchGetDMPTagResultURI, bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "http newRequest")
	}
	httpReq = httpReq.WithContext(ctx)
	for key, value := range dmpHeaders {
		httpReq.Header.Set(key, value)
	}
	authHeader(httpReq)
	httpReq.Header.Set(KeyToken, internal.C.SecretKey)
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "http do")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("invalid http status:%s", resp.Status)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "ioutil readAll")
	}
	result := &protocdmpproxyserver.BatchGetDMPTagResultResp{}
	err = proto.Unmarshal(respBody, result)
	if err != nil {
		return nil, errors.Wrap(err, "proto unmarshal")
	}
	return result, nil
}

var dmpHeaders = map[string]string{
	"Content-Type":          "application/proto",
	"X-Tab-Rpc-ServiceName": "opensource.tab.dmp_proxy_server.APIServer",
}

// IsHitDMP TODO
func IsHitDMP(ctx context.Context, req *protocdmpproxyserver.GetDMPTagResultReq, dmpTag string) (bool, error) {
	if req == nil {
		return false, nil
	}
	result, err := DC.BatchGetDMPTagResult(ctx, &protocdmpproxyserver.BatchGetDMPTagResultReq{
		ReqList: []*protocdmpproxyserver.GetDMPTagResultReq{
			req,
		},
	})
	if err != nil {
		return false, errors.Wrap(err, "batchGetDMPTagResult")
	}
	for _, getDMPTagResultResp := range result.RespList {
		statusCode, ok := getDMPTagResultResp.DmpResult[dmpTag]
		if !ok {
			return false, nil
		}
		return statusCode == protocdmpproxyserver.StatusCode_STATUS_CODE_HIT, nil
	}
	return false, nil
}
