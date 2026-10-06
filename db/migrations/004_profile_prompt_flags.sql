SET XACT_ABORT ON;
BEGIN TRANSACTION;

IF COL_LENGTH('dbo.UserProfiles', 'NamePrompted') IS NULL
BEGIN
    ALTER TABLE dbo.UserProfiles
    ADD NamePrompted BIT NOT NULL
        CONSTRAINT DF_UserProfiles_NamePrompted DEFAULT 0;
END;

IF COL_LENGTH('dbo.UserProfiles', 'MobilePrompted') IS NULL
BEGIN
    ALTER TABLE dbo.UserProfiles
    ADD MobilePrompted BIT NOT NULL
        CONSTRAINT DF_UserProfiles_MobilePrompted DEFAULT 0;
END;

COMMIT TRANSACTION;
