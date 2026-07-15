# Changelog

All notable changes to the Oracle Collector will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.10.0] - 2026-07-07

### Added
- **SSL Support**: Added support for the `MITM_DB_SSLMODE` environment variable. The collector now respects this setting and applies it to the MitM PostgreSQL connection string.
- **Robust Audit Logging**: Connection errors to both the MitM database and the Oracle source database are now explicitly logged to the central Audit Log via `ipc.SendAudit()`.

### Fixed
- **Oracle Connection DSN**: Switched to using `go_ora.BuildUrl()` instead of `net/url` to generate the connection string. This fixes a bug where URL-encoding of special characters in the password (e.g., `+` to `%2B`) caused authentication failures (Ping-Timeout) because the `go-ora` driver does not URL-decode credentials.

## [v0.9.0] - 2026-06-30

### Changed
- **Config Restructuring**: Updated the MitM database connection logic to correctly parse the JSON configuration (`MITM_DB_CONFIG_JSON`) provided by the scheduler, successfully unpacking the nested `"db"` object format.
- **Database Connection**: The collector now strictly prioritizes the JSON configuration (`MITM_DB_CONFIG_JSON`) over direct environment variables. Direct environment variables (`MITM_DB_HOST`, etc.) now serve only as a fallback.
- **Audit Logging**: Added IPC audit logging (`ipc.SendAudit`) during initialization to accurately log the source of the database configuration (`JSON Config (MITM_DB_CONFIG_JSON)` vs `Environment Variables`).

## [v0.8.0] - 2026-06-24

### Added
- **extending logging**


## [v0.7.0] - 2026-06-21

### Added
- **Stateful Aggregation**: Replaced `raw_ingestion_id` with a deterministic `correlation_id` (UUIDv5).
- **Business Keys**: Introduced a `business_key_column` configuration property. The collector dynamically hashes this column's value (or a fallback) to compute stable correlation IDs, allowing cross-system fragment joins in the Transformation Layer.

## [v0.6.0] - 2026-06-15

### Added
- **Centralized App Info**: Added `appName` and `version` globally. The component now broadcasts its name and version via IPC when starting.

## [v0.5.0] - 2026-06-10

### Added
- **Component Identifier**: Upgraded `IPCClient` to include the `Component` identifier `"mitm_collector_ora"` in all audit and status events sent to the scheduler.

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
