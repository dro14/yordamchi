import 'dart:io';
import 'dart:ui';
import 'dart:math';

import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart' show showDialog, AlertDialog, ListTile;

import 'package:yordamchi/core/repository/repository.dart';
import 'components.dart';
import 'fonts.dart';
import 'icons.dart';
import 'themes.dart';

Future<void> showColorPicker(
  BuildContext context,
  int initialIndex,
  void Function(int) onChanged,
) async {
  if (Platform.isIOS) {
    await showCupertinoModalPopup(
      context: context,
      builder: (context) {
        return Container(
          height: 216.0,
          padding: const EdgeInsets.only(top: 6.0),
          margin: EdgeInsets.only(
            bottom: MediaQuery.viewInsetsOf(context).bottom,
          ),
          decoration: BoxDecoration(
            color: surfaceContainer,
            borderRadius: const BorderRadius.only(
              topLeft: Radius.circular(12.0),
              topRight: Radius.circular(12.0),
            ),
          ),
          child: SafeArea(
            top: false,
            child: CupertinoPicker(
              magnification: 1.22,
              squeeze: 1.2,
              useMagnifier: true,
              itemExtent: 32.0,
              scrollController: FixedExtentScrollController(
                initialItem: initialIndex,
              ),
              onSelectedItemChanged: onChanged,
              children: [
                for (int index = 0; index < colors.length; index++)
                  Row(
                    children: [
                      SizedBox(width: 0.4 * MediaQuery.sizeOf(context).width),
                      Container(
                        width: 20.0,
                        height: 20.0,
                        decoration: BoxDecoration(
                          shape: BoxShape.circle,
                          color: colors.values.elementAt(index),
                        ),
                      ),
                      const SizedBox(width: 8.0),
                      Text(
                        getColorName(index),
                        style: TextStyle(color: colors.values.elementAt(index)),
                      ),
                    ],
                  ),
              ],
            ),
          ),
        );
      },
    );
  } else {
    final scrollController = ScrollController();

    final result = await showDialog(
      context: context,
      builder: (context) {
        WidgetsBinding.instance.addPostFrameCallback((duration) {
          scrollController.animateTo(
            max(initialIndex - 2, 0) * 56.0,
            duration: const Duration(seconds: 1),
            curve: Curves.easeInOut,
          );
        });

        return AlertDialog(
          contentPadding: const EdgeInsets.symmetric(horizontal: 24.0),
          backgroundColor: surfaceContainer,
          content: SizedBox(
            width: double.maxFinite,
            height: 300.0,
            child: ListView.builder(
              padding: const EdgeInsets.symmetric(vertical: 12.0),
              controller: scrollController,
              itemExtent: 56.0,
              itemCount: colors.length,
              itemBuilder: (context, index) {
                final entry = colors.entries.elementAt(index);
                final isSelected = index == initialIndex;

                return ListTile(
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12.0),
                  ),
                  leading: Container(
                    width: 24.0,
                    height: 24.0,
                    decoration: BoxDecoration(
                      shape: BoxShape.circle,
                      color: entry.value,
                    ),
                  ),
                  title: Text(getColorName(index)),
                  titleTextStyle: TextStyle(
                    color: entry.value,
                    fontWeight: isSelected ? FontWeight.bold : null,
                    fontFamily: 'CupertinoSystemText',
                  ),
                  tileColor: isSelected ? entry.value.withAlpha(40) : null,
                  onTap: () {
                    Navigator.pop(context, index);
                  },
                );
              },
            ),
          ),
        );
      },
    );
    if (result != null) {
      onChanged(result);
    }
  }
}

Future<int?> showPicker({
  required BuildContext context,
  required List<String> options,
  int? initialIndex,
  required String title,
}) async {
  return showCupertinoModalPopup(
    context: context,
    builder: (context) {
      return CupertinoPopupSurface(
        blurSigma: 0.0,
        child: SingleChildScrollView(
          child: Container(
            color: surfaceContainer,
            child: SafeArea(
              top: false,
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Padding(
                    padding: const EdgeInsets.only(
                      left: 20.0,
                      top: 16.0,
                      right: 20.0,
                    ),
                    child: Text(
                      title,
                      style: Fonts.titleMedium.copyWith(color: textColor),
                      textAlign: TextAlign.center,
                    ),
                  ),
                  SizedBox(
                    height: 44.0 * options.length + 44.0,
                    child: CupertinoListSection.insetGrouped(
                      backgroundColor: surfaceContainer,
                      children: [
                        for (int index = 0; index < options.length; index++)
                          CupertinoListTile(
                            title: Text(options[index]),
                            leading: CupertinoRadio<int>(
                              value: index,
                              groupValue: initialIndex,
                              onChanged: (index) {
                                Navigator.pop(context, index);
                              },
                            ),
                            onTap: () {
                              Navigator.pop(context, index);
                            },
                          ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      );
    },
  );
}

Future<bool> showDestructiveDialog({
  required BuildContext context,
  required String title,
  required String message,
  required String cancel,
  required String confirm,
}) async {
  final bool? result;
  if (Platform.isIOS) {
    result = await showCupertinoDialog(
      barrierDismissible: true,
      context: context,
      builder: (context) {
        return CupertinoAlertDialog(
          title: Text(title),
          content: Text(message),
          actions: [
            CupertinoDialogAction(
              onPressed: () => Navigator.pop(context, false),
              child: Text(cancel),
            ),
            CupertinoDialogAction(
              isDestructiveAction: true,
              onPressed: () => Navigator.pop(context, true),
              child: Text(confirm),
            ),
          ],
        );
      },
    );
  } else {
    result = await showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: Text(title),
          titleTextStyle: TextStyle(
            color: textColor,
            fontFamily: 'CupertinoSystemText',
            fontSize: 24.0,
          ),
          content: Text(message),
          contentTextStyle: TextStyle(
            color: textColor,
            fontFamily: 'CupertinoSystemText',
            fontSize: 14.0,
          ),
          backgroundColor: surfaceContainer,
          actions: [
            CupertinoButton(
              sizeStyle: CupertinoButtonSize.small,
              onPressed: () => Navigator.pop(context, false),
              child: Text(cancel),
            ),
            CupertinoButton(
              sizeStyle: CupertinoButtonSize.small,
              onPressed: () => Navigator.pop(context, true),
              child: Text(
                confirm,
                style: TextStyle(color: errorColor(context)),
              ),
            ),
          ],
        );
      },
    );
  }
  return result ?? false;
}

Future<bool> showInfoDialog({
  required BuildContext context,
  required String title,
  required String message,
  String? cancel,
  String? confirm,
}) async {
  final bool? result;
  if (Platform.isIOS) {
    result = await showCupertinoDialog(
      context: context,
      builder: (context) {
        return CupertinoAlertDialog(
          title: Text(title),
          content: Text(message),
          actions: [
            if (cancel != null)
              CupertinoDialogAction(
                onPressed: () => Navigator.pop(context, false),
                child: Text(cancel),
              ),
            CupertinoDialogAction(
              onPressed: () => Navigator.pop(context, true),
              child: Text(confirm ?? 'OK'),
            ),
          ],
        );
      },
    );
  } else {
    result = await showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: Text(title),
          titleTextStyle: TextStyle(
            color: textColor,
            fontFamily: 'CupertinoSystemText',
            fontSize: 24.0,
          ),
          content: Text(message),
          contentTextStyle: TextStyle(
            color: textColor,
            fontFamily: 'CupertinoSystemText',
            fontSize: 14.0,
          ),
          backgroundColor: surfaceContainer,
          actions: [
            if (cancel != null)
              CupertinoButton(
                sizeStyle: CupertinoButtonSize.small,
                onPressed: () => Navigator.pop(context, false),
                child: Text(cancel),
              ),
            CupertinoButton(
              sizeStyle: CupertinoButtonSize.small,
              onPressed: () => Navigator.pop(context, true),
              child: Text(confirm ?? 'OK'),
            ),
          ],
        );
      },
    );
  }
  return result ?? false;
}

Future<int?> showActionSheet({
  required BuildContext context,
  required List<IconData> icons,
  required List<String> actions,
}) async {
  assert(icons.length == actions.length);

  return showCupertinoModalPopup(
    context: context,
    builder: (context) {
      return CupertinoPopupSurface(
        blurSigma: 0.0,
        child: Container(
          color: surfaceContainer,
          child: SafeArea(
            top: false,
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                const SizedBox(height: 12.0),
                for (int index = 0; index < actions.length; index++)
                  CupertinoButton(
                    onPressed: () => Navigator.pop(context, index),
                    child: Row(
                      children: [
                        Icon(icons[index], size: 24.0),
                        const SizedBox(width: 12.0),
                        Text(actions[index]),
                      ],
                    ),
                  ),
              ],
            ),
          ),
        ),
      );
    },
  );
}

Future<String?> showInputDialog({required BuildContext context}) async {
  final controller = TextEditingController();
  final focusNode = FocusNode();

  WidgetsBinding.instance.addPostFrameCallback((duration) {
    focusNode.requestFocus();
  });

  return showCupertinoModalPopup(
    context: context,
    builder: (context) {
      var padding = MediaQuery.paddingOf(context);
      padding = EdgeInsets.only(
        left: padding.left + 12.0,
        top: 8.0,
        right: padding.right + 12.0,
        bottom: padding.bottom + MediaQuery.viewInsetsOf(context).bottom + 8.0,
      );

      return CupertinoPopupSurface(
        blurSigma: 0.0,
        child: ClipRect(
          child: BackdropFilter(
            filter: ImageFilter.blur(sigmaX: 10.0, sigmaY: 10.0),
            child: Container(
              color: surfaceContainer.withAlpha(200),
              padding: padding,
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Expanded(
                    child: MyTextField(
                      controller: controller,
                      placeholder: l10n.giveThisChatName,
                      focusNode: focusNode,
                    ),
                  ),
                  const SizedBox(width: 8.0),
                  SizedBox(
                    height: 36.0,
                    child: CupertinoButton(
                      sizeStyle: CupertinoButtonSize.small,
                      padding: EdgeInsets.zero,
                      onPressed: () {
                        final text = controller.text.trim();
                        Navigator.pop(context, text.isEmpty ? null : text);
                      },
                      child: Icon(MyIcons.checkCircleFill, size: 28.0),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      );
    },
  );
}
