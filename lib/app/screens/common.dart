import 'dart:ui';

import 'package:flutter/cupertino.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:yordamchi/core/core.dart';

class MyScaffold extends StatefulWidget {
  final ObstructingPreferredSizeWidget? navigationBar;
  final Widget child;

  const MyScaffold({this.navigationBar, required this.child, super.key});

  @override
  State<MyScaffold> createState() => _MyScaffoldState();
}

class _MyScaffoldState extends State<MyScaffold> {
  double _dragOffset = 0.0;

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: hideKeyboard,
      behavior: HitTestBehavior.opaque,
      onHorizontalDragUpdate: (details) {
        _dragOffset += details.delta.dx;
      },
      onHorizontalDragEnd: (details) {
        if (_dragOffset > 60.0 && Navigator.canPop(context)) {
          Navigator.pop(context);
        } else {
          _dragOffset = 0.0;
        }
      },
      child: CupertinoPageScaffold(
        navigationBar: widget.navigationBar,
        child: widget.child,
      ),
    );
  }
}

class MessageInput extends ConsumerWidget {
  final TextEditingController controller;
  final String placeholder;
  final FocusNode focusNode;
  final EdgeInsets padding;
  final Widget sendIcon;
  final List<String>? images;
  final Widget? header;
  final void Function(int)? onImageTap;
  final VoidCallback? onImageSelect;
  final void Function(String)? onChanged;
  final String? error;
  final VoidCallback? onSend;

  const MessageInput({
    required this.controller,
    required this.placeholder,
    required this.focusNode,
    required this.padding,
    required this.sendIcon,
    this.header,
    this.images,
    this.onImageTap,
    required this.onImageSelect,
    this.onChanged,
    this.error,
    this.onSend,
    super.key,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return ClipRect(
      child: BackdropFilter(
        filter: ImageFilter.blur(sigmaX: 10.0, sigmaY: 10.0),
        child: Container(
          decoration: BoxDecoration(
            border: Border(
              top: BorderSide(color: separator(context), width: 0.4),
            ),
            color: surfaceContainer.withAlpha(200),
          ),
          padding: EdgeInsets.only(top: 8.0, bottom: padding.bottom + 8.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              if (header != null) header!,
              if (error != null)
                Padding(
                  padding: EdgeInsets.only(
                    left: padding.left + 12.0,
                    right: padding.right + 12.0,
                    bottom: 8.0,
                  ),
                  child: Center(
                    child: Text(
                      error!,
                      style: Fonts.bodyMedium.copyWith(
                        color: errorColor(context),
                      ),
                    ),
                  ),
                ),
              Padding(
                padding: EdgeInsets.only(
                  left: padding.left + 12.0,
                  right: padding.right + 12.0,
                ),
                child: Row(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    SizedBox(
                      height: 36.0,
                      child: CupertinoButton(
                        sizeStyle: CupertinoButtonSize.small,
                        padding: EdgeInsets.zero,
                        onPressed: onImageSelect,
                        child: Icon(MyIcons.plusCircleFill, size: 28.0),
                      ),
                    ),
                    const SizedBox(width: 8.0),
                    Expanded(
                      child: MyTextField(
                        controller: controller,
                        placeholder: placeholder,
                        focusNode: focusNode,
                        onChanged: onChanged,
                      ),
                    ),
                    const SizedBox(width: 8.0),
                    SizedBox(
                      height: 36.0,
                      child: CupertinoButton(
                        sizeStyle: CupertinoButtonSize.small,
                        padding: EdgeInsets.zero,
                        onPressed: onSend,
                        child: sendIcon,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class MessageInputHeader extends StatelessWidget {
  final IconData icon;
  final String text;
  final VoidCallback onPressed;
  final VoidCallback onClear;

  const MessageInputHeader({
    required this.icon,
    required this.text,
    required this.onPressed,
    required this.onClear,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        const SizedBox(width: 12.0),
        CupertinoButton(
          padding: EdgeInsets.zero,
          sizeStyle: CupertinoButtonSize.small,
          onPressed: onPressed,
          child: Icon(icon, size: 20.0),
        ),
        const SizedBox(width: 8.0),
        Container(
          decoration: BoxDecoration(
            color: primary,
            borderRadius: BorderRadius.circular(12.0),
          ),
          width: 2.0,
          height: 30.0,
        ),
        const SizedBox(width: 8.0),
        Expanded(
          child: GestureDetector(
            onTap: onPressed,
            behavior: HitTestBehavior.opaque,
            child: MyText(
              text,
              style: Fonts.bodyMedium,
              overflow: TextOverflow.ellipsis,
              isSelectable: false,
            ),
          ),
        ),
        const SizedBox(width: 8.0),
        Container(
          decoration: BoxDecoration(
            color: primary,
            borderRadius: BorderRadius.circular(12.0),
          ),
          width: 2.0,
          height: 30.0,
        ),
        const SizedBox(width: 8.0),
        CupertinoButton(
          padding: EdgeInsets.zero,
          sizeStyle: CupertinoButtonSize.small,
          onPressed: onClear,
          child: Icon(MyIcons.clear, size: 20.0),
        ),
        const SizedBox(width: 12.0),
      ],
    );
  }
}
