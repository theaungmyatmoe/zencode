#!/data/data/com.termux/files/usr/bin/bash
# Zencode Termux bootstrap (early version)
# Usage: curl -fsSL <raw-url-to-this-script> | bash
# Or run locally after cloning.

set -euo pipefail

echo "==> Zencode Termux bootstrap (dev)"

pkg update -y || true
pkg install -y git ripgrep || true

# We prefer prebuilt binaries. This script can also help set up for source builds.
echo "==> Checking for existing zencode binary in PATH..."
if command -v zencode >/dev/null 2>&1; then
    echo "zencode is already on PATH: $(command -v zencode)"
else
    echo "No zencode in PATH yet."
    echo "After building or downloading the binary, put it in \$PREFIX/bin or ~/bin and chmod +x it."
fi

echo
echo "Next steps (manual for now while we are in early dev):"
echo "  1. export XAI_API_KEY=your-xai-key"
echo "  2. cd /path/to/your/project"
echo "  3. zencode"
echo
echo "Tip: termux-setup-storage   # if you want easy access to ~/storage/shared"
echo "Tip: pkg install golang     # only if you plan to build from source inside Termux (quite heavy)"
echo
echo "Bootstrap done."
