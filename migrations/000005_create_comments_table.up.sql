PRAGMA foreign_keys = ON;

CREATE TABLE Comment (
    CommentId INTEGER NOT NULL PRIMARY KEY,
    UserId INTEGER NOT NULL,
    ArticleId INTEGER NOT NULL,
    Body TEXT NOT NULL,
    CreatedAt TEXT NOT NULL,
    UpdatedAt TEXT NOT NULL,
    FOREIGN KEY (UserId) REFERENCES User (UserId) ON DELETE CASCADE
    FOREIGN KEY (ArticleId) REFERENCES Article (ArticleId) ON DELETE CASCADE
);

CREATE INDEX idx_comments_user_id ON Comment (UserId);