SET XACT_ABORT ON;
BEGIN TRANSACTION;

CREATE TABLE dbo.Users (
    Id BIGINT IDENTITY(1,1) NOT NULL CONSTRAINT PK_Users PRIMARY KEY,
    TelegramUserId BIGINT NOT NULL,
    Username NVARCHAR(255) NULL,
    FirstSeenAt DATETIME2 NOT NULL CONSTRAINT DF_Users_FirstSeenAt DEFAULT SYSUTCDATETIME(),
    LastSeenAt DATETIME2 NOT NULL CONSTRAINT DF_Users_LastSeenAt DEFAULT SYSUTCDATETIME(),
    CONSTRAINT UQ_Users_TelegramUserId UNIQUE (TelegramUserId)
);

CREATE TABLE dbo.UserProfiles (
    UserId BIGINT NOT NULL CONSTRAINT PK_UserProfiles PRIMARY KEY,
    Name NVARCHAR(200) NULL,
    Age INT NULL,
    Gender NVARCHAR(50) NULL,
    Mobile NVARCHAR(50) NULL,
    ProfileCompletionLevel INT NOT NULL CONSTRAINT DF_UserProfiles_Level DEFAULT 0,
    UpdatedAt DATETIME2 NOT NULL CONSTRAINT DF_UserProfiles_UpdatedAt DEFAULT SYSUTCDATETIME(),
    CONSTRAINT FK_UserProfiles_Users FOREIGN KEY (UserId) REFERENCES dbo.Users(Id)
);

CREATE TABLE dbo.TestCategories (
    Id BIGINT IDENTITY(1,1) NOT NULL CONSTRAINT PK_TestCategories PRIMARY KEY,
    Code NVARCHAR(100) NOT NULL,
    Title NVARCHAR(200) NOT NULL,
    SortOrder INT NOT NULL CONSTRAINT DF_TestCategories_Sort DEFAULT 0,
    IsActive BIT NOT NULL CONSTRAINT DF_TestCategories_Active DEFAULT 1,
    CONSTRAINT UQ_TestCategories_Code UNIQUE (Code)
);

CREATE TABLE dbo.Tests (
    Id BIGINT IDENTITY(1,1) NOT NULL CONSTRAINT PK_Tests PRIMARY KEY,
    CategoryId BIGINT NOT NULL,
    Code NVARCHAR(100) NOT NULL,
    Title NVARCHAR(250) NOT NULL,
    Description NVARCHAR(1000) NULL,
    SortOrder INT NOT NULL CONSTRAINT DF_Tests_Sort DEFAULT 0,
    IsActive BIT NOT NULL CONSTRAINT DF_Tests_Active DEFAULT 1,
    CreatedAt DATETIME2 NOT NULL CONSTRAINT DF_Tests_CreatedAt DEFAULT SYSUTCDATETIME(),
    CONSTRAINT UQ_Tests_Code UNIQUE (Code),
    CONSTRAINT FK_Tests_Categories FOREIGN KEY (CategoryId) REFERENCES dbo.TestCategories(Id)
);

CREATE TABLE dbo.Questions (
    Id BIGINT IDENTITY(1,1) NOT NULL CONSTRAINT PK_Questions PRIMARY KEY,
    TestId BIGINT NOT NULL,
    [Text] NVARCHAR(1000) NOT NULL,
    [Order] INT NOT NULL,
    QuestionType NVARCHAR(50) NOT NULL,
    IsActive BIT NOT NULL CONSTRAINT DF_Questions_Active DEFAULT 1,
    CONSTRAINT FK_Questions_Tests FOREIGN KEY (TestId) REFERENCES dbo.Tests(Id),
    CONSTRAINT UQ_Questions_Test_Order UNIQUE (TestId, [Order])
);

CREATE TABLE dbo.QuestionOptions (
    Id BIGINT IDENTITY(1,1) NOT NULL CONSTRAINT PK_QuestionOptions PRIMARY KEY,
    QuestionId BIGINT NOT NULL,
    [Text] NVARCHAR(500) NOT NULL,
    [Order] INT NOT NULL,
    Score DECIMAL(9,4) NULL,
    TraitKey NVARCHAR(100) NULL,
    CONSTRAINT FK_QuestionOptions_Questions FOREIGN KEY (QuestionId) REFERENCES dbo.Questions(Id)
);

CREATE TABLE dbo.TestSessions (
    Id BIGINT IDENTITY(1,1) NOT NULL CONSTRAINT PK_TestSessions PRIMARY KEY,
    UserId BIGINT NOT NULL,
    TestId BIGINT NOT NULL,
    CurrentQuestionId BIGINT NULL,
    [Status] NVARCHAR(30) NOT NULL,
    StartedAt DATETIME2 NOT NULL CONSTRAINT DF_TestSessions_StartedAt DEFAULT SYSUTCDATETIME(),
    CompletedAt DATETIME2 NULL,
    CONSTRAINT FK_TestSessions_Users FOREIGN KEY (UserId) REFERENCES dbo.Users(Id),
    CONSTRAINT FK_TestSessions_Tests FOREIGN KEY (TestId) REFERENCES dbo.Tests(Id),
    CONSTRAINT FK_TestSessions_CurrentQuestion FOREIGN KEY (CurrentQuestionId) REFERENCES dbo.Questions(Id)
);

CREATE TABLE dbo.TestAnswers (
    Id BIGINT IDENTITY(1,1) NOT NULL CONSTRAINT PK_TestAnswers PRIMARY KEY,
    SessionId BIGINT NOT NULL,
    QuestionId BIGINT NOT NULL,
    OptionId BIGINT NULL,
    AnswerValue NVARCHAR(1000) NULL,
    ResponseTimeMs BIGINT NULL,
    AnsweredAt DATETIME2 NOT NULL CONSTRAINT DF_TestAnswers_AnsweredAt DEFAULT SYSUTCDATETIME(),
    CONSTRAINT FK_TestAnswers_Sessions FOREIGN KEY (SessionId) REFERENCES dbo.TestSessions(Id),
    CONSTRAINT FK_TestAnswers_Questions FOREIGN KEY (QuestionId) REFERENCES dbo.Questions(Id),
    CONSTRAINT FK_TestAnswers_Options FOREIGN KEY (OptionId) REFERENCES dbo.QuestionOptions(Id)
);

CREATE TABLE dbo.TestResults (
    Id BIGINT IDENTITY(1,1) NOT NULL CONSTRAINT PK_TestResults PRIMARY KEY,
    SessionId BIGINT NOT NULL,
    ResultType NVARCHAR(100) NOT NULL,
    ScoreJson NVARCHAR(MAX) NULL,
    Summary NVARCHAR(MAX) NULL,
    ShareImageUrl NVARCHAR(1000) NULL,
    CreatedAt DATETIME2 NOT NULL CONSTRAINT DF_TestResults_CreatedAt DEFAULT SYSUTCDATETIME(),
    CONSTRAINT UQ_TestResults_Session UNIQUE (SessionId),
    CONSTRAINT FK_TestResults_Sessions FOREIGN KEY (SessionId) REFERENCES dbo.TestSessions(Id),
    CONSTRAINT CK_TestResults_ScoreJson CHECK (ScoreJson IS NULL OR ISJSON(ScoreJson) = 1)
);

CREATE TABLE dbo.UserTraits (
    UserId BIGINT NOT NULL,
    TraitKey NVARCHAR(100) NOT NULL,
    Score DECIMAL(9,4) NOT NULL,
    Confidence DECIMAL(9,4) NOT NULL,
    UpdatedAt DATETIME2 NOT NULL CONSTRAINT DF_UserTraits_UpdatedAt DEFAULT SYSUTCDATETIME(),
    CONSTRAINT PK_UserTraits PRIMARY KEY (UserId, TraitKey),
    CONSTRAINT FK_UserTraits_Users FOREIGN KEY (UserId) REFERENCES dbo.Users(Id)
);

CREATE TABLE dbo.UserEvents (
    Id BIGINT IDENTITY(1,1) NOT NULL CONSTRAINT PK_UserEvents PRIMARY KEY,
    UserId BIGINT NOT NULL,
    EventType NVARCHAR(100) NOT NULL,
    TestId BIGINT NULL,
    QuestionId BIGINT NULL,
    MetadataJson NVARCHAR(MAX) NULL,
    CreatedAt DATETIME2 NOT NULL CONSTRAINT DF_UserEvents_CreatedAt DEFAULT SYSUTCDATETIME(),
    CONSTRAINT FK_UserEvents_Users FOREIGN KEY (UserId) REFERENCES dbo.Users(Id),
    CONSTRAINT FK_UserEvents_Tests FOREIGN KEY (TestId) REFERENCES dbo.Tests(Id),
    CONSTRAINT FK_UserEvents_Questions FOREIGN KEY (QuestionId) REFERENCES dbo.Questions(Id),
    CONSTRAINT CK_UserEvents_MetadataJson CHECK (MetadataJson IS NULL OR ISJSON(MetadataJson) = 1)
);

CREATE INDEX IX_TestSessions_User_Status ON dbo.TestSessions(UserId, [Status]);
CREATE INDEX IX_TestAnswers_Session ON dbo.TestAnswers(SessionId, QuestionId);
CREATE INDEX IX_UserEvents_User_CreatedAt ON dbo.UserEvents(UserId, CreatedAt DESC);

COMMIT TRANSACTION;
