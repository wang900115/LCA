package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/allegro/bigcache/v3"
)

type BigCacheMethod interface {
	Get(key string) (any, error)
	Set(key string, value any) error
	Delete(key string) error

	serialize(value any) ([]byte, error)
	deserialize(bytes []byte) (any, error)
}

type BigCache struct {
	bc *bigcache.BigCache
}

func NewBigCache(ctx context.Context, life time.Duration) (BigCacheMethod, error) {
	bc, err := bigcache.New(ctx, bigcache.DefaultConfig(life))
	if err != nil {
		return nil, err
	}
	return &BigCache{
		bc: bc,
	}, nil
}

func (b *BigCache) Get(key string) (any, error) {
	v, err := b.bc.Get(key)
	if err != nil {
		return nil, err
	}
	value, err := b.deserialize(v)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (b *BigCache) Set(key string, value any) error {
	bytes, err := b.serialize(value)
	if err != nil {
		return err
	}
	return b.bc.Set(key, bytes)
}

func (b *BigCache) Delete(key string) error {
	return b.bc.Delete(key)
}

func (b *BigCache) serialize(value any) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	gob.Register(value)

	err := enc.Encode(&value)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (b *BigCache) deserialize(byte []byte) (any, error) {
	var value any
	buf := bytes.NewBuffer(byte)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&value)
	if err != nil {
		return nil, err
	}
	return value, nil
}
