// Package preparedSQL handles statements preparing
package preparedsql

import "database/sql"

var queryRegistry = map[string]string{}

// Add query to global queries string registry mapped to query names
func Add(name, query string) {
	queryRegistry[name] = query
}

// Registry keeps prepared statements linked to db for each query of global queryRegistry
type Registry struct {
	storage map[string]*sql.Stmt
}

func New(db *sql.DB) (*Registry, error) {
	registry := &Registry{storage: make(map[string]*sql.Stmt, len(queryRegistry))}
	err := registry.Prepare(db)
	if err != nil {
		return nil, err
	}
	return registry, nil
}

func (m *Registry) Prepare(db *sql.DB) (err error) {
	for name, query := range queryRegistry {
		m.storage[name], err = db.Prepare(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Registry) Get(query string) *sql.Stmt {
	return m.storage[query]
}
