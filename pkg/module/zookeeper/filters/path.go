package filters

import (
	"encoding/binary"

	"mosn.io/mosn/pkg/module/zookeeper"
	"mosn.io/pkg/log"
)

func NewPathFilter() zookeeper.Filter {
	return new(path)
}

type path struct {
}

func (p path) HandleRequest(ctx *zookeeper.Context) {
	if ctx.OpCode == zookeeper.Undefined {
		return
	}
	switch ctx.OpCode {
	case zookeeper.OpGetChildren, zookeeper.OpGetAcl, zookeeper.OpSync:
	case zookeeper.OpCreate, zookeeper.OpCreateTTL, zookeeper.OpDelete, zookeeper.OpSetData, zookeeper.OpSetAcl:
		handlePath(ctx)
	case zookeeper.OpGetChildren2, zookeeper.OpExists, zookeeper.OpGetData:
		handlePathAndWatch(ctx)
	default:
	}
}

const (
	pathEndIndex = 4 * zookeeper.Uint32Size
)

func handlePath(ctx *zookeeper.Context) {
	handlePathInternal(ctx)
}

func handlePathAndWatch(ctx *zookeeper.Context) {
	failed, index, buffer := handlePathInternal(ctx)
	if failed {
		return
	}
	if len(buffer) < index {
		log.DefaultLogger.Errorf("zookeeper.filters.path.handlePathAndWatch, buffer too short")
		return
	}
	ctx.Watch = buffer[index] != 0
}

func handlePathInternal(ctx *zookeeper.Context) (failed bool, pathEnd int, buffer []byte) {
	buffer = ctx.RawPayload
	length := len(buffer)
	if length < pathEndIndex {
		failed = true
		log.DefaultLogger.Errorf("zookeeper.filters.path.handlePath, buffer too short")
		return
	}
	pathLenght := int(binary.BigEndian.Uint32(buffer[3*zookeeper.Uint32Size : pathEndIndex]))
	if pathEndIndex+pathLenght > length {
		failed = true
		log.DefaultLogger.Errorf("zookeeper.filters.path.handlePath, buffer too short for path length")
		return
	}
	ctx.PathBegin = pathEndIndex
	pathEnd = pathEndIndex + pathLenght
	ctx.PathEnd = pathEnd
	ctx.Path = string(buffer[ctx.PathBegin:ctx.PathEnd])
	ctx.OriginalPath = ctx.Path
	return
}

func (path) HandleResponse(ctx *zookeeper.Context) {
}
