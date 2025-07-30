import 'package:flutter/material.dart' show Material;
import 'package:flutter/cupertino.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:dots_indicator/dots_indicator.dart';
import 'package:carousel_slider/carousel_slider.dart';

import 'package:yordamchi/core/core.dart';

class WelcomeScreen extends ConsumerStatefulWidget {
  const WelcomeScreen({super.key});

  static Future<T?> route<T>(BuildContext context) {
    return Navigator.of(context, rootNavigator: true).push<T>(
      CupertinoPageRoute<T>(builder: (context) => const WelcomeScreen()),
    );
  }

  @override
  ConsumerState<WelcomeScreen> createState() => _WelcomeScreenState();
}

class _WelcomeScreenState extends ConsumerState<WelcomeScreen> {
  final _controller = CarouselSliderController();
  var _index = 0.0;

  Widget _buildItem({
    required String title,
    required String label,
    required Widget child,
  }) {
    return SizedBox(
      width: MediaQuery.sizeOf(context).width * 0.8,
      child: Column(
        spacing: 10.0,
        mainAxisSize: MainAxisSize.min,
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          child,
          Text(title, style: Fonts.titleLarge, textAlign: TextAlign.center),
          Text(label, style: Fonts.bodyLarge, textAlign: TextAlign.center),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final color = ref.read(configsProvider).valueOrNull!.getColor;

    return CupertinoPageScaffold(
      child: Center(
        child: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: MediaQuery.sizeOf(context).height * 0.1),
              CarouselSlider(
                carouselController: _controller,
                options: CarouselOptions(
                  aspectRatio: 1.25,
                  viewportFraction: 1.0,
                  pageViewKey: const PageStorageKey('welcome'),
                  autoPlay: true,
                  autoPlayCurve: Curves.easeOutCubic,
                  onPageChanged: (index, reason) {
                    setState(() => _index = index.toDouble());
                  },
                ),
                items: [
                  _buildItem(
                    title: appName,
                    label: l10n.tagline1,
                    child: Logos.yordamchi(color: color, size: 150.0),
                  ),
                  _buildItem(
                    title: l10n.title2,
                    label: l10n.tagline2,
                    child: Icon(
                      CupertinoIcons.question,
                      color: color,
                      size: 150.0,
                    ),
                  ),
                  _buildItem(
                    title: l10n.title3,
                    label: l10n.tagline3,
                    child: Icon(
                      CupertinoIcons.clock,
                      color: color,
                      size: 150.0,
                    ),
                  ),
                  _buildItem(
                    title: l10n.title4,
                    label: l10n.tagline4,
                    child: Icon(
                      CupertinoIcons.lightbulb,
                      color: color,
                      size: 150.0,
                    ),
                  ),
                  _buildItem(
                    title: l10n.title5,
                    label: l10n.tagline5,
                    child: Icon(
                      CupertinoIcons.camera,
                      color: color,
                      size: 150.0,
                    ),
                  ),
                ],
              ),
              Material(
                color: transparent,
                child: DotsIndicator(
                  dotsCount: 5,
                  animate: true,
                  position: _index,
                  decorator: DotsDecorator(
                    activeColor: color,
                    size: const Size.square(9.0),
                    activeSize: const Size(18.0, 9.0),
                    activeShape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(5.0),
                    ),
                  ),
                  onTap: (index) {
                    _controller.animateToPage(index);
                  },
                ),
              ),
              SizedBox(height: MediaQuery.sizeOf(context).height * 0.2),
              SizedBox(
                width: MediaQuery.sizeOf(context).width * 0.8,
                child: CupertinoButton.filled(
                  color: primary,
                  onPressed: () async {
                    final userId = await newUser();
                    ref.read(configsProvider.notifier).updateId(userId);
                  },
                  child: Text(
                    l10n.continueButton,
                    style: TextStyle(color: onPrimary),
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
