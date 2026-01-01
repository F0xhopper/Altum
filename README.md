# Altum

A CLI tool for managing your daily notes.

## Installation (macOS)

Download the universal binary (works on both Intel and Apple Silicon Macs):

```sh
# Download the latest release
curl -LO https://github.com/F0xhopper/Altum/releases/latest/download/altum_darwin_all.tar.gz

# Extract
tar -xzf altum_darwin_all.tar.gz

# Make executable and install globally
chmod +x altum
sudo mv altum /usr/local/bin/altum

# Verify installation
altum --version

# Cleanup
rm altum_darwin_all.tar.gz
```

**One-liner installation:**

```sh
curl -LO https://github.com/F0xhopper/Altum/releases/latest/download/altum_darwin_all.tar.gz && \
tar -xzf altum_darwin_all.tar.gz && \
chmod +x altum && \
sudo mv altum /usr/local/bin/altum && \
rm altum_darwin_all.tar.gz && \
altum --version
```

**Alternative: Install to user directory (no sudo required):**

```sh
curl -LO https://github.com/F0xhopper/Altum/releases/latest/download/altum_darwin_all.tar.gz
tar -xzf altum_darwin_all.tar.gz
chmod +x altum
mkdir -p ~/bin
mv altum ~/bin/altum

# Add ~/bin to PATH if not already there (add to ~/.zshrc or ~/.bash_profile)
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

altum --version
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