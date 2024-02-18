package helper_test

import (
	"app/utils/helper"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fastjson"
)

func Test_GetJSONValue(t *testing.T) {
		var p fastjson.Parser
		obj, _ := p.Parse(`{"a": {"b1": {"c1": 1, "c2": 1.2}, "b2": "demo\nnewline"}}`)

		// 获取整数
		ret, _ := helper.GetJSONInt64(obj, "a.b1.c1")
		assert.Equal(t, ret, int64(1))

		ret, err := helper.GetJSONInt64(obj, "a.b1.c2")
		assert.NotNil(t, err)

		ret2, err := helper.GetJSONFloat64(obj, "a.b1.c2")
		assert.Equal(t, ret2, 1.2)

		// 测试字符串的转义
		ret3 := helper.GetJSONString(obj, "a.b2")
		assert.Equal(t, ret3, "demo\nnewline")
}

func Test_ParseJSONQuery(t *testing.T) {
	qs := "names[0]=1&names[1]=2&key=a&key=b&ids[]=4&ids[]=5&ids[]=6&array-Item[0][name][op]=eq&array-Item[0][name][value]=12&array-Item[1][name][op]=gt&array-Item[1][name][value]=15&dict[name]=foo&dict[age]=18&dict[ids][]=11&dict[ids][]=12&nullable="
	// {"names":{"1":"2","0":"1"},"ids":["4","5","6"],"dict":{"age":"18","ids":["11","12"],"name":"foo"},"array-Item":{"0":{"name":{"op":"eq","value":"12"}},"1":{"name":{"op":"gt","value":"15"}}},"key":"b"}
	ret := helper.ParseJSONQuery(qs)

	assert.Equal(t, helper.GetJSONString(ret, "key"), "b")
	assert.Equal(t, ret.Exists(("nullable")), false)		// 空字段不解析
	assert.Equal(t, string(ret.Get("ids").MarshalTo(nil)), `["4","5","6"]`)

	assert.Equal(t, helper.GetJSONString(ret, "names.0"), "1")
	assert.Equal(t, helper.GetJSONString(ret, "names.1"), "2")

	assert.Equal(t, helper.GetJSONString(ret, "array-Item.0.name.op"), "eq")
	assert.Equal(t, helper.GetJSONString(ret, "array-Item.1.name.op"), "gt")
}