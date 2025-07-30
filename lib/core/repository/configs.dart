import 'dart:io';
import 'dart:ui';

import 'package:flutter/widgets.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_dynamic_icon_plus/flutter_dynamic_icon_plus.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'package:yordamchi/core/models/models.dart';
import 'package:yordamchi/core/utils/utils.dart';
import 'package:yordamchi/core/l10n/l10n.dart';

late int id;
late AppLocalizations l10n;

final configsProvider = AsyncNotifierProvider<ConfigsNotifier, Configs>(() {
  return ConfigsNotifier();
});

class ConfigsNotifier extends AsyncNotifier<Configs> {
  PlatformDispatcher? _platformDispatcher;
  SharedPreferencesWithCache? _prefs;

  @override
  Future<Configs> build() async {
    _platformDispatcher ??= WidgetsBinding.instance.platformDispatcher;
    _prefs ??= await SharedPreferencesWithCache.create(
      cacheOptions: const SharedPreferencesWithCacheOptions(),
    );
    final configs = Configs(
      languageCode: _prefs!.getString('languageCode') ?? '',
      scriptCode: _prefs!.getString('scriptCode'),
      isDark: _prefs!.getBool('isDark'),
      color: _prefs!.getString('color'),
    );
    id = _prefs!.getInt('id') ?? 0;
    minimumVersion = _prefs!.getString('minimumVersion') ?? '1.0.0';
    latestVersion = _prefs!.getString('latestVersion') ?? '1.0.0';
    l10n = await lookupAppLocalizations(configs.getLocale);
    if (configs.isDark == null) {
      _platformDispatcher!.onPlatformBrightnessChanged = () {
        state = AsyncData(state.valueOrNull!);
      };
    }
    return configs;
  }

  Future<void> updateId(int userId) async {
    id = userId;
    _prefs!.setInt('id', userId);
    state = AsyncData(state.valueOrNull!);
  }

  Future<void> updateVersion() async {
    _prefs!.setString('minimumVersion', minimumVersion);
    _prefs!.setString('latestVersion', latestVersion);
  }

  Future<void> updateConfigs({
    String? languageCode,
    String? scriptCode,
    bool? isDark,
    String? color,
    bool removeScriptCode = false,
    bool removeIsDark = false,
  }) async {
    final configs = state.valueOrNull!.copyWith(
      languageCode: languageCode,
      scriptCode: scriptCode,
      isDark: isDark,
      color: color,
      removeScriptCode: removeScriptCode,
      removeIsDark: removeIsDark,
    );
    if (languageCode != null) {
      _prefs!.setString('languageCode', languageCode);
      l10n = await lookupAppLocalizations(configs.getLocale);
    }
    if (scriptCode != null) {
      _prefs!.setString('scriptCode', scriptCode);
    } else if (removeScriptCode && _prefs!.containsKey('scriptCode')) {
      _prefs!.remove('scriptCode');
    }
    if (isDark != null) {
      _prefs!.setBool('isDark', isDark);
      _platformDispatcher!.onPlatformBrightnessChanged = null;
    } else if (removeIsDark && _prefs!.containsKey('isDark')) {
      _prefs!.remove('isDark');
      _platformDispatcher!.onPlatformBrightnessChanged = () {
        state = AsyncData(state.valueOrNull!);
      };
    }
    if (color != null) {
      _prefs!.setString('color', color);
      await _changeAppIcon(color);
    }
    state = AsyncData(configs);
  }

  Future<void> _changeAppIcon(String color) async {
    if (Platform.isIOS) {
      await FlutterDynamicIconPlus.setAlternateIconName(
        iconName: color,
        isSilent: true,
      );
    } else if (Platform.isAndroid) {
      final oldIcon = state.valueOrNull!.color ?? 'default';
      const channel = MethodChannel('uz.chuqurtech.yordamchi/dynamic_icon');
      final args = {'newIcon': color, 'oldIcon': oldIcon};
      await channel.invokeMethod('setAlternateIconName', args);
    }
  }
}
