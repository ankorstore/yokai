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

func TestNewFxModuleInfoRegistry(t *testing.T) {
	t.Parallel()

	registry, err := prepareTestFxModuleInfoRegistry()
	assert.NoError(t, err)

	assert.IsType(t, &fxcore.FxModuleInfoRegistry{}, registry)
}

func TestAll(t *testing.T) {
	t.Parallel()

	registry, err := prepareTestFxModuleInfoRegistry()
	assert.NoError(t, err)

	assert.Len(t, registry.All(), 2)
}

func TestNames(t *testing.T) {
	t.Parallel()

	registry, err := prepareTestFxModuleInfoRegistry()
	assert.NoError(t, err)

	assert.Equal(t, []string{fxcore.ModuleName, "test"}, registry.Names())
}

func TestFind(t *testing.T) {
	t.Parallel()

	registry, err := prepareTestFxModuleInfoRegistry()
	assert.NoError(t, err)

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
}

func prepareTestFxModuleInfoRegistry() (*fxcore.FxModuleInfoRegistry, error) {
	cfg, err := config.NewDefaultConfigFactory().Create(
		config.WithFilePaths("./testdata/config"),
	)
	if err != nil {
		return nil, err
	}

	return fxcore.NewFxModuleInfoRegistry(fxcore.FxModuleInfoRegistryParam{
		Infos: []interface{}{
			&testModuleInfo{},
			fxcore.NewFxCoreModuleInfo(fxcore.FxCoreModuleInfoParam{
				Config:     cfg,
				ExtraInfos: []fxcore.FxExtraInfo{},
			}),
			"invalid",
		},
	}), nil
}
