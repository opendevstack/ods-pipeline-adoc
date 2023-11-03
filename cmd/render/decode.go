package main

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	"sigs.k8s.io/yaml"
)

// decoder is an interface implemented by *json.Decoder
// and by yamlDecoder. It allows the rendering logic to
// work with both decoders interchangeably.
type decoder interface {
	Decode(v any) error
}

// yamlDecoder implements the decoder interface.
type yamlDecoder struct {
	r io.Reader
}

// Decode decodes data coming from an io.Reader into v.
func (d *yamlDecoder) Decode(v any) error {
	b, err := io.ReadAll(d.r)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, v)
}

type plainDecoder struct {
	r io.Reader
}

func (d *plainDecoder) Decode(v any) error {
	b, err := io.ReadAll(d.r)
	if err != nil {
		return err
	}
	if x, ok := v.(*map[string]interface{}); ok {
		z := *x
		z["value"] = strings.TrimSpace(string(b))
	} else {
		return errors.New("unexpected type")
	}
	return nil
}

// newDecoderFunc creates a decoder for given io.Reader r.
type newDecoderFunc func(r io.Reader) decoder

// selectNewDecoderFunc chooses the decoder to use for each file extension.
func selectNewDecoderFunc(ext string) newDecoderFunc {
	switch ext {
	case ".json":
		return func(r io.Reader) decoder { return json.NewDecoder(r) }
	case ".yaml", ".yml":
		return func(r io.Reader) decoder { return &yamlDecoder{r} }
	}
	return func(r io.Reader) decoder { return &plainDecoder{r} }
}
