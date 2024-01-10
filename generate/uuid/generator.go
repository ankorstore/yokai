package uuid

import googleuuid "github.com/google/uuid"

// UuidGenerator is the interface for UUID generators.
type UuidGenerator interface {
	Generate() string
}

// DefaultUuidGenerator is the default [UuidGenerator] implementation.
type DefaultUuidGenerator struct{}

// NewDefaultUuidGenerator returns a [DefaultUuidGenerator], implementing [UuidGenerator].
func NewDefaultUuidGenerator() *DefaultUuidGenerator {
	return &DefaultUuidGenerator{}
}

// Generate returns a new UUID, using [Google UUID].
//
// [Google UUID]: https://github.com/google/uuid
func (g *DefaultUuidGenerator) Generate() string {
	return googleuuid.New().String()
}
