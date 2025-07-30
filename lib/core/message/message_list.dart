import 'dart:math';
import 'dart:async';

import 'package:flutter/cupertino.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:visibility_detector/visibility_detector.dart';
import 'package:scrollable_positioned_list/scrollable_positioned_list.dart';

import 'package:yordamchi/core/repository/repository.dart';
import 'package:yordamchi/core/models/models.dart';
import 'package:yordamchi/core/utils/utils.dart';
import 'package:yordamchi/core/ui/ui.dart';
import 'message_widget.dart';

class MessageList extends ConsumerStatefulWidget {
  final int chatId;
  final EdgeInsets padding;
  final ItemScrollController scrollController;
  final void Function(Message) onReply;
  final void Function(Message) onEdit;
  final void Function(List<String>, int) onImageTap;

  const MessageList({
    required this.chatId,
    required this.padding,
    required this.scrollController,
    required this.onReply,
    required this.onEdit,
    required this.onImageTap,
    super.key,
  });

  @override
  ConsumerState<MessageList> createState() => MessageListState();
}

class MessageListState extends ConsumerState<MessageList>
    with TickerProviderStateMixin {
  final _animControllers = <AnimationController>[];
  final _key = GlobalKey<ScrollButtonState>();
  var _lastIndex = 0;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _key.currentState?.changeState(true);
    });
  }

  @override
  void dispose() {
    for (final animController in _animControllers) {
      animController.dispose();
    }
    super.dispose();
  }

  double get _alignment {
    return (widget.padding.top + 6.0) / MediaQuery.sizeOf(context).height;
  }

  Future<void> scrollTo(int messageId, {bool animate = true}) async {
    final chat = ref.read(chatProvider(widget.chatId)).valueOrNull!;
    final index = chat.messages.indexWhere((m) => m.id == messageId);
    if (index != -1) {
      await widget.scrollController.scrollTo(
        index: min(index, _lastIndex) + 1,
        duration: const Duration(milliseconds: 500),
        curve: Curves.easeInOut,
        alignment: _alignment,
      );
      if (animate) {
        Future.delayed(const Duration(milliseconds: 100), () {
          _animControllers[index].forward();
        });
      }
    }
  }

  Future<void> animate(int messageId) async {
    final chat = ref.read(chatProvider(widget.chatId)).valueOrNull!;
    final index = chat.messages.indexWhere((m) => m.id == messageId);
    if (index != -1) {
      _animControllers[index].forward();
    }
  }

  Future<void> scrollToLastIndex() async {
    await widget.scrollController.scrollTo(
      index: _lastIndex + 1,
      duration: const Duration(milliseconds: 500),
      curve: Curves.easeInOut,
      alignment: _alignment,
    );
  }

  Future<void> scrollToBottom() async {
    await widget.scrollController.scrollTo(
      index: _lastIndex + 2,
      duration: const Duration(milliseconds: 500),
      curve: Curves.easeInOut,
      alignment: 1.0 - 1.0 / MediaQuery.sizeOf(context).height,
    );
  }

  Future<void> onRetry(Message message, Message previous) async {
    final notifier = ref.read(chatProvider(widget.chatId).notifier);
    if (message.inReplyTo == 0) {
      onFollowUps(previous);
    } else {
      notifier.edit(notifier.get(message.inReplyTo)!);
    }
  }

  Future<void> onEdit(Message message, Message previous) async {
    final notifier = ref.read(chatProvider(widget.chatId).notifier);
    if (message.inReplyTo == 0) {
      widget.onEdit(notifier.get(previous.inReplyTo)!);
    } else {
      widget.onEdit(notifier.get(message.inReplyTo)!);
    }
  }

  Future<void> onFollowUps(Message message) async {
    final notifier = ref.read(chatProvider(widget.chatId).notifier);
    await notifier.followUps(message);
    Future.delayed(const Duration(milliseconds: 100), scrollToLastIndex);
  }

  Future<void> onFollowUpTap(String text) async {
    final message = Message(
      id: DateTime.now().millisecondsSinceEpoch,
      userId: id,
      chatId: widget.chatId,
      role: 'user',
      createdAt: DateTime.now().millisecondsSinceEpoch,
      text: text,
    );
    Future.delayed(const Duration(milliseconds: 100), scrollToLastIndex);
    final notifier = ref.read(chatProvider(widget.chatId).notifier);
    await notifier.create(message);
  }

  @override
  Widget build(BuildContext context) {
    final asyncChat = ref.watch(chatProvider(widget.chatId));
    if (asyncChat.hasError) {
      Crashlytics.log(asyncChat.error!, asyncChat.stackTrace);
    }

    final Chat chat;
    if (asyncChat.hasValue) {
      chat = asyncChat.value!;
    } else {
      return const SizedBox.shrink();
    }

    final messages = chat.messages;
    for (int i = _animControllers.length; i < messages.length + 1; i++) {
      final animController = AnimationController(
        vsync: this,
        duration: const Duration(milliseconds: 500),
      );
      animController.addStatusListener((status) {
        if (status == AnimationStatus.completed) {
          animController.reverse();
        }
      });
      _animControllers.add(animController);
    }

    for (int i = 0; i < messages.length; i++) {
      if ((messages[i].role == 'user' && messages[i].displayText.isNotEmpty) ||
          messages[i].followUps.isNotEmpty) {
        _lastIndex = i;
      }
    }

    final notifier = ref.read(chatProvider(widget.chatId).notifier);
    final chunk = ref.watch(chunkProvider(widget.chatId));

    return Stack(
      children: [
        ScrollablePositionedList.builder(
          padding: EdgeInsets.only(
            left: widget.padding.left,
            right: widget.padding.right,
          ),
          initialAlignment: _alignment,
          initialScrollIndex: _lastIndex + 1,
          itemCount: _lastIndex + 3,
          itemScrollController: widget.scrollController,
          itemBuilder: (context, index) {
            if (index == 0) {
              return SizedBox(height: widget.padding.top + 6.0);
            } else if (index == _lastIndex + 2) {
              return VisibilityDetector(
                key: Key('${widget.chatId}_$index'),
                onVisibilityChanged: (visibilityInfo) {
                  _key.currentState?.changeState(
                    visibilityInfo.visibleFraction == 0.0,
                  );
                },
                child: const SizedBox(height: 1.0),
              );
            }

            final message = messages[index - 1];
            final previous = index == 1 ? message : messages[index - 2];
            final isLast = index == _lastIndex + 1;
            final minHeight = isLast
                ? MediaQuery.sizeOf(context).height - widget.padding.top - 7.0
                : 0.0;

            return ConstrainedBox(
              constraints: BoxConstraints(minHeight: minHeight),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  MessageWidget(
                    animController: _animControllers[index - 1],
                    message: message,
                    inReplyTo: notifier.get(message.inReplyTo),
                    onReply: () => widget.onReply(message),
                    onRetry: () => onRetry(message, previous),
                    onEdit: () => onEdit(message, previous),
                    onReplyTap: () => scrollTo(message.inReplyTo),
                    onFollowUps: () => onFollowUps(message),
                    onFollowUpTap: onFollowUpTap,
                    onImageTap: widget.onImageTap,
                  ),
                  if (isLast) ...[
                    for (int i = _lastIndex + 1; i < messages.length; i++) ...[
                      MessageWidget(
                        animController: _animControllers[i],
                        message: messages[i],
                        inReplyTo: notifier.get(messages[i].inReplyTo),
                        onReply: () => widget.onReply(messages[i]),
                        onRetry: () => onRetry(messages[i], messages[i - 1]),
                        onEdit: () => onEdit(messages[i], messages[i - 1]),
                        onReplyTap: () => scrollTo(messages[i].inReplyTo),
                        onFollowUps: () => onFollowUps(messages[i]),
                        onFollowUpTap: onFollowUpTap,
                        onImageTap: widget.onImageTap,
                      ),
                    ],
                    if (chunk != null)
                      MessageWidget(
                        animController: _animControllers[messages.length],
                        chunk: chunk,
                      ),
                    SizedBox(height: widget.padding.bottom),
                  ],
                ],
              ),
            );
          },
        ),
        Positioned(
          right: widget.padding.right + 8.0,
          bottom: widget.padding.bottom + 8.0,
          child: ScrollButton(
            key: _key,
            icon: MyIcons.down,
            onPressed: scrollToBottom,
          ),
        ),
      ],
    );
  }
}
