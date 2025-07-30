import 'dart:io';
import 'dart:math';
import 'dart:convert';

import 'package:flutter/cupertino.dart';
import 'package:flutter/services.dart';
import 'package:flutter_timezone/flutter_timezone.dart';
import 'package:image_picker/image_picker.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:package_info_plus/package_info_plus.dart';
import 'package:intl/intl.dart' show DateFormat;

import 'package:yordamchi/core/repository/repository.dart';
import 'package:yordamchi/core/ui/ui.dart';
import 'constants.dart';

bool isSameDay(DateTime a, DateTime b) {
  return a.year == b.year && a.month == b.month && a.day == b.day;
}

bool isDayBefore(DateTime a, DateTime b) {
  b = b.subtract(const Duration(days: 1));
  return a.year == b.year && a.month == b.month && a.day == b.day;
}

bool isWithinYear(DateTime a, DateTime b) {
  final year = b.year - 1;
  final month = b.month;
  final day = min(b.day, DateTime(year, month + 1, 0).day);
  return a.isAfter(DateTime(year, month, day));
}

String formatDate(DateTime dateTime) {
  final now = DateTime.now();
  if (isSameDay(dateTime, now)) {
    return l10n.today;
  } else if (isDayBefore(dateTime, now)) {
    return l10n.yesterday;
  } else {
    final withYear = !isWithinYear(dateTime, now);
    final languageCode = l10n.localeName.split('_').first;
    final String pattern;
    switch (languageCode) {
      case 'uz':
        pattern = 'd-MMMM${withYear ? ' yyyy' : ''}';
      case 'ru':
        pattern = 'd MMMM${withYear ? ' yyyy' : ''}';
      case 'en':
        pattern = 'MMMM d${withYear ? ', yyyy' : ''}';
      default:
        throw Exception('Unsupported locale: $languageCode');
    }
    return DateFormat(pattern, l10n.localeName).format(dateTime);
  }
}

String formatImageUrl(String url) {
  return url.startsWith('http') ? url : 'https://$host/rasmlar/$url';
}

Future<void> hideKeyboard() async {
  FocusManager.instance.primaryFocus?.unfocus();
  await SystemChannels.textInput.invokeMethod('TextInput.hide');
}

Future<void> hapticFeedback() async {
  await HapticFeedback.mediumImpact();
}

Future<void> copyToClipboard(String text) async {
  await Clipboard.setData(ClipboardData(text: text));
}

String getBrightnessLabel(bool? isDark) {
  return isDark == null
      ? l10n.system
      : isDark == false
      ? l10n.light
      : l10n.dark;
}

double getTextWidth({required String text, required TextStyle style}) {
  final textPainter = TextPainter(
    textDirection: TextDirection.ltr,
    text: TextSpan(text: text, style: style),
  );
  textPainter.layout(maxWidth: double.infinity);
  return textPainter.width;
}

Future<void> initAppVersion() async {
  final packageInfo = await PackageInfo.fromPlatform();
  currentVersion = packageInfo.version;
}

Future<void> initTimezone() async {
  timezone = await FlutterTimezone.getLocalTimezone();
}

Future<List<String>?> pickImage({int limit = 1, ImageSource? source}) async {
  final picker = ImagePicker();
  if (limit > 1) {
    final images = await picker.pickMultiImage(limit: limit);
    if (images.isNotEmpty) {
      var paths = images.map((image) => image.path).toList();
      if (paths.length > limit) {
        paths = paths.sublist(0, limit);
      }
      return paths;
    } else {
      return null;
    }
  } else {
    final image = await picker.pickImage(source: source ?? ImageSource.gallery);
    if (image != null) {
      return [image.path];
    } else {
      return null;
    }
  }
}

Future<void> openTermsOfService() async {
  final url = 'https://terms.yordamchi.chuqurtech.uz/#${l10n.localeName}';
  await launchUrl(Uri.parse(url));
}

Future<void> openPrivacyPolicy() async {
  final url = 'https://privacy.yordamchi.chuqurtech.uz/#${l10n.localeName}';
  await launchUrl(Uri.parse(url));
}

Future<void> openContact() async {
  const url = 'https://t.me/yordamchi_dasturchisi';
  await launchUrl(Uri.parse(url));
}

Future<void> openUpdate() async {
  final url = Platform.isIOS
      ? 'https://apps.apple.com/us/app/yordamchi/id6747711302'
      : 'https://play.google.com/store/apps/details?id=uz.chuqurtech.yordamchi';
  await launchUrl(Uri.parse(url));
}

Widget? backButton(BuildContext context) {
  return Platform.isIOS
      ? null
      : CupertinoButton(
          padding: const EdgeInsets.only(left: 12.0),
          sizeStyle: CupertinoButtonSize.small,
          onPressed: () => Navigator.pop(context),
          child: Icon(MyIcons.back, size: 24.0),
        );
}

List<dynamic> decodeList(dynamic json) {
  return json is String ? jsonDecode(json) : json ?? [];
}
