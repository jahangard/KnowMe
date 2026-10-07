SET XACT_ABORT ON;
BEGIN TRANSACTION;

IF OBJECT_ID('dbo.LifeLessons', 'U') IS NULL
BEGIN
    CREATE TABLE dbo.LifeLessons (
        Id BIGINT IDENTITY(1,1) NOT NULL CONSTRAINT PK_LifeLessons PRIMARY KEY,
        Code NVARCHAR(100) NOT NULL CONSTRAINT UQ_LifeLessons_Code UNIQUE,
        Title NVARCHAR(250) NOT NULL,
        Summary NVARCHAR(500) NOT NULL,
        Goal NVARCHAR(MAX) NOT NULL CONSTRAINT DF_LifeLessons_Goal DEFAULT N'',
        Explanation NVARCHAR(MAX) NOT NULL CONSTRAINT DF_LifeLessons_Explanation DEFAULT N'',
        Consequences NVARCHAR(MAX) NOT NULL CONSTRAINT DF_LifeLessons_Consequences DEFAULT N'',
        AgeRange NVARCHAR(250) NOT NULL CONSTRAINT DF_LifeLessons_AgeRange DEFAULT N'',
        Method NVARCHAR(MAX) NOT NULL CONSTRAINT DF_LifeLessons_Method DEFAULT N'',
        SeekHelp NVARCHAR(MAX) NOT NULL CONSTRAINT DF_LifeLessons_SeekHelp DEFAULT N'',
        KeyPoint NVARCHAR(500) NOT NULL CONSTRAINT DF_LifeLessons_KeyPoint DEFAULT N'',
        Content NVARCHAR(MAX) NOT NULL,
        SortOrder INT NOT NULL CONSTRAINT DF_LifeLessons_Sort DEFAULT 0,
        IsActive BIT NOT NULL CONSTRAINT DF_LifeLessons_Active DEFAULT 1
    );
END;

IF COL_LENGTH('dbo.LifeLessons', 'Goal') IS NULL ALTER TABLE dbo.LifeLessons ADD Goal NVARCHAR(MAX) NOT NULL CONSTRAINT DF_LifeLessons_Goal DEFAULT N'';
IF COL_LENGTH('dbo.LifeLessons', 'Explanation') IS NULL ALTER TABLE dbo.LifeLessons ADD Explanation NVARCHAR(MAX) NOT NULL CONSTRAINT DF_LifeLessons_Explanation DEFAULT N'';
IF COL_LENGTH('dbo.LifeLessons', 'Consequences') IS NULL ALTER TABLE dbo.LifeLessons ADD Consequences NVARCHAR(MAX) NOT NULL CONSTRAINT DF_LifeLessons_Consequences DEFAULT N'';
IF COL_LENGTH('dbo.LifeLessons', 'AgeRange') IS NULL ALTER TABLE dbo.LifeLessons ADD AgeRange NVARCHAR(250) NOT NULL CONSTRAINT DF_LifeLessons_AgeRange DEFAULT N'';
IF COL_LENGTH('dbo.LifeLessons', 'Method') IS NULL ALTER TABLE dbo.LifeLessons ADD Method NVARCHAR(MAX) NOT NULL CONSTRAINT DF_LifeLessons_Method DEFAULT N'';
IF COL_LENGTH('dbo.LifeLessons', 'SeekHelp') IS NULL ALTER TABLE dbo.LifeLessons ADD SeekHelp NVARCHAR(MAX) NOT NULL CONSTRAINT DF_LifeLessons_SeekHelp DEFAULT N'';
IF COL_LENGTH('dbo.LifeLessons', 'KeyPoint') IS NULL ALTER TABLE dbo.LifeLessons ADD KeyPoint NVARCHAR(500) NOT NULL CONSTRAINT DF_LifeLessons_KeyPoint DEFAULT N'';

COMMIT TRANSACTION;
