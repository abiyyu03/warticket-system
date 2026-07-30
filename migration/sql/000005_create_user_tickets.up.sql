CREATE TABLE user_tickets (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT       NOT NULL,
    event_id    BIGINT       NOT NULL,
    code        VARCHAR(255) NOT NULL,
    status      VARCHAR(50)  NOT NULL DEFAULT 'ACTIVE',
    valid_until TIMESTAMPTZ  NOT NULL,

    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_user_tickets_user  FOREIGN KEY (user_id)  REFERENCES users (id)  ON DELETE RESTRICT,
    CONSTRAINT fk_user_tickets_event FOREIGN KEY (event_id) REFERENCES events (id) ON DELETE RESTRICT,

    -- Redeem melakukan lookup lewat kolom ini. UNIQUE sekaligus menyediakan
    -- index-nya dan mencegah satu kode menandai banyak tiket.
    CONSTRAINT uq_user_tickets_code UNIQUE (code),

    CONSTRAINT ck_user_tickets_status CHECK (status IN ('ACTIVE', 'REDEEMED', 'EXPIRED', 'CANCELLED'))
);

CREATE INDEX idx_user_tickets_user_id  ON user_tickets (user_id);
CREATE INDEX idx_user_tickets_event_id ON user_tickets (event_id);
CREATE INDEX idx_user_tickets_status   ON user_tickets (status);
