CREATE TABLE user_event_tickets (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    event_detail_id BIGINT NOT NULL REFERENCES event_details(id),
    seat_id BIGINT NOT NULL REFERENCES seats(id),
    ticket_code VARCHAR(255) NOT NULL,
    is_redeem BOOLEAN NOT NULL DEFAULT FALSE,
    redeemed_at TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);
alter table user_event_tickets ADD COLUMN seat_id BIGINT NOT NULL REFERENCES seats(id);
CREATE INDEX idx_user_event_ticket_user_id ON user_event_tickets(user_id);
CREATE INDEX idx_user_event_ticket_event_detail_id ON user_event_tickets(event_detail_id);