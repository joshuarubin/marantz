// Code generated by protoc-gen-gogo.
// source: cmds.proto
// DO NOT EDIT!

package msg

import proto "code.google.com/p/gogoprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type Cmds struct {
	Cmd              []*Cmd `protobuf:"bytes,1,rep,name=cmd" json:"cmd,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *Cmds) Reset()         { *m = Cmds{} }
func (m *Cmds) String() string { return proto.CompactTextString(m) }
func (*Cmds) ProtoMessage()    {}

func (m *Cmds) GetCmd() []*Cmd {
	if m != nil {
		return m.Cmd
	}
	return nil
}

func init() {
}