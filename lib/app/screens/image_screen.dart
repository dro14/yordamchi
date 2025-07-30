import 'dart:io';
import 'dart:ui';
import 'dart:async';

import 'package:flutter/cupertino.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:photo_view/photo_view.dart';
import 'package:photo_view/photo_view_gallery.dart';
import 'package:cached_network_image/cached_network_image.dart';

import 'package:yordamchi/core/core.dart';

class ImageScreen extends ConsumerStatefulWidget {
  final List<String> images;
  final int initialIndex;
  final EdgeInsets padding;

  const ImageScreen(this.images, this.initialIndex, this.padding, {super.key});

  static Future<T?> route<T>(
    BuildContext context,
    List<String> images, {
    int initialIndex = 0,
  }) {
    final padding = MediaQuery.of(context).padding;
    return Navigator.of(context, rootNavigator: true).push<T>(
      CupertinoPageRoute<T>(
        builder: (context) {
          return ImageScreen(images, initialIndex, padding);
        },
      ),
    );
  }

  @override
  ConsumerState<ImageScreen> createState() => _ImageScreenState();
}

class _ImageScreenState extends ConsumerState<ImageScreen> {
  late final PageController _pageController;
  late int _currentIndex;
  double _dragOffset = 0.0;
  bool _showActions = true;
  Timer? _timer;

  @override
  void initState() {
    super.initState();
    Analytics.image();
    _pageController = PageController(initialPage: widget.initialIndex);
    _currentIndex = widget.initialIndex;
    _timer = Timer(const Duration(seconds: 3), () {
      setState(() => _showActions = false);
    });

    SystemChrome.setEnabledSystemUIMode(SystemUiMode.immersive);
  }

  @override
  void dispose() {
    _pageController.dispose();
    _timer?.cancel();
    SystemChrome.setEnabledSystemUIMode(SystemUiMode.edgeToEdge);
    super.dispose();
  }

  ImageProvider _getProvider(int index) {
    final image = widget.images[index];
    return image.startsWith('/')
        ? FileImage(File(image))
        : CachedNetworkImageProvider(
            formatImageUrl(image),
            headers: image.startsWith('http') ? null : headers,
          );
  }

  @override
  Widget build(BuildContext context) {
    return CupertinoTheme(
      data: CupertinoThemeData(
        brightness: Brightness.dark,
        primaryColor: white,
        primaryContrastingColor: black,
        barBackgroundColor: black.withAlpha(100),
        scaffoldBackgroundColor: black.withAlpha(100),
        applyThemeToAll: true,
      ),
      child: CupertinoPageScaffold(
        resizeToAvoidBottomInset: false,
        child: Stack(
          children: [
            GestureDetector(
              onTap: () {
                setState(() => _showActions = !_showActions);
                _timer?.cancel();
                if (_showActions) {
                  _timer = Timer(const Duration(seconds: 3), () {
                    setState(() => _showActions = false);
                  });
                }
              },
              onVerticalDragUpdate: (details) {
                setState(() => _dragOffset += details.primaryDelta ?? 0);
              },
              onVerticalDragEnd: (details) {
                if (details.velocity.pixelsPerSecond.dy.abs() > 300.0 ||
                    _dragOffset.abs() > 100.0) {
                  Navigator.pop(context);
                } else {
                  setState(() => _dragOffset = 0.0);
                }
              },
              child: Transform.translate(
                offset: Offset(0, _dragOffset),
                child: PhotoViewGallery.builder(
                  pageController: _pageController,
                  itemCount: widget.images.length,
                  onPageChanged: (index) {
                    setState(() => _currentIndex = index);
                  },
                  builder: (context, index) {
                    if (index > 0) {
                      precacheImage(_getProvider(index - 1), context);
                    }
                    if (index + 1 < widget.images.length) {
                      precacheImage(_getProvider(index + 1), context);
                    }
                    return PhotoViewGalleryPageOptions(
                      imageProvider: _getProvider(index),
                      initialScale: PhotoViewComputedScale.contained,
                      minScale: PhotoViewComputedScale.contained,
                      maxScale: PhotoViewComputedScale.covered * 2.0,
                      heroAttributes: PhotoViewHeroAttributes(
                        tag: '$index:${widget.images[index]}',
                        transitionOnUserGestures: true,
                      ),
                    );
                  },
                  loadingBuilder: (context, loadingProgress) {
                    return const CupertinoActivityIndicator();
                  },
                ),
              ),
            ),
            Positioned(
              left: widget.padding.left + 8.0,
              top: widget.padding.top,
              right: widget.padding.right + 8.0,
              child: Row(
                children: [
                  Expanded(
                    child: Align(
                      alignment: Alignment.centerLeft,
                      child: _Blur(
                        isVisible: _showActions,
                        child: CupertinoButton(
                          padding: const EdgeInsets.all(6.0),
                          sizeStyle: CupertinoButtonSize.small,
                          onPressed: () => Navigator.pop(context),
                          child: Icon(MyIcons.back, color: white, size: 24.0),
                        ),
                      ),
                    ),
                  ),
                  if (widget.images.length > 1) ...[
                    Expanded(
                      child: Center(
                        child: _Blur(
                          isVisible: _showActions,
                          child: Padding(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 12.0,
                              vertical: 6.0,
                            ),
                            child: Text(
                              '${_currentIndex + 1}/${widget.images.length}',
                              style: Fonts.titleMedium.copyWith(color: white),
                            ),
                          ),
                        ),
                      ),
                    ),
                    const Spacer(),
                  ],
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _Blur extends StatelessWidget {
  final bool isVisible;
  final Widget child;

  const _Blur({required this.isVisible, required this.child});

  @override
  Widget build(BuildContext context) {
    return AnimatedScale(
      scale: isVisible ? 1.0 : 0.0,
      duration: const Duration(milliseconds: 300),
      child: ClipRRect(
        borderRadius: BorderRadius.circular(100.0),
        child: BackdropFilter(
          filter: ImageFilter.blur(sigmaX: 10.0, sigmaY: 10.0),
          child: Container(
            decoration: BoxDecoration(
              color: black.withAlpha(100),
              borderRadius: BorderRadius.circular(100.0),
            ),
            child: child,
          ),
        ),
      ),
    );
  }
}
