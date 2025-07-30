import 'package:flutter/foundation.dart';

class Response {
  final String id;
  final String name;
  final Map<String, dynamic> response;

  const Response({
    required this.id,
    required this.name,
    required this.response,
  });

  factory Response.fromJson(Map<String, dynamic> json) {
    return Response(
      id: json['id'] ?? '',
      name: json['name'] ?? '',
      response: json['response'] ?? {},
    );
  }

  Response copyWith({
    String? id,
    String? name,
    Map<String, dynamic>? response,
  }) {
    return Response(
      id: id ?? this.id,
      name: name ?? this.name,
      response: response ?? this.response,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      if (id.isNotEmpty) 'id': id,
      if (name.isNotEmpty) 'name': name,
      if (response.isNotEmpty) 'response': response,
    };
  }

  @override
  String toString() {
    var string = '';
    string += 'id: $id\n';
    string += 'name: $name\n';
    string += 'response: $response\n';
    return string;
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is Response &&
        other.runtimeType == runtimeType &&
        other.id == id &&
        other.name == name &&
        mapEquals(other.response, response);
  }

  @override
  int get hashCode {
    return Object.hash(id, name, response);
  }
}
