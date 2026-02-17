# Installation Guide

This guide will help you set up all the prerequisites for working with this Bazel tutorial.

## Prerequisites Overview

- **Bazel** (6.4.0 or later)
- **Git** for cloning the repository
- **Python 3.11+** (optional, managed by Bazel)
- **Go 1.21+** (optional, managed by Bazel)
- **Node.js 18+** (for Economic Dashboard UI)

## Installing Bazel

### Option 1: Using Bazelisk (Recommended)

Bazelisk automatically downloads and uses the correct Bazel version specified in `.bazelversion`.

#### macOS
```bash
brew install bazelisk
```

#### Linux
```bash
# Using npm
npm install -g @bazel/bazelisk

# Or download binary
wget https://github.com/bazelbuild/bazelisk/releases/latest/download/bazelisk-linux-amd64 -O /usr/local/bin/bazel
chmod +x /usr/local/bin/bazel
```

#### Windows
```bash
choco install bazelisk
```

### Option 2: Direct Bazel Installation

Follow the official guide for your platform: [https://bazel.build/install](https://bazel.build/install)

#### macOS
```bash
brew install bazel
```

#### Ubuntu/Debian
```bash
sudo apt update && sudo apt install apt-transport-https curl gnupg
curl -fsSL https://bazel.build/bazel-release.pub.gpg | gpg --dearmor >bazel-archive-keyring.gpg
sudo mv bazel-archive-keyring.gpg /usr/share/keyrings
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/bazel-archive-keyring.gpg] https://storage.googleapis.com/bazel-apt stable jdk1.8" | sudo tee /etc/apt/sources.list.d/bazel.list
sudo apt update && sudo apt install bazel
```

### Verify Installation

```bash
bazel --version
# Should show: bazel 6.4.0 or later
```

## Installing Node.js (For Economic Dashboard)

The Economic Dashboard web UI requires Node.js 18 or later.

### macOS
```bash
brew install node
```

### Linux
```bash
# Using NodeSource
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs
```

### Windows
Download from [nodejs.org](https://nodejs.org)

### Verify Installation

```bash
node --version  # Should be v18.0.0 or higher
npm --version
```

## Getting a FRED API Key (For Economic Dashboard)

The Economic Dashboard fetches data from the Federal Reserve Economic Data (FRED) API.

1. Visit [https://fred.stlouisfed.org/](https://fred.stlouisfed.org/)
2. Create a free account (takes 1 minute)
3. Go to **My Account → API Keys**
4. Click **Request API Key**
5. Copy your API key

Set it in your environment:

```bash
export FRED_API_KEY=your_actual_api_key_here

# Add to your shell profile for persistence
echo 'export FRED_API_KEY=your_actual_api_key_here' >> ~/.bashrc  # or ~/.zshrc
```

## Clone the Repository

```bash
git clone https://github.com/amulya-labs/bazel-tutorial.git
cd bazel-tutorial
```

## Verify Your Setup

Run the validation script to check everything is working:

```bash
./validate.sh
```

This will:
- Verify Bazel is installed
- Check repository structure
- Build all targets
- Run all tests

## Common Installation Issues

### "command not found: bazel"

**Problem:** Bazel is not installed or not in your PATH.

**Solution:**
1. Make sure you followed the installation steps above
2. Check if Bazel is in your PATH: `which bazel`
3. Try using Bazelisk instead of direct Bazel installation

### "Permission denied" when installing

**Problem:** Insufficient permissions to install packages.

**Solution:**
- Use `sudo` for system-wide installations (Linux/macOS)
- Run as Administrator (Windows)
- Use user-level package managers (Homebrew, npm global packages)

### "Incompatible Bazel version"

**Problem:** Your Bazel version doesn't match `.bazelversion`.

**Solution:**
- If using Bazelisk: It should auto-download the correct version
- If using direct Bazel: Install the version specified in `.bazelversion`
- Check with: `cat .bazelversion`

### Node.js version too old

**Problem:** Economic Dashboard requires Node.js 18+.

**Solution:**
```bash
# Check your version
node --version

# Upgrade using your package manager
brew upgrade node  # macOS
# or use nvm
nvm install 18
nvm use 18
```

## Next Steps

Once installation is complete:

1. **[Quick Start Guide](quickstart.md)** - Run your first builds
2. **[Basic Tutorial](basic-tutorial.md)** - Learn with simple services
3. **[Economic Dashboard](economic-dashboard.md)** - Build the full-stack app

## Additional Tools (Optional)

### Buildifier

Format Bazel files:
```bash
go install github.com/bazelbuild/buildtools/buildifier@latest
```

### VS Code Extensions

- **Bazel** - Syntax highlighting for BUILD files
- **Python** - Python language support
- **Go** - Go language support

### Remote Caching (Advanced)

For team collaboration, consider setting up remote caching:
- [Bazel Remote Cache Documentation](https://bazel.build/remote/caching)

---

Need help? [Open an issue](https://github.com/amulya-labs/bazel-tutorial/issues)
