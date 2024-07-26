package fxcore_test

import (
	"testing"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxcore"
	"github.com/stretchr/testify/assert"
)

func TestNewFxCoreModuleInfo(t *testing.T) {
	t.Setenv("APP_ENV", "test")

	cfg, err := config.NewDefaultConfigFactory().Create(
		config.WithFilePaths("./testdata/config"),
	)
	assert.NoError(t, err)

	info := fxcore.NewFxCoreModuleInfo(
		fxcore.FxCoreModuleInfoParam{
			Config: cfg,
			ExtraInfos: []fxcore.FxExtraInfo{
				fxcore.NewFxExtraInfo("foo", "bar"),
				fxcore.NewFxExtraInfo("foo", "baz"),
			},
		},
	)
	assert.IsType(t, &fxcore.FxCoreModuleInfo{}, info)
	assert.Implements(t, (*fxcore.FxModuleInfo)(nil), info)

	assert.Equal(t, fxcore.ModuleName, info.Name())
	assert.Equal(
		t,
		map[string]interface{}{
			"app": map[string]interface{}{
				"name":        "core-app",
				"description": "core app description",
				"env":         "test",
				"debug":       true,
				"version":     "0.1.0",
			},
			"log": map[string]interface{}{
				"level":  "debug",
				"output": "test",
			},
			"trace": map[string]interface{}{
				"processor": "test",
				"sampler":   "parent-based-always-on",
			},
			"extra": map[string]string{
				"foo": "baz",
			},
		},
		info.Data(),
	)
}
