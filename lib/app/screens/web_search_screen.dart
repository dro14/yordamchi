// import 'dart:io';

// import 'package:flutter/cupertino.dart';
// import 'package:flutter/material.dart'
//     show LinearProgressIndicator, SelectionArea;
// import 'package:path_provider/path_provider.dart';
// import 'package:share_plus/share_plus.dart';
// import 'package:intl/intl.dart';

// import 'package:yordamchi/core/core.dart';

// class WebSearchScreen extends StatefulWidget {
//   const WebSearchScreen({super.key});

//   static Future<T?> route<T>(BuildContext context) {
//     return Navigator.of(context, rootNavigator: true).push<T>(
//       CupertinoPageRoute<T>(
//         builder: (context) {
//           return const WebSearchScreen();
//         },
//       ),
//     );
//   }

//   @override
//   State<WebSearchScreen> createState() => _WebSearchScreenState();
// }

// class _WebSearchScreenState extends State<WebSearchScreen> {
//   bool _isLoading = false;

//   Future<void> _search(String query) async {
//     query = query.trim();
//     if (query.isEmpty) return;
//     setState(() => _isLoading = true);
//     final response = await getSearchResults(query);
//     setState(() => _isLoading = false);
//     if (mounted) {
//       showCupertinoModalPopup(
//         context: context,
//         builder: (context) {
//           return CupertinoPopupSurface(
//             blurSigma: 0.0,
//             child: Container(
//               height: 500.0,
//               color: surfaceContainer,
//               child: SingleChildScrollView(
//                 padding: const EdgeInsets.all(16.0),
//                 child: SelectionArea(child: Text(response)),
//               ),
//             ),
//           );
//         },
//       );
//     }
//   }

//   Future<void> _sharePageSource() async {
//     final result = await getPageSource();
//     if (result.isNotEmpty) {
//       final appDir = await getApplicationDocumentsDirectory();
//       final timestamp = DateFormat('yyyyMMdd_HHmmss').format(DateTime.now());
//       final filename = 'search_results_$timestamp.html';
//       final file = File('${appDir.path}/$filename');
//       await file.writeAsString(result);
//       await SharePlus.instance.share(ShareParams(files: [XFile(file.path)]));
//     }
//   }

//   @override
//   Widget build(BuildContext context) {
//     return CupertinoPageScaffold(
//       navigationBar: CupertinoNavigationBar(
//         padding: const EdgeInsetsDirectional.only(end: 12.0),
//         middle: CupertinoTextField(
//           placeholder: 'Google',
//           textInputAction: TextInputAction.go,
//           onSubmitted: _search,
//         ),
//         trailing: CupertinoButton(
//           padding: EdgeInsets.zero,
//           sizeStyle: CupertinoButtonSize.small,
//           onPressed: _sharePageSource,
//           child: Icon(MyIcons.share, size: 24.0),
//         ),
//         bottom: _isLoading
//             ? const PreferredSize(
//                 preferredSize: Size.fromHeight(4.0),
//                 child: LinearProgressIndicator(),
//               )
//             : null,
//       ),
//       child: const SafeArea(
//         child: Center(
//           child: Column(
//             mainAxisAlignment: MainAxisAlignment.center,
//             children: [
//               Icon(CupertinoIcons.globe, size: 64.0),
//               SizedBox(height: 16.0),
//               Text('Enter a query to search the web', style: Fonts.bodyLarge),
//             ],
//           ),
//         ),
//       ),
//     );
//   }
// }
