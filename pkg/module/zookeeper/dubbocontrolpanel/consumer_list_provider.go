package dubbocontrolpanel

import (
	"mosn.io/mosn/pkg/module/dubbo"
	"mosn.io/mosn/pkg/module/zookeeper"
)

func handleConsumerListProviders(upstream zookeeper.Upstream, request *zookeeper.Context) {
	var interfaceName string
	request.MustGetParam("interface", &interfaceName)
	downstream, response := upstream.DirectForward(request)
	if response.Error != nil {
		downstream.DirectReply(response)
		return
	}
	children := zookeeper.ParseChildren(response.RawPayload)
	response.Value = children
	dubbo.UpdateApplicationsByInterface(interfaceName, children)
	downstream.DirectReply(response)
}
