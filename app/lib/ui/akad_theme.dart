import 'package:material_ui/material_ui.dart';

/// Shared Akad Web brand tokens. Keep interface colors out of VPN core logic.
abstract final class AkadTheme {
  static const navy = Color(0xFF060F22);
  static const deepBlue = Color(0xFF0B1C3C);
  static const surface = Color(0xFF10264A);
  static const blue = Color(0xFF147CFF);
  static const cyan = Color(0xFF49D2FF);
  static const pale = Color(0xFFF0F6FF);
  static const ink = Color(0xFF102443);
  static const muted = Color(0xFFA0B8D5);

  static ColorScheme colors(Brightness brightness) {
    final dark = brightness == Brightness.dark;
    return ColorScheme.fromSeed(
      seedColor: blue,
      brightness: brightness,
    ).copyWith(
      primary: dark ? cyan : blue,
      onPrimary: dark ? navy : Colors.white,
      secondary: dark ? blue : const Color(0xFF006FCF),
      surface: dark ? navy : pale,
      onSurface: dark ? Colors.white : ink,
      onSurfaceVariant: dark ? muted : const Color(0xFF526985),
      outlineVariant: dark ? const Color(0xFF29476B) : const Color(0xFFD5E4F5),
      primaryContainer: dark ? const Color(0xFF123A69) : const Color(0xFFD9EBFF),
      onPrimaryContainer: dark ? const Color(0xFFD9F5FF) : ink,
    );
  }
}
