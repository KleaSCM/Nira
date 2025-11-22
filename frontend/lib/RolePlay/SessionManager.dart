import 'package:flutter/material.dart';
import 'roleplay_models.dart';
import 'roleplay_repository.dart';
import 'RPChatScreen.dart';
import 'CharacterList.dart';
import 'StoryCardList.dart';

class SessionManager extends StatefulWidget {
  final String? world;
  const SessionManager({super.key, this.world});

  @override
  State<SessionManager> createState() => _SessionManagerState();
}

class _SessionManagerState extends State<SessionManager> {
  final List<RPSession> _sessions = [];
  final RolePlayRepository _repo = RolePlayRepository();
  RPCharacter? _selectedCharacter;
  RPStoryCard? _selectedCard;
  final TextEditingController _nameController = TextEditingController();
  final TextEditingController _rulesController = TextEditingController();
  List<RPCharacter> _allCharacters = [];
  List<RPStoryCard> _allCards = [];

  @override
  void initState() {
    super.initState();
    _load();
    _loadLibrary();
  }

  Future<void> _loadLibrary() async {
    final chars = await _repo.getCharacters();
    final cards = await _repo.getStoryCards();
    setState(() {
      _allCharacters = chars;
      _allCards = cards;
    });

    // Preselect story card matching the provided world, if any
    if (widget.world != null && widget.world!.isNotEmpty) {
      final match = cards.firstWhere((c) => c.world == widget.world, orElse: () => RPStoryCard(id: null, title: '', content: '', world: ''));
      if (match.id != null && match.id! > 0) {
        setState(() => _selectedCard = match);
      }
    }

    // If there's exactly one character in the library, preselect it
    if (chars.length == 1) setState(() => _selectedCharacter = chars.first);
  }

  Future<void> _load() async {
    final s = await _repo.getSessions();
    setState(() => _sessions
      ..clear()
      ..addAll(s));
  }
  Future<void> _createSession() async {
    final now = DateTime.now().millisecondsSinceEpoch;
    final name = _nameController.text.trim().isNotEmpty ? _nameController.text.trim() : 'Session ${_sessions.length + 1}';
    final world = widget.world ?? _selectedCard?.world ?? '';
    final rules = _rulesController.text.trim();
    final s = RPSession(
      name: name,
      world: world,
      characterId: _selectedCharacter?.id,
      storyCardId: _selectedCard?.id,
      rules: rules,
      createdAt: now,
    );
    final id = await _repo.insertSession(s);
    // create session object with returned id and navigate straight to chat
    final created = RPSession(id: id, name: name, world: world, characterId: _selectedCharacter?.id, storyCardId: _selectedCard?.id, rules: rules, createdAt: now);
    await _load();
    // open RP chat for the new session
    if (!mounted) return;
    await Navigator.push(context, MaterialPageRoute(builder: (_) => RPChatScreen(session: created)));
    // show a small confirmation after returning from chat (or immediately)
    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Session created and opened')));
    }
  }

  Future<void> _onStartPressed() async {
    // basic validation
    final name = _nameController.text.trim().isNotEmpty ? _nameController.text.trim() : 'Session ${_sessions.length + 1}';
    final world = widget.world ?? _selectedCard?.world ?? '';
    final charName = _selectedCharacter?.name ?? 'None';
    final cardTitle = _selectedCard?.title ?? 'None';
    final rules = _rulesController.text.trim();

    final confirmed = await showDialog<bool>(
      context: context,
      builder: (c) => AlertDialog(
        title: const Text('Create Session?'),
        content: SingleChildScrollView(
          child: ListBody(children: [
            Text('Name: $name'),
            if (world.isNotEmpty) Text('World: $world'),
            Text('Character: $charName'),
            Text('Story Card: $cardTitle'),
            if (rules.isNotEmpty) const SizedBox(height: 8),
            if (rules.isNotEmpty) Text('Rules: $rules', maxLines: 5, overflow: TextOverflow.ellipsis),
          ]),
        ),
        actions: [
          TextButton(onPressed: () => Navigator.of(c).pop(false), child: const Text('Cancel')),
          ElevatedButton(onPressed: () => Navigator.of(c).pop(true), child: const Text('Create')),
        ],
      ),
    );

    if (confirmed == true) {
      await _createSession();
    }
  }

  @override
  void dispose() {
    _nameController.dispose();
    _rulesController.dispose();
    super.dispose();
  }

  void _openSession(RPSession s) {
    Navigator.push(context, MaterialPageRoute(builder: (_) => RPChatScreen(session: s)));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Session Manager'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            TextField(
              controller: _nameController,
              decoration: const InputDecoration(labelText: 'Session Name', hintText: 'Enter a name or leave default'),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _rulesController,
              decoration: const InputDecoration(labelText: 'Rules (optional)', hintText: 'Session-specific rules or notes'),
              maxLines: 3,
            ),
            const SizedBox(height: 12),
            Row(
              children: [
                Expanded(
                  child: ElevatedButton.icon(
                    onPressed: () async {
                      final res = await Navigator.push(context, MaterialPageRoute(builder: (_) => const CharacterList()));
                      if (res != null && res is RPCharacter) {
                        setState(() => _selectedCharacter = res);
                      }
                    },
                    icon: const Icon(Icons.person),
                    label: const Text('Pick Character'),
                  ),
                ),
                const SizedBox(width: 8),
                Expanded(
                  child: ElevatedButton.icon(
                    onPressed: () async {
                      final res = await Navigator.push(context, MaterialPageRoute(builder: (_) => StoryCardList(world: widget.world)));
                      if (res != null && res is RPStoryCard) {
                        setState(() => _selectedCard = res);
                      }
                    },
                    icon: const Icon(Icons.menu_book),
                    label: const Text('Pick Story Card'),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),
            if (_selectedCharacter != null) ListTile(title: Text('Character: ${_selectedCharacter!.name}'), subtitle: Text(_selectedCharacter!.description)),
            if (_selectedCard != null) ListTile(title: Text('Story Card: ${_selectedCard!.title}'), subtitle: Text(_selectedCard!.world.isNotEmpty ? _selectedCard!.world : 'No world')),
            const SizedBox(height: 12),
            Row(
              children: [
                Expanded(child: ElevatedButton.icon(onPressed: _onStartPressed, icon: const Icon(Icons.play_arrow), label: const Text('Start Session'))),
                const SizedBox(width: 8),
                OutlinedButton.icon(
                  onPressed: () {
                    _nameController.clear();
                    _rulesController.clear();
                    setState(() {
                      _selectedCharacter = null;
                      _selectedCard = null;
                    });
                  },
                  icon: const Icon(Icons.clear),
                  label: const Text('Reset'),
                ),
              ],
            ),
            const SizedBox(height: 12),
            const Divider(),
            const SizedBox(height: 8),
            Expanded(
              child: ListView.builder(
                itemCount: _sessions.length,
                itemBuilder: (context, i) => ListTile(
                  title: Text(_sessions[i].name),
                  subtitle: Text(_sessions[i].world ?? ''),
                  trailing: const Icon(Icons.play_arrow),
                  onTap: () => _openSession(_sessions[i]),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}