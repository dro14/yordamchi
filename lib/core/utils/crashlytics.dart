import 'package:firebase_crashlytics/firebase_crashlytics.dart';

class Crashlytics {
  static final _crashlytics = FirebaseCrashlytics.instance;

  static void log(Object error, StackTrace? stack, {bool fatal = true}) {
    _crashlytics.recordError(error, stack, fatal: fatal);
  }
}
