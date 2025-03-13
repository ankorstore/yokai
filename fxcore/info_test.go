package fxcore_test

import (
	"testing"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxcore"
	"github.com/stretchr/testify/assert"
)

type testModuleInfo struct{}

func (i *testModuleInfo) Name() string {
	return "test"
}

func (i *testModuleInfo) Data() map[string]interface{} {
	return map[string]interface{}{}
}

func TestFxCoreModuleInfo(t *testing.T) {
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

func TestFxModuleInfoRegistry(t *testing.T) {
	t.Parallel()

	createRegistry := func(tb testing.TB) *fxcore.FxModuleInfoRegistry {
		tb.Helper()

		cfg, err := config.NewDefaultConfigFactory().Create(
			config.WithFilePaths("./testdata/config"),
		)
		assert.NoError(tb, err)

		return fxcore.NewFxModuleInfoRegistry(fxcore.FxModuleInfoRegistryParam{
			Infos: []interface{}{
				&testModuleInfo{},
				fxcore.NewFxCoreModuleInfo(fxcore.FxCoreModuleInfoParam{
					Config:     cfg,
					ExtraInfos: []fxcore.FxExtraInfo{},
				}),
				"invalid",
			},
		})
	}

	t.Run("test type", func(t *testing.T) {
		t.Parallel()

		registry := createRegistry(t)

		assert.IsType(t, &fxcore.FxModuleInfoRegistry{}, registry)
	})

	t.Run("test all", func(t *testing.T) {
		t.Parallel()

		registry := createRegistry(t)

		assert.Len(t, registry.All(), 2)
	})

	t.Run("test names", func(t *testing.T) {
		t.Parallel()

		registry := createRegistry(t)

		assert.Equal(t, []string{fxcore.ModuleName, "test"}, registry.Names())
	})

	t.Run("test find", func(t *testing.T) {
		t.Parallel()

		registry := createRegistry(t)

		testInfo, err := registry.Find("test")
		assert.NoError(t, err)
		assert.Equal(t, "test", testInfo.Name())

		coreInfo, err := registry.Find(fxcore.ModuleName)
		assert.NoError(t, err)
		assert.Equal(t, fxcore.ModuleName, coreInfo.Name())

		invalidInfo, err := registry.Find("invalid")
		assert.Error(t, err)
		assert.Equal(t, "fx module info with name invalid was not found", err.Error())
		assert.Nil(t, invalidInfo)
	})
}
