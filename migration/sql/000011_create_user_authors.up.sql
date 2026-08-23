-- Akun author (penyelenggara) terpisah dari users (pembeli). Login via JWT.
CREATE TABLE user_authors (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    email      VARCHAR(255) NOT NULL,
    password   VARCHAR(255) NOT NULL, -- hash bcrypt
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_user_authors_email UNIQUE (email)
);
