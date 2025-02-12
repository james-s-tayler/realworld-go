PRAGMA foreign_keys = ON;

CREATE TABLE Follower (
    UserId INTEGER NOT NULL,
    FollowUserId INTEGER NOT NULL,
    PRIMARY KEY (UserId, FollowUserId),
    FOREIGN KEY (UserId) REFERENCES User(UserId) ON DELETE CASCADE,
    FOREIGN KEY (FollowUserId) REFERENCES User(UserId) ON DELETE CASCADE
);