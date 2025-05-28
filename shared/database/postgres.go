package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// PostgresDB wrapper untuk database connection yang mudah digunakan
type PostgresDB struct {
	Connection *sqlx.DB
	logger     *logrus.Logger
}

// DatabaseConfig konfigurasi database yang user-friendly
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DatabaseName    string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// NewPostgresConnection membuat koneksi baru ke PostgreSQL
func NewPostgresConnection(config DatabaseConfig, logger *logrus.Logger) (*PostgresDB, error) {
	// Build connection string yang mudah dibaca
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DatabaseName,
		config.SSLMode,
	)

	logger.WithFields(logrus.Fields{
		"host":     config.Host,
		"port":     config.Port,
		"database": config.DatabaseName,
		"user":     config.User,
	}).Info("Connecting to PostgreSQL database...")

	// Buka koneksi database
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		logger.WithError(err).Error("Failed to connect to database")
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings untuk performance optimal
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	// Test koneksi
	if err := db.Ping(); err != nil {
		logger.WithError(err).Error("Failed to ping database")
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Successfully connected to PostgreSQL database")

	return &PostgresDB{
		Connection: db,
		logger:     logger,
	}, nil
}

// Close menutup koneksi database dengan graceful
func (db *PostgresDB) Close() error {
	if db.Connection != nil {
		db.logger.Info("Closing database connection...")
		return db.Connection.Close()
	}
	return nil
}

// HealthCheck mengecek kesehatan database connection
func (db *PostgresDB) HealthCheck() error {
	if db.Connection == nil {
		return fmt.Errorf("database connection is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return db.Connection.PingContext(ctx)
}

// GetStats mengembalikan statistik connection pool
func (db *PostgresDB) GetStats() sql.DBStats {
	if db.Connection == nil {
		return sql.DBStats{}
	}
	return db.Connection.Stats()
}

// Transaction helper untuk menjalankan operasi dalam transaction
func (db *PostgresDB) Transaction(fn func(*sqlx.Tx) error) error {
	tx, err := db.Connection.Beginx()
	if err != nil {
		db.logger.WithError(err).Error("Failed to begin transaction")
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-throw panic setelah rollback
		} else if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				db.logger.WithError(rollbackErr).Error("Failed to rollback transaction")
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				db.logger.WithError(commitErr).Error("Failed to commit transaction")
				err = commitErr
			}
		}
	}()

	err = fn(tx)
	return err
}

// ExecuteInTransaction helper untuk execute query dalam transaction
func (db *PostgresDB) ExecuteInTransaction(queries []string) error {
	return db.Transaction(func(tx *sqlx.Tx) error {
		for _, query := range queries {
			if _, err := tx.Exec(query); err != nil {
				db.logger.WithError(err).WithField("query", query).Error("Failed to execute query in transaction")
				return fmt.Errorf("failed to execute query: %w", err)
			}
		}
		return nil
	})
}

// MigrateSchema helper untuk menjalankan database migrations
func (db *PostgresDB) MigrateSchema(migrationQueries []string) error {
	db.logger.Info("Starting database migration...")

	for i, query := range migrationQueries {
		db.logger.WithField("migration_step", i+1).Info("Executing migration step")

		if _, err := db.Connection.Exec(query); err != nil {
			db.logger.WithError(err).WithField("migration_step", i+1).Error("Migration failed")
			return fmt.Errorf("migration step %d failed: %w", i+1, err)
		}
	}

	db.logger.Info("Database migration completed successfully")
	return nil
}

// LogConnectionStats mencatat statistik connection pool
func (db *PostgresDB) LogConnectionStats() {
	stats := db.GetStats()
	db.logger.WithFields(logrus.Fields{
		"open_connections":    stats.OpenConnections,
		"in_use":              stats.InUse,
		"idle":                stats.Idle,
		"wait_count":          stats.WaitCount,
		"wait_duration":       stats.WaitDuration,
		"max_idle_closed":     stats.MaxIdleClosed,
		"max_lifetime_closed": stats.MaxLifetimeClosed,
	}).Info("Database connection pool statistics")
}
