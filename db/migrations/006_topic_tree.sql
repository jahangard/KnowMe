SET XACT_ABORT ON;
BEGIN TRANSACTION;

IF COL_LENGTH('dbo.TestCategories', 'ParentId') IS NULL
    ALTER TABLE dbo.TestCategories ADD ParentId BIGINT NULL;
IF COL_LENGTH('dbo.TestCategories', 'CatalogType') IS NULL
    ALTER TABLE dbo.TestCategories ADD CatalogType NVARCHAR(30) NOT NULL CONSTRAINT DF_TestCategories_CatalogType DEFAULT N'tests';
IF COL_LENGTH('dbo.LifeLessons', 'TopicId') IS NULL
    ALTER TABLE dbo.LifeLessons ADD TopicId BIGINT NOT NULL CONSTRAINT DF_LifeLessons_TopicId DEFAULT 0;

COMMIT TRANSACTION;
