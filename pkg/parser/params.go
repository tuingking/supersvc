package parser

import (
	"github.com/gorilla/schema"
)

type ParamParser interface {
	Encode(src interface{}, dest map[string][]string) error
	Decode(dest interface{}, src map[string][]string) error
}

type paramparser struct {
	encoder *schema.Encoder
	decoder *schema.Decoder
}

func InitParamParser() ParamParser {

	p := &paramparser{
		encoder: schema.NewEncoder(),
		decoder: schema.NewDecoder(),
	}

	p.InitDecoder()
	p.InitEncoder()

	p.decoder.SetAliasTag("param")
	p.decoder.IgnoreUnknownKeys(true)
	p.decoder.ZeroEmpty(false)

	p.encoder.SetAliasTag("param")

	return p
}

func (p *paramparser) Encode(src interface{}, dest map[string][]string) error {
	return p.encoder.Encode(src, dest)
}

func (p *paramparser) Decode(dest interface{}, src map[string][]string) error {
	return p.decoder.Decode(dest, src)
}
