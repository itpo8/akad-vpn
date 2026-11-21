# AkadVPN for iOS — development status

This branch contains the generated Flutter iOS host project. It is **not yet a working iOS VPN**.

## Verified starting point

- iOS Xcode host project generated from the existing Flutter app (Flutter 3.50.0).
- Bundle namespace: `com.akadweb`; display name: `AkadVPN`.
- A macOS GitHub Actions job runs `flutter build ios --release --no-codesign` to identify host and plugin integration failures.
- The Android and desktop core cannot be launched unchanged on iOS. `app/lib/main.dart` must receive an iOS-specific bootstrap.

## Required implementation work

1. **Packet Tunnel Provider extension.** Create an Xcode extension target with `NEPacketTunnelProvider`, its Info.plist, entitlements and an App Group for safe app/extension coordination.
2. **Core integration.** Provide an iOS-compatible build of the Go tunnel engine (e.g. `gomobile bind` or an equivalent static library), then attach it to the packet tunnel flow. Neither the Android MethodChannel service nor desktop executable/service launching is usable on iOS.
3. **Flutter bridge.** Add a dedicated iOS platform API for VPN permission, connection start/stop, status, diagnostics, and persisted settings. Route `main.dart` through that adapter.
4. **Network routing.** Configure IP routes, DNS, MTU, and packet flow with `NEPacketTunnelNetworkSettings` and test real traffic, reconnects, sleep/wake, and IPv4/IPv6 on a physical device.
5. **Branding.** Use the existing Akad Web icon and the unchanged transparent white in-app mark; generate the iOS AppIcon from the original supplied asset.
6. **Signing and installation.** Add a paid Apple Developer Program team, provisioning for both app and extension, and the Network Extensions capability. Store certificates/profiles in GitHub Actions secrets; do not commit them.

## Important restrictions

- A Personal Team (free Apple Account) generally cannot sign the Network Extensions entitlement required for a functioning packet tunnel.
- An unsigned simulator/device build proves compilation only; it is not an installable or usable VPN.
- GitHub Actions macOS can perform builds without the user's Windows PC needing a local Mac, but physical-device installation still needs Apple-authorized provisioning.
- No iOS release should be advertised before physical-device tunnel traffic testing and code-signing verification.

References: [Apple Packet Tunnel Provider](https://developer.apple.com/documentation/NetworkExtension/NEPacketTunnelProvider), [Flutter iOS deployment](https://docs.flutter.dev/deployment/ios).
