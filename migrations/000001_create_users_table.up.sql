CREATE TABLE User (
    Id INTEGER PRIMARY KEY AUTOINCREMENT,
    Email TEXT NOT NULL UNIQUE,
    PasswordHash TEXT NOT NULL,
    Username TEXT NOT NULL UNIQUE,
    Bio TEXT,
    Image TEXT
);