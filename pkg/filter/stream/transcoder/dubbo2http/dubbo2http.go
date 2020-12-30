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

package dubbo2http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

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
	transcoder.MustRegister(TRANSCODER_NAME, &dubbo2http{})
}

type dubbo2http struct{}

//accept return when head has transcoder key and value is equal to TRANSCODER_NAME
func (t *dubbo2http) Accept(ctx context.Context, headers types.HeaderMap, buf types.IoBuffer, trailers types.HeaderMap) bool {
	if get, ok := headers.(http.RequestHeader).Get("transcoder"); ok && get == TRANSCODER_NAME {
		return true
	}
	return false
}

//dubbo request 2 http request
func (t *dubbo2http) TranscodingRequest(ctx context.Context, headers types.HeaderMap, buf types.IoBuffer, trailers types.HeaderMap) (types.HeaderMap, types.IoBuffer, types.HeaderMap, error) {
	// 1. set sub protocol
	mosnctx.WithValue(ctx, types.ContextSubProtocol, string(protocol.HTTP1))
	sourceHeader := headers.(*dubbo.Frame)
	//// 2. assemble target request
	byteData, err := DeocdeWorkLoad(headers, buf)
	if err != nil {
		return nil, nil, nil, err
	}
	reqHeaderImpl := &fasthttp.RequestHeader{}
	sourceHeader.Header.CommonHeader.Range(func(key, value string) bool {
		if key != fasthttp.HeaderContentLength {
			reqHeaderImpl.SetCanonical([]byte(key), []byte(value))
		}
		return true
	})
	//set request id
	reqHeaderImpl.Set(dubbo.HTTP_DUBBO_REQUEST_ID_NAME, strconv.FormatUint(sourceHeader.Id, 10))
	reqHeaders := http.RequestHeader{reqHeaderImpl, nil}
	return reqHeaders, buffer.NewIoBufferBytes(byteData), nil, nil
}

// encode the dubbo request body 2 http request body
func DeocdeWorkLoad(headers types.HeaderMap, buf types.IoBuffer) ([]byte, error) {
	var paramsTypes string
	sourceRequest := headers.(*dubbo.Frame)
	dataArr := bytes.Split(sourceRequest.GetData().Bytes(), []byte{10})
	err := json.Unmarshal(dataArr[4], &paramsTypes)
	if err != nil {
		return nil, fmt.Errorf("[xprotocol][dubbo] decode params fail")
	}
	count := dubbo.GetArgumentCount(paramsTypes)
	//skip useless dubbo flags
	arrs := dataArr[5 : 5+count]
	// decode dubbo body
	params, err := dubbo.DecodeParams(paramsTypes, arrs)
	if err != nil {
		return nil, fmt.Errorf("[xprotocol][dubbo] decode params fail")
	}
	attachments := map[string]string{}
	err = json.Unmarshal(dataArr[5+count], &attachments)
	if err != nil {
		return nil, fmt.Errorf("[xprotocol][dubbo] decode attachments fail")
	}
	//encode to http budy
	content := map[string]interface{}{}
	content["attachments"] = attachments
	content["parameters"] = params
	byte, _ := json.Marshal(content)
	return byte, nil
}

//http2dubbo
func (t *dubbo2http) TranscodingResponse(ctx context.Context, headers types.HeaderMap, buf types.IoBuffer, trailers types.HeaderMap) (types.HeaderMap, types.IoBuffer, types.HeaderMap, error) {
	targetRequest, err := DecodeHttp2Dubbo(headers, buf)
	if err != nil {
		return nil, nil, nil, err
	}
	return targetRequest.GetHeader(), targetRequest.GetData(), trailers, nil
}

// decode http response to dubbo response
func DecodeHttp2Dubbo(headers types.HeaderMap, buf types.IoBuffer) (*dubbo.Frame, error) {

	sourceHeader := headers.(http.ResponseHeader)
	//header
	allHeaders := map[string]string{}
	sourceHeader.Range(func(key, value string) bool {
		//skip for Content-Length,the Content-Length may effect the value decode when transcode more one time
		if key != "Content-Length" && key != "Accept:" {
			allHeaders[key] = value
		}
		return true
	})
	frame := &dubbo.Frame{
		Header: dubbo.Header{
			CommonHeader: protocol.CommonHeader(allHeaders),
		},
	}
	// convert data to dubbo frame
	workLoad, err := EncodeWorkLoad(headers, buf)
	if err != nil {
		return nil, err
	}

	//magic
	frame.Magic = dubbo.MagicTag

	//  flag
	frame.Flag = 0x46
	// status when http return not ok, return error
	if sourceHeader.StatusCode() != http.OK {
		//BAD_RESPONSE
		frame.Status = 0x32
	}
	// decode request id
	if id, ok := allHeaders[dubbo.HTTP_DUBBO_REQUEST_ID_NAME]; !ok {
		return nil, fmt.Errorf("[xprotocol][dubbo] decode dubbo id missed")
	} else {
		frameId, _ := strconv.ParseInt(id, 10, 64)
		frame.Id = uint64(frameId)
	}

	// event
	frame.IsEvent = false
	// twoway
	frame.IsTwoWay = true
	// direction
	frame.Direction = dubbo.EventResponse
	// serializationId json
	frame.SerializationId = 6
	frame.SetData(buffer.NewIoBufferBytes(workLoad))
	return frame, nil
}

// http response body example: {"attachments":null,"value":"22222","exception":null}
// make dubbo workload
func EncodeWorkLoad(headers types.HeaderMap, buf types.IoBuffer) ([]byte, error) {
	responseBody := dubbo.DubboHttpResponseBody{}
	workload := [][]byte{}
	if buf == nil {
		return nil, fmt.Errorf("no workload to decode")
	}
	err := json.Unmarshal(buf.Bytes(), &responseBody)
	if err != nil {
		return nil, err
	}
	if responseBody.Exception == "" {
		//out.writeByte(
		if responseBody.Value == nil {
			resType, _ := json.Marshal(dubbo.RESPONSE_NULL_VALUE)
			workload = append(workload, resType)
		} else {
			resType, _ := json.Marshal(dubbo.RESPONSE_VALUE)
			workload = append(workload, resType)
			ret, _ := json.Marshal(responseBody.Value)
			workload = append(workload, ret)

		}
	} else {
		resType, _ := json.Marshal(dubbo.RESPONSE_WITH_EXCEPTION)
		workload = append(workload, resType)
		ret, _ := json.Marshal(responseBody.Exception)
		workload = append(workload, ret)
	}
	workloadByte := bytes.Join(workload, []byte{10})

	return workloadByte, nil
}
