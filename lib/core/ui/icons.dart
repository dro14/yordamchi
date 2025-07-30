import 'dart:io';

import 'package:flutter/material.dart' show Icons;
import 'package:flutter/cupertino.dart';

class MyIcons {
  static bool isIOS = Platform.isIOS;

  static IconData settings = Icons.settings_outlined;

  static IconData plus = Icons.add_rounded;

  static IconData history = Icons.history_rounded;

  static IconData plusCircleFill = CupertinoIcons.plus_circle_fill;

  static IconData clear = CupertinoIcons.xmark_circle_fill;

  static IconData send = CupertinoIcons.arrow_up_circle_fill;

  static IconData checkCircleFill = CupertinoIcons.checkmark_alt_circle_fill;

  static IconData up = isIOS
      ? CupertinoIcons.chevron_up
      : Icons.arrow_upward_rounded;

  static IconData down = isIOS
      ? CupertinoIcons.chevron_down
      : Icons.arrow_downward_rounded;

  static IconData delete = CupertinoIcons.trash;

  static IconData error = CupertinoIcons.exclamationmark_octagon;

  static IconData language = CupertinoIcons.globe;

  static IconData brightness = CupertinoIcons.sun_max;

  static IconData theme = CupertinoIcons.paintbrush;

  static IconData termsOfService = CupertinoIcons.doc_text;

  static IconData privacyPolicy = CupertinoIcons.lock_shield;

  static IconData contact = CupertinoIcons.paperplane;

  static IconData check = Icons.check_rounded;

  static IconData reply = isIOS
      ? CupertinoIcons.arrowshape_turn_up_left
      : Icons.reply_rounded;

  static IconData copy = isIOS
      ? CupertinoIcons.doc_on_doc
      : Icons.content_copy_rounded;

  static IconData retry = isIOS
      ? CupertinoIcons.arrow_counterclockwise
      : Icons.refresh_rounded;

  static IconData edit = isIOS
      ? CupertinoIcons.square_pencil
      : Icons.edit_rounded;

  static IconData followUps = isIOS
      ? CupertinoIcons.question
      : Icons.question_mark_rounded;

  static IconData back = isIOS
      ? CupertinoIcons.chevron_left
      : Icons.arrow_back_rounded;

  static IconData options = isIOS
      ? CupertinoIcons.ellipsis_circle
      : Icons.more_vert_rounded;

  static IconData share = isIOS
      ? CupertinoIcons.square_arrow_up
      : Icons.share_outlined;
}
