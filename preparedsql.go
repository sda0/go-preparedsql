// Package preparedSQL handles statements preparing
package preparedsql

import (
	"database/sql"
	"fmt"
)

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

func (r *Registry) Prepare(db *sql.DB) (err error) {
	for name, query := range queryRegistry {
		r.storage[name], err = db.Prepare(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Registry) Get(query string) (*sql.Stmt, error) {
	if r.storage[query] == nil {
		if queryRegistry[query] == "" {
			return nil, fmt.Errorf("prepared query '%s' is not added", query)
		}
		return nil, fmt.Errorf("query '%s' is not prepared", query)
	}
	return r.storage[query], nil
}
