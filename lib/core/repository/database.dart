import 'dart:convert';

import 'package:sqflite/sqflite.dart';

import 'package:yordamchi/core/models/models.dart';

// Version 1
const chats = '''
CREATE TABLE IF NOT EXISTS chats (
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL,
  name TEXT NOT NULL
);
''';
const idxChatList =
    'CREATE INDEX IF NOT EXISTS idx_chats_updated_at ON chats (updated_at DESC);';
const messages = '''
CREATE TABLE IF NOT EXISTS messages (
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL,
  chat_id INTEGER NOT NULL,
  role TEXT NOT NULL,
  created_at INTEGER NOT NULL,
  in_reply_to INTEGER,
  text TEXT,
  images TEXT,
  follow_ups TEXT
);
''';
const idxMessageList =
    'CREATE INDEX IF NOT EXISTS idx_messages_chat_id ON messages (chat_id);';

// Version 3
const callsColumn = 'ALTER TABLE messages ADD COLUMN calls TEXT;';
const responsesColumn = 'ALTER TABLE messages ADD COLUMN responses TEXT;';

late final Database _db;

Future<void> initDatabase() async {
  _db = await openDatabase(
    'yordamchi.db',
    version: 3,
    onUpgrade: (db, oldVersion, newVersion) async {
      if (oldVersion < 1) {
        await db.execute(chats);
        await db.execute(idxChatList);
        await db.execute(messages);
        await db.execute(idxMessageList);
      }
      if (oldVersion < 3) {
        await db.execute(callsColumn);
        await db.execute(responsesColumn);
      }
    },
  );
}

Future<List<int>> getChatIds() async {
  final results = await _db.query('chats', orderBy: 'updated_at DESC');
  return results.map((result) => result['id'] as int).toList();
}

Future<void> createChat(Chat chat) async {
  await _db.insert('chats', {
    'id': chat.id,
    'user_id': chat.userId,
    'created_at': chat.createdAt,
    'updated_at': chat.updatedAt,
    'name': chat.name,
  });
}

Future<Chat?> getChat(int chatId) async {
  var results = await _db.query('chats', where: 'id = ?', whereArgs: [chatId]);
  if (results.isEmpty) return null;
  final chat = Chat.fromJson(results.first);
  results = await _db.query(
    'messages',
    where: 'chat_id = ?',
    whereArgs: [chatId],
    orderBy: 'id',
  );
  final messages = results.map((result) => Message.fromJson(result)).toList();
  return chat.copyWith(messages: messages);
}

Future<void> updateChat(Chat chat) async {
  await _db.update(
    'chats',
    {'updated_at': chat.updatedAt, 'name': chat.name},
    where: 'id = ?',
    whereArgs: [chat.id],
  );
}

Future<void> deleteChats(List<int> chatIds) async {
  final batch = _db.batch();
  for (final chatId in chatIds) {
    batch.delete('chats', where: 'id = ?', whereArgs: [chatId]);
    batch.delete('messages', where: 'chat_id = ?', whereArgs: [chatId]);
  }
  await batch.commit();
}

Future<void> createMessage(Message message) async {
  await _db.insert('messages', {
    'id': message.id,
    'user_id': message.userId,
    'chat_id': message.chatId,
    'role': message.role,
    'created_at': message.createdAt,
    if (message.inReplyTo != 0) 'in_reply_to': message.inReplyTo,
    if (message.text.isNotEmpty) 'text': message.text,
    if (message.images.isNotEmpty) 'images': jsonEncode(message.images),
    if (message.followUps.isNotEmpty)
      'follow_ups': jsonEncode(message.followUps),
    if (message.calls.isNotEmpty) 'calls': jsonEncode(message.calls),
    if (message.responses.isNotEmpty)
      'responses': jsonEncode(message.responses),
  });
}

Future<void> deleteMessages(List<int> messageIds) async {
  if (messageIds.isEmpty) return;
  final batch = _db.batch();
  for (final messageId in messageIds) {
    batch.delete('messages', where: 'id = ?', whereArgs: [messageId]);
  }
  await batch.commit();
}
