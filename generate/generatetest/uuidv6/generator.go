package uuidv6

import googleuuid "github.com/google/uuid"

// TestUuidV6Generator is a [UuidV6Generator] implementation allowing deterministic generations (for testing).
type TestUuidV6Generator struct {
	value string
}

// NewTestUuidV6Generator returns a [TestUuidGenerator], implementing [UuidGenerator].
//
// It accepts a value that will be used for deterministic generation results.
func NewTestUuidV6Generator(value string) (*TestUuidV6Generator, error) {
	err := googleuuid.Validate(value)
	if err != nil {
		return nil, err
	}

	return &TestUuidV6Generator{
		value: value,
	}, nil
}

// SetValue sets the value to use for deterministic generations.
func (g *TestUuidV6Generator) SetValue(value string) error {
	err := googleuuid.Validate(value)
	if err != nil {
		return err
	}

	g.value = value

	return nil
}

// Generate returns the configured deterministic value.
func (g *TestUuidV6Generator) Generate() (googleuuid.UUID, error) {
	return googleuuid.Parse(g.value)
}
