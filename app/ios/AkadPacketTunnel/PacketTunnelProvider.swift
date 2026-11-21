import Foundation
import NetworkExtension

/// Packet Tunnel entry point. This is intentionally fail-closed until the
/// iOS Go packet engine and packetFlow reader/writer are implemented.
/// Do not set network routes without a functioning engine: doing so can
/// silently black-hole the phone's traffic.
final class PacketTunnelProvider: NEPacketTunnelProvider {
  override func startTunnel(
    options: [String: NSObject]?,
    completionHandler: @escaping (Error?) -> Void
  ) {
    completionHandler(NSError(
      domain: "com.akadweb.akadvpn.packet-tunnel",
      code: 1,
      userInfo: [NSLocalizedDescriptionKey:
        "AkadVPN iOS packet engine is not integrated yet."]
    ))
  }

  override func stopTunnel(
    with reason: NEProviderStopReason,
    completionHandler: @escaping () -> Void
  ) {
    completionHandler()
  }
}
