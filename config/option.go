package config

import (
	"os"
	"path"
)

// Options are options for the [ConfigFactory] implementations.
type Options struct {
	FileName  string
	FilePaths []string
}

// DefaultConfigOptions are the default options used in the [DefaultConfigFactory].
func DefaultConfigOptions() Options {
	opts := Options{
		FileName: "config",
		FilePaths: []string{
			".",
			"./configs",
		},
	}

	// check if the OCI image has been build using ko and use KO_DATA_PATH as config root
	if val, ok := os.LookupEnv("KO_DATA_PATH"); ok {
		opts.FilePaths = append(opts.FilePaths, val, path.Join(val, "configs"))
	}

	return opts
}

// ConfigOption are functional options for the [ConfigFactory] implementations.
type ConfigOption func(o *Options)

// WithFileName is used to specify the file base name (without extension) of the config file to load.
func WithFileName(n string) ConfigOption {
	return func(o *Options) {
		o.FileName = n
	}
}

// WithFilePaths is used to specify the list of file paths to lookup config files to load.
func WithFilePaths(p ...string) ConfigOption {
	return func(o *Options) {
		o.FilePaths = p
	}
}
