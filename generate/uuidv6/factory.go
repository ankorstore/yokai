package uuidv6

// UuidV6GeneratorFactory is the interface for [UuidV6Generator] factories.
type UuidV6GeneratorFactory interface {
	Create() UuidV6Generator
}

// DefaultUuidV6GeneratorFactory is the default [UuidV6GeneratorFactory] implementation.
type DefaultUuidV6GeneratorFactory struct{}

// NewDefaultUuidV6GeneratorFactory returns a [DefaultUuidV6GeneratorFactory], implementing [UuidV6GeneratorFactory].
func NewDefaultUuidV6GeneratorFactory() UuidV6GeneratorFactory {
	return &DefaultUuidV6GeneratorFactory{}
}

// Create returns a new [UuidV6Generator].
func (g *DefaultUuidV6GeneratorFactory) Create() UuidV6Generator {
	return NewDefaultUuidV6Generator()
}
