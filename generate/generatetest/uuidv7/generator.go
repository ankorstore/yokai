package uuid

import googleuuid "github.com/google/uuid"

// TestUuidV7Generator is a [UuidV7Generator] implementation allowing deterministic generations (for testing).
type TestUuidV7Generator struct {
	value string
}

// NewTestUuidV7Generator returns a [TestUuidGenerator], implementing [UuidGenerator].
//
// It accepts a value that will be used for deterministic generation results.
func NewTestUuidV7Generator(value string) (*TestUuidV7Generator, error) {
	err := googleuuid.Validate(value)
	if err != nil {
		return nil, err
	}

	return &TestUuidV7Generator{
		value: value,
	}, nil
}

// SetValue sets the value to use for deterministic generations.
func (g *TestUuidV7Generator) SetValue(value string) error {
	err := googleuuid.Validate(value)
	if err != nil {
		return err
	}

	g.value = value

	return nil
}

// Generate returns the configured deterministic value.
func (g *TestUuidV7Generator) Generate() (googleuuid.UUID, error) {
	return googleuuid.Parse(g.value)
}
