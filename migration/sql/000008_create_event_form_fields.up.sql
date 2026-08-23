CREATE TABLE event_form_fields (
    id          BIGSERIAL PRIMARY KEY,
    event_id    BIGINT       NOT NULL,
    label       VARCHAR(255) NOT NULL,
    field_type  VARCHAR(50)  NOT NULL,
    required    BOOLEAN      NOT NULL DEFAULT FALSE,
    options     JSONB,
    position    INT          NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    -- field adalah komposisi dari event: ikut terhapus saat event dihapus.
    CONSTRAINT fk_event_form_fields_event FOREIGN KEY (event_id)
        REFERENCES events (id) ON DELETE CASCADE,

    CONSTRAINT ck_event_form_fields_type CHECK (field_type IN ('text', 'select', 'checkbox')),

    -- field pilihan wajib punya options berupa array tidak kosong; text tidak.
    CONSTRAINT ck_event_form_fields_options CHECK (
        field_type = 'text'
        OR (options IS NOT NULL
            AND jsonb_typeof(options) = 'array'
            AND jsonb_array_length(options) > 0)
    )
);

CREATE INDEX idx_event_form_fields_event_id ON event_form_fields (event_id);
