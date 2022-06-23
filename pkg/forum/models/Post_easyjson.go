// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	time "time"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson6aa74c22DecodeGithubComRflbanParkmailDbmsPkgForumModels(in *jlexer.Lexer, out *Post) {
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
		case "id":
			if in.IsNull() {
				in.Skip()
				out.Id = nil
			} else {
				if out.Id == nil {
					out.Id = new(int64)
				}
				*out.Id = int64(in.Int64())
			}
		case "parent":
			if in.IsNull() {
				in.Skip()
				out.Parent = nil
			} else {
				if out.Parent == nil {
					out.Parent = new(int64)
				}
				*out.Parent = int64(in.Int64())
			}
		case "author":
			out.Author = string(in.String())
		case "message":
			out.Message = string(in.String())
		case "isEdited":
			if in.IsNull() {
				in.Skip()
				out.IsEdited = nil
			} else {
				if out.IsEdited == nil {
					out.IsEdited = new(bool)
				}
				*out.IsEdited = bool(in.Bool())
			}
		case "forum":
			if in.IsNull() {
				in.Skip()
				out.Forum = nil
			} else {
				if out.Forum == nil {
					out.Forum = new(string)
				}
				*out.Forum = string(in.String())
			}
		case "thread":
			if in.IsNull() {
				in.Skip()
				out.Thread = nil
			} else {
				if out.Thread == nil {
					out.Thread = new(int32)
				}
				*out.Thread = int32(in.Int32())
			}
		case "created":
			if in.IsNull() {
				in.Skip()
				out.Created = nil
			} else {
				if out.Created == nil {
					out.Created = new(time.Time)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.Created).UnmarshalJSON(data))
				}
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
func easyjson6aa74c22EncodeGithubComRflbanParkmailDbmsPkgForumModels(out *jwriter.Writer, in Post) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Id != nil {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.Int64(int64(*in.Id))
	}
	if in.Parent != nil {
		const prefix string = ",\"parent\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(*in.Parent))
	}
	{
		const prefix string = ",\"author\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Author))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	if in.IsEdited != nil {
		const prefix string = ",\"isEdited\":"
		out.RawString(prefix)
		out.Bool(bool(*in.IsEdited))
	}
	if in.Forum != nil {
		const prefix string = ",\"forum\":"
		out.RawString(prefix)
		out.String(string(*in.Forum))
	}
	if in.Thread != nil {
		const prefix string = ",\"thread\":"
		out.RawString(prefix)
		out.Int32(int32(*in.Thread))
	}
	if in.Created != nil {
		const prefix string = ",\"created\":"
		out.RawString(prefix)
		out.Raw((*in.Created).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Post) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6aa74c22EncodeGithubComRflbanParkmailDbmsPkgForumModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Post) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6aa74c22EncodeGithubComRflbanParkmailDbmsPkgForumModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Post) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6aa74c22DecodeGithubComRflbanParkmailDbmsPkgForumModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Post) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6aa74c22DecodeGithubComRflbanParkmailDbmsPkgForumModels(l, v)
}
