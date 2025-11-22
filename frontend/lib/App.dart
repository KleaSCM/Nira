/**
 * Root application widget.
 *
 * Sets up the Material app theme and navigation structure for NIRA.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 * File: App.dart
 * Description: Main application widget with theme configuration.
 */

import 'package:flutter/material.dart';
import 'package:nira_frontend/ChatScreen.dart';

class NiraApp extends StatelessWidget {
	const NiraApp({super.key});

	@override
	Widget build(BuildContext context) {
		return MaterialApp(
			title: 'NIRA',
			theme: ThemeData(
				useMaterial3: true,
				colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
			),
			home: const ChatScreen(),
		);
	}
}

