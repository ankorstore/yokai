package fxvalidator

// AliasDefinition is the interface for aliases definitions.
type AliasDefinition interface {
	Alias() string
	Tags() string
}

type aliasDefinition struct {
	alias string
	tags  string
}

// NewAliasDefinition returns a new [CronJobDefinition].
func NewAliasDefinition(alias string, tags string) AliasDefinition {
	return &aliasDefinition{
		alias: alias,
		tags:  tags,
	}
}

// Alias returns the definition alias.
func (a *aliasDefinition) Alias() string {
	return a.alias
}

// Tags returns the definition tags.
func (a *aliasDefinition) Tags() string {
	return a.tags
}
