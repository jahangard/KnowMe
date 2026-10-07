# KnowMe Product & Content Principles

> این سند اصول ثابت طراحی محصول، محتوا و تجربه کاربر KnowMe است. هر Feature جدید باید با این اصول سازگار باشد.

## 1. تست باید جذاب باشد، نه شبیه فرم

KnowMe یک محصول تعاملی است، نه فرم ارزیابی اداری.

اصل:

- سؤال کوتاه
- جواب سریع
- زبان طبیعی
- حس کنجکاوی
- Completion Rate بالا
- نتیجه‌ای که کاربر بخواهد تا آخر ببیند

هر Test در MVP بهتر است معمولاً حدود 7 تا 12 سؤال داشته باشد.

## 2. نتیجه باید Share-worthy باشد

Result فقط یک Score خام نیست.

یک Result خوب باید:

- اسم جذاب داشته باشد
- یک جمله Hook داشته باشد
- توضیح کوتاه و قابل‌فهم بدهد
- امکان مقایسه ایجاد کند
- برای Screenshot/Share مناسب باشد
- بیش از حد منفی یا تحقیرآمیز نباشد

در آینده Share Card باید بخشی جدی از Growth loop محصول باشد.

## 3. زن و مرد تا جای ممکن از یک Test مشترک استفاده کنند

به‌طور پیش‌فرض Testها Gender-neutral طراحی شوند.

فقط زمانی سؤال یا نتیجه Gender-specific شود که واقعاً ماهیت موضوع نیاز داشته باشد.

## 4. Entertainment ≠ Diagnosis

بخش مهمی از Testهای KnowMe سرگرمی/خودشناسی هستند.

نباید ادعا شود:

- تشخیص روان‌شناختی قطعی
- تشخیص پزشکی
- اختلال
- طبیعی/غیرطبیعی
- پیش‌بینی قطعی آینده

برای Testهای سرگرمی، Disclaimer روشن کافی است.

## 5. تست ۱۸+ باید مسئولانه باشد

قواعد:

- فقط Age >= 18
- موضوعات به‌صورت نگرش، ترجیح، سبک و طیف بیان شوند
- از Label تحقیرآمیز یا پزشکی اجتناب شود
- Mobile و اطلاعات هویتی نباید برای نمایش Result لازم باشند
- محتوای حساس باید Optional باشد

## 6. Progressive Profiling

اطلاعات شخصی یک‌جا جمع نمی‌شوند.

اصل:

**هر بار فقط اطلاعاتی بخواه که همین لحظه برای Experience ارزش دارد.**

Flow فعلی:

- Gender
- Age
- Test
- Name after first completion
- Mobile later and optional

## 7. Database is the Content Source of Truth

Test content باید در DB باشد:

- Categories
- Tests
- Questions
- Options
- Scores
- Traits
- Result profiles

کد فقط Engine عمومی است.

بعد از Bootstrap، Restart نباید تغییرات Content در DB را بازنویسی کند.

## 8. No Test-specific Hard-code

ممنوع:

- switch روی TestCode برای متن Result
- Questionهای ثابت در Telegram code
- Optionهای ثابت در UI
- Result title ثابت در UI
- Scoring branch مخصوص یک Test

اگر یک قابلیت جدید فقط برای یک Test لازم است، ابتدا بررسی شود آیا می‌توان آن را به Rule عمومی تبدیل کرد.

## 9. TraitKey باید Namespace داشته باشد

برای Testهای جدید:

```text
attraction:magnetic
dating:planner
boundaries:firm
truth:approval
```

این کار از Collision جلوگیری می‌کند.

Exception:

`love_style_v1` برای backward compatibility از:

- words
- time
- care
- touch

استفاده می‌کند.

## 10. Raw Answer مهم‌تر از Derived Result است

هر Answer باید ذخیره شود.

Result و UserTrait داده مشتق‌شده هستند و باید در صورت نیاز از Answer history قابل بازسازی باشند.

## 11. Roadmap باید هوشمند شود، ولی قابل توضیح بماند

در نسخه فعلی Roadmap ساده است.

در آینده Recommendation باید از مواردی مثل:

- completed tests
- trait profile
- confidence
- abandon history
- age
- gender applicability
- category diversity

استفاده کند.

اما Recommendation نباید Black Box غیرقابل توضیح باشد.

## 12. دو مسیر ورود به Test

User باید همیشه دو انتخاب داشته باشد:

### انتخاب آزاد
موضوع → Test → Start

### نقشه راه
سیستم Test مناسب بعدی را انتخاب کند.

هیچ‌کدام نباید دیگری را حذف کند.

## 13. Category Browser باید ساده باشد

Categoryهای اصلی:

- ❤️ عشق و رابطه
- 🎯 چالشی و باحال
- 🧠 خودشناسی عمیق‌تر
- 🔞 ۱۸+

و کنار آن‌ها:

- 🧭 نقشه راه من

Roadmap یک Virtual Category است.

## 14. Chat نباید شلوغ شود

در Telegram:

- Inline navigation با Edit Message
- Question transition با Edit Message
- Back/Home با Edit در صورت امکان
- فقط Promptهایی که واقعاً نیاز به Text/Contact دارند Message جدا بدهند

## 15. هر تست باید یک وعده روشن داشته باشد

عنوان و Description باید به سؤال کاربر جواب دهند:

**«اگر این تست را انجام بدهم، چه چیز جالبی درباره خودم می‌فهمم؟»**

عنوان‌های کنجکاوی‌برانگیز بر عنوان‌های دانشگاهی خشک اولویت دارند.

مثال بهتر:

«وقتی کسی تأییدت نمی‌کند، چقدر هنوز خودت را دوست داری؟»

به‌جای:

«مقیاس عزت نفس»

## 16. تست باید نتیجه متنوع داشته باشد

Resultها نباید صرفاً Good/Bad باشند.

ترجیح:

- Archetype
- Style
- Pattern
- Spectrum
- Tendency

این باعث Shareability و کاهش حس قضاوت می‌شود.

## 17. سؤال‌ها نباید جواب «درست» را لو بدهند

اگر User بتواند خیلی راحت بفهمد کدام گزینه «بهتر» است، Test تبدیل به Social desirability survey می‌شود.

سؤال‌ها بهتر است:

- موقعیتی
- trade-off based
- غیرمستقیم
- بدون گزینه واضحاً اخلاقی/غیراخلاقی

باشند.

## 18. Result باید از Answerها قابل دفاع باشد

حتی در Test سرگرمی، Result نباید تصادفی به نظر برسد.

User باید وقتی Result را می‌خواند حس کند:

«آره، از جواب‌هایی که دادم منطقیه.»

## 19. Seed فقط Bootstrap است

Seed برای راه‌اندازی دیتابیس خالی است.

Seed نباید:

- DB Content موجود را overwrite کند
- Editor آینده را بی‌اثر کند
- Source of Truth دوم بسازد

## 20. قابلیت آینده: Admin/Test Studio

هدف طبیعی معماری Data-driven این است که بعداً بتوان Admin UI ساخت برای:

- ساخت Category
- ساخت Test
- سؤال‌ها
- Optionها
- Score
- Result Profile
- Preview
- Publish/Unpublish

بدون Deploy مجدد Bot.
