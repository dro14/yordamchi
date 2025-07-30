import 'package:flutter/foundation.dart';

class Call {
  final String id;
  final String name;
  final Map<String, dynamic> args;

  const Call({required this.id, required this.name, required this.args});

  factory Call.fromJson(Map<String, dynamic> json) {
    return Call(
      id: json['id'] ?? '',
      name: json['name'] ?? '',
      args: json['args'] ?? {},
    );
  }

  Call copyWith({String? id, String? name, Map<String, dynamic>? args}) {
    return Call(
      id: id ?? this.id,
      name: name ?? this.name,
      args: args ?? this.args,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      if (id.isNotEmpty) 'id': id,
      if (name.isNotEmpty) 'name': name,
      if (args.isNotEmpty) 'args': args,
    };
  }

  @override
  String toString() {
    var string = '';
    string += 'id: $id\n';
    string += 'name: $name\n';
    string += 'args: $args\n';
    return string;
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is Call &&
        other.runtimeType == runtimeType &&
        other.id == id &&
        other.name == name &&
        mapEquals(other.args, args);
  }

  @override
  int get hashCode {
    return Object.hash(id, name, args);
  }
}
