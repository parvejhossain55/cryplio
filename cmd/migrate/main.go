package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"cryplio/pkg/database"
)

// CLI holds parsed command-line flags
type CLI struct {
	CreateDB      bool
	MigrationsDir string
	StatusOnly    bool
	DownOne       bool
	Steps         int
	Config        *database.Config
}

// ParseFlags parses command-line arguments into a CLI struct
func ParseFlags() *CLI {
	cfg := database.DefaultConfig()

	createDB := flag.Bool("create-db", false, "Create database if not exists")
	migrationsDir := flag.String("dir", "./migrations", "Migrations directory")
	statusOnly := flag.Bool("status", false, "Show migration status only (don't apply)")
	downOne := flag.Bool("down", false, "Rollback the last migration")
	steps := flag.Int("steps", 0, "Run a specific number of migration steps (negative for rollback)")

	flag.StringVar(&cfg.Host, "host", cfg.Host, "PostgreSQL host")
	flag.IntVar(&cfg.Port, "port", cfg.Port, "PostgreSQL port")
	flag.StringVar(&cfg.User, "user", cfg.User, "PostgreSQL user")
	flag.StringVar(&cfg.Password, "password", cfg.Password, "PostgreSQL password")
	flag.StringVar(&cfg.DBName, "dbname", cfg.DBName, "Database name")
	flag.Parse()

	return &CLI{
		CreateDB:      *createDB,
		MigrationsDir: *migrationsDir,
		StatusOnly:    *statusOnly,
		DownOne:       *downOne,
		Steps:         *steps,
		Config:        cfg,
	}
}

func main() {
	cli := ParseFlags()

	if cli.StatusOnly {
		runStatus(cli)
		return
	}

	if cli.DownOne {
		cli.Steps = -1
	}

	runMigrate(cli)
}

func runStatus(cli *CLI) {
	printer := NewStatusPrinter(cli.MigrationsDir)
	connector := NewDBConnector(cli.Config)

	if err := printer.PrintStatus(connector); err != nil {
		log.Fatalf("Status error: %v", err)
	}
}

func runMigrate(cli *CLI) {
	connector := NewDBConnector(cli.Config)

	// Ensure DB exists if requested
	if cli.CreateDB {
		if err := connector.EnsureDatabase(); err != nil {
			log.Fatalf("Failed to ensure database: %v", err)
		}
		fmt.Printf("✓ Database '%s' ready\n", cli.Config.DBName)
	}

	// Connect
	db, err := connector.Connect()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("Cannot ping database: %v", err)
	}

	fmt.Printf("✓ Connected to '%s'\n\n", cli.Config.DBName)

	// Migrate
	migrator, err := database.NewMigrator(cli.Config, cli.MigrationsDir)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v", err)
	}
	if cli.Steps != 0 {
		fmt.Printf("Running migration steps: %d...\n", cli.Steps)
		if err := migrator.Steps(cli.Steps); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("✓ Migration steps completed successfully.")
		return
	}

	fmt.Println("Applying migrations...")

	if err := migrator.Apply(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("✓ All migrations applied successfully.")
}

// DBConnector abstracts database connection operations
type DBConnector struct {
	cfg *database.Config
}

func NewDBConnector(cfg *database.Config) *DBConnector {
	return &DBConnector{cfg: cfg}
}

func (c *DBConnector) Connect() (*sql.DB, error) {
	return database.Open(c.cfg)
}

func (c *DBConnector) EnsureDatabase() error {
	return database.EnsureDatabase(c.cfg)
}

// StatusPrinter handles printing migration status
type StatusPrinter struct {
	migrationsDir string
}

func NewStatusPrinter(migrationsDir string) *StatusPrinter {
	return &StatusPrinter{migrationsDir: migrationsDir}
}

func (sp *StatusPrinter) PrintStatus(connector *DBConnector) error {
	// Try DB first
	db, err := connector.Connect()
	if err == nil {
		err = db.Ping()
	}
	if err != nil {
		// Fallback to file-only mode
		sp.printFileOnly()
		return nil
	}
	defer db.Close()

	// Get data
	applied, err := database.AppliedVersionsDB(db)
	if err != nil {
		return fmt.Errorf("query applied versions: %w", err)
	}

	migs, err := sp.listMigrations()
	if err != nil {
		return fmt.Errorf("list migrations: %w", err)
	}

	// Print
	fmt.Println("Migration Status")
	fmt.Println("=================")
	fmt.Printf("Database: %s@%s:%d/%s\n",
		connector.cfg.User, connector.cfg.Host, connector.cfg.Port, connector.cfg.DBName)

	appliedSet := make(map[int]bool)
	for _, v := range applied {
		appliedSet[v] = true
	}

	pending := 0
	for _, m := range migs {
		if appliedSet[m.version] {
			fmt.Printf("  [✓] %03d  %s\n", m.version, m.name)
		} else {
			fmt.Printf("  [ ] %03d  %s\n", m.version, m.name)
			pending++
		}
	}

	fmt.Println()
	if pending == 0 {
		fmt.Println("✓ Database is up to date")
	} else {
		fmt.Printf("⚠  %d migration(s) pending\n", pending)
	}
	return nil
}

func (sp *StatusPrinter) printFileOnly() {
	migs, err := sp.listMigrations()
	if err != nil {
		log.Fatalf("Failed to list migrations: %v", err)
	}

	fmt.Println("Migration Status (files only)")
	fmt.Println("==============================")
	for _, m := range migs {
		fmt.Printf("  [ ] %03d  %s\n", m.version, m.name)
	}
	fmt.Printf("\n⚠  %d migration(s) pending\n", len(migs))
	fmt.Println("\nTo apply: go run ./cmd/db-migrate -create-db \\")
	fmt.Println("  -host=localhost -port=5432 -user=postgres -password=... -dbname=cryplio_db")
}

// listMigrations returns sorted list of migration files
func (sp *StatusPrinter) listMigrations() ([]migration, error) {
	files, err := os.ReadDir(sp.migrationsDir)
	if err != nil {
		return nil, err
	}

	var migs []migration
	for _, f := range files {
		name := f.Name()
		if !f.IsDir() && strings.HasSuffix(name, ".up.sql") {
			var v int
			if _, err := fmt.Sscanf(name, "%03d", &v); err == nil {
				migs = append(migs, migration{v, name})
			}
		}
	}

	sort.Slice(migs, func(i, j int) bool {
		return migs[i].version < migs[j].version
	})

	return migs, nil
}

type migration struct {
	version int
	name    string
}
