-- Log request/response ke payment gateway untuk audit & debugging.
-- Body dan header disimpan sebagai TEXT (bukan JSONB) supaya response yang
-- malformed/non-JSON (HTML error page, body kepotong) tetap bisa terekam.
CREATE TABLE gateway_requests (
    id              BIGSERIAL PRIMARY KEY,
    providers       VARCHAR(100) NOT NULL,
    request         TEXT,
    response        TEXT,
    response_code   INTEGER,
    request_header  TEXT,
    response_header TEXT,

    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_gateway_requests_providers     ON gateway_requests (providers);
CREATE INDEX idx_gateway_requests_response_code ON gateway_requests (response_code);
