/**
 * Chat interface screen.
 *
 * Provides the main chat UI for interacting with NIRA, including message
 * display, input field, and WebSocket connection management.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 * File: ChatScreen.dart
 * Description: Main chat interface component.
 */

import 'package:flutter/material.dart';
import 'package:nira_frontend/WebSocketService.dart';

class ChatScreen extends StatefulWidget {
	const ChatScreen({super.key});

	@override
	State<ChatScreen> createState() => _ChatScreenState();
}

class _ChatScreenState extends State<ChatScreen> {
	final List<Map<String, String>> Messages = [];
	final TextEditingController MessageController = TextEditingController();
	final WebSocketService WsService = WebSocketService();
	int CurrentAssistantIndex = -1;

	@override
	void initState() {
		super.initState();
		_connectWebSocket();
	}

	void _connectWebSocket() {
		final channel = WsService.connect('ws://localhost:8080/ws');
		if (channel != null) {
			WsService.messageStream?.listen((msg) {
				setState(() {
					if (msg.type == MessageType.chunk) {
						if (CurrentAssistantIndex == -1) {
							CurrentAssistantIndex = Messages.length;
							Messages.add({
								'text': msg.content,
								'sender': 'NIRA',
							});
						} else {
							Messages[CurrentAssistantIndex]['text'] =
								(Messages[CurrentAssistantIndex]['text'] ?? '') + msg.content;
						}
					} else if (msg.type == MessageType.assistant) {
						CurrentAssistantIndex = -1;
					} else if (msg.type == MessageType.error) {
						Messages.add({
							'text': 'Error: ${msg.content}',
							'sender': 'System',
						});
					}
				});
			});
		}
	}

	@override
	void dispose() {
		MessageController.dispose();
		WsService.disconnect();
		super.dispose();
	}

	void _sendMessage() {
		final text = MessageController.text.trim();
		if (text.isEmpty) {
			return;
		}

		setState(() {
			Messages.add({
				'text': text,
				'sender': 'You',
			});
			CurrentAssistantIndex = -1;
			MessageController.clear();
		});

		WsService.sendMessage(text);
	}

	@override
	Widget build(BuildContext context) {
		return Scaffold(
			appBar: AppBar(
				title: const Text('NIRA'),
			),
			body: Column(
				children: [
					Expanded(
						child: ListView.builder(
							itemCount: Messages.length,
							itemBuilder: (context, index) {
								final msg = Messages[index];
								final isUser = msg['sender'] == 'You';
								return Align(
									alignment: isUser ? Alignment.centerRight : Alignment.centerLeft,
									child: Container(
										margin: const EdgeInsets.symmetric(vertical: 4, horizontal: 8),
										padding: const EdgeInsets.all(12),
										decoration: BoxDecoration(
											color: isUser ? Colors.blue : Colors.grey[300],
											borderRadius: BorderRadius.circular(12),
										),
										child: Column(
											crossAxisAlignment: CrossAxisAlignment.start,
											children: [
												Text(
													msg['sender'] ?? '',
													style: TextStyle(
														fontSize: 12,
														fontWeight: FontWeight.bold,
														color: isUser ? Colors.white70 : Colors.black54,
													),
												),
												const SizedBox(height: 4),
												Text(
													msg['text'] ?? '',
													style: TextStyle(
														color: isUser ? Colors.white : Colors.black87,
													),
												),
											],
										),
									),
								);
							},
						),
					),
					Padding(
						padding: const EdgeInsets.all(8.0),
						child: Row(
							children: [
								Expanded(
									child: TextField(
										controller: MessageController,
										decoration: const InputDecoration(
											hintText: 'Type a message...',
											border: OutlineInputBorder(),
										),
										onSubmitted: (_) => _sendMessage(),
									),
								),
								const SizedBox(width: 8),
								IconButton(
									icon: const Icon(Icons.send),
									onPressed: _sendMessage,
								),
							],
						),
					),
				],
			),
		);
	}
}

