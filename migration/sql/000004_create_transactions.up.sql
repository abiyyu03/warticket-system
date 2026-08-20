CREATE TABLE transactions (
    id               BIGSERIAL PRIMARY KEY,
    tx_id            VARCHAR(255)   NOT NULL,
    user_id          BIGINT         NOT NULL,
    event_id         BIGINT         NOT NULL,
    author_id        BIGINT         NOT NULL,
    status           VARCHAR(50)    NOT NULL DEFAULT 'PENDING',
    amount           NUMERIC(12, 2) NOT NULL,
    amount_deduction NUMERIC(12, 2),
    promo_id         BIGINT,
    tax              NUMERIC(12, 2) NOT NULL DEFAULT 0,
    admin_fee        NUMERIC(12, 2) NOT NULL DEFAULT 0,
    payment_at       TIMESTAMPTZ,

    created_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_transactions_user   FOREIGN KEY (user_id)   REFERENCES users (id)  ON DELETE RESTRICT,
    CONSTRAINT fk_transactions_event  FOREIGN KEY (event_id)  REFERENCES events (id) ON DELETE RESTRICT,
    -- author_id = pemilik/penyelenggara event. Diarahkan ke users karena itu
    -- satu-satunya tabel identitas saat ini.
    CONSTRAINT fk_transactions_author FOREIGN KEY (author_id) REFERENCES users (id)  ON DELETE RESTRICT,

    -- Kunci idempotensi. Payment gateway melakukan retry pada callback;
    -- tanpa UNIQUE ini satu pembayaran bisa terproses berkali-kali.
    CONSTRAINT uq_transactions_tx_id UNIQUE (tx_id),

    CONSTRAINT ck_transactions_amount_non_negative     CHECK (amount >= 0),
    CONSTRAINT ck_transactions_deduction_non_negative  CHECK (amount_deduction IS NULL OR amount_deduction >= 0),
    CONSTRAINT ck_transactions_tax_non_negative        CHECK (tax >= 0),
    CONSTRAINT ck_transactions_admin_fee_non_negative  CHECK (admin_fee >= 0),
    CONSTRAINT ck_transactions_status CHECK (status IN ('PENDING', 'SUCCESSFUL', 'CANCELLED','REFUNDED'))
);

-- promo_id belum diberi FK karena tabel promos belum ada. Tambahkan
-- fk_transactions_promo saat tabel promos dibuat.

CREATE INDEX idx_transactions_user_id   ON transactions (user_id);
CREATE INDEX idx_transactions_event_id  ON transactions (event_id);
CREATE INDEX idx_transactions_author_id ON transactions (author_id);
CREATE INDEX idx_transactions_status    ON transactions (status);
