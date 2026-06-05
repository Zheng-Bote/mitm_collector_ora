# Oracle Table Data Collector

The **Oracle Table Data Collector** is an autonomous Go program designed to run as a scheduled job. It dynamically retrieves all records from a specified database table inside an Oracle database instance using the `"github.com/sijms/go-ora/v2"` pure-Go driver, encrypts the records using AES-GCM Envelope Encryption (with the database storage DEK), and writes the encrypted payloads to the central MitM database's `raw_ingestion` table.

For code details, refer to:

- [main.go](file:///home/zb_bamboo/DEV/__NEW__/Go/mitm-2/collector-layer/mitm_collector_ora-employee/main.go) - Dynamic row reading, encryption, and ingestion logic.
- [go.mod](file:///home/zb_bamboo/DEV/__NEW__/Go/mitm-2/collector-layer/mitm_collector_ora-employee/go.mod) - Dependency definition.

---

## 🏗️ How It Works

1.  **Bootstrapping**: Expects the MitM database connection configuration passed as a JSON string argument (`os.Args[1]`), an optional JSON arguments string (`os.Args[2]`) injected by the scheduler to override settings (like table name, source name, cursor column, and topic), and environment settings.
2.  **Envelope Decryption**:
    - Reads the Key Encryption Key (KEK) from the `MASTER_KEY` environment variable.
    - Retrieves the encrypted Oracle source DB config and wrapped Data Encryption Key (DEK) from the MitM PostgreSQL database.
    - Decrypts the DEK using the KEK, then decrypts the Oracle connection credentials using the DEK.
3.  **Dynamic Extraction**:
    - Connects to the source Oracle database using `go-ora`.
    - Loads the last processed cursor offset from `ingestion_cursors`.
    - Queries new records from the specified table using a configurable cursor column (`cursor_column > lastCursor`).
    - **Dynamic Scanning**: Scans columns dynamically without knowing the schema at compile time, resolving columns to a map of strings/values.
4.  **Ingestion**:
    - Serializes each database row map into a JSON string.
    - Encrypts the JSON payload via AES-GCM using the DEK and a fresh random nonce.
    - Inserts the encrypted records into the target `raw_ingestion` table with a status of `pending`.
    - Updates the cursor offset to the highest processed ID.
5.  **IPC Event Reporting**: Reports events (`started`, `processing`, `finished`, `failed`, and `audit`) to the scheduler via Unix Domain Socket.

---

## ⚙️ Configuration & Environment

### Environment Variables

- `MASTER_KEY` (Required): The base64-encoded 32-byte Master Key (KEK) used to unwrap DEKs.
- `RUN_ID` (Optional): Run ID injected by the scheduler to identify this execution.
- `SCHEDULER_SOCKET_PATH` (Optional): Path to the Unix socket for IPC event logging.

### JSON CLI Arguments

The collector accepts up to two JSON parameters as command-line arguments:

#### 1. MitM Target DB Connection (`os.Args[1]`)

A JSON string detailing the connection parameters to the central MitM target database.

Example:

```json
{
  "host": "localhost",
  "port": 5432,
  "user": "mitm_user",
  "password": "mitm_password",
  "database": "mitm",
  "source_name": "ORA_EMPLOYEE"
}
```

#### 2. Optional Job Overrides (`os.Args[2]`)

An optional JSON string passed by the scheduler to override the default ingestion behaviour.

Example:

```json
{
  "source_name": "ORA_EMPLOYEE",
  "table": "EMPLOYEES",
  "cursor_column": "ID",
  "topic": "employee.data"
}
```

---

## 🛠️ Build Instructions

### Prerequisites

- Go 1.25.0 or later installed.

### Compiling the Binary

To compile the collector into a standalone executable, navigate to the collector directory and build:

```bash
cd /home/zb_bamboo/DEV/__NEW__/Go/mitm-2/collector-layer/mitm_collector_ora-employee
go build -o bin/mitm-collector-ora-employee main.go
```

This compiles a static executable `mitm-collector-ora-employee` inside the local `bin/` directory.

---

## 🚀 Execution Example

To test the binary manually from the command line:

```bash
# 1. Export the Master Key (must match the one used during DB initialization)
export MASTER_KEY="Y29uZmlkZW50aWFsX21hc3Rlcl9rZXlfMzJfYnl0ZXM="

# 2. Run the collector binary, passing the MitM connection details and optional overrides
./bin/mitm-collector-ora-employee '{
  "host": "127.0.0.1",
  "port": 1521,
  "user": "orauser",
  "password": "yourpassword",
  "database": "hr",
  "source_name": "ORA_EMPLOYEE"
}' '{"source_name": "ORA_EMPLOYEE", "table": "EMPLOYEES", "cursor_column": "ID", "topic": "employee.data"}'
```
