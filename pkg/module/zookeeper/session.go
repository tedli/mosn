package zookeeper

import "encoding/binary"

func NewSessionFilter() Filter {
	return new(session)
}

type session struct{}

func (session) HandleRequest(ctx *Context) {
	ctx.session.Store(ctx.Xid, ctx)
}

func (session) HandleResponse(ctx *Context) {
	request, loaded := ctx.session.LoadAndDelete(ctx.Xid)
	if !loaded {
		return
	}
	ctx.Request = request.(*Context)
	ctx.Watch = ctx.Request.Watch
}

const (
	childrenCountEnd = 6 * Uint32Size
)

func ParseChildren(content []byte) (children []string) {
	childrenCount := int(binary.BigEndian.Uint32(content[5*Uint32Size : childrenCountEnd]))
	children = make([]string, childrenCount)
	var end int
	for i, begin := 0, childrenCountEnd; i < childrenCount; i++ {
		end = begin + Uint32Size
		pathLength := int(binary.BigEndian.Uint32(content[begin:end]))
		begin = end
		end = begin + pathLength
		path := string(content[begin:end])
		children[i] = path
		begin = end
	}
	return
}
