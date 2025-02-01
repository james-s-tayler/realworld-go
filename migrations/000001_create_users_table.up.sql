PRAGMA foreign_keys = ON;

CREATE TABLE User (
    UserId INTEGER PRIMARY KEY AUTOINCREMENT,
    Email TEXT NOT NULL UNIQUE,
    PasswordHash TEXT NOT NULL,
    Username TEXT NOT NULL UNIQUE,
    Bio TEXT,
    Image TEXT
);