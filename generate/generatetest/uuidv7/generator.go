package uuid

// TestUuidV7Generator is a [UuidV7Generator] implementation allowing deterministic generations (for testing).
type TestUuidV7Generator struct {
	value string
}

// NewTestUuidV7Generator returns a [TestUuidGenerator], implementing [UuidGenerator].
//
// It accepts a value that will be used for deterministic generation results.
func NewTestUuidV7Generator(value string) *TestUuidV7Generator {
	return &TestUuidV7Generator{
		value: value,
	}
}

// SetValue sets the value to use for deterministic generations.
func (g *TestUuidV7Generator) SetValue(value string) *TestUuidV7Generator {
	g.value = value

	return g
}

// Generate returns the configured deterministic value.
func (g *TestUuidV7Generator) Generate() (string, error) {
	return g.value, nil
}
