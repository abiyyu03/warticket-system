CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    cover VARCHAR(255) NOT NULL,
    location_id BIGINT NOT NULL REFERENCES locations(id),
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_event_name ON events(name);
CREATE INDEX idx_event_cover ON events(cover);
CREATE INDEX idx_event_status ON events(status);
CREATE INDEX idx_event_start_date ON events(start_date);
