import Flutter
import NetworkExtension
import UIKit

@main
@objc class AppDelegate: FlutterAppDelegate, FlutterImplicitEngineDelegate {
  override func application(
    _ application: UIApplication,
    didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?
  ) -> Bool {
    return super.application(application, didFinishLaunchingWithOptions: launchOptions)
  }

  func didInitializeImplicitFlutterEngine(_ engineBridge: FlutterImplicitEngineBridge) {
    GeneratedPluginRegistrant.register(with: engineBridge.pluginRegistry)
    if let registrar = engineBridge.pluginRegistry.registrar(forPlugin: "AkadVPNiOS") {
      AkadVPNiOSBridge.register(with: registrar)
    }
  }
}

/// Native boundary for the eventual Packet Tunnel Provider. Connection actions
/// deliberately fail closed until the signed extension and Go core are wired.
/// This avoids accidentally presenting a connected state without a VPN tunnel.
private final class AkadVPNiOSBridge: NSObject {
  private static var activeBridge: AkadVPNiOSBridge?
  private let channel: FlutterMethodChannel

  private init(messenger: FlutterBinaryMessenger) {
    channel = FlutterMethodChannel(name: "akad_vpn/ios", binaryMessenger: messenger)
    super.init()
    channel.setMethodCallHandler { [weak self] call, result in
      self?.handle(call, result: result)
    }
  }

  static func register(with registrar: FlutterPluginRegistrar) {
    activeBridge = AkadVPNiOSBridge(messenger: registrar.messenger())
  }

  private func handle(_ call: FlutterMethodCall, result: @escaping FlutterResult) {
    switch call.method {
    case "status":
      NETunnelProviderManager.loadAllFromPreferences { managers, error in
        DispatchQueue.main.async {
          if let error = error {
            result(FlutterError(code: "VPN_STATUS", message: error.localizedDescription, details: nil))
            return
          }
          let status = managers?.first?.connection.status
          result(["status": self.statusName(status), "engineReady": false])
        }
      }
    case "connect":
      result(FlutterError(
        code: "PACKET_TUNNEL_NOT_READY",
        message: "The signed Packet Tunnel extension and Go packet engine are not integrated yet.",
        details: nil
      ))
    case "disconnect":
      NETunnelProviderManager.loadAllFromPreferences { managers, error in
        DispatchQueue.main.async {
          if let error = error {
            result(FlutterError(code: "VPN_STOP", message: error.localizedDescription, details: nil))
            return
          }
          managers?.forEach { $0.connection.stopVPNTunnel() }
          result(nil)
        }
      }
    default:
      result(FlutterMethodNotImplemented)
    }
  }

  private func statusName(_ status: NEVPNStatus?) -> String {
    switch status {
    case .connected: return "connected"
    case .connecting: return "connecting"
    case .disconnecting: return "disconnecting"
    case .reasserting: return "reasserting"
    case .invalid: return "invalid"
    default: return "disconnected"
    }
  }
}
