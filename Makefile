MIGRATE_DIR := migration/sql

.PHONY: migrate-up migrate-down migrate-step-up migrate-step-down migrate-version migrate-force migrate-create migrate-reset run build

## migrate-up: jalankan seluruh migration yang belum dipakai
migrate-up:
	go run ./cmd/migrate up

## migrate-down: rollback seluruh migration
migrate-down:
	go run ./cmd/migrate down

## migrate-step-up: maju N langkah, contoh: make migrate-step-up n=1
migrate-step-up:
	go run ./cmd/migrate up $(n)

## migrate-step-down: mundur N langkah, contoh: make migrate-step-down n=1
migrate-step-down:
	go run ./cmd/migrate down $(n)

## migrate-version: tampilkan versi migration aktif
migrate-version:
	go run ./cmd/migrate version

## migrate-force: tandai versi sebagai bersih, contoh: make migrate-force n=3
migrate-force:
	go run ./cmd/migrate force $(n)

## migrate-reset: turunkan semua lalu naikkan lagi dari nol
migrate-reset:
	go run ./cmd/migrate down
	go run ./cmd/migrate up

## migrate-create: buat pasangan file up/down baru, contoh: make migrate-create name=add_refund_column
migrate-create:
ifndef name
	$(error name wajib diisi, contoh: make migrate-create name=add_refund_column)
endif
	@next=$$(printf "%06d" $$(( $$(ls $(MIGRATE_DIR) 2>/dev/null \
		| sed -n 's/^\([0-9]\{6\}\)_.*/\1/p' | sort -n | tail -1 | sed 's/^0*//' | grep . || echo 0) + 1 ))); \
	touch $(MIGRATE_DIR)/$${next}_$(name).up.sql $(MIGRATE_DIR)/$${next}_$(name).down.sql; \
	echo "dibuat:"; \
	echo "  $(MIGRATE_DIR)/$${next}_$(name).up.sql"; \
	echo "  $(MIGRATE_DIR)/$${next}_$(name).down.sql"

## run: jalankan aplikasi
run:
	go run ./app

## build: compile binary aplikasi dan migrator
build:
	go build -o bin/app ./app
	go build -o bin/migrate ./cmd/migrate
