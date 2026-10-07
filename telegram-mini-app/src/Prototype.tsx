import { useEffect, useMemo, useState } from "react";
import { MobileScroll } from "./mobile";

type Trait = string;
type Screen = "home" | "topics" | "category" | "test-items" | "roadmap" | "lessons" | "lesson-list" | "lesson-detail" | "profile" | "quiz" | "result";
type CategoryCode = "love" | "challenge" | "self_knowledge" | "adult";

type AnswerOption = { text: string; trait: Trait };
type Question = { text: string; options: AnswerOption[] };
type TestItem = {
  id: string;
  title: string;
  description: string;
  category: CategoryCode;
  available?: boolean;
  method?: string;
  specification?: string;
};
type LifeLesson = { id: string; title: string; summary: string; goal: string; explanation: string; consequences: string; ageRange: string; method: string[]; help: string; takeaway: string };
type TopicNode = { id: string; title: string; caption: string; children?: TopicNode[]; testIds?: string[]; lessonIds?: string[] };

const lifeLessons: LifeLesson[] = [
  {
    id: "grief", title: "چطور با غم‌های بزرگ کنار بیاییم؟", summary: "برای روزهایی که اندوه، همه‌چیز را سنگین می‌کند.",
    goal: "قرار نیست غم را یک‌شبه حذف کنی؛ هدف این است که در کنار این احساس، قدم‌های کوچک و امن برداری.",
    explanation: "سوگ بعد از فقدان، جدایی، بیماری یا تغییرهای بزرگ می‌تواند سراغ هرکسی بیاید. گریه، خشم، خستگی یا بی‌حسی، شکل‌های متفاوت این تجربه‌اند و زمان‌بندی ثابتی برای بهتر شدن وجود ندارد.",
    consequences: "نادیده گرفتن طولانی‌مدت احساس‌ها ممکن است روی خواب، تمرکز و کارهای روزمره اثر بگذارد. این نشانه ضعف نیست؛ شاید یعنی به حمایت بیشتری نیاز داری.",
    ageRange: "برای نوجوانان و بزرگسالان؛ نوجوانان بهتر است یک بزرگسال امن را در جریان بگذارند.",
    method: ["احساست را نام ببر، بدون اینکه خودت را قضاوت کنی.", "امروز فقط یک نیاز پایه را انجام بده: آب، غذای ساده یا کمی استراحت.", "از یک آدم امن درخواست مشخصی بکن؛ مثلاً «می‌شود کمی کنارم بمانی؟»", "راه شخصی خودت برای یادآوری یا خداحافظی را پیدا کن."],
    help: "اگر انجام کارهای روزمره برای مدتی دشوار مانده، با مشاور یا پزشک صحبت کن. اگر فکر آسیب زدن به خودت داری، تنها نمان و فوراً از فردی مورد اعتماد و خدمات اضطراری محل زندگی‌ات کمک بگیر.",
    takeaway: "لازم نیست امروز خوب شوی؛ فقط قدم بعدیِ امن را بردار.",
  },
  {
    id: "letting-go", title: "چطور کسی را فراموش کنیم؟", summary: "راهی مهربانانه برای عبور از دلتنگی و پایان رابطه.",
    goal: "هدف پاک کردن خاطره نیست؛ کمک به توست تا خاطره کم‌کم اختیار امروزت را کمتر بگیرد.",
    explanation: "بعد از پایان رابطه، دلتنگی و دوگانگی طبیعی است. هرکس با سرعت خودش از این دوره عبور می‌کند.",
    consequences: "پیگیری مداوم صفحه‌ها یا امیدهای مبهم ممکن است پذیرش پایان را دشوارتر کند و خواب و تمرکز را تحت فشار بگذارد؛ این موضوع به معنی سرزنش تو نیست.",
    ageRange: "برای نوجوانان و بزرگسالان؛ در رابطه‌ی آسیب‌زا یا ناامن، از یک بزرگسال قابل اعتماد کمک بگیر.",
    method: ["واقعیت پایان را با خودت روشن و مهربانانه مرور کن.", "برای مدتی دیدن صفحه‌ها و پیام‌ها را محدود کن.", "وقتی میل ناگهانی به پیام دادن می‌آید، کمی مکث کن و با دوستی امن حرف بزن.", "خوبی‌ها و دشواری‌های رابطه را کنار هم ببین و برنامه‌های خودت را آرام‌آرام بساز."],
    help: "اگر جدایی برای مدتی طولانی خواب، تحصیل یا احساس امنیتت را مختل کرده، با مشاور یا فردی قابل اعتماد صحبت کن. در خطر فوری با خدمات اضطراری تماس بگیر.",
    takeaway: "قرار نیست یک‌شبه فراموش کنی؛ امروز می‌توانی یک مرز کوچک برای خودت بسازی.",
  },
];

declare global {
  interface Window {
    Telegram?: {
      WebApp?: {
        initDataUnsafe?: { user?: { first_name?: string } };
        ready?: () => void;
        expand?: () => void;
        enableClosingConfirmation?: () => void;
        disableClosingConfirmation?: () => void;
        BackButton?: {
          show: () => void;
          hide: () => void;
          onClick: (handler: () => void) => void;
          offClick: (handler: () => void) => void;
        };
        HapticFeedback?: { selectionChanged: () => void; notificationOccurred: (kind: "success") => void };
        openTelegramLink?: (url: string) => void;
        themeParams?: { bg_color?: string; text_color?: string; secondary_bg_color?: string };
        colorScheme?: "light" | "dark";
        onEvent?: (event: "themeChanged", handler: () => void) => void;
        offEvent?: (event: "themeChanged", handler: () => void) => void;
      };
    };
  }
}

type TraitProfile = { label: string; title: string; subtitle: string; description: string };

const traitProfiles: Record<Trait, TraitProfile> = {
  words: {
    label: "کلمات",
    title: "عاشقِ کلمات",
    subtitle: "برای تو، حرف خوب فقط حرف نیست.",
    description: "ابراز مستقیم احساس، تعریف و جمله‌های صمیمی خیلی زود به قلبت راه پیدا می‌کنن. احتمالاً خودت هم وقتی کسی برات مهمه، از کلمات برای نشان‌دادن علاقه استفاده می‌کنی.",
  },
  time: {
    label: "حضور",
    title: "عاشقِ حضور",
    subtitle: "برای تو، وقت گذاشتن یعنی انتخاب کردن.",
    description: "حضور واقعی، توجه بدون حواس‌پرتی و وقت دونفره بیشتر از کارهای نمایشی روی تو اثر می‌ذاره. وقتی کسی زمانش رو به تو می‌ده، احساس ارزشمندی بیشتری می‌کنی.",
  },
  care: {
    label: "عمل و مراقبت",
    title: "عاشقِ عمل",
    subtitle: "برای تو، دوست داشتن باید دیده بشه.",
    description: "کمک کردن، مسئولیت برداشتن و کارهای کوچک واقعی برای تو معنی زیادی دارن. احتمالاً بیشتر به رفتار نگاه می‌کنی تا وعده‌ها.",
  },
  touch: {
    label: "نزدیکی",
    title: "عاشقِ نزدیکی",
    subtitle: "برای تو، فاصله کم یعنی احساس بیشتر.",
    description: "آغوش، تماس و نزدیکی فیزیکی محترمانه برای تو یکی از روشن‌ترین نشانه‌های محبت و امنیت عاطفیه.",
  },
};

const loveQuestions: Question[] = [
  {
    text: "وقتی کسی که دوستش داری ناراحته، معمولاً اولین واکنش تو چیه؟",
    options: [
      { text: "با حرف زدن آرومش می‌کنم", trait: "words" },
      { text: "کنارش می‌مونم و وقتم رو بهش می‌دم", trait: "time" },
      { text: "یه کاری براش انجام می‌دم که حالش بهتر شه", trait: "care" },
      { text: "با آغوش و نزدیکی بهش آرامش می‌دم", trait: "touch" },
    ],
  },
  {
    text: "اگر بخوای خیلی واضح نشون بدی که کسی برات مهمه، بیشتر چه کار می‌کنی؟",
    options: [
      { text: "بهش می‌گم چقدر برام مهمه", trait: "words" },
      { text: "یه زمان مخصوص فقط برای دوتامون می‌ذارم", trait: "time" },
      { text: "کارش رو سبک می‌کنم یا کمکش می‌کنم", trait: "care" },
      { text: "بیشتر بغلش می‌کنم و نزدیکش می‌مونم", trait: "touch" },
    ],
  },
  {
    text: "در یک روز شلوغ، کدوم رفتار طرف مقابل بیشتر به دلت می‌شینه؟",
    options: [
      { text: "یک پیام محبت‌آمیز و صمیمی", trait: "words" },
      { text: "اینکه با وجود شلوغی، وقت برای من باز کنه", trait: "time" },
      { text: "اینکه بدون گفتن، یک کارم رو انجام بده", trait: "care" },
      { text: "یک بغل گرم وقتی همدیگه رو می‌بینیم", trait: "touch" },
    ],
  },
  {
    text: "وقتی دلخور می‌شی، کدوم کار بیشتر کمک می‌کنه دوباره احساس نزدیکی کنی؟",
    options: [
      { text: "اینکه حرف دلش رو واضح بگه", trait: "words" },
      { text: "اینکه بشینیم و باهم وقت بگذرونیم", trait: "time" },
      { text: "اینکه برای جبران، یک کار واقعی انجام بده", trait: "care" },
      { text: "اینکه با یک آغوش صمیمی فاصله رو کم کنه", trait: "touch" },
    ],
  },
  {
    text: "برای یک مناسبت خاص، کدوم برنامه بیشتر تو رو خوشحال می‌کنه؟",
    options: [
      { text: "یک نامه یا پیام خاص و احساسی", trait: "words" },
      { text: "یک روز کامل فقط با هم بودن", trait: "time" },
      { text: "یک کار غافلگیرکننده که زندگی‌م رو راحت‌تر کنه", trait: "care" },
      { text: "یک شب صمیمی و پر از نزدیکی", trait: "touch" },
    ],
  },
  {
    text: "وقتی از کسی خوشت میاد، خودت ناخودآگاه بیشتر کدوم رفتار رو انجام می‌دی؟",
    options: [
      { text: "زیاد تعریف و ابراز احساس می‌کنم", trait: "words" },
      { text: "دنبال فرصت می‌گردم باهاش تنها باشم", trait: "time" },
      { text: "کارهای کوچیکش رو انجام می‌دم", trait: "care" },
      { text: "با تماس و نزدیکی علاقه‌م رو نشون می‌دم", trait: "touch" },
    ],
  },
  {
    text: "در یک رابطه، کدوم کمبود بیشتر اذیتت می‌کنه؟",
    options: [
      { text: "کمبود حرف‌های محبت‌آمیز", trait: "words" },
      { text: "کمبود وقت دونفره", trait: "time" },
      { text: "اینکه همه چیز فقط در حد حرف بمونه", trait: "care" },
      { text: "فاصله و سردی فیزیکی", trait: "touch" },
    ],
  },
  {
    text: "اگر فقط یکی رو انتخاب کنی، کدوم جمله بیشتر حس دوست‌داشتن بهت می‌ده؟",
    options: [
      { text: "«دوستت دارم و بهت افتخار می‌کنم»", trait: "words" },
      { text: "«امروز رو کامل برای تو خالی کردم»", trait: "time" },
      { text: "«نگران نباش، من انجامش دادم»", trait: "care" },
      { text: "«بیا یه بغل طولانی»", trait: "touch" },
    ],
  },
];

const tests: TestItem[] = [
  { id: "love_style_v1", title: "سبک عشق‌ورزی", description: "شیوه غالب ابراز علاقه و دریافت محبتت را پیدا کن.", category: "love" },
  { id: "attraction_style_v1", title: "تیپ جذابیت", description: "ببین جذابیت تو بیشتر از چه جنسیه؛ حضور، گرما، رازآلودگی یا انرژی.", category: "love" },
  { id: "dating_scenarios_v1", title: "سناریوهای قرار", description: "در موقعیت‌های واقعی قرار و آشنایی، غریزه تو چطور تصمیم می‌گیره؟", category: "love" },
  { id: "personal_boundaries_v1", title: "مرزهای شخصی", description: "ببین وقتی پای نه گفتن، احترام و فضای شخصی وسطه، سبک تو چیه.", category: "love" },
  { id: "bitter_truth_v1", title: "حقیقت تلخ", description: "یک تست چالشی برای پیدا کردن الگویی که شاید درباره خودت کمتر دوست داشته باشی ببینی.", category: "challenge" },
  { id: "mbti_v1", title: "تست شخصیت‌شناسی MBTI", description: "ترجیحاتت در دریافت انرژی، توجه، تصمیم‌گیری و سبک زندگی را بررسی کن.", category: "self_knowledge", available: false, method: "پرسش‌نامه‌ی خودشناسی بر پایه‌ی چهار ترجیح: درون‌گرایی/برون‌گرایی، حسی/شهودی، منطقی/احساسی و قضاوت‌گر/ادراکی؛ نتیجه یکی از ۱۶ تیپ است.", specification: "۶۰ سؤال · حدود ۱۵ دقیقه" },
  { id: "disc_v1", title: "تست شخصیت DISC", description: "الگوی رفتاری غالب تو در تعامل با دیگران و موقعیت‌های کاری.", category: "self_knowledge", available: false, method: "ترجیحت در دو محورِ کار در برابر افراد و سرعت بالا در برابر سرعت متوسط بررسی می‌شود؛ خروجی چهار گرایش D، I، S و C و ترکیب‌های شخصیتی را نشان می‌دهد." },
  { id: "cattell_16_v1", title: "تست شخصیت ۱۶ عاملی کتل", description: "نگاهی چندبعدی به ویژگی‌های شخصیتی و تفاوت‌های فردی.", category: "self_knowledge", available: false },
  { id: "neo_v1", title: "تست شخصیت نئو (NEO-FFI)", description: "شناخت پنج بُعد اصلی شخصیت و بررسی ارتباط آن‌ها با ترجیحات شغلی.", category: "self_knowledge", available: false, method: "پرسش‌نامه‌ی خودگزارشی؛ پنج عامل روان‌رنجوری، برون‌گرایی، توافق‌پذیری، گشودگی به تجربه و وظیفه‌شناسی سنجیده می‌شوند.", specification: "۶۰ سؤال · حدود ۱۵ دقیقه" },
  { id: "organizational_commitment_v1", title: "تست تعهد سازمانی", description: "میزان پیوند و تعهدت به سازمان و محیط کاری را بررسی کن.", category: "self_knowledge", available: false },
  { id: "job_satisfaction_v1", title: "تست رضایت شغلی", description: "نگاهت به وظایف، محیط کار و تجربه‌ی شغلی‌ات را ارزیابی کن.", category: "self_knowledge", available: false },
  { id: "job_burnout_v1", title: "تست فرسودگی شغلی", description: "نشانه‌های خستگی و فشار مزمن در تجربه‌ی کاری را بررسی کن.", category: "self_knowledge", available: false },
  { id: "entrepreneur_personality_v1", title: "تست شخصیت کارآفرین", description: "گرایش‌ها و توانمندی‌های شخصیتی مرتبط با کارآفرینی.", category: "self_knowledge", available: false },
  { id: "vocational_interest_v1", title: "تست رغبت‌سنج شغلی", description: "زمینه‌های کاری و فعالیت‌هایی را پیدا کن که بیشتر به آن‌ها علاقه داری.", category: "self_knowledge", available: false },
  { id: "raven_adult_v1", title: "تست هوش ریون بزرگسالان", description: "استدلال قیاسی، درک مفاهیم انتزاعی و ادراک را بررسی کن.", category: "self_knowledge", available: false, method: "۶۰ پرسش چندگزینه‌ای تصویری از آسان به دشوار؛ با کامل‌کردن ماتریس‌های هندسی، استدلال غیرکلامی و تشخیص رابطه‌ی میان شکل‌ها سنجیده می‌شود.", specification: "۶۰ سؤال · حدود ۴۵ دقیقه" },
  { id: "gardner_multiple_intelligences_v1", title: "تست هوش‌های چندگانه گاردنر", description: "توانمندی‌های برجسته‌ات را در حوزه‌های مختلف هوش بشناس.", category: "self_knowledge", available: false },
  { id: "cattell_adult_a_v1", title: "تست هوش کتل بزرگسال A", description: "ارزیابی توانایی استدلال و هوش سیال بزرگسالان، فرم A.", category: "self_knowledge", available: false },
  { id: "cattell_adult_b_v1", title: "تست هوش کتل بزرگسال B", description: "ارزیابی توانایی استدلال و هوش سیال بزرگسالان، فرم B.", category: "self_knowledge", available: false },
  { id: "bonnardel_v1", title: "تست هوش بوناردل", description: "توانایی تحلیل و استدلال در الگوهای تصویری را بررسی کن.", category: "self_knowledge", available: false },
  { id: "baron_emotional_intelligence_v1", title: "تست هوش هیجانی بارآن", description: "مجموعه‌ای از شایستگی‌های هیجانی و اجتماعی را بررسی کن.", category: "self_knowledge", available: false, method: "پرسش‌نامه‌ی خودگزارشی در ۱۵ بُعد، از جمله خودآگاهی هیجانی، همدلی، روابط بین‌فردی، حل مسئله و تحمل فشار روانی.", specification: "۹۰ سؤال · حدود ۲۵ دقیقه" },
  { id: "schutte_emotional_intelligence_v1", title: "تست هوش هیجانی شات", description: "برداشت و مدیریت هیجان‌ها در خودت و دیگران را ارزیابی کن.", category: "self_knowledge", available: false },
  { id: "child_intelligence_v1", title: "تست هوش کودک", description: "ارزیابی متناسب با سن کودک؛ جزئیات اجرا و گروه سنی پیش از انتشار تکمیل می‌شود.", category: "self_knowledge", available: false },
  { id: "sternberg_triangular_love_v1", title: "تست مثلث عشق استرنبرگ", description: "صمیمیت، شور و تعهد را در تجربه‌ی رابطه بررسی کن.", category: "self_knowledge", available: false },
  { id: "premarital_fears_v1", title: "تست ترس‌های قبل از ازدواج", description: "نگرانی‌ها و دغدغه‌های پیش از تصمیم به ازدواج را بررسی کن.", category: "self_knowledge", available: false },
  { id: "enrich_marital_satisfaction_v1", title: "تست رضایت زناشویی انریچ (ENRICH)", description: "ابعاد گوناگون رضایت و کیفیت رابطه‌ی زناشویی را ارزیابی کن.", category: "self_knowledge", available: false },
  { id: "emotional_intelligence_v1", title: "تست هوش هیجانی", description: "شناخت و مدیریت هیجان‌ها را در زندگی روزمره بررسی کن.", category: "self_knowledge", available: false },
  { id: "loneliness_v1", title: "تست احساس تنهایی", description: "تجربه‌ی تنهایی و احساس پیوند اجتماعی‌ات را بررسی کن.", category: "self_knowledge", available: false },
  { id: "general_health_v1", title: "تست سلامت عمومی (GHQ)", description: "وضعیت سلامت عمومی‌ات را در بازه‌ی یک ماه گذشته مرور کن.", category: "self_knowledge", available: false, method: "پرسش‌نامه‌ی خودگزارشی درباره‌ی تجربه‌ی سلامت عمومی در یک ماه اخیر." },
  { id: "social_anxiety_v1", title: "تست اضطراب اجتماعی", description: "میزان نگرانی و دشواری در موقعیت‌های اجتماعی را بررسی کن.", category: "self_knowledge", available: false },
  { id: "depression_v1", title: "تست افسردگی", description: "نشانه‌های خلقی را برای خودآگاهی اولیه بررسی کن.", category: "self_knowledge", available: false },
  { id: "adapto_career_path_v1", title: "تست مسیر شغلی ادپتو (Career Path)", description: "ترکیبی از چند ارزیابی برای پیدا کردن مسیرها و شغل‌های متناسب با تو.", category: "self_knowledge", available: false, method: "چهار بُعد شخصیت شغلی، علایق، استعدادها و ارزش‌های شغلی را هم‌زمان می‌سنجد و بر اساس تطبیق آن‌ها مسیرها و شغل‌های پیشنهادی می‌دهد.", specification: "۴ ارزیابی · حدود ۴۰ دقیقه" },
  { id: "dmsi_motivation_v1", title: "تست سبک انگیزشی غالب (DMSI)", description: "محرک اصلی انگیزه‌ات را در کار و رشد شناسایی کن.", category: "self_knowledge", available: false, method: "بر پایه‌ی نیازهای انگیزشی مک‌کللند؛ پاسخ‌ها گرایش غالب به پیشرفت، قدرت یا پیوندجویی را نشان می‌دهند.", specification: "۱۵ سؤال" },
  { id: "tki_conflict_style_v1", title: "تست سبک حل تعارض توماس–کیلمن (TKI)", description: "واکنش معمولت را هنگام اختلاف و تعارض با دیگران بررسی کن.", category: "self_knowledge", available: false, method: "پرسش‌نامه‌ی موقعیت‌محور درباره‌ی تعارض‌های بین‌فردی؛ نتیجه پنج سبک همکاری، رقابت، مصالحه، سازش و اجتناب را می‌سنجد.", specification: "۳۰ سؤال" },
  { id: "csi_stress_coping_v1", title: "تست سبک‌های مقابله با استرس (CSI)", description: "روش‌هایی را که برای کنار آمدن با موقعیت‌های پراسترس به کار می‌بری بشناس.", category: "self_knowledge", available: false, method: "پرسش‌نامه‌ی خودگزارشی درباره‌ی واکنش به موقعیت‌های استرس‌زا؛ هشت سبک از حمایت اجتماعی و حل مسئله تا اجتناب اجتماعی را بررسی می‌کند.", specification: "۷۲ سؤال · حدود ۱۰ دقیقه" },
  { id: "cognitive_abilities_v1", title: "تست توانایی‌های شناختی (Cognitive Abilities)", description: "توانایی‌های شناختی مرتبط با حل مسئله و عملکرد حرفه‌ای را بررسی کن.", category: "self_knowledge", available: false, method: "مجموعه‌ای از سؤال‌ها برای ارزیابی هفت حوزه، از جمله انعطاف‌پذیری شناختی، توجه، برنامه‌ریزی، تصمیم‌گیری، کنترل مهاری و حافظه.", specification: "۳۰ سؤال" },
  { id: "hartman_personality_v1", title: "تست شخصیت هارتمن", description: "ارزش‌های محوری و الگوی شخصیتت را در چهار رنگ بشناس.", category: "self_knowledge", available: false, method: "پرسش‌نامه‌ی خودشناسی مبتنی بر ارزش‌های محوری؛ نتیجه در چهار رنگ قرمز، آبی، سفید و زرد و چهار گرایش قدرت، صمیمیت، آرامش و شوخ‌طبعی گزارش می‌شود.", specification: "۴۵ سؤال · حدود ۱۰ دقیقه" },
  { id: "career_anchors_v1", title: "تست لنگرگاه‌های شغلی (Career Anchors)", description: "ارزش‌ها و اولویت‌هایی را پیدا کن که جهت انتخاب‌های شغلی‌ات را می‌دهند.", category: "self_knowledge", available: false, method: "خودارزیابیِ شایستگی‌ها، ارزش‌ها و انگیزه‌ها؛ نتیجه نشان می‌دهد کدام‌یک از هشت گرایش شغلی برایت پررنگ‌تر است.", specification: "۴۰ سؤال · حدود ۱۰ دقیقه" },
  { id: "neo_pir_v1", title: "تست شخصیت نئو (NEO PI-R)", description: "پنج عامل بزرگ شخصیت را با ارزیابی گسترده‌تر بررسی کن.", category: "self_knowledge", available: false, method: "پرسش‌نامه‌ی خودگزارشی پنج‌عاملی؛ روان‌رنجوری، برون‌گرایی، توافق‌پذیری، گشودگی به تجربه و وظیفه‌شناسی را می‌سنجد.", specification: "۲۴۰ سؤال · حدود ۶۰ دقیقه" },
  { id: "clifton_strengths_v1", title: "تست استعدادیابی کلیفتون (گالوپ)", description: "نقاط قوت و استعدادهای برجسته‌ات را برای رشد فردی و کاری شناسایی کن.", category: "self_knowledge", available: false, method: "ارزیابی نقاط قوت که ۳۴ استعداد را رتبه‌بندی می‌کند؛ گزارش، استعدادهای برتر و کاربردها و نقش‌های سازگار را توضیح می‌دهد." },
  { id: "holland_vocational_interest_v1", title: "تست رغبت شغلی هالند (RIASEC)", description: "علاقه‌های کاری‌ات را با محیط‌ها و مسیرهای شغلی هماهنگ کن.", category: "self_knowledge", available: false, method: "پرسش‌نامه‌ی رغبت‌سنجی بر پایه‌ی شش تیپ هالند: واقع‌گرا، جست‌وجوگر، هنری، اجتماعی، متهور و قراردادی." },
];

const extraTestContent: Record<string, { traits: string[]; questions: Question[]; profiles: Record<string, TraitProfile> }> = {
  attraction_style_v1: {
    traits: ["magnetic", "warm", "mysterious", "playful"],
    profiles: {
      magnetic: { label: "مغناطیسی", title: "جذابیت مغناطیسی", subtitle: "لازم نیست زیاد تلاش کنی؛ حضورت دیده می‌شه.", description: "اعتمادبه‌نفس، قاطعیت و حضور تو باعث می‌شه دیگران متوجهت بشن. جذابیتت بیشتر از جنس حضور واقعیه تا نمایش." },
      warm: { label: "گرم", title: "جذابیت گرم و صمیمی", subtitle: "کنارت بودن حس راحتی می‌ده.", description: "لبخند و توجه تو فضایی می‌سازه که آدم‌ها زودتر خود واقعی‌شون می‌شن. صمیمیت و پذیرا بودنت نقطه قوت توئه." },
      mysterious: { label: "رازآلود", title: "جذابیت رازآلود", subtitle: "همه‌چیز را همان اول رو نمی‌کنی.", description: "تو با استقلال و شناخت تدریجی کنجکاوی ایجاد می‌کنی. جذابیتت در اینه که برای شناختنت باید کمی بیشتر وقت گذاشت." },
      playful: { label: "بازیگوش", title: "جذابیت پرانرژی", subtitle: "با تو فضا زنده‌تر می‌شه.", description: "شوخ‌طبعی و انرژی تو شروع ارتباط را آسان می‌کنه. دیگران کنار تو حس می‌کنن می‌تونن راحت‌تر و خودمانی‌تر باشن." },
    },
    questions: [
      q("وقتی وارد جمع تازه‌ای می‌شی، معمولاً چطور شروع می‌کنی؟", "با آرامش و اعتمادبه‌نفس خودم رو معرفی می‌کنم", "زود لبخند می‌زنم و گرم می‌گیرم", "اول فضا رو می‌سنجم و کم‌کم وارد می‌شم", "با شوخی یخ جمع رو آب می‌کنم"),
      q("دیگران معمولاً بابت چی ازت تعریف می‌کنن؟", "کاریزما و حضورم", "راحت بودن کنارم", "خاص بودنم", "انرژی و شوخ‌طبعی‌م"),
      q("در گفت‌وگوی دونفره کدوم سبک به تو نزدیک‌تره؟", "مستقیم و مطمئن حرف می‌زنم", "با توجه و مهربونی گوش می‌دم", "همه چیز رو یک‌باره نمی‌گم", "با شوخی و بازی کلامی پیش می‌رم"),
      q("در پیام دادن چطور علاقه‌ات رو نشون می‌دی؟", "کم ولی روشن و مطمئن", "صمیمی و پیگیر حال طرف", "کم‌حرف و بافاصله", "با شوخی و جواب‌های پرانرژی"),
      q("وقتی توجه جمع به تو جلب می‌شه؟", "راحت فضا رو مدیریت می‌کنم", "توجه رو با بقیه تقسیم می‌کنم", "ترجیح می‌دم بخشی از خودم خصوصی بمونه", "ازش برای بامزه‌تر کردن فضا استفاده می‌کنم"),
      q("جذابیت خودت رو چطور توصیف می‌کنی؟", "قاطع و اثرگذار", "دوست‌داشتنی و پذیرا", "خاص و کمی غیرقابل‌پیش‌بینی", "خلاق و پرانرژی"),
      q("اگر کسی که دوستش داری سر صحبت رو باز کنه؟", "واضح و مطمئن جواب می‌دم", "با توجه و مهربونی جواب می‌دم", "کمی زمان می‌دم تا مطمئن شم", "با شیطنت و شوخی واکنش می‌دم"),
      q("دوست داری بعد از آشنایی چه چیزی از تو یادش بمونه؟", "اعتمادبه‌نفس و حضورم", "حس خوب و آرامشی که دادم", "کنجکاوی برای شناخت بیشترم", "خنده و انرژی‌ای که به فضا دادم"),
    ],
  },
  dating_scenarios_v1: {
    traits: ["planner", "spontaneous", "deep", "selective"],
    profiles: {
      planner: { label: "برنامه‌ریز", title: "قرارساز حسابگر", subtitle: "دوست داری بدانی کجا داری می‌ری.", description: "توجه، برنامه و قابل‌اتکا بودن کیفیت قرار را برای تو بالا می‌بره. غافلگیری خوبه، وقتی چارچوب کلی روشن باشه." },
      spontaneous: { label: "ماجراجو", title: "ماجراجوی لحظه‌ای", subtitle: "بهترین قرار شاید از قبل قرار نباشه.", description: "تجربه تازه و تصمیم‌های لحظه‌ای برات جذابن. وقتی فضا طبیعی و بدون فشار پیش می‌ره، بیشتر خودت می‌شی." },
      deep: { label: "عمیق", title: "جوینده اتصال عمیق", subtitle: "قرار خوب یعنی یک گفت‌وگوی واقعی.", description: "تو بیشتر از ظاهر برنامه دنبال ارتباط ذهنی و عاطفی هستی. گفت‌وگوی واقعی برات از یک برنامه پرزرق‌وبرق مهم‌تره." },
      selective: { label: "انتخاب‌گر", title: "انتخاب‌گر دقیق", subtitle: "رفتار واقعی از حرف‌های قشنگ مهم‌تره.", description: "تو با دقت و احترام به نشانه‌ها جلو می‌ری. برایت امنیت، صداقت و هماهنگی حرف و عمل پایه‌های یک آشنایی خوبن." },
    },
    questions: [
      q("برای قرار اول، کدوم برنامه رو ترجیح می‌دی؟", "زمان و مکان مشخص و از قبل هماهنگ‌شده", "بریم بیرون و در مسیر تصمیم بگیریم", "جایی آروم برای حرف زدن", "فعالیتی که رفتار واقعی همدیگه رو ببینیم"),
      q("بهترین بخش یک قرار چیه؟", "اینکه همه‌چیز مرتب و بافکر پیش بره", "یک اتفاق غیرمنتظره و باحال", "گفت‌وگویی که از سطح معمولی رد بشه", "رفتار محترمانه و بدون نشانه نگران‌کننده"),
      q("بعد از قرار خوب چه کار می‌کنی؟", "برای قرار بعدی زمان پیشنهاد می‌دم", "می‌ذارم حس و حال خودش جلو بره", "یک پیام واقعی درباره حسم می‌فرستم", "کمی صبر می‌کنم تا رفتار بعدی رو ببینم"),
      q("اگر برنامه قرار ناگهان عوض بشه؟", "می‌پذیرم ولی برای بعد هماهنگی بیشتر می‌خوام", "اشکالی نداره، برنامه رو عوض می‌کنیم", "اگر توضیح صادقانه باشه مهم‌تره", "این رفتار رو کنار نشانه‌های دیگه می‌سنجم"),
      q("برای شناختن طرف مقابل کدوم قرار بهتره؟", "رستوران خوب با رزرو و برنامه", "قدم‌زدن و یک کافه تصادفی", "جایی آروم برای ساعت‌ها حرف زدن", "فعالیتی که رفتار واقعی رو نشان بده"),
      q("اگر رابطه سریع پیش بره؟", "سرعت رو مرحله‌به‌مرحله تنظیم می‌کنم", "اگر حسش خوب باشه می‌ذارم اتفاق بیفته", "می‌خوام بفهمم پشت احساس چی هست", "محتاط می‌شم تا حرف و عملش رو ببینم"),
      q("چه چیزی درباره یک قرار بیشتر خیالت رو راحت می‌کنه؟", "جزئیات و زمان‌بندی روشن", "اینکه لازم نباشه همه‌چیز برنامه‌ریزی بشه", "بتونیم راحت و صادقانه حرف بزنیم", "ثبات و احترام در رفتارش"),
      q("اگر بین دو نفر مردد باشی، چه چیزی تعیین‌کننده‌تره؟", "کسی که برای دیدارها برنامه و توجه می‌ذاره", "کسی که باهاش تجربه‌های تازه دارم", "کسی که گفت‌وگوی عمیق‌تری داریم", "کسی که رفتارش امن و قابل‌اعتماده"),
    ],
  },
  personal_boundaries_v1: {
    traits: ["firm", "flexible", "peacekeeper", "guarded"],
    profiles: {
      firm: { label: "روشن", title: "مرزبان روشن", subtitle: "نه گفتن برای تو بی‌احترامی نیست.", description: "معمولاً می‌دونی چه چیزی برات قابل‌قبوله و می‌تونی شفاف بیانش کنی. وضوح نقطه قوته؛ انعطاف رو با نادیده گرفتن نیازهای خودت اشتباه نگیر." },
      flexible: { label: "منعطف", title: "منعطف اما آگاه", subtitle: "هم مرز داری، هم جا برای مذاکره.", description: "تو شرایط و آدم‌ها رو در نظر می‌گیری و بعد تصمیم می‌گیری. این انعطاف نقطه قوته، تا وقتی نیازهای خودت گم نشن." },
      peacekeeper: { label: "صلح‌طلب", title: "صلح‌طلبِ بیش‌ازحد", subtitle: "گاهی آرام نگه داشتن فضا را به خودت ترجیح می‌دی.", description: "ممکنه برای حفظ رابطه بیشتر از حد لازم کوتاه بیای. صلح واقعی وقتی دوام میاره که مرزهای خودت هم محترم بمونن." },
      guarded: { label: "محتاط", title: "محافظ محتاط", subtitle: "برای اعتماد کردن به زمان نیاز داری.", description: "تو بااحتیاط از فضای شخصی و احساساتت مراقبت می‌کنی. این محافظت می‌تونه مفید باشه؛ ارتباط امن گاهی با گفتن تدریجی نیازها ساخته می‌شه." },
    },
    questions: [
      q("کسی ازت کاری می‌خواد که واقعاً وقتش رو نداری.", "محترمانه و روشن نه می‌گم", "اگر راهی باشه زمان دیگه‌ای پیشنهاد می‌دم", "قبول می‌کنم که ناراحت نشه", "جواب رو عقب می‌ندازم تا مطمئن شم"),
      q("دوستت بدون هماهنگی برنامه‌ات رو تغییر می‌ده.", "می‌گم این تغییر برای من مناسب نیست", "شرایط رو می‌سنجم و شاید همراهی کنم", "برای جلوگیری از بحث چیزی نمی‌گم", "دفعه بعد کمتر برنامه‌هام رو باهاش درمیون می‌ذارم"),
      q("کسی درباره موضوع خصوصی ازت سؤال می‌پرسه.", "می‌گم ترجیح می‌دم درباره‌اش حرف نزنم", "به اندازه‌ای که راحتم پاسخ می‌دم", "جواب می‌دم که فضا معذب نشه", "موضوع رو عوض می‌کنم"),
      q("در رابطه به کمی زمان تنهایی نیاز داری.", "نیازم رو واضح توضیح می‌دم", "زمانی پیدا می‌کنیم که برای هر دو خوب باشه", "صبر می‌کنم تا طرف خودش متوجه بشه", "فاصله می‌گیرم بدون اینکه توضیح بدم"),
      q("کسی با شوخی از خط قرمزت رد می‌شه.", "همون لحظه می‌گم این شوخی برام خوب نبود", "می‌گم می‌دونم قصد بدی نداشتی ولی تکرارش نکن", "می‌خندم تا تنش درست نشه", "بعد از این کمتر باهاش صمیمی می‌شم"),
      q("برای یک تصمیم مهم با تو مشورت نشده.", "می‌گم انتظار داشتم نظر من هم پرسیده بشه", "درباره دلیلش گفت‌وگو می‌کنم", "می‌پذیرم تا اختلاف پیش نیاد", "فعلاً چیزی نمی‌گم و فاصله می‌گیرم"),
      q("نه گفتن به کسی که دوستش داری برات سخته.", "با مهربونی اما قاطع نه می‌گم", "راه‌حل دیگری پیشنهاد می‌دم", "معمولاً قبول می‌کنم تا ناراحت نشه", "از جواب دادن طفره می‌رم"),
      q("چه چیزی در مرزهای رابطه برات سخت‌تره؟", "اینکه نه گفتن بی‌احترامی برداشت بشه", "اینکه نتونیم به توافق منصفانه برسیم", "اینکه نه گفتن باعث از دست دادن آدم‌ها بشه", "اینکه نزدیک شدن باعث آسیب بشه"),
    ],
  },
  bitter_truth_v1: {
    traits: ["approval", "control", "avoidance", "intensity"],
    profiles: {
      approval: { label: "تأییدطلب", title: "حقیقت تلخ: زیادی دنبال تأییدی", subtitle: "گاهی نظر بقیه صدای خودت را کم‌رنگ می‌کنه.", description: "ممکنه واکنش دیگران بیشتر از چیزی که فکر می‌کنی روی تصمیم‌هات اثر بذاره. پیش از پرسیدن نظر همه، از خودت بپرس واقعاً چی می‌خوای." },
      control: { label: "کنترل‌گر", title: "حقیقت تلخ: دوست داری کنترل دست تو باشه", subtitle: "ابهام و بی‌برنامگی زود خسته‌ات می‌کنه.", description: "توان مدیریت تو نقطه قوته؛ اما وقتی همه‌چیز باید طبق نقشه پیش بره، آزادی و خودجوشی دیگران کمتر جا پیدا می‌کنه." },
      avoidance: { label: "اجتنابی", title: "حقیقت تلخ: بعضی چیزها را عقب می‌اندازی", subtitle: "سکوت کوتاه‌مدت آرامش می‌ده، ولی مسئله محو نمی‌شه.", description: "شاید در تنش‌ها اول دنبال کم‌کردن فشار باشی و گفت‌وگوی لازم رو عقب بندازی. روبه‌رو شدن به‌موقع می‌تونه مسئله رو کوچک‌تر نگه داره." },
      intensity: { label: "همه یا هیچ", title: "حقیقت تلخ: گاهی همه‌چیز برایت صفر یا صده", subtitle: "وقتی چیزی برات مهم شه، نصفه‌نیمه بودن سخته.", description: "شدت احساس و تعهدت می‌تونه نقطه قوت باشه؛ نگاه همه‌یا‌هیچ گاهی اجازه نمی‌ده خاکستری‌ها و تغییر تدریجی رو ببینی." },
    },
    questions: [
      q("پیامت دیده شده ولی چند ساعت جوابی نیومده. ذهنت چی می‌گه؟", "نکنه چیزی گفتم که بد برداشت کرده؟", "دوست دارم دقیقاً بدونم کی جواب می‌ده", "بی‌خیال، منم فعلاً جواب نمی‌دم", "یا علاقه داره یا نداره؛ این وسط‌بازی رو دوست ندارم"),
      q("وقتی کسی ازت انتقاد می‌کنه؟", "فکر می‌کنم نکنه واقعاً بد دیده شدم", "دنبال دلیل می‌گردم تا بفهمم حق با کیه", "ترجیح می‌دم بحث رو عوض کنم", "اگر ناعادلانه باشه شدید واکنش می‌دم"),
      q("برنامه‌ای که براش ذوق داشتی ناگهان عوض می‌شه.", "شاید برای بقیه مهم نبودم", "نظم و کنترل ماجرا از دست رفت", "ترجیح می‌دم کلاً بی‌خیالش بشم", "از ذوق کامل به ناامیدی کامل می‌رسم"),
      q("کدوم رفتار خودت بعداً بیشتر حرصت می‌ده؟", "برای خوشحال کردن بقیه چیزی رو قبول کردم", "روی جزئیات و نتیجه زیادی پافشاری کردم", "حرف لازم رو نزدم و گذاشتم بگذره", "از روی احساس تصمیم قطعی گرفتم"),
      q("در یک اختلاف رابطه‌ای کدوم جمله شبیه ذهن توئه؟", "فقط می‌خوام مطمئن شم هنوز دوستم داره", "باید مشخص کنیم دقیقاً چطور ادامه بدیم", "الان حوصله بحث ندارم؛ بعداً", "اگر اینطوریه شاید اصلاً ادامه ندیم"),
      q("هنگام تصمیم مهم کدوم دام بیشتر سراغت میاد؟", "نظر همه رو می‌پرسم و گیج‌تر می‌شم", "آنقدر اطلاعات جمع می‌کنم که تصمیم عقب می‌افته", "تا مجبور نشم تصمیم رو عقب می‌ندازم", "یک لحظه مطمئن می‌شم و سریع می‌پرم"),
      q("کدوم تعریف پنهانی بیشتر خوشحالت می‌کنه؟", "همه دوستت دارن", "همیشه می‌دونی باید چی کار کرد", "هیچ‌وقت وارد تنش نمی‌شی", "برای چیزهایی که می‌خوای می‌جنگی"),
      q("کدوم عادت رو دوست داری کمتر کنی؟", "دنبال رضایت همه بودن", "کنترل کردن همه‌چیز", "عقب انداختن گفت‌وگوهای سخت", "تصمیم قطعی گرفتن بدون دیدن خاکستری‌ها"),
    ],
  },
};

function q(text: string, ...answers: [string, string, string, string]): Question {
  const traits = ["magnetic", "warm", "mysterious", "playful"];
  return { text, options: answers.map((answer, index) => ({ text: answer, trait: traits[index] })) };
}

const categories: { code: CategoryCode; title: string; caption: string }[] = [
  { code: "love", title: "عشق و رابطه", caption: "از دل تا رفتار" },
  { code: "challenge", title: "چالشی و باحال", caption: "برای وقت‌های متفاوت" },
  { code: "self_knowledge", title: "خودشناسی", caption: "سفر به دنیای من" },
  { code: "adult", title: "۱۸+", caption: "ویژه بزرگسالان" },
];

const testTopicTree: TopicNode[] = [
  { id: "love", title: "عشق و رابطه", caption: "از دل تا رفتار", children: [
    { id: "relationship-style", title: "شناخت رابطه و جذابیت", caption: "الگوهای ارتباط و نزدیک‌شدن", testIds: ["love_style_v1", "attraction_style_v1"] },
    { id: "dating", title: "آشنایی و قرار", caption: "انتخاب و رفتار در شروع رابطه", children: [
      { id: "dating-scenarios", title: "سناریوهای قرار", caption: "واکنش تو در موقعیت‌های واقعی", testIds: ["dating_scenarios_v1"] },
    ] },
    { id: "boundaries", title: "مرزهای شخصی", caption: "احترام، نه گفتن و فضای امن", testIds: ["personal_boundaries_v1"] },
  ] },
  { id: "challenge", title: "چالشی و باحال", caption: "برای وقت‌های متفاوت", children: [
    { id: "self-challenge", title: "روبه‌رو شدن با خود", caption: "الگوهایی که شاید کمتر ببینی", testIds: ["bitter_truth_v1"] },
  ] },
  { id: "self_knowledge", title: "خودشناسی", caption: "سفر به دنیای من", children: [
    { id: "emotions", title: "احساسات و تصمیم‌ها", caption: "این مسیر به‌زودی کامل‌تر می‌شود", testIds: [] },
    { id: "personality-types", title: "تیپ‌های شخصیتی", caption: "الگوهای متفاوت شخصیت", children: [
      { id: "personality-assessments", title: "آزمون‌های شخصیت", caption: "شناخت الگوها و ویژگی‌های شخصیتی", testIds: ["mbti_v1", "disc_v1", "cattell_16_v1", "neo_v1"] },
      { id: "personality-archetypes", title: "کهن‌الگوها", caption: "", testIds: [] },
    ] },
    { id: "career-self-knowledge", title: "خودشناسی شغلی", caption: "شناخت تجربه، گرایش و سبک کاری", testIds: ["adapto_career_path_v1", "organizational_commitment_v1", "job_satisfaction_v1", "job_burnout_v1", "entrepreneur_personality_v1", "disc_v1", "mbti_v1", "vocational_interest_v1", "holland_vocational_interest_v1", "neo_v1", "neo_pir_v1", "dmsi_motivation_v1", "tki_conflict_style_v1", "csi_stress_coping_v1", "cognitive_abilities_v1", "career_anchors_v1", "clifton_strengths_v1"] },
    { id: "intelligence-self-knowledge", title: "خودشناسی هوش", caption: "آشنایی با توانمندی‌های شناختی و هیجانی", testIds: ["raven_adult_v1", "gardner_multiple_intelligences_v1", "cattell_adult_a_v1", "cattell_adult_b_v1", "bonnardel_v1", "baron_emotional_intelligence_v1", "schutte_emotional_intelligence_v1", "child_intelligence_v1", "cognitive_abilities_v1"] },
    { id: "love-marriage-self-knowledge", title: "عشق و ازدواج", caption: "شناخت رابطه و آمادگی برای ازدواج", testIds: ["sternberg_triangular_love_v1", "premarital_fears_v1", "enrich_marital_satisfaction_v1", "cattell_16_v1"] },
    { id: "social-self-knowledge", title: "خودشناسی اجتماعی", caption: "هیجان‌ها، ارتباط و بهزیستی اجتماعی", testIds: ["emotional_intelligence_v1", "loneliness_v1", "premarital_fears_v1", "general_health_v1", "social_anxiety_v1", "depression_v1", "tki_conflict_style_v1", "csi_stress_coping_v1"] },
  ] },
  { id: "adult", title: "۱۸+", caption: "ویژه بزرگسالان", children: [
    { id: "adult-relationships", title: "رابطه‌ی بزرگسالان", caption: "محتوای ویژه‌ی بزرگسالان", testIds: [] },
  ] },
];

const lessonTopicTree: TopicNode[] = [
  { id: "feelings", title: "احساسات و سوگ", caption: "برای روزهایی که همه‌چیز سنگین می‌شود", children: [
    { id: "grief", title: "غم و فقدان", caption: "کنار آمدن با سوگ و تغییرهای بزرگ", lessonIds: ["grief"] },
  ] },
  { id: "relationships", title: "رابطه و دلتنگی", caption: "مراقبت از خود در رابطه‌ها", children: [
    { id: "breakups", title: "جدایی و دل‌کندن", caption: "عبور آرام از پایان رابطه", lessonIds: ["letting-go"] },
  ] },
];

function webApp() {
  return window.Telegram?.WebApp;
}

function knowmeAsset(file: string) {
  return `${import.meta.env.BASE_URL}assets/knowme/${file}`;
}

function faNumber(value: number | string) {
  return String(value).replace(/[0-9]/g, (digit) => "۰۱۲۳۴۵۶۷۸۹"[Number(digit)]);
}

export default function Prototype() {
  const [screen, setScreen] = useState<Screen>("home");
  const [category, setCategory] = useState<CategoryCode>("love");
  const [testTopicPath, setTestTopicPath] = useState<TopicNode[]>([]);
  const [lessonTopicPath, setLessonTopicPath] = useState<TopicNode[]>([]);
  const [selectedLessonTopic, setSelectedLessonTopic] = useState<TopicNode | null>(null);
  const [selectedLesson, setSelectedLesson] = useState<LifeLesson>(lifeLessons[0]);
  const [activeTestId, setActiveTestId] = useState(tests[0].id);
  const [questionIndex, setQuestionIndex] = useState(0);
  const [answers, setAnswers] = useState<Array<Trait | null>>([]);
  const activeTest = tests.find((test) => test.id === activeTestId) ?? tests[0];
  const extraContent = extraTestContent[activeTest.id];
  const activeTraits = extraContent?.traits ?? ["words", "time", "care", "touch"];
  const activeProfiles = extraContent?.profiles ?? traitProfiles;
  const activeQuestions = extraContent
    ? extraContent.questions.map((question) => ({ ...question, options: question.options.map((option, index) => ({ ...option, trait: activeTraits[index] })) }))
    : loveQuestions;
  const scores = useMemo(() => answers.reduce<Record<Trait, number>>((total, trait) => {
    if (trait) total[trait] = (total[trait] ?? 0) + 1;
    return total;
  }, Object.fromEntries(activeTraits.map((trait) => [trait, 0]))), [answers, activeTestId]);
  const [completedCount, setCompletedCount] = useState(0);
  const [toast, setToast] = useState("");
  const [ageNotice, setAgeNotice] = useState(false);
  const [previewTheme, setPreviewTheme] = useState<"light" | "dark">(() =>
    window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light",
  );

  const firstName = webApp()?.initDataUnsafe?.user?.first_name?.trim() || "دوست من";
  const result = useMemo(() => {
    return (Object.keys(scores) as Trait[]).sort((left, right) => scores[right] - scores[left])[0];
  }, [scores]);
  const progress = screen === "quiz" ? ((questionIndex + 1) / activeQuestions.length) * 100 : 0;

  useEffect(() => {
    document.documentElement.lang = "fa";
    document.documentElement.dir = "rtl";
    document.title = "KnowMe | مینی‌اپ تلگرام";
    const app = webApp();
    app?.ready?.();
    app?.expand?.();
    const applyTheme = () => {
      const scheme = app?.colorScheme ?? previewTheme;
      document.documentElement.dataset.theme = scheme;
      const params = app?.themeParams;
      if (params?.bg_color) document.documentElement.style.setProperty("--km-bg", params.bg_color);
      if (params?.text_color) document.documentElement.style.setProperty("--km-text", params.text_color);
      if (params?.secondary_bg_color) document.documentElement.style.setProperty("--km-surface", params.secondary_bg_color);
    };
    applyTheme();
    app?.onEvent?.("themeChanged", applyTheme);
    const media = window.matchMedia("(prefers-color-scheme: dark)");
    const syncSystemTheme = () => setPreviewTheme(media.matches ? "dark" : "light");
    media.addEventListener("change", syncSystemTheme);
    return () => {
      app?.offEvent?.("themeChanged", applyTheme);
      media.removeEventListener("change", syncSystemTheme);
      document.documentElement.removeAttribute("data-theme");
      document.documentElement.style.removeProperty("--km-bg");
      document.documentElement.style.removeProperty("--km-text");
      document.documentElement.style.removeProperty("--km-surface");
    };
  }, [previewTheme]);

  useEffect(() => {
    const app = webApp();
    const backButton = app?.BackButton;
    if (!backButton) return;
    const goBack = () => {
      if (screen === "quiz") {
        if (questionIndex > 0) setQuestionIndex((current) => current - 1);
        else setScreen("category");
      } else if (screen === "result") setScreen("home");
      else if (screen === "lesson-detail") setScreen("lesson-list");
      else if (screen === "lesson-list") {
        setLessonTopicPath((path) => path.slice(0, -1));
        setScreen("lessons");
      } else if (screen === "lessons" && lessonTopicPath.length > 0) {
        setLessonTopicPath((path) => path.slice(0, -1));
      } else if (screen === "test-items") {
        setTestTopicPath((path) => path.slice(0, -1));
        setScreen("category");
      } else if (screen === "category" && testTopicPath.length > 1) {
        setTestTopicPath((path) => path.slice(0, -1));
      } else if (screen === "category") setScreen("topics");
      else setScreen("home");
    };
    if (screen !== "home") {
      backButton.show();
      backButton.onClick(goBack);
      return () => backButton.offClick(goBack);
    }
    backButton.hide();
  }, [screen, questionIndex]);

  useEffect(() => {
    if (!toast) return;
    const timeout = window.setTimeout(() => setToast(""), 2500);
    return () => window.clearTimeout(timeout);
  }, [toast]);

  function openCategory(code: CategoryCode) {
    webApp()?.HapticFeedback?.selectionChanged();
    setCategory(code);
    if (code === "adult") {
      setAgeNotice(true);
      return;
    }
    const root = testTopicTree.find((topic) => topic.id === code);
    setTestTopicPath(root ? [root] : []);
    setScreen("category");
  }

  function openTestTopic(topic: TopicNode) {
    setTestTopicPath((path) => [...path, topic]);
    if (topic.children?.length) return;
    setScreen("test-items");
  }

  function openLessonTopic(topic: TopicNode) {
    setLessonTopicPath((path) => [...path, topic]);
    if (topic.children?.length) return;
    setSelectedLessonTopic(topic);
    setScreen("lesson-list");
  }

  function startTest(test: TestItem) {
    setActiveTestId(test.id);
    setAnswers([]);
    setQuestionIndex(0);
    setScreen("quiz");
    webApp()?.HapticFeedback?.selectionChanged();
  }

  function chooseAnswer(option: AnswerOption) {
    webApp()?.HapticFeedback?.selectionChanged();
    const updated = [...answers];
    updated[questionIndex] = option.trait;
    setAnswers(updated);
    if (questionIndex + 1 === activeQuestions.length) {
      setCompletedCount((count) => count + 1);
      setScreen("result");
      webApp()?.HapticFeedback?.notificationOccurred("success");
      return;
    }
    setQuestionIndex((current) => current + 1);
  }

  async function shareResult() {
    const profile = activeProfiles[result];
    const text = `KnowMe | ${profile.title}\n${profile.subtitle}\nیک تست کوتاه برای سرگرمی و خودشناسی.`;
    try {
      if (navigator.share) await navigator.share({ title: "نتیجه تست KnowMe", text });
      else if (navigator.clipboard) {
        await navigator.clipboard.writeText(text);
        setToast("متن نتیجه کپی شد؛ حالا می‌تونی در تلگرام بفرستیش.");
      } else setToast("از نتیجه اسکرین‌شات بگیر و در تلگرام بفرست.");
    } catch {
      setToast("برای اشتراک‌گذاری، از نتیجه اسکرین‌شات بگیر.");
    }
  }

  const hideNavigation = screen === "quiz" || screen === "result";

  return (
    <>
      <MobileScroll className="app-screen">
        <main className={`screen-content screen-${screen}`} dir="rtl">
          {screen === "home" && (
            <HomeScreen
              firstName={firstName}
              completedCount={completedCount}
              onStart={() => startTest(tests[0])}
              onTopics={() => setScreen("topics")}
              onProfile={() => setScreen("profile")}
              onRoadmap={() => setScreen("roadmap")}
              onCategory={openCategory}
            />
          )}
          {screen === "topics" && (
            <TopicsScreen onBack={() => setScreen("home")} onCategory={openCategory} />
          )}
          {screen === "category" && (
            <CategoryScreen
              path={testTopicPath}
              onBack={() => testTopicPath.length > 1 ? setTestTopicPath((path) => path.slice(0, -1)) : setScreen("topics")}
              onSelect={openTestTopic}
            />
          )}
          {screen === "test-items" && (
            <TestItemsScreen topic={testTopicPath[testTopicPath.length - 1]} onBack={() => { setTestTopicPath((path) => path.slice(0, -1)); setScreen("category"); }} onStart={startTest} />
          )}
          {screen === "roadmap" && (
            <RoadmapScreen onBack={() => setScreen("home")} onStart={() => startTest(tests[0])} />
          )}
          {screen === "lessons" && (
            <LifeLessonsScreen
              topics={lessonTopicPath.length ? (lessonTopicPath[lessonTopicPath.length - 1].children ?? []) : lessonTopicTree}
              title={lessonTopicPath.length ? lessonTopicPath[lessonTopicPath.length - 1].title : "آموزه‌های زندگی"}
              onBack={() => lessonTopicPath.length ? setLessonTopicPath((path) => path.slice(0, -1)) : setScreen("home")}
              onSelect={openLessonTopic}
            />
          )}
          {screen === "lesson-list" && selectedLessonTopic && (
            <LessonItemsScreen topic={selectedLessonTopic} onBack={() => { setLessonTopicPath((path) => path.slice(0, -1)); setScreen("lessons"); }} onSelect={(lesson) => { setSelectedLesson(lesson); setScreen("lesson-detail"); }} />
          )}
          {screen === "lesson-detail" && (
            <LifeLessonDetailScreen lesson={selectedLesson} onBack={() => setScreen("lesson-list")} />
          )}
          {screen === "profile" && (
            <ProfileScreen firstName={firstName} completedCount={completedCount} onBack={() => setScreen("home")} />
          )}
          {screen === "quiz" && (
            <QuizScreen
              title={activeTest.title}
              question={activeQuestions[questionIndex]}
              questionIndex={questionIndex}
              total={activeQuestions.length}
              progress={progress}
              onBack={() => {
                if (questionIndex > 0) setQuestionIndex((current) => current - 1);
                else setScreen("category");
              }}
              onChoose={chooseAnswer}
            />
          )}
          {screen === "result" && (
            <ResultScreen
              trait={result}
              testTitle={activeTest.title}
              scores={scores}
              profiles={activeProfiles}
              onShare={shareResult}
              onNext={() => setScreen("roadmap")}
              onHome={() => setScreen("home")}
            />
          )}
        </main>
      </MobileScroll>
      {!hideNavigation && (
        <BottomNavigation
          active={screen === "roadmap" ? "roadmap" : screen === "profile" ? "profile" : ["lessons", "lesson-list", "lesson-detail"].includes(screen) ? "lessons" : ["topics", "category", "test-items"].includes(screen) ? "tests" : "home"}
          onHome={() => setScreen("home")}
          onRoadmap={() => setScreen("roadmap")}
          onTests={() => { setTestTopicPath([]); setScreen("topics"); }}
          onLessons={() => { setLessonTopicPath([]); setScreen("lessons"); }}
          onProfile={() => setScreen("profile")}
        />
      )}
      {ageNotice && <AgeNotice onClose={() => setAgeNotice(false)} />}
      {toast && <div className="toast" role="status">{toast}</div>}
      {!webApp() && <div className="theme-preview-switch" role="group" aria-label="پیش‌نمایش تم">
        <button type="button" aria-pressed={previewTheme === "light"} onClick={() => { setPreviewTheme("light"); document.documentElement.dataset.theme = "light"; }}>روشن</button>
        <button type="button" aria-pressed={previewTheme === "dark"} onClick={() => { setPreviewTheme("dark"); document.documentElement.dataset.theme = "dark"; }}>تیره</button>
      </div>}
    </>
  );
}

function BrandHeader({ onProfile }: { onProfile?: () => void }) {
  return (
    <header className="brand-header">
      <div className="brand-lockup" aria-label="KnowMe">
        <span className="brand-mark">Know<span className="brand-accent">Me</span></span>
      </div>
      <button className="avatar-button" type="button" onClick={onProfile}>پروفایل</button>
    </header>
  );
}

function HomeScreen(props: {
  firstName: string;
  completedCount: number;
  onStart: () => void;
  onTopics: () => void;
  onProfile: () => void;
  onRoadmap: () => void;
  onCategory: (code: CategoryCode) => void;
}) {
  return (
    <>
      <section className="hero" aria-labelledby="home-title">
        <img className="hero-image" src={knowmeAsset("night-window.png")} alt="پنجره‌ای رو به شب و چراغ‌های شهر" />
        <div className="hero-shade" />
        <BrandHeader onProfile={props.onProfile} />
        <div className="hero-content">
          <p className="welcome-line">سلام {props.firstName} <span className="wave-dot">•</span></p>
          <h1 id="home-title">یه لایه عمیق‌تر<br /><span>خودتو بشناس</span></h1>
          <p className="hero-copy">با چند سؤال ساده، دریچه‌های تازه‌ای به ذهن و رابطه‌هات باز کن.</p>
          <button className="primary-button" type="button" onClick={props.onStart}>
            <span>شروع تست پیشنهادی</span>
          </button>
          <div className="hero-note"><span className="note-line" /> کوتاه، شخصی و برای خودشناسی</div>
        </div>
        <span className="hero-glow" aria-hidden="true" />
      </section>

      <section className="recommendation-section" aria-labelledby="recommendation-title">
        <div className="section-heading">
          <div><span className="eyebrow">تست پیشنهادی امروز</span><h2 id="recommendation-title">برای شروع مسیرت</h2></div>
          <button className="text-link" type="button" onClick={props.onRoadmap}>نقشه راه</button>
        </div>
        <button className="recommendation" type="button" onClick={props.onStart}>
          <img src={knowmeAsset("love-style.png")} alt="قلب شیشه‌ای زیر گنبد" />
          <span className="recommendation-copy">
            <strong>سبک عشق‌ورزی</strong>
            <span>زبان نزدیک‌شدن تو در رابطه</span>
            <span className="test-meta"><span>۸ سؤال</span><i /><span>حدود ۳ دقیقه</span></span>
          </span>
        </button>
      </section>

      <section className="topic-section" aria-labelledby="topics-title">
        <div className="section-heading compact-heading">
          <div><span className="eyebrow">از هر زاویه‌ای</span><h2 id="topics-title">موضوع‌ها</h2></div>
          <button className="text-link" type="button" onClick={props.onTopics}>همه موضوع‌ها</button>
        </div>
        <div className="topic-grid">
          {categories.map((item) => <TopicButton key={item.code} item={item} onClick={() => props.onCategory(item.code)} />)}
        </div>
      </section>
      {props.completedCount > 0 && <p className="home-progress-note">تا اینجا {props.completedCount} تست کامل کردی؛ مسیرت کم‌کم دقیق‌تر می‌شه.</p>}
    </>
  );
}

function TopicButton({ item, onClick }: { item: (typeof categories)[number]; onClick: () => void }) {
  return (
    <button className={`topic-button topic-${item.code}`} type="button" onClick={onClick}>
      <span className="topic-name">{item.title}</span>
      <span className="topic-caption">{item.caption}</span>
    </button>
  );
}

function PageHeader({ title, onBack, eyebrow }: { title: string; onBack: () => void; eyebrow?: string }) {
  return (
    <header className="page-header">
      <button className="back-button" onClick={onBack} type="button">برگشت</button>
      <div>{eyebrow && <span className="eyebrow">{eyebrow}</span>}<h1>{title}</h1></div>
      <span className="header-spacer" />
    </header>
  );
}

function TopicsScreen({ onBack, onCategory }: { onBack: () => void; onCategory: (code: CategoryCode) => void }) {
  return (
    <>
      <PageHeader title="موضوع‌ها" eyebrow="یک مسیر رو انتخاب کن" onBack={onBack} />
      <p className="page-intro">از کنجکاوی‌ات شروع کن؛ هر تست فقط چند دقیقه زمان می‌بره.</p>
      <div className="topic-list">
        {categories.map((item) => {
          const count = tests.filter((test) => test.category === item.code).length;
          return (
            <button className={`topic-row topic-row-${item.code}`} key={item.code} type="button" onClick={() => onCategory(item.code)}>
              <span className="row-copy"><strong>{item.title}</strong><small>{item.caption}{count > 0 ? ` · ${count} تست` : ""}</small></span>
              <span className="row-action">{item.code === "adult" ? "محدودیت سنی" : "مشاهده تست‌ها"}</span>
            </button>
          );
        })}
      </div>
      <div className="quiet-note"><span>نقشه راه با هر تستی که انجام می‌دی، قدم بعدی رو پیشنهاد می‌کنه.</span></div>
    </>
  );
}

function CategoryScreen({ path, onBack, onSelect }: { path: TopicNode[]; onBack: () => void; onSelect: (topic: TopicNode) => void }) {
  const current = path[path.length - 1];
  const children = current?.children ?? [];
  return (
    <>
      <PageHeader title={current?.title ?? "موضوع‌ها"} eyebrow={path.slice(0, -1).map((node) => node.title).join(" / ") || "موضوع تست‌ها"} onBack={onBack} />
      <p className="page-intro">{current?.caption ?? "یک زیرموضوع رو انتخاب کن."}؛ برای دیدن تست‌ها وارد زیرموضوع شو.</p>
      <div className="topic-list">{children.map((topic) => <TopicRow key={topic.id} topic={topic} onClick={() => onSelect(topic)} />)}</div>
      {!children.length && <div className="empty-state"><h2>این مسیر در حال شکل‌گرفتنه</h2><p>به‌زودی زیرموضوع‌ها و تست‌های تازه اضافه می‌شن.</p></div>}
    </>
  );
}

function TestItemsScreen({ topic, onBack, onStart }: { topic: TopicNode; onBack: () => void; onStart: (test: TestItem) => void }) {
  const topicTests = tests.filter((test) => topic.testIds?.includes(test.id));
  return <>
    <PageHeader title={topic.title} eyebrow="تست‌های این موضوع" onBack={onBack} />
    <p className="page-intro">{topic.caption}</p>
    {topicTests.length ? <div className="test-list">{topicTests.map((test, index) => <button className="test-row" key={test.id} type="button" disabled={test.available === false} onClick={() => onStart(test)}>
      <span className="test-number">{faNumber(String(index + 1).padStart(2, "0"))}</span>
      <span className="row-copy"><strong>{test.title}</strong><small>{test.description}</small>{test.method && <small>روش اجرا: {test.method}</small>}{test.specification && <small className="test-meta">{test.specification}</small>}{test.available === false && <small className="test-meta">سؤال‌های KnowMe در حال آماده‌سازی</small>}</span>
      <span className="coming-soon">{test.available === false ? "به‌زودی" : "شروع تست"}</span>
    </button>)}</div> : <div className="empty-state"><h2>تست‌های این موضوع در راه‌اند</h2><p>این شاخه آماده است و با انتشار تست‌های تازه کامل می‌شود.</p></div>}
  </>;
}

function TopicRow({ topic, onClick }: { topic: TopicNode; onClick: () => void }) {
  const hasChildren = Boolean(topic.children?.length);
  return <button className="topic-row tree-topic-row" type="button" onClick={onClick}>
    <span className="tree-topic-marker">{hasChildren ? "⌄" : "•"}</span>
    <span className="row-copy"><strong>{topic.title}</strong><small>{topic.caption}</small></span>
    <span className="row-action">{hasChildren ? "زیرموضوع‌ها" : "مشاهده آیتم‌ها"} ←</span>
  </button>;
}

function RoadmapScreen({ onBack, onStart }: { onBack: () => void; onStart: () => void }) {
  return (
    <>
      <PageHeader title="نقشه راه من" eyebrow="قدم بعدی تو" onBack={onBack} />
      <div className="roadmap-hero">
        <h2>هر بار، یک کشف تازه</h2>
        <p>با شناخت بیشتر از جواب‌هات، مسیر بعدی متناسب‌تر می‌شه.</p>
      </div>
      <span className="eyebrow roadmap-label">پیشنهاد بعدی برای تو</span>
      <button className="roadmap-test" type="button" onClick={onStart}>
        <img src={knowmeAsset("love-style.png")} alt="" />
        <span><strong>سبک عشق‌ورزی</strong><small>برای شناخت زبان ابراز علاقه‌ات</small><small className="test-meta">۸ سؤال · حدود ۳ دقیقه</small></span>
        <span className="coming-soon">شروع تست</span>
      </button>
      <p className="disclaimer">این مسیر سرگرمی و خودشناسی است و تشخیص روان‌شناختی محسوب نمی‌شود.</p>
    </>
  );
}

function ProfileScreen({ firstName, completedCount, onBack }: { firstName: string; completedCount: number; onBack: () => void }) {
  return (
    <>
      <PageHeader title="پروفایل من" eyebrow="KnowMe" onBack={onBack} />
      <div className="profile-card">
        <div><span className="eyebrow">خوش اومدی</span><h2>{firstName}</h2><p>این مسیر فقط برای خودته.</p></div>
      </div>
      <div className="profile-stat"><span><strong>{completedCount} تست کامل‌شده</strong><small>هر پاسخ، یک تکه از شناخت تو</small></span></div>
      <div className="privacy-note"><p>این پیش‌نمایش اطلاعات را روی همین دستگاه نگه نمی‌دارد. نسخه متصل از شناسه معتبر Telegram برای شناسایی استفاده می‌کند.</p></div>
      <p className="disclaimer">نتیجه‌ها ابزار تشخیص پزشکی یا روان‌شناختی نیستند.</p>
    </>
  );
}

function QuizScreen(props: { title: string; question: Question; questionIndex: number; total: number; progress: number; onBack: () => void; onChoose: (option: AnswerOption) => void }) {
  return (
    <div className="quiz-shell">
      <PageHeader title={props.title} eyebrow="تست کوتاه تو" onBack={props.onBack} />
      <div className="quiz-progress-label"><span>سؤال {faNumber(props.questionIndex + 1)} از {faNumber(props.total)}</span><span>{faNumber(Math.round(props.progress))}٪</span></div>
      <div className="progress-track" aria-label={`پیشرفت ${faNumber(Math.round(props.progress))} درصد`}><span style={{ width: `${props.progress}%` }} /></div>
      <div className="quiz-question"><span className="eyebrow">بدون جواب درست یا غلط</span><h2>{props.question.text}</h2></div>
      <div className="answer-list">
        {props.question.options.map((option, index) => (
          <button className="answer-option" type="button" key={option.text} onClick={() => props.onChoose(option)}>
            <span className="answer-index">{faNumber(String(index + 1).padStart(2, "0"))}</span>
            <span>{option.text}</span>
          </button>
        ))}
      </div>
      <p className="quiz-footnote">جوابت فقط برای ساخت نتیجه همین تست استفاده می‌شه.</p>
    </div>
  );
}

function ResultScreen(props: { trait: Trait; testTitle: string; scores: Record<Trait, number>; profiles: Record<Trait, TraitProfile>; onShare: () => void; onNext: () => void; onHome: () => void }) {
  const profile = props.profiles[props.trait];
  const total = Object.values(props.scores).reduce((sum, value) => sum + value, 0) || 1;
  return (
    <div className="result-shell">
      <div className="result-kicker"><span>یک کشف تازه از تو</span></div>
      <img className="result-image" src={knowmeAsset("love-style.png")} alt="قلب شیشه‌ای در یک گنبد" />
      <section className="result-copy">
        <span className="eyebrow">نتیجهٔ {props.testTitle}</span>
        <h1>{profile.title}</h1>
        <p className="result-subtitle">{profile.subtitle}</p>
        <p className="result-description">{profile.description}</p>
      </section>
      <section className="score-section" aria-label="نقشه سبک‌های تو">
        <h2>نقشه سبک‌های تو</h2>
        {(Object.keys(props.profiles) as Trait[]).map((trait) => {
          const value = Math.round((props.scores[trait] / total) * 100);
          return <div className="score-row" key={trait}><span>{props.profiles[trait].label}</span><div className="score-track"><i style={{ width: `${value}%` }} /></div><span>{value}٪</span></div>;
        })}
      </section>
      <button className="share-button" type="button" onClick={props.onShare}>اشتراک نتیجه</button>
      <button className="primary-button result-next" type="button" onClick={props.onNext}><span>دیدن قدم بعدی</span></button>
      <button className="subtle-button" type="button" onClick={props.onHome}>بازگشت به خانه</button>
      <p className="disclaimer">این نتیجه برای سرگرمی و خودشناسی است؛ تشخیص روان‌شناختی یا پزشکی نیست.</p>
    </div>
  );
}

function LifeLessonsScreen({ topics, title, onBack, onSelect }: { topics: TopicNode[]; title: string; onBack: () => void; onSelect: (topic: TopicNode) => void }) {
  return (
    <>
      <PageHeader title={title} eyebrow="موضوع‌های آموزه‌های زندگی" onBack={onBack} />
      <p className="page-intro">از موضوع کلی شروع کن؛ هر شاخه تو را به راهنمای مرتبط می‌رساند.</p>
      <div className="topic-list">{topics.map((topic) => <TopicRow key={topic.id} topic={topic} onClick={() => onSelect(topic)} />)}</div>
      <div className="quiet-note"><span>این آموزه‌ها جایگزین کمک تخصصی نیستند؛ اگر شرایط سخت یا ناامن است، از یک فرد قابل اعتماد یا متخصص کمک بگیر.</span></div>
    </>
  );
}

function LessonItemsScreen({ topic, onBack, onSelect }: { topic: TopicNode; onBack: () => void; onSelect: (lesson: LifeLesson) => void }) {
  const lessons = lifeLessons.filter((lesson) => topic.lessonIds?.includes(lesson.id));
  return <>
    <PageHeader title={topic.title} eyebrow="آموزه‌های این موضوع" onBack={onBack} />
    <p className="page-intro">{topic.caption}</p>
    <div className="lesson-list">
      {lessons.map((lesson, index) => <button className="lesson-card" key={lesson.id} type="button" onClick={() => onSelect(lesson)}>
        <span className="lesson-index">{faNumber(String(index + 1).padStart(2, "0"))}</span>
        <span className="lesson-card-copy"><strong>{lesson.title}</strong><small>{lesson.summary}</small><em>خواندن راهنما ←</em></span>
      </button>)}
    </div>
  </>;
}

function LifeLessonDetailScreen({ lesson, onBack }: { lesson: LifeLesson; onBack: () => void }) {
  return (
    <>
      <PageHeader title="آموزه‌ی زندگی" eyebrow="یک قدم آرام‌تر" onBack={onBack} />
      <section className="lesson-detail-hero"><span>🌱 راهنمای زندگی</span><h2>{lesson.title}</h2><p>{lesson.summary}</p></section>
      <div className="lesson-sections">
        <section className="lesson-section lesson-goal"><span>🎯 هدف این آموزه</span><p>{lesson.goal}</p></section>
        <section className="lesson-section"><span>🪞 موضوع چیست؟</span><p>{lesson.explanation}</p></section>
        <section className="lesson-section lesson-consequences"><span>⚠️ اگر نادیده‌اش بگیریم</span><p>{lesson.consequences}</p></section>
        <section className="lesson-section"><span>👥 رده‌ی سنی</span><p>{lesson.ageRange}</p></section>
        <section className="lesson-section lesson-method"><span>🧭 روش‌های عملی</span><ol>{lesson.method.map((step) => <li key={step}>{step}</li>)}</ol></section>
        <section className="lesson-section lesson-help"><span>🤝 چه زمانی کمک بگیرم؟</span><p>{lesson.help}</p></section>
        <section className="lesson-takeaway"><span>💚 یادت بماند</span><p>{lesson.takeaway}</p></section>
      </div>
    </>
  );
}

function BottomNavigation(props: { active: "home" | "roadmap" | "tests" | "lessons" | "profile"; onHome: () => void; onRoadmap: () => void; onTests: () => void; onLessons: () => void; onProfile: () => void }) {
  const items = [
    { id: "home", label: "خانه", onClick: props.onHome },
    { id: "roadmap", label: "نقشه راه", onClick: props.onRoadmap },
    { id: "tests", label: "تست‌ها", onClick: props.onTests },
    { id: "lessons", label: "آموزه‌ها", onClick: props.onLessons },
    { id: "profile", label: "پروفایل", onClick: props.onProfile },
  ] as const;
  return <nav className="bottom-nav" aria-label="ناوبری اصلی">{items.map(({ id, label, onClick }) => <button key={id} type="button" className={`nav-item ${props.active === id ? "is-active" : ""}`} aria-current={props.active === id ? "page" : undefined} onClick={onClick}><span>{label}</span></button>)}</nav>;
}

function AgeNotice({ onClose }: { onClose: () => void }) {
  return (
    <div className="notice-overlay" role="presentation" onClick={onClose}>
      <section className="notice-sheet" role="dialog" aria-modal="true" aria-labelledby="age-title" onClick={(event) => event.stopPropagation()}>
        <span className="eyebrow">محتوای محدود</span>
        <h2 id="age-title">بخش ۱۸+ با تأیید سن باز می‌شود</h2>
        <p>برای حفظ محدودیت سنی، دسترسی باید با سن ثبت‌شده در پروفایل KnowMe تأیید شود.</p>
        <button className="primary-button" type="button" onClick={onClose}>متوجه شدم</button>
      </section>
    </div>
  );
}
