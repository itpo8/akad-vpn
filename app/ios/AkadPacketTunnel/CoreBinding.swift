import Foundation
import AkadCore

/// Owns the Go core within the Packet Tunnel process, not the Flutter host.
/// Start is deliberately *not* called by PacketTunnelProvider yet: packetFlow
/// ingress/egress and protected upstream sockets must be ready first.
enum AkadCoreBinding {
  private static let appGroup = "group.com.akadweb.akadvpn"

  static func dataDirectory() throws -> String {
    guard let groupURL = FileManager.default.containerURL(
      forSecurityApplicationGroupIdentifier: appGroup
    ) else {
      throw NSError(
        domain: "com.akadweb.akadvpn.core", code: 1,
        userInfo: [NSLocalizedDescriptionKey:
          "AkadVPN App Group is unavailable; check signing and provisioning."]
      )
    }
    let url = groupURL.appendingPathComponent("CoreData", isDirectory: true)
    try FileManager.default.createDirectory(
      at: url, withIntermediateDirectories: true
    )
    // VPN extensions run in the background after the first device unlock.
    try FileManager.default.setAttributes(
      [.protectionKey: FileProtectionType.completeUntilFirstUserAuthentication],
      ofItemAtPath: url.path
    )
    return url.path
  }

  /// Run the daemon with proxy disabled until packet forwarding is ready.
  /// Any failure propagates back to the Packet Tunnel lifecycle controller.
  static func startDisabled() throws {
    try MobileStart(dataDirectory(), false)
  }

  static func stop() {
    MobileStop()
  }

  static var running: Bool {
    MobileRunning()
  }
}
