package config_test

import (
	"testing"

	"github.com/ankorstore/yokai/config"
	"github.com/stretchr/testify/assert"
)

func TestWithFileName(t *testing.T) {
	option := config.WithFileName("test")

	opts := &config.Options{}
	option(opts)

	assert.Equal(t, "test", opts.FileName)
}

func TestWithFilePaths(t *testing.T) {
	option := config.WithFilePaths("path1", "path2")

	opts := &config.Options{}
	option(opts)

	assert.Equal(t, []string{"path1", "path2"}, opts.FilePaths)
}
