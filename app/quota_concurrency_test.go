package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	eventRepo "go-projects/hexagonal-example/internal/adapter/outbound/repository/event"
	"go-projects/hexagonal-example/pkg"
	"go-projects/hexagonal-example/pkg/di"
)

// TestE2E_ConcurrentDecrementNoOversell membuktikan DecrementQuota aman terhadap
// konkurensi: dari stok terbatas, jumlah decrement sukses tepat = stok, sisanya
// ditolak ErrInsufficientQuota, dan quota_remaining tidak pernah negatif.
func TestE2E_ConcurrentDecrementNoOversell(t *testing.T) {
	container, err := di.Container()
	if err != nil {
		t.Fatalf("build container: %v", err)
	}
	var p pkg.Package
	if err := container.Invoke(func(pp pkg.Package) error { p = pp; return nil }); err != nil {
		t.Fatalf("invoke container: %v", err)
	}
	repo := eventRepo.New(p)
	ctx := context.Background()

	const (
		stock   = 5  // sisa kuota awal
		workers = 50 // pembeli yang balapan, jauh > stok
	)

	// seed event langsung (start_date wajib; harga/nama tak relevan untuk uji ini).
	start := time.Now().UTC().Add(48 * time.Hour)
	var eventID int64
	if err := p.DB.Raw(
		`INSERT INTO events (name, price, quota, quota_remaining, start_date)
		 VALUES (?, 0, ?, ?, ?) RETURNING id`,
		fmt.Sprintf("E2E Concurrency %d", time.Now().UnixNano()), stock, stock, start,
	).Scan(&eventID).Error; err != nil {
		t.Fatalf("seed event: %v", err)
	}
	t.Cleanup(func() {
		p.DB.Exec(`DELETE FROM events WHERE id = ?`, eventID)
	})

	var (
		success int64
		soldOut int64
		start2  = make(chan struct{})
		wg      sync.WaitGroup
	)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start2 // lepas serempak biar benar-benar balapan
			err := repo.DecrementQuota(ctx, p.DB.DB, eventID, 1)
			switch {
			case err == nil:
				atomic.AddInt64(&success, 1)
			case errors.Is(err, eventRepo.ErrInsufficientQuota):
				atomic.AddInt64(&soldOut, 1)
			default:
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}
	close(start2)
	wg.Wait()

	if success != stock {
		t.Fatalf("sukses decrement = %d, want %d (tepat sejumlah stok)", success, stock)
	}
	if soldOut != workers-stock {
		t.Fatalf("ditolak sold out = %d, want %d", soldOut, workers-stock)
	}

	var remaining int64
	p.DB.Raw(`SELECT quota_remaining FROM events WHERE id = ?`, eventID).Scan(&remaining)
	if remaining != 0 {
		t.Fatalf("quota_remaining akhir = %d, want 0 (tidak boleh negatif / oversell)", remaining)
	}

	t.Logf("OK: %d worker rebutan %d stok -> %d sukses, %d ditolak, sisa=%d",
		workers, stock, success, soldOut, remaining)
}
