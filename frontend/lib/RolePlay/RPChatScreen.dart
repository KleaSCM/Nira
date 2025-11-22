import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:nira_frontend/WebSocketService.dart';
import 'roleplay_models.dart';
import 'roleplay_repository.dart';

class RPChatScreen extends StatefulWidget {
  final RPSession session;
  const RPChatScreen({super.key, required this.session});

  @override
  State<RPChatScreen> createState() => _RPChatScreenState();
}

class _RPChatScreenState extends State<RPChatScreen> {
  final List<Map<String, String>> _messages = [];
  final TextEditingController _ctrl = TextEditingController();
  final ScrollController _scroll = ScrollController();
  final WebSocketService _ws = WebSocketService();
  bool _started = false;
  final RolePlayRepository _repo = RolePlayRepository();
  List<RPCharacter> _characters = [];
  List<RPStoryCard> _cards = [];
  RPCharacter? _selectedCharacter;
  RPStoryCard? _selectedCard;

  @override
  void initState() {
    super.initState();
    // connect if not already connected (WebSocketService.connect should be idempotent)
    _ws.connect('ws://localhost:8080/ws');
    // Optionally listen for incoming messages if WebSocketService exposes a stream
    _ws.messageStream?.listen((msg) {
      setState(() {
        _messages.add({'sender': 'NIRA', 'text': msg.content});
      });
      _scrollToBottom();
    });
    _loadLibrary();
  }

  Future<void> _loadLibrary() async {
    final chars = await _repo.getCharacters();
    final cards = await _repo.getStoryCards();
    setState(() {
      _characters = chars;
      _cards = cards;
    });

    // Initialize selected character/card from session if present
    if (widget.session.characterId != null) {
      final match = chars.firstWhere((c) => c.id == widget.session.characterId, orElse: () => RPCharacter(id: null, name: '', description: ''));
      if (match.id != null && match.id! > 0) setState(() => _selectedCharacter = match);
    }
    if (widget.session.storyCardId != null) {
      final match = cards.firstWhere((c) => c.id == widget.session.storyCardId, orElse: () => RPStoryCard(id: null, title: '', content: '', world: ''));
      if (match.id != null && match.id! > 0) setState(() => _selectedCard = match);
    }
  }

  @override
  void dispose() {
    _ctrl.dispose();
    _scroll.dispose();
    super.dispose();
  }

  void _scrollToBottom() {
    Future.microtask(() {
      if (_scroll.hasClients) {
        _scroll.animateTo(_scroll.position.maxScrollExtent, duration: const Duration(milliseconds: 250), curve: Curves.easeOut);
      }
    });
  }

  void _startSession() {
    if (_started) return;
    final payload = {
      'type': 'rp_start',
      'session_id': widget.session.id,
      'session_name': widget.session.name,
      'created_at': widget.session.createdAt,
      'character': _selectedCharacter == null
          ? null
          : {
              'id': _selectedCharacter!.id,
              'name': _selectedCharacter!.name,
              'description': _selectedCharacter!.description,
            },
      'story_card': _selectedCard == null
          ? null
          : {
              'id': _selectedCard!.id,
              'title': _selectedCard!.title,
              'content': _selectedCard!.content,
              'world': _selectedCard!.world,
            },
    };
    _ws.sendRawJson(payload);
    setState(() {
      _messages.add({'sender': 'System', 'text': 'Session started: ${widget.session.name}'});
      _started = true;
    });
    _scrollToBottom();
  }

  void _sendMessage() {
    final text = _ctrl.text.trim();
    if (text.isEmpty) return;
    // ensure session started
    if (!_started) _startSession();

    final payload = {
      'type': 'rp_message',
      'session_id': widget.session.id,
      'text': text,
    };
    _ws.sendRawJson(payload);

    setState(() {
      _messages.add({'sender': 'You', 'text': text});
      _ctrl.clear();
    });
    _scrollToBottom();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(widget.session.name, style: GoogleFonts.quicksand()),
            if (widget.session.world != null && widget.session.world!.isNotEmpty)
              Text('World: ${widget.session.world}', style: GoogleFonts.quicksand(fontSize: 12)),
            if (widget.session.rules != null && widget.session.rules!.isNotEmpty)
              Text('Rules: ${widget.session.rules}', style: GoogleFonts.quicksand(fontSize: 12)),
            if (_selectedCharacter != null || _selectedCard != null)
              Text(
                '${_selectedCharacter != null ? 'Char: ${_selectedCharacter!.name}' : ''}${_selectedCharacter != null && _selectedCard != null ? ' · ' : ''}${_selectedCard != null ? 'Card: ${_selectedCard!.title}' : ''}',
                style: GoogleFonts.quicksand(fontSize: 12),
              ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () {
              // quick re-start / send session context again
              _started = false;
              _startSession();
            },
            child: Text('Restart', style: GoogleFonts.quicksand(color: Colors.white)),
          )
        ],
      ),
      body: Column(
        children: [
          // Drawer-like quick selector area
          SizedBox(
            height: 80,
            child: ListView(
              scrollDirection: Axis.horizontal,
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
              children: [
                Row(children: [Text('Characters: ', style: GoogleFonts.quicksand(fontWeight: FontWeight.w700))]),
                for (final c in _characters)
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 6.0),
                    child: ActionChip(
                      label: Text(c.name),
                      onPressed: () => setState(() => _selectedCharacter = c),
                      backgroundColor: _selectedCharacter?.id == c.id ? Colors.pink[100] : null,
                    ),
                  ),
                const SizedBox(width: 16),
                Row(children: [Text('Cards: ', style: GoogleFonts.quicksand(fontWeight: FontWeight.w700))]),
                for (final card in _cards)
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 6.0),
                    child: ActionChip(
                      label: Text(card.title),
                      onPressed: () => setState(() => _selectedCard = card),
                      backgroundColor: _selectedCard?.id == card.id ? Colors.purple[100] : null,
                    ),
                  ),
              ],
            ),
          ),
          Expanded(
            child: ListView.builder(
              controller: _scroll,
              padding: const EdgeInsets.all(12),
              itemCount: _messages.length,
              itemBuilder: (c, i) {
                final m = _messages[i];
                final isUser = m['sender'] == 'You';
                return Align(
                  alignment: isUser ? Alignment.centerRight : Alignment.centerLeft,
                  child: Container(
                    margin: const EdgeInsets.symmetric(vertical: 6),
                    padding: const EdgeInsets.all(12),
                    decoration: BoxDecoration(
                      color: isUser ? Colors.pink[50] : Colors.grey[200],
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        if (m['sender'] != null) Text('${m['sender']}', style: GoogleFonts.quicksand(fontSize: 12, color: Colors.grey[700])),
                        const SizedBox(height: 4),
                        Text(m['text'] ?? '', style: GoogleFonts.quicksand()),
                      ],
                    ),
                  ),
                );
              },
            ),
          ),
          SafeArea(
            child: Padding(
              padding: const EdgeInsets.fromLTRB(12, 8, 12, 12),
              child: Row(
                children: [
                  Expanded(
                    child: TextField(
                      controller: _ctrl,
                      decoration: InputDecoration(
                        hintText: 'Send message to session...',
                        border: OutlineInputBorder(borderRadius: BorderRadius.circular(10)),
                        contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
                      ),
                      minLines: 1,
                      maxLines: 4,
                    ),
                  ),
                  const SizedBox(width: 8),
                  ElevatedButton(
                    onPressed: _sendMessage,
                    child: const Icon(Icons.send),
                  )
                ],
              ),
            ),
          )
        ],
      ),
    );
  }
}