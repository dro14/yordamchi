import 'package:flutter/cupertino.dart';
import 'package:shimmer/shimmer.dart';

import 'package:yordamchi/core/repository/repository.dart';
import 'package:yordamchi/core/models/models.dart';
import 'package:yordamchi/core/utils/utils.dart';
import 'package:yordamchi/core/ui/ui.dart';

class MessageWidget extends StatefulWidget {
  final Message? message;
  final Message? inReplyTo;
  final String? chunk;
  final AnimationController animController;
  final VoidCallback? onReply;
  final VoidCallback? onRetry;
  final VoidCallback? onEdit;
  final VoidCallback? onReplyTap;
  final VoidCallback? onFollowUps;
  final void Function(String)? onFollowUpTap;
  final void Function(List<String>, int)? onImageTap;

  const MessageWidget({
    required this.animController,
    this.message,
    this.inReplyTo,
    this.chunk,
    this.onReply,
    this.onRetry,
    this.onEdit,
    this.onReplyTap,
    this.onFollowUps,
    this.onFollowUpTap,
    this.onImageTap,
    super.key,
  }) : assert(
         message != null || chunk != null,
         'Either message or chunk must be provided',
       );

  @override
  State<MessageWidget> createState() => _MessageWidgetState();
}

class _MessageWidgetState extends State<MessageWidget> {
  late final _animController = widget.animController;
  late final Animation<Color?> _animation;

  @override
  void initState() {
    super.initState();
    _animation = ColorTween(
      begin: widget.message?.role == 'user' ? primaryContainer : scaffold,
      end: primary,
    ).animate(_animController);
  }

  @override
  Widget build(BuildContext context) {
    final message = widget.message;
    final inReplyTo = widget.inReplyTo;
    final chunk = widget.chunk;
    final Widget child;
    if (chunk != null) {
      child = Padding(
        padding: const EdgeInsets.only(
          left: 12.0,
          top: 8.0,
          right: 12.0,
          bottom: 14.0,
        ),
        child: chunk.isEmpty
            ? Shimmer.fromColors(
                baseColor: textColor,
                highlightColor: scaffold,
                child: Text(l10n.thinking, style: Fonts.bodyLarge),
              )
            : MyText(chunk, style: Fonts.bodyLarge),
      );
    } else if (message != null && message.followUps.isNotEmpty) {
      child = Padding(
        padding: const EdgeInsets.symmetric(horizontal: 12.0),
        child: Column(
          children: [
            const SizedBox(height: 12.0),
            for (final followUp in message.followUps) ...[
              ConstrainedBox(
                constraints: const BoxConstraints(minHeight: buttonHeight),
                child: SizedBox(
                  width: double.infinity,
                  child: CupertinoButton.tinted(
                    padding: const EdgeInsets.all(12.0),
                    sizeStyle: CupertinoButtonSize.small,
                    onPressed: () => widget.onFollowUpTap?.call(followUp),
                    child: MyText(
                      followUp,
                      style: Fonts.bodyLarge,
                      textAlign: TextAlign.center,
                      isLinkBlue: false,
                      isSelectable: false,
                    ),
                  ),
                ),
              ),
              const SizedBox(height: 12.0),
            ],
          ],
        ),
      );
    } else if (message != null && message.role == 'model') {
      if (message.text.isEmpty || message.images.isEmpty) {
        child = Padding(
          padding: const EdgeInsets.only(left: 12.0, top: 8.0, right: 12.0),
          child: Text(
            l10n.webSearchResults,
            style: Fonts.titleMedium.copyWith(color: primary),
          ),
        );
      } else {
        child = Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            AnimatedBuilder(
              animation: _animation,
              builder: (context, child) {
                return Container(
                  decoration: BoxDecoration(
                    color: _animController.isAnimating
                        ? _animation.value
                        : scaffold,
                    borderRadius: BorderRadius.circular(12.0),
                  ),
                  child: child,
                );
              },
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const SizedBox(height: 8.0),
                  if (message.images.isNotEmpty) ...[
                    ImageCarousel(
                      images: message.images,
                      onTap: (index) {
                        widget.onImageTap?.call(message.images, index);
                      },
                    ),
                    const SizedBox(height: 6.0),
                  ],
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 12.0),
                    child: MyText(message.text, style: Fonts.bodyLarge),
                  ),
                  const SizedBox(height: 6.0),
                ],
              ),
            ),
            const SizedBox(height: 2.0),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 12.0),
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  ActionButton(
                    onTap: widget.onReply,
                    icon: MyIcons.reply,
                    label: l10n.reply,
                  ),
                  ActionButton(
                    onTap: () => copyToClipboard(message.text),
                    icon: MyIcons.copy,
                    label: l10n.copy,
                  ),
                  ActionButton(
                    onTap: widget.onRetry,
                    icon: MyIcons.retry,
                    label: l10n.retry,
                  ),
                  ActionButton(
                    onTap: widget.onEdit,
                    icon: MyIcons.edit,
                    label: l10n.edit,
                  ),
                  ActionButton(
                    onTap: widget.onFollowUps,
                    icon: MyIcons.followUps,
                    label: l10n.followUps,
                  ),
                ],
              ),
            ),
            const SizedBox(height: 14.0),
          ],
        );
      }
    } else if (message != null && message.role == 'user') {
      if (message.text.isEmpty && message.images.isEmpty) {
        child = const SizedBox.shrink();
      } else {
        child = Padding(
          padding: const EdgeInsets.only(left: 8.0, right: 8.0, bottom: 6.0),
          child: LayoutBuilder(
            builder: (context, constraints) {
              return Align(
                alignment: Alignment.centerRight,
                child: AnimatedBuilder(
                  animation: _animation,
                  builder: (context, child) {
                    return Container(
                      constraints: BoxConstraints(
                        maxWidth: 0.9 * constraints.maxWidth,
                      ),
                      decoration: BoxDecoration(
                        color: _animController.isAnimating
                            ? _animation.value
                            : primaryContainer,
                        borderRadius: BorderRadius.circular(12.0),
                      ),
                      child: child,
                    );
                  },
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const SizedBox(height: 8.0),
                      if (inReplyTo != null) ...[
                        GestureDetector(
                          onTap: widget.onReplyTap,
                          behavior: HitTestBehavior.opaque,
                          child: Container(
                            margin: const EdgeInsets.symmetric(horizontal: 8.0),
                            decoration: BoxDecoration(
                              color: onInverseSurface,
                              borderRadius: BorderRadius.circular(4.0),
                            ),
                            child: Row(
                              mainAxisSize: MainAxisSize.min,
                              children: [
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
                                  child: MyText(
                                    inReplyTo.displayText,
                                    style: Fonts.bodyLarge,
                                    overflow: TextOverflow.ellipsis,
                                    isSelectable: false,
                                  ),
                                ),
                                const SizedBox(width: 8.0),
                              ],
                            ),
                          ),
                        ),
                        const SizedBox(height: 6.0),
                      ],
                      if (message.images.isNotEmpty) ...[
                        ImageCarousel(
                          images: message.images,
                          onTap: (index) {
                            widget.onImageTap?.call(message.images, index);
                          },
                        ),
                        const SizedBox(height: 6.0),
                      ],
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 8.0),
                        child: MyText(
                          message.text,
                          style: Fonts.bodyLarge.copyWith(
                            color: onPrimaryContainer,
                          ),
                          isLinkBlue: false,
                        ),
                      ),
                      const SizedBox(height: 8.0),
                    ],
                  ),
                ),
              );
            },
          ),
        );
      }
    } else {
      child = const SizedBox.shrink();
    }
    return RepaintBoundary(
      child: AnimatedSize(
        duration: const Duration(milliseconds: 150),
        curve: Curves.easeInOutCubic,
        alignment: Alignment.topLeft,
        child: child,
      ),
    );
  }
}
