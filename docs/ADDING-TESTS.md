# افزودن تست جدید بدون تغییر کد

KnowMe از این مرحله به بعد برای تست‌های استاندارد Single Choice کاملاً Data-driven است. یعنی برای اضافه‌کردن تست جدید نیازی به تغییر Telegram Handler، Test Engine یا UI نیست؛ کافی است داده‌های مربوط به تست در دیتابیس ثبت شوند.

## جداول موردنیاز

برای یک تست جدید معمولاً فقط این جداول درگیر هستند:

1. `TestCategories`
2. `Tests`
3. `Questions`
4. `QuestionOptions`
5. `TestResultProfiles`

در زمان اجرا:

- `TestSessions` اجرای تست را ثبت می‌کند.
- `TestAnswers` پاسخ‌های کاربر را ذخیره می‌کند.
- `TestResults` نتیجه نهایی هر Session را ذخیره می‌کند.
- `UserTraits` امتیاز Traitهای کاربر را به‌روزرسانی می‌کند.

## مدل امتیازدهی فعلی

هر گزینه می‌تواند یک `TraitKey` و یک `Score` داشته باشد.

پس از پایان تست:

1. امتیاز گزینه‌های انتخاب‌شده بر اساس `TraitKey` جمع می‌شوند.
2. Trait با بالاترین امتیاز ResultType اصلی می‌شود.
3. عنوان و توضیح نتیجه از `TestResultProfiles` خوانده می‌شود.
4. هیچ Result Title یا Trait Presentation در کد Telegram هاردکد نیست.

مثال:

```text
QuestionOption
  TraitKey = attraction:magnetic
  Score    = 1

TestResultProfile
  TraitKey    = attraction:magnetic
  Label       = مغناطیسی
  Title       = 🧲 جذابیت مغناطیسی
  Subtitle    = ...
  Description = ...
```

## ساختار پیشنهادی TraitKey

برای جلوگیری از تداخل بین تست‌ها، Traitهای جدید را Namespace کنید:

```text
attraction:magnetic
dating:planner
boundaries:firm
truth:approval
```

Traitهای قدیمی تست `love_style_v1` برای حفظ سازگاری همان `words`, `time`, `care`, `touch` باقی مانده‌اند.

## مراحل افزودن یک تست

### 1. ایجاد Test

یک رکورد در `Tests` با موارد زیر ثبت شود:

- `CategoryId`
- `Code` یکتا
- `Title`
- `Description`
- `SortOrder`
- `IsActive = 1`

### 2. ایجاد Questionها

هر سؤال در `Questions`:

- `TestId`
- `Text`
- `Order`
- `QuestionType = single_choice`
- `IsActive = 1`

### 3. ایجاد Optionها

هر گزینه در `QuestionOptions`:

- `QuestionId`
- `Text`
- `Order`
- `Score`
- `TraitKey`

### 4. ایجاد Result Profileها

برای هر Trait ممکن، یک رکورد در `TestResultProfiles`:

- `TestId`
- `TraitKey`
- `Label`
- `Title`
- `Subtitle`
- `Description`
- `SortOrder`
- `IsActive = 1`

## چه چیزهایی خودکار هستند؟

بعد از ثبت داده‌ها:

- تست خودکار در Category مربوطه دیده می‌شود.
- Start Test بدون تغییر کد کار می‌کند.
- سؤال‌ها و گزینه‌ها از DB خوانده می‌شوند.
- امتیازها از DB خوانده می‌شوند.
- نتیجه و متن Result از DB خوانده می‌شود.
- Score Breakdown با Labelهای DB ساخته می‌شود.
- تست به‌صورت خودکار وارد Roadmap می‌شود، مگر بعداً Rule دیگری برای Recommendation اضافه شود.

## محدودیت فعلی

این مدل برای تست‌های `single_choice + dominant trait` طراحی شده است.

اگر در آینده بخواهیم تست‌هایی با منطق کاملاً متفاوت مثل:

- چندانتخابی
- Slider
- فرمول وزنی چندمرحله‌ای
- شرط‌گذاری بر اساس چند Trait
- Branching Questions

اضافه کنیم، باید موتور Rule/Scoring توسعه پیدا کند. اما برای تمام تست‌های فعلی و بخش بزرگی از بانک تست برنامه، فقط تغییر دیتابیس کافی است.
