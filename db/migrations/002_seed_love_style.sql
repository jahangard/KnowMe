SET XACT_ABORT ON;
BEGIN TRANSACTION;

DECLARE @CategoryId BIGINT;
SELECT @CategoryId = Id FROM dbo.TestCategories WHERE Code = N'love';

IF @CategoryId IS NULL
BEGIN
    INSERT INTO dbo.TestCategories (Code, Title, SortOrder, IsActive)
    VALUES (N'love', N'عشق و رابطه', 10, 1);
    SET @CategoryId = SCOPE_IDENTITY();
END;

DECLARE @TestId BIGINT;
SELECT @TestId = Id FROM dbo.Tests WHERE Code = N'love_style_v1';

IF @TestId IS NULL
BEGIN
    INSERT INTO dbo.Tests (CategoryId, Code, Title, Description, SortOrder, IsActive)
    VALUES (
        @CategoryId,
        N'love_style_v1',
        N'💞 سبک عشق‌ورزی',
        N'یک تست کوتاه سرگرمی برای شناخت شیوه غالب ابراز علاقه در رابطه.',
        10,
        1
    );
    SET @TestId = SCOPE_IDENTITY();
END;

IF NOT EXISTS (SELECT 1 FROM dbo.Questions WHERE TestId = @TestId)
BEGIN
    INSERT INTO dbo.Questions (TestId, [Text], [Order], QuestionType, IsActive)
    VALUES
    (@TestId, N'وقتی کسی که دوستش داری ناراحته، معمولاً اولین واکنش تو چیه؟', 1, N'single_choice', 1),
    (@TestId, N'اگر بخوای خیلی واضح نشون بدی که کسی برات مهمه، بیشتر چه کار می‌کنی؟', 2, N'single_choice', 1),
    (@TestId, N'در یک روز شلوغ، کدوم رفتار طرف مقابل بیشتر به دلت می‌شینه؟', 3, N'single_choice', 1),
    (@TestId, N'وقتی دلخور می‌شی، کدوم کار بیشتر کمک می‌کنه دوباره احساس نزدیکی کنی؟', 4, N'single_choice', 1),
    (@TestId, N'برای یک مناسبت خاص، کدوم برنامه بیشتر تو رو خوشحال می‌کنه؟', 5, N'single_choice', 1),
    (@TestId, N'وقتی از کسی خوشت میاد، خودت ناخودآگاه بیشتر کدوم رفتار رو انجام می‌دی؟', 6, N'single_choice', 1),
    (@TestId, N'در یک رابطه، کدوم کمبود بیشتر اذیتت می‌کنه؟', 7, N'single_choice', 1),
    (@TestId, N'اگر فقط یکی رو انتخاب کنی، کدوم جمله بیشتر حس دوست‌داشتن بهت می‌ده؟', 8, N'single_choice', 1);

    DECLARE @Q1 BIGINT = (SELECT Id FROM dbo.Questions WHERE TestId=@TestId AND [Order]=1);
    DECLARE @Q2 BIGINT = (SELECT Id FROM dbo.Questions WHERE TestId=@TestId AND [Order]=2);
    DECLARE @Q3 BIGINT = (SELECT Id FROM dbo.Questions WHERE TestId=@TestId AND [Order]=3);
    DECLARE @Q4 BIGINT = (SELECT Id FROM dbo.Questions WHERE TestId=@TestId AND [Order]=4);
    DECLARE @Q5 BIGINT = (SELECT Id FROM dbo.Questions WHERE TestId=@TestId AND [Order]=5);
    DECLARE @Q6 BIGINT = (SELECT Id FROM dbo.Questions WHERE TestId=@TestId AND [Order]=6);
    DECLARE @Q7 BIGINT = (SELECT Id FROM dbo.Questions WHERE TestId=@TestId AND [Order]=7);
    DECLARE @Q8 BIGINT = (SELECT Id FROM dbo.Questions WHERE TestId=@TestId AND [Order]=8);

    INSERT INTO dbo.QuestionOptions (QuestionId, [Text], [Order], Score, TraitKey) VALUES
    (@Q1, N'با حرف زدن آرومش می‌کنم', 1, 1, N'words'),
    (@Q1, N'کنارش می‌مونم و وقتم رو بهش می‌دم', 2, 1, N'time'),
    (@Q1, N'یه کاری براش انجام می‌دم که حالش بهتر شه', 3, 1, N'care'),
    (@Q1, N'با آغوش و نزدیکی بهش آرامش می‌دم', 4, 1, N'touch'),

    (@Q2, N'بهش می‌گم چقدر برام مهمه', 1, 1, N'words'),
    (@Q2, N'یه زمان مخصوص فقط برای دوتامون می‌ذارم', 2, 1, N'time'),
    (@Q2, N'کارش رو سبک می‌کنم یا کمکش می‌کنم', 3, 1, N'care'),
    (@Q2, N'بیشتر بغلش می‌کنم و نزدیکش می‌مونم', 4, 1, N'touch'),

    (@Q3, N'یک پیام محبت‌آمیز و صمیمی', 1, 1, N'words'),
    (@Q3, N'اینکه با وجود شلوغی، وقت برای من باز کنه', 2, 1, N'time'),
    (@Q3, N'اینکه بدون گفتن، یک کارم رو انجام بده', 3, 1, N'care'),
    (@Q3, N'یک بغل گرم وقتی همدیگه رو می‌بینیم', 4, 1, N'touch'),

    (@Q4, N'اینکه حرف دلش رو واضح بگه', 1, 1, N'words'),
    (@Q4, N'اینکه بشینیم و باهم وقت بگذرونیم', 2, 1, N'time'),
    (@Q4, N'اینکه برای جبران، یک کار واقعی انجام بده', 3, 1, N'care'),
    (@Q4, N'اینکه با یک آغوش صمیمی فاصله رو کم کنه', 4, 1, N'touch'),

    (@Q5, N'یک نامه یا پیام خاص و احساسی', 1, 1, N'words'),
    (@Q5, N'یک روز کامل فقط با هم بودن', 2, 1, N'time'),
    (@Q5, N'یک کار غافلگیرکننده که زندگی‌م رو راحت‌تر کنه', 3, 1, N'care'),
    (@Q5, N'یک شب صمیمی و پر از نزدیکی', 4, 1, N'touch'),

    (@Q6, N'زیاد تعریف و ابراز احساس می‌کنم', 1, 1, N'words'),
    (@Q6, N'دنبال فرصت می‌گردم باهاش تنها باشم', 2, 1, N'time'),
    (@Q6, N'کارهای کوچیکش رو انجام می‌دم', 3, 1, N'care'),
    (@Q6, N'با تماس و نزدیکی علاقه‌م رو نشون می‌دم', 4, 1, N'touch'),

    (@Q7, N'کمبود حرف‌های محبت‌آمیز', 1, 1, N'words'),
    (@Q7, N'کمبود وقت دونفره', 2, 1, N'time'),
    (@Q7, N'اینکه همه چیز فقط در حد حرف بمونه', 3, 1, N'care'),
    (@Q7, N'فاصله و سردی فیزیکی', 4, 1, N'touch'),

    (@Q8, N'«دوستت دارم و بهت افتخار می‌کنم»', 1, 1, N'words'),
    (@Q8, N'«امروز رو کامل برای تو خالی کردم»', 2, 1, N'time'),
    (@Q8, N'«نگران نباش، من انجامش دادم»', 3, 1, N'care'),
    (@Q8, N'«بیا یه بغل طولانی»', 4, 1, N'touch');
END;

COMMIT TRANSACTION;
