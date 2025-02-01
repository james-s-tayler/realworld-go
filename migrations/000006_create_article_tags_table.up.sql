PRAGMA foreign_keys = ON;

CREATE TABLE ArticleTag (
    ArticleId INTEGER NOT NULL,
    TagId INTEGER NOT NULL,
    PRIMARY KEY (ArticleId, TagId),
    FOREIGN KEY (ArticleId) REFERENCES Article (ArticleId) ON DELETE CASCADE,
    FOREIGN KEY (TagId) REFERENCES Tag (TagId) ON DELETE CASCADE
);