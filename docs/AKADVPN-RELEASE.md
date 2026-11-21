# AkadVPN — first public release

This fork inherits the GPLv3 license and cross-platform Flutter/Go build pipeline from [Arman VPN](https://github.com/armanlife/arman-vpn-client). Keep the LICENSE and source attribution in published distributions. No build or release should be represented as endorsed by Mozilla or Arman VPN.

## 1. GitHub repository

Create a GitHub repository for AkadVPN. The initial assumed repository is `itpo8/akad-vpn`. If your actual owner or repository differs, replace **every occurrence** of `itpo8/akad-vpn` and `github.com/itpo8/akad-vpn` in the sources, CI scripts, update URLs, and documentation before building. Do not leave update URLs pointing at another publisher.

Upload the full source including `.github/workflows` and push the `main` branch. Enable GitHub Actions with read/write permissions appropriate for release publishing. The upstream `Build` workflow requires `contents: write` for its release job and uses the repository's default `GITHUB_TOKEN`; Docker publishing also requires package permissions.

## 2. Release signing

Configure GitHub Actions secrets before releasing:

- `ANDROID_KEYSTORE_BASE64`: base64 of your own release keystore, never the upstream keystore.
- `ANDROID_KEYSTORE_PASSWORD`: its keystore password.
- `FIREFOX_VPN_SIGNING_KEY`: your own signing key for signed desktop/core updates, as described in `scripts/signing-key.ps1`.

Securely back up release keys. An Android app signed with a different certificate cannot upgrade an existing installation of the same package. AkadVPN uses a distinct Android package ID, so it installs separately from Arman VPN.

## 3. Build and publish

1. Check `app/pubspec.yaml` version, currently `1.0.19+157`.
2. Ensure `.github/updates.json` is not ahead of that version and update it when offering a new release.
3. Push and let CI tests succeed.
4. Push a tag matching the version exactly: `git tag v1.0.19` then `git push origin v1.0.19`.
5. Inspect Actions: Android, Windows, macOS and Linux artifacts and Docker image must all succeed before announcing the release.
6. Verify downloads, installation and VPN connectivity on real devices. This fork has not been verified on all four operating systems in this environment.

Do not upload or distribute an unsigned/debug Android APK as a production release. Do not publish the original project's signing keys or credentials.

## 4. Before claiming a production-ready release

- Review all network destinations and data-handling behavior; this VPN authenticates against third-party Firefox/Mozilla services and depends on their availability and terms.
- Check privacy notices, license attribution, and your own project support contacts.
- Run `go test ./...` from `core/`, `flutter analyze` and `flutter test` from `app/`, then test release builds and connectivity.
- Supply brand artwork, app store metadata and actual publisher identity. This draft includes a navy/cyan icon variant, not a finalized brand identity.
