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

func easyjsonC49b5281DecodeGithubComRflbanParkmailDbmsPkgForumModels(in *jlexer.Lexer, out *Forum) {
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
		case "title":
			out.Title = string(in.String())
		case "username":
			out.User = string(in.String())
		case "slug":
			out.Slug = string(in.String())
		case "posts":
			if in.IsNull() {
				in.Skip()
				out.Posts = nil
			} else {
				if out.Posts == nil {
					out.Posts = new(int64)
				}
				*out.Posts = int64(in.Int64())
			}
		case "threads":
			if in.IsNull() {
				in.Skip()
				out.Threads = nil
			} else {
				if out.Threads == nil {
					out.Threads = new(int32)
				}
				*out.Threads = int32(in.Int32())
			}
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
func easyjsonC49b5281EncodeGithubComRflbanParkmailDbmsPkgForumModels(out *jwriter.Writer, in Forum) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix[1:])
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"username\":"
		out.RawString(prefix)
		out.String(string(in.User))
	}
	{
		const prefix string = ",\"slug\":"
		out.RawString(prefix)
		out.String(string(in.Slug))
	}
	if in.Posts != nil {
		const prefix string = ",\"posts\":"
		out.RawString(prefix)
		out.Int64(int64(*in.Posts))
	}
	if in.Threads != nil {
		const prefix string = ",\"threads\":"
		out.RawString(prefix)
		out.Int32(int32(*in.Threads))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Forum) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC49b5281EncodeGithubComRflbanParkmailDbmsPkgForumModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Forum) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC49b5281EncodeGithubComRflbanParkmailDbmsPkgForumModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Forum) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC49b5281DecodeGithubComRflbanParkmailDbmsPkgForumModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Forum) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC49b5281DecodeGithubComRflbanParkmailDbmsPkgForumModels(l, v)
}
