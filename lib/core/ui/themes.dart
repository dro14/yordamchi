import 'dart:io';

import 'package:flutter/material.dart'
    show Colors, ColorScheme, DynamicSchemeVariant, Theme;
import 'package:flutter/cupertino.dart';

import 'package:yordamchi/core/repository/repository.dart';

const Color linkColor = CupertinoColors.link;
const Color white = Colors.white;
const Color black = Colors.black;
const Color transparent = Colors.transparent;
Color primary = transparent;
Color onPrimary = transparent;
Color primaryContainer = transparent;
Color onPrimaryContainer = transparent;
Color onInverseSurface = transparent;
Color scaffold = transparent;
Color surfaceContainer = transparent;
Color outline = transparent;
Color textColor = transparent;

Map<String, Color> colors = {
  'pink': Colors.pink,
  'red': Colors.red,
  'deepOrange': Colors.deepOrange,
  'orange': Colors.orange,
  'amber': Colors.amber,
  'yellow': Colors.yellow,
  'lime': Colors.lime,
  'lightGreen': Colors.lightGreen,
  'green': Colors.green,
  'teal': Colors.teal,
  'cyan': Colors.cyan,
  'lightBlue': Colors.lightBlue,
  'blue': Colors.blue,
  'indigo': Colors.indigo,
  'purple': Colors.purple,
  'deepPurple': Colors.deepPurple,
  'blueGrey': Colors.blueGrey,
  'brown': Colors.brown,
  'grey': Colors.grey,
};

String getColorName(int index) {
  final colorNames = [
    l10n.pink,
    l10n.red,
    l10n.deepOrange,
    l10n.orange,
    l10n.amber,
    l10n.yellow,
    l10n.lime,
    l10n.lightGreen,
    l10n.green,
    l10n.teal,
    l10n.cyan,
    l10n.lightBlue,
    l10n.blue,
    l10n.indigo,
    l10n.purple,
    l10n.deepPurple,
    l10n.blueGrey,
    l10n.brown,
    l10n.grey,
  ];
  return colorNames[index];
}

Color errorColor(BuildContext context) {
  return CupertinoDynamicColor.resolve(CupertinoColors.systemRed, context);
}

Color clearButton(BuildContext context) {
  return CupertinoDynamicColor.resolve(CupertinoColors.secondaryLabel, context);
}

Color textField(BuildContext context) {
  return CupertinoDynamicColor.resolve(
    CupertinoColors.tertiarySystemFill,
    context,
  );
}

Color separator(BuildContext context) {
  return CupertinoDynamicColor.resolve(
    const CupertinoDynamicColor.withBrightness(
      color: Color(0xFFA9A9AF),
      darkColor: Color(0xFF57585A),
    ),
    context,
  );
}

ColorScheme getColorScheme(Color color, Brightness brightness) {
  return ColorScheme.fromSeed(
    seedColor: color,
    brightness: brightness,
    dynamicSchemeVariant: color == Colors.grey
        ? DynamicSchemeVariant.monochrome
        : DynamicSchemeVariant.vibrant,
  );
}

Brightness brightness(BuildContext context) {
  return Platform.isIOS
      ? CupertinoTheme.of(context).brightness!
      : Theme.of(context).brightness;
}
