# Changelog

All notable changes to the Oracle Collector will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.3.0] - 2026-06-06

### Changed
- Changed MitM database credentials initialization: Credentials are now read from `MITM_DB_*` environment variables instead of `os.Args[1]`.
- Job argument configuration (`CollectorArgs`) is now read from `os.Args[1]` instead of `os.Args[2]`.

## [v0.2.0] - 2026-06-05

### Added
- Fully table-independent and dynamic query engine using standard SQL `rows.Columns()` metadata.
- Dynamic row scanning into generic maps using column pointers to support arbitrary Oracle schemas.
- Support for runtime configuration overrides (source name, table, cursor column, and destination topic) passed as a JSON string via `os.Args[2]`.

### Changed
- Replaced hardcoded `Employee` struct scan logic with generic map serialization.
- Updated database insertion query to route records to dynamic topics (defaults to `oracle.<table_name>.data`).
- Updated cursor persistence to support generic string-based cursor values (`maxCursorValue`) instead of numeric IDs.

## [v0.1.0] - 2026-06-04

### Added
- Initial release of the Oracle Employee Collector.
- Automated extraction of raw employee data from Oracle using envelope encryption (AES-256-GCM).
- State-based pagination using the `ingestion_cursors` table.
- Nil-safe IPC reporting for status, audit, and progress logging to the scheduler.
