import 'package:firebase_analytics/firebase_analytics.dart';

class Analytics {
  static final _analytics = FirebaseAnalytics.instance;

  static void _logScreenView(String screenName, String screenClass) {
    _analytics.logScreenView(screenName: screenName, screenClass: screenClass);
  }

  static void chat() {
    _logScreenView('Chat', 'ChatScreen');
  }

  static void history() {
    _logScreenView('History', 'HistoryScreen');
  }

  static void settings() {
    _logScreenView('Settings', 'SettingsScreen');
  }

  static void image() {
    _logScreenView('Image', 'ImageScreen');
  }
}
