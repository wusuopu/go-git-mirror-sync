package helper

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/valyala/fastjson"
)

var pPool fastjson.ParserPool
var aPool fastjson.ArenaPool
/*
	动态解析 body 为 json
*/
func GetJSONBody(c *gin.Context) *fastjson.Value  {
	a := aPool.Get()
	p := pPool.Get()
	defer func ()  {
		aPool.Put(a)
		pPool.Put(p)
	}()

	body, ok := c.Get("body")
	if ok {
		return body.(*fastjson.Value)
	}
	empty := a.NewObject()

	rawBody, _ := c.Get("rawBody")
	if rawBody == nil {
		return empty
	}
	maxLength := 1024 * 1024 * 2		// 最大支持 2M 长度字符，以免内存不足 
	data, _ := rawBody.([]byte)
	ret, err := p.ParseBytes(data[:min(len(data), maxLength)])
	if err != nil {
		ret = empty
	}
	// 将解析结果缓存下来
	c.Set("body", ret)

	return ret
}
/*
	动态解析 url query 为 json
*/
func GetJSONQuery(c *gin.Context) *fastjson.Value {
	body, ok := c.Get("query")
	if ok {
		return body.(*fastjson.Value)
	}

	obj := ParseJSONQuery(c.Request.URL.RawQuery)
	c.Set("query", obj)
	return obj
}
func ParseJSONQuery(qs string) *fastjson.Value {
	a := aPool.Get()
	defer func ()  {
		aPool.Put(a)
	}()

	obj := a.NewObject()
	query, err := url.ParseQuery(qs)
	if err != nil {
		return obj
	}
	for k, v := range query {
		if len(v) == 1 && len(v[0]) == 0 {
			//该字段为空
			continue
		}

		isArray := false
		isNested := false
		if strings.HasSuffix(k, "[]") {
			isArray = true
			k = k[:(len(k)-2)]
		}
		if strings.HasSuffix(k, "]") {
			// 嵌套结构
			isNested = true
		}
		var val *fastjson.Value
		if isArray {
			val = a.NewArray()
			for idx, el := range v {
				val.SetArrayItem(idx, a.NewString(el))
			}
		} else {
			val = a.NewString(v[len(v) - 1])
		}

		if !isNested {
			// 单层结构，直接赋值
			obj.Set(k, val)
			continue
		}

		parent := obj
		keys := strings.Split(k, "[")
		for idx, key := range keys {
			if strings.HasSuffix(key, "]") {
				key = key[:len(key)-1]
			}
			if key == "" {
				continue
			}

			if idx == (len(keys) - 1) {
				parent.Set(key, val)
			} else {
				if !parent.Exists(key) {
					subItem := a.NewObject()
					parent.Set(key, subItem)

				}
				parent = parent.Get(key)
			}
		}
	}

	return obj
}


// lodash 风格的 GET 操作
func GetJSONInt64(val *fastjson.Value, keys string) (int64, error) {
	k := strings.Split(keys, ".")
	if !val.Exists(k...) {
		return 0, fmt.Errorf("key %s not exists", keys)
	}
	item := val.Get(k...)
	if item.Type() == fastjson.TypeNumber {
		return item.Int64()
	}
	if item.Type() == fastjson.TypeString {
	 strconv.ParseInt(item.String(), 10, 64)
	}
	return 0, fmt.Errorf("key %s in not number", keys)
}
func GetJSONFloat64(val *fastjson.Value, keys string) (float64, error) {
	k := strings.Split(keys, ".")
	if !val.Exists(k...) {
		return 0, fmt.Errorf("key %s not exists", keys)
	}
	item := val.Get(k...)
	if item.Type() == fastjson.TypeNumber {
		return item.Float64()
	}
	if item.Type() == fastjson.TypeString {
	 strconv.ParseFloat(item.String(), 64)
	}
	return 0, fmt.Errorf("key %s in not number", keys)
}
func GetJSONItem(val *fastjson.Value, keys string) *fastjson.Value {
	return val.Get(strings.Split(keys, ".")...)
}
func GetJSONString(val *fastjson.Value, keys string) string {
	v := GetJSONItem(val, keys)
	if v.Type() != fastjson.TypeString {
		return ""
	}

	ret := v.String()
	if len(ret) < 2 { return ret }
	if ret[0] == '"' && ret[len(ret)-1] == '"' {
			// 去掉前后的双引号
			return ret[1:len(ret)-1]
	}
	return ret
}