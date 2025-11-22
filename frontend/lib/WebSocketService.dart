/**
 * WebSocket service module.
 *
 * Manages WebSocket connection to the NIRA backend, handles message
 * sending and receiving, and provides streaming support for responses.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 * File: WebSocketService.dart
 * Description: WebSocket connection and message handling.
 */

import 'dart:convert';
import 'package:web_socket_channel/web_socket_channel.dart';

enum MessageType {
	user,
	assistant,
	system,
	error,
	chunk,
}

class WSMessage {
	final MessageType type;
	final String content;
	final String? id;

	WSMessage({
		required this.type,
		required this.content,
		this.id,
	});

	Map<String, dynamic> toJson() {
		return {
			'type': type.name,
			'content': content,
			if (id != null) 'id': id,
		};
	}

	factory WSMessage.fromJson(Map<String, dynamic> json) {
		return WSMessage(
			type: MessageType.values.firstWhere(
				(e) => e.name == json['type'],
				orElse: () => MessageType.system,
			),
			content: json['content'] ?? '',
			id: json['id'],
		);
	}
}

class WebSocketService {
	WebSocketChannel? Channel;
	String CurrentMessage = '';
	bool IsConnected = false;

	WebSocketChannel? connect(String url) {
		try {
			Channel = WebSocketChannel.connect(Uri.parse(url));
			IsConnected = true;
			return Channel;
		} catch (e) {
			IsConnected = false;
			return null;
		}
	}

	void disconnect() {
		Channel?.sink.close();
		IsConnected = false;
	}

	void sendMessage(String content) {
		if (Channel == null || !IsConnected) {
			return;
		}

		final msg = WSMessage(
			type: MessageType.user,
			content: content,
		);

		Channel!.sink.add(jsonEncode(msg.toJson()));
	}

	Stream<WSMessage>? get messageStream {
		if (Channel == null) {
			return null;
		}
		return Channel!.stream.map((data) {
			final json = jsonDecode(data);
			return WSMessage.fromJson(json);
		});
	}
}

