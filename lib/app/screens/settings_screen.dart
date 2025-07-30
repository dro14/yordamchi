import 'dart:io';

import 'package:flutter/cupertino.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:yordamchi/core/core.dart';
import 'common.dart';

class SettingsScreen extends ConsumerStatefulWidget {
  const SettingsScreen({super.key});

  static Future<T?> route<T>(BuildContext context) {
    return Navigator.of(context, rootNavigator: true).push<T>(
      CupertinoPageRoute<T>(builder: (context) => const SettingsScreen()),
    );
  }

  @override
  ConsumerState<SettingsScreen> createState() => _SettingsScreenState();
}

class _SettingsScreenState extends ConsumerState<SettingsScreen> {
  @override
  void initState() {
    super.initState();
    Analytics.settings();
  }

  @override
  Widget build(BuildContext context) {
    final padding = MediaQuery.paddingOf(context);
    final configs = ref.watch(configsProvider).valueOrNull!;
    final colorIndex = colors.values.toList().indexOf(configs.getColor);
    const locales = AppLocalizations.supportedLocales;
    final localeIndex = configs.languageCode.isNotEmpty
        ? locales.indexOf(configs.getLocale)
        : 0;

    return MyScaffold(
      navigationBar: CupertinoNavigationBar(
        padding: const EdgeInsetsDirectional.only(end: 12.0),
        leading: backButton(context),
        middle: Text(l10n.settings, style: Fonts.titleMedium),
        border: Border(
          bottom: BorderSide(color: separator(context), width: 0.4),
        ),
      ),
      child: Center(
        child: SingleChildScrollView(
          child: Padding(
            padding: padding.copyWith(
              left: padding.left + 30.0,
              top: padding.top + 44.0,
              right: padding.right + 30.0,
            ),
            child: Column(
              spacing: 10.0,
              children: [
                const YordamchiLogo(),
                SizedBox(
                  height: buttonHeight,
                  child: CupertinoButton.tinted(
                    sizeStyle: CupertinoButtonSize.medium,
                    onPressed: () async {
                      final result = await showPicker(
                        context: context,
                        options: languages,
                        initialIndex: localeIndex,
                        title: l10n.language,
                      );
                      if (result != null) {
                        final notifier = ref.read(configsProvider.notifier);
                        notifier.updateConfigs(
                          languageCode: locales[result].languageCode,
                          scriptCode: locales[result].scriptCode,
                          removeScriptCode: locales[result].scriptCode == null,
                        );
                      }
                    },
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Icon(MyIcons.language),
                            const SizedBox(width: 8.0),
                            Text(
                              l10n.language,
                              style: Fonts.bodyMedium.copyWith(
                                color: textColor,
                              ),
                            ),
                          ],
                        ),
                        Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Text(
                              languages[localeIndex],
                              style: Fonts.bodyMedium.copyWith(
                                color: textColor,
                              ),
                            ),
                            const SizedBox(width: 8.0),
                          ],
                        ),
                      ],
                    ),
                  ),
                ),
                SizedBox(
                  height: buttonHeight,
                  child: CupertinoButton.tinted(
                    sizeStyle: CupertinoButtonSize.medium,
                    onPressed: () async {
                      final result = await showPicker(
                        context: context,
                        options: [l10n.system, l10n.light, l10n.dark],
                        initialIndex: switch (configs.isDark) {
                          null => 0,
                          false => 1,
                          true => 2,
                        },
                        title: l10n.brightness,
                      );
                      if (result != null) {
                        final notifier = ref.read(configsProvider.notifier);
                        switch (result) {
                          case 0:
                            notifier.updateConfigs(removeIsDark: true);
                          case 1:
                            notifier.updateConfigs(isDark: false);
                          case 2:
                            notifier.updateConfigs(isDark: true);
                        }
                      }
                    },
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Icon(MyIcons.brightness),
                            const SizedBox(width: 8.0),
                            Text(
                              l10n.brightness,
                              style: Fonts.bodyMedium.copyWith(
                                color: textColor,
                              ),
                            ),
                          ],
                        ),
                        Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Text(
                              getBrightnessLabel(configs.isDark),
                              style: Fonts.bodyMedium.copyWith(
                                color: textColor,
                              ),
                            ),
                            const SizedBox(width: 8.0),
                          ],
                        ),
                      ],
                    ),
                  ),
                ),
                SizedBox(
                  height: buttonHeight,
                  child: CupertinoButton.tinted(
                    sizeStyle: CupertinoButtonSize.medium,
                    onPressed: () {
                      showColorPicker(context, colorIndex, (index) async {
                        final notifier = ref.read(configsProvider.notifier);
                        final color = colors.keys.toList()[index];
                        await notifier.updateConfigs(color: color);
                        if (!Platform.isIOS && context.mounted) {
                          showInfoDialog(
                            context: context,
                            title: l10n.colorChangeTitle,
                            message: l10n.colorChangeMessage,
                          );
                        }
                      });
                    },
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Icon(MyIcons.theme),
                            const SizedBox(width: 8.0),
                            Text(
                              l10n.appColor,
                              style: Fonts.bodyMedium.copyWith(
                                color: textColor,
                              ),
                            ),
                            const SizedBox(width: 8.0),
                          ],
                        ),
                        Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Text(
                              getColorName(colorIndex),
                              style: Fonts.bodyMedium.copyWith(
                                color: textColor,
                              ),
                            ),
                            const SizedBox(width: 8.0),
                          ],
                        ),
                      ],
                    ),
                  ),
                ),
                const SizedBox(height: 10.0),
                SizedBox(
                  height: buttonHeight,
                  child: CupertinoButton.tinted(
                    sizeStyle: CupertinoButtonSize.medium,
                    onPressed: openTermsOfService,
                    child: Row(
                      children: [
                        Icon(MyIcons.termsOfService),
                        const SizedBox(width: 8.0),
                        Text(
                          l10n.termsOfService,
                          style: Fonts.bodyMedium.copyWith(color: textColor),
                        ),
                      ],
                    ),
                  ),
                ),
                SizedBox(
                  height: buttonHeight,
                  child: CupertinoButton.tinted(
                    sizeStyle: CupertinoButtonSize.medium,
                    onPressed: openPrivacyPolicy,
                    child: Row(
                      children: [
                        Icon(MyIcons.privacyPolicy),
                        const SizedBox(width: 8.0),
                        Text(
                          l10n.privacyPolicy,
                          style: Fonts.bodyMedium.copyWith(color: textColor),
                        ),
                      ],
                    ),
                  ),
                ),
                SizedBox(
                  height: buttonHeight,
                  child: CupertinoButton.tinted(
                    sizeStyle: CupertinoButtonSize.medium,
                    onPressed: openContact,
                    child: Row(
                      children: [
                        Icon(MyIcons.contact),
                        const SizedBox(width: 8.0),
                        Text(
                          l10n.contact,
                          style: Fonts.bodyMedium.copyWith(color: textColor),
                        ),
                      ],
                    ),
                  ),
                ),
                const SizedBox(height: 10.0),
                SizedBox(
                  height: buttonHeight,
                  child: CupertinoButton.tinted(
                    sizeStyle: CupertinoButtonSize.medium,
                    onPressed: () async {
                      final chatIds = await getChatIds();
                      if (chatIds.isEmpty && context.mounted) {
                        await showInfoDialog(
                          context: context,
                          title: l10n.noChats,
                          message: l10n.noChatsMessage,
                        );
                      } else if (context.mounted) {
                        final result = await showDestructiveDialog(
                          context: context,
                          title: l10n.deleteAllChatsTitle,
                          message: l10n.deleteAllChatsMessage,
                          cancel: l10n.cancel,
                          confirm: l10n.delete,
                        );
                        if (result) {
                          await clearChats(chatIds);
                          ref.invalidate(chatsProvider);
                          ref.read(currentChatIdProvider.notifier).state = 0;
                          if (context.mounted) {
                            Navigator.pop(context);
                          }
                        }
                      }
                    },
                    child: Row(
                      children: [
                        Icon(MyIcons.delete),
                        const SizedBox(width: 8.0),
                        Text(
                          l10n.deleteAllChatsTitle,
                          style: Fonts.bodyMedium.copyWith(color: textColor),
                        ),
                      ],
                    ),
                  ),
                ),
                const SizedBox(height: 10.0),
                Text(
                  '© ${DateTime.now().year} ChuqurTech\nv$currentVersion',
                  style: Fonts.bodyMedium.copyWith(color: outline),
                  textAlign: TextAlign.center,
                ),
                const SizedBox(height: 10.0),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
