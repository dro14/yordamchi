import 'package:flutter/cupertino.dart';

import 'package:yordamchi/core/ui/ui.dart';

@immutable
class Configs {
  final String languageCode;
  final String? scriptCode;
  final bool? isDark;
  final String? color;

  const Configs({
    required this.languageCode,
    required this.scriptCode,
    required this.isDark,
    required this.color,
  });

  Locale get getLocale {
    return languageCode.isNotEmpty
        ? Locale.fromSubtags(languageCode: languageCode, scriptCode: scriptCode)
        : const Locale('uz');
  }

  Brightness? get getBrightness {
    switch (isDark) {
      case null:
        return null;
      case false:
        return Brightness.light;
      case true:
        return Brightness.dark;
    }
  }

  Color get getColor {
    return colors[color ?? 'blue']!;
  }

  Configs copyWith({
    String? languageCode,
    String? scriptCode,
    bool? isDark,
    String? color,
    bool removeScriptCode = false,
    bool removeIsDark = false,
  }) {
    return Configs(
      languageCode: languageCode ?? this.languageCode,
      scriptCode: removeScriptCode ? null : scriptCode ?? this.scriptCode,
      isDark: removeIsDark ? null : isDark ?? this.isDark,
      color: color ?? this.color,
    );
  }

  @override
  String toString() {
    var string = '';
    string += 'languageCode: $languageCode\n';
    string += 'scriptCode: $scriptCode\n';
    string += 'isDark: $isDark\n';
    string += 'color: $color\n';
    return string;
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is Configs &&
        other.runtimeType == runtimeType &&
        other.languageCode == languageCode &&
        other.scriptCode == scriptCode &&
        other.isDark == isDark &&
        other.color == color;
  }

  @override
  int get hashCode {
    return Object.hash(languageCode, scriptCode, isDark, color);
  }
}
