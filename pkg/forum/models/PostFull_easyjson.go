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

func easyjson5077f799DecodeGithubComRflbanParkmailDbmsPkgForumModels(in *jlexer.Lexer, out *PostFull) {
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
		case "post":
			if in.IsNull() {
				in.Skip()
				out.Post = nil
			} else {
				if out.Post == nil {
					out.Post = new(Post)
				}
				(*out.Post).UnmarshalEasyJSON(in)
			}
		case "author":
			if in.IsNull() {
				in.Skip()
				out.Author = nil
			} else {
				if out.Author == nil {
					out.Author = new(User)
				}
				(*out.Author).UnmarshalEasyJSON(in)
			}
		case "thread":
			if in.IsNull() {
				in.Skip()
				out.Thread = nil
			} else {
				if out.Thread == nil {
					out.Thread = new(Thread)
				}
				(*out.Thread).UnmarshalEasyJSON(in)
			}
		case "forum":
			if in.IsNull() {
				in.Skip()
				out.Forum = nil
			} else {
				if out.Forum == nil {
					out.Forum = new(Forum)
				}
				easyjson5077f799DecodeGithubComRflbanParkmailDbmsPkgForumModels1(in, out.Forum)
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
func easyjson5077f799EncodeGithubComRflbanParkmailDbmsPkgForumModels(out *jwriter.Writer, in PostFull) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Post != nil {
		const prefix string = ",\"post\":"
		first = false
		out.RawString(prefix[1:])
		(*in.Post).MarshalEasyJSON(out)
	}
	if in.Author != nil {
		const prefix string = ",\"author\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.Author).MarshalEasyJSON(out)
	}
	if in.Thread != nil {
		const prefix string = ",\"thread\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.Thread).MarshalEasyJSON(out)
	}
	if in.Forum != nil {
		const prefix string = ",\"forum\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		easyjson5077f799EncodeGithubComRflbanParkmailDbmsPkgForumModels1(out, *in.Forum)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PostFull) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5077f799EncodeGithubComRflbanParkmailDbmsPkgForumModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostFull) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5077f799EncodeGithubComRflbanParkmailDbmsPkgForumModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostFull) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5077f799DecodeGithubComRflbanParkmailDbmsPkgForumModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostFull) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5077f799DecodeGithubComRflbanParkmailDbmsPkgForumModels(l, v)
}
func easyjson5077f799DecodeGithubComRflbanParkmailDbmsPkgForumModels1(in *jlexer.Lexer, out *Forum) {
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
		case "user":
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
func easyjson5077f799EncodeGithubComRflbanParkmailDbmsPkgForumModels1(out *jwriter.Writer, in Forum) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix[1:])
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"user\":"
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
