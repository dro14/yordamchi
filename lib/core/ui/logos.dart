import 'package:flutter/cupertino.dart';
import 'package:flutter_svg/flutter_svg.dart';

class Logos {
  static const _path = 'assets/logos';

  static SvgPicture yordamchi({required Color color, required double size}) {
    return SvgPicture.asset(
      '$_path/yordamchi.svg',
      colorFilter: ColorFilter.mode(color, BlendMode.srcIn),
      width: size,
    );
  }
}
