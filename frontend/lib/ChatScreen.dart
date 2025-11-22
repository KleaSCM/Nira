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
import 'package:flutter/services.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:nira_frontend/WebSocketService.dart';

const Color primaryColor = Color(0xFFFFB3D9);  // Bright pink
const Color secondaryColor = Color(0xFFFFF0F5); // Lavender blush
const Color accentColor = Color(0xFFFF69B4);   // Hot pink
const Color softPurple = Color(0xFFE6B3FF);    // Soft purple
const Color textColor = Color(0xFF8B4789);     // Deep purple
const Color userBubbleColor = Color(0xFFFFE4F0); // Pale pink
const Color niraBubbleColor = Color(0xFFFFF8DC); // Cornsilk
const Color errorColor = Color(0xFFFFCCE5);     // Light pink red

class ChatScreen extends StatefulWidget {
  const ChatScreen({super.key});

  @override
  State<ChatScreen> createState() => _ChatScreenState();
}

class _ChatScreenState extends State<ChatScreen> with SingleTickerProviderStateMixin {
  final List<Map<String, String>> Messages = [];
  final TextEditingController MessageController = TextEditingController();
  final WebSocketService WsService = WebSocketService();
  final ScrollController _scrollController = ScrollController();
  final FocusNode _textFocusNode = FocusNode();
  int CurrentAssistantIndex = -1;
  late AnimationController _animationController;
  bool _isComposing = false;

  @override
  void initState() {
    super.initState();
    _animationController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 300),
    );
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
        _scrollToBottom();
      });
    }
  }

  @override
  void dispose() {
    MessageController.dispose();
    WsService.disconnect();
    _animationController.dispose();
    _textFocusNode.dispose();
    super.dispose();
  }

  void _sendMessage() {
    final text = MessageController.text.trim();
    if (text.isEmpty) return;

    setState(() {
      Messages.add({
        'text': text,
        'sender': 'You',
      });
      CurrentAssistantIndex = -1;
      MessageController.clear();
      _isComposing = false;
    });

    _scrollToBottom();
    WsService.sendMessage(text);
  }

  void _scrollToBottom() {
    Future.delayed(const Duration(milliseconds: 100), () {
      if (_scrollController.hasClients) {
        _scrollController.animateTo(
          _scrollController.position.maxScrollExtent,
          duration: const Duration(milliseconds: 300),
          curve: Curves.easeOut,
        );
      }
    });
  }

  Widget _buildAvatar(String sender) {
    final isNira = sender == 'NIRA';
    return Container(
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        gradient: LinearGradient(
          colors: isNira 
            ? [Color(0xFFFF69B4), Color(0xFFFF1493)]
            : [Color(0xFFE6B3FF), Color(0xFFDA70D6)],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        boxShadow: [
          BoxShadow(
            color: (isNira ? Colors.pink : Colors.purple).withOpacity(0.4),
            blurRadius: 8,
            offset: const Offset(0, 3),
          ),
        ],
      ),
      child: CircleAvatar(
        backgroundColor: Colors.transparent,
        radius: 22,
        child: Text(
          isNira ? '💖' : '✨',
          style: const TextStyle(fontSize: 24),
        ),
      ),
    );
  }

  Widget _buildMessageBubble(Map<String, String> message, int index) {
    final isUser = message['sender'] == 'You';
    final isError = message['sender'] == 'System';
    
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 6.0, horizontal: 12.0),
      child: Row(
        mainAxisAlignment: isUser ? MainAxisAlignment.end : MainAxisAlignment.start,
        crossAxisAlignment: CrossAxisAlignment.end,
        children: [
          if (!isUser) _buildAvatar(message['sender']!),
          const SizedBox(width: 10),
          Flexible(
            child: Container(
              padding: const EdgeInsets.symmetric(vertical: 12, horizontal: 18),
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  colors: isError 
                    ? [errorColor, errorColor]
                    : isUser 
                      ? [Color(0xFFFFE4F0), Color(0xFFFFD6E8)]
                      : [Color(0xFFFFF8DC), Color(0xFFFFF0DB)],
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                ),
                borderRadius: BorderRadius.only(
                  topLeft: const Radius.circular(24),
                  topRight: const Radius.circular(24),
                  bottomLeft: Radius.circular(isUser ? 24 : 6),
                  bottomRight: Radius.circular(isUser ? 6 : 24),
                ),
                boxShadow: [
                  BoxShadow(
                    color: Colors.pink.withOpacity(0.2),
                    blurRadius: 8,
                    offset: const Offset(0, 3),
                  ),
                ],
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  if (!isUser && !isError)
                    Padding(
                      padding: const EdgeInsets.only(bottom: 4),
                      child: Text(
                        '${message['sender']!} 💕',
                        style: GoogleFonts.comfortaa(
                          fontWeight: FontWeight.w700,
                          color: Color(0xFFFF1493),
                          fontSize: 13,
                        ),
                      ),
                    ),
                  Text(
                    message['text']!,
                    style: GoogleFonts.quicksand(
                      color: isError ? Colors.red[800] : textColor,
                      fontSize: 15,
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                ],
              ),
            ).animate().fadeIn(duration: 300.ms).slideY(
              begin: 0.2,
              end: 0,
              curve: Curves.easeOutCubic,
            ).scale(
              begin: const Offset(0.9, 0.9),
              end: const Offset(1, 1),
              curve: Curves.easeOutBack,
            ),
          ),
          if (isUser) const SizedBox(width: 10),
          if (isUser) _buildAvatar('You'),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: secondaryColor,
      appBar: PreferredSize(
        preferredSize: const Size.fromHeight(80),
        child: Container(
          decoration: BoxDecoration(
            gradient: LinearGradient(
              colors: [
                Color(0xFFFF69B4),
                Color(0xFFFF1493),
                Color(0xFFE6B3FF),
              ],
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
            ),
            boxShadow: [
              BoxShadow(
                color: Colors.pink.withOpacity(0.4),
                blurRadius: 20,
                offset: const Offset(0, 5),
              ),
            ],
          ),
          child: SafeArea(
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
              child: Row(
                children: [
                  Container(
                    padding: const EdgeInsets.all(8),
                    decoration: BoxDecoration(
                      color: Colors.white.withOpacity(0.9),
                      shape: BoxShape.circle,
                      boxShadow: [
                        BoxShadow(
                          color: Colors.pink.withOpacity(0.3),
                          blurRadius: 8,
                        ),
                      ],
                    ),
                    child: const Icon(
                      Icons.favorite,
                      color: Color(0xFFFF1493),
                      size: 28,
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Text(
                          'NIRA Assistant ✨',
                          style: GoogleFonts.pacifico(
                            color: Colors.white,
                            fontSize: 26,
                            shadows: [
                              Shadow(
                                color: Colors.black.withOpacity(0.2),
                                offset: const Offset(0, 2),
                                blurRadius: 4,
                              ),
                            ],
                          ),
                        ),
                        Text(
                          'lkhjsdf',
                          style: GoogleFonts.quicksand(
                            color: Colors.white.withOpacity(0.9),
                            fontSize: 13,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                      ],
                    ),
                  ),
                  Container(
                    padding: const EdgeInsets.all(8),
                    decoration: BoxDecoration(
                      color: Colors.white.withOpacity(0.9),
                      shape: BoxShape.circle,
                    ),
                    child: const Icon(
                      Icons.more_vert,
                      color: Color(0xFFFF1493),
                      size: 24,
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
      body: Column(
        children: [
          Expanded(
            child: Container(
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  begin: Alignment.topCenter,
                  end: Alignment.bottomCenter,
                  colors: [
                    Color(0xFFFFF0F5),
                    Color(0xFFFFE4F0),
                    Color(0xFFFFF8F0),
                  ],
                ),
              ),
              child: ListView.builder(
                controller: _scrollController,
                padding: const EdgeInsets.symmetric(vertical: 12.0),
                itemCount: Messages.length,
                itemBuilder: (context, index) {
                  return _buildMessageBubble(Messages[index], index);
                },
              ),
            ),
          ),
          Container(
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: const BorderRadius.vertical(top: Radius.circular(30)),
              boxShadow: [
                BoxShadow(
                  color: Colors.pink.withOpacity(0.15),
                  blurRadius: 20,
                  offset: const Offset(0, -5),
                ),
              ],
            ),
            padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
            child: SafeArea(
              child: Column(
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                    children: [
                      _buildToolButton(Icons.image_outlined, 'Photo', Colors.pink),
                      _buildToolButton(Icons.attachment_outlined, 'File', Colors.purple),
                      _buildToolButton(Icons.camera_alt_outlined, 'Camera', Colors.deepPurple),
                      _buildToolButton(Icons.location_on_outlined, 'Location', Colors.pinkAccent),
                    ],
                  ),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      Container(
                        decoration: BoxDecoration(
                          gradient: LinearGradient(
                            colors: [Color(0xFFFF69B4), Color(0xFFFF1493)],
                          ),
                          shape: BoxShape.circle,
                        ),
                        child: IconButton(
                          icon: const Icon(
                            Icons.emoji_emotions,
                            color: Colors.white,
                            size: 26,
                          ),
                          onPressed: () {
                            // TODO: Add emoji picker
                          },
                        ),
                      ),
                      const SizedBox(width: 8),
                      Expanded(
                        child: Container(
                          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
                          decoration: BoxDecoration(
                            gradient: LinearGradient(
                              colors: [
                                Color(0xFFFFF0F5),
                                Color(0xFFFFE4F0),
                              ],
                            ),
                            borderRadius: BorderRadius.circular(28),
                            border: Border.all(
                              color: Color(0xFFFFB3D9),
                              width: 2,
                            ),
                          ),
                          child: KeyboardListener(
                            focusNode: FocusNode(),
                            onKeyEvent: (KeyEvent event) {
                              if (event is KeyDownEvent) {
                                // Enter without Shift sends message
                                if (event.logicalKey == LogicalKeyboardKey.enter &&
                                    !HardwareKeyboard.instance.isShiftPressed) {
                                  _sendMessage();
                                }
                              }
                            },
                            child: TextField(
                              controller: MessageController,
                              focusNode: _textFocusNode,
                              onChanged: (text) {
                                setState(() {
                                  _isComposing = text.trim().isNotEmpty;
                                });
                              },
                              style: GoogleFonts.quicksand(
                                color: textColor,
                                fontSize: 15,
                                fontWeight: FontWeight.w500,
                              ),
                              decoration: InputDecoration(
                                hintText: 'Type something..',
                                hintStyle: GoogleFonts.quicksand(
                                  color: Colors.pink[300],
                                  fontSize: 15,
                                ),
                                border: InputBorder.none,
                                contentPadding: const EdgeInsets.symmetric(vertical: 12),
                              ),
                              maxLines: 4,
                              minLines: 1,
                              textInputAction: TextInputAction.newline,
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(width: 8),
                      AnimatedSwitcher(
                        duration: const Duration(milliseconds: 250),
                        transitionBuilder: (child, animation) {
                          return ScaleTransition(
                            scale: animation,
                            child: RotationTransition(
                              turns: animation,
                              child: child,
                            ),
                          );
                        },
                        child: _isComposing
                            ? Container(
                                key: const ValueKey('send'),
                                decoration: BoxDecoration(
                                  gradient: LinearGradient(
                                    colors: [
                                      Color(0xFFFF69B4),
                                      Color(0xFFFF1493),
                                    ],
                                    begin: Alignment.topLeft,
                                    end: Alignment.bottomRight,
                                  ),
                                  shape: BoxShape.circle,
                                  boxShadow: [
                                    BoxShadow(
                                      color: Colors.pink.withOpacity(0.5),
                                      blurRadius: 12,
                                      offset: const Offset(0, 4),
                                    ),
                                  ],
                                ),
                                child: IconButton(
                                  icon: const Icon(
                                    Icons.send_rounded,
                                    color: Colors.white,
                                    size: 26,
                                  ),
                                  onPressed: _sendMessage,
                                ),
                              )
                            : Container(
                                key: const ValueKey('mic'),
                                decoration: BoxDecoration(
                                  gradient: LinearGradient(
                                    colors: [Color(0xFFE6B3FF), Color(0xFFDA70D6)],
                                  ),
                                  shape: BoxShape.circle,
                                ),
                                child: IconButton(
                                  icon: const Icon(
                                    Icons.mic_none_rounded,
                                    color: Colors.white,
                                    size: 26,
                                  ),
                                  onPressed: () {
                                    // TODO: Add voice input
                                  },
                                ),
                              ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildToolButton(IconData icon, String label, Color color) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        Container(
          decoration: BoxDecoration(
            color: color.withOpacity(0.1),
            shape: BoxShape.circle,
            border: Border.all(
              color: color.withOpacity(0.3),
              width: 2,
            ),
          ),
          child: IconButton(
            icon: Icon(icon, color: color, size: 22),
            onPressed: () {
              // TODO: Implement tool functionality
            },
          ),
        ),
        const SizedBox(height: 4),
        Text(
          label,
          style: GoogleFonts.quicksand(
            fontSize: 11,
            fontWeight: FontWeight.w600,
            color: color,
          ),
        ),
      ],
    );
  }
}