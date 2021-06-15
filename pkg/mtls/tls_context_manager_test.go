package mtls

import "testing"

func TestIsMatchAlpns(t *testing.T) {
	result := IsMatchAlpns("dubbo", "ingress_dubbo")
	if !result {
		t.Errorf("IsMatchAlpns error")
	}

	maps := make(map[string]string)
	maps["ingress_dubbo"] = "dubbo"
	SetListenerNameWithAlpnProtocol(&maps)

	result = IsMatchAlpns("dubbo", "ingress_dubbo")
	if !result {
		t.Errorf("IsMatchAlpns error")
	}

	result = IsMatchAlpns("dubbo", "ingress_springcloud")
	if result {
		t.Errorf("IsMatchAlpns error")
	}

	result = IsMatchAlpns("springcloud", "ingress_dubbo")
	if result {
		t.Errorf("IsMatchAlpns error")
	}
}
