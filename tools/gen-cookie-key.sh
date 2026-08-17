#!/usr/bin/env bash
# Generate a Base64-encoded 32-byte AES-256 key for cookie encryption.
# Usage: ./tools/gen-cookie-key.sh
#
# The generated key is safe for use as the COOKIE_ENCRYPTION_KEY environment
# variable. Set it in your .env file or as a system/Fly.io secret.

set -euo pipefail

KEY=$(openssl rand -base64 32)

echo ""
echo "COOKIE_ENCRYPTION_KEY=${KEY}"
echo ""
echo "Add this to your .env file:"
echo ""
echo "  COOKIE_ENCRYPTION_KEY=${KEY}"
echo ""
echo "Or set it as a Fly.io secret:"
echo ""
echo "  fly secrets set COOKIE_ENCRYPTION_KEY=${KEY}"
