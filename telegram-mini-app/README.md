# KnowMe Telegram Mini App

A mobile-first, Persian RTL Mini App prototype for the existing KnowMe quiz product. The UI uses the selected midnight-plum and mint direction and the Vazirmatn typeface.

When opened inside Telegram, the app follows Telegram's light/dark `colorScheme`, listens for `themeChanged`, and uses the host `bg_color`, `text_color`, and `secondary_bg_color` where provided. In a regular browser it follows the system color preference; the floating preview switch can override the theme for visual review.

## Run locally

```bash
npm ci
npm run dev -- --host 0.0.0.0 --port 4173
```

## Prototype flow

- Browse the current KnowMe topics and quiz catalog.
- Open the roadmap or profile screens.
- Complete the eight seeded «سبک عشق‌ورزی» questions and view the seeded result profiles and score breakdown.
- Open the adult-content notice for the 18+ topic.
- Use Telegram's Web App ready, expand, back-button, and haptic APIs when running inside Telegram; browser preview uses local fallbacks.

The other catalog entries are shown from the current seed data but aren't executable in this standalone preview. The prototype doesn't store quiz sessions or write profiles to the Go service. Telegram's `initDataUnsafe` is used only for the greeting; production identity and age checks must be validated on the server from signed `initData`.

## Production connection still needed

To launch this from the bot, configure a public HTTPS Mini App URL and add a Telegram Web App entry point. Then connect catalog, profile, session, answer, result, and age-gate reads/writes to server endpoints that validate Telegram `initData`. Never trust client-provided Telegram IDs or age for authorization.
