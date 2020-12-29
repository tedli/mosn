package filters

import (
	"mosn.io/mosn/pkg/module/dubbo"
	"mosn.io/mosn/pkg/module/zookeeper"
)

func NewProvidersFilter() zookeeper.Filter {
	return new(providers)
}

type listAndWatchProviderInstances interface {
	IsListProviders() bool
	ApplicationName() string
}

type listProvidersMark struct {
	applicationName string
}

func (listProvidersMark) IsListProviders() bool {
	return true
}

func (lpm listProvidersMark) ApplicationName() string {
	return lpm.applicationName
}

type providers struct{}

func (providers) HandleRequest(ctx *zookeeper.Context) {
	if ctx.OpCode != zookeeper.OpGetChildren2 {
		return
	}
	matches := dubbo.InstancesPathPattern.FindStringSubmatch(ctx.Path)
	if matches == nil {
		return
	}
	applicationName := dubbo.GetProviderApplicationFromMatches(matches)
	ctx.Value = &listProvidersMark{
		applicationName: applicationName,
	}
}

func (providers) HandleResponse(ctx *zookeeper.Context) {
	if ctx.Request == nil {
		return
	}
	var theMark listAndWatchProviderInstances
	var ok bool
	if mark := ctx.Request.Value; mark == nil {
		return
	} else if theMark, ok = mark.(listAndWatchProviderInstances); !ok || !theMark.IsListProviders() {
		return
	}
	content := ctx.RawPayload
	children := zookeeper.ParseChildren(content)
	ctx.Value = children
	applicationName := theMark.ApplicationName()
	dubbo.UpdateEndpointsByApplication(applicationName, children)
}
