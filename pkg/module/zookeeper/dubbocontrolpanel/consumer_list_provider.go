package dubbocontrolpanel

import (
	"mosn.io/mosn/pkg/module/dubbo"
	"mosn.io/mosn/pkg/module/zookeeper"
)

func handleConsumerListProviders(upstream zookeeper.Upstream, request *zookeeper.Context) {
	var interfaceName string
	request.MustGetParam("interface", &interfaceName)
	_, response := upstream.DirectForward(request)
	if response.Error != nil {
		return
	}
	children := zookeeper.ParseChildren(response.RawPayload)
	response.Value = children
	dubbo.UpdateApplicationsByInterface(interfaceName, children)
}

func handleConsumerListProviderApplications(upstream zookeeper.Upstream, request *zookeeper.Context) {
	_, response := upstream.DirectForward(request)
	content := response.RawPayload
	children := zookeeper.ParseChildren(content)
	request.Value = children
	var application string
	request.MustGetParam("application", &application)
	dubbo.UpdateEndpointsByApplication(application, children)
}
