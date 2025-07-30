import 'package:flutter/cupertino.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:yordamchi/core/repository/repository.dart';
import 'package:yordamchi/core/models/models.dart';
import 'package:yordamchi/core/utils/utils.dart';
import 'package:yordamchi/core/ui/ui.dart';
import 'chat_widget.dart';

class ChatList extends ConsumerStatefulWidget {
  final List<Chat> chats;
  final EdgeInsets padding;
  final bool hasTrailing;

  const ChatList({
    required this.chats,
    required this.padding,
    this.hasTrailing = true,
    super.key,
  });

  @override
  ConsumerState<ChatList> createState() => ChatListState();
}

class ChatListState extends ConsumerState<ChatList> {
  final _key = GlobalKey<ScrollButtonState>();
  final _controller = ScrollController();

  Future<void> scrollToTop() async {
    await _controller.animateTo(
      0.0,
      duration: const Duration(milliseconds: 500),
      curve: Curves.easeInOut,
    );
  }

  @override
  Widget build(BuildContext context) {
    if (widget.chats.isEmpty) {
      return Center(child: Text(l10n.noChats, style: Fonts.titleMedium));
    }

    return Stack(
      children: [
        NotificationListener<ScrollUpdateNotification>(
          onNotification: (notification) {
            _key.currentState?.changeState(notification.metrics.pixels > 100);
            return true;
          },
          child: ListView.builder(
            padding: widget.padding,
            controller: _controller,
            itemCount: widget.chats.length,
            itemBuilder: (context, index) {
              final chat = widget.chats[index];
              final isFirst = index == 0;
              final prev = isFirst ? chat : widget.chats[index - 1];
              final String? date;
              if (isFirst || !isSameDay(prev.dateTime, chat.dateTime)) {
                date = formatDate(chat.dateTime);
              } else {
                date = null;
              }
              return Column(
                children: [
                  if (date != null)
                    Padding(
                      padding: const EdgeInsets.symmetric(vertical: 4.0),
                      child: Text(date, style: Fonts.bodyMedium),
                    ),
                  ChatWidget(chat: chat, hasTrailing: widget.hasTrailing),
                ],
              );
            },
          ),
        ),
        Positioned(
          right: widget.padding.right + 8.0,
          bottom: widget.padding.bottom + 8.0,
          child: ScrollButton(
            key: _key,
            icon: MyIcons.up,
            onPressed: scrollToTop,
          ),
        ),
      ],
    );
  }
}
