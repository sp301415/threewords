CREATE TABLE fileDB (
  id TEXT PRIMARY KEY,
  file_path TEXT NOT NULL,
  original_name TEXT NOT NULL,
  expire_date INTEGER NOT NULL
);
