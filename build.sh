#!/bin/zsh

flutter build ipa --obfuscate --split-debug-info=debug_info/ios
flutter build appbundle --obfuscate --split-debug-info=debug_info/android
firebase crashlytics:symbols:upload --app=1:561132327391:android:92eab46f0d0ef4ba32f46a debug_info/android
