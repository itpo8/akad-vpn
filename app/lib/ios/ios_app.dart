import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

/// Entry point for the native iOS integration stage. The existing desktop
/// executable and Android foreground-service channels cannot run on iOS.
/// The UI explicitly prevents a false "connected" state until the packet
/// tunnel extension and Go core have been tested on a physical device.
Future<void> runIOS() async {
  await SystemChrome.setPreferredOrientations([
    DeviceOrientation.portraitUp,
    DeviceOrientation.portraitDown,
  ]);
  runApp(const _AkadIOSApp());
}

class _AkadIOSApp extends StatefulWidget {
  const _AkadIOSApp();

  @override
  State<_AkadIOSApp> createState() => _AkadIOSAppState();
}

class _AkadIOSAppState extends State<_AkadIOSApp> {
  static const _bridge = MethodChannel('akad_vpn/ios');
  String _status = 'disconnected';

  @override
  void initState() {
    super.initState();
    _refreshStatus();
  }

  Future<void> _refreshStatus() async {
    try {
      final state = await _bridge.invokeMapMethod<String, dynamic>('status');
      if (!mounted) return;
      setState(() => _status = state?['status'] as String? ?? 'disconnected');
    } on PlatformException {
      if (mounted) setState(() => _status = 'unavailable');
    }
  }

  @override
  Widget build(BuildContext context) {
    const navy = Color(0xFF081739);
    const cyan = Color(0xFF49D2FF);
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      title: 'AkadVPN',
      theme: ThemeData.dark(useMaterial3: true),
      home: Scaffold(
        backgroundColor: navy,
        body: Container(
          decoration: const BoxDecoration(
            gradient: LinearGradient(
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
              colors: [Color(0xFF09142E), Color(0xFF102E67), navy],
            ),
          ),
          child: SafeArea(
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 28),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Image.asset('assets/logo-white.png', width: 96, height: 96),
                  const SizedBox(height: 20),
                  const Text('AkadVPN', style: TextStyle(
                    fontFamily: 'Inter', fontSize: 30, fontWeight: FontWeight.w700)),
                  const SizedBox(height: 12),
                  const Text('iOS development build', style: TextStyle(
                    color: cyan, fontSize: 16, fontWeight: FontWeight.w600)),
                  const SizedBox(height: 32),
                  Container(
                    width: double.infinity,
                    padding: const EdgeInsets.all(22),
                    decoration: BoxDecoration(
                      color: Colors.white.withValues(alpha: 0.07),
                      borderRadius: BorderRadius.circular(22),
                      border: Border.all(color: cyan.withValues(alpha: 0.22)),
                    ),
                    child: Column(children: [
                      const Icon(Icons.shield_outlined, color: cyan, size: 40),
                      const SizedBox(height: 12),
                      Text('VPN status: $_status', style: const TextStyle(fontSize: 15)),
                      const SizedBox(height: 14),
                      const Text(
                        'The Packet Tunnel extension and VPN engine are being integrated. '
                        'Connections are disabled in this development build.',
                        textAlign: TextAlign.center,
                        style: TextStyle(color: Color(0xFFBAC8DF), height: 1.5),
                      ),
                    ]),
                  ),
                  const SizedBox(height: 20),
                  OutlinedButton.icon(
                    onPressed: _refreshStatus,
                    icon: const Icon(Icons.refresh),
                    label: const Text('Refresh status'),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
