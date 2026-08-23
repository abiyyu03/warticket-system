-- email wajib pada pendaftaran: karena belum ada login, email menjadi identitas
-- peserta sekaligus alamat tujuan distribusi tiket. Terpisah dari field kustom
-- author sehingga selalu dijamin ada.
ALTER TABLE user_registrations ADD COLUMN email VARCHAR(255);

-- backfill baris lama (dev) supaya bisa di-set NOT NULL.
UPDATE user_registrations SET email = '' WHERE email IS NULL;

ALTER TABLE user_registrations ALTER COLUMN email SET NOT NULL;

CREATE INDEX idx_user_registrations_email ON user_registrations (email);
