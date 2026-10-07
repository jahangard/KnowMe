# KnowMe Architecture

> وضعیت معماری تا 2026-10-07. این سند مرجع معماری اجرایی پروژه است.

## 1. هدف معماری

KnowMe باید یک موتور تست سبک، قابل‌گسترش و Data-driven باشد که در MVP از Telegram استفاده می‌کند و بدون وابستگی به UI یا دیتابیس خاص بتواند توسعه پیدا کند.

اصل اصلی معماری:

**محتوای تست باید داده باشد، نه کد.**

برای تست‌های استاندارد Single Choice، افزودن یا ویرایش تست نباید نیازمند تغییر Telegram Handler، Test Engine یا Result UI باشد.

## 2. Runtime

- Language: Go
- Telegram transport: Long Polling
- ORM: GORM
- Development database: SQLite
- Production database: SQL Server
- Public inbound port: لازم نیست
- Configuration: environment variables + optional local `.env`
- Logging: Go `slog` + file rotation
- Target OS: Windows 10/11, Windows Server, Linux و در آینده ARM/container

Long Polling باعث می‌شود Bot برای دریافت Update به IP ثابت، Domain یا Port Forward نیاز نداشته باشد.

## 3. لایه‌ها

### Telegram Transport

مسئول:

- دریافت Update
- Callback Query
- نمایش منوها
- ویرایش Message برای جلوگیری از شلوغی Chat
- تبدیل Interaction تلگرام به فراخوانی Application Service

Telegram نباید منطق امتیازدهی یا محتوای تست را در خود نگه دارد.

### Test Catalog

مسئول:

- خواندن Categoryهای فعال
- خواندن Testهای فعال هر Category
- نمایش Testها بدون Hard-code

منبع داده: `TestCategories` و `Tests`.

### Test Engine

مسئول:

- Start/Resume Session
- تعیین سؤال بعدی
- ثبت فوری هر Answer
- پایان Session
- محاسبه Score
- ساخت Result
- Update کردن UserTrait

Test Engine به عنوان یا متن یک تست خاص وابسته نیست.

### Scoring

مدل فعلی:

`single_choice + weighted trait + dominant trait`

هر Option دارای:

- `TraitKey`
- `Score`

در پایان، Score هر Trait جمع می‌شود و Trait غالب نتیجه اصلی است.

Tie-break:

1. `TestResultProfiles.SortOrder`
2. سپس `TraitKey` برای Deterministic بودن

### Result Presentation

نتیجه از جدول `TestResultProfiles` خوانده می‌شود:

- Label
- Title
- Subtitle
- Description
- SortOrder

بنابراین Result متناظر با Trait در Telegram Hard-code نیست.

### Roadmap Engine

نسخه فعلی:

اولین Test فعال که User هنوز Completed نکرده است.

مرتب‌سازی فعلی بر اساس:

1. `Tests.SortOrder`
2. `Tests.Id`

نسخه آینده باید Candidateها را بر اساس اطلاعات بیشتری امتیازدهی کند:

- Testهای Completed
- UserTraits
- Confidence
- Answer history
- Abandonment
- Age/Gender applicability
- Category diversity
- Exploration vs personalization

Roadmap باید تا حد ممکن Explainable و Deterministic باقی بماند.

### Progressive Profile

اصل محصول این است که کاربر در شروع با فرم بلند مواجه نشود.

جریان فعلی:

1. Gender
2. Age
3. ورود به Test experience
4. بعد از اولین Test، Name/Nickname به‌صورت اختیاری
5. Mobile به‌صورت اختیاری و در مرحله بعدی

Profile باید به مرور کامل شود.

### Event Tracking

`UserEvents` برای ثبت Interactionهای قابل تحلیل در نظر گرفته شده است.

Eventهای هدف:

- test_started
- question_answered
- test_completed
- test_abandoned
- roadmap_recommended
- profile_field_saved
- profile_field_skipped

ثبت Event باید مستقل از UI و مناسب تحلیل Product باشد.

### Share Result

Share Card Generator هنوز Future module است.

هدف:

- Telegram share
- Instagram Story
- Title
- Result type
- Bars/Stars
- Short insight
- Branding
- CTA/deep link

## 4. Data Model

جداول اصلی:

### Identity / Profile

- `Users`
- `UserProfiles`

### Test Definition

- `TestCategories`
- `Tests`
- `Questions`
- `QuestionOptions`
- `TestResultProfiles`

### Test Runtime

- `TestSessions`
- `TestAnswers`
- `TestResults`

### Derived Knowledge / Analytics

- `UserTraits`
- `UserEvents`

## 5. Data-driven Test Contract

برای یک Test استاندارد، کد Runtime فقط Contract عمومی را می‌شناسد.

تعریف تست در DB شامل:

```text
TestCategory
  └── Test
       ├── Questions
       │    └── QuestionOptions
       │          ├── TraitKey
       │          └── Score
       └── TestResultProfiles
             ├── TraitKey
             ├── Label
             ├── Title
             ├── Subtitle
             └── Description
```

در نتیجه برای افزودن Test جدید در Category موجود:

**فقط DB Update کافی است.**

## 6. Database Providers

### SQLite

برای Development پیش‌فرض است.

ویژگی‌ها:

- local file
- مسیر پیش‌فرض: `data/knowme.db`
- WAL
- Foreign Keys
- Busy Timeout
- Max Open Connection محدود برای رفتار پایدار

### SQL Server

برای Production و محیط Remote پشتیبانی می‌شود.

Connection از `SQLSERVER_DSN` خوانده می‌شود.

TLS باید در محیط Production امن نگه داشته شود. `TrustServerCertificate=true` فقط زمانی استفاده شود که Certificate chain معتبر در دسترس نیست و ریسک آن پذیرفته شده باشد.

## 7. Migration و Seed

Schema با GORM `AutoMigrate` مدیریت می‌شود.

Seed فعلی نقش **Bootstrap** دارد:

- Categoryهای اولیه را فقط اگر وجود ندارند می‌سازد.
- پنج Test اولیه را فقط اگر وجود ندارند می‌سازد.
- Question/Option/ResultProfile موجود را در Restart بازنویسی نمی‌کند.

اصل مهم:

**بعد از ایجاد اولیه، Database منبع حقیقت Content است.**

بنابراین تغییر مستقیم Title، Question، Option، Score یا Result Profile در DB با Restart از بین نمی‌رود.

## 8. پنج Test اولیه

1. `love_style_v1` — 💞 سبک عشق‌ورزی
2. `attraction_style_v1` — ✨ تیپ جذابیت
3. `dating_scenarios_v1` — 💘 سناریوهای قرار
4. `personal_boundaries_v1` — 🛡 مرزهای شخصی
5. `bitter_truth_v1` — 🪞 حقیقت تلخ

هر کدام در نسخه اولیه ۸ سؤال دارند.

## 9. Category Browser

UI فعلی:

```text
Home
  ├── 🧩 موضوعات تست‌ها
  │    ├── 🧭 نقشه راه من
  │    ├── ❤️ عشق و رابطه
  │    ├── 🎯 چالشی و باحال
  │    ├── 🧠 خودشناسی عمیق‌تر
  │    └── 🔞 ۱۸+
  ├── 🧭 نقشه راه من
  ├── 👤 پروفایل من
  └── ✨ درباره KnowMe
```

Categoryها از DB خوانده می‌شوند.

«نقشه راه من» یک System/Virtual entry است و Category عادی DB نیست.

بخش ۱۸+ Age Gate دارد و برای User زیر ۱۸ سال باز نمی‌شود.

## 10. UX Navigation Principle

برای Interactionهای Inline، ترجیح با Edit همان Message است تا Chat شلوغ نشود.

قاعده:

- Question → edit current message
- Category → edit current message
- Test list → edit current message
- Back/Home → تا جای ممکن edit
- Promptهای مبتنی بر Text/Contact در صورت نیاز Message جداگانه

## 11. Logging

Logging باید هم برای Development و هم Post-mortem مناسب باشد.

فعلی:

- Structured logging با `slog`
- Terminal output
- File output
- Rotation با Lumberjack
- Configurable level
- Startup/Shutdown
- DB connection/migration
- Telegram update handling
- Error duration
- panic stack

فایل پیش‌فرض:

`logs/knowme.log`

Secrets نباید Log شوند.

## 12. Security / Privacy

- `.env` Commit نمی‌شود.
- Bot token و DB password نباید Log شوند.
- Mobile اختیاری است.
- شماره موبایل در Profile UI به‌صورت مقدار واقعی نمایش داده نمی‌شود.
- Testهای ۱۸+ باید Age Gate داشته باشند.
- Testهای سرگرمی نباید به‌عنوان Diagnosis پزشکی یا روان‌شناختی معرفی شوند.
- در Testهای ۱۸+ از Labelهایی مثل «طبیعی/غیرطبیعی» اجتناب شود.
- در آینده Delete My Data و Privacy/Consent flow باید اضافه شوند.

## 13. CI

GitHub Actions:

- Setup Go از `go.mod`
- `go mod tidy`
- `gofmt` check
- `go test ./...`

هر تغییر معماری باید حداقل Compile و Format check را پاس کند.

## 14. محدودیت فعلی Engine

بدون تغییر کد، Engine فعلی از این خانواده تست پشتیبانی می‌کند:

- Single Choice
- یک TraitKey برای هر Option
- Score عددی
- Dominant Trait result

برای قابلیت‌های زیر نیاز به توسعه Engine داریم:

- Multi Select
- Slider
- Matrix
- Branching Questions
- Conditional Question
- Multi-dimensional formula
- threshold-based result
- composite result
- adaptive question selection

این توسعه باید به‌صورت Rule/Scoring abstraction انجام شود، نه Hard-code برای Test خاص.
