-- Membuang skema lama yang dibuat manual sebelum ada migration tracking.
-- Aman dijalankan di database kosong (IF EXISTS) maupun di database lama.
-- Urutan drop mengikuti dependency FK dari anak ke induk.

DROP TABLE IF EXISTS user_event_tickets CASCADE;
DROP TABLE IF EXISTS event_details CASCADE;
DROP TABLE IF EXISTS seats CASCADE;
DROP TABLE IF EXISTS events CASCADE;
DROP TABLE IF EXISTS locations CASCADE;
DROP TABLE IF EXISTS users CASCADE;
