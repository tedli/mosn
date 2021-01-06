package dubbocontrolpanel

import (
	"mosn.io/mosn/pkg/module/dubbo"
	"mosn.io/mosn/pkg/module/zookeeper"
	"mosn.io/pkg/log"
)

func handleConsumerCreateMetadata(upstream zookeeper.Upstream, request *zookeeper.Context) {
	var application string
	var interfaceName string
	request.MustGetParam("application", &application)
	request.MustGetParam("interface", &interfaceName)
	var parameters map[string]string
	if err := json.Unmarshal(request.Data, &parameters); err != nil {
		log.DefaultLogger.Errorf("zookeeper.filters.metadata.Invoke, unmarshal consumer parameters failed, %s", err)
		upstream.Passthrough()
		return
	}
	request.Value = parameters
	if parameters != nil {
		application = parameters[dubbo.ParameterApplicationKey]
	}
	dubbo.SetApplicationName(application)
	if log.DefaultLogger.GetLogLevel() >= log.TRACE {
		log.DefaultLogger.Tracef("zookeeper.filters.metadata.Invoke, application name: %s", application)
	}
	downstream, response := upstream.DirectForward(request)
	downstream.DirectReply(response)
}
