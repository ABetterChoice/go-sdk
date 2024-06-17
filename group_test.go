// Package abc provides a set of APIs for external use, including APIs for ABC system initialization.
// It also encompasses functionalities such as traffic distribution for A/B experiments, user configuration data retrieval,
// user feature flag management, exposure data reporting, and logger registration.
package abc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroup_ParamAPI(t *testing.T) {
	var param = map[string]string{
		"float64": "10.98",
		"int64":   "991119",
		"string":  "stringValue",
		"bool":    "1",
		"jsonMap": `{"key1": "value1"}`,
	}
	g := &Group{
		params:      param,
		sceneIDList: nil,
	}
	t.Run("normal", func(t *testing.T) {
		assert.Equal(t, param, g.Params())
		assert.Nil(t, g.SceneIDList())
		g.sceneIDList = []int64{1, 2, 3}
		assert.Equal(t, []int64{1, 2, 3}, g.SceneIDList())
		assert.Equal(t, 10.98, g.GetFloat64WithDefault("float64", 10))
		assert.Equal(t, float64(10), g.GetFloat64WithDefault("string", 10))
		assert.Equal(t, 10.98, g.MustGetFloat64("float64"))
		float64Value, err := g.GetFloat64("float64")
		assert.Nil(t, err)
		assert.Equal(t, float64Value, 10.98)
		float64Value, err = g.GetFloat64("string")
		assert.NotNil(t, err)
		assert.Equal(t, float64Value, float64(0))
		float64Value, err = g.GetFloat64("empty")
		assert.NotNil(t, err)
		assert.Equal(t, float64Value, float64(0))
		assert.Equal(t, int64(991119), g.GetInt64WithDefault("int64", int64(95)))
		assert.Equal(t, int64(95), g.GetInt64WithDefault("string", int64(95)))
		assert.Equal(t, int64(991119), g.MustGetInt64("int64"))
		int64Value, err := g.GetInt64("int64")
		assert.Nil(t, err)
		assert.Equal(t, int64Value, int64(991119))
		int64Value, err = g.GetInt64("empty")
		assert.NotNil(t, err)
		assert.Equal(t, int64Value, int64(0))
		assert.Equal(t, "stringValue", g.GetStringWithDefault("string", "defaultString"))
		assert.Equal(t, "", g.MustGetString("empty"))
		assert.Equal(t, "xx", g.GetStringWithDefault("empty", "xx"))
		assert.Equal(t, true, g.GetBoolWithDefault("bool", false))
		assert.Equal(t, false, g.GetBoolWithDefault("int64", false))
		boolValue, err := g.GetBool("empty")
		assert.NotNil(t, err)
		assert.False(t, boolValue)
		assert.True(t, g.MustGetBool("bool"))
		assert.Equal(t, map[string]interface{}{"key1": "value1"}, g.GetJSONMapWithDefault("jsonMap", nil))
		assert.Nil(t, g.MustGetJSONMap("empty"))
		assert.Nil(t, g.MustGetJSONMap("string"))
		assert.Nil(t, g.GetJSONMapWithDefault("string", nil))
	})
}
