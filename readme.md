# Slurp

![Slurp Diagram](sipload.png)

A token-bucket-based SIP load-testing tool written in Go.  
This project provides a CLI for generating SIP-like traffic at a controlled rate and concurrency level, optionally simulating SIP calls and registration sequences.

---

## Features

- **Token Bucket Rate Limiting**  
  Control calls-per-second (CPS) to throttle the load on your SIP infrastructure.

- **Concurrency Management**  
  Limit the number of simultaneous calls (sessions) to emulate real-world usage.

- **Mock SIP Logic**  
  Sends mock SIP calls (simulated) with randomized outcomes (for demonstration). Integrate a real SIP library to test production scenarios.

- **Optional SIP REGISTER**  
  Send a SIP REGISTER message before placing calls if your environment requires registration/authentication.

- **YAML Configuration**  
  Use the central file `configs/config.yaml` for defaults, and override with CLI flags or environment variables.

- **Zap Logging**  
  Structured logging with various log levels (info, error, debug, etc.).

- **Metrics Collection**  
  Basic stats on total calls, failures, and elapsed time, displayed at the end of each test.

- **Extensible**  
  Designed with modular packages (`internal/load`, `internal/sip`, `internal/stats`, etc.) for easy customization.

---

## Table of Contents

1. [Prerequisites](#prerequisites)  
2. [Installation](#installation)  
3. [Configuration](#configuration)  
4. [Usage](#usage)  
5. [Commands](#commands)  
   - [test](#test-command)  
   - [version](#version-command)  
6. [Examples](#examples)  
7. [Testing the Code](#testing-the-code)  
8. [Roadmap](#roadmap)  
9. [License](#license)  

---

## Prerequisites

- **Go 1.20+** (older versions may work, but 1.20 or newer is recommended)
- **Git** (optional, if you need to clone the repo)

## Installation

1. **Clone the repository** (or download the source code):
   ```bash
   git clone https://github.com/yourusername/Slurp.git
   cd Slurp
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Build the binary**:
   ```bash
   go build -o Slurp .
   ```
   This produces an executable named **`Slurp`**.

---

## Configuration

- The default configuration file is located at `./configs/config.yaml`.  
- Example contents:
  ```yaml
  # configs/config.yaml

  target_uri: "sip:echo@sip.testserver.com"
  calls_per_second: 5
  concurrency: 2
  duration: 10
  local_contact: "sip:mytestclient@127.0.0.1:5060"
  register_first: false
  ```
- **Environment Variables**:  
  You can override settings by using environment variables with the prefix `Slurp_`. For example:
  ```bash
  export Slurp_TARGET_URI="sip:echo@other.testserver.com"
  export Slurp_CALLS_PER_SECOND=20
  ```
- **CLI Flags**:  
  Values from the config file can be overridden by flags like `--target`, `--calls-per-second`, etc.

---

## Usage

After building, you can run the CLI via:

```bash
./Slurp [command] [flags...]
```

---

## Commands

### Test Command

```bash
./Slurp test [flags...]
```

- **Description**: Run a SIP load test with token-bucket rate limiting and optional registration.  
- **Flags**:  
  - `--target <uri>`: SIP target URI (overrides `target_uri` in config)  
  - `--calls-per-second <n>`: Desired call generation rate (CPS)  
  - `--concurrency <n>`: Max number of simultaneous calls  
  - `--duration <time>`: Test duration (e.g., `10s`, `30s`, or `0` for infinite until Ctrl+C)  
  - `--contact <uri>`: Local SIP contact  
  - `--register-first`: If set, send a mock REGISTER first  

### Version Command

```bash
./Slurp version
```

- **Description**: Shows the current version of **Slurp**.

---

## Examples

1. **Use defaults from config.yaml**:
   ```bash
   ./Slurp test
   ```
2. **Override with flags**:
   ```bash
   ./Slurp test \
     --target "sip:echo@otherserver.com" \
     --calls-per-second 10 \
     --concurrency 5 \
     --duration 15s \
     --register-first
   ```
3. **Check version**:
   ```bash
   ./Slurp version
   ```

---

## Testing the Code

- **Unit tests** live in `_test.go` files under each package (for example, `internal/rng/rng_test.go`).  
- Run them with:
  ```bash
  go test ./...
  ```
- The tests include basic coverage for rate limiting, concurrency, stats collection, and mock SIP calls.  
- Note: The mock SIP logic might fail randomly. In a real environment, you’d integrate a true SIP library or mock out the random failure for deterministic tests.

---

## Roadmap

- **Integration with Real SIP Library**: Replace the mock calls in `internal/sip` with an actual SIP stack.  
- **Distributed/Clustered Testing**: Scale across multiple machines or containers.  
- **Prometheus Metrics**: Export real-time metrics for visualizations (e.g., Grafana).  
- **Scenario Scripting**: Support advanced call flows beyond a basic `INVITE` or `REGISTER`.
