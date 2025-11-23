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

import 'dart:async';
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
    // Singleton
    static final WebSocketService instance = WebSocketService._internal();
    factory WebSocketService() => instance;
    WebSocketService._internal();

    /// Send a tool call to the backend (e.g. {"name": "web_search", "arguments": {"query": "..."}})
    void sendToolCall(Map<String, dynamic> toolCall) {
        if (Channel == null || !IsConnected) return;
        Channel!.sink.add(jsonEncode(toolCall));
    }

    WebSocketChannel? Channel;
    late Stream<dynamic> _incoming; // broadcast stream of raw frames
    String CurrentMessage = '';
    bool IsConnected = false;

    // Pending RPC-style tool calls (by ID)
    final Map<String, Completer<WSMessage>> _pending = {};
    StreamSubscription<WSMessage>? _routerSub;

    WebSocketChannel? connect(String url) {
        try {
            // If already connected, return existing channel
            if (Channel != null && IsConnected) {
                return Channel;
            }
            Channel = WebSocketChannel.connect(Uri.parse(url));
            IsConnected = true;
            _incoming = Channel!.stream.asBroadcastStream();
            // Ensure router is listening to dispatch responses to pending completers
            _ensureRouter();
            return Channel;
        } catch (e) {
            IsConnected = false;
            return null;
        }
    }

 void disconnect() {
        Channel?.sink.close();
        IsConnected = false;
        _routerSub?.cancel();
        _routerSub = null;
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

	/// Send a typed message (allows sending assistant/system/error types)
	void sendTyped(MessageType type, String content) {
		if (Channel == null || !IsConnected) return;
		final msg = WSMessage(type: type, content: content);
		Channel!.sink.add(jsonEncode(msg.toJson()));
	}

	/// Send a raw JSON-like object directly over the socket.
	/// Use this when you need to send structured events (e.g. rp_start) that the backend
	/// expects as top-level JSON objects.
	void sendRawJson(Map<String, dynamic> map) {
		if (Channel == null || !IsConnected) return;
		Channel!.sink.add(jsonEncode(map));
	}

    void _ensureRouter() {
        _routerSub ??= messageStream?.listen((msg) {
            final id = msg.id;
            if (id != null && _pending.containsKey(id)) {
                // Resolve pending completer
                final c = _pending.remove(id)!;
                c.complete(msg);
            }
        });
    }

    // Broadcast message stream for UI subscribers
    Stream<WSMessage>? get messageStream {
        if (Channel == null) {
            return null;
        }
        return _incoming.map((data) {
            final json = jsonDecode(data);
            return WSMessage.fromJson(json);
        });
    }

    /// RPC-style helper: call a tool and await a single structured JSON response.
    /// Sends a direct tool call with an ID and requests silent handling.
    /// Returns the decoded JSON object/array from the tool result.
    Future<dynamic> callToolJson(String name, Map<String, dynamic> arguments, {Duration timeout = const Duration(seconds: 10)}) async {
        if (!IsConnected) {
            // Attempt default connection
            connect('ws://localhost:8080/ws');
        }
        if (Channel == null || !IsConnected) {
            throw Exception('WebSocket not connected');
        }
        // Prepare call ID and arguments
        final id = DateTime.now().microsecondsSinceEpoch.toString();
        final args = Map<String, dynamic>.from(arguments);
        args['_silent'] = true; // instruct backend not to stream side messages
        final payload = {
            'id': id,
            'name': name,
            'arguments': args,
        };
        final completer = Completer<WSMessage>();
        _pending[id] = completer;
        sendToolCall(payload);

        final msg = await completer.future.timeout(timeout, onTimeout: () {
            _pending.remove(id);
            throw TimeoutException('Tool call timed out: $name');
        });
        if (msg.type == MessageType.error) {
            throw Exception(msg.content);
        }
        // Content is a JSON string produced by backend for silent replies
        try {
            final decoded = jsonDecode(msg.content);
            return decoded;
        } catch (e) {
            // If not JSON, return raw string
            return msg.content;
        }
    }
}

