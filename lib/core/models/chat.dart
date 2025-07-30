import 'package:flutter/foundation.dart';

import 'message.dart';

@immutable
class Chat {
  final int id;
  final int userId;
  final int createdAt;
  final int updatedAt;
  final String name;
  final List<Message> messages;

  const Chat({
    this.id = 0,
    this.userId = 0,
    this.createdAt = 0,
    this.updatedAt = 0,
    this.name = '',
    this.messages = const [],
  });

  DateTime get dateTime => DateTime.fromMillisecondsSinceEpoch(updatedAt);

  factory Chat.fromJson(Map<String, dynamic> json) {
    return Chat(
      id: json['id'] ?? 0,
      userId: json['user_id'] ?? 0,
      createdAt: json['created_at'] ?? 0,
      updatedAt: json['updated_at'] ?? 0,
      name: json['name'] ?? '',
      messages: List<Message>.from(json['messages'] ?? []),
    );
  }

  Chat copyWith({
    int? id,
    int? userId,
    int? createdAt,
    int? updatedAt,
    String? name,
    List<Message>? messages,
  }) {
    return Chat(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
      name: name ?? this.name,
      messages: messages ?? this.messages,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      if (id != 0) 'id': id,
      if (userId != 0) 'user_id': userId,
      if (createdAt != 0) 'created_at': createdAt,
      if (updatedAt != 0) 'updated_at': updatedAt,
      if (name != '') 'name': name,
    };
  }

  @override
  String toString() {
    var string = '';
    string += 'id: $id\n';
    string += 'userId: $userId\n';
    string += 'createdAt: $createdAt\n';
    string += 'updatedAt: $updatedAt\n';
    string += 'name: $name\n';
    string += 'messages: $messages\n';
    return string;
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is Chat &&
        other.runtimeType == runtimeType &&
        other.id == id &&
        other.userId == userId &&
        other.createdAt == createdAt &&
        other.updatedAt == updatedAt &&
        other.name == name &&
        listEquals(other.messages, messages);
  }

  @override
  int get hashCode {
    return Object.hash(id, userId, createdAt, updatedAt, name, messages);
  }
}
