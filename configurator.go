package config

import (
	"context"
	"github.com/gti-blue-print/config/core/value"
	"github.com/gti-blue-print/config/errors"
	"log"
	"strconv"
	"strings"
	"sync/atomic"
)

type Configurator interface {
	// Get 获取配置值
	Get(pattern string, def ...interface{}) value.Value
	// Set 设置配置值
	Set(pattern string, value interface{}) error
	// Watch 设置监听回调
	Watch(cb WatchCallbackFunc, names ...string)
	// Load 加载配置项
	Load(ctx context.Context, source string, file ...string) ([]*Configuration, error)
	// Store 保存配置项
	Store(ctx context.Context, source string, file string, content interface{}, override ...bool) error
	// Close 关闭配置监听
	Close()
}

type WatchCallbackFunc func(names ...string)

type defaultConfigurator struct {
	opts    *options
	idx     int64 // 搭配 values 字段使用，形成环形缓冲
	ctx     context.Context
	cancel  context.CancelFunc
	sources map[string]Source         // 配置文件源
	values  [2]map[string]interface{} // 配置环形缓冲区
}

func NewConfigurator(opts ...Option) Configurator {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	r := &defaultConfigurator{}
	r.opts = o
	r.ctx, r.cancel = context.WithCancel(o.ctx)
	r.init()
	return r
}

// 初始化配置源
func (c *defaultConfigurator) init() {
	c.sources = make(map[string]Source, len(c.opts.sources))
	for _, s := range c.opts.sources {
		c.sources[s.Name()] = s
	}

	values := make(map[string]interface{})
	for _, s := range c.opts.sources {
		cs, err := s.Load(c.ctx)
		if err != nil {
			log.Printf("load configure failed: %v", err)
			continue
		}

		for _, cc := range cs {
			if len(cc.Content) == 0 {
				continue
			}

			v, err := c.opts.decoder(cc.Format, cc.Content)
			if err != nil {
				if err != errors.ErrInvalidFormat {
					log.Printf("decode configure failed: %v", err)
				}
				continue
			}

			values[cc.Name] = v
		}
	}

	c.store(values)
}

// 保存配置，循环存储机制，len(c.values) 定长，确保 idx 在 0 和 1 之间循环
func (c *defaultConfigurator) store(values map[string]interface{}) {
	idx := atomic.AddInt64(&c.idx, 1) % int64(len(c.values))
	c.values[idx] = values
}

// Get 获取配置值
func (c *defaultConfigurator) Get(pattern string, def ...interface{}) value.Value {
	if val, ok := c.doGet(pattern); ok {
		return val
	}

	return value.NewValue(def...)
}

// 执行获取配置操作
func (c *defaultConfigurator) doGet(pattern string) (value.Value, bool) {
	var (
		keys   = strings.Split(pattern, ".")
		node   interface{}
		found  = true
		values = c.load()
	)

	if values == nil || len(values) == 0 {
		goto NOTFOUND
	}

	keys = reviseKeys(keys, values)
	node = values
	for _, key := range keys {
		switch vs := node.(type) {
		case map[string]interface{}:
			if v, ok := vs[key]; ok {
				node = v
			} else {
				found = false
			}
		case []interface{}:
			i, err := strconv.Atoi(key)
			if err != nil {
				found = false
			} else if len(vs) > i {
				node = vs[i]
			} else {
				found = false
			}
		default:
			found = false
		}

		if !found {
			break
		}
	}

	if found {
		return value.NewValue(node), true
	}

NOTFOUND:
	return nil, false
}

// 重组 keys
// e.g: 如果输入 keys 是 ["a", "b", "c"]，而 values map 中存在 "a.b" 这个键，那么函数会将 keys 修改为 ["a.b", "c"]
func reviseKeys(keys []string, values map[string]interface{}) []string {
	for i := 1; i < len(keys); i++ {
		key := strings.Join(keys[:i+1], ".")
		if _, ok := values[key]; ok {
			keys[0] = key
			temp := keys[i+1:]
			copy(keys[1:], temp)
			keys = keys[:len(temp)+1]
			break
		}
	}

	return keys
}

func (c *defaultConfigurator) Set(pattern string, value interface{}) error {
	return nil
}

func (c *defaultConfigurator) Watch(cb WatchCallbackFunc, names ...string) {
	return
}

func (c *defaultConfigurator) Load(ctx context.Context, source string, file ...string) ([]*Configuration, error) {
	return nil, nil
}

// 加载配置
func (c *defaultConfigurator) load() map[string]interface{} {
	idx := atomic.LoadInt64(&c.idx) % int64(len(c.values))
	return c.values[idx]
}

func (c *defaultConfigurator) Store(ctx context.Context, source string, file string, content interface{}, override ...bool) error {
	return nil
}

func (c *defaultConfigurator) Close() {
	return
}
