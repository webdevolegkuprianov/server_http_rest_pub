package sqlstore

import (
	"database/sql"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/webdevolegkuprianov/server_http_rest/app/apiserver/model"
	"github.com/webdevolegkuprianov/server_http_rest/app/apiserver/store"
)

//Store
type Store struct {
	dbPostgres     *pgxpool.Pool
	dbMssql        *sql.DB
	config         *model.Service
	userRepository *UserRepository
	dataRepository *DataRepository
}

//New_db
func New(db *pgxpool.Pool, dbmssql *sql.DB, config *model.Service) *Store {
	return &Store{
		dbPostgres: db,
		dbMssql:    dbmssql,
		config:     config,
	}
}

//User
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

//Data
func (s *Store) Data() store.DataRepository {
	if s.dataRepository != nil {
		return s.dataRepository
	}

	s.dataRepository = &DataRepository{
		store: s,
	}

	return s.dataRepository
}
