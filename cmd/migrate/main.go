// Command migrate menjalankan migration database.
//
// Penggunaan:
//
//	go run ./cmd/migrate up             jalankan seluruh migration yang belum dipakai
//	go run ./cmd/migrate down           rollback seluruh migration
//	go run ./cmd/migrate up 1           maju 1 langkah
//	go run ./cmd/migrate down 1         mundur 1 langkah
//	go run ./cmd/migrate version        tampilkan versi aktif
//	go run ./cmd/migrate force 3        tandai versi 3 sebagai bersih (perbaiki state dirty)
package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"

	"go-projects/hexagonal-example/config"
	"go-projects/hexagonal-example/migration"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up|down|version|force> [n]")
	}

	source, err := iofs.New(migration.SQL, "sql")
	if err != nil {
		log.Fatalf("gagal membaca file migration: %s", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, dsn())
	if err != nil {
		log.Fatalf("gagal terhubung ke database: %s", err)
	}
	defer m.Close()

	if err := run(m, os.Args[1], os.Args[2:]); err != nil {
		log.Fatalf("migration gagal: %s", err)
	}

	printVersion(m)
}

func run(m *migrate.Migrate, cmd string, args []string) error {
	steps, err := optionalSteps(args)
	if err != nil {
		return err
	}

	switch cmd {
	case "up":
		if steps > 0 {
			return ignoreNoChange(m.Steps(steps))
		}
		return ignoreNoChange(m.Up())

	case "down":
		if steps > 0 {
			return ignoreNoChange(m.Steps(-steps))
		}
		return ignoreNoChange(m.Down())

	case "version":
		return nil

	case "force":
		if len(args) != 1 {
			return errors.New("force membutuhkan nomor versi, contoh: migrate force 3")
		}
		return m.Force(steps)

	default:
		return fmt.Errorf("perintah tidak dikenal: %q", cmd)
	}
}

// optionalSteps membaca argumen jumlah langkah. Tidak adanya argumen berarti 0
// (semua langkah), bukan error.
func optionalSteps(args []string) (int, error) {
	if len(args) == 0 {
		return 0, nil
	}

	n, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, fmt.Errorf("jumlah langkah harus berupa angka: %q", args[0])
	}

	return n, nil
}

// ignoreNoChange memperlakukan "tidak ada migration baru" sebagai sukses supaya
// perintah ini aman dipanggil berulang, misalnya dari skrip startup.
func ignoreNoChange(err error) error {
	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("tidak ada migration baru")
		return nil
	}
	return err
}

func printVersion(m *migrate.Migrate) {
	version, dirty, err := m.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		log.Println("versi: belum ada migration yang dijalankan")
		return
	}
	if err != nil {
		log.Fatalf("gagal membaca versi: %s", err)
	}

	state := "clean"
	if dirty {
		state = "DIRTY - perbaiki dengan: migrate force <versi>"
	}
	log.Printf("versi: %d (%s)", version, state)
}

func dsn() string {
	cfg := config.LoadPostgresConfig()

	sslMode := cfg.SslMode
	if sslMode == "" {
		sslMode = "disable"
	}

	q := url.Values{}
	q.Set("sslmode", sslMode)

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Path:     cfg.DBName,
		RawQuery: q.Encode(),
	}

	return u.String()
}
