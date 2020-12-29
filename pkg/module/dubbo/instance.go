package dubbo

import (
	"fmt"
	"regexp"

	jsoniter "github.com/json-iterator/go"
	"mosn.io/pkg/log"
)

var (
	dubboExposePort = 20888
	dubboOriginalPort = 20880

	MetadataEndpointKey = "dubbo.endpoints"
	MetadataRevisionKey = "dubbo.metadata.revision"

	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func SetDubboExposePort(port int) {
	dubboExposePort = port
}

func GetDubboExposePort() int {
	return dubboExposePort
}

func SetDubboOriginalPort(port int) {
	dubboOriginalPort = port
}

func GetDubboOriginalPort() int {
	return dubboOriginalPort
}

func ModifyInstanceInfo(path string, serviceInstance *ServiceInstance) (
	modifiedPath string, modifiedServiceInstance *ServiceInstance) {

	portPattern := fmt.Sprintf("$1:%d", dubboExposePort)
	pattern, err := regexp.Compile(fmt.Sprintf(`^(\S+):%d$`, serviceInstance.Port))
	if err != nil {
		log.DefaultLogger.Errorf("dubbo.instance.ModifyInstanceInfo, compile path port replace pattern failed, %s", err)
	} else {
		modifiedPath = pattern.ReplaceAllString(path, portPattern)
	}

	modifiedServiceInstance = serviceInstance
	if pattern != nil {
		modifiedServiceInstance.ID = pattern.ReplaceAllString(modifiedServiceInstance.ID, portPattern)
	}
	modifiedServiceInstance.Port = dubboExposePort
	if metadata := modifiedServiceInstance.Payload.Metadata; metadata != nil {
		if rawEndpoints, exist := metadata[MetadataEndpointKey]; exist {
			var endpoints []*Endpoint
			if err := json.UnmarshalFromString(rawEndpoints, &endpoints); err != nil {
				log.DefaultLogger.Errorf("dubbo.instance.ModifyInstanceInfo, unmarshal endpoints failed, endpoints: %s, %s", rawEndpoints, err)
				return
			}
			for _, endpoint := range endpoints {
				if endpoint.Protocol == "dubbo" {
					endpoint.Port = dubboExposePort
				}
			}
			if modified, err := json.MarshalToString(endpoints); err != nil {
				log.DefaultLogger.Errorf("dubbo.instance.ModifyInstanceInfo, marshal modified endpoints failed, %s", err)
				return
			} else {
				metadata[MetadataEndpointKey] = modified
			}
		}
	}
	return
}
