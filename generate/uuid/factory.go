package uuid

// UuidGeneratorFactory is the interface for [UuidGenerator] factories.
type UuidGeneratorFactory interface {
	Create() UuidGenerator
}

// DefaultUuidGeneratorFactory is the default [UuidGeneratorFactory] implementation.
type DefaultUuidGeneratorFactory struct{}

// NewDefaultUuidGeneratorFactory returns a [DefaultUuidGeneratorFactory], implementing [UuidGeneratorFactory].
func NewDefaultUuidGeneratorFactory() UuidGeneratorFactory {
	return &DefaultUuidGeneratorFactory{}
}

// Create returns a new [UuidGenerator].
// For example:
//
//	var generator, _ = uuid.NewDefaultConfigFactory().Create()
func (g *DefaultUuidGeneratorFactory) Create() UuidGenerator {
	return NewDefaultUuidGenerator()
}
