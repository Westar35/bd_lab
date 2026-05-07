package db

import (
	"database/sql"
	"fmt"
)

// DBType identifies a supported DBMS.
type DBType string

const (
	DBPostgres DBType = "postgres"
	DBMySQL    DBType = "mysql"
)

func ParseDBType(value string, fallback DBType) DBType {
	switch DBType(value) {
	case DBPostgres:
		return DBPostgres
	case DBMySQL:
		return DBMySQL
	default:
		return fallback
	}
}

func (t DBType) DisplayName() string {
	switch t {
	case DBMySQL:
		return "MySQL"
	default:
		return "PostgreSQL"
	}
}

func (t DBType) Placeholder(n int) string {
	if t == DBMySQL {
		return "?"
	}
	return fmt.Sprintf("$%d", n)
}

// Manager owns connections to all supported databases.
type Manager struct {
	defaultDB DBType
	conns     map[DBType]*sql.DB
}

func NewManager(defaultDB DBType, postgresDB, mysqlDB *sql.DB) *Manager {
	return &Manager{
		defaultDB: defaultDB,
		conns: map[DBType]*sql.DB{
			DBPostgres: postgresDB,
			DBMySQL:    mysqlDB,
		},
	}
}

func (m *Manager) Default() DBType {
	return m.defaultDB
}

func (m *Manager) Get(kind DBType) (*sql.DB, DBType, error) {
	kind = ParseDBType(string(kind), m.defaultDB)
	conn := m.conns[kind]
	if conn == nil {
		return nil, kind, fmt.Errorf("соединение с %s не настроено", kind.DisplayName())
	}
	return conn, kind, nil
}

func (m *Manager) Close() error {
	var firstErr error
	for _, conn := range m.conns {
		if conn == nil {
			continue
		}
		if err := conn.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
