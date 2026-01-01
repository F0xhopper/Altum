# Altum

A CLI tool for managing your daily notes.

## Installation (macOS)

### For Apple Silicon (arm64)

```sh
# Download the latest release
curl -LO https://github.com/F0xhopper/Altum/releases/latest/download/Altum_Darwin_arm64.tar.gz
# Extract
tar -xzf Altum_Darwin_arm64.tar.gz

# Make executable and install globally
chmod +x Altum
sudo mv Altum /usr/local/bin/altum

# Verify installation
altum --version

# Cleanup
rm Altum_Darwin_arm64.tar.gz
```

**One-liner installation (arm64):**

```sh
curl -LO https://github.com/F0xhopper/Altum/releases/latest/download/Altum_Darwin_arm64.tar.gz && \
tar -xzf Altum_Darwin_arm64.tar.gz && \
chmod +x Altum && \
sudo mv Altum /usr/local/bin/altum && \
rm Altum_Darwin_arm64.tar.gz && \
altum --version
```

### For Intel (x86_64)

```sh
# Download the latest release
curl -LO https://github.com/F0xhopper/Altum/releases/latest/download/Altum_Darwin_x86_64.tar.gz
# Extract
tar -xzf Altum_Darwin_x86_64.tar.gz

# Make executable and install globally
chmod +x Altum
sudo mv Altum /usr/local/bin/altum

# Verify installation
altum --version

# Cleanup
rm Altum_Darwin_x86_64.tar.gz
```

**Alternative: Install to user directory (no sudo required):**

```sh
# For arm64
curl -LO https://github.com/F0xhopper/Altum/releases/latest/download/Altum_Darwin_arm64.tar.gz
tar -xzf Altum_Darwin_arm64.tar.gz
chmod +x Altum
mkdir -p ~/bin
mv Altum ~/bin/altum

# Add ~/bin to PATH if not already there (add to ~/.zshrc or ~/.bash_profile)
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

altum --version
```

### Verifying Checksums

You can verify the integrity of the downloaded files using the checksums file:

```sh
# Download the checksums file
curl -LO https://github.com/F0xhopper/Altum/releases/latest/download/Altum_1.0.6_checksums.txt

# Verify the downloaded archive (replace with the appropriate file name)
sha256sum -c Altum_1.0.6_checksums.txt --ignore-missing Altum_Darwin_arm64.tar.gz
# or on macOS
shasum -a 256 -c Altum_1.0.6_checksums.txt --ignore-missing Altum_Darwin_arm64.tar.gz
```

## Usage

Run `altum` from anywhere in your terminal:

```sh
altum
```

## Configuration

Altum uses a configuration file located at `~/.config/altum/config.yaml`.

You can also specify a custom config file using the `--config` flag:

```sh
altum --config /path/to/config.yaml
```

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.