CREATE TABLE events (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(255)   NOT NULL,
    description TEXT,
    image_file  VARCHAR(255),
    price       NUMERIC(12, 2) NOT NULL,
    quota       BIGINT         NOT NULL,
    start_date  TIMESTAMPTZ    NOT NULL,
    end_date    TIMESTAMPTZ,

    created_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),

    CONSTRAINT ck_events_price_non_negative CHECK (price >= 0),

    -- >= 0, BUKAN > 0. Kalau pakai > 0, UPDATE yang menurunkan kuota terakhir
    -- ke nol akan ditolak dan event tidak akan pernah bisa sold out.
    CONSTRAINT ck_events_quota_non_negative CHECK (quota >= 0),

    CONSTRAINT ck_events_date_range CHECK (end_date IS NULL OR end_date > start_date)
);

CREATE INDEX idx_events_start_date ON events (start_date);
