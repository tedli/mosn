package filters

import (
	"encoding/binary"

	"mosn.io/mosn/pkg/module/zookeeper"
	"mosn.io/pkg/log"
)

func NewDataFilter() zookeeper.Filter {
	return new(data)
}

type data struct{}

func (data) HandleRequest(ctx *zookeeper.Context) {
	if ctx.OpCode == zookeeper.Undefined || len(ctx.Path) <= 0 {
		return
	}
	switch ctx.OpCode {
	case zookeeper.OpCreate, zookeeper.OpSetData:
		handleRequestData(ctx)
	default:
	}
}

func handleRequestData(ctx *zookeeper.Context) {
	pathLength := len(ctx.Path)
	begin := 4 * zookeeper.Uint32Size + pathLength
	handleDataInternal(ctx, begin)
}

func (data) HandleResponse(ctx *zookeeper.Context) {
	if ctx.Error != nil {
		return
	}
	var request *zookeeper.Context
	if request = ctx.Request; request == nil {
		return
	}
	if request.OpCode == zookeeper.Undefined || len(request.Path) <= 0 {
		return
	}
	switch request.OpCode {
	case zookeeper.OpGetData:
		handleResponseData(ctx)
	default:
	}
}

func handleResponseData(ctx *zookeeper.Context) {
	handleDataInternal(ctx, 5*zookeeper.Uint32Size)
}

func handleDataInternal(ctx *zookeeper.Context, begin int) {
	buffer := ctx.RawPayload
	length := len(buffer)
	if length < begin {
		log.DefaultLogger.Errorf("zookeeper.filters.path.handleRequestData, buffer too short")
		return
	}
	dataLenght := int(binary.BigEndian.Uint32(buffer[begin : begin + zookeeper.Uint32Size]))
	if begin+zookeeper.Uint32Size+dataLenght > length {
		log.DefaultLogger.Errorf("zookeeper.filters.path.handleRequestData, buffer too short for data")
		return
	}
	ctx.DataBegin = begin+zookeeper.Uint32Size
	ctx.DataEnd = begin+zookeeper.Uint32Size+dataLenght
	ctx.Data = buffer[ctx.DataBegin : ctx.DataEnd]
}
