package dubbo

import (
	"context"

	"mosn.io/api"
	v2 "mosn.io/mosn/pkg/config/v2"
	mosnctx "mosn.io/mosn/pkg/context"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/protocol"
	"mosn.io/mosn/pkg/protocol/xprotocol/dubbo"
	"mosn.io/mosn/pkg/types"
	"mosn.io/mosn/pkg/variable"
	"mosn.io/pkg/buffer"
)

func init() {
	api.RegisterStream(v2.DubboStream, buildStream)
}

type factory struct{}

func buildStream(conf map[string]interface{}) (api.StreamFilterChainFactory, error) {
	Init(conf)
	return &factory{}, nil
}

func (f *factory) CreateFilterChain(ctx context.Context, callbacks api.StreamFilterChainFactoryCallbacks) {
	filter := buildDubboFilter(ctx)
	callbacks.AddStreamReceiverFilter(filter, api.BeforeRoute)
	callbacks.AddStreamSenderFilter(filter, api.BeforeSend)
}

type dubboFilter struct {
	handler api.StreamReceiverFilterHandler
}

func buildDubboFilter(ctx context.Context) *dubboFilter {
	return &dubboFilter{}
}

func (d *dubboFilter) OnDestroy() {}

func (d *dubboFilter) OnReceive(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {

	listener := mosnctx.Get(ctx, types.ContextKeyListenerName).(string)

	service, ok := headers.Get(dubbo.ServiceNameHeader)
	if !ok {
		log.DefaultLogger.Errorf("%s is empty, may be the protocol is not dubbo", dubbo.ServiceNameHeader)
		return api.StreamFiltertermination
	}

	// adapte dubbo service to http host
	variable.SetVariableValue(ctx, protocol.MosnHeaderHostKey, service)
	// because use http rule, so should add default path
	variable.SetVariableValue(ctx, protocol.MosnHeaderPathKey, "/")

	method, _ := headers.Get(dubbo.MethodNameHeader)
	stats := getStats(listener, service, method)
	if stats != nil {
		stats.RequestServiceInfo.Inc(1)

		variable.SetVariableValue(ctx, VarDubboRequestService, service)
		variable.SetVariableValue(ctx, VarDubboRequestMethod, method)
	}

	for k, v := range types.GetPodLabels() {
		headers.Set(k, v)
	}

	return api.StreamFilterContinue
}

func (d *dubboFilter) SetReceiveFilterHandler(handler api.StreamReceiverFilterHandler) {
	d.handler = handler
}

func (d *dubboFilter) Append(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	listener := mosnctx.Get(ctx, types.ContextKeyListenerName).(string)
	service, err := variable.GetVariableValue(ctx, VarDubboRequestService)
	if err != nil {
		log.DefaultLogger.Warnf("Get request service info failed: %+v", err)
		return api.StreamFilterContinue
	}
	method, err := variable.GetVariableValue(ctx, VarDubboRequestMethod)
	if err != nil {
		log.DefaultLogger.Warnf("Get request method info failed: %+v", err)
		return api.StreamFilterContinue
	}

	stats := getStats(listener, service, method)
	if stats == nil {
		return api.StreamFilterContinue
	}

	var isSuccess bool
	switch frame := headers.(type) {
	case *dubbo.Frame:
		isSuccess = frame.GetStatusCode() == dubbo.RespStatusOK
	default:
		log.DefaultLogger.Errorf("this filter {%s} just for dubbo protocol, please check your config.", v2.DubboStream)
		return api.StreamFiltertermination
	}

	if isSuccess {
		stats.ResponseSucc.Inc(1)
	} else {
		stats.ResponseFail.Inc(1)
	}
	return api.StreamFilterContinue
}

func (d *dubboFilter) SetSenderFilterHandler(handler api.StreamSenderFilterHandler) {}
