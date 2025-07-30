import 'dart:io';
import 'dart:async';
import 'dart:convert';

import 'package:http/http.dart' as http;
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:yordamchi/core/models/models.dart';
import 'package:yordamchi/core/utils/utils.dart';
import 'calls.dart';
import 'configs.dart';
import 'database.dart';

final currentChatIdProvider = StateProvider<int>((ref) => 0);

final chunkProvider = StateProviderFamily<String?, int>((ref, chatId) => null);

final chatsProvider = FutureProvider<List<Chat>>((ref) async {
  final chatIds = await getChatIds();
  final chats = <Chat>[];
  for (final chatId in chatIds) {
    final chat = await ref.read(chatProvider(chatId).future);
    final messages = <Message>[];
    for (final message in chat.messages.reversed) {
      if (message.displayText.isNotEmpty) messages.insert(0, message);
      if (messages.length == 2) break;
    }
    chats.add(chat.copyWith(messages: messages));
  }
  return chats;
});

final chatProvider = AsyncNotifierProviderFamily<ChatNotifier, Chat, int>(
  () => ChatNotifier(),
);

class ChatNotifier extends FamilyAsyncNotifier<Chat, int> {
  final _map = <int, Message>{};
  var _event = '';

  @override
  Future<Chat> build(int chatId) async {
    if (chatId == 0) {
      return const Chat();
    }
    final chat = (await getChat(chatId))!;
    for (final message in chat.messages) {
      _map[message.id] = message;
    }
    return chat;
  }

  Future<void> rename(String name) async {
    var chat = state.valueOrNull!;
    chat = chat.copyWith(name: name);
    state = AsyncData(chat);
    ref.invalidate(chatsProvider);
    await renameChat(chat);
  }

  Message? get(int id) {
    return _map[id];
  }

  Future<void> create(Message message) async {
    _addMessage(message);
    try {
      await _newMessage('POST', message);
    } on http.ClientException {
      _addMessage(_errorMessage(withInReplyTo: true, noInternet: true));
    } on SocketException {
      _addMessage(_errorMessage(withInReplyTo: true, noInternet: true));
    } catch (error, stack) {
      _addMessage(_errorMessage(withInReplyTo: true));
      Crashlytics.log(error, stack);
    }
  }

  Future<void> edit(Message message) async {
    _deleteUntil(message);
    try {
      await _newMessage('PUT', message);
    } on http.ClientException {
      _addMessage(_errorMessage(withInReplyTo: true, noInternet: true));
    } on SocketException {
      _addMessage(_errorMessage(withInReplyTo: true, noInternet: true));
    } catch (error, stack) {
      _addMessage(_errorMessage(withInReplyTo: true));
      Crashlytics.log(error, stack);
    }
  }

  Future<void> followUps(Message message) async {
    final chat = state.valueOrNull!;
    if (ref.read(chunkProvider(chat.id)) != '') {
      ref.read(chunkProvider(chat.id).notifier).state = '';
    }
    final contents = <Message>[];
    var count = 0;
    for (final message in chat.messages.reversed) {
      if (message.followUps.isEmpty) contents.insert(0, message);
      if (message.text.isNotEmpty || message.images.isNotEmpty) count++;
      if (count == 4) break;
    }
    Message response;
    try {
      response = await newFollowUps({
        'user_id': chat.userId,
        'chat_id': chat.id,
        'language': l10n.localeName,
        'contents': contents,
      });
    } on http.ClientException {
      response = _errorMessage(noInternet: true);
    } on SocketException {
      response = _errorMessage(noInternet: true);
    } catch (error, stack) {
      response = _errorMessage();
      Crashlytics.log(error, stack);
    }
    _addMessage(response);
    ref.read(chunkProvider(chat.id).notifier).state = null;
  }

  Message _errorMessage({bool withInReplyTo = false, bool noInternet = false}) {
    return Message(
      id: DateTime.now().millisecondsSinceEpoch,
      userId: state.valueOrNull!.userId,
      chatId: state.valueOrNull!.id,
      role: 'model',
      createdAt: DateTime.now().millisecondsSinceEpoch,
      inReplyTo: withInReplyTo ? state.valueOrNull!.messages.last.id : 0,
      text: noInternet ? l10n.noInternet : l10n.errorMessage,
    );
  }

  Future<void> _newMessage(String method, Message message) async {
    final chat = state.valueOrNull!;
    ref.read(chunkProvider(chat.id).notifier).state = '';
    final request = http.Request(method, Uri.https(host, '/v1/message'));
    request.headers['Authorization'] = authorization;
    request.body = jsonEncode({
      'user_id': chat.userId,
      'chat_id': chat.id,
      'language': l10n.localeName,
      'contents': _getContents(message),
    });
    final response = await client.send(request);
    response.stream
        .transform(utf8.decoder)
        .transform(const LineSplitter())
        .listen((line) async {
          if (line.startsWith('event:')) {
            _event = line.substring(6);
          } else if (line.startsWith('data:')) {
            final data = line.substring(5);
            switch (_event) {
              case 'request':
                final request = Message.fromJson(jsonDecode(data));
                _addMessage(request, replace: true);
              case 'chunk':
                final chunk = jsonDecode(data)['chunk'];
                ref.read(chunkProvider(chat.id).notifier).state = chunk;
              case 'response':
                final response = Message.fromJson(jsonDecode(data));
                _addMessage(response);
                if (response.calls.isEmpty) {
                  ref.read(chunkProvider(chat.id).notifier).state = null;
                } else {
                  final responses = <Response>[];
                  for (final call in response.calls) {
                    if (call.name == 'web_search' ||
                        call.name == 'google_search') {
                      final result = await getSearchResults(call.args['query']);
                      responses.add(
                        Response(
                          id: call.id,
                          name: call.name,
                          response: {'result': result},
                        ),
                      );
                    }
                  }
                  if (responses.isNotEmpty) {
                    final message = Message(
                      id: DateTime.now().millisecondsSinceEpoch,
                      userId: chat.userId,
                      chatId: chat.id,
                      role: 'user',
                      createdAt: DateTime.now().millisecondsSinceEpoch,
                      responses: responses,
                    );
                    _addMessage(message);
                    try {
                      await _newMessage('POST', message);
                    } on http.ClientException {
                      _addMessage(
                        _errorMessage(withInReplyTo: true, noInternet: true),
                      );
                    } on SocketException {
                      _addMessage(
                        _errorMessage(withInReplyTo: true, noInternet: true),
                      );
                    } catch (error, stack) {
                      _addMessage(_errorMessage(withInReplyTo: true));
                      Crashlytics.log(error, stack);
                    }
                  }
                }
              case 'error':
                _addMessage(_errorMessage());
                ref.read(chunkProvider(chat.id).notifier).state = null;
                Crashlytics.log(data, StackTrace.current);
            }
          }
        });
  }

  void _addMessage(Message message, {bool replace = false}) {
    var chat = state.valueOrNull!;
    final messages = List<Message>.from(chat.messages);
    if (replace) {
      final last = messages.removeLast();
      deleteMessages([last.id]);
      _map.remove(last.id);
    }
    messages.add(message);
    createMessage(message);
    _map[message.id] = message;
    chat = chat.copyWith(messages: messages, updatedAt: message.createdAt);
    state = AsyncData(chat);
    updateChat(chat);
    ref.invalidate(chatsProvider);
  }

  void _deleteUntil(Message message) {
    var chat = state.valueOrNull!;
    final messages = List<Message>.from(chat.messages);
    final deleteIds = <int>[];
    for (int i = messages.length - 1; i >= 0; i--) {
      deleteIds.add(messages[i].id);
      if (messages[i].id == message.id) break;
    }
    if (deleteIds.isNotEmpty) {
      messages.removeWhere((message) => deleteIds.contains(message.id));
      deleteMessages(deleteIds);
      _map.removeWhere((messageId, message) => deleteIds.contains(messageId));
    }
    messages.add(message);
    createMessage(message);
    _map[message.id] = message;
    chat = chat.copyWith(messages: messages, updatedAt: message.createdAt);
    state = AsyncData(chat);
    updateChat(chat);
    ref.invalidate(chatsProvider);
  }

  List<Message> _getContents(Message message) {
    final contents = <Message>[];
    if (message.inReplyTo == 0) {
      var count = 0;
      for (final message in state.valueOrNull!.messages.reversed) {
        if (message.followUps.isEmpty) contents.insert(0, message);
        if (message.text.isNotEmpty || message.images.isNotEmpty) count++;
        if (count == 3) break;
      }
    } else {
      final response = get(message.inReplyTo);
      if (response != null) {
        final request = get(response.inReplyTo);
        if (request != null) {
          contents.add(request);
        }
        contents.add(response);
      }
      contents.add(message);
    }
    return contents;
  }
}
