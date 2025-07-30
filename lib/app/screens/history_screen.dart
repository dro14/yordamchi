import 'package:flutter/cupertino.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:yordamchi/core/core.dart';
import 'common.dart';

class HistoryScreen extends ConsumerStatefulWidget {
  const HistoryScreen({super.key});

  static Future<T?> route<T>(BuildContext context) {
    return Navigator.of(context, rootNavigator: true).push<T>(
      CupertinoPageRoute<T>(
        builder: (context) {
          return const HistoryScreen();
        },
      ),
    );
  }

  @override
  ConsumerState<HistoryScreen> createState() => _HistoryScreenState();
}

class _HistoryScreenState extends ConsumerState<HistoryScreen> {
  @override
  void initState() {
    super.initState();
    Analytics.history();
  }

  @override
  Widget build(BuildContext context) {
    final padding = MediaQuery.paddingOf(context);
    final chats = ref.watch(chatsProvider).valueOrNull!;

    return MyScaffold(
      navigationBar: CupertinoNavigationBar(
        padding: const EdgeInsetsDirectional.only(end: 12.0),
        leading: backButton(context),
        middle: Text(l10n.history, style: Fonts.titleMedium),
        border: Border(
          bottom: BorderSide(color: separator(context), width: 0.4),
        ),
      ),
      child: ChatList(
        chats: chats,
        padding: padding.copyWith(top: padding.top + 44.0),
      ),
    );
  }
}
