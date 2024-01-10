package uuid

// TestUuidGenerator is a [UuidGenerator] implementation allowing deterministic generations (for testing).
type TestUuidGenerator struct {
	value string
}

// NewTestUuidGenerator returns a [TestUuidGenerator], implementing [UuidGenerator].
//
// It accepts a value that will be used for deterministic generation results.
func NewTestUuidGenerator(value string) *TestUuidGenerator {
	return &TestUuidGenerator{
		value: value,
	}
}

// SetValue sets the value to use for deterministic generations.
func (g *TestUuidGenerator) SetValue(value string) *TestUuidGenerator {
	g.value = value

	return g
}

// Generate returns the configured deterministic value.
func (g *TestUuidGenerator) Generate() string {
	return g.value
}
