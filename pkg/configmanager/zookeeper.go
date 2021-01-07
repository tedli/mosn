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

package configmanager

import (
	"encoding/json"
	"net"

	"mosn.io/mosn/pkg/config/v2"
	"mosn.io/mosn/pkg/log"
)

func ZookeeperConfigLoad(zookeerAddress string) *v2.MOSNConfig {
	log.StartLogger.Infof("[config] [zookeeper load] Generating config by ZookeeperConfigLoad")
	return &v2.MOSNConfig{
		Servers: []v2.ServerConfig{
			{
				DefaultLogPath:  "stdout",
				DefaultLogLevel: "DEBUG",
				Listeners: []v2.Listener{
					{
						Addr: &net.TCPAddr{
							IP:   net.ParseIP("0.0.0.0"),
							Port: 2182,
						},
						ListenerConfig: v2.ListenerConfig{
							Name:       "litener",
							AddrConfig: "0.0.0.0:2182",
							BindToPort: true,
							FilterChains: []v2.FilterChain{
								{
									FilterChainConfig: v2.FilterChainConfig{
										Filters: []v2.Filter{
											{
												Type: "zookeeper",
											},
											{
												Type: v2.TCP_PROXY,
												Config: map[string]interface{}{
													"cluster": "cluster",
													"routes": map[string]interface{}{
														"cluster": "cluster",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		ClusterManager: v2.ClusterManagerConfig{
			Clusters: []v2.Cluster{
				{
					Name:                 "cluster",
					ClusterType:          v2.SIMPLE_CLUSTER,
					LbType:               v2.LB_RANDOM,
					MaxRequestPerConn:    1024,
					ConnBufferLimitBytes: 32768,
					Hosts: []v2.Host{
						{
							HostConfig: v2.HostConfig{
								Address: zookeerAddress,
							},
						},
					},
				},
			},
		},
		RawAdmin: json.RawMessage(`{"address": {"socket_address": {"address": "127.0.0.1","port_value": 34901}}}`),
	}
}
