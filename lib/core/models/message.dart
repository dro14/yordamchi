import 'package:flutter/foundation.dart';

import 'package:yordamchi/core/repository/repository.dart';
import 'package:yordamchi/core/utils/utils.dart';
import 'call.dart';
import 'response.dart';

@immutable
class Message {
  final int id;
  final int userId;
  final int chatId;
  final String role;
  final int createdAt;
  final int inReplyTo;
  final String text;
  final List<String> images;
  final List<String> followUps;
  final List<Call> calls;
  final List<Response> responses;

  const Message({
    this.id = 0,
    this.userId = 0,
    this.chatId = 0,
    this.role = '',
    this.createdAt = 0,
    this.inReplyTo = 0,
    this.text = '',
    this.images = const [],
    this.followUps = const [],
    this.calls = const [],
    this.responses = const [],
  });

  String get displayText {
    if (text.isNotEmpty) {
      return text.replaceAll('\n', ' ');
    } else if (images.isNotEmpty) {
      return l10n.photo;
    } else if (followUps.isNotEmpty) {
      return followUps.map((e) => e.replaceAll('\n', ' ')).join(' ');
    } else {
      return '';
    }
  }

  factory Message.fromJson(Map<String, dynamic> json) {
    return Message(
      id: json['id'] ?? 0,
      userId: json['user_id'] ?? 0,
      chatId: json['chat_id'] ?? 0,
      role: json['role'] ?? '',
      createdAt: json['created_at'] ?? 0,
      inReplyTo: json['in_reply_to'] ?? 0,
      text: json['text'] ?? '',
      images: List<String>.from(decodeList(json['images'])),
      followUps: List<String>.from(decodeList(json['follow_ups'])),
      calls: List<Call>.from(
        decodeList(json['calls']).map((e) => Call.fromJson(e)),
      ),
      responses: List<Response>.from(
        decodeList(json['responses']).map((e) => Response.fromJson(e)),
      ),
    );
  }

  Message copyWith({
    int? id,
    int? userId,
    int? chatId,
    String? role,
    int? createdAt,
    int? inReplyTo,
    String? text,
    List<String>? images,
    List<String>? followUps,
    List<Call>? calls,
    List<Response>? responses,
  }) {
    return Message(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      chatId: chatId ?? this.chatId,
      role: role ?? this.role,
      createdAt: createdAt ?? this.createdAt,
      inReplyTo: inReplyTo ?? this.inReplyTo,
      text: text ?? this.text,
      images: images ?? this.images,
      followUps: followUps ?? this.followUps,
      calls: calls ?? this.calls,
      responses: responses ?? this.responses,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      if (id != 0) 'id': id,
      if (userId != 0) 'user_id': userId,
      if (chatId != 0) 'chat_id': chatId,
      if (role.isNotEmpty) 'role': role,
      if (createdAt != 0) 'created_at': createdAt,
      if (inReplyTo != 0) 'in_reply_to': inReplyTo,
      if (text.isNotEmpty) 'text': text,
      if (images.isNotEmpty) 'images': images,
      if (followUps.isNotEmpty) 'follow_ups': followUps,
      if (calls.isNotEmpty) 'calls': calls,
      if (responses.isNotEmpty) 'responses': responses,
    };
  }

  @override
  String toString() {
    var string = '';
    string += 'id: $id\n';
    string += 'userId: $userId\n';
    string += 'chatId: $chatId\n';
    string += 'role: $role\n';
    string += 'createdAt: $createdAt\n';
    string += 'inReplyTo: $inReplyTo\n';
    string += 'text: $text\n';
    string += 'images: $images\n';
    string += 'followUps: $followUps\n';
    string += 'calls: $calls\n';
    string += 'responses: $responses\n';
    return string;
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is Message &&
        other.runtimeType == runtimeType &&
        other.id == id &&
        other.userId == userId &&
        other.chatId == chatId &&
        other.role == role &&
        other.createdAt == createdAt &&
        other.inReplyTo == inReplyTo &&
        other.text == text &&
        listEquals(other.images, images) &&
        listEquals(other.followUps, followUps) &&
        listEquals(other.calls, calls) &&
        listEquals(other.responses, responses);
  }

  @override
  int get hashCode {
    return Object.hash(
      id,
      userId,
      chatId,
      role,
      createdAt,
      inReplyTo,
      text,
      images,
      followUps,
      calls,
      responses,
    );
  }
}
