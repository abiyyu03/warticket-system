CREATE TABLE transactions (
    id         BIGSERIAL PRIMARY KEY,
    trx_id     VARCHAR(255)   NOT NULL,
    user_id    BIGINT         NOT NULL,
    event_id   BIGINT         NOT NULL,
    amount     NUMERIC(12, 2) NOT NULL,
    status     VARCHAR(50)    NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ    NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_transactions_user  FOREIGN KEY (user_id)  REFERENCES users (id)  ON DELETE RESTRICT,
    CONSTRAINT fk_transactions_event FOREIGN KEY (event_id) REFERENCES events (id) ON DELETE RESTRICT,

    -- Kunci idempotensi. Payment gateway melakukan retry pada callback;
    -- tanpa UNIQUE ini satu pembayaran bisa terproses berkali-kali.
    CONSTRAINT uq_transactions_trx_id UNIQUE (trx_id),

    CONSTRAINT ck_transactions_amount_non_negative CHECK (amount >= 0),
    CONSTRAINT ck_transactions_status CHECK (status IN ('PENDING', 'PAID', 'FAILED', 'EXPIRED', 'REFUNDED'))
);

CREATE INDEX idx_transactions_user_id  ON transactions (user_id);
CREATE INDEX idx_transactions_event_id ON transactions (event_id);
CREATE INDEX idx_transactions_status   ON transactions (status);
