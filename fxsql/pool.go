package fxsql

import "fmt"

type DatabasePool struct {
	primary     *Database
	auxiliaries map[string]*Database
}

func NewDatabasePool(primary *Database, auxiliaries ...*Database) *DatabasePool {
	auxiliariesMap := make(map[string]*Database, len(auxiliaries))
	for _, aux := range auxiliaries {
		auxiliariesMap[aux.Name()] = aux
	}

	return &DatabasePool{
		primary:     primary,
		auxiliaries: auxiliariesMap,
	}
}

func (p *DatabasePool) Primary() *Database {
	return p.primary
}

func (p *DatabasePool) Auxiliary(name string) (*Database, error) {
	if db, ok := p.auxiliaries[name]; ok {
		return db, nil
	}

	return nil, fmt.Errorf("database with name %s was not found", name)
}

func (p *DatabasePool) Auxiliaries() map[string]*Database {
	return p.auxiliaries
}
