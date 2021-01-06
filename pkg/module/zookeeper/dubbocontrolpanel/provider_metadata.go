package dubbocontrolpanel

import (
	"mosn.io/mosn/pkg/module/dubbo"
	"mosn.io/mosn/pkg/module/zookeeper"
	"mosn.io/pkg/log"
)

func handleProviderCreateMetadata(upstream zookeeper.Upstream, request *zookeeper.Context) {
	var application string
	var interfaceName string
	request.MustGetParam("application", &application)
	request.MustGetParam("interface", &interfaceName)
	var serviceDefinition dubbo.FullServiceDefinition
	if err := json.Unmarshal(request.Data, &serviceDefinition); err != nil {
		log.DefaultLogger.Errorf("zookeeper.filters.metadata.Invoke, unmarshal service definition failed, %s", err)
		upstream.Passthrough()
		return
	}
	dubbo.UpdateApplicationsByInterface(interfaceName, []string{application})
	request.Value = &serviceDefinition
	if applicationParameter, exist := serviceDefinition.Parameters[dubbo.ParameterApplicationKey]; exist && applicationParameter != "" {
		application = applicationParameter
	}
	dubbo.SetApplicationName(application)
	if log.DefaultLogger.GetLogLevel() >= log.TRACE {
		log.DefaultLogger.Tracef("zookeeper.filters.metadata.Invoke, application name: %s", application)
	}
	downstream, response := upstream.DirectForward(request)
	downstream.DirectReply(response)
}
