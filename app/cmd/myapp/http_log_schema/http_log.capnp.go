// Code generated by capnpc-go. DO NOT EDIT.

package http_log_schema

import (
	capnp "capnproto.org/go/capnp/v3"
	text "capnproto.org/go/capnp/v3/encoding/text"
	schemas "capnproto.org/go/capnp/v3/schemas"
)

type HttpLogRecord capnp.Struct

// HttpLogRecord_TypeID is the unique identifier for the type HttpLogRecord.
const HttpLogRecord_TypeID = 0xc5c66f22213ada4b

func NewHttpLogRecord(s *capnp.Segment) (HttpLogRecord, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 40, PointerCount: 4})
	return HttpLogRecord(st), err
}

func NewRootHttpLogRecord(s *capnp.Segment) (HttpLogRecord, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 40, PointerCount: 4})
	return HttpLogRecord(st), err
}

func ReadRootHttpLogRecord(msg *capnp.Message) (HttpLogRecord, error) {
	root, err := msg.Root()
	return HttpLogRecord(root.Struct()), err
}

func (s HttpLogRecord) String() string {
	str, _ := text.Marshal(0xc5c66f22213ada4b, capnp.Struct(s))
	return str
}

func (s HttpLogRecord) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Struct(s).EncodeAsPtr(seg)
}

func (HttpLogRecord) DecodeFromPtr(p capnp.Ptr) HttpLogRecord {
	return HttpLogRecord(capnp.Struct{}.DecodeFromPtr(p))
}

func (s HttpLogRecord) ToPtr() capnp.Ptr {
	return capnp.Struct(s).ToPtr()
}
func (s HttpLogRecord) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s HttpLogRecord) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s HttpLogRecord) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s HttpLogRecord) TimestampEpochMilli() uint64 {
	return capnp.Struct(s).Uint64(0)
}

func (s HttpLogRecord) SetTimestampEpochMilli(v uint64) {
	capnp.Struct(s).SetUint64(0, v)
}

func (s HttpLogRecord) ResourceId() uint64 {
	return capnp.Struct(s).Uint64(8)
}

func (s HttpLogRecord) SetResourceId(v uint64) {
	capnp.Struct(s).SetUint64(8, v)
}

func (s HttpLogRecord) BytesSent() uint64 {
	return capnp.Struct(s).Uint64(16)
}

func (s HttpLogRecord) SetBytesSent(v uint64) {
	capnp.Struct(s).SetUint64(16, v)
}

func (s HttpLogRecord) RequestTimeMilli() uint64 {
	return capnp.Struct(s).Uint64(24)
}

func (s HttpLogRecord) SetRequestTimeMilli(v uint64) {
	capnp.Struct(s).SetUint64(24, v)
}

func (s HttpLogRecord) ResponseStatus() uint16 {
	return capnp.Struct(s).Uint16(32)
}

func (s HttpLogRecord) SetResponseStatus(v uint16) {
	capnp.Struct(s).SetUint16(32, v)
}

func (s HttpLogRecord) CacheStatus() (string, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return p.Text(), err
}

func (s HttpLogRecord) HasCacheStatus() bool {
	return capnp.Struct(s).HasPtr(0)
}

func (s HttpLogRecord) CacheStatusBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return p.TextBytes(), err
}

func (s HttpLogRecord) SetCacheStatus(v string) error {
	return capnp.Struct(s).SetText(0, v)
}

func (s HttpLogRecord) Method() (string, error) {
	p, err := capnp.Struct(s).Ptr(1)
	return p.Text(), err
}

func (s HttpLogRecord) HasMethod() bool {
	return capnp.Struct(s).HasPtr(1)
}

func (s HttpLogRecord) MethodBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(1)
	return p.TextBytes(), err
}

func (s HttpLogRecord) SetMethod(v string) error {
	return capnp.Struct(s).SetText(1, v)
}

func (s HttpLogRecord) RemoteAddr() (string, error) {
	p, err := capnp.Struct(s).Ptr(2)
	return p.Text(), err
}

func (s HttpLogRecord) HasRemoteAddr() bool {
	return capnp.Struct(s).HasPtr(2)
}

func (s HttpLogRecord) RemoteAddrBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(2)
	return p.TextBytes(), err
}

func (s HttpLogRecord) SetRemoteAddr(v string) error {
	return capnp.Struct(s).SetText(2, v)
}

func (s HttpLogRecord) Url() (string, error) {
	p, err := capnp.Struct(s).Ptr(3)
	return p.Text(), err
}

func (s HttpLogRecord) HasUrl() bool {
	return capnp.Struct(s).HasPtr(3)
}

func (s HttpLogRecord) UrlBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(3)
	return p.TextBytes(), err
}

func (s HttpLogRecord) SetUrl(v string) error {
	return capnp.Struct(s).SetText(3, v)
}

// HttpLogRecord_List is a list of HttpLogRecord.
type HttpLogRecord_List = capnp.StructList[HttpLogRecord]

// NewHttpLogRecord creates a new list of HttpLogRecord.
func NewHttpLogRecord_List(s *capnp.Segment, sz int32) (HttpLogRecord_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 40, PointerCount: 4}, sz)
	return capnp.StructList[HttpLogRecord](l), err
}

// HttpLogRecord_Future is a wrapper for a HttpLogRecord promised by a client call.
type HttpLogRecord_Future struct{ *capnp.Future }

func (f HttpLogRecord_Future) Struct() (HttpLogRecord, error) {
	p, err := f.Future.Ptr()
	return HttpLogRecord(p.Struct()), err
}

const schema_f42cd342ff520eca = "x\xdaL\xcc1\x8b\xd4@\x18\xc6\xf1\xe7\x99I6\xb7" +
	"\xb0z\x172\x82\x85\xc2y(\x9c\xb0\x0a\x0b\x16\xb2\xcd" +
	"\xa9 (z`v\xac\xae\x91\\2\xdc\x06\x92\x9d\x98" +
	"\xcc\x16\xfa\x11\xfc\x18~\x11\xab\xb3\x10,\xe4:?\x82" +
	" XX(\x08#S\xec\xb1\xdd\xfb\xfe\xf8\xf3\xec]" +
	"<\x8afW>\x11\"W\xf1\xc8\xbf\xf8>\xbfu`" +
	"?\x9f#\xdfe\xec\xbf\\]\xf8'\x17\xd3\xdf\x88\xa3" +
	"\x04H\xcf\x7f\xa4\xdf\x12`\xf6\xd5\x13\xf7\xfc\xd2\xb9\xee" +
	"Mc\xcf\xc4\xfd\xb2\xe8V\xdd\xfc\x99s\xddK{\xb6" +
	"0\xbb\xa5\xed\xabWd~[F@D \xfd\xf9\x11" +
	"\xc8\x7fI\xe6\xff\x04SR1\xe0\xdf\x13 \xff#\xa9" +
	"#\x0a\xa6B(\x0a #\x17\xc0\x82\x92z\x12XJ" +
	"E\x09dc~\x00\xf4$\xf8\xf5\xe0\xd1\x9eb\x04d" +
	"\xd7\xf8\x1e\xd0*\xf8>\x05\x19+\xc6@v\x93\xa7\x80" +
	"\xbe\x11\xf80\xe4#*\x8e\x80\xec\x0e\xe7\x80\xde\x0f>" +
	"\x0d\x9e\x08\xc5\x04\xc8\xee\xf2\x04\xd0\x87\xc1\x1f\x04\xdf\x91" +
	"\x8a;@6\xe3\x01\xa0\xa7\xc1\x1fR\xd0\xbb\xba5\x83" +
	"+ZvO;[.\x8f\xeb\xa4ij\x8e!8\x06" +
	"}o\x06\xbb\xeeK\x03\xf9\xbc\xba\xc4\xd3w\xce\x0c\xda" +
	"\xac@\xb7\x15\xbe]\x9b\xc1\xbdf\xdd\x9a\xe3\xbaij" +
	"`{\xa4\xb3\xab\xc1\xe0H\xbb\xc2\xad\x07&\x10L@" +
	"_\x16\xe5\xd2hW \x09:\x81\xe0\x04<j\x8d[" +
	"\xdaj\xf3\xfa\xde\xb4\xd6\x99\xc7\x15d\xd5o0Y\xf7" +
	"\xcd\xe6\xfe\x1f\x00\x00\xff\xff\xf63b\x8c"

func RegisterSchema(reg *schemas.Registry) {
	reg.Register(&schemas.Schema{
		String: schema_f42cd342ff520eca,
		Nodes: []uint64{
			0xc5c66f22213ada4b,
		},
		Compressed: true,
	})
}