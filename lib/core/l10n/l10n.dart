import 'dart:async';

import 'package:flutter/widgets.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:intl/intl.dart' as intl;

import 'l10n_en.dart' deferred as l10n_en;
import 'l10n_ru.dart' deferred as l10n_ru;
import 'l10n_uz.dart' deferred as l10n_uz;

// ignore_for_file: type=lint

/// Callers can lookup localized strings with an instance of AppLocalizations
/// returned by `AppLocalizations.of(context)`.
///
/// Applications need to include `AppLocalizations.delegate()` in their app's
/// `localizationDelegates` list, and the locales they support in the app's
/// `supportedLocales` list. For example:
///
/// ```dart
/// import 'l10n/l10n.dart';
///
/// return MaterialApp(
///   localizationsDelegates: AppLocalizations.localizationsDelegates,
///   supportedLocales: AppLocalizations.supportedLocales,
///   home: MyApplicationHome(),
/// );
/// ```
///
/// ## Update pubspec.yaml
///
/// Please make sure to update your pubspec.yaml to include the following
/// packages:
///
/// ```yaml
/// dependencies:
///   # Internationalization support.
///   flutter_localizations:
///     sdk: flutter
///   intl: any # Use the pinned version from flutter_localizations
///
///   # Rest of dependencies
/// ```
///
/// ## iOS Applications
///
/// iOS applications define key application metadata, including supported
/// locales, in an Info.plist file that is built into the application bundle.
/// To configure the locales supported by your app, you’ll need to edit this
/// file.
///
/// First, open your project’s ios/Runner.xcworkspace Xcode workspace file.
/// Then, in the Project Navigator, open the Info.plist file under the Runner
/// project’s Runner folder.
///
/// Next, select the Information Property List item, select Add Item from the
/// Editor menu, then select Localizations from the pop-up menu.
///
/// Select and expand the newly-created Localizations item then, for each
/// locale your application supports, add a new item and select the locale
/// you wish to add from the pop-up menu in the Value field. This list should
/// be consistent with the languages listed in the AppLocalizations.supportedLocales
/// property.
abstract class AppLocalizations {
  AppLocalizations(String locale)
    : localeName = intl.Intl.canonicalizedLocale(locale.toString());

  final String localeName;

  static AppLocalizations? of(BuildContext context) {
    return Localizations.of<AppLocalizations>(context, AppLocalizations);
  }

  static const LocalizationsDelegate<AppLocalizations> delegate =
      _AppLocalizationsDelegate();

  /// A list of this localizations delegate along with the default localizations
  /// delegates.
  ///
  /// Returns a list of localizations delegates containing this delegate along with
  /// GlobalMaterialLocalizations.delegate, GlobalCupertinoLocalizations.delegate,
  /// and GlobalWidgetsLocalizations.delegate.
  ///
  /// Additional delegates can be added by appending to this list in
  /// MaterialApp. This list does not have to be used at all if a custom list
  /// of delegates is preferred or required.
  static const List<LocalizationsDelegate<dynamic>> localizationsDelegates =
      <LocalizationsDelegate<dynamic>>[
        delegate,
        GlobalMaterialLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
      ];

  /// A list of this localizations delegate's supported locales.
  static const List<Locale> supportedLocales = <Locale>[
    Locale('uz'),
    Locale('en'),
    Locale('ru'),
  ];

  /// Label for the 'Update available' title
  ///
  /// In en, this message translates to:
  /// **'Update available'**
  String get updateTitle;

  /// Label for the 'Update message prefix' message
  ///
  /// In en, this message translates to:
  /// **'A new version of Yordamchi is available!'**
  String get updateMessagePrefix;

  /// Label for the 'Update message suffix' message
  ///
  /// In en, this message translates to:
  /// **' To continue, please update the app.'**
  String get updateMessageSuffix;

  /// Label for the 'Update' button
  ///
  /// In en, this message translates to:
  /// **'Update'**
  String get update;

  /// Label for the 'Not now' button
  ///
  /// In en, this message translates to:
  /// **'Not now'**
  String get notNow;

  /// The first tagline for the welcome screen
  ///
  /// In en, this message translates to:
  /// **'Your personal assistant in daily life'**
  String get tagline1;

  /// The second title for the welcome screen
  ///
  /// In en, this message translates to:
  /// **'Anything'**
  String get title2;

  /// The second tagline for the welcome screen
  ///
  /// In en, this message translates to:
  /// **'Yordamchi answers to your questions on any topic'**
  String get tagline2;

  /// The third title for the welcome screen
  ///
  /// In en, this message translates to:
  /// **'Anytime'**
  String get title3;

  /// The third tagline for the welcome screen
  ///
  /// In en, this message translates to:
  /// **'Yordamchi is always ready to help you'**
  String get tagline3;

  /// The fourth title for the welcome screen
  ///
  /// In en, this message translates to:
  /// **'Knowledge'**
  String get title4;

  /// The fourth tagline for the welcome screen
  ///
  /// In en, this message translates to:
  /// **'Learn new things with Yordamchi'**
  String get tagline4;

  /// The fifth title for the welcome screen
  ///
  /// In en, this message translates to:
  /// **'Camera'**
  String get title5;

  /// The fifth tagline for the welcome screen
  ///
  /// In en, this message translates to:
  /// **'Yordamchi helps you understand what you see'**
  String get tagline5;

  /// Label for the 'Continue' button
  ///
  /// In en, this message translates to:
  /// **'Continue'**
  String get continueButton;

  /// Leading for the 'Agree to terms' message
  ///
  /// In en, this message translates to:
  /// **'By continuing, you agree to the '**
  String get agreeToTermsLeading;

  /// Label for the 'Terms of service' link
  ///
  /// In en, this message translates to:
  /// **'terms of service'**
  String get agreeToTermsToS;

  /// Middle for the 'Agree to terms' message
  ///
  /// In en, this message translates to:
  /// **' and '**
  String get agreeToTermsMiddle;

  /// Label for the 'Privacy policy' link
  ///
  /// In en, this message translates to:
  /// **'privacy policy'**
  String get agreeToTermsPrivacyPolicy;

  /// Trailing for the 'Agree to terms' message
  ///
  /// In en, this message translates to:
  /// **'.'**
  String get agreeToTermsTrailing;

  /// Welcoming message in the assistant screen
  ///
  /// In en, this message translates to:
  /// **'How can I help you today?'**
  String get howCanIHelpYou;

  /// Placeholder for the assistant input field
  ///
  /// In en, this message translates to:
  /// **'Ask anything...'**
  String get askAnything;

  /// Label for the 'Thinking' message
  ///
  /// In en, this message translates to:
  /// **'Thinking...'**
  String get thinking;

  /// Label for the 'Web search results' message
  ///
  /// In en, this message translates to:
  /// **'Web search results'**
  String get webSearchResults;

  /// Error message when request is too long
  ///
  /// In en, this message translates to:
  /// **'Request is too long'**
  String get errorTooLong;

  /// Label for the 'Camera' action
  ///
  /// In en, this message translates to:
  /// **'Camera'**
  String get camera;

  /// Label for the 'Photos' action
  ///
  /// In en, this message translates to:
  /// **'Photos'**
  String get photos;

  /// Label for the 'Gallery' action
  ///
  /// In en, this message translates to:
  /// **'Gallery'**
  String get gallery;

  /// Label for the 'Photo' action
  ///
  /// In en, this message translates to:
  /// **'Photo'**
  String get photo;

  /// Label for the 'Today' message
  ///
  /// In en, this message translates to:
  /// **'Today'**
  String get today;

  /// Label for the 'Yesterday' message
  ///
  /// In en, this message translates to:
  /// **'Yesterday'**
  String get yesterday;

  /// Label for the 'Copy' button
  ///
  /// In en, this message translates to:
  /// **'Copy'**
  String get copy;

  /// Label for the 'Copied' message
  ///
  /// In en, this message translates to:
  /// **'Copied'**
  String get copied;

  /// Label for the 'Save' button
  ///
  /// In en, this message translates to:
  /// **'Save'**
  String get save;

  /// Label for the 'Share' button
  ///
  /// In en, this message translates to:
  /// **'Share'**
  String get share;

  /// Label for the 'Reply' button
  ///
  /// In en, this message translates to:
  /// **'Reply'**
  String get reply;

  /// Label for the 'Retry' button
  ///
  /// In en, this message translates to:
  /// **'Retry'**
  String get retry;

  /// Label for the 'Edit' button
  ///
  /// In en, this message translates to:
  /// **'Edit'**
  String get edit;

  /// Label for the 'Follow-ups' button
  ///
  /// In en, this message translates to:
  /// **'Follow-ups'**
  String get followUps;

  /// The error message to be shown when an error occurs
  ///
  /// In en, this message translates to:
  /// **'An error occurred, please retry the request'**
  String get errorMessage;

  /// The error message to be shown when there is no internet connection
  ///
  /// In en, this message translates to:
  /// **'No connection to the Internet'**
  String get noInternet;

  /// Label for the 'Give this chat a name' message
  ///
  /// In en, this message translates to:
  /// **'Give this chat a name'**
  String get giveThisChatName;

  /// Title for the settings screen
  ///
  /// In en, this message translates to:
  /// **'Settings'**
  String get settings;

  /// Label for the 'Language' setting
  ///
  /// In en, this message translates to:
  /// **'Language'**
  String get language;

  /// Label for the 'Brightness' setting
  ///
  /// In en, this message translates to:
  /// **'Brightness'**
  String get brightness;

  /// Label for the 'System' brightness option
  ///
  /// In en, this message translates to:
  /// **'System'**
  String get system;

  /// Label for the 'Light' brightness option
  ///
  /// In en, this message translates to:
  /// **'Light'**
  String get light;

  /// Label for the 'Dark' brightness option
  ///
  /// In en, this message translates to:
  /// **'Dark'**
  String get dark;

  /// Label for the 'App color' setting
  ///
  /// In en, this message translates to:
  /// **'App color'**
  String get appColor;

  /// Label for the 'Terms of service' button
  ///
  /// In en, this message translates to:
  /// **'Terms of service'**
  String get termsOfService;

  /// Label for the 'Privacy policy' button
  ///
  /// In en, this message translates to:
  /// **'Privacy policy'**
  String get privacyPolicy;

  /// Label for the 'Contact the developer' button
  ///
  /// In en, this message translates to:
  /// **'Contact the developer'**
  String get contact;

  /// Label for the 'Color change' title
  ///
  /// In en, this message translates to:
  /// **'App color change'**
  String get colorChangeTitle;

  /// Label for the 'Color change' message
  ///
  /// In en, this message translates to:
  /// **'App may stop to change the color. Restart the app to continue.'**
  String get colorChangeMessage;

  /// Label for the 'No chats' message
  ///
  /// In en, this message translates to:
  /// **'You don\'t have any chats.'**
  String get noChatsMessage;

  /// Title for the 'History' screen
  ///
  /// In en, this message translates to:
  /// **'History'**
  String get history;

  /// Label for the 'No chats' message
  ///
  /// In en, this message translates to:
  /// **'No chats'**
  String get noChats;

  /// Label for the 'Delete chat' title
  ///
  /// In en, this message translates to:
  /// **'Delete chat'**
  String get deleteChatTitle;

  /// Label for the 'Delete chat' message
  ///
  /// In en, this message translates to:
  /// **'Do you want to delete this chat?'**
  String get deleteChatMessage;

  /// Label for the 'Delete all chats' title
  ///
  /// In en, this message translates to:
  /// **'Delete all chats'**
  String get deleteAllChatsTitle;

  /// Label for the 'Delete all chats' message
  ///
  /// In en, this message translates to:
  /// **'Are you sure you want to delete all chats? This action is irreversible.'**
  String get deleteAllChatsMessage;

  /// Label for the 'Delete' button
  ///
  /// In en, this message translates to:
  /// **'Delete'**
  String get delete;

  /// Label for the 'Cancel' button
  ///
  /// In en, this message translates to:
  /// **'Cancel'**
  String get cancel;

  /// Label for the 'Rename chat' button
  ///
  /// In en, this message translates to:
  /// **'Rename chat'**
  String get rename;

  /// The first default prompt
  ///
  /// In en, this message translates to:
  /// **'What can you do?'**
  String get prompt1;

  /// The second default prompt
  ///
  /// In en, this message translates to:
  /// **'My mood is great!'**
  String get prompt2;

  /// The third default prompt
  ///
  /// In en, this message translates to:
  /// **'Tell me something surprising'**
  String get prompt3;

  /// The fourth default prompt
  ///
  /// In en, this message translates to:
  /// **'Translate this text into English...'**
  String get prompt4;

  /// The fifth default prompt
  ///
  /// In en, this message translates to:
  /// **'Check this for grammar mistakes...'**
  String get prompt5;

  /// The sixth default prompt
  ///
  /// In en, this message translates to:
  /// **'Explain this image to me...'**
  String get prompt6;

  /// The seventh default prompt
  ///
  /// In en, this message translates to:
  /// **'Recommend me a movie about adventure'**
  String get prompt7;

  /// The eighth default prompt
  ///
  /// In en, this message translates to:
  /// **'Give me the top 10 productivity books'**
  String get prompt8;

  /// The ninth default prompt
  ///
  /// In en, this message translates to:
  /// **'Summarize this article...'**
  String get prompt9;

  /// The tenth default prompt
  ///
  /// In en, this message translates to:
  /// **'Write an essay about technology'**
  String get prompt10;

  /// Label for the 'Pink' color
  ///
  /// In en, this message translates to:
  /// **'Pink'**
  String get pink;

  /// Label for the 'Red' color
  ///
  /// In en, this message translates to:
  /// **'Red'**
  String get red;

  /// Label for the 'Deep orange' color
  ///
  /// In en, this message translates to:
  /// **'Deep orange'**
  String get deepOrange;

  /// Label for the 'Orange' color
  ///
  /// In en, this message translates to:
  /// **'Orange'**
  String get orange;

  /// Label for the 'Amber' color
  ///
  /// In en, this message translates to:
  /// **'Amber'**
  String get amber;

  /// Label for the 'Yellow' color
  ///
  /// In en, this message translates to:
  /// **'Yellow'**
  String get yellow;

  /// Label for the 'Lime' color
  ///
  /// In en, this message translates to:
  /// **'Lime'**
  String get lime;

  /// Label for the 'Light green' color
  ///
  /// In en, this message translates to:
  /// **'Light green'**
  String get lightGreen;

  /// Label for the 'Green' color
  ///
  /// In en, this message translates to:
  /// **'Green'**
  String get green;

  /// Label for the 'Teal' color
  ///
  /// In en, this message translates to:
  /// **'Teal'**
  String get teal;

  /// Label for the 'Cyan' color
  ///
  /// In en, this message translates to:
  /// **'Cyan'**
  String get cyan;

  /// Label for the 'Light blue' color
  ///
  /// In en, this message translates to:
  /// **'Light blue'**
  String get lightBlue;

  /// Label for the 'Blue' color
  ///
  /// In en, this message translates to:
  /// **'Blue'**
  String get blue;

  /// Label for the 'Indigo' color
  ///
  /// In en, this message translates to:
  /// **'Indigo'**
  String get indigo;

  /// Label for the 'Purple' color
  ///
  /// In en, this message translates to:
  /// **'Purple'**
  String get purple;

  /// Label for the 'Deep purple' color
  ///
  /// In en, this message translates to:
  /// **'Deep purple'**
  String get deepPurple;

  /// Label for the 'Blue grey' color
  ///
  /// In en, this message translates to:
  /// **'Blue grey'**
  String get blueGrey;

  /// Label for the 'Brown' color
  ///
  /// In en, this message translates to:
  /// **'Brown'**
  String get brown;

  /// Label for the 'Grey' color
  ///
  /// In en, this message translates to:
  /// **'Grey'**
  String get grey;
}

class _AppLocalizationsDelegate
    extends LocalizationsDelegate<AppLocalizations> {
  const _AppLocalizationsDelegate();

  @override
  Future<AppLocalizations> load(Locale locale) {
    return lookupAppLocalizations(locale);
  }

  @override
  bool isSupported(Locale locale) =>
      <String>['en', 'ru', 'uz'].contains(locale.languageCode);

  @override
  bool shouldReload(_AppLocalizationsDelegate old) => false;
}

Future<AppLocalizations> lookupAppLocalizations(Locale locale) {
  // Lookup logic when only language code is specified.
  switch (locale.languageCode) {
    case 'en':
      return l10n_en.loadLibrary().then(
        (dynamic _) => l10n_en.AppLocalizationsEn(),
      );
    case 'ru':
      return l10n_ru.loadLibrary().then(
        (dynamic _) => l10n_ru.AppLocalizationsRu(),
      );
    case 'uz':
      return l10n_uz.loadLibrary().then(
        (dynamic _) => l10n_uz.AppLocalizationsUz(),
      );
  }

  throw FlutterError(
    'AppLocalizations.delegate failed to load unsupported locale "$locale". This is likely '
    'an issue with the localizations generation tool. Please file an issue '
    'on GitHub with a reproducible sample app and the gen-l10n configuration '
    'that was used.',
  );
}
