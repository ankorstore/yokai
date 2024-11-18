package uuidv6

import googleuuid "github.com/google/uuid"

// UuidV6Generator is the interface for UUID v6 generators.
type UuidV6Generator interface {
	Generate() (googleuuid.UUID, error)
}

// DefaultUuidV6Generator is the default [UuidGenerator] implementation.
type DefaultUuidV6Generator struct{}

// NewDefaultUuidV6Generator returns a [DefaultUuidGenerator], implementing [UuidGenerator].
func NewDefaultUuidV6Generator() *DefaultUuidV6Generator {
	return &DefaultUuidV6Generator{}
}

// Generate returns a new UUID V6, using [Google UUID].
//
// [Google UUID]: https://github.com/google/uuid
func (g *DefaultUuidV6Generator) Generate() (googleuuid.UUID, error) {
	return googleuuid.NewV6()
}
