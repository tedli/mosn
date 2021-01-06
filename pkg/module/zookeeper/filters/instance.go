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

package filters

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"mosn.io/mosn/pkg/module/dubbo"
	"mosn.io/mosn/pkg/module/zookeeper"
	"mosn.io/pkg/log"
)

func NewInstanceFilter() zookeeper.Filter {
	return new(instance)
}

type instance struct{}

type getDataMark interface {
	IsGetInstance() bool
	GetApplicationName() string
	GetIP() string
	GetPort() int
}

type getDataMarkImpl struct {
	applicationName string
	ip              string
	port            int
	endpoint        *dubbo.ProviderEndpoint
}

func (gdm getDataMarkImpl) GetPort() int {
	return gdm.port
}

func (gdm getDataMarkImpl) GetIP() string {
	return gdm.ip
}

func (gdm getDataMarkImpl) GetApplicationName() string {
	return gdm.applicationName
}

func (getDataMarkImpl) IsGetInstance() bool {
	return true
}

func (instance) HandleRequest(ctx *zookeeper.Context) {
	isGetData := ctx.OpCode == zookeeper.OpGetData
	if !isGetData && ctx.OpCode != zookeeper.OpCreate {
		return
	}
	matches := dubbo.ServicePathPattern.FindStringSubmatch(ctx.Path)
	if matches == nil {
		return
	}
	if isGetData {
		// magic happens in response
		application, host, portStr := dubbo.GetHostPortFromMatches(matches)
		port, err := strconv.Atoi(portStr)
		if err != nil {
			log.DefaultLogger.Errorf("zookeeper.filters.instance.HandleRequest, parse endpoint port failed, port str: %s, %s", portStr, err)
			return
		}
		mark := &getDataMarkImpl{applicationName: application, ip: host, port: port}
		ctx.Value = mark
		return
	}
	var serviceInstance dubbo.ServiceInstance
	if err := json.Unmarshal(ctx.Data, &serviceInstance); err != nil {
		log.DefaultLogger.Errorf("zookeeper.filters.instance.Invoke, unmarshal service instance failed, %s", err)
		return
	}
	var application string
	var address string
	var port int
	if serviceInstance.Port != 0 {
		port = serviceInstance.Port
	}
	var portStr string
	application, address, portStr = dubbo.GetHostPortFromMatches(matches)
	if port == 0 {
		if p, err := strconv.ParseInt(portStr, 10, 32); err != nil {
			log.DefaultLogger.Errorf("zookeeper.filters.instance.Invoke, parse port failed, port string: %s, %s", portStr, err)
		} else {
			port = int(p)
		}
	}
	if serviceInstance.Address != "" {
		address = serviceInstance.Address
	}
	ctx.Value = &serviceInstance
	dubbo.SetDubboOriginalPort(port)
	dubbo.UpdateEndpointsByApplication(application, []string{fmt.Sprintf("127.0.0.1:%d", port)})
	if log.DefaultLogger.GetLogLevel() >= log.TRACE {
		log.DefaultLogger.Tracef("zookeeper.filters.instance.Invoke, application name: %s, address: %s, port: %d",
			serviceInstance.Name, address, port)
	}
	ctx.Path, ctx.Value = dubbo.ModifyInstanceInfo(ctx.Path, &serviceInstance)
	ctx.Modified = true
	ctx.Payload = &zookeeper.CreateRequest{
		XidAndOpCode: zookeeper.RequestHeader(ctx.RawPayload),
		TheRest:      zookeeper.OriginalContentView(ctx.RawPayload, ctx.DataEnd, zookeeper.Undefined),
	}
}

func (instance) HandleResponse(ctx *zookeeper.Context) {
	if ctx.Request == nil {
		return
	}
	if mark := ctx.Request.Value; mark != nil {
		if theMark, ok := mark.(getDataMark); ok && theMark.IsGetInstance() {
			handleGetDataResponse(ctx, theMark)
			return
		}
	}
	if !ctx.Request.Modified || ctx.Request.OriginalPath == "" {
		return
	}
	if ctx.Request.OpCode != zookeeper.OpCreate {
		return
	}
	matches := dubbo.ServicePathPattern.FindStringSubmatch(ctx.Request.OriginalPath)
	if matches == nil {
		return
	}
	if ctx.Error != nil {
		return
	}
	ctx.Path = ctx.Request.OriginalPath
	content := ctx.RawPayload
	// set path begin and end to let serialization get old path length
	ctx.PathBegin = 0
	ctx.PathEnd = int(binary.BigEndian.Uint32(content[5*zookeeper.Uint32Size : 6*zookeeper.Uint32Size]))
	ctx.Payload = &zookeeper.CreateResponse{
		XidZxidAndErrCode: zookeeper.ResponseHeader(content),
	}
	ctx.Modified = true
}

func handleGetDataResponse(ctx *zookeeper.Context, mark getDataMark) {
	if ctx.Error != nil {
		return
	}
	var serviceInstance dubbo.ServiceInstance
	if err := json.Unmarshal(ctx.Data, &serviceInstance); err != nil {
		log.DefaultLogger.Errorf("zookeeper.filters.instance.handleGetDataResponse, unmarshal data failed, %s", err)
		return
	}
	application := mark.GetApplicationName()
	revision := serviceInstance.Payload.Metadata[dubbo.MetadataRevisionKey]
	endpoint := dubbo.GetEndpointByApplication(application, mark.GetIP(), mark.GetPort())
	if endpoint != nil {
		endpoint.Revision = revision
	}
	exposePort := dubbo.GetDubboExposePort()
	serviceInstance.Address = "127.0.0.1"
	serviceInstance.ID = fmt.Sprintf("127.0.0.1:%d", exposePort)
	serviceInstance.Port = exposePort
	var endpoints []dubbo.Endpoint
	json.UnmarshalFromString(serviceInstance.Payload.Metadata[dubbo.MetadataEndpointKey], &endpoints)
	for i := range endpoints {
		endpoints[i].Port = exposePort
	}
	if metadataEndpoints, err := json.MarshalToString(endpoints); err != nil {
		log.DefaultLogger.Errorf("zookeeper.filters.instance.handleGetDataResponse, marshal modified endpoints failed, %s", err)
		return
	} else {
		serviceInstance.Payload.Metadata[dubbo.MetadataEndpointKey] = metadataEndpoints
	}
	ctx.Value = &serviceInstance
	ctx.Payload = &zookeeper.GetDataResponse{
		XidZxidAndErrCode: zookeeper.ResponseHeader(ctx.RawPayload),
		TheRest:           zookeeper.OriginalContentView(ctx.RawPayload, ctx.DataEnd, zookeeper.Undefined),
	}
	ctx.Modified = true
}
