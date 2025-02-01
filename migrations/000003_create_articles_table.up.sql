PRAGMA foreign_keys = ON;

CREATE TABLE Article (
    ArticleId INTEGER NOT NULL PRIMARY KEY,
    UserId INTEGER NOT NULL,
    Title TEXT NOT NULL UNIQUE,
    Slug TEXT NOT NULL UNIQUE,
    Description TEXT NOT NULL,
    Body TEXT NOT NULL,
    CreatedAt TEXT NOT NULL,
    UpdatedAt TEXT NOT NULL,
    FOREIGN KEY (UserId) REFERENCES User(UserId) ON DELETE CASCADE
);

CREATE INDEX idx_articles_user_id ON Article (UserId);