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
