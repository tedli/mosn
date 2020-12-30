/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package http2dubbo

import (
	"bytes"
	"context"
	"encoding/json"
	"math/rand"

	transcoder "mosn.io/mosn/pkg/filter/stream/transcoder"

	"github.com/valyala/fasthttp"
	mosnctx "mosn.io/mosn/pkg/context"
	"mosn.io/mosn/pkg/protocol"
	"mosn.io/mosn/pkg/protocol/http"
	"mosn.io/mosn/pkg/protocol/xprotocol/dubbo"
	"mosn.io/mosn/pkg/types"
	"mosn.io/pkg/buffer"
)

func init() {
	transcoder.MustRegister(TRANSCODER_NAME, &http2dubbo{})
}

type http2dubbo struct{}

//accept return when head has transcoder key and value is equal to TRANSCODER_NAME
func (t *http2dubbo) Accept(ctx context.Context, headers types.HeaderMap, buf types.IoBuffer, trailers types.HeaderMap) bool {
	if get, ok := headers.(http.RequestHeader).Get("transcoder"); ok && get == TRANSCODER_NAME {
		return true
	}
	return false
}

// transcode dubbp request to http request
func (t *http2dubbo) TranscodingRequest(ctx context.Context, headers types.HeaderMap, buf types.IoBuffer, trailers types.HeaderMap) (types.HeaderMap, types.IoBuffer, types.HeaderMap, error) {
	// 1. set sub protocol
	mosnctx.WithValue(ctx, types.ContextSubProtocol, string(dubbo.ProtocolName))
	// 2. assemble target request
	targetRequest, err := EncodeHttp2Dubbo(headers, buf)
	if err != nil {
		return nil, nil, nil, err
	}
	return targetRequest.GetHeader(), targetRequest.GetData(), trailers, nil
}

// transcode dubbo response to http response
func (t *http2dubbo) TranscodingResponse(ctx context.Context, headers types.HeaderMap, buf types.IoBuffer, trailers types.HeaderMap) (types.HeaderMap, types.IoBuffer, types.HeaderMap, error) {

	response, err := DecodeDubbo2Http(ctx, headers, buf, trailers)
	if err != nil {
		return nil, nil, nil, err
	}
	return http.ResponseHeader{ResponseHeader: &response.Header}, buffer.NewIoBufferBytes(response.Body()), trailers, nil
}

// decode dubbo response to http response
func DecodeDubbo2Http(ctx context.Context, headers types.HeaderMap, buf types.IoBuffer, trailers types.HeaderMap) (fasthttp.Response, error) {
	sourceResponse := headers.(*dubbo.Frame)
	targetResponse := fasthttp.Response{}

	// 1. headers
	sourceResponse.Range(func(key, Value string) bool {
		if key == "dubbo" || key == "method" || key == "service" {
			targetResponse.Header.Set(key, Value)
		}
		return true
	})
	//body
	dataArr := bytes.Split(sourceResponse.GetData().Bytes(), []byte{10})
	var resType int
	json.Unmarshal(dataArr[0], &resType)
	responseBody := dubbo.DubboHttpResponseBody{}
	exception, value, attachments := []byte{}, []byte{}, []byte{}
	switch resType {
	case dubbo.RESPONSE_WITH_EXCEPTION:
		// error
		exception = dataArr[0]
	case dubbo.RESPONSE_VALUE:
		value = dataArr[1]
	case dubbo.RESPONSE_NULL_VALUE:
	case dubbo.RESPONSE_WITH_EXCEPTION_WITH_ATTACHMENTS:
		exception = dataArr[1]
		value = dataArr[3]
	case dubbo.RESPONSE_VALUE_WITH_ATTACHMENTS:
		value = dataArr[1]
		attachments = dataArr[3]
	case dubbo.RESPONSE_NULL_VALUE_WITH_ATTACHMENTS:
		attachments = dataArr[2]
	}
	json.Unmarshal(value, &responseBody.Value)
	json.Unmarshal(exception, &responseBody.Exception)
	json.Unmarshal(attachments, &responseBody.Attachments)

	body, _ := json.Marshal(responseBody)
	targetResponse.SetBody(body)
	return targetResponse, nil
}

//encode http require to dubbo require
func EncodeHttp2Dubbo(headers types.HeaderMap, buf types.IoBuffer) (*dubbo.Frame, error) {
	//workload
	payLoadByteFin, err := dubbo.EncodeWorkLoad(headers, buf)
	if err != nil {
		return nil, err
	}
	//header
	allHeaders := map[string]string{}
	headers.Range(func(key, value string) bool {
		if key == "dubbo" || key == "method" || key == "service" {
			allHeaders[key] = value
		}
		return true
	})
	// convert data to dubbo frame
	frame := &dubbo.Frame{
		Header: dubbo.Header{
			CommonHeader: protocol.CommonHeader(allHeaders),
		},
	}
	//magic
	frame.Magic = dubbo.MagicTag
	//  flag
	frame.Flag = 0xc6
	// status
	frame.Status = 0
	// decode request id
	frame.Id = rand.Uint64()
	// event
	frame.IsEvent = (frame.Flag & (1 << 5)) != 0
	// twoway
	frame.IsTwoWay = (frame.Flag & (1 << 6)) != 0
	frame.Direction = dubbo.EventRequest
	// serializationId
	frame.SerializationId = int(frame.Flag & 0x1f)

	frame.SetData(buffer.NewIoBufferBytes(payLoadByteFin))

	return frame, nil
}
