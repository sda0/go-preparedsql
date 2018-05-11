# About
Package stores list of map queryName=>queryPreparedStatement linked to the sql/db object.

## Usage example

```go
import (
	"github.com/sda0/go-preparedsql"
)

const (
	sqlGetFilesTerms = "sqlGetFilesTerms"
	sqlDeleteFilesTerms = "sqlDeleteFilesTerms"
)
type SQLStorage struct {
	connect   *sql.DB
	prepQuery *preparedsql.Registry
}

func init() {
	// Add queries to registry
	preparedsql.Add(sqlGetFilesTerms, "SELECT * FROM term_file WHERE fid::text=ANY($1)")
	preparedsql.Add(sqlDeleteFilesTerms, "DELETE FROM term_file tf WHERE fid=$1")
}

func NewSQLStorage() (*SQLStorage, error) {
	pg := &SQLStorage{}
	_, err := pg.getConnect()
	if err != nil {
		return nil, err
	}
	err = pg.MigrateUp()
	if err != nil {
		pg.connect.Close()
		return nil, err
	}
	// Prepare all registered queries
	pg.prepQuery, err = preparedsql.New(pg.connect)
	if err != nil {
		pg.connect.Close()
		return nil, err
	}
	return pg, nil
}

// GetFilesTerms returns all terms associated for requested files
func (s *SQLStorage) GetFilesTerms(ctx context.Context, files []string) ([]*storage.FileTerms, error) {
	getFilesTermsQuery, err := s.prepQuery.Get(sqlGetFilesTerms)
	if err != nil {
		return nil, err
	}
	rows, err := getFilesTermsQuery.QueryContext(ctx, pg2.Array(files))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*storage.FileTerms
	for rows.Next() {
		var r storage.FileTerms
		err = rows.StructScan(&r)
		if err != nil {
			return nil, err
		}
		result = append(result, &r)
	}
	return result, nil
}

// RemoveFilesTerms uses transactions 
func (s *SQLStorage) RemoveFilesTerms(ctx context.Context, vv []int) (error) {
	tx, err := s.connect.Begin()
	if err != nil {
		return nil
	}
	deleteFilesTermsQuery, err := s.prepQuery.GetTx(tx, sqlDeleteFilesTerms)
	if err != nil {
		tx.Rollback()
		return nil
	}
	for _, v := range vv {
		r, err := deleteFilesTermsQuery.ExecContext(ctx, v)
		if err != nil {
			tx.Rollback()
			return nil
		}
	}
	return tx.Commit()
}
