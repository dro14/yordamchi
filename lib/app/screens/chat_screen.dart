import 'dart:io';
import 'dart:async';

import 'package:flutter/cupertino.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:http/http.dart';
import 'package:scrollable_positioned_list/scrollable_positioned_list.dart';
import 'package:image_picker/image_picker.dart';

import 'package:yordamchi/core/core.dart';
import 'common.dart';
import 'history_screen.dart';
import 'image_screen.dart';
import 'settings_screen.dart';

class ChatScreen extends ConsumerStatefulWidget {
  const ChatScreen({super.key});

  static Future<T?> route<T>(BuildContext context) {
    return Navigator.of(context, rootNavigator: true).push<T>(
      CupertinoPageRoute<T>(
        builder: (context) {
          return const ChatScreen();
        },
      ),
    );
  }

  @override
  ConsumerState<ChatScreen> createState() => _ChatScreenState();
}

class _ChatScreenState extends ConsumerState<ChatScreen> {
  final _scrollController = ItemScrollController();
  final _list = GlobalKey<MessageListState>();
  final _focusNode = FocusNode();
  final _input = GlobalKey();
  double? _messageInputHeight;
  bool _isLoading = false;
  String? _error;

  final _controller = TextEditingController();
  final _images = <String>[];
  bool _isEmpty = true;
  Message? _inReplyTo;
  Message? _message;

  int _currentPromptIndex = 0;
  Timer? _promptTimer;

  @override
  void initState() {
    super.initState();
    Analytics.chat();
    _startPromptAnimation();
    WidgetsBinding.instance.addPostFrameCallback((_) async {
      if (!ref.exists(chatsProvider)) {
        await ref.read(chatsProvider.future);
      }
      _focusNode.requestFocus();
      while (!hasCheckedVersion) {
        final bool? haveTo;
        try {
          haveTo = await haveToUpdate();
        } catch (error) {
          await Future.delayed(const Duration(seconds: 3));
          continue;
        }
        ref.read(configsProvider.notifier).updateVersion();
        if (haveTo != null && mounted) {
          var message = l10n.updateMessagePrefix;
          if (haveTo) {
            message += l10n.updateMessageSuffix;
          }
          var willUpdate = await showInfoDialog(
            context: context,
            title: l10n.updateTitle,
            message: message,
            cancel: haveTo ? null : l10n.notNow,
            confirm: l10n.update,
          );
          while (haveTo && !willUpdate) {
            await Future.delayed(const Duration(milliseconds: 500));
            if (mounted) {
              willUpdate = await showInfoDialog(
                context: context,
                title: l10n.updateTitle,
                message: message,
                confirm: l10n.update,
              );
            }
          }
          if (willUpdate) openUpdate();
        }
        hasCheckedVersion = true;
      }
    });
  }

  @override
  void dispose() {
    _controller.dispose();
    _promptTimer?.cancel();
    super.dispose();
  }

  void _startPromptAnimation() {
    _promptTimer = Timer.periodic(const Duration(seconds: 4), (timer) {
      if (mounted) {
        setState(() {
          _currentPromptIndex = (_currentPromptIndex + 1) % 10;
        });
      }
    });
  }

  void _onError(Object error, StackTrace stack) {
    Crashlytics.log(error, stack);
    setState(() {
      _isLoading = false;
      _error = error is ClientException || error is SocketException
          ? l10n.noInternet
          : l10n.errorMessage;
    });
  }

  void _onSend(int chatId, [String? prompt]) async {
    setState(() => _isLoading = true);

    if (id == 0) {
      final int userId;
      try {
        userId = await newUser();
      } catch (error, stack) {
        _onError(error, stack);
        return;
      }
      ref.read(configsProvider.notifier).updateId(userId);
    }

    if (chatId == 0) {
      try {
        chatId = await newChat();
      } catch (error, stack) {
        _onError(error, stack);
        return;
      }
      await ref.read(chatProvider(chatId).future);
      ref.read(currentChatIdProvider.notifier).state = chatId;
    }

    final Future<void> future;
    if (_message != null) {
      final message = _message!.copyWith(
        createdAt: DateTime.now().millisecondsSinceEpoch,
        text: prompt ?? _controller.text.trim(),
        images: List<String>.from(_images),
      );
      future = ref.read(chatProvider(chatId).notifier).edit(message);
    } else {
      final message = Message(
        id: DateTime.now().millisecondsSinceEpoch,
        userId: id,
        chatId: chatId,
        role: 'user',
        createdAt: DateTime.now().millisecondsSinceEpoch,
        inReplyTo: _inReplyTo?.id ?? 0,
        text: prompt ?? _controller.text.trim(),
        images: List<String>.from(_images),
      );
      future = ref.read(chatProvider(chatId).notifier).create(message);
    }

    _controller.clear();
    _images.clear();
    _focusNode.unfocus();
    setState(() {
      _isEmpty = true;
      _inReplyTo = null;
      _message = null;
    });
    await future;
    setState(() => _isLoading = false);

    Future.delayed(const Duration(milliseconds: 100), () async {
      if (_scrollController.isAttached) {
        await _list.currentState?.scrollToLastIndex();
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    final color = ref.watch(configsProvider).valueOrNull!.getColor;
    final chatId = ref.watch(currentChatIdProvider);
    final asyncChat = ref.watch(chatProvider(chatId));
    final Chat chat;
    if (asyncChat.hasValue) {
      chat = asyncChat.value!;
    } else {
      if (asyncChat.hasError) {
        Crashlytics.log(asyncChat.error!, asyncChat.stackTrace);
      }
      chat = const Chat();
    }
    final messages = chat.messages;
    final isActive = ref.watch(chunkProvider(chatId)) != null;

    WidgetsBinding.instance.addPostFrameCallback((_) {
      final renderBox = _input.currentContext?.findRenderObject() as RenderBox?;
      if (renderBox != null && renderBox.size.height != _messageInputHeight) {
        setState(() => _messageInputHeight = renderBox.size.height);
      }
    });

    final limit = maxImages - _images.length;
    final prompts = [
      l10n.prompt1,
      l10n.prompt2,
      l10n.prompt3,
      l10n.prompt4,
      l10n.prompt5,
      l10n.prompt6,
      l10n.prompt7,
      l10n.prompt8,
      l10n.prompt9,
      l10n.prompt10,
    ];

    var padding = MediaQuery.paddingOf(context);
    padding = padding.copyWith(
      bottom: padding.bottom + MediaQuery.viewInsetsOf(context).bottom,
    );

    return GestureDetector(
      onTap: hideKeyboard,
      child: CupertinoPageScaffold(
        resizeToAvoidBottomInset: false,
        navigationBar: CupertinoNavigationBar(
          padding: const EdgeInsetsDirectional.symmetric(horizontal: 12.0),
          leading: CupertinoButton(
            padding: EdgeInsets.zero,
            sizeStyle: CupertinoButtonSize.small,
            onPressed: () {
              SettingsScreen.route(context);
            },
            child: Icon(MyIcons.settings, size: 24.0),
          ),
          middle: const Text(appName, style: Fonts.titleMedium),
          trailing: Row(
            spacing: 12.0,
            mainAxisSize: MainAxisSize.min,
            children: [
              CupertinoButton(
                padding: EdgeInsets.zero,
                sizeStyle: CupertinoButtonSize.small,
                onPressed: () {
                  ref.read(currentChatIdProvider.notifier).state = 0;
                  _focusNode.requestFocus();
                },
                child: Icon(MyIcons.plus, size: 24.0),
              ),
              CupertinoButton(
                padding: EdgeInsets.zero,
                sizeStyle: CupertinoButtonSize.small,
                onPressed: () async {
                  final result = await HistoryScreen.route(context);
                  if (result ?? false) {
                    Future.delayed(const Duration(milliseconds: 500), () async {
                      if (_scrollController.isAttached) {
                        await _list.currentState?.scrollToLastIndex();
                      }
                    });
                  }
                },
                child: Icon(MyIcons.history, size: 24.0),
              ),
            ],
          ),
          border: Border(
            bottom: BorderSide(color: separator(context), width: 0.4),
          ),
        ),
        child: Stack(
          fit: StackFit.expand,
          children: [
            if (messages.isNotEmpty)
              MessageList(
                key: _list,
                chatId: chatId,
                padding: padding.copyWith(
                  top: padding.top + 44.0,
                  bottom: _messageInputHeight ?? padding.bottom + 52.3,
                ),
                scrollController: _scrollController,
                onReply: (message) {
                  setState(() {
                    _inReplyTo = message;
                    _message = null;
                  });
                  _focusNode.requestFocus();
                  _list.currentState?.animate(message.id);
                },
                onEdit: (message) {
                  _controller.text = message.text;
                  _images.clear();
                  _images.addAll(message.images);
                  setState(() {
                    _isEmpty = false;
                    _inReplyTo = null;
                    _message = message;
                  });
                  _focusNode.requestFocus();
                  Future.delayed(const Duration(milliseconds: 100), () async {
                    if (_scrollController.isAttached) {
                      await _list.currentState?.scrollTo(message.id);
                    }
                  });
                },
                onImageTap: (images, index) {
                  ImageScreen.route(context, images, initialIndex: index);
                },
              )
            else
              Center(
                child: SingleChildScrollView(
                  padding: padding.copyWith(
                    left: padding.left + 30.0,
                    top: padding.top + 44.0,
                    right: padding.right + 30.0,
                    bottom: _messageInputHeight ?? padding.bottom + 52.4,
                  ),
                  child: SizedBox(
                    height: 300.0,
                    child: Column(
                      children: [
                        Logos.yordamchi(color: color, size: 100.0),
                        const SizedBox(height: 20.0),
                        Text(
                          l10n.howCanIHelpYou,
                          style: Fonts.titleMedium,
                          textAlign: TextAlign.center,
                        ),
                        const SizedBox(height: 20.0),
                        AnimatedSwitcher(
                          duration: const Duration(milliseconds: 500),
                          transitionBuilder: (child, animation) {
                            return FadeTransition(
                              opacity: animation,
                              child: SlideTransition(
                                position: Tween(
                                  begin: const Offset(0.0, 1.0),
                                  end: Offset.zero,
                                ).animate(animation),
                                child: child,
                              ),
                            );
                          },
                          child: ConstrainedBox(
                            key: ValueKey(_currentPromptIndex),
                            constraints: const BoxConstraints(
                              minHeight: buttonHeight,
                            ),
                            child: CupertinoButton.tinted(
                              padding: const EdgeInsets.all(12.0),
                              sizeStyle: CupertinoButtonSize.small,
                              onPressed: () {
                                _onSend(chatId, prompts[_currentPromptIndex]);
                              },
                              child: MyText(
                                prompts[_currentPromptIndex],
                                style: Fonts.bodyLarge,
                                textAlign: TextAlign.center,
                                isLinkBlue: false,
                                isSelectable: false,
                              ),
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
            Positioned(
              left: 0.0,
              right: 0.0,
              bottom: 0.0,
              child: MessageInput(
                key: _input,
                controller: _controller,
                placeholder: l10n.askAnything,
                focusNode: _focusNode,
                padding: padding,
                sendIcon: _isLoading || isActive
                    ? const CupertinoActivityIndicator()
                    : Icon(
                        _message != null
                            ? MyIcons.checkCircleFill
                            : MyIcons.send,
                        size: 28.0,
                      ),
                header: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    if (_inReplyTo != null) ...[
                      MessageInputHeader(
                        icon: MyIcons.reply,
                        text: _inReplyTo!.displayText,
                        onPressed: () {
                          _list.currentState?.scrollTo(_inReplyTo!.id);
                        },
                        onClear: () {
                          setState(() => _inReplyTo = null);
                        },
                      ),
                      const SizedBox(height: 8.0),
                    ],
                    if (_message != null) ...[
                      MessageInputHeader(
                        icon: MyIcons.edit,
                        text: _message!.displayText,
                        onPressed: () {
                          _list.currentState?.scrollTo(_message!.id);
                        },
                        onClear: () {
                          _controller.clear();
                          _images.clear();
                          setState(() {
                            _isEmpty = true;
                            _message = null;
                          });
                        },
                      ),
                      const SizedBox(height: 8.0),
                    ],
                    if (_images.isNotEmpty) ...[
                      ImageCarousel(
                        images: _images,
                        onTap: (index) async {
                          if (_images[index].isEmpty) return;
                          await ImageScreen.route(
                            context,
                            _images,
                            initialIndex: index,
                          );
                        },
                        onDelete: (index) {
                          setState(() {
                            _images.removeAt(index);
                            if (_images.isEmpty) {
                              _isEmpty = _controller.text.trim().isEmpty;
                            }
                          });
                        },
                        size: 100.0,
                      ),
                      const SizedBox(height: 8.0),
                    ],
                  ],
                ),
                onImageSelect: limit > 0
                    ? () async {
                        final result = await showActionSheet(
                          context: context,
                          icons: [CupertinoIcons.camera, CupertinoIcons.photo],
                          actions: [
                            l10n.camera,
                            Platform.isIOS ? l10n.photos : l10n.gallery,
                          ],
                        );
                        if (result == null) return;
                        final paths = await pickImage(
                          limit: limit,
                          source: result == 0
                              ? ImageSource.camera
                              : ImageSource.gallery,
                        );
                        if (paths == null || paths.isEmpty) return;
                        setState(() {
                          _images.addAll([for (final _ in paths) '']);
                          _isLoading = true;
                        });
                        final List<String> images;
                        try {
                          images = await Future.wait(paths.map(newImage));
                        } catch (error, stack) {
                          _onError(error, stack);
                          return;
                        }
                        final length = _images.length;
                        setState(() {
                          _images.removeRange(length - paths.length, length);
                          _images.addAll(images);
                          _isLoading = false;
                          _isEmpty = false;
                        });
                      }
                    : null,
                onChanged: (text) {
                  text = text.trim();
                  if (text.isEmpty != _isEmpty) {
                    setState(() {
                      _isEmpty = text.isEmpty;
                    });
                  }
                  if (text.characters.length > 5000) {
                    if (_error != l10n.errorTooLong) {
                      setState(() {
                        _error = l10n.errorTooLong;
                      });
                    }
                  } else if (_error != null) {
                    setState(() {
                      _error = null;
                    });
                  }
                },
                error: _error,
                onSend: _isLoading || _isEmpty || isActive || _error != null
                    ? null
                    : () => _onSend(chatId),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
