PRAGMA foreign_keys = ON;

CREATE TABLE ArticleFavorite (
    ArticleId INTEGER NOT NULL,
    UserId INTEGER NOT NULL,
    PRIMARY KEY (ArticleId, UserId),
    FOREIGN KEY (ArticleId) REFERENCES Article (ArticleId) ON DELETE CASCADE,
    FOREIGN KEY (UserId) REFERENCES User (UserId) ON DELETE CASCADE
);