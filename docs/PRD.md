# PRD — WarTicket System

> **Status dokumen:** Draft / hidup (mengikuti kode yang masih WIP)
> **Terakhir diperbarui:** 2026-07-31
> **Pemilik:** Tim WarTicket

---

## 1. Ringkasan

WarTicket adalah backend penjualan tiket event dengan karakteristik **high-contention** —
banyak pembeli merebut kuota terbatas pada waktu yang hampir bersamaan (fenomena "war tiket").
Masalah inti bukan fitur, melainkan **contention**: menjaga agar kuota tidak terjual melebihi batas
(oversell) di tengah lonjakan request serentak, tanpa mengorbankan latensi.

Sistem menyediakan:

- Registrasi & data user.
- **Pengelolaan event oleh author** (penyelenggara): membuat & mengatur event beserta kuota & harga.
- Alur pembelian tiket berbasis **reservasi + pembayaran + callback**.
- Penukaran (*redeem*) tiket saat hari-H.
- (Rencana) integrasi payment gateway dan notifikasi email.

Prinsip desain kuota: **Redis sebagai gerbang admission (authoritative anti-oversell) + reservasi
ber-TTL**, didukung **PostgreSQL sebagai catatan kebenaran final** dengan decrement atomic.

---

## 2. Latar Belakang & Masalah

Penjualan tiket event populer menghadapi lonjakan traffic ekstrem saat penjualan dibuka. Tantangan:

1. **Anti-oversell** — kuota tidak boleh terjual melebihi batas walau request serentak.
2. **Idempotensi pembayaran** — payment gateway retry callback (at-least-once); satu pembayaran
   tidak boleh terproses ganda.
3. **Reservasi sementara** — pembeli butuh jeda menyelesaikan pembayaran tanpa kehilangan slot,
   tapi slot tak boleh terkunci selamanya (cart abandonment normal 60–80%).
4. **Kompensasi** — slot yang direservasi tapi tidak jadi dibayar **harus dikembalikan**, kalau
   tidak kuota bocor dan event tampak sold out padahal masih ada.
5. **Latensi** — keputusan admission harus cepat; tidak semua bisa dibebankan ke database.

---

## 3. Tujuan & Non-Tujuan

### Tujuan
- Pembeli dapat memesan & membeli tiket tanpa risiko oversell.
- Author dapat membuat & mengelola event miliknya.
- Setiap transaksi tercatat dengan status jelas dan dapat diaudit.
- Kuota terkelola cepat via Redis, dengan database sebagai kebenaran final.

### Non-Tujuan (iterasi ini)
- Pemilihan kursi spesifik (seat map). Model kursi lama **ditinggalkan**; tiket berbasis **kuota**.
- Multi-currency.
- Panel admin platform (moderasi lintas author).

---

## 4. Persona & Peran

| Peran | Deskripsi | Kapabilitas utama |
|---|---|---|
| **Buyer (User)** | Pembeli tiket | Register, init order, purchase, lihat tiket, redeem |
| **Author (Organizer)** | Penyelenggara event | CRUD event miliknya, atur kuota & harga, lihat penjualan |
| **Payment Gateway** | Sistem eksternal | Memproses pembayaran, mengirim callback status (harus diverifikasi) |
| **(Rencana) Platform Admin** | Pengelola platform | Moderasi, laporan lintas author |

---

## 5. Lingkup Fitur

### 5.1 Autentikasi & User — *sebagian ada*
- Registrasi user (email unik, password, nama), flag `is_banned`.
- **Rencana:** auth (JWT/session). Sementara identitas dibawa lewat header `x-user-id` (belum aman).

### 5.2 Pengelolaan Event oleh Author — *direncanakan*
- Create event: nama, deskripsi, gambar (`image_file`, multipart), harga, kuota, tanggal mulai/selesai.
  - Saat dibuat, counter kuota Redis `tickets:event:<id>` diinisialisasi dari `quota`.
- List/detail event milik author; update (dengan sinkronisasi counter Redis); nonaktifkan/hapus.
- Lihat penjualan: transaksi & tiket terbit untuk event miliknya.
- **Otorisasi kepemilikan**: author hanya boleh menyentuh event miliknya.

> **Gap:** tabel `events` **belum punya `author_id`**. `transactions.author_id` sudah ada tapi
> diisi placeholder. Fitur ini mensyaratkan `events.author_id` (FK → users). Lihat §12.

### 5.3 Alur Transaksi Tiket — *inti sistem* (detail di §6)
Tiga fase: **Init → Purchase (PENDING) → Callback (SUCCESSFUL/CANCELLED)**. Tiket **diterbitkan
saat callback sukses**, bukan saat purchase.

### 5.4 Redeem Tiket — *WIP*
- Tukar tiket berdasarkan `code`; status `ACTIVE` → `REDEEMED`, anti-double-redeem di level DB.

### 5.5 Notifikasi — *rencana*
- Email konfirmasi & e-ticket (`externals/mail` masih stub).

---

## 6. Alur Transaksi & Kuota (Authoritative)

### 6.1 Model kuota (dua lapis)
- **Redis `tickets:event:<id>` = gerbang admission.** Di-`DECR` saat init. Kalau hasil `< 0`,
  reservasi ditolak (**sold out**). Inilah penahan oversell di lini depan.
- **DB `events.quota` = catatan final.** Di-decrement saat callback sukses dengan UPDATE atomic
  bersyarat; berfungsi sebagai jaring pengaman kedua (CHECK `quota >= 0`).

### 6.2 Fase

**Fase 1 — Init Order** (`POST /v1/api/tickets/init-order`)
- **Idempotent per (user, event):** kalau reservasi aktif sudah ada, jangan `DECR` lagi.
- `DECR` Redis sebanyak `quantity`. Kalau `< 0` → `INCR` balik & tolak (sold out).
- Validasi tanggal dalam rentang event.
- Simpan payload order ke cache **ber-TTL** (mis. 10 menit) sebagai reservasi sementara.

**Fase 2 — Purchase** (`POST /v1/api/tickets/claim`)
- Ambil payload reservasi dari cache.
- **Persist transaksi `PENDING` (dengan `tx_id`) DULU**, baru panggil payment gateway
  (hindari payment yatim kalau crash).
- Catat request/response gateway ke `gateway_requests`.
- Kembalikan status `PENDING` + instruksi/URL bayar.
- *Belum menerbitkan tiket, belum decrement DB quota.*

**Fase 3a — Callback Sukses**
- **Verifikasi signature/HMAC** dari gateway (tolak kalau tidak sah).
- **Idempoten:** `UPDATE transactions SET status='SUCCESSFUL' WHERE tx_id=? AND status='PENDING'`.
  Lanjut **hanya jika `RowsAffected == 1`** (callback ulang = no-op).
- Dalam **satu DB transaction**:
  - Decrement kuota atomic: `UPDATE events SET quota = quota - $qty WHERE id=? AND quota >= $qty`.
  - Terbitkan `qty` baris `user_tickets` (status `ACTIVE`).
- **Safety net:** kalau decrement gagal (`RowsAffected == 0`, artinya drift/oversold) padahal
  sudah dibayar → **REFUND** (jalur kompensasi wajib) dan tandai transaksi sesuai.

**Fase 3b — Callback Gagal**
- Verifikasi signature. `UPDATE ... status='CANCELLED' WHERE tx_id=? AND status='PENDING'`.
- **`INCR` Redis quota balik** sebanyak `qty` (kembalikan reservasi).

**Fase 4 — Sweeper (rencana)**
- Transaksi `PENDING` yang melewati batas waktu (gateway tak pernah callback) → `CANCELLED`
  + `INCR` Redis quota balik. Mencegah kebocoran kuota.

### 6.3 Diagram

```
Buyer            Service              Redis                 Postgres            Gateway
 | init-order ----->|                   |                      |                  |
 |                  |-- DECR(qty) ------>| (<0? INCR & tolak)   |                  |
 |                  |-- validate event ------------------------>|                  |
 |                  |-- cache order(TTL)>|                      |                  |
 | <- INITIATED ----|                   |                      |                  |
 |                  |                   |                      |                  |
 | claim ---------->|                   |                      |                  |
 |                  |-- insert tx PENDING --------------------->|                  |
 |                  |-- request bayar --------------------------------------------->|
 |                  |-- log gateway_requests ------------------>|                  |
 | <- PENDING ------|                   |                      |                  |
 |                  |                   |                      |   callback(sukses)|
 |                  |<---------------------------------------------------------- verify sig
 |                  |-- PENDING->SUCCESSFUL (WHERE status=PENDING, RowsAffected==1)|
 |                  |-- BEGIN tx: UPDATE quota WHERE quota>=qty; insert tickets; COMMIT
 |                  |   (gagal & sudah bayar -> REFUND)          |                  |
 |                  |                   |                      |   callback(gagal) |
 |                  |-- PENDING->CANCELLED; INCR(qty) --------->|                  |
```

### 6.4 Invariant yang dijaga
- Jumlah tiket `ACTIVE`+`REDEEMED` untuk sebuah event **≤ kuota awal**.
- Satu `tx_id` menghasilkan **paling banyak satu** transisi terminal (SUCCESSFUL/CANCELLED).
- Setiap `DECR` Redis punya pasangan `INCR` pada jalur batal/expire (tidak ada kebocoran).
- Callback tidak pernah menerbitkan tiket dua kali untuk `tx_id` yang sama.

---

## 7. Model Data

Skema dikelola via migration bernomor (`golang-migrate`, `migration/sql`).

### `users`
`id`, `name`, `email` (UNIQUE), `password`, `is_banned`, timestamps.

### `events`
`id`, `name`, `description`, `image_file`, `price` (≥0), `quota` (≥0), `start_date`, `end_date`,
timestamps. CHECK: `end_date > start_date` (bila ada), `price ≥ 0`, `quota ≥ 0` (agar bisa sold out).
- **Rencana:** `author_id` (FK → users).

### `transactions`
`id`, `tx_id` (UNIQUE — idempotensi), `user_id`, `event_id`, `author_id`, `status`, `amount`,
`amount_deduction` (nullable), `promo_id` (nullable), `tax`, `admin_fee`, `payment_at` (nullable),
timestamps.
- **`status ∈ {PENDING, SUCCESSFUL, CANCELLED}`** (expiry & refund dipetakan ke `CANCELLED`).
- FK `user_id`, `event_id`, `author_id` → users/events (ON DELETE RESTRICT).

### `user_tickets`
`id`, `user_id`, `event_id`, `code` (UNIQUE), `status`, `valid_until`, timestamps.
- `status ∈ {ACTIVE, REDEEMED, EXPIRED, CANCELLED}`.

### `gateway_requests`
`id`, `providers`, `request`, `response`, `response_code`, `request_header`, `response_header`,
timestamps. Body & header `TEXT` agar response malformed/non-JSON tetap terekam.

### Kunci Redis
- `tickets:event:<event_id>` — counter kuota (SET saat create event; `DECR` init; `INCR` batal/expire).
- `tickets:order:<user_id>:event:<event_id>` — reservasi order (payload + qty), ber-TTL.

---

## 8. Kontrak API

| Method | Path | Deskripsi | Status |
|---|---|---|---|
| GET | `/health-check` | Health check | ✅ |
| POST | `/v1/api/register` | Registrasi user | ✅ |
| GET | `/v1/api/users` | Daftar user | ✅ |
| POST | `/v1/api/events` | Buat event (multipart, upload gambar) | ✅ |
| POST | `/v1/api/tickets/init-order` | Reservasi + DECR kuota Redis | ✅ |
| POST | `/v1/api/tickets/claim` | Buat transaksi PENDING + ke gateway | ⚠️ gateway distub |
| POST | `/v1/api/tickets/redeem` | Tukar tiket | 🚧 WIP |
| — | *(rencana)* `/v1/api/payments/callback` | Callback status (verifikasi signature) | 🗓️ |
| — | *(rencana)* endpoint event milik author (list/update/delete) | Author mgmt | 🗓️ |

Catatan: identitas user sementara lewat header `x-user-id` (belum ada auth middleware).

---

## 9. Kebutuhan Non-Fungsional

- **Anti-oversell (2 lapis):** `DECR` atomic Redis sebagai gerbang + `quota >= qty` UPDATE atomic DB.
- **Idempotensi pembayaran:** `tx_id` UNIQUE + transisi status bersyarat (`WHERE status='PENDING'`,
  cek `RowsAffected`).
- **Keamanan callback:** verifikasi signature/HMAC gateway sebelum memproses.
- **Kompensasi kuota:** setiap `DECR` init punya pasangan `INCR` pada batal/gagal/expire.
- **Konsistensi:** update transaksi + decrement kuota + terbit tiket dalam satu DB transaction.
- **Refund safety net:** callback sukses tapi kuota DB tak cukup → refund otomatis.
- **Anti-double-redeem:** update `user_tickets` difilter status di level DB.
- **Auditability:** semua panggilan gateway tercatat di `gateway_requests`.
- **Timezone-safe:** seluruh kolom waktu `TIMESTAMPTZ`.

---

## 10. Arsitektur

- **Bahasa/Framework:** Go + Fiber (HTTP), GORM (ORM).
- **Datastore:** PostgreSQL (kebenaran final) + Redis (gerbang kuota & reservasi cepat).
- **Pola:** Hexagonal (ports & adapters): `inbound/rest` → handler; `service` → use case;
  `outbound/repository` (Postgres), `outbound/cache` (Redis), `outbound/externals` (payment, mail).
- **DI:** Uber `dig`. **Migration:** `golang-migrate`, di-embed, dijalankan via `cmd/migrate` / `Makefile`.

---

## 11. Roadmap & Status

| Fase | Cakupan | Status |
|---|---|---|
| **0 — Fondasi** | Skema baru, migration tooling, entity, DI | ✅ Sebagian besar |
| **1 — Event & Reservasi** | Create event + quota Redis, init order | ✅ |
| **2 — Purchase & Callback** | Transaksi PENDING, callback SUCCESSFUL/CANCELLED, decrement DB atomic, kompensasi | ⚠️ purchase ada; callback belum |
| **3 — Redeem** | Tukar tiket hari-H | 🚧 WIP |
| **4 — Author Management** | `events.author_id`, CRUD event milik author, otorisasi, dashboard | 🗓️ |
| **5 — Payment Gateway** | Integrasi provider, verifikasi signature, logging, refund | 🗓️ |
| **6 — Notifikasi** | Email konfirmasi & e-ticket | 🗓️ |
| **7 — Promo & Pajak** | `promo_id`, `amount_deduction`, `tax`, `admin_fee` end-to-end | 🗓️ |
| **8 — Sweeper** | Expire PENDING + kembalikan kuota | 🗓️ |

---

## 12. Risiko & Pertanyaan Terbuka

1. **`events.author_id` belum ada.** Prasyarat Fase 4; `transactions.author_id` kini placeholder.
   Perlu migration tambahan + backfill.
2. **Enum status transaksi belum sinkron.** PRD menetapkan `PENDING/SUCCESSFUL/CANCELLED`, tetapi
   migration `000004` + konstanta entity masih `PAID/FAILED/EXPIRED/REFUNDED`. Perlu diselaraskan.
3. **Tiket masih diterbitkan saat purchase.** Implementasi `purchase.go` sekarang menerbitkan tiket
   di fase claim; sesuai alur otoritatif ini harus **dipindah ke callback sukses**.
4. **Callback belum ada.** Butuh endpoint + verifikasi signature + idempotensi + refund path.
5. **Kompensasi & sweeper belum ada.** Tanpa `INCR` balik dan expiry, kuota Redis bocor.
6. **Autentikasi.** `x-user-id` via header belum aman; butuh auth sebelum produksi.
7. **Model `promo`.** `promo_id` nullable tanpa tabel `promos`; FK menyusul.
8. **`gateway_requests` belum tertaut transaksi.** Pertimbangkan `tx_id` untuk korelasi.
