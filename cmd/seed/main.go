// Command seed mengisi data awal (dev): satu akun author, satu buyer, dan
// beberapa event contoh (gratis + berbayar) sekaligus counter kuota di Redis.
//
// Penggunaan:
//
//	go run ./cmd/seed
//
// Idempoten: dijalankan berkali-kali aman (author/user pakai ON CONFLICT,
// event dilewati bila namanya sudah ada).
package main

import (
	"context"
	"log"
	"time"

	"go-projects/hexagonal-example/internal/service"
	authorUc "go-projects/hexagonal-example/internal/service/entity/author"
	eventUc "go-projects/hexagonal-example/internal/service/entity/event"
	"go-projects/hexagonal-example/pkg"
	"go-projects/hexagonal-example/pkg/di"
)

func main() {
	log.SetFlags(0)

	container, err := di.Container()
	if err != nil {
		log.Fatalf("build container: %v", err)
	}

	err = container.Invoke(func(p pkg.Package, svc service.Service) error {
		ctx := context.Background()

		seedAuthor(ctx, svc)
		seedBuyer(p)
		seedEvents(ctx, p, svc)

		log.Println("\nseeder selesai.")
		return nil
	})
	if err != nil {
		log.Fatalf("seed: %v", err)
	}
}

// seedAuthor membuat akun author untuk login JWT (password di-hash bcrypt).
func seedAuthor(ctx context.Context, svc service.Service) {
	const (
		name  = "Admin Penyelenggara"
		email = "author@warticket.local"
		pass  = "password123"
	)
	err := svc.Author.Register(ctx, authorUc.RegisterRequest{Name: name, Email: email, Password: pass})
	switch {
	case err == nil:
		log.Printf("author dibuat  : %s / %s", email, pass)
	case err.Error() == "email sudah terdaftar":
		log.Printf("author ada     : %s (dilewati)", email)
	default:
		log.Fatalf("seed author: %v", err)
	}
}

// seedBuyer menyisipkan user pembeli dengan id=1 (dipakai header x-user-id).
func seedBuyer(p pkg.Package) {
	if err := p.DB.Exec(
		`INSERT INTO users (id, name, email, password) VALUES (1, ?, ?, ?) ON CONFLICT DO NOTHING`,
		"Budi Santoso", "buyer@warticket.local", "secret",
	).Error; err != nil {
		log.Fatalf("seed buyer: %v", err)
	}
	log.Printf("buyer          : id=1 buyer@warticket.local (x-user-id)")
}

// seedEvents membuat event contoh lewat service (sekalian set kuota Redis).
func seedEvents(ctx context.Context, p pkg.Package, svc service.Service) {
	start := time.Now().UTC().Add(7 * 24 * time.Hour)
	start = time.Date(start.Year(), start.Month(), start.Day(), 9, 0, 0, 0, time.UTC)
	end := start.Add(6 * time.Hour)

	events := []eventUc.CreateEventRequest{
		{
			Name:        "Konser Gratis WarTicket",
			Description: "Event gratis dengan formulir pendaftaran.",
			Price:       0,
			Quota:       100,
			StartDate:   start.Format(time.RFC3339),
			EndDate:     end.Format(time.RFC3339),
			FormFields: []eventUc.FormFieldInput{
				{Label: "Nama Lengkap", FieldType: "text", Required: true, Position: 1},
				{Label: "Domisili", FieldType: "select", Required: true, Options: []string{"Jakarta", "Bandung", "Surabaya"}, Position: 2},
			},
		},
		{
			Name:        "Konser Berbayar WarTicket",
			Description: "Event berbayar tanpa formulir.",
			Price:       150000,
			Quota:       50,
			StartDate:   start.Format(time.RFC3339),
			EndDate:     end.Format(time.RFC3339),
		},
	}

	for _, ev := range events {
		var count int64
		p.DB.Raw(`SELECT count(*) FROM events WHERE name = ?`, ev.Name).Scan(&count)
		if count > 0 {
			log.Printf("event ada      : %q (dilewati)", ev.Name)
			continue
		}
		if err := svc.Event.CreateEvent(ctx, ev); err != nil {
			log.Fatalf("seed event %q: %v", ev.Name, err)
		}

		var id int64
		p.DB.Raw(`SELECT id FROM events WHERE name = ?`, ev.Name).Scan(&id)
		log.Printf("event dibuat   : id=%d price=%.0f quota=%d %q", id, ev.Price, ev.Quota, ev.Name)
	}
}
