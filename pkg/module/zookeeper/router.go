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

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"mosn.io/mosn/pkg/log"
)

func NewRouteFilter() Filter {
	return new(routerFilter)
}

type upstream struct {
	handleRequestFinished *sync.WaitGroup
	startToHandleResponse *sync.WaitGroup
	response              *downstream
}

func newUpstream() *upstream {
	handleRequestFinished := new(sync.WaitGroup)
	handleRequestFinished.Add(1)
	startToHandleResponse := new(sync.WaitGroup)
	startToHandleResponse.Add(1)
	return &upstream{
		handleRequestFinished: handleRequestFinished,
		startToHandleResponse: startToHandleResponse,
	}
}

func (u *upstream) DirectForward(*Context) (Downstream, *Context) {
	u.handleRequestFinished.Done()
	u.startToHandleResponse.Wait()
	return u.response, u.response.response
}

func (u *upstream) ModifyAndSend(ctx *Context, message interface{}) (Downstream, *Context) {
	ctx.Payload = message
	ctx.Modified = true
	return u.DirectForward(ctx)
}

type downstream struct {
	response               *Context
	handleResponseFinished *sync.WaitGroup
}

func newDownstream(response *Context) *downstream {
	handleResponseFinished := new(sync.WaitGroup)
	handleResponseFinished.Add(1)
	return &downstream{
		response:               response,
		handleResponseFinished: handleResponseFinished,
	}
}

func (d *downstream) DirectReply(ctx *Context) {
	d.handleResponseFinished.Done()
}

func (d *downstream) ModifyAndReply(ctx *Context, message interface{}) {
	ctx.Payload = message
	ctx.Modified = true
	d.DirectReply(ctx)
}

type routerFilter struct {
}

func (routerFilter) HandleRequest(ctx *Context) {
	route, matches := resolveRoute(ctx)
	if route == nil {
		return
	}
	params := prepareParams(route, matches)
	ctx.params = params
	u := newUpstream()
	ctx.route = u
	go route.handler(u, ctx)
	u.handleRequestFinished.Wait()
}

func (routerFilter) HandleResponse(ctx *Context) {
	if ctx.Request == nil || ctx.Request.route == nil {
		return
	}
	rc := ctx.Request.route
	d := newDownstream(ctx)
	rc.response = d
	ctx.Request.route = nil
	rc.startToHandleResponse.Done()
	d.handleResponseFinished.Wait()
}

func prepareParams(route *entry, matches []string) (params map[string]interface{}) {
	matcheLength := len(matches)
	params = make(map[string]interface{}, matcheLength)
	for name, param := range route.captures {
		if param.matchIndex >= matcheLength {
			log.DefaultLogger.Errorf("zookeeper.router.prepareParams, capture index out of bound")
			continue
		}
		value := matches[param.matchIndex]
		if param.paramType == paramTypeInt {
			i, e := strconv.Atoi(value)
			if e != nil {
				log.DefaultLogger.Errorf("zookeeper.route.prepareParams, value is not int, value: %s, %s", value, e)
			} else {
				params[name] = i
				continue
			}
		}
		params[name] = value
	}
	return
}

func resolveRoute(ctx *Context) (e *entry, matches []string) {
	patterns, exist := routers[ctx.OpCode]
	if !exist {
		return
	}
	for _, pattern := range patterns {
		if matches = pattern.pattern.FindStringSubmatch(ctx.Path); matches != nil {
			e = pattern
			break
		}
	}
	return
}

type Downstream interface {
	DirectReply(*Context)
	ModifyAndReply(ctx *Context, message interface{})
}

type Upstream interface {
	DirectForward(*Context) (Downstream, *Context)
	ModifyAndSend(ctx *Context, message interface{}) (Downstream, *Context)
}

type Handler func(Upstream, *Context)

type entry struct {
	handler  Handler
	captures map[string]*capture
	pattern  *regexp.Regexp
}

var (
	routers = make(map[OpCode]map[string]*entry)
)

func MustRegister(opcode OpCode, path string, handler Handler) {
	if err := Register(opcode, path, handler); err != nil {
		panic(err)
	}
}

func Register(opcode OpCode, path string, handler Handler) error {
	mapping, exist := routers[opcode]
	if !exist {
		mapping = make(map[string]*entry)
		routers[opcode] = mapping
	}
	patternStr, pattern, captures, err := pathToRegexp(path)
	if err != nil {
		return err
	}
	if _, exist := mapping[patternStr]; exist {
		return ErrRouterConflict
	}
	mapping[patternStr] = &entry{
		handler:  handler,
		captures: captures,
		pattern:  pattern,
	}
	return nil
}

const (
	pathElement          = `[^ \f\n\r\t\v\\}/]`
	stringCapturePattern = `(?P<%s>[^ \f\n\r\t\v/]+)`
	intCapturePattern    = `(?P<%s>\d+)`
)

var (
	paramPattern = regexp.MustCompile(fmt.Sprintf("\\{%s+\\}", pathElement))

	ErrUnsupportedParamType      = errors.New("unsupported param type")
	ErrSameParamNameMoreThanOnce = errors.New("same param name more than once")
	ErrRouterConflict            = errors.New("router conflict")
)

type capture struct {
	paramName  string
	paramType  string
	matchIndex int
}

const (
	paramTypeString = "string"
	paramTypeInt    = "int"
)

func pathToRegexp(path string) (patternStr string, pattern *regexp.Regexp, captures map[string]*capture, err error) {
	captures = make(map[string]*capture)
	rawPattern := paramPattern.ReplaceAllStringFunc(path, func(param string) string {
		param = strings.Trim(param, "{}")
		parts := strings.Split(param, ":")
		c := new(capture)
		if lp := len(parts); lp < 2 {
			c.paramName = param
		} else if lp == 2 {
			c.paramName = parts[1]
			if t := parts[0]; t != paramTypeString && t != paramTypeInt {
				err = ErrUnsupportedParamType
				return ""
			}
			c.paramType = parts[0]
		}
		if _, exist := captures[c.paramName]; exist {
			err = ErrSameParamNameMoreThanOnce
			return ""
		}
		captures[c.paramName] = c
		if c.paramType == paramTypeInt {
			return fmt.Sprintf(intCapturePattern, c.paramName)
		}
		return fmt.Sprintf(stringCapturePattern, c.paramName)
	})
	if err != nil {
		return
	}
	patternStr = fmt.Sprintf("^%s$", rawPattern)
	if pattern, err = regexp.Compile(patternStr); err != nil {
		return
	}
	names := pattern.SubexpNames()
	for index, name := range names {
		if c, wanted := captures[name]; wanted {
			c.matchIndex = index
		}
	}
	return
}
