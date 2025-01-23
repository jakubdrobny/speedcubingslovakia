BEGIN;

CREATE TABLE IF NOT EXISTS results_status (
  results_status_id BIGSERIAL PRIMARY KEY,
  approvalfinished BOOLEAN NOT NULL,
  approved BOOLEAN NOT NULL,
  visible BOOLEAN NOT NULL,
  displayname TEXT NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT results_status_unique UNIQUE (approvalfinished, approved, visible, displayname)
);

COMMIT;
