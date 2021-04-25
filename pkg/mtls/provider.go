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

package mtls

import (
	"os"
	"strconv"

	"mosn.io/mosn/pkg/config/v2"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/types"
)

// GlobalMTLS is global mtls cofnig
var GlobalMTLS v2.GlobalMTLSConfig

// staticProvider is an implementation of types.Provider
// staticProvider stored a static certificate
type staticProvider struct {
	*tlsContext
}

func (p *staticProvider) Ready() bool {
	return true
}

func (p *staticProvider) Empty() bool {
	return p.tlsContext.server == nil
}

// IsGlobalMTLS return false by default
func IsGlobalMTLS() bool {
	isEnableMTLS, err := strconv.ParseBool(os.Getenv("ENABLE_GLOBAL_MTLS"))
	if err != nil {
		return false
	}
	if isEnableMTLS && GlobalMTLS.ClientTLSContext.Status && GlobalMTLS.ServerTLSContext.Status {
		return true
	}
	return false
}

// NewProvider returns a types.Provider.
// we support sds provider and static provider.
func NewProvider(cfg *v2.TLSConfig) (types.TLSProvider, error) {
	log.DefaultLogger.Infof("[mtls] Provider TLS Config status: %v", cfg.Status)
	log.DefaultLogger.Infof("[mtls] [NewProvider] config:%v", cfg)
	if !cfg.Status /*|| !IsGlobalMTLS()*/ {
		return nil, nil
	}
	if cfg.SdsConfig != nil {
		log.DefaultLogger.Infof("[mtls] Provider TLS sds config: %v", cfg.SdsConfig )
		log.DefaultLogger.Infof("[mtls]SdsConfig config Valid, c != nil %v, CertificateConfig %v, ValidationConfig %v", cfg.SdsConfig != nil, cfg.SdsConfig.CertificateConfig, cfg.SdsConfig.ValidationConfig)
		if !cfg.SdsConfig.Valid() {
			return nil, ErrorNoCertConfigure
		}
		return getOrCreateProvider(cfg), nil
	} else {
		// static provider
		secret := &secretInfo{
			Certificate: cfg.CertChain,
			PrivateKey:  cfg.PrivateKey,
			Validation:  cfg.CACert,
		}
		ctx, err := newTLSContext(cfg, secret)
		if err != nil {
			return nil, err
		}
		return &staticProvider{
			tlsContext: ctx,
		}, nil
	}
}
