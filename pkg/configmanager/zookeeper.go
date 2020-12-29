package configmanager

import (
	"net"
	"encoding/json"

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
							Port: 2181,
						},
						ListenerConfig: v2.ListenerConfig{
							Name:       "litener",
							AddrConfig: "0.0.0.0:2181",
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
