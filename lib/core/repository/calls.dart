import 'dart:io';
import 'dart:convert';
import 'dart:typed_data';

import 'package:http/http.dart' as http;
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:flutter_image_compress/flutter_image_compress.dart';

import 'package:yordamchi/core/models/models.dart';
import 'package:yordamchi/core/utils/utils.dart';
import 'configs.dart';
import 'database.dart';

final client = http.Client();
final host = dotenv.env['ENDPOINT']!;
final _apiKey = dotenv.env['API_KEY']!;
final _encoded = base64Encode(utf8.encode('foydalanuvchi:$_apiKey'));
final authorization = 'Basic $_encoded';
final headers = {'Authorization': authorization};

Future<int> newUser() async {
  final json = await _call('POST', 'user');
  return json['id'];
}

Future<int> newChat() async {
  final json = await _call('POST', 'chat', body: {'user_id': id});
  final chat = Chat.fromJson(json);
  await createChat(chat);
  return chat.id;
}

Future<void> renameChat(Chat chat) async {
  await updateChat(chat);
  await _call('PATCH', 'chat', body: chat.toJson());
}

Future<void> clearChats(List<int> chatIds) async {
  await deleteChats(chatIds);
  await _call('DELETE', 'chat', body: {'chat_ids': chatIds});
}

Future<Message> newFollowUps(Map<String, dynamic> body) async {
  final json = await _call('POST', 'follow-up', body: body);
  return Message.fromJson(json);
}

Future<String> newImage(String image) async {
  if (!await File(image).exists()) return image;
  final response = await client.post(
    Uri.https(host, '/v1/image'),
    headers: {
      'Authorization': authorization,
      'Content-Type': 'application/octet-stream',
    },
    body: await _compressImage(image),
  );
  _throwIfError(response);
  return response.headers['x-filename'] ?? '';
}

Future<bool?> haveToUpdate() async {
  final json = await _call('GET', 'version');
  minimumVersion = json['minimum'];
  final current = currentVersion.split('.').map(int.parse).toList();
  final minimum = minimumVersion.split('.').map(int.parse).toList();
  for (int i = 0; i < 3; i++) {
    if (current[i] > minimum[i]) break;
    if (current[i] < minimum[i]) return true;
  }
  if (latestVersion == json['latest']) return null;
  latestVersion = json['latest'];
  final latest = latestVersion.split('.').map(int.parse).toList();
  for (int i = 0; i < 3; i++) {
    if (current[i] > latest[i]) break;
    if (current[i] < latest[i]) return false;
  }
  return null;
}

Future<dynamic> _call(String method, String path, {Object? body}) async {
  final request = http.Request(method, Uri.https(host, '/v1/$path'));
  request.headers['Authorization'] = authorization;
  if (body != null) {
    request.headers['Content-Type'] = 'application/json';
    request.body = jsonEncode(body);
  }
  final response = await http.Response.fromStream(await client.send(request));
  _throwIfError(response);
  return response.body.isNotEmpty ? jsonDecode(response.body) : null;
}

void _throwIfError(http.Response response) {
  if (response.statusCode != 200) {
    var message = '${response.statusCode}';
    try {
      final data = jsonDecode(response.body);
      message += ': ${data['error']}';
    } catch (error) {
      if (response.body.isNotEmpty) {
        message += ': ${response.body}';
      }
    }
    throw Exception(message);
  }
}

Future<Uint8List?> _compressImage(String image) {
  return FlutterImageCompress.compressWithFile(
    image,
    minWidth: 1080,
    minHeight: 1080,
    quality: 80,
  );
}
