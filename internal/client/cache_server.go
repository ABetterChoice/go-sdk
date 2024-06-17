// Package client TODO
package client

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/abetterchoice/go-sdk/env"
	"github.com/abetterchoice/go-sdk/internal"
	protoctabcacheserver "github.com/abetterchoice/protoc_cache_server"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

// Client The cache service client abstract class can be used to customize
// the background cache service by implementing client
//
//go:generate mockgen -source=cache_server.go -package=client -destination cache_server_mock.go
type Client interface {
	// GetTabConfigData Get cache data, including experiments, configurations, switches, etc.
	GetTabConfigData(ctx context.Context, req *protoctabcacheserver.GetTabConfigReq) (
		*protoctabcacheserver.GetTabConfigResp, error)
	// BatchGetExperimentBucketInfo Get the experimental bucket information.
	// Only experiments on the double hash type layer have experimental bucket information.
	BatchGetExperimentBucketInfo(ctx context.Context, req *protoctabcacheserver.BatchGetExperimentBucketReq) (
		*protoctabcacheserver.BatchGetExperimentBucketResp, error)
	// BatchGetGroupBucketInfo Get the experimental group bucket information
	BatchGetGroupBucketInfo(ctx context.Context, req *protoctabcacheserver.BatchGetGroupBucketReq) (
		*protoctabcacheserver.BatchGetGroupBucketResp, error)
}

var (
	getTabConfigURI                 = "/opensource.tab.cache_server.APIServer/GetTabConfig"
	batchGetExperimentBucketInfoURI = "/opensource.tab.cache_server.APIServer/BatchGetExperimentBucket"
	batchGetGroupBucketInfoURI      = "/opensource.tab.cache_server.APIServer/BatchGetGroupBucket"
)

const (
	// KeyToken TODO
	KeyToken = "X-Token"
	// KeyAK TODO
	KeyAK = "X-AK"
	// KeyET TODO
	KeyET = "X-ET"
	// KeyES TODO
	KeyES = "X-ES"
)

var (
	// CacheClient Default background cache service client,
	// Init defaults to using tab formal environment background cache service
	CacheClient Client = nil
)

func transport(addr string) http.RoundTripper {
	if len(addr) == 0 {
		return nil
	}
	addrURL, err := url.Parse(addr)
	if err != nil {
		return nil
	}
	return &http.Transport{
		Proxy: http.ProxyURL(addrURL),
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

// NewTABCacheClient create tab server client
func NewTABCacheClient(opts ...Option) Client {
	client := &tabCacheClient{
		httpClient: &http.Client{ // The specific client timeout can also be reset by ctx
			Transport: transport(env.CacheServerSocket5Addr),
			Timeout:   10 * time.Second, // The default timeout is 10s, and data is refreshed asynchronously,
			// which has little impact on the user. It can be customized through WithHTTPClient
		},
		addr: env.GetAddr(env.TypePrd), // Get the official environment address by default
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

// RegisterCacheClient register client
func RegisterCacheClient(client Client) {
	CacheClient = client
}

// Option client option
type Option func(client *tabCacheClient)

// WithHTTPClient Set http client, customize timeout and proxy
func WithHTTPClient(client *http.Client) Option {
	return func(cacheClient *tabCacheClient) {
		cacheClient.httpClient = client
	}
}

// WithEnvType Setting up the environment
func WithEnvType(envType env.Type) Option {
	return func(client *tabCacheClient) {
		client.addr = env.GetAddr(envType)
	}
}

// tabCacheClient Background cache service implementation
type tabCacheClient struct {
	httpClient *http.Client
	addr       string // http request addr=scheme+host，eg: https://openapi.abetterchoice.ai
}

// BatchGetExperimentBucketInfo Get experimental bucket information in batches
func (c *tabCacheClient) BatchGetExperimentBucketInfo(ctx context.Context,
	req *protoctabcacheserver.BatchGetExperimentBucketReq) (*protoctabcacheserver.BatchGetExperimentBucketResp, error) {
	if req == nil || len(req.BucketVersionIndex) == 0 {
		return &protoctabcacheserver.BatchGetExperimentBucketResp{
			Code:        protoctabcacheserver.Code_CODE_SUCCESS,
			Message:     "empty resp",
			BucketIndex: map[int64]*protoctabcacheserver.BucketInfo{},
		}, nil
	}
	body, err := proto.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "proto marshal")
	}
	httpReq, err := http.NewRequest(http.MethodPost, c.addr+batchGetExperimentBucketInfoURI, bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "http newRequest")
	}
	httpReq = httpReq.WithContext(ctx)
	for key, value := range headers {
		httpReq.Header.Set(key, value)
	}
	authHeader(httpReq)
	httpReq.Header.Set(KeyToken, internal.C.SecretKey)
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "http do")
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "ioutil readAll")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("invalid http status:%s, body=%s", resp.Status, respBody)
	}
	result := &protoctabcacheserver.BatchGetExperimentBucketResp{}
	err = proto.Unmarshal(respBody, result)
	if err != nil {
		return nil, errors.Wrap(err, "proto unmarshal")
	}
	return result, nil
}

// BatchGetGroupBucketInfo Get experimental group bucket information in batches
func (c *tabCacheClient) BatchGetGroupBucketInfo(ctx context.Context,
	req *protoctabcacheserver.BatchGetGroupBucketReq) (*protoctabcacheserver.BatchGetGroupBucketResp, error) {
	if req == nil || len(req.BucketVersionIndex) == 0 {
		return nil, errors.Errorf("bucketVersionIndex is required")
	}
	body, err := proto.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "proto marshal")
	}
	httpReq, err := http.NewRequest(http.MethodPost, c.addr+batchGetGroupBucketInfoURI,
		bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "http newRequest")
	}
	httpReq = httpReq.WithContext(ctx)
	for key, value := range headers {
		httpReq.Header.Set(key, value)
	}
	authHeader(httpReq)
	httpReq.Header.Set(KeyToken, internal.C.SecretKey)
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "http do")
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "ioutil readAll")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("invalid http status:%s, body=%s", resp.Status, respBody)
	}
	result := &protoctabcacheserver.BatchGetGroupBucketResp{}
	err = proto.Unmarshal(respBody, result)
	if err != nil {
		return nil, errors.Wrap(err, "proto unmarshal")
	}
	return result, nil
}

// GetTabConfigData Get cache data
func (c *tabCacheClient) GetTabConfigData(ctx context.Context, req *protoctabcacheserver.GetTabConfigReq) (
	*protoctabcacheserver.GetTabConfigResp, error) {
	body, err := proto.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "proto marshal")
	}
	httpReq, err := http.NewRequest(http.MethodPost, c.addr+getTabConfigURI,
		bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "http newRequest")
	}
	httpReq = httpReq.WithContext(ctx)
	for key, value := range headers {
		httpReq.Header.Set(key, value)
	}
	authHeader(httpReq)
	httpReq.Header.Set(KeyToken, internal.C.SecretKey)
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "http do")
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "ioutil readAll")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("invalid http status:%s, body=%s", resp.Status, respBody)
	}
	result := &protoctabcacheserver.GetTabConfigResp{}
	err = proto.Unmarshal(respBody, result)
	if err != nil {
		return nil, errors.Wrap(err, "proto unmarshal")
	}
	return result, nil
}

// Request a fixed header for the tab background
var headers = map[string]string{
	"Content-Type":          "application/proto",
	"X-Tab-Rpc-ServiceName": "opensource.tab.cache_server.APIServer",
}

func mustGetAK(secretKey string) string {
	ak, _ := getAK(secretKey)
	return ak
}

func getAK(secretKey string) (string, error) {
	pairs := strings.Split(secretKey, ".")
	if len(pairs) != 3 {
		return "", errors.Errorf("invalid secretKey=%s", secretKey)
	}
	payload, err := base64.RawStdEncoding.DecodeString(pairs[1])
	if err != nil {
		return "", errors.Wrapf(err, "decodeString")
	}
	type data struct {
		TokenName string `json:"tokenName"`
	}
	d := &data{}
	err = json.Unmarshal(payload, d)
	return d.TokenName, errors.Wrapf(err, "json unmarshal")
}

func genSign(secretKey, ak, timestamp string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(secretKey+ak+timestamp)))
}

func authHeader(req *http.Request) {
	ak := mustGetAK(internal.C.SecretKey)
	now := strconv.FormatInt(time.Now().Unix(), 10)
	req.Header.Set(KeyAK, ak)
	req.Header.Set(KeyET, now)
	req.Header.Set(KeyES, genSign(internal.C.SecretKey, ak, now))
}
