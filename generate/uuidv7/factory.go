package uuidv7

// UuidV7GeneratorFactory is the interface for [UuidV7Generator] factories.
type UuidV7GeneratorFactory interface {
	Create() UuidV7Generator
}

// DefaultUuidV7GeneratorFactory is the default [UuidV7GeneratorFactory] implementation.
type DefaultUuidV7GeneratorFactory struct{}

// NewDefaultUuidV7GeneratorFactory returns a [DefaultUuidV7GeneratorFactory], implementing [UuidV7GeneratorFactory].
func NewDefaultUuidV7GeneratorFactory() UuidV7GeneratorFactory {
	return &DefaultUuidV7GeneratorFactory{}
}

// Create returns a new [UuidV7Generator].
func (g *DefaultUuidV7GeneratorFactory) Create() UuidV7Generator {
	return NewDefaultUuidV7Generator()
}
