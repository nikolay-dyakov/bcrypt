# bcrypt

A small command-line tool for generating bcrypt password hashes.

## Requirements

- No local Go installation is required when using a prebuilt release.
- Building or installing with `go install` requires Go 1.27 or newer.

## Install a prebuilt release

Download the archive for your operating system and processor from the
[latest GitHub release](https://github.com/nikolay-dyakov/bcrypt/releases/latest):

| Operating system | Intel/AMD 64-bit | ARM 64-bit |
| --- | --- | --- |
| Linux | `linux_amd64.tar.gz` | `linux_arm64.tar.gz` |
| macOS | `darwin_amd64.tar.gz` | `darwin_arm64.tar.gz` |
| Windows | `windows_amd64.zip` | `windows_arm64.zip` |

Extract the archive and place `bcrypt` (or `bcrypt.exe` on Windows) in a
directory included in your `PATH`. On Linux and macOS, for example:

```bash
tar -xzf bcrypt_VERSION_OS_ARCH.tar.gz
sudo install -m 0755 bcrypt /usr/local/bin/bcrypt
bcrypt version
```

Each release includes `checksums.txt` so downloads can be verified before
installation.

## Install from GitHub

Install the latest version directly from GitHub:

```bash
go install github.com/nikolay-dyakov/bcrypt@latest
```

Go installs the executable in `GOBIN` when it is configured, or in
`$(go env GOPATH)/bin` otherwise. Make sure that directory is included in your
`PATH`. For example, Bash and Zsh users can add this to their shell profile:

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
```

Open a new terminal and verify the installation:

```bash
bcrypt version
```

Run the same `go install` command again whenever you want to upgrade to the
latest version.

## Build from source

Clone the repository and build it locally:

```bash
git clone https://github.com/nikolay-dyakov/bcrypt.git
cd bcrypt
go build -o bcrypt .
```

## Usage

Run without `-p` to enter the password using a hidden terminal prompt:

```bash
bcrypt
```

You can also pipe a password into the tool:

```bash
printf '%s' 'correct horse battery staple' | bcrypt
```

Copy the generated hash directly to the clipboard:

```bash
bcrypt -clipboard
```

Display the installed version using any of these forms:

```bash
bcrypt version
bcrypt -version
bcrypt --version
```

The default bcrypt cost is 13. Override it when needed:

```bash
bcrypt -cost 12
```

Costs below 10 are rejected. The best cost depends on the hardware and workload of
the system that will verify the hashes; choose the highest cost that keeps password
verification within the application's latency budget. Bcrypt accepts passwords of
at most 72 bytes.

The `-p`/`-password` option is retained for scripting, but it can expose the password
in shell history and process listings. Prefer the hidden prompt or stdin.

```text
Usage of bcrypt:
  -c, -cost int
        bcrypt cost (10-31) (default 13)
  -cc, -clipboard
        copy the bcrypt hash to the clipboard
  -p, -password string
        password to hash (prefer the hidden prompt or stdin)
  -version
        print version information
```
