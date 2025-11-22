/**
 * NIRA Frontend - Main entry point.
 *
 * Flutter application providing the GUI for interacting with NIRA.
 * Handles chat interface, tool logs, settings, and WebSocket communication
 * with the Go backend.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 * File: main.dart
 * Description: Flutter app initialization and main widget.
 */

import 'package:flutter/material.dart';
import 'package:nira_frontend/App.dart';

void main() {
	runApp(const NiraApp());
}

