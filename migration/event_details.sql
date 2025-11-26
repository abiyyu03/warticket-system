CREATE TABLE event_details (
    id BIGSERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL REFERENCES events(id),
    ticket_code VARCHAR(255) NOT NULL UNIQUE,
    ticket_type VARCHAR(50) NOT NULL DEFAULT 'REGULAR',
    price NUMERIC(12,2) NOT NULL CHECK(price >= 0),
    available_slot BIGINT NOT NULL CHECK(available_slot > 0),
    status VARCHAR(50) NOT NULL DEFAULT 'active',

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_event_detail_event_id ON event_details(event_id);
CREATE INDEX idx_event_detail_status ON event_details(status);