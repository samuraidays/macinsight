# macinsight

macOS Security Audit CLI tool that checks various security settings on macOS systems.

## Features

- **SIP (System Integrity Protection)**: Check if SIP is enabled
- **Gatekeeper**: Verify Gatekeeper status and settings
- **FileVault**: Check FileVault encryption status
- **Firewall**: Verify firewall configuration

## Installation

```bash
go install github.com/samuraidays/macinsight/cmd/macinsight@latest
```

## Usage

### Basic audit
```bash
macinsight audit
```

### JSON output
```bash
macinsight audit --json
```

### Run specific checks
```bash
macinsight audit --only filevault,gatekeeper
```

### Exclude specific checks
```bash
macinsight audit --exclude sip,firewall
```

### Set timeout
```bash
macinsight audit --timeout 5s
```

### List available checks
```bash
macinsight list-checks
```

### Show version
```bash
macinsight version
```

## Available Checks

- `sip`: System Integrity Protection status
- `gatekeeper`: Gatekeeper configuration
- `filevault`: FileVault encryption status
- `firewall`: Firewall settings

## Output Formats

- **Table format** (default): Human-readable table output
- **JSON format**: Machine-readable JSON output with `--json` flag

## Requirements

- macOS 10.12 or later
- Go 1.19 or later (for building from source)

## Building from Source

```bash
git clone https://github.com/samuraidays/macinsight.git
cd macinsight
go build -o macinsight cmd/macinsight/main.go
```

## License

MIT License - see [LICENSE](LICENSE) file for details.
