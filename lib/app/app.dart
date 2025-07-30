import 'package:flutter/cupertino.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:yordamchi/core/core.dart';
import 'screens/chat_screen.dart';
import 'screens/language_screen.dart';
import 'screens/welcome_screen.dart';

const _localizationsDelegates = AppLocalizations.localizationsDelegates;
const _supportedLocales = AppLocalizations.supportedLocales;

class Yordamchi extends ConsumerWidget {
  const Yordamchi({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return ref
        .watch(configsProvider)
        .when(
          data: (configs) {
            return CupertinoApp(
              debugShowCheckedModeBanner: false,
              title: appName,
              color: configs.getColor,
              theme: _getTheme(configs.getColor, configs.getBrightness),
              locale: configs.getLocale,
              localizationsDelegates: _localizationsDelegates,
              supportedLocales: _supportedLocales,
              home: configs.languageCode.isEmpty
                  ? const LanguageScreen()
                  : id == 0
                  ? const WelcomeScreen()
                  : const ChatScreen(),
            );
          },
          error: (error, stack) {
            Crashlytics.log(error, stack);
            return CupertinoApp(
              debugShowCheckedModeBanner: false,
              title: appName,
              color: colors['blue'],
              theme: _getTheme(colors['blue']!),
              locale: const Locale('uz'),
              localizationsDelegates: _localizationsDelegates,
              supportedLocales: _supportedLocales,
              home: const CupertinoPageScaffold(
                child: Center(child: CupertinoActivityIndicator()),
              ),
            );
          },
          loading: () {
            return CupertinoApp(
              debugShowCheckedModeBanner: false,
              title: appName,
              color: colors['blue'],
              theme: _getTheme(colors['blue']!),
              locale: const Locale('uz'),
              localizationsDelegates: _localizationsDelegates,
              supportedLocales: _supportedLocales,
              home: const CupertinoPageScaffold(
                child: Center(child: CupertinoActivityIndicator()),
              ),
            );
          },
        );
  }
}

CupertinoThemeData _getTheme(Color color, [Brightness? brightness]) {
  brightness ??= WidgetsBinding.instance.platformDispatcher.platformBrightness;
  final colorScheme = getColorScheme(color, brightness);
  primary = colorScheme.primary;
  onPrimary = colorScheme.onPrimary;
  primaryContainer = colorScheme.primaryContainer;
  onPrimaryContainer = colorScheme.onPrimaryContainer;
  onInverseSurface = colorScheme.onInverseSurface;
  scaffold = brightness == Brightness.light ? white : black;
  surfaceContainer = colorScheme.surfaceContainer;
  outline = brightness == Brightness.light
      ? const Color(0xFF57585A)
      : const Color(0xFFA9A9AF);
  textColor = brightness == Brightness.light ? black : white;
  return CupertinoThemeData(
    brightness: brightness,
    primaryColor: primary,
    primaryContrastingColor: onPrimary,
    barBackgroundColor: surfaceContainer.withAlpha(200),
    scaffoldBackgroundColor: scaffold,
    applyThemeToAll: true,
  );
}
