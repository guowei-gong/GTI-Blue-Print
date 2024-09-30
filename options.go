package config

import (
	"context"
	"github.com/gti-blue-print/config/encoding/toml"
	"github.com/gti-blue-print/config/errors"
	"strings"
)

type Option func(o *options)

type Encoder func(format string, content interface{}) ([]byte, error)
type Decoder func(format string, content []byte) (interface{}, error)
type Scanner func(format string, content []byte, dest interface{}) error

type options struct {
	ctx     context.Context
	sources []Source // 配置源
	encoder Encoder  // 编码器
	decoder Decoder  // 解码器
	scanner Scanner  // 扫描器
}

func defaultOptions() *options {
	return &options{
		ctx:     context.Background(),
		encoder: defaultEncoder,
		decoder: defaultDecoder,
		scanner: defaultScanner,
	}
}

// 默认编码器
func defaultEncoder(format string, content interface{}) ([]byte, error) {
	switch strings.ToLower(format) {
	case toml.Name:
		return toml.Marshal(content)
	default:
		return nil, errors.ErrInvalidFormat
	}
}

// 默认解码器
func defaultDecoder(format string, content []byte) (interface{}, error) {
	switch strings.ToLower(format) {
	case toml.Name:
		return unmarshal(content, toml.Unmarshal)
	default:
		return nil, errors.ErrInvalidFormat
	}
}

// 默认扫描器
func defaultScanner(format string, content []byte, dest interface{}) error {
	switch strings.ToLower(format) {
	case toml.Name:
		return toml.Unmarshal(content, dest)
	default:
		return errors.ErrInvalidFormat
	}
}

// WithSources 设置配置源
func WithSources(sources ...Source) Option {
	// sources[:]，浅拷贝防止外部修改影响 options.sources 的值
	return func(o *options) { o.sources = sources[:] }
}

// 返回解码后的数据
func unmarshal(content []byte, fn func(data []byte, v interface{}) error) (dest interface{}, err error) {
	dest = make(map[string]interface{})
	if err = fn(content, &dest); err == nil {
		return
	}
	return
}
