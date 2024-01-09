package config

// Options are options for the [ConfigFactory] implementations.
type Options struct {
	FileName  string
	FilePaths []string
}

// DefaultConfigOptions are the default options used in the [DefaultConfigFactory].
func DefaultConfigOptions() Options {
	return Options{
		FileName: "config",
		FilePaths: []string{
			".",
			"./configs",
		},
	}
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
