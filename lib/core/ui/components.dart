import 'dart:io';

import 'package:flutter/cupertino.dart';
import 'package:flutter/gestures.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:carousel_slider/carousel_slider.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:like_button/like_button.dart';
import 'package:shimmer/shimmer.dart';

import 'package:yordamchi/core/utils/utils.dart';
import 'package:yordamchi/core/repository/repository.dart';
import 'fonts.dart';
import 'icons.dart';
import 'logos.dart';
import 'themes.dart';

class YordamchiLogo extends ConsumerWidget {
  const YordamchiLogo({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final color = ref.watch(configsProvider).valueOrNull!.getColor;

    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        const SizedBox(height: 20.0),
        Logos.yordamchi(color: color, size: 100.0),
        const SizedBox(height: 10.0),
        const Text(appName, style: Fonts.titleLarge),
        const SizedBox(height: 40.0),
      ],
    );
  }
}

class AgreeToTerms extends ConsumerWidget {
  const AgreeToTerms({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Text.rich(
      TextSpan(
        children: [
          TextSpan(text: l10n.agreeToTermsLeading),
          TextSpan(
            text: l10n.agreeToTermsToS,
            style: const TextStyle(color: linkColor),
            recognizer: TapGestureRecognizer()..onTap = openTermsOfService,
          ),
          TextSpan(text: l10n.agreeToTermsMiddle),
          TextSpan(
            text: l10n.agreeToTermsPrivacyPolicy,
            style: const TextStyle(color: linkColor),
            recognizer: TapGestureRecognizer()..onTap = openPrivacyPolicy,
          ),
          TextSpan(text: l10n.agreeToTermsTrailing),
        ],
      ),
      textAlign: TextAlign.center,
      style: Fonts.bodyMedium.copyWith(color: textColor),
    );
  }
}

class ImageCarousel extends StatelessWidget {
  final List<String> images;
  final void Function(int) onTap;
  final void Function(int)? onDelete;
  final double? size;

  const ImageCarousel({
    required this.images,
    required this.onTap,
    this.onDelete,
    this.size,
    super.key,
  });

  Widget _buildImage(int index, double size) {
    final image = images[index];

    final imageWidget = Hero(
      tag: '$index:${images[index]}',
      transitionOnUserGestures: true,
      child: ClipRRect(
        borderRadius: BorderRadius.circular(12.0),
        clipBehavior: Clip.hardEdge,
        child: GestureDetector(
          onTap: () {
            onTap(index);
          },
          child: SizedBox(
            width: size,
            height: size,
            child: ImageWidget(image: image),
          ),
        ),
      ),
    );

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 6.0),
      child: onDelete != null
          ? SizedBox(
              width: size + 4.0,
              height: size + 4.0,
              child: Stack(
                children: [
                  Positioned(left: 0.0, bottom: 0.0, child: imageWidget),
                  Positioned(
                    top: 0.0,
                    right: 0.0,
                    child: Container(
                      decoration: BoxDecoration(
                        color: primary,
                        shape: BoxShape.circle,
                      ),
                      child: CupertinoButton(
                        padding: EdgeInsets.zero,
                        sizeStyle: CupertinoButtonSize.small,
                        onPressed: () => onDelete!.call(index),
                        child: Icon(
                          MyIcons.delete,
                          size: 16.0,
                          color: onPrimary,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            )
          : imageWidget,
    );
  }

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        final size =
            this.size ??
            (images.length == 1 ? constraints.maxWidth - 12.0 : 200.0);

        return (size + 12.0) * images.length < constraints.maxWidth
            ? Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  for (int index = 0; index < images.length; index++)
                    _buildImage(index, size),
                ],
              )
            : NotificationListener<ScrollUpdateNotification>(
                onNotification: (notification) => true,
                child: CarouselSlider.builder(
                  options: CarouselOptions(
                    pageViewKey: PageStorageKey(images),
                    padEnds: false,
                    pageSnapping: false,
                    enableInfiniteScroll: false,
                    viewportFraction: (size + 12.0) / constraints.maxWidth,
                    height: size,
                  ),
                  itemCount: images.length,
                  itemBuilder: (context, index, idx) {
                    return _buildImage(index, size);
                  },
                ),
              );
      },
    );
  }
}

class ImageWidget extends StatelessWidget {
  final String image;

  const ImageWidget({required this.image, super.key});

  @override
  Widget build(BuildContext context) {
    if (image.isEmpty) {
      return Shimmer.fromColors(
        baseColor: scaffold,
        highlightColor: primaryContainer,
        child: Container(color: scaffold),
      );
    }

    return image.startsWith('/')
        ? Image.file(
            File(image),
            fit: BoxFit.cover,
            frameBuilder: (context, child, frame, wasSynchronouslyLoaded) {
              return frame != null
                  ? child
                  : Shimmer.fromColors(
                      baseColor: scaffold,
                      highlightColor: primaryContainer,
                      child: Container(color: scaffold),
                    );
            },
          )
        : CachedNetworkImage(
            imageUrl: formatImageUrl(image),
            httpHeaders: image.startsWith('http') ? null : headers,
            fit: BoxFit.cover,
            fadeOutDuration: Duration.zero,
            fadeInDuration: Duration.zero,
            progressIndicatorBuilder: (context, url, progress) {
              return Shimmer.fromColors(
                baseColor: scaffold,
                highlightColor: primaryContainer,
                child: Container(color: scaffold),
              );
            },
          );
  }
}

class ScrollButton extends StatefulWidget {
  final IconData icon;
  final VoidCallback onPressed;

  const ScrollButton({required this.icon, required this.onPressed, super.key});

  @override
  State<ScrollButton> createState() => ScrollButtonState();
}

class ScrollButtonState extends State<ScrollButton> {
  bool isVisible = false;

  void changeState(bool show) {
    if (show != isVisible) {
      setState(() => isVisible = show);
    }
  }

  @override
  Widget build(BuildContext context) {
    return AnimatedScale(
      scale: isVisible ? 1.0 : 0.0,
      duration: const Duration(milliseconds: 300),
      child: CupertinoButton.filled(
        padding: const EdgeInsets.all(7.0),
        sizeStyle: CupertinoButtonSize.small,
        onPressed: widget.onPressed,
        child: Icon(widget.icon, size: Platform.isIOS ? 27.0 : 24.0),
      ),
    );
  }
}

class ActionButton extends StatefulWidget {
  final VoidCallback? onTap;
  final IconData icon;
  final String label;

  const ActionButton({
    required this.onTap,
    required this.icon,
    required this.label,
    super.key,
  });

  @override
  State<ActionButton> createState() => ActionButtonState();
}

class ActionButtonState extends State<ActionButton> {
  var _isLiked = false;

  @override
  Widget build(BuildContext context) {
    return Expanded(
      child: Column(
        children: [
          LikeButton(
            size: 18.0,
            isLiked: _isLiked,
            likeBuilder: (isLiked) {
              return Icon(widget.icon, color: primary, size: 20.0);
            },
            onTap: (isLiked) async {
              widget.onTap?.call();
              Future.delayed(const Duration(seconds: 1), () async {
                if (mounted) {
                  setState(() => _isLiked = true);
                  await Future.delayed(const Duration(milliseconds: 100));
                  if (mounted) {
                    setState(() => _isLiked = false);
                  }
                }
              });
              return true;
            },
            likeCountPadding: EdgeInsets.zero,
            padding: const EdgeInsets.all(4.0),
            circleColor: CircleColor(start: primary, end: primary),
            bubblesColor: BubblesColor(
              dotPrimaryColor: primary,
              dotSecondaryColor: primary,
            ),
          ),
          Text(
            widget.label,
            style: TextStyle(fontSize: 10.0, color: primary),
            textAlign: TextAlign.center,
          ),
        ],
      ),
    );
  }
}

class MyTextField extends StatelessWidget {
  final TextEditingController controller;
  final String placeholder;
  final FocusNode focusNode;
  final void Function(String)? onChanged;

  const MyTextField({
    required this.controller,
    required this.placeholder,
    required this.focusNode,
    this.onChanged,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return CupertinoTextField(
      controller: controller,
      placeholder: placeholder,
      focusNode: focusNode,
      maxLines: 15,
      minLines: 1,
      padding: const EdgeInsets.symmetric(horizontal: 12.0, vertical: 8.0),
      decoration: BoxDecoration(
        color: textField(context),
        borderRadius: BorderRadius.circular(12.0),
      ),
      onChanged: onChanged,
      suffix: controller.text.isNotEmpty
          ? CupertinoButton(
              padding: const EdgeInsets.symmetric(horizontal: 8.0),
              sizeStyle: CupertinoButtonSize.small,
              onPressed: () {
                controller.clear();
                onChanged?.call('');
              },
              child: Icon(
                MyIcons.clear,
                color: clearButton(context),
                size: 20.0,
              ),
            )
          : null,
    );
  }
}
