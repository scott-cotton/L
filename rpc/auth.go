package rpc

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io"
)

type msg[T any] struct {
	Payload   []byte
	Signature []byte
}

var ErrSignature = errors.New("ErrSignature")

func fromReader[T any](key []byte, r io.Reader) (*T, error) {
	var m msg[T]
	if err := json.NewDecoder(r).Decode(&m); err != nil {
		return nil, err
	}
	return m.Decode(key)
}
func toWriter[T any](key []byte, w io.Writer, v *T) error {
	var (
		m   msg[T]
		err error
	)
	m.Payload, err = json.Marshal(v)
	if err != nil {
		return err
	}
	mac := hmac.New(sha256.New, key)
	mac.Write(m.Payload)
	m.Signature = mac.Sum(nil)
	return json.NewEncoder(w).Encode(m)
}

func (m *msg[T]) Decode(key []byte) (*T, error) {
	mac := hmac.New(sha256.New, key)
	mac.Write(m.Payload)
	sig := mac.Sum(nil)
	var err error
	if !hmac.Equal(sig, m.Signature) {
		err = ErrSignature
	}
	var t T
	if err := json.Unmarshal(m.Payload, &t); err != nil {
		return nil, err
	}
	return &t, err
}
