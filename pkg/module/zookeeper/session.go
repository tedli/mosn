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

package zookeeper

import "encoding/binary"

func NewSessionFilter() Filter {
	return new(session)
}

type session struct{}

func (session) HandleRequest(ctx *Context) {
	ctx.session.Store(ctx.Xid, ctx)
}

func (session) HandleResponse(ctx *Context) {
	request, loaded := ctx.session.LoadAndDelete(ctx.Xid)
	if !loaded {
		return
	}
	ctx.Request = request.(*Context)
	ctx.Watch = ctx.Request.Watch
}

const (
	childrenCountEnd = 6 * Uint32Size
)

func ParseChildren(content []byte) (children []string) {
	childrenCount := int(binary.BigEndian.Uint32(content[5*Uint32Size : childrenCountEnd]))
	children = make([]string, childrenCount)
	var end int
	for i, begin := 0, childrenCountEnd; i < childrenCount; i++ {
		end = begin + Uint32Size
		pathLength := int(binary.BigEndian.Uint32(content[begin:end]))
		begin = end
		end = begin + pathLength
		path := string(content[begin:end])
		children[i] = path
		begin = end
	}
	return
}