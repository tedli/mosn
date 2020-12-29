package dubbo

import (
	"regexp"
)

var (
	ServicePathPattern = regexp.MustCompile(`^/services/(?P<application>[^ \f\n\r\t\v/]+)/(?P<host>[^ \f\n\r\t\v/]+):(?P<port>\d+)$`)
	InstancesPathPattern = regexp.MustCompile(`^/services/(?P<application>[^ \f\n\r\t\v/]+)$`)

	applicationNameIndex, hostIndex, portIndex = func() (application, host, port int) {
		names := ServicePathPattern.SubexpNames()
		for index, name := range names {
			if name == "application" {
				application = index
			} else if name == "host" {
				host = index
			} else if name == "port" {
				port = index
			}
		}
		return
	}()

	providerApplicationNameIndex = func() (application int) {
		names := InstancesPathPattern.SubexpNames()
		for index, name := range names {
			if name == "application" {
				application = index
				break
			}
		}
		return
	}()
)

func GetProviderApplicationFromMatches(matches []string) (application string) {
	if length := len(matches); length <= providerApplicationNameIndex {
		return
	}
	application = matches[providerApplicationNameIndex]
	return
}

func GetHostPortFromMatches(matches []string) (application, host, port string) {
	if length := len(matches); length <= portIndex {
		return
	}
	application = matches[applicationNameIndex]
	host = matches[hostIndex]
	port = matches[portIndex]
	return
}

// C:\Users\lizhen\.m2\repository\org\apache\curator\curator-x-discovery\4.0.1\curator-x-discovery-4.0.1.jar!\org\apache\curator\x\discovery\ServiceInstance.class
type ServiceInstance struct {
	Name                string                   `json:"name"`
	ID                  string                   `json:"id"`
	Address             string                   `json:"address"`
	Port                int                      `json:"port"`
	SslPort             *int                     `json:"sslPort"`
	Payload             ZookeeperInstancePayload `json:"payload"`
	RegistrationTimeUTC int64                    `json:"registrationTimeUTC"`
	ServiceType         string                   `json:"serviceType"`
	URISpec             *URISpec                 `json:"uriSpec"`
	Enabled             *bool                    `json:"enabled"`
}

type ZookeeperInstancePayload struct {
	Class             string `json:"@class"`
	ZookeeperInstance `json:",inline"`
}

// dubbo-registry/dubbo-registry-zookeeper/src/main/java/org/apache/dubbo/registry/zookeeper/ZookeeperInstance.java
type ZookeeperInstance struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Metadata map[string]string `json:"metadata"`
}

// C:\Users\lizhen\.m2\repository\org\apache\curator\curator-x-discovery\4.0.1\curator-x-discovery-4.0.1.jar!\org\apache\curator\x\discovery\UriSpec.class
type URISpecPart struct {
	Value    string `json:"value"`
	Variable bool   `json:"variable"`
}

type URISpec struct {
	Parts []URISpecPart `json:"parts"`
}

type URLParams struct {
	Dubbo struct{
		Version string `json:"version"`
		Dubbo string `json:"dubbo"`
		Port string `json:"port"`
	} `json:"dubbo"`
}

type Endpoint struct {
	Port int `json:"port"`
	Protocol string `json:"protocol"`
}
