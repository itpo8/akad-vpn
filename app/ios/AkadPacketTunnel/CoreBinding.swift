import Foundation
import AkadCore

/// Verifies that the gomobile-generated AkadCore module is accessible inside
/// the Packet Tunnel target. Actual Start/Stop orchestration is intentionally
/// deferred until packetFlow forwarding and protected egress are implemented.
enum AkadCoreBinding {
  static let frameworkLinked = true
}
