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

package hooks

import (
	"context"
	"mosn.io/mosn/pkg/types"
)

type Request interface {
	GetHeader(key string) string
	GetBody() []byte
}

type Hook interface {
	Invoke(ctx context.Context, headers types.HeaderMap, serviceName string, request Request) (err error)
}

type serviceNameBuilder func(headers types.HeaderMap, ctx context.Context) (serviceID, serviceType string)

var (
	beforeEncodeHooks                    = make([]Hook, 0)
	afterDecodeHooks                     = make([]Hook, 0)
	buildServiceName  serviceNameBuilder = func(headers types.HeaderMap, ctx context.Context) (serviceID, serviceType string) {
		return
	}
)

func RegisterBeforeEncodeHook(hook Hook) {
	beforeEncodeHooks = append(beforeEncodeHooks, hook)
}

func RegisterAfterDecodeHook(hook Hook) {
	afterDecodeHooks = append(afterDecodeHooks, hook)
}

func RegisterServiceNameBuilder(builder serviceNameBuilder) {
	buildServiceName = builder
}

func invoke(
	hooks []Hook, ctx context.Context, headers types.HeaderMap, serviceName string, request Request) (err error) {

	for _, hook := range hooks {
		if err = hook.Invoke(ctx, headers, serviceName, request); err != nil {
			return
		}
	}
	return
}

func BeforeEncode(ctx context.Context, headers types.HeaderMap, serviceName string, request Request) (err error) {
	err = invoke(beforeEncodeHooks, ctx, headers, serviceName, request)
	return
}

func AfterDecode(ctx context.Context, headers types.HeaderMap, serviceName string, request Request) (err error) {
	err = invoke(afterDecodeHooks, ctx, headers, serviceName, request)
	return
}

func BuildServiceName(headers types.HeaderMap, ctx context.Context) (serviceID, serviceType string) {
	serviceID, serviceType = buildServiceName(toReadonly(headers), ctx)
	return
}
