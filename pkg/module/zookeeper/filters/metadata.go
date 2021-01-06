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
	jsoniter "github.com/json-iterator/go"
	"mosn.io/mosn/pkg/module/dubbo"
	"mosn.io/mosn/pkg/module/zookeeper"
	"mosn.io/pkg/log"
)

func NewMetadataFilters() zookeeper.Filter {
	return new(metadata)
}

type metadata struct {
}

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type interfaceApplicationMapping interface {
	IsInterfaceApplicationMapping() bool
	GetInterfaceName() string
}

type mapping struct {
	interfaceName string
}

func (m mapping) GetInterfaceName() string {
	return m.interfaceName
}

func (mapping) IsInterfaceApplicationMapping() bool {
	return true
}

type getMetadataMapping interface {
	IsGetMetadataMapping() bool
}

type metadataMappingRequest struct{}

func (metadataMappingRequest) IsGetMetadataMapping() bool {
	return true
}

func (metadata) HandleRequest(ctx *zookeeper.Context) {
	isGetChildren := ctx.OpCode == zookeeper.OpGetChildren2
	isGetData := ctx.OpCode == zookeeper.OpGetData
	if !isGetChildren && !isGetData && ctx.OpCode != zookeeper.OpCreate {
		return
	}
	if isGetChildren {
		matches := dubbo.InterfaceApplicationMappingPathPattern.FindStringSubmatch(ctx.Path)
		if matches == nil {
			return
		}
		interfaceName := dubbo.GetInterfaceFromMatches(matches)
		ctx.Value = &mapping{interfaceName: interfaceName}
		return
	}
	matches := dubbo.MetadataPathPattern.FindStringSubmatch(ctx.Path)
	if matches != nil {
		handleApplicationMetadata(ctx, matches)
		return
	}
	matches = dubbo.MetadataMappingPathPattern.FindStringSubmatch(ctx.Path)
	if matches != nil {
		if isGetData {
			ctx.Value = new(metadataMappingRequest)
			return
		}
		handleServiceApplicationMapping(ctx)
		return
	}
}

func handleServiceApplicationMapping(ctx *zookeeper.Context) {
	var metainfo dubbo.MetadataInfo
	if err := json.Unmarshal(ctx.Data, &metainfo); err != nil {
		log.DefaultLogger.Errorf("zookeeper.filter.metadat.handleServiceApplicationMapping, unmarshal metainfo failed, %s", err)
		return
	}
	services := make([]dubbo.ServiceInfo, 0, len(metainfo.Services))
	for _, service := range metainfo.Services {
		services = append(services, service)
	}
	dubbo.UpdateClustersByProvider(metainfo.Application, metainfo.Revision, services)
}

func handleApplicationMetadata(ctx *zookeeper.Context, matches []string) {
	serviceApplicationRef := dubbo.GetServiceApplicationRef(matches)
	var application string
	if serviceApplicationRef.Side == dubbo.ServiceSideConsumer {
		application = handleConsumer(ctx, serviceApplicationRef)
		return
	} else {
		application = handleProvider(ctx, serviceApplicationRef)
	}
	dubbo.SetApplicationName(application)
	if log.DefaultLogger.GetLogLevel() >= log.TRACE {
		log.DefaultLogger.Tracef("zookeeper.filters.metadata.Invoke, application name: %s", application)
	}
}

func handleConsumer(ctx *zookeeper.Context, serviceApplicationRef *dubbo.ServiceApplicationRef) (application string) {
	var parameters map[string]string
	if err := json.Unmarshal(ctx.Data, &parameters); err != nil {
		log.DefaultLogger.Errorf("zookeeper.filters.metadata.Invoke, unmarshal consumer parameters failed, %s", err)
		return
	}
	ctx.Value = parameters
	if parameters == nil {
		return
	}
	application = parameters[dubbo.ParameterApplicationKey]
	return
}

func handleProvider(ctx *zookeeper.Context, serviceApplicationRef *dubbo.ServiceApplicationRef) (application string) {
	var serviceDefinition dubbo.FullServiceDefinition
	if err := json.Unmarshal(ctx.Data, &serviceDefinition); err != nil {
		log.DefaultLogger.Errorf("zookeeper.filters.metadata.Invoke, unmarshal service definition failed, %s", err)
		return
	}
	dubbo.UpdateApplicationsByInterface(serviceApplicationRef.Interface, []string{serviceApplicationRef.Application})
	ctx.Value = &serviceDefinition
	if applicationParameter, exist := serviceDefinition.Parameters[dubbo.ParameterApplicationKey]; exist && applicationParameter != "" {
		application = applicationParameter
	} else {
		application = serviceApplicationRef.Application
	}
	return
}

func (metadata) HandleResponse(ctx *zookeeper.Context) {
	var mark interfaceApplicationMapping
	var getMapping getMetadataMapping
	var ok bool
	if ctx.Request == nil {
		return
	} else if value := ctx.Request.Value; value == nil {
		return
	} else if mark, ok = value.(interfaceApplicationMapping); ok && mark.IsInterfaceApplicationMapping() {
		children := zookeeper.ParseChildren(ctx.RawPayload)
		ctx.Value = children
		dubbo.UpdateApplicationsByInterface(mark.GetInterfaceName(), children)
	} else if getMapping, ok = value.(getMetadataMapping); ok && getMapping.IsGetMetadataMapping() {
		handleServiceApplicationMapping(ctx)
	}
}
