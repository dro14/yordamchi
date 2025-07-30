import 'package:flutter/cupertino.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:yordamchi/core/repository/repository.dart';
import 'package:yordamchi/core/models/models.dart';
import 'package:yordamchi/core/ui/ui.dart';

class ChatWidget extends ConsumerWidget {
  final Chat chat;
  final bool hasTrailing;

  const ChatWidget({required this.chat, this.hasTrailing = true, super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    if (chat.messages.isEmpty) {
      return const SizedBox.shrink();
    }

    final isSelected = chat.id == ref.watch(currentChatIdProvider);

    return GestureDetector(
      onTap: () {
        Navigator.pop(context, true);
        ref.read(currentChatIdProvider.notifier).state = chat.id;
      },
      behavior: HitTestBehavior.opaque,
      child: Container(
        margin: const EdgeInsets.symmetric(horizontal: 12.0, vertical: 4.0),
        padding: const EdgeInsets.all(8.0),
        decoration: BoxDecoration(
          color: isSelected ? primaryContainer : surfaceContainer,
          borderRadius: BorderRadius.circular(12.0),
        ),
        child: Row(
          spacing: 8.0,
          children: [
            Expanded(
              child: Column(
                spacing: 4.0,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  if (chat.name.isNotEmpty)
                    Padding(
                      padding: const EdgeInsets.symmetric(vertical: 4.0),
                      child: MyText(
                        chat.name.replaceAll('\n', ' '),
                        overflow: TextOverflow.ellipsis,
                        style: Fonts.titleMedium.copyWith(
                          color: isSelected ? onPrimaryContainer : textColor,
                        ),
                        isSelectable: false,
                      ),
                    )
                  else ...[
                    MyText(
                      chat.messages.first.displayText,
                      overflow: TextOverflow.ellipsis,
                      style: Fonts.titleSmall.copyWith(
                        color: isSelected ? onPrimaryContainer : textColor,
                      ),
                      isSelectable: false,
                    ),
                    MyText(
                      chat.messages.last.displayText,
                      overflow: TextOverflow.ellipsis,
                      style: Fonts.bodySmall.copyWith(
                        color: isSelected ? onPrimaryContainer : outline,
                      ),
                      isSelectable: false,
                    ),
                  ],
                ],
              ),
            ),
            if (hasTrailing)
              CupertinoButton(
                padding: EdgeInsets.zero,
                sizeStyle: CupertinoButtonSize.small,
                onPressed: () async {
                  final result = await showActionSheet(
                    context: context,
                    icons: [MyIcons.edit, MyIcons.delete],
                    actions: [l10n.rename, l10n.delete],
                  );
                  if (result != null && context.mounted) {
                    switch (result) {
                      case 0:
                        final name = await showInputDialog(context: context);
                        if (name != null) {
                          ref.read(chatProvider(chat.id).notifier).rename(name);
                        }
                      case 1:
                        final result = await showDestructiveDialog(
                          context: context,
                          title: l10n.deleteChatTitle,
                          message: l10n.deleteChatMessage,
                          cancel: l10n.cancel,
                          confirm: l10n.delete,
                        );
                        if (result) {
                          await clearChats([chat.id]);
                          ref.invalidate(chatsProvider);
                          if (ref.read(currentChatIdProvider) == chat.id) {
                            ref.read(currentChatIdProvider.notifier).state = 0;
                          }
                        }
                    }
                  }
                },
                child: Icon(
                  MyIcons.options,
                  color: isSelected ? onPrimaryContainer : primary,
                  size: 20.0,
                ),
              ),
          ],
        ),
      ),
    );
  }
}
