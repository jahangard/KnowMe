INSERT OR IGNORE INTO TestCategories (Code, Title, SortOrder, IsActive)
VALUES ('love', 'عشق و رابطه', 10, 1);

INSERT OR IGNORE INTO Tests (CategoryId, Code, Title, Description, SortOrder, IsActive)
SELECT Id, 'love_style_v1', '💞 سبک عشق‌ورزی',
       'یک تست کوتاه سرگرمی برای شناخت شیوه غالب ابراز علاقه در رابطه.',
       10, 1
FROM TestCategories
WHERE Code = 'love';

INSERT OR IGNORE INTO Questions (TestId, Text, "Order", QuestionType, IsActive)
SELECT Id, 'وقتی کسی که دوستش داری ناراحته، معمولاً اولین واکنش تو چیه؟', 1, 'single_choice', 1 FROM Tests WHERE Code='love_style_v1';
INSERT OR IGNORE INTO Questions (TestId, Text, "Order", QuestionType, IsActive)
SELECT Id, 'اگر بخوای خیلی واضح نشون بدی که کسی برات مهمه، بیشتر چه کار می‌کنی؟', 2, 'single_choice', 1 FROM Tests WHERE Code='love_style_v1';
INSERT OR IGNORE INTO Questions (TestId, Text, "Order", QuestionType, IsActive)
SELECT Id, 'در یک روز شلوغ، کدوم رفتار طرف مقابل بیشتر به دلت می‌شینه؟', 3, 'single_choice', 1 FROM Tests WHERE Code='love_style_v1';
INSERT OR IGNORE INTO Questions (TestId, Text, "Order", QuestionType, IsActive)
SELECT Id, 'وقتی دلخور می‌شی، کدوم کار بیشتر کمک می‌کنه دوباره احساس نزدیکی کنی؟', 4, 'single_choice', 1 FROM Tests WHERE Code='love_style_v1';
INSERT OR IGNORE INTO Questions (TestId, Text, "Order", QuestionType, IsActive)
SELECT Id, 'برای یک مناسبت خاص، کدوم برنامه بیشتر تو رو خوشحال می‌کنه؟', 5, 'single_choice', 1 FROM Tests WHERE Code='love_style_v1';
INSERT OR IGNORE INTO Questions (TestId, Text, "Order", QuestionType, IsActive)
SELECT Id, 'وقتی از کسی خوشت میاد، خودت ناخودآگاه بیشتر کدوم رفتار رو انجام می‌دی؟', 6, 'single_choice', 1 FROM Tests WHERE Code='love_style_v1';
INSERT OR IGNORE INTO Questions (TestId, Text, "Order", QuestionType, IsActive)
SELECT Id, 'در یک رابطه، کدوم کمبود بیشتر اذیتت می‌کنه؟', 7, 'single_choice', 1 FROM Tests WHERE Code='love_style_v1';
INSERT OR IGNORE INTO Questions (TestId, Text, "Order", QuestionType, IsActive)
SELECT Id, 'اگر فقط یکی رو انتخاب کنی، کدوم جمله بیشتر حس دوست‌داشتن بهت می‌ده؟', 8, 'single_choice', 1 FROM Tests WHERE Code='love_style_v1';

INSERT OR IGNORE INTO QuestionOptions (QuestionId, Text, "Order", Score, TraitKey)
SELECT q.Id, v.Text, v.Ord, 1, v.TraitKey
FROM Questions q
JOIN Tests t ON t.Id=q.TestId
JOIN (
    SELECT 1 QOrd, 1 Ord, 'با حرف زدن آرومش می‌کنم' Text, 'words' TraitKey UNION ALL
    SELECT 1,2,'کنارش می‌مونم و وقتم رو بهش می‌دم','time' UNION ALL
    SELECT 1,3,'یه کاری براش انجام می‌دم که حالش بهتر شه','care' UNION ALL
    SELECT 1,4,'با آغوش و نزدیکی بهش آرامش می‌دم','touch' UNION ALL
    SELECT 2,1,'بهش می‌گم چقدر برام مهمه','words' UNION ALL
    SELECT 2,2,'یه زمان مخصوص فقط برای دوتامون می‌ذارم','time' UNION ALL
    SELECT 2,3,'کارش رو سبک می‌کنم یا کمکش می‌کنم','care' UNION ALL
    SELECT 2,4,'بیشتر بغلش می‌کنم و نزدیکش می‌مونم','touch' UNION ALL
    SELECT 3,1,'یک پیام محبت‌آمیز و صمیمی','words' UNION ALL
    SELECT 3,2,'اینکه با وجود شلوغی، وقت برای من باز کنه','time' UNION ALL
    SELECT 3,3,'اینکه بدون گفتن، یک کارم رو انجام بده','care' UNION ALL
    SELECT 3,4,'یک بغل گرم وقتی همدیگه رو می‌بینیم','touch' UNION ALL
    SELECT 4,1,'اینکه حرف دلش رو واضح بگه','words' UNION ALL
    SELECT 4,2,'اینکه بشینیم و باهم وقت بگذرونیم','time' UNION ALL
    SELECT 4,3,'اینکه برای جبران، یک کار واقعی انجام بده','care' UNION ALL
    SELECT 4,4,'اینکه با یک آغوش صمیمی فاصله رو کم کنه','touch' UNION ALL
    SELECT 5,1,'یک نامه یا پیام خاص و احساسی','words' UNION ALL
    SELECT 5,2,'یک روز کامل فقط با هم بودن','time' UNION ALL
    SELECT 5,3,'یک کار غافلگیرکننده که زندگی‌م رو راحت‌تر کنه','care' UNION ALL
    SELECT 5,4,'یک شب صمیمی و پر از نزدیکی','touch' UNION ALL
    SELECT 6,1,'زیاد تعریف و ابراز احساس می‌کنم','words' UNION ALL
    SELECT 6,2,'دنبال فرصت می‌گردم باهاش تنها باشم','time' UNION ALL
    SELECT 6,3,'کارهای کوچیکش رو انجام می‌دم','care' UNION ALL
    SELECT 6,4,'با تماس و نزدیکی علاقه‌م رو نشون می‌دم','touch' UNION ALL
    SELECT 7,1,'کمبود حرف‌های محبت‌آمیز','words' UNION ALL
    SELECT 7,2,'کمبود وقت دونفره','time' UNION ALL
    SELECT 7,3,'اینکه همه چیز فقط در حد حرف بمونه','care' UNION ALL
    SELECT 7,4,'فاصله و سردی فیزیکی','touch' UNION ALL
    SELECT 8,1,'«دوستت دارم و بهت افتخار می‌کنم»','words' UNION ALL
    SELECT 8,2,'«امروز رو کامل برای تو خالی کردم»','time' UNION ALL
    SELECT 8,3,'«نگران نباش، من انجامش دادم»','care' UNION ALL
    SELECT 8,4,'«بیا یه بغل طولانی»','touch'
) v ON v.QOrd=q."Order"
WHERE t.Code='love_style_v1';
