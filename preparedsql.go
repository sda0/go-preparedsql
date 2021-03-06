// Package preparedSQL handles statements preparing
package preparedsql

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
)

var queryRegistry = map[string]string{}

// Add query to global queries string registry mapped to query names
func Add(name, query string) {
	queryRegistry[name] = query
}

// Registry keeps prepared statements linked to db for each query of global queryRegistry
type Registry struct {
	db *sql.DB
	storage map[string]*sql.Stmt
}

// New prepares all queries listed in registry to provided DB and returns Registry instance
func New(db *sql.DB) (*Registry, error) {
	registry := &Registry{db: db, storage: make(map[string]*sql.Stmt, len(queryRegistry))}
	err := registry.Prepare(db)
	if err != nil {
		return nil, err
	}
	return registry, nil
}

// Prepare all queries listed in registry to provided DB
func (r *Registry) Prepare(db *sql.DB) (err error) {
	for name, query := range queryRegistry {
		r.storage[name], err = db.Prepare(query)
		if err != nil {
			return errors.Wrapf(err, "cannot prepare query %q", name)
		}
	}
	return nil
}

// Get returns statement to be executed
func (r *Registry) Get(query string) (*sql.Stmt, error) {
	if r.storage[query] == nil {
		if queryRegistry[query] == "" {
			return nil, fmt.Errorf("prepared query '%s' is not added", query)
		}
		stmt, err := r.db.Prepare(queryRegistry[query])
		if err != nil {
			return nil, err
		}
		r.storage[query] = stmt
	}
	return r.storage[query], nil
}

// GetTx returns statement to be executed linked to provided transaction
func (r *Registry) GetTx(tx *sql.Tx, query string) (*sql.Stmt, error) {
	pg, err := r.Get(query)
	if err != nil {
		return nil, err
	}
	return tx.Stmt(pg), nil
}
