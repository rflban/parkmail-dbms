// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson9889753aDecodeGithubComRflbanParkmailDbmsPkgForumModels(in *jlexer.Lexer, out *Status) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "user":
			out.User = int32(in.Int32())
		case "forum":
			out.Forum = int32(in.Int32())
		case "thread":
			out.Thread = int32(in.Int32())
		case "post":
			out.Post = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson9889753aEncodeGithubComRflbanParkmailDbmsPkgForumModels(out *jwriter.Writer, in Status) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"user\":"
		out.RawString(prefix[1:])
		out.Int32(int32(in.User))
	}
	{
		const prefix string = ",\"forum\":"
		out.RawString(prefix)
		out.Int32(int32(in.Forum))
	}
	{
		const prefix string = ",\"thread\":"
		out.RawString(prefix)
		out.Int32(int32(in.Thread))
	}
	{
		const prefix string = ",\"post\":"
		out.RawString(prefix)
		out.Int64(int64(in.Post))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Status) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9889753aEncodeGithubComRflbanParkmailDbmsPkgForumModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Status) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9889753aEncodeGithubComRflbanParkmailDbmsPkgForumModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Status) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9889753aDecodeGithubComRflbanParkmailDbmsPkgForumModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Status) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9889753aDecodeGithubComRflbanParkmailDbmsPkgForumModels(l, v)
}
