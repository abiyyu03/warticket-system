-- Pendaftaran peserta terhadap formulir kustom sebuah event (§5.3, buyer-side).
-- Satu baris per (user, event); jawaban disimpan JSONB keyed by field_id.
CREATE TABLE user_registrations (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT      NOT NULL,
    event_id    BIGINT      NOT NULL,
    -- {"<field_id>": ["nilai", ...]} — array agar seragam untuk semua tipe field.
    answers     JSONB       NOT NULL DEFAULT '{}',

    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_user_registrations_user  FOREIGN KEY (user_id)
        REFERENCES users (id)  ON DELETE RESTRICT,
    CONSTRAINT fk_user_registrations_event FOREIGN KEY (event_id)
        REFERENCES events (id) ON DELETE CASCADE,

    -- idempotensi: satu peserta hanya boleh mendaftar sekali per event.
    CONSTRAINT uq_user_registrations_user_event UNIQUE (user_id, event_id)
);

CREATE INDEX idx_user_registrations_event_id ON user_registrations (event_id);
