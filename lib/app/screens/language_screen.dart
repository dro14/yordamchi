import 'package:flutter/cupertino.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:yordamchi/core/core.dart';

class LanguageScreen extends ConsumerWidget {
  const LanguageScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final padding = MediaQuery.paddingOf(context);
    const locales = AppLocalizations.supportedLocales;

    return CupertinoPageScaffold(
      child: Center(
        child: SingleChildScrollView(
          child: Padding(
            padding: padding.copyWith(
              left: padding.left + 50.0,
              right: padding.right + 50.0,
            ),
            child: Column(
              spacing: 20.0,
              children: [
                for (int index = 0; index < languages.length; index++)
                  SizedBox(
                    width: double.infinity,
                    height: buttonHeight,
                    child: CupertinoButton.tinted(
                      sizeStyle: CupertinoButtonSize.medium,
                      onPressed: () {
                        final notifier = ref.read(configsProvider.notifier);
                        notifier.updateConfigs(
                          languageCode: locales[index].languageCode,
                          scriptCode: locales[index].scriptCode,
                          removeScriptCode: locales[index].scriptCode == null,
                        );
                      },
                      child: Text(languages[index], style: Fonts.titleSmall),
                    ),
                  ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
