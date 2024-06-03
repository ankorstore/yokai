package uuidv7

import googleuuid "github.com/google/uuid"

// UuidV7Generator is the interface for UUID v7 generators.
type UuidV7Generator interface {
	Generate() (googleuuid.UUID, error)
}

// DefaultUuidV7Generator is the default [UuidGenerator] implementation.
type DefaultUuidV7Generator struct{}

// NewDefaultUuidV7Generator returns a [DefaultUuidGenerator], implementing [UuidGenerator].
func NewDefaultUuidV7Generator() *DefaultUuidV7Generator {
	return &DefaultUuidV7Generator{}
}

// Generate returns a new UUID V7, using [Google UUID].
//
// [Google UUID]: https://github.com/google/uuid
func (g *DefaultUuidV7Generator) Generate() (googleuuid.UUID, error) {
	return googleuuid.NewV7()
}
