# Architecture

## Runtime

KnowMe v0 is a small Go process designed to run as a Windows service or as a Linux/ARM64 binary.

- Telegram transport: Long Polling
- Database: remote SQL Server
- No public inbound port is required
- Secrets are provided only through environment variables

## Core modules

- **Telegram transport**: receives updates and sends messages.
- **Test Engine**: owns sessions, questions, options and answers.
- **Scoring Engine**: converts answers to result dimensions/traits.
- **Roadmap Engine**: recommends the next test based on previous results and completion history.
- **Profile Engine**: progressive profiling; personal fields are requested gradually, not in one onboarding form.
- **Event Tracking**: stores meaningful bot interactions such as test start, answer, abandonment and completion.
- **Share Card Generator**: future module for result images suitable for Telegram/Instagram sharing.

## Data principle

Every answer is persisted immediately. Aggregate traits are derived data and can be recalculated from raw answers.

## Roadmap v0

The first implementation recommends the first active test that the user has not completed.

Later versions should score candidate tests using:

1. completed tests,
2. trait confidence,
3. previous answers,
4. abandoned tests,
5. age/gender applicability,
6. category diversity,
7. exploration versus personalization.

The roadmap must remain deterministic enough to explain why a test was recommended.
