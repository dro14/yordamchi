import 'package:flutter/material.dart' show SelectionArea;
import 'package:flutter/cupertino.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_math_fork/flutter_math.dart';
import 'package:flutter_highlight/flutter_highlight.dart';
import 'package:flutter_highlight/theme_map.dart';
import 'package:gpt_markdown/gpt_markdown.dart';
import 'package:gpt_markdown/custom_widgets/markdown_config.dart';
import 'package:gpt_markdown/custom_widgets/selectable_adapter.dart';
import 'package:url_launcher/url_launcher.dart';

import 'package:yordamchi/core/repository/repository.dart';
import 'package:yordamchi/core/utils/utils.dart';
import 'package:yordamchi/core/ui/ui.dart';

final List<MarkdownComponent> components = [
  CodeBlockMd(),
  NewLines(),
  BlockQuote(),
  ImageMd(),
  ATagMd(),
  TableMd(),
  HTag(),
  UnOrderedList(),
  OrderedList(),
  RadioButtonMd(),
  CheckBoxMd(),
  HrLine(),
  StrikeMd(),
  BoldMd(),
  ItalicMd(),
  LatexMath(),
  LatexMathMultiLine(),
  HighlightedText(),
  SourceTag(),
  IndentMd(),
];

final List<MarkdownComponent> inlineComponents = [
  ImageMd(),
  ATagMd(),
  TableMd(),
  StrikeMd(),
  BoldMd(),
  ItalicMd(),
  LatexMath(),
  LatexMathMultiLine(),
  HighlightedText(),
  SourceTag(),
];

final _linkRegExp = RegExp(
  r'(\[.*?\]\(.*?\))'
  r'|(?<![\w`*_])'
  r'((?:https?:\/\/)'
  r'(?:www\.)?'
  r'(?:(?:[-a-z0-9]{2,63}\.)+[-a-z0-9]{2,63}'
  r'|localhost'
  r'|(?:\d{1,3}\.){3}\d{1,3})'
  r"[-a-z0-9._~:/?#\[\]@!$&'()*+,;=%]*)",
  caseSensitive: false,
);

class MyText extends StatelessWidget {
  final String data;
  final TextStyle style;
  final TextOverflow? overflow;
  final TextAlign? textAlign;
  final bool isLinkBlue;
  final bool isSelectable;

  const MyText(
    this.data, {
    required this.style,
    this.overflow,
    this.textAlign,
    this.isLinkBlue = true,
    this.isSelectable = true,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    final child = NotificationListener<ScrollUpdateNotification>(
      onNotification: (notification) => true,
      child: GptMarkdown(
        data.replaceAllMapped(_linkRegExp, (match) {
          final s = match.group(0)!;
          return match.group(1) != null ? s : '[$s]($s)';
        }),
        linkBuilder: (context, title, url, style) {
          return GestureDetector(
            onTap: () => launchUrl(Uri.parse(url)),
            child: Text.rich(
              title,
              overflow: TextOverflow.ellipsis,
              style: style.copyWith(
                color: isLinkBlue ? linkColor : null,
                decoration: TextDecoration.underline,
              ),
            ),
          );
        },
        latexBuilder: (context, text, textStyle, isInline) {
          return Latex(text: text, textStyle: textStyle);
        },
        codeBuilder: (context, name, code, closed) {
          return Code(name: name, code: code);
        },
        style: style,
        overflow: overflow,
        textAlign: textAlign,
        components: components,
        inlineComponents: inlineComponents,
        useDollarSignsForLatex: true,
      ),
    );

    return isSelectable ? SelectionArea(child: child) : child;
  }
}

class Latex extends StatefulWidget {
  final String text;
  final TextStyle textStyle;

  const Latex({required this.text, required this.textStyle, super.key});

  @override
  State<Latex> createState() => _LatexState();
}

class _LatexState extends State<Latex> {
  static final _isOverflowing = <String, bool>{};
  final _key = GlobalKey();

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        final child = Padding(
          key: _key,
          padding: const EdgeInsets.symmetric(vertical: 6.0),
          child: SelectableAdapter(
            selectedText: widget.text,
            child: Math.tex(
              widget.text,
              textScaleFactor: 1,
              textStyle: widget.textStyle,
              mathStyle: MathStyle.display,
              settings: const TexParserSettings(strict: Strict.ignore),
              options: MathOptions(
                style: MathStyle.display,
                sizeUnderTextStyle: MathSize.large,
                color: widget.textStyle.color ?? textColor,
                fontSize: widget.textStyle.fontSize ?? 16.0,
                mathFontOptions: FontOptions(
                  fontFamily: 'Main',
                  fontWeight: widget.textStyle.fontWeight ?? FontWeight.normal,
                  fontShape: FontStyle.normal,
                ),
                textFontOptions: FontOptions(
                  fontFamily: 'Main',
                  fontWeight: widget.textStyle.fontWeight ?? FontWeight.normal,
                  fontShape: FontStyle.normal,
                ),
              ),
              onErrorFallback: (error) {
                return Text(
                  widget.text,
                  textDirection: TextDirection.ltr,
                  style: widget.textStyle.copyWith(
                    color: (kDebugMode) ? errorColor(context) : null,
                  ),
                );
              },
            ),
          ),
        );

        final id = '${constraints.maxWidth}_${widget.text}';

        if (!_isOverflowing.containsKey(id)) {
          WidgetsBinding.instance.addPostFrameCallback((_) {
            final box = _key.currentContext?.findRenderObject() as RenderBox?;
            if (box != null) {
              setState(() {
                _isOverflowing[id] = box.size.width > constraints.maxWidth;
              });
            }
          });
          return SingleChildScrollView(
            scrollDirection: Axis.horizontal,
            child: child,
          );
        }

        return _isOverflowing[id]!
            ? Container(
                margin: const EdgeInsets.symmetric(vertical: 6.0),
                decoration: BoxDecoration(
                  color: onInverseSurface,
                  borderRadius: BorderRadius.circular(8.0),
                ),
                child: SingleChildScrollView(
                  scrollDirection: Axis.horizontal,
                  child: Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 8.0),
                    child: child,
                  ),
                ),
              )
            : child;
      },
    );
  }
}

class Code extends StatefulWidget {
  final String name;
  final String code;

  const Code({required this.name, required this.code, super.key});

  @override
  State<Code> createState() => _CodeState();
}

class _CodeState extends State<Code> {
  var _copied = false;

  void _copy() async {
    hapticFeedback();
    await copyToClipboard(widget.code);
    setState(() {
      _copied = true;
    });
    await Future.delayed(const Duration(seconds: 2));
    setState(() {
      _copied = false;
    });
  }

  Widget get _icon {
    return Icon(
      _copied ? MyIcons.check : MyIcons.copy,
      color: textColor,
      size: 15,
    );
  }

  Widget get _label {
    return Text(
      _copied ? l10n.copied : l10n.copy,
      style: Fonts.bodyMedium.copyWith(color: textColor),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Map<String, TextStyle>.from(
      brightness(context) == Brightness.light
          ? themeMap['vs']!
          : themeMap['vs2015']!,
    );
    theme['root'] = TextStyle(
      color: theme['root']!.color,
      backgroundColor: onInverseSurface,
    );

    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8.0),
      child: ClipRRect(
        borderRadius: BorderRadius.circular(8.0),
        clipBehavior: Clip.hardEdge,
        child: Container(
          color: onInverseSurface,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 12.0,
                      vertical: 4,
                    ),
                    child: Text(widget.name, style: Fonts.bodyMedium),
                  ),
                  const Spacer(),
                  CupertinoButton(
                    padding: EdgeInsets.zero,
                    sizeStyle: CupertinoButtonSize.small,
                    onPressed: _copy,
                    child: Row(
                      children: [
                        const SizedBox(width: 8.0),
                        _icon,
                        const SizedBox(width: 4.0),
                        _label,
                        const SizedBox(width: 12.0),
                      ],
                    ),
                  ),
                ],
              ),
              Container(color: separator(context), height: 1.0),
              SingleChildScrollView(
                scrollDirection: Axis.horizontal,
                child: HighlightView(
                  widget.code,
                  language: widget.name,
                  theme: theme,
                  padding: const EdgeInsets.all(12.0),
                  textStyle: const TextStyle(
                    fontFamily: 'RobotoMono',
                    fontSize: 14.0,
                    fontWeight: FontWeight.w400,
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class CheckBoxMd extends BlockMd {
  @override
  String get expString => r'\[((?:\X|\x|\ ))\]\ (\S[^\n]*?)$';
  dynamic get onLinkTab => null;

  @override
  Widget build(BuildContext context, String text, GptMarkdownConfig config) {
    final match = exp.firstMatch(text.trim());
    return CustomCb(
      value: '${match?[1]}'.toLowerCase() == 'x',
      textDirection: config.textDirection,
      child: MdWidget(context, '${match?[2]}', false, config: config),
    );
  }
}

class CustomCb extends StatelessWidget {
  final Widget child;
  final bool value;
  final double spacing;
  final TextDirection textDirection;

  const CustomCb({
    super.key,
    this.spacing = 5,
    required this.child,
    this.textDirection = TextDirection.ltr,
    required this.value,
  });

  @override
  Widget build(BuildContext context) {
    return Directionality(
      textDirection: textDirection,
      child: Row(
        mainAxisSize: MainAxisSize.min,
        textBaseline: TextBaseline.alphabetic,
        crossAxisAlignment: CrossAxisAlignment.baseline,
        children: [
          Text.rich(
            WidgetSpan(
              alignment: PlaceholderAlignment.middle,
              child: Padding(
                padding: EdgeInsetsDirectional.only(
                  start: spacing,
                  end: spacing,
                ),
                child: CupertinoCheckbox(value: value, onChanged: (value) {}),
              ),
            ),
          ),
          Flexible(child: child),
        ],
      ),
    );
  }
}

class RadioButtonMd extends BlockMd {
  @override
  String get expString => r'\(((?:\X|\x|\ ))\)\ (\S[^\n]*)$';
  dynamic get onLinkTab => null;

  @override
  Widget build(BuildContext context, String text, GptMarkdownConfig config) {
    final match = exp.firstMatch(text.trim());
    return CustomRb(
      value: '${match?[1]}'.toLowerCase() == 'x',
      textDirection: config.textDirection,
      child: MdWidget(context, '${match?[2]}', false, config: config),
    );
  }
}

class CustomRb extends StatelessWidget {
  const CustomRb({
    super.key,
    this.spacing = 5,
    required this.child,
    this.textDirection = TextDirection.ltr,
    required this.value,
  });
  final Widget child;
  final bool value;
  final double spacing;
  final TextDirection textDirection;

  @override
  Widget build(BuildContext context) {
    return Directionality(
      textDirection: textDirection,
      child: Row(
        mainAxisSize: MainAxisSize.min,
        textBaseline: TextBaseline.alphabetic,
        crossAxisAlignment: CrossAxisAlignment.baseline,
        children: [
          Text.rich(
            WidgetSpan(
              alignment: PlaceholderAlignment.middle,
              child: Padding(
                padding: EdgeInsetsDirectional.only(
                  start: spacing,
                  end: spacing,
                ),
                child: CupertinoRadio(
                  value: value,
                  groupValue: true,
                  onChanged: (value) {},
                ),
              ),
            ),
          ),
          Flexible(child: child),
        ],
      ),
    );
  }
}

class TableMd extends BlockMd {
  @override
  String get expString =>
      r'(((\|[^\n\|]+\|)((([^\n\|]+\|)+)?)\ *)(\n\ *(((\|[^\n\|]+\|)(([^\n\|]+\|)+)?))\ *)+)$';
  @override
  Widget build(BuildContext context, String text, GptMarkdownConfig config) {
    final List<Map<int, String>> value = text
        .split('\n')
        .map<Map<int, String>>(
          (e) => e
              .trim()
              .split('|')
              .where((element) => element.isNotEmpty)
              .toList()
              .asMap(),
        )
        .toList();
    bool heading = RegExp(
      r'^\|.*?\|\n\|-[-\\ |]*?-\|$',
      multiLine: true,
    ).hasMatch(text.trim());
    int maxCol = 0;
    for (final each in value) {
      if (maxCol < each.keys.length) {
        maxCol = each.keys.length;
      }
    }
    if (maxCol == 0) {
      return Text('', style: config.style);
    }
    return SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      child: Container(
        color: onInverseSurface,
        child: Table(
          textDirection: config.textDirection,
          defaultColumnWidth: CustomTableColumnWidth(),
          defaultVerticalAlignment: TableCellVerticalAlignment.middle,
          border: TableBorder.all(width: 1, color: textColor),
          children: value
              .asMap()
              .entries
              .map<TableRow>(
                (entry) => TableRow(
                  decoration: (heading)
                      ? BoxDecoration(
                          color: (entry.key == 0) ? primaryContainer : null,
                        )
                      : null,
                  children: List.generate(maxCol, (index) {
                    var e = entry.value;
                    String data = e[index] ?? '';
                    if (RegExp(r'^:?--+:?$').hasMatch(data.trim()) ||
                        data.trim().isEmpty) {
                      return const SizedBox();
                    }

                    return Center(
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 8,
                          vertical: 4,
                        ),
                        child: MdWidget(
                          context,
                          (e[index] ?? '').trim(),
                          false,
                          config: config,
                        ),
                      ),
                    );
                  }),
                ),
              )
              .toList(),
        ),
      ),
    );
  }
}

class ItalicMd extends InlineMd {
  @override
  RegExp get exp =>
      RegExp(r'(?:(?<!\*)\*(?<!\s)(.+?)(?<!\s)\*(?!\*))', dotAll: true);

  @override
  InlineSpan span(BuildContext context, String text, GptMarkdownConfig config) {
    var match = exp.firstMatch(text.trim());
    var data = match?[1] ?? match?[2];
    var conf = config.copyWith(
      style: (config.style ?? const TextStyle()).copyWith(
        fontStyle: FontStyle.italic,
        fontFamily: 'Roboto',
      ),
    );
    return TextSpan(
      children: MarkdownComponent.generate(context, '$data', conf, false),
      style: conf.style,
    );
  }
}

class SourceTag extends InlineMd {
  @override
  RegExp get exp => RegExp(r'(?:【.*?)?\[(\d+?)\]');

  @override
  InlineSpan span(BuildContext context, String text, GptMarkdownConfig config) {
    var match = exp.firstMatch(text.trim());
    var content = match?[1];
    if (content == null) {
      return const TextSpan();
    }
    return WidgetSpan(
      alignment: PlaceholderAlignment.middle,
      child: Padding(
        padding: const EdgeInsets.all(2),
        child:
            config.sourceTagBuilder?.call(
              context,
              content,
              const TextStyle(),
            ) ??
            SizedBox(
              width: 20.0,
              height: 20.0,
              child: Container(
                decoration: BoxDecoration(
                  color: onInverseSurface,
                  shape: BoxShape.circle,
                ),
                child: FittedBox(
                  fit: BoxFit.scaleDown,
                  child: Text(
                    content,
                    style: (config.style ?? const TextStyle()).copyWith(),
                    textDirection: config.textDirection,
                  ),
                ),
              ),
            ),
      ),
    );
  }
}
