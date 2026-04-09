package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sdsouz12/SER594-Team24-GateAndCrownB2B/backend-service/pkg/envload"
)

// toPgx5MigrateURL converts DATABASE_URL for golang-migrate's pgx5 driver.
// Do not set x-multi-statement=true globally: naive splitting on ";"
// breaks PL/pgSQL blocks (DO $$ ... $$). Each .sql file is sent as one batch by default.
func toPgx5MigrateURL(dbURL string) string {
	u := strings.TrimSpace(dbURL)
	switch {
	case strings.HasPrefix(u, "postgresql://"):
		u = "pgx5://" + strings.TrimPrefix(u, "postgresql://")
	case strings.HasPrefix(u, "postgres://"):
		u = "pgx5://" + strings.TrimPrefix(u, "postgres://")
	default:
		if !strings.HasPrefix(u, "pgx5://") {
			u = "pgx5://" + u
		}
	}
	return u
}

// safeDatabaseName allows only unquoted PostgreSQL identifiers (letters, digits, underscore).
func safeDatabaseName(name string) bool {
	if name == "" || len(name) > 63 {
		return false
	}
	for i, r := range name {
		if i == 0 && unicode.IsDigit(r) {
			return false
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}

// ensureTargetDatabase creates the database named in DATABASE_URL if it does not exist.
// Connects using POSTGRES_MAINTENANCE_DATABASE (default "postgres"). Requires a user with CREATEDB (or superuser).
func ensureTargetDatabase(dbURL string) error {
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("parse DATABASE_URL: %w", err)
	}
	targetDB := strings.TrimSpace(cfg.ConnConfig.Database)
	if targetDB == "" {
		return fmt.Errorf("DATABASE_URL must include a database name (path after host/port)")
	}
	if !safeDatabaseName(targetDB) {
		return fmt.Errorf("database name %q must be letters, digits, underscore only, max 63 chars, and not start with a digit", targetDB)
	}

	maintenance := strings.TrimSpace(os.Getenv("POSTGRES_MAINTENANCE_DATABASE"))
	if maintenance == "" {
		maintenance = "postgres"
	}
	if !safeDatabaseName(maintenance) {
		return fmt.Errorf("POSTGRES_MAINTENANCE_DATABASE %q is invalid", maintenance)
	}

	cfg.ConnConfig.Database = maintenance
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return fmt.Errorf("connect to maintenance database %q: %w", maintenance, err)
	}
	defer pool.Close()

	var exists bool
	if err := pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)`, targetDB).Scan(&exists); err != nil {
		return fmt.Errorf("check database exists: %w", err)
	}
	if exists {
		return nil
	}

	// Identifier is validated; safe to interpolate.
	if _, err := pool.Exec(ctx, fmt.Sprintf(`CREATE DATABASE %s`, targetDB)); err != nil {
		return fmt.Errorf("CREATE DATABASE %q: %w (does this role have CREATEDB?)", targetDB, err)
	}
	log.Printf("created database %q", targetDB)
	return nil
}

func migrationsDir(wd string) string {
	sub := filepath.Join(wd, "backend-service", "migrations")
	if fi, err := os.Stat(sub); err == nil && fi.IsDir() {
		return sub
	}
	return filepath.Join(wd, "migrations")
}

func main() {
	envload.Load()

	dbURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required (e.g. in backend-service/.env)")
	}

	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	migDir := migrationsDir(root)
	src := "file://" + filepath.ToSlash(migDir)

	cmd := "up"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	if cmd == "up" {
		if err := ensureTargetDatabase(dbURL); err != nil {
			log.Fatalf("ensure database: %v", err)
		}
	}

	m, err := migrate.New(src, toPgx5MigrateURL(dbURL))
	if err != nil {
		log.Fatal(err)
	}
	defer m.Close()

	switch cmd {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(err)
		}
		v, dirty, verr := m.Version()
		if verr != nil {
			log.Fatal(verr)
		}
		fmt.Printf("migrations ok: version=%d dirty=%v\n", v, dirty)
	case "down":
		if err := m.Steps(-1); err != nil {
			log.Fatal(err)
		}
		fmt.Println("migrations: stepped down 1")
	case "version":
		v, dirty, err := m.Version()
		if err != nil {
			if errors.Is(err, migrate.ErrNilVersion) {
				fmt.Println("no migrations applied")
				return
			}
			log.Fatal(err)
		}
		fmt.Printf("version=%d dirty=%v\n", v, dirty)
	case "force":
		if len(os.Args) < 3 {
			log.Fatalf("usage: %s force VERSION  (set schema_migrations without running SQL; use to baseline a DB)", filepath.Base(os.Args[0]))
		}
		n, err := strconv.Atoi(os.Args[2])
		if err != nil || n < 0 {
			log.Fatal("force VERSION: VERSION must be a non-negative integer (e.g. 30)")
		}
		if err := m.Force(n); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("forced schema_migrations version=%d dirty=false\n", n)
	case "unwind-dirty":
		v, dirty, err := m.Version()
		if err != nil {
			log.Fatal(err)
		}
		if !dirty {
			fmt.Println("database is not dirty; nothing to unwind")
			return
		}
		prev := int(v) - 1
		if prev < 0 {
			prev = 0
		}
		if err := m.Force(prev); err != nil {
			log.Fatal(err)
		}
		self := filepath.Base(os.Args[0])
		fmt.Printf("dirty flag cleared: schema_migrations now version=%d. Re-run: %s up\n", prev, self)
		fmt.Printf("(This re-opens migration %d so it can run again with fixed SQL.)\n", v)
	default:
		log.Fatalf("usage: %s [up|down|version|force|unwind-dirty]", filepath.Base(os.Args[0]))
	}
}
