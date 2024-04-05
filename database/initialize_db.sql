/* create tables if neccessary */

CREATE TABLE IF NOT EXISTS users (
  user_id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  country TEXT NOT NULL,
  sex TEXT NOT NULL,
  wcaid TEXT NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  isadmin BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS competitions (
  competition_id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  startdate TIMESTAMP NOT NULL,
  enddate TIMESTAMP NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS events (
  event_id BIGSERIAL PRIMARY KEY,
  displayname TEXT NOT NULL,
  format TEXT NOT NULL,
  iconcode TEXT NOT NULL,
  puzzlecode TEXT NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT event_unique UNIQUE (displayname, format, iconcode, puzzlecode)
);

CREATE TABLE IF NOT EXISTS competition_events (
  competition_events_id BIGSERIAL PRIMARY KEY,
  competition_id INTEGER REFERENCES competitions (competition_id) NOT NULL,
  event_id INTEGER REFERENCES events (event_id) NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS results_status (
  results_status_id BIGSERIAL PRIMARY KEY,
  approvalfinished BOOLEAN NOT NULL,
  approved BOOLEAN NOT NULL,
  visible BOOLEAN NOT NULL,
  displayname TEXT NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT results_status_unique UNIQUE (approvalfinished, approved, visible, displayname)
);

CREATE TABLE IF NOT EXISTS results (
  id BIGSERIAL PRIMARY KEY,
  competition_id INTEGER REFERENCES competitions (competition_id) NOT NULL,
  user_id INTEGER REFERENCES users (user_id) NOT NULL,
  event_id INTEGER REFERENCES events (event_id) NOT NULL,
  solve1 TEXT,
  solve2 TEXT,
  solve3 TEXT,
  solve4 TEXT,
  solve5 TEXT,
  comment TEXT NOT NULL,
  status_id INTEGER REFERENCES results_status (results_status_id) NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE scrambles (
  scramble_id BIGSERIAL PRIMARY KEY,
  scramble TEXT NOT NULL,
  event_id INTEGER REFERENCES events (event_id) NOT NULL,
  competition_id INTEGER REFERENCES competitions (competition_id) NOT NULL,
  "order" INTEGER NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


/* insert stock data */

/* events */
INSERT INTO events (displayname, format, iconcode, puzzlecode)
VALUES ('3x3x3', 'ao5', '333', '3x3x3')
ON CONFLICT (displayname, format, iconcode, puzzlecode) DO NOTHING;

INSERT INTO events (displayname, format, iconcode, puzzlecode)
VALUES ('2x2x2', 'ao5', '222', '2x2x2')
ON CONFLICT (displayname, format, iconcode, puzzlecode) DO NOTHING;

INSERT INTO events (displayname, format, iconcode, puzzlecode)
VALUES ('4x4x4', 'ao5', '444', '4x4x4')
ON CONFLICT (displayname, format, iconcode, puzzlecode) DO NOTHING;

INSERT INTO events (displayname, format, iconcode, puzzlecode)
VALUES ('5x5x5', 'ao5', '555', '5x5x5')
ON CONFLICT (displayname, format, iconcode, puzzlecode) DO NOTHING;

INSERT INTO events (displayname, format, iconcode, puzzlecode)
VALUES ('6x6x6', 'mo3', '666', '6x6x6')
ON CONFLICT (displayname, format, iconcode, puzzlecode) DO NOTHING;

INSERT INTO events (displayname, format, iconcode, puzzlecode)
VALUES ('7x7x7', 'mo3', '777', '7x7x7')
ON CONFLICT (displayname, format, iconcode, puzzlecode) DO NOTHING;

INSERT INTO events (displayname, format, iconcode, puzzlecode)
VALUES ('3BLD', 'bo3', '333bf', '3x3x3')
ON CONFLICT (displayname, format, iconcode, puzzlecode) DO NOTHING;

INSERT INTO events (displayname, format, iconcode, puzzlecode)
VALUES ('FMC', 'mo3', '333fm', '3x3x3')
ON CONFLICT (displayname, format, iconcode, puzzlecode) DO NOTHING;

INSERT INTO events (displayname, format, iconcode, puzzlecode)
VALUES ('OH', 'ao5', '333oh', '3x3x3')
ON CONFLICT (displayname, format, iconcode, puzzlecode) DO NOTHING;

INSERT INTO events (displayname, format, iconcode, puzzlecode)
VALUES ('Clock', 'ao5', 'clock', 'clock')
ON CONFLICT (displayname, format, iconcode, puzzlecode) DO NOTHING;

INSERT INTO events (displayname, format, iconcode, puzzlecode)
VALUES ('Mega', 'ao5', 'megaminx', 'megaminx')
ON CONFLICT (displayname, format, iconcode, puzzlecode) DO NOTHING;

INSERT INTO events (displayname, format, iconcode, puzzlecode)
VALUES ('Pyra', 'ao5', 'pyraminx', 'pyraminx')
ON CONFLICT (displayname, format, iconcode, puzzlecode) DO NOTHING;

INSERT INTO events (displayname, format, iconcode, puzzlecode)
VALUES ('Skewb', 'ao5', 'skewb', 'skewb')
ON CONFLICT (displayname, format, iconcode, puzzlecode) DO NOTHING;

INSERT INTO events (displayname, format, iconcode, puzzlecode)
VALUES ('Sq-1', 'ao5', 'sq1', 'square1')
ON CONFLICT (displayname, format, iconcode, puzzlecode) DO NOTHING;

INSERT INTO events (displayname, format, iconcode, puzzlecode)
VALUES ('4BLD', 'bo3', '444bf', '4x4x4')
ON CONFLICT (displayname, format, iconcode, puzzlecode) DO NOTHING;

INSERT INTO events (displayname, format, iconcode, puzzlecode)
VALUES ('5BLD', 'bo3', '555bf', '5x5x5')
ON CONFLICT (displayname, format, iconcode, puzzlecode) DO NOTHING;


/* results status */
INSERT INTO results_status (approvalfinished, approved, visible, displayname)
VALUES (false, false, false, 'Waiting for approval')
ON CONFLICT (approvalfinished, approved, visible, displayname) DO NOTHING;

INSERT INTO results_status (approvalfinished, approved, visible, displayname)
VALUES (true, false, false, 'Denied')
ON CONFLICT (approvalfinished, approved, visible, displayname) DO NOTHING;

INSERT INTO results_status (approvalfinished, approved, visible, displayname)
VALUES (true, true, true, 'Approved')
ON CONFLICT (approvalfinished, approved, visible, displayname) DO NOTHING;
