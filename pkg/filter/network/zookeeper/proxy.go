package zookeeper

import (
	"bytes"
	"context"
	"sync"

	"mosn.io/api"
	module "mosn.io/mosn/pkg/module/zookeeper"
	zkapi "mosn.io/mosn/pkg/module/zookeeper/api"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/log"
)

func init() {
	api.RegisterNetwork("zookeeper", func(_ map[string]interface{}) (api.NetworkFilterChainFactory, error) {
		// TODO: there should be a config, to control what filters should be used
		return new(factory), nil
	})
}

type factory struct {
}

func (f factory) CreateFilterChain(_ context.Context, callbacks api.NetWorkFilterChainFactoryCallbacks) {
	// TODO: the filters should be configured through a config
	filter := newInterceptor([]module.Filter{
		zkapi.HeaderFilter,
		zkapi.SessionFilter,
		zkapi.PathFilter,
		zkapi.DataFilter,
		zkapi.MetadataFilter,
		zkapi.InstanceFilter,
		zkapi.ProvidersFilter,
		zkapi.DebugFilter,
	})
	callbacks.AddReadFilter(filter)
	callbacks.AddWriteFilter(filter)
}

func (i *interceptor) OnWrite(content []buffer.IoBuffer) api.FilterStatus {
	var condition module.IOStateCondition

	if len(content) == 1 {
		condition = i.response.OnData(content[0], false)
	} else {
		buf := new(bytes.Buffer)
		for _, c := range content {
			if _, err := c.WriteTo(buf); err != nil {
				log.DefaultLogger.Errorf("network.zookeeper.responseInterceptor.OnWrite, aggravate response data failed, %s", err)
				return api.Continue
			}
		}
		condition = i.response.OnData(buf, false)
	}

	switch condition {
	case module.ZookeeperNeedMoreData:
		return api.Stop
	case module.ZookeeperFinishedFiltering:
		return api.Continue
	}
	return api.Continue
}

func newInterceptor(filters []module.Filter) *interceptor {
	session := new(sync.Map)
	return &interceptor{
		request:  module.NewManager(filters, session),
		response: module.NewManager(filters, session),
	}
}

type interceptor struct {
	request, response module.Manager
}

func (i *interceptor) OnData(content buffer.IoBuffer) api.FilterStatus {

	condition := i.request.OnData(content, true)
	switch condition {
	case module.ZookeeperNeedMoreData:
		return api.Stop
	case module.ZookeeperFinishedFiltering:
		return api.Continue
	}
	return api.Continue
}

func (i interceptor) OnNewConnection() api.FilterStatus {
	return api.Continue
}

func (i interceptor) InitializeReadFilterCallbacks(_ api.ReadFilterCallbacks) {
}

func (i interceptor) ReadDisableUpstream(_ bool) {
}

func (i interceptor) ReadDisableDownstream(_ bool) {
}
