# AkadVPN — signing keys and recovery

This document contains **no private keys or passwords**. Never commit or paste
`ANDROID_KEYSTORE_BASE64`, `ANDROID_KEYSTORE_PASSWORD`, or
`FIREFOX_VPN_SIGNING_KEY` into a document, issue, chat, or source file.

## GitHub Actions release secrets

| Secret | Purpose |
| --- | --- |
| `ANDROID_KEYSTORE_BASE64` | Base64 encoding of the Android release JKS file |
| `ANDROID_KEYSTORE_PASSWORD` | Password for the JKS keystore and its `release` key alias |
| `FIREFOX_VPN_SIGNING_KEY` | Base64 **32-byte Ed25519 seed** for signed application updates (historical secret name) |

The workflow currently requires all three. `FIREFOX_VPN_SIGNING_KEY` is an
internal CI compatibility name and does not alter the AkadVPN application name.

## Recovery plan

1. Generate both keys on a computer you control using Java `keytool` and
   `go run ./core/tools/sign -keygen` (from within `core/`).
2. Create an **encrypted, offline backup** of the JKS, the base64 Ed25519 seed,
   and the Android keystore password. Keep its decryption password in your own
   password manager, separate from the backup.
3. Add the three values in GitHub repository Settings → Secrets and variables
   → Actions → Repository secrets. Do not use repository Variables for secrets.
4. Verify the public Ed25519 key and Android signing certificate fingerprint,
   and note them here (public values are safe to document).
5. Test GitHub Actions on a non-release branch before tagging a release.

## Public verification metadata

- Android signing certificate SHA-256: _to be filled after generation_
- Ed25519 update public key (base64): _to be filled after generation_
- First signed release: _pending_
- Secure backup location: _record location only, not password or key_

## Important

Losing the Android signing key prevents seamless upgrades of the Android app
using its existing application ID. Losing the Ed25519 seed can break the
desktop auto-update trust chain. Do not rotate these casually.

This repository also retains GPLv3 notices and proper upstream attribution.
