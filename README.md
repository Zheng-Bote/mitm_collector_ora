# Oracle Table Data Collector

The **Oracle Table Data Collector** is an autonomous Go program designed to run as a scheduled job. It dynamically retrieves all records from a specified database table inside an Oracle database instance using the `"github.com/sijms/go-ora/v2"` pure-Go driver, encrypts the records using AES-GCM Envelope Encryption (with the database storage DEK), and writes the encrypted payloads to the central MitM database's `raw_ingestion` table.

For code details, refer to:

- [main.go](main.go) - Dynamic row reading, encryption, and ingestion logic.
- [go.mod](go.mod) - Dependency definition.

---

## 🏗️ How It Works

1.  **Bootstrapping**: Expects the MitM database connection configuration passed via `MITM_DB_*` environment variables, an optional JSON arguments string (`os.Args[1]`) injected by the scheduler to override settings (like table name, source name, cursor column, and topic), and other environment settings.
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
    - Generates a deterministic `correlation_id` (UUIDv5) based on the specified `business_key_column` or fallback column.
    - Inserts the encrypted records into the target `raw_ingestion` table with a status of `pending`.
    - Updates the cursor offset to the highest processed ID.
5.  **IPC Event Reporting**: Reports events (`started`, `processing`, `finished`, `failed`, and `audit`) to the scheduler via Unix Domain Socket.

---

## ⚙️ Configuration & Environment

### Environment Variables

- `MASTER_KEY` (Required): The base64-encoded 32-byte Master Key (KEK) used to unwrap DEKs.
- `MITM_DB_CONFIG_JSON` (**Preferred**): JSON-encoded credentials containing a nested `"db"` object for the MitM PostgreSQL database.
- `MITM_DB_HOST`, `MITM_DB_PORT`, `MITM_DB_USER`, `MITM_DB_PASSWORD`, `MITM_DB_NAME` (**Fallback**): The connection parameters for the central target MitM database.
- `MITM_DB_SSLMODE` (Optional): If set to `"true"`, enforces SSL connections (`sslmode=require`) to the MitM target database.
- `RUN_ID` (Optional): Run ID injected by the scheduler to identify this execution.
- `SCHEDULER_SOCKET_PATH` (Optional): Path to the Unix socket for IPC event logging.

### JSON CLI Arguments

The collector accepts an optional JSON parameter as command-line argument:

#### 1. Optional Job Overrides (`os.Args[1]`)

An optional JSON string passed by the scheduler to override the default ingestion behaviour.

Example:

```json
{
  "source_name": "ORA_EMPLOYEE",
  "table": "EMPLOYEES",
  "cursor_column": "ID",
  "topic": "employee.data",
  "business_key_column": "EMPLOYEE_ID"
}
```

- `business_key_column`: (Optional) Specifies the column to be used for generating the deterministic `correlation_id`. This is critical for **Stateful Aggregation**, allowing the Transformation Layer to join data from multiple sources. If omitted, it defaults to the `cursor_column` or "UNKNOWN".

---

## 🛠️ Build Instructions

### Prerequisites

- Go 1.25.0 or later installed.

### Compiling the Binary

To compile the collector into a standalone executable, navigate to the collector directory and build:

```bash
cd /home/zb_bamboo/DEV/__NEW__/Go/mitm-2/collector-layer/mitm_collector_ora
go build -o bin/mitm-collector-ora main.go
```

This compiles a static executable `mitm-collector-ora` inside the local `bin/` directory.

---

## 🚀 Execution Example

To test the binary manually from the command line:

```bash
# 1. Export the Master Key (must match the one used during DB initialization)
export MASTER_KEY="Y29uZmlkZW50aWFsX21hc3Rlcl9rZXlfMzJfYnl0ZXM="

# 2. Provide MitM Database configuration (Preferred JSON format)
export MITM_DB_CONFIG_JSON='{"db":{"host":"127.0.0.1","port":5432,"user":"mitm_user","password":"mitm_password","database":"mitm_db"}}'

# 3. Run the collector binary, passing the optional JSON arguments (Overrides)
# This example reads from the 'EMPLOYEES' table, uses 'EMPLOYEE_ID' as correlation key,
# and tracks the 'ID' column for incremental extraction.
./bin/mitm-collector-ora '{
  "source_name": "ORA_EMPLOYEE",
  "table": "EMPLOYEES",
  "cursor_column": "ID",
  "topic": "employee.data",
  "business_key_column": "EMPLOYEE_ID"
}'
```

### Alternative: Direct Environment Variables (Fallback)

If you prefer not to use `MITM_DB_CONFIG_JSON`, you can configure the MitM connection explicitly:

```bash
export MITM_DB_HOST="127.0.0.1"
export MITM_DB_PORT="5432"
export MITM_DB_USER="mitm_user"
export MITM_DB_PASSWORD="mitm_password"
export MITM_DB_NAME="mitm_db"
export MITM_DB_SSLMODE="true"

./bin/mitm-collector-ora '{"source_name": "ORA_EMPLOYEE", "table": "EMPLOYEES", "cursor_column": "ID", "topic": "employee.data", "business_key_column": "EMPLOYEE_ID"}'
```
