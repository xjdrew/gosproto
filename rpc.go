package sproto

import (
	"errors"
	"reflect"
	"sync"
)

var (
	ErrRepeatedRpc     = errors.New("sproto rpc: repeated rpc")
	ErrUnknownProtocol = errors.New("sproto rpc: unknown protocol")
	ErrUnknownSession  = errors.New("sproto rpc: unknown session")
)

type RpcMode int

const (
	RpcRequestMode RpcMode = iota
	RpcResponseMode
)

type rpcHeader struct {
	Type    *int `sproto:"integer,0,name=type"`
	Session *int `sproto:"integer,1,name=session"`
}

type Protocol struct {
	Type     int
	Name     string
	Request  reflect.Type
	Response reflect.Type
}

type Rpc struct {
	protocols    []*Protocol
	idMap        map[int]int
	nameMap      map[string]int
	sessionMutex sync.Mutex
	sessions     map[int]int
}

func getRpcSprotoType(typ reflect.Type) (*SprotoType, error) {
	if typ == nil {
		return nil, nil
	}

	if typ.Kind() != reflect.Ptr {
		return nil, ErrNonPtr
	}

	return GetSprotoType(typ.Elem())
}

func (rpc *Rpc) Dispatch(packed []byte) (mode RpcMode, name string, session int, sp interface{}, err error) {
	var unpacked []byte
	if unpacked, err = Unpack(packed); err != nil {
		return
	}

	var used int
	header := rpcHeader{}
	if used, err = Decode(unpacked, &header); err != nil {
		return
	}

	var proto *Protocol
	if header.Type != nil {
		index, ok := rpc.idMap[*header.Type]
		if !ok {
			err = ErrUnknownProtocol
			return
		}
		proto = rpc.protocols[index]
		if proto.Request != nil {
			sp = reflect.New(proto.Request.Elem()).Interface()
			if _, err = Decode(unpacked[used:], sp); err != nil {
				return
			}
		}
		mode = RpcRequestMode
		if header.Session != nil {
			session = *header.Session
		}
	} else {
		if header.Session == nil {
			err = ErrUnknownSession
			return
		}
		session = *header.Session
		rpc.sessionMutex.Lock()
		defer rpc.sessionMutex.Unlock()
		index, ok := rpc.sessions[session]
		if !ok {
			err = ErrUnknownSession
			return
		}
		delete(rpc.sessions, session)

		proto = rpc.protocols[index]
		if proto.Response != nil {
			sp = reflect.New(proto.Response.Elem()).Interface()
			if _, err = Decode(unpacked[used:], sp); err != nil {
				return
			}
		}
		mode = RpcResponseMode
	}
	name = proto.Name
	return
}

func (rpc *Rpc) ResponseEncode(name string, session int, response interface{}) (data []byte, err error) {
	index, ok := rpc.nameMap[name]
	if !ok {
		err = ErrUnknownProtocol
		return
	}

	protocol := rpc.protocols[index]
	if protocol.Response != nil {
		if data, err = Encode(response); err != nil {
			return
		}
	}

	header, _ := Encode(&rpcHeader{Session: &session})
	data = Pack(Append(header, data))
	return
}

// session > 0: need response
func (rpc *Rpc) RequestEncode(name string, session int, req interface{}) (data []byte, err error) {
	index, ok := rpc.nameMap[name]
	if !ok {
		err = ErrUnknownProtocol
		return
	}

	protocol := rpc.protocols[index]
	if protocol.Request != nil {
		if data, err = Encode(req); err != nil {
			return
		}
	}

	header, _ := Encode(&rpcHeader{
		Type:    &protocol.Type,
		Session: &session,
	})

	if session > 0 {
		rpc.sessionMutex.Lock()
		defer rpc.sessionMutex.Unlock()
		rpc.sessions[session] = index
	}
	data = Pack(Append(header, data))
	return
}

func NewRpc(protocols []*Protocol) (*Rpc, error) {
	idMap := make(map[int]int)
	nameMap := make(map[string]int)
	for i, protocol := range protocols {
		if _, err := getRpcSprotoType(protocol.Request); err != nil {
			return nil, err
		}
		if _, err := getRpcSprotoType(protocol.Response); err != nil {
			return nil, err
		}
		if _, ok := idMap[protocol.Type]; ok {
			return nil, ErrRepeatedRpc
		}
		if _, ok := nameMap[protocol.Name]; ok {
			return nil, ErrRepeatedRpc
		}
		idMap[protocol.Type] = i
		nameMap[protocol.Name] = i
	}
	rpc := &Rpc{
		protocols: protocols,
		idMap:     idMap,
		nameMap:   nameMap,
		sessions:  make(map[int]int),
	}
	return rpc, nil
}
