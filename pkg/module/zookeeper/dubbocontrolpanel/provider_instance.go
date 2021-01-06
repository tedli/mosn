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

package dubbocontrolpanel

import (
	"encoding/binary"
	"fmt"

	"mosn.io/mosn/pkg/module/dubbo"
	"mosn.io/mosn/pkg/module/zookeeper"
	"mosn.io/pkg/log"
)

func handleProviderInstanceRegister(upstream zookeeper.Upstream, request *zookeeper.Context) {
	var serviceInstance dubbo.ServiceInstance
	if err := json.Unmarshal(request.Data, &serviceInstance); err != nil {
		log.DefaultLogger.Errorf("zookeeper.filters.instance.Invoke, unmarshal service instance failed, %s", err)
		return
	}
	var application string
	var address string
	var port int
	if serviceInstance.Port != 0 {
		port = serviceInstance.Port
	}
	if port == 0 {
		request.MustGetParam("port", &port)
	}
	if serviceInstance.Address != "" {
		address = serviceInstance.Address
	}
	request.Value = &serviceInstance
	dubbo.SetDubboOriginalPort(port)
	dubbo.UpdateEndpointsByApplication(application, []string{fmt.Sprintf("127.0.0.1:%d", port)})
	if log.DefaultLogger.GetLogLevel() >= log.TRACE {
		log.DefaultLogger.Tracef("zookeeper.filters.instance.Invoke, application name: %s, address: %s, port: %d",
			serviceInstance.Name, address, port)
	}
	request.Path, request.Value = dubbo.ModifyInstanceInfo(request.Path, &serviceInstance)
	downstream, response := upstream.ModifyAndSend(request, &zookeeper.CreateRequest{
		XidAndOpCode: zookeeper.RequestHeader(request.RawPayload),
		TheRest:      zookeeper.OriginalContentView(request.RawPayload, request.DataEnd, zookeeper.Undefined),
	})
	if response.Error != nil {
		downstream.DirectReply(response)
		return
	}
	response.Path = response.Request.OriginalPath
	content := response.RawPayload
	// set path begin and end to let serialization get old path length
	response.PathBegin = 0
	response.PathEnd = int(binary.BigEndian.Uint32(content[5*zookeeper.Uint32Size : 6*zookeeper.Uint32Size]))
	downstream.ModifyAndReply(response, &zookeeper.CreateResponse{
		XidZxidAndErrCode: zookeeper.ResponseHeader(content),
	})
}
