package filters

import (
	"bytes"
	"encoding/binary"

	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/module/zookeeper"
)

func NewHeaderFilter() zookeeper.Filter {
	return new(header)
}

type header struct {}

const (
	requestHeaderLength = 3 * zookeeper.Uint32Size
	responseHeaderLength = 5 * zookeeper.Uint32Size
)

func (header) HandleRequest(ctx *zookeeper.Context) {
	buffer := ctx.RawPayload
	if length := len(buffer); length < requestHeaderLength {
		log.DefaultLogger.Errorf("zookeeper.filters.requestHeader.Invoke, content too short, length: %d", length)
		return
	}
	var xid int32
	if err := binary.Read(
		bytes.NewBuffer(buffer[zookeeper.Uint32Size: zookeeper.Uint32Size + zookeeper.Uint32Size]),
		binary.BigEndian,
		&xid); err != nil {

		log.DefaultLogger.Errorf("zookeeper.filters.requestHeader.Invoke, decode xid failed, %s", err)
	} else {
		ctx.Xid = int(xid)
	}

	ctx.OpCode = zookeeper.OpCode(binary.BigEndian.Uint32(
		buffer[zookeeper.Uint32Size + zookeeper.Uint32Size: requestHeaderLength]))
}

func (header) HandleResponse(ctx *zookeeper.Context) {
	buffer := ctx.RawPayload
	if length := len(buffer); length < responseHeaderLength {
		log.DefaultLogger.Errorf("zookeeper.filters.responseHeader.Invoke, content too short, length: %d", length)
		return
	}

	var xid int32
	if err := binary.Read(
		bytes.NewReader(buffer[zookeeper.Uint32Size:zookeeper.Uint32Size+zookeeper.Uint32Size]),
		binary.BigEndian,
		&xid); err != nil {

		log.DefaultLogger.Errorf("zookeeper.filters.responseHeader.Invoke, decode xid failed, %s", err)
	} else {
		ctx.Xid = int(xid)
	}

	if err := binary.Read(
		bytes.NewReader(buffer[zookeeper.Uint32Size+zookeeper.Uint32Size:4*zookeeper.Uint32Size]),
		binary.BigEndian,
		&(ctx.Zxid)); err != nil {

		log.DefaultLogger.Errorf("zookeeper.filters.responseHeader.Invoke, decode zxid failed, %s", err)
	}

	var epoch, counter, errCode int32

	if log.DefaultLogger.GetLogLevel() >= log.TRACE {
		if err := binary.Read(
			bytes.NewReader(buffer[zookeeper.Uint32Size+zookeeper.Uint32Size:3*zookeeper.Uint32Size]),
			binary.BigEndian,
			&epoch); err != nil {

			log.DefaultLogger.Errorf("zookeeper.filters.responseHeader.Invoke, decode epoch failed, %s", err)
		}
		if err := binary.Read(
			bytes.NewReader(buffer[3*zookeeper.Uint32Size:4*zookeeper.Uint32Size]),
			binary.BigEndian,
			&counter); err != nil {

			log.DefaultLogger.Errorf("zookeeper.filters.responseHeader.Invoke, decode counter failed, %s", err)
		}
		log.DefaultLogger.Tracef("zookeeper.filters.responseHeader.Invoke, Xid: %d, Zxid: %d, Epoch: %d, Counter: %d", ctx.Xid, ctx.Zxid, epoch, counter)
	}

	if err := binary.Read(
		bytes.NewReader(buffer[4*zookeeper.Uint32Size:responseHeaderLength]),
		binary.BigEndian,
		&errCode); err != nil {

		log.DefaultLogger.Errorf("zookeeper.filters.responseHeader.Invoke, decode error code failed, %s", err)
	}

	ctx.Error = zookeeper.ParseErrorCode(int(errCode))
}
