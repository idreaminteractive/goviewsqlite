
CREATE TABLE count (
    id integer  primary key autoincrement,
    value int not null default 0    
);
CREATE TABLE user (
    id INTEGER NOT NULL  PRIMARY KEY AUTOINCREMENT,
    idp_user_id text not null,
    email text not null,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
, default_organization_id INTEGER DEFAULT 0, full_name TEXT NOT NULL DEFAULT '');
CREATE UNIQUE INDEX user_email_idx on user (email);
CREATE TRIGGER user_update_updated_at_trigger
AFTER UPDATE ON user
FOR EACH ROW
BEGIN
    UPDATE user SET updated_at = CURRENT_TIMESTAMP WHERE rowid = NEW.rowid;
END;
CREATE TABLE IF NOT EXISTS "organization" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER update_organization_updated_at
AFTER UPDATE ON organization
FOR EACH ROW
BEGIN
    UPDATE organization
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.id;
END;
CREATE TABLE casbin_rule(
    p_type VARCHAR(32)  DEFAULT '' NOT NULL,
    v0     VARCHAR(255) DEFAULT '' NOT NULL,
    v1     VARCHAR(255) DEFAULT '' NOT NULL,
    v2     VARCHAR(255) DEFAULT '' NOT NULL,
    v3     VARCHAR(255) DEFAULT '' NOT NULL,
    v4     VARCHAR(255) DEFAULT '' NOT NULL,
    v5     VARCHAR(255) DEFAULT '' NOT NULL,
    CHECK (TYPEOF("p_type") = "text" AND
           LENGTH("p_type") <= 32),
    CHECK (TYPEOF("v0") = "text" AND
           LENGTH("v0") <= 255),
    CHECK (TYPEOF("v1") = "text" AND
           LENGTH("v1") <= 255),
    CHECK (TYPEOF("v2") = "text" AND
           LENGTH("v2") <= 255),
    CHECK (TYPEOF("v3") = "text" AND
           LENGTH("v3") <= 255),
    CHECK (TYPEOF("v4") = "text" AND
           LENGTH("v4") <= 255),
    CHECK (TYPEOF("v5") = "text" AND
           LENGTH("v5") <= 255)
);
CREATE INDEX idx_casbin_rule ON casbin_rule (p_type,v0,v1);
CREATE TABLE game (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    organization_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, description TEXT DEFAULT '', region TEXT NOT NULL DEFAULT 'na', db_url TEXT DEFAULT '',
    FOREIGN KEY (organization_id) REFERENCES organization (id) ON DELETE CASCADE
);
CREATE TRIGGER update_game_updated_at
AFTER UPDATE ON game
FOR EACH ROW
BEGIN
    UPDATE game
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.id;
END;
CREATE UNIQUE INDEX unique_idp_user_id ON user (idp_user_id);
CREATE TABLE goqite (
  id text primary key default ('m_' || lower(hex(randomblob(16)))),
  created text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')),
  updated text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')),
  queue text not null,
  body blob not null,
  timeout text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')),
  received integer not null default 0
) strict;
CREATE TRIGGER goqite_updated_timestamp after update on goqite begin
  update goqite set updated = strftime('%Y-%m-%dT%H:%M:%fZ') where id = old.id;
end;
CREATE INDEX goqite_queue_created_idx on goqite (queue, created);
