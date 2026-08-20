-- quota_remaining: stok yang bisa di-decrement (jalur transaksi), sementara
-- quota tetap sebagai kapasitas asli (immutable). sold = quota - quota_remaining.
ALTER TABLE events ADD COLUMN quota_remaining BIGINT NOT NULL DEFAULT 0;

-- backfill baris lama: sisa = kapasitas awal.
UPDATE events SET quota_remaining = quota;

-- backstop anti-oversell di level DB.
ALTER TABLE events ADD CONSTRAINT ck_events_quota_remaining_non_negative
    CHECK (quota_remaining >= 0);

-- sisa tidak boleh melebihi kapasitas.
ALTER TABLE events ADD CONSTRAINT ck_events_quota_remaining_le_quota
    CHECK (quota_remaining <= quota);
