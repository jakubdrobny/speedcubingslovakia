/* create tables if neccessary */

CREATE TABLE IF NOT EXISTS continents (
    continent_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    recordName TEXT NOT NULL,
    CONSTRAINT continents_unique UNIQUE (continent_id, name, recordName)
);

CREATE TABLE IF NOT EXISTS countries (
  country_id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  continent_id TEXT REFERENCES continents (continent_id) NOT NULL,
  iso2 TEXT NOT NULL,
  CONSTRAINT countries_unique UNIQUE (country_id, name, continent_id, iso2)
);

CREATE TABLE IF NOT EXISTS users (
  user_id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  country_id TEXT REFERENCES countries (country_id) NOT NULL,
  sex TEXT NOT NULL,
  wcaid TEXT NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  isadmin BOOLEAN NOT NULL,
  url TEXT NOT NULL,
  avatarurl TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS competitions (
  competition_id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  startdate TIMESTAMP NOT NULL,
  enddate TIMESTAMP NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS events (
  event_id BIGSERIAL PRIMARY KEY,
  fulldisplayname TEXT NOT NULL,
  displayname TEXT NOT NULL,
  format TEXT NOT NULL,
  iconcode TEXT NOT NULL,
  scramblingcode TEXT NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT event_unique UNIQUE (fulldisplayname, displayname, format, iconcode, scramblingcode)
);

CREATE TABLE IF NOT EXISTS competition_events (
  competition_events_id BIGSERIAL PRIMARY KEY,
  competition_id TEXT REFERENCES competitions (competition_id) NOT NULL,
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
  result_id BIGSERIAL PRIMARY KEY,
  competition_id TEXT REFERENCES competitions (competition_id) NOT NULL,
  user_id INTEGER REFERENCES users (user_id) NOT NULL,
  event_id INTEGER REFERENCES events (event_id) NOT NULL,
  solve1 TEXT NOT NULL,
  solve2 TEXT NOT NULL,
  solve3 TEXT NOT NULL,
  solve4 TEXT NOT NULL,
  solve5 TEXT NOT NULL,
  comment TEXT NOT NULL,
  status_id INTEGER REFERENCES results_status (results_status_id) NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS scrambles (
  scramble_id BIGSERIAL PRIMARY KEY,
  scramble TEXT NOT NULL,
  event_id INTEGER REFERENCES events (event_id) NOT NULL,
  competition_id TEXT REFERENCES competitions (competition_id) NOT NULL,
  "order" INTEGER NOT NULL,
  img TEXT NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tags (
    tag_id BIGSERIAL PRIMARY KEY,
    label TEXT NOT NULL,
    color TEXT NOT NULL CHECK (color IN ('primary', 'warning', 'success', 'danger')),
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS announcements (
    announcement_id BIGSERIAL PRIMARY KEY,
    author_id INTEGER REFERENCES users (user_id) NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS announcement_tags (
    announcement_tags_id BIGSERIAL PRIMARY KEY,
    announcement_id INTEGER REFERENCES announcements (announcement_id) NOT NULL,
    tag_id INTEGER REFERENCES tags (tag_id) NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS announcement_read (
    announcement_read_id BIGSERIAL PRIMARY KEY,
    announcement_id INTEGER REFERENCES announcements (announcement_id) NOT NULL,
    user_id INTEGER REFERENCES users (user_id) NOT NULL,
    read BOOLEAN NOT NULL DEFAULT FALSE,
    read_timestamp TIMESTAMP,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS announcement_reaction (
    announcement_reaction_id BIGSERIAL PRIMARY KEY,
    announcement_id INTEGER REFERENCES announcements (announcement_id) ON DELETE CASCADE NOT NULL,
    user_id INTEGER REFERENCES users (user_id) ON DELETE CASCADE NOT NULl,
    emoji TEXT NOT NULL,
    "by" TEXT NOT NULL,
    "set" BOOLEAN NOT NULL
);

/* insert stock data */


/* events */
INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('3x3x3 Cube', '3x3x3', 'ao5', '333', '333')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('2x2x2 Cube', '2x2x2', 'ao5', '222', '222so')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('4x4x4 Cube', '4x4x4', 'ao5', '444', '444wca')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('5x5x5 Cube', '5x5x5', 'ao5', '555', '555wca')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('6x6x6 Cube', '6x6x6', 'mo3', '666', '666wca')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('7x7x7 Cube', '7x7x7', 'mo3', '777', '777wca')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('3x3x3 Blindfolded', '3BLD', 'bo3', '333bf', '333ni')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('3x3x3 Fewest Moves', 'FMC', 'mo3', '333fm', '333fm')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('3x3x3 One-Handed', 'OH', 'ao5', '333oh', '333oh')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('Clock', 'Clock', 'ao5', 'clock', 'clkwca')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('Megaminx', 'Mega', 'ao5', 'minx', 'mgmp')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('Pyraminx', 'Pyra', 'ao5', 'pyram', 'pyrso')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('Skewb', 'Skewb', 'ao5', 'skewb', 'skbso')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('Square-1', 'Sq-1', 'ao5', 'sq1', 'sqrs')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('4x4x4 Blindfolded', '4BLD', 'bo3', '444bf', '444bld')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('5x5x5 Blindfolded', '5BLD', 'bo3', '555bf', '555bld')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('3x3x3 Multi-Blind', 'MBLD', 'bo3', '333mbf', '333ni')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('2x2x2 Blindfolded', '2BLD', 'bo3', 'unofficial-222bf', '222so')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('6x6x6 Blindfolded', '6BLD', 'bo1', 'unofficial-666bf', '666wca')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('7x7x7 Blindfolded', '7BLD', 'bo1', 'unofficial-777bf', '777wca')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('3x3x3 With Feet', 'Feet', 'ao5', '333ft', '333')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('3x3x3 Match The Scramble', 'Match', 'ao5', 'unofficial-333mts', '333')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('2x2 - 4x4 Relay', '2-4 Relay', 'bo1', 'unofficial-234relay', 'r234w')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('2x2 - 5x5 Relay', '2-5 Relay', 'bo1', 'unofficial-2345relay', 'r2345w')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('2x2 - 6x6 Relay', '2-6 Relay', 'bo1', 'unofficial-23456relay', 'r23456w')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('2x2 - 7x7 Relay', '2-7 Relay', 'bo1', 'unofficial-234567relay', 'r234567w')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('Kilominx', 'Kilo', 'ao5', 'unofficial-kilominx', 'klmso')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('Mini Guildford', 'Mini Guild', 'bo1', 'unofficial-miniguild', 'rmngf')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('Redi Cube', 'Redi', 'ao5', 'unofficial-redi', 'rediso')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('Master Pyraminx', 'Master Pyra', 'ao5', 'unofficial-mpyram', 'mpyrso')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('15 Puzzle', '15 Puzzle', 'ao5', 'unofficial-15puzzle', '15prp')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('Mirror Blocks', 'Mirror', 'ao5', 'unofficial-mirror', '333')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

INSERT INTO events (fulldisplayname, displayname, format, iconcode, scramblingcode)
VALUES ('Face Turning Octahedron', 'FTO', 'ao5', 'unofficial-fto', 'ftoso')
ON CONFLICT (fulldisplayname, displayname, format, iconcode, scramblingcode) DO NOTHING;

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




/* continents */
INSERT INTO continents (continent_id, name, recordName)
VALUES ('_Africa','Africa','AfR')
ON CONFLICT (continent_id, name, recordName) DO NOTHING;

INSERT INTO continents (continent_id, name, recordName)
VALUES ('_Asia','Asia','AsR')
ON CONFLICT (continent_id, name, recordName) DO NOTHING;

INSERT INTO continents (continent_id, name, recordName)
VALUES ('_Europe','Europe','ER')
ON CONFLICT (continent_id, name, recordName) DO NOTHING;

INSERT INTO continents (continent_id, name, recordName)
VALUES ('_North America','North America','NAR')
ON CONFLICT (continent_id, name, recordName) DO NOTHING;

INSERT INTO continents (continent_id, name, recordName)
VALUES ('_Oceania','Oceania','OcR')
ON CONFLICT (continent_id, name, recordName) DO NOTHING;

INSERT INTO continents (continent_id, name, recordName)
VALUES ('_South America','South America','SAR')
ON CONFLICT (continent_id, name, recordName) DO NOTHING;


/* countries */
INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Afghanistan','Afghanistan','_Asia','AF')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Albania','Albania','_Europe','AL')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Algeria','Algeria','_Africa','DZ')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Andorra','Andorra','_Europe','AD')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Angola','Angola','_Africa','AO')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Antigua and Barbuda','Antigua and Barbuda','_North America','AG')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Argentina','Argentina','_South America','AR')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Armenia','Armenia','_Europe','AM')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Australia','Australia','_Oceania','AU')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Austria','Austria','_Europe','AT')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Azerbaijan','Azerbaijan','_Europe','AZ')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Bahamas','Bahamas','_North America','BS')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Bahrain','Bahrain','_Asia','BH')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Bangladesh','Bangladesh','_Asia','BD')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Barbados','Barbados','_North America','BB')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Belarus','Belarus','_Europe','BY')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Belgium','Belgium','_Europe','BE')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Belize','Belize','_North America','BZ')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Benin','Benin','_Africa','BJ')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Bhutan','Bhutan','_Asia','BT')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Bolivia','Bolivia','_South America','BO')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Bosnia and Herzegovina','Bosnia and Herzegovina','_Europe','BA')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Botswana','Botswana','_Africa','BW')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Brazil','Brazil','_South America','BR')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Brunei','Brunei','_Asia','BN')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Bulgaria','Bulgaria','_Europe','BG')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Burkina Faso','Burkina Faso','_Africa','BF')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Burundi','Burundi','_Africa','BI')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Cabo Verde','Cabo Verde','_Africa','CV')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Cambodia','Cambodia','_Asia','KH')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Cameroon','Cameroon','_Africa','CM')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Canada','Canada','_North America','CA')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Central African Republic','Central African Republic','_Africa','CF')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Chad','Chad','_Africa','TD')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Chile','Chile','_South America','CL')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('China','China','_Asia','CN')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Colombia','Colombia','_South America','CO')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Comoros','Comoros','_Africa','KM')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Congo','Congo','_Africa','CG')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Costa Rica','Costa Rica','_North America','CR')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Cote d_Ivoire','Côte d''Ivoire','_Africa','CI')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Croatia','Croatia','_Europe','HR')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Cuba','Cuba','_North America','CU')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Cyprus','Cyprus','_Europe','CY')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Czech Republic','Czech Republic','_Europe','CZ')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Democratic People_s Republic of Korea','Democratic People''s Republic of Korea','_Asia','KP')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Democratic Republic of the Congo','Democratic Republic of the Congo','_Africa','CD')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Denmark','Denmark','_Europe','DK')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Djibouti','Djibouti','_Africa','DJ')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Dominica','Dominica','_North America','DM')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Dominican Republic','Dominican Republic','_North America','DO')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Ecuador','Ecuador','_South America','EC')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Egypt','Egypt','_Africa','EG')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('El Salvador','El Salvador','_North America','SV')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Equatorial Guinea','Equatorial Guinea','_Africa','GQ')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Eritrea','Eritrea','_Africa','ER')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Estonia','Estonia','_Europe','EE')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Eswatini','Eswatini','_Africa','SZ')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Ethiopia','Ethiopia','_Africa','ET')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Federated States of Micronesia','Federated States of Micronesia','_Oceania','FM')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Fiji','Fiji','_Oceania','FJ')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Finland','Finland','_Europe','FI')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('France','France','_Europe','FR')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Gabon','Gabon','_Africa','GA')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Gambia','Gambia','_Africa','GM')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Georgia','Georgia','_Europe','GE')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Germany','Germany','_Europe','DE')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Ghana','Ghana','_Africa','GH')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Greece','Greece','_Europe','GR')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Grenada','Grenada','_North America','GD')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Guatemala','Guatemala','_North America','GT')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Guinea','Guinea','_Africa','GN')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Guinea Bissau','Guinea Bissau','_Africa','GW')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Guyana','Guyana','_South America','GY')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Haiti','Haiti','_North America','HT')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Honduras','Honduras','_North America','HN')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Hong Kong','Hong Kong, China','_Asia','HK')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Hungary','Hungary','_Europe','HU')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Iceland','Iceland','_Europe','IS')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('India','India','_Asia','IN')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Indonesia','Indonesia','_Asia','ID')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Iran','Iran','_Asia','IR')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Iraq','Iraq','_Asia','IQ')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Ireland','Ireland','_Europe','IE')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Israel','Israel','_Europe','IL')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Italy','Italy','_Europe','IT')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Jamaica','Jamaica','_North America','JM')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Japan','Japan','_Asia','JP')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Jordan','Jordan','_Asia','JO')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Kazakhstan','Kazakhstan','_Asia','KZ')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Kenya','Kenya','_Africa','KE')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Kiribati','Kiribati','_Oceania','KI')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Korea','Republic of Korea','_Asia','KR')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Kosovo','Kosovo','_Europe','XK')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Kuwait','Kuwait','_Asia','KW')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Kyrgyzstan','Kyrgyzstan','_Asia','KG')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Laos','Laos','_Asia','LA')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Latvia','Latvia','_Europe','LV')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Lebanon','Lebanon','_Asia','LB')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Lesotho','Lesotho','_Africa','LS')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Liberia','Liberia','_Africa','LR')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Libya','Libya','_Africa','LY')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Liechtenstein','Liechtenstein','_Europe','LI')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Lithuania','Lithuania','_Europe','LT')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Luxembourg','Luxembourg','_Europe','LU')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Macau','Macau, China','_Asia','MO')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Madagascar','Madagascar','_Africa','MG')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Malawi','Malawi','_Africa','MW')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Malaysia','Malaysia','_Asia','MY')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Maldives','Maldives','_Asia','MV')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Mali','Mali','_Africa','ML')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Malta','Malta','_Europe','MT')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Marshall Islands','Marshall Islands','_Oceania','MH')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Mauritania','Mauritania','_Africa','MR')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Mauritius','Mauritius','_Africa','MU')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Mexico','Mexico','_North America','MX')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Moldova','Moldova','_Europe','MD')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Monaco','Monaco','_Europe','MC')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Mongolia','Mongolia','_Asia','MN')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Montenegro','Montenegro','_Europe','ME')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Morocco','Morocco','_Africa','MA')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Mozambique','Mozambique','_Africa','MZ')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Myanmar','Myanmar','_Asia','MM')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Namibia','Namibia','_Africa','NA')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Nauru','Nauru','_Oceania','NR')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Nepal','Nepal','_Asia','NP')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Netherlands','Netherlands','_Europe','NL')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('New Zealand','New Zealand','_Oceania','NZ')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Nicaragua','Nicaragua','_North America','NI')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Niger','Niger','_Africa','NE')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Nigeria','Nigeria','_Africa','NG')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('North Macedonia','North Macedonia','_Europe','MK')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Norway','Norway','_Europe','NO')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Oman','Oman','_Asia','OM')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Pakistan','Pakistan','_Asia','PK')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Palau','Palau','_Oceania','PW')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Palestine','Palestine','_Asia','PS')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Panama','Panama','_North America','PA')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Papua New Guinea','Papua New Guinea','_Oceania','PG')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Paraguay','Paraguay','_South America','PY')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Peru','Peru','_South America','PE')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Philippines','Philippines','_Asia','PH')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Poland','Poland','_Europe','PL')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Portugal','Portugal','_Europe','PT')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Qatar','Qatar','_Asia','QA')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Romania','Romania','_Europe','RO')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Russia','Russia','_Europe','RU')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Rwanda','Rwanda','_Africa','RW')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Saint Kitts and Nevis','Saint Kitts and Nevis','_North America','KN')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Saint Lucia','Saint Lucia','_North America','LC')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Saint Vincent and the Grenadines','Saint Vincent and the Grenadines','_North America','VC')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Samoa','Samoa','_Oceania','WS')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('San Marino','San Marino','_Europe','SM')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Sao Tome and Principe','São Tomé and Príncipe','_Africa','ST')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Saudi Arabia','Saudi Arabia','_Asia','SA')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Senegal','Senegal','_Africa','SN')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Serbia','Serbia','_Europe','RS')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Seychelles','Seychelles','_Africa','SC')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Sierra Leone','Sierra Leone','_Africa','SL')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Singapore','Singapore','_Asia','SG')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Slovakia','Slovakia','_Europe','SK')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Slovenia','Slovenia','_Europe','SI')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Solomon Islands','Solomon Islands','_Oceania','SB')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Somalia','Somalia','_Africa','SO')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('South Africa','South Africa','_Africa','ZA')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('South Sudan','South Sudan','_Africa','SS')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Spain','Spain','_Europe','ES')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Sri Lanka','Sri Lanka','_Asia','LK')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Sudan','Sudan','_Africa','SD')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Suriname','Suriname','_South America','SR')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Sweden','Sweden','_Europe','SE')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Switzerland','Switzerland','_Europe','CH')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Syria','Syria','_Asia','SY')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Taiwan','Chinese Taipei','_Asia','TW')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Tajikistan','Tajikistan','_Asia','TJ')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Tanzania','Tanzania','_Africa','TZ')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Thailand','Thailand','_Asia','TH')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Timor-Leste','Timor-Leste','_Asia','TL')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Togo','Togo','_Africa','TG')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Tonga','Tonga','_Oceania','TO')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Trinidad and Tobago','Trinidad and Tobago','_North America','TT')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Tunisia','Tunisia','_Africa','TN')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Turkey','Turkey','_Europe','TR')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Turkmenistan','Turkmenistan','_Asia','TM')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Tuvalu','Tuvalu','_Oceania','TV')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Uganda','Uganda','_Africa','UG')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Ukraine','Ukraine','_Europe','UA')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('United Arab Emirates','United Arab Emirates','_Asia','AE')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('United Kingdom','United Kingdom','_Europe','GB')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Uruguay','Uruguay','_South America','UY')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('USA','United States','_North America','US')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Uzbekistan','Uzbekistan','_Asia','UZ')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Vanuatu','Vanuatu','_Oceania','VU')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Vatican City','Vatican City','_Europe','VA')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Venezuela','Venezuela','_South America','VE')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Vietnam','Vietnam','_Asia','VN')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Yemen','Yemen','_Asia','YE')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Zambia','Zambia','_Africa','ZM')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;

INSERT INTO countries (country_id, name, continent_id, iso2)
VALUES ('Zimbabwe','Zimbabwe','_Africa','ZW')
ON CONFLICT (country_id, name, continent_id, iso2) DO NOTHING;
