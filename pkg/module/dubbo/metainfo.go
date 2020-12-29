package dubbo

// dubbo-metadata/dubbo-metadata-api/src/main/java/org/apache/dubbo/metadata/MetadataInfo.java
type MetadataInfo struct {
	Application string                 `json:"app"`
	Revision    string                 `json:"revision"`
	Services    map[string]ServiceInfo `json:"services"`
	// ExtendParams map[string]string `json:""`
	// and others fields
}

// dubbo-metadata/dubbo-metadata-api/src/main/java/org/apache/dubbo/metadata/MetadataInfo.java
type ServiceInfo struct {
	Name     string            `json:"name"`
	Group    string            `json:"group"`
	Version  string            `json:"version"`
	Protocol string            `json:"protocol"`
	Path     string            `json:"path"`
	Params   map[string]string `json:"params"`
	// and other fields
}
