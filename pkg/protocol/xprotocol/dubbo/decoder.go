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

package dubbo

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"
	"runtime/debug"
	"sync"

	"mosn.io/mosn/pkg/trace"

	hessian "github.com/apache/dubbo-go-hessian2"
	"mosn.io/mosn/pkg/protocol"
	"mosn.io/mosn/pkg/types"
	"mosn.io/pkg/buffer"
)

// Decoder is heavy and caches to improve performance.
// Avoid allocating 4k memory every time you create an object
var (
	decodePoolCheap = &sync.Pool{
		New: func() interface{} {
			return hessian.NewCheapDecoderWithSkip([]byte{})
		},
	}
	decodePool = &sync.Pool{
		New: func() interface{} {
			return hessian.NewDecoderWithSkip([]byte{})
		},
	}
)

func decodeFrame(ctx context.Context, data types.IoBuffer) (cmd interface{}, err error) {
	// convert data to dubbo frame
	dataBytes := data.Bytes()
	frame := &Frame{
		Header: Header{
			CommonHeader: protocol.CommonHeader{},
		},
	}
	// decode magic
	frame.Magic = dataBytes[MagicIdx:FlagIdx]
	// decode flag
	frame.Flag = dataBytes[FlagIdx]
	// decode status
	frame.Status = dataBytes[StatusIdx]
	// decode request id
	reqIDRaw := dataBytes[IdIdx:(IdIdx + IdLen)]
	frame.Id = binary.BigEndian.Uint64(reqIDRaw)
	// decode data length
	frame.DataLen = binary.BigEndian.Uint32(dataBytes[DataLenIdx:(DataLenIdx + DataLenSize)])

	// decode event
	frame.IsEvent = (frame.Flag & (1 << 5)) != 0

	// decode twoway
	frame.IsTwoWay = (frame.Flag & (1 << 6)) != 0

	// decode direction
	directionBool := frame.Flag & (1 << 7)
	if directionBool != 0 {
		frame.Direction = EventRequest
	} else {
		frame.Direction = EventResponse
	}
	// decode serializationId
	frame.SerializationId = int(frame.Flag & 0x1f)

	frameLen := HeaderLen + frame.DataLen
	// decode payload
	body := make([]byte, frameLen)
	copy(body, dataBytes[:frameLen])
	frame.payload = body[HeaderLen:]
	frame.content = buffer.NewIoBufferBytes(frame.payload)

	// not heartbeat & is request
	if !frame.IsEvent && frame.Direction == EventRequest {
		// service aware
		meta, err := getServiceAwareMeta(ctx, frame)
		if err != nil {
			return nil, err
		}
		for k, v := range meta {
			frame.Set(k, v)
		}
	}

	frame.rawData = body
	frame.data = buffer.NewIoBufferBytes(frame.rawData)
	data.Drain(int(frameLen))
	return frame, nil
}

func getServiceAwareMeta(ctx context.Context, frame *Frame) (map[string]string, error) {
	meta := make(map[string]string, 8)
	switch frame.SerializationId {
	//dubbo encode by Hessian
	case 2:
		m, err2 := decodeHessian(ctx, frame, meta)
		if err2 != nil {
			return m, err2
		}
	// dubbo decode by fastson
	case 6:
		m, err2 := decodeFastjosn(ctx, frame, meta)
		if err2 != nil {
			return m, err2
		}
	default:
		return meta, fmt.Errorf("[xprotocol][dubbo] type do not support")
	}
	return meta, nil
}
func decodeFastjosn(ctx context.Context, frame *Frame, meta map[string]string) (map[string]string, error) {
	var (
		err              error
		frameworkVersion string
		path             string
		version          string
		method           string
		paramsTypes      string
		attachmentsMap   map[string]string
	)
	arr := bytes.Split(frame.payload[:], []byte{10})
	err = json.Unmarshal(arr[0], &frameworkVersion)
	if err != nil {
		return meta, fmt.Errorf("[xprotocol][dubbo] fastjson decode framework version fail")
	}
	meta[FrameworkVersionNameHeader] = frameworkVersion

	err = json.Unmarshal(arr[1], &path)
	if err != nil {
		return meta, fmt.Errorf("[xprotocol][dubbo] fastjson decode service path fail")
	}
	meta[ServiceNameHeader] = path

	// get version name
	err = json.Unmarshal(arr[2], &version)
	if err != nil {
		return nil, fmt.Errorf("[xprotocol][dubbo] fastjson decode method version fail")
	}
	meta[VersionNameHeader] = version
	//method
	err = json.Unmarshal(arr[3], &method)
	if err != nil {
		return nil, fmt.Errorf("[xprotocol][dubbo] fastjson decode method fail")
	}
	meta[MethodNameHeader] = method
	//params
	err = json.Unmarshal(arr[4], &paramsTypes)
	if err != nil {
		return nil, fmt.Errorf("[xprotocol][dubbo] fastjson decode paramsTypes fail")
	}

	if ctx != nil {
		listener := ctx.Value(types.ContextKeyListenerName)

		var (
			node    *Node
			matched bool
		)

		// for better performance.
		// If the ingress scenario is not using group,
		// we can skip parsing attachment to improve performance
		if listener == IngressDubbo {
			if node, matched = DubboPubMetadata.Find(path, version); matched {
				meta[ServiceNameHeader] = node.Service
				meta[GroupNameHeader] = node.Group
			}
		} else if listener == EgressDubbo {
			// for better performance.
			// If the egress scenario is not using group,
			// we can skip parsing attachment to improve performance
			if node, matched = DubboSubMetadata.Find(path, version); matched {
				meta[ServiceNameHeader] = node.Service
				meta[GroupNameHeader] = node.Group
			}
		}

		// decode the attachment to get the real service and group parameters
		if !matched && (listener == EgressDubbo || listener == IngressDubbo) || trace.IsEnabled() {

			count := GetArgumentCount(paramsTypes)
			attachments := arr[5+count]
			err = json.Unmarshal(attachments, &attachmentsMap)
			if err != nil {
				return nil, fmt.Errorf("[xprotocol][dubbo] fastjosn decode dubbo attachments error, %v", err)
			}
			// we loop all attachments and check element type,
			// we should only read string types.
			for k, v := range attachmentsMap {
				meta[k] = v
				// we should use interface value,
				// convenient for us to do service discovery.
				if k == InterfaceNameHeader {
					meta[ServiceNameHeader] = v
				}
			}
		}
	}

	return meta, nil
}

func DecodeParams(paramsTypes string, i [][]byte) ([]Parameter, error) {
	params := make([]Parameter, 0, len(i))
	types := getArguments(paramsTypes)
	if len(types) == 0 {
		return params, nil
	}
	for key, val := range types {
		pa := Parameter{}
		pa.Type = val
		err := json.Unmarshal(i[key], &pa.Value)
		if err != nil {
			return nil, err
		}
		params = append(params, pa)
	}
	return params, nil
}

func decodeHessian(ctx context.Context, frame *Frame, meta map[string]string) (map[string]string, error) {

	// Recycle decode
	var (
		decoder  *hessian.Decoder
		listener interface{}
	)

	if ctx != nil {
		listener = ctx.Value(types.ContextKeyListenerName)
	}
	if listener == IngressDubbo || listener == EgressDubbo {
		decoder = decodePool.Get().(*hessian.Decoder)
		defer decodePool.Put(decoder)
	} else {
		decoder = decodePoolCheap.Get().(*hessian.Decoder)
		defer decodePoolCheap.Put(decoder)
	}
	decoder.Reset(frame.payload[:])

	var (
		field            interface{}
		err              error
		ok               bool
		frameworkVersion string
		path             string
		version          string
		method           string
	)

	// framework version + path + version + method
	// get service name
	field, err = decoder.Decode()
	if err != nil {
		return meta, fmt.Errorf("[xprotocol][dubbo] hessiandecode framework version fail: %v", err)
	}
	frameworkVersion, ok = field.(string)
	if !ok {
		return meta, fmt.Errorf("[xprotocol][dubbo] hessiandecode framework version {%v} type error", field)
	}
	meta[FrameworkVersionNameHeader] = frameworkVersion

	field, err = decoder.Decode()
	if err != nil {
		return meta, fmt.Errorf("[xprotocol][dubbo] decode service path fail: %v", err)
	}
	path, ok = field.(string)
	if !ok {
		return meta, fmt.Errorf("[xprotocol][dubbo] service path {%v} type error", field)
	}
	meta[ServiceNameHeader] = path

	// get method name
	field, err = decoder.Decode()
	if err != nil {
		return nil, fmt.Errorf("[xprotocol][dubbo] decode method version fail: %v", err)
	}
	// callback maybe return nil
	if field != nil {
		version, ok = field.(string)
		if !ok {
			return nil, fmt.Errorf("[xprotocol][dubbo] method version {%v} type fail", field)
		}
	}
	meta[VersionNameHeader] = version

	field, err = decoder.Decode()
	if err != nil {
		return nil, fmt.Errorf("[xprotocol][dubbo] decode method fail: %v", err)
	}
	method, ok = field.(string)
	if !ok {
		return nil, fmt.Errorf("[xprotocol][dubbo] method {%v} type error", field)
	}
	meta[MethodNameHeader] = method

	if ctx != nil {
		var (
			node    *Node
			matched bool
		)

		// for better performance.
		// If the ingress scenario is not using group,
		// we can skip parsing attachment to improve performance
		if listener == IngressDubbo {
			if node, matched = DubboPubMetadata.Find(path, version); matched {
				meta[ServiceNameHeader] = node.Service
				meta[GroupNameHeader] = node.Group
			}
		} else if listener == EgressDubbo {
			// for better performance.
			// If the egress scenario is not using group,
			// we can skip parsing attachment to improve performance
			if node, matched = DubboSubMetadata.Find(path, version); matched {
				meta[ServiceNameHeader] = node.Service
				meta[GroupNameHeader] = node.Group
			}
		}

		// decode the attachment to get the real service and group parameters
		if !matched && (listener == EgressDubbo || listener == IngressDubbo) || trace.IsEnabled() {

			// decode arguments maybe panic, when dubbo payload have complex struct
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("decode arguments error :%v\n%s", r, debug.Stack())
				}
			}()

			field, err = decoder.Decode()
			if err != nil {
				return nil, fmt.Errorf("[xprotocol][dubbo] decode dubbo argument types error: %v", err)
			}

			arguments := GetArgumentCount(field.(string))
			// we must skip all method arguments.
			for i := 0; i < arguments; i++ {
				_, err = decoder.Decode()
				if err != nil {
					return nil, fmt.Errorf("[xprotocol][dubbo] decode dubbo argument error: %v", err)
				}
			}

			field, err = decoder.Decode()
			if err != nil {
				return nil, fmt.Errorf("[xprotocol][dubbo] decode dubbo attachments error: %v", err)
			}

			if field != nil {
				if origin, ok := field.(map[interface{}]interface{}); ok {
					// we loop all attachments and check element type,
					// we should only read string types.
					for k, v := range origin {
						if key, ok := k.(string); ok {
							if val, ok := v.(string); ok {
								meta[key] = val
								// we should use interface value,
								// convenient for us to do service discovery.
								if key == InterfaceNameHeader {
									meta[ServiceNameHeader] = val
								}
							}
						}
					}
				}
			}
		}
	}

	return meta, nil
}

//  more unit test:
// https://github.com/zonghaishang/dubbo/commit/e0fd702825a274379fb609229bdb06ca0586122e
func GetArgumentCount(desc string) int {
	lens := len(desc)
	if lens == 0 {
		return 0
	}

	var args, next = 0, false
	for _, ch := range desc {

		// is array ?
		if ch == '[' {
			continue
		}

		// is object ?
		if next && ch != ';' {
			continue
		}

		switch ch {
		case 'V', // void
			'Z', // boolean
			'B', // byte
			'C', // char
			'D', // double
			'F', // float
			'I', // int
			'J', // long
			'S': // short
			args++
		default:
			// we found object
			if ch == 'L' {
				args++
				next = true
				// end of object ?
			} else if ch == ';' {
				next = false
			}
		}

	}
	return args
}

// get types arr
//desc example: Ljava/lang/String;ILtest/bean/TestBean;[I"
//return example ["java/lang/String","int","test/bean/TestBean","[int"]
func getArguments(desc string) []string {
	typesArr := []string{}
	lens := len(desc)
	if lens == 0 {
		return typesArr
	}

	var next, tmp = false, ""
	for _, ch := range desc {
		// is array ?
		if ch == '[' {
			tmp += "["
			continue
		}

		// is object ?
		if next && ch != ';' {
			if ch == '/' {
				tmp += string('.')
			} else {
				tmp += string(ch)
			}
			continue
		}

		switch ch {
		case 'V': // void
			tmp += "void"
			typesArr = append(typesArr, tmp)
			tmp = ""
		case 'Z': // boolean
			tmp += "boolean"
			typesArr = append(typesArr, tmp)
			tmp = ""
		case 'B': // byte
			tmp += "byte"
			typesArr = append(typesArr, tmp)
			tmp = ""
		case 'C': // char
			tmp += "char"
			typesArr = append(typesArr, tmp)
			tmp = ""
		case 'D': // double
			tmp += "double"
			typesArr = append(typesArr, tmp)
			tmp = ""
		case 'F': // float
			tmp += "float"
			typesArr = append(typesArr, tmp)
			tmp = ""
		case 'I': // int
			tmp += "int"
			typesArr = append(typesArr, tmp)
			tmp = ""
		case 'J': // long
			tmp += "long"
			typesArr = append(typesArr, tmp)
			tmp = ""
		case 'S': // short
			tmp += "short"
			typesArr = append(typesArr, tmp)
			tmp = ""
		default:
			// we found object
			if ch == 'L' {
				next = true
				// end of object ?
			} else if ch == ';' {
				next = false
				typesArr = append(typesArr, tmp)
				tmp = ""
			}
		}

	}
	return typesArr
}

//decode the simple type to dubbo transmission type
// exmple: int -> I
func EncodeRequestType(tp string) string {
	var res string
	if strings.HasPrefix(tp, "[") {
		res += "["
		tp = tp[1:]
	}
	switch tp {
	case "void":
		res += "V"
	case "boolean":
		res += "Z"
	case "byte":
		res += "B"
	case "char":
		res += "C"
	case "double":
		res += "D"
	case "float":
		res += "F"
	case "int":
		res += "I"
	case "long":
		res += "J"
	case "short":
		res += "S"
	default:
		res += "L" + strings.ReplaceAll(tp, ".", "/") + ";"
	}
	return res
}

// encode dubbo.freame to fastjson dubbo workload
// return []byte set in fream.workload to transmit over the network
func EncodeWorkLoad(headers types.HeaderMap, buf types.IoBuffer) ([]byte, error) {
	//body
	var reqBody DubboHttpRequestParams
	if buf == nil {
		return nil, fmt.Errorf("nil buf data error")
	}
	if err := json.Unmarshal(buf.Bytes(), &reqBody); err != nil {
		return nil, err
	}
	if reqBody.Attachments == nil {
		reqBody.Attachments = make(map[string]string)
	}

	//设置playload
	headers.Range(func(key, value string) bool {
		reqBody.Attachments[key] = value
		return true
	})

	//service
	serviceName := HeadGetDefault(headers, "service", "")
	reqBody.Attachments["interface"] = serviceName

	dubboVersion := HeadGetDefault(headers, "dubbo", "2.6.5")
	serviceVersion := HeadGetDefault(headers, "version", "0.0.0")
	reqBody.Attachments["version"] = serviceVersion

	serviceMethod := HeadGetDefault(headers, "method", "")
	serviceGroup := HeadGetDefault(headers, "group", "")
	if serviceGroup != "" {
		reqBody.Attachments["group"] = serviceGroup
	}

	dubboVersionByte := []byte(`"` + dubboVersion + `"`)
	serviceNameByte := []byte(`"` + serviceName + `"`)
	verionByte := []byte(`"` + serviceVersion + `"`)
	methodNameByte := []byte(`"` + serviceMethod + `"`)

	//有几个类型写一个类型，直接跟在字符串后面，数组的前面加[，基本类型要转义,不换行
	//有几个参数写几个参数，要写在byte[]后面,换行
	paramesByte := []byte{}
	var paramesTypeStr string
	for i := 0; i < len(reqBody.Parameters); i++ {
		paramesTypeStr += EncodeRequestType(reqBody.Parameters[i].Type)
		valByte, _ := json.Marshal(reqBody.Parameters[i].Value)
		paramesByte = append(paramesByte, valByte...)
		if i < len(reqBody.Parameters)-1 {
			paramesByte = append(paramesByte, []byte{10}...)
		}
	}

	paramesTypeByte, _ := json.Marshal(paramesTypeStr)
	attachmentsByte, _ := json.Marshal(reqBody.Attachments)

	payLoadByteFin := bytes.Join([][]byte{dubboVersionByte, serviceNameByte, verionByte, methodNameByte, paramesTypeByte, paramesByte, attachmentsByte}, []byte{10})

	return payLoadByteFin, nil
}

func HeadGetDefault(headers types.HeaderMap, key string, defaultValue string) string {
	if value, ok := headers.Get(key); !ok {
		return defaultValue
	} else {
		return value
	}
}
