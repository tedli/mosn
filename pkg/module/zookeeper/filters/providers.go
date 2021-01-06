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
	"mosn.io/mosn/pkg/module/dubbo"
	"mosn.io/mosn/pkg/module/zookeeper"
)

func NewProvidersFilter() zookeeper.Filter {
	return new(providers)
}

type listAndWatchProviderInstances interface {
	IsListProviders() bool
	ApplicationName() string
}

type listProvidersMark struct {
	applicationName string
}

func (listProvidersMark) IsListProviders() bool {
	return true
}

func (lpm listProvidersMark) ApplicationName() string {
	return lpm.applicationName
}

type providers struct{}

func (providers) HandleRequest(ctx *zookeeper.Context) {
	if ctx.OpCode != zookeeper.OpGetChildren2 {
		return
	}
	matches := dubbo.InstancesPathPattern.FindStringSubmatch(ctx.Path)
	if matches == nil {
		return
	}
	applicationName := dubbo.GetProviderApplicationFromMatches(matches)
	ctx.Value = &listProvidersMark{
		applicationName: applicationName,
	}
}

func (providers) HandleResponse(ctx *zookeeper.Context) {
	if ctx.Request == nil {
		return
	}
	var theMark listAndWatchProviderInstances
	var ok bool
	if mark := ctx.Request.Value; mark == nil {
		return
	} else if theMark, ok = mark.(listAndWatchProviderInstances); !ok || !theMark.IsListProviders() {
		return
	}
	content := ctx.RawPayload
	children := zookeeper.ParseChildren(content)
	ctx.Value = children
	applicationName := theMark.ApplicationName()
	dubbo.UpdateEndpointsByApplication(applicationName, children)
}
