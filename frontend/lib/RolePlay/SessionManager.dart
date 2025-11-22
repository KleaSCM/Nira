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
  List<RPCharacter> _selectedCharacters = [];
  List<RPStoryCard> _selectedCards = [];
  RPCharacter? _playerCharacter;
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

    // Preselect story cards matching the provided world, if any
    if (widget.world != null && widget.world!.isNotEmpty) {
      final matches = cards.where((c) => c.world == widget.world).toList();
      if (matches.isNotEmpty) setState(() => _selectedCards = matches);
    }

    // If there's exactly one character in the library, preselect it as a participant and player
    if (chars.length == 1) setState(() {
      _selectedCharacters = chars;
      _playerCharacter = chars.first;
    });
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
    final world = widget.world ?? (_selectedCards.isNotEmpty ? _selectedCards.first.world : '') ?? '';
    final rules = _rulesController.text.trim();
    final s = RPSession(
      name: name,
      world: world,
      playerCharacterId: _playerCharacter?.id,
      characterIds: _selectedCharacters.map((c) => c.id!).toList(),
      storyCardIds: _selectedCards.map((c) => c.id!).toList(),
      rules: rules,
      createdAt: now,
    );
    final id = await _repo.insertSession(s);
    // create session object with returned id and navigate straight to chat
    final created = RPSession(
      id: id,
      name: name,
      world: world,
      playerCharacterId: _playerCharacter?.id,
      characterIds: _selectedCharacters.map((c) => c.id!).toList(),
      storyCardIds: _selectedCards.map((c) => c.id!).toList(),
      rules: rules,
      createdAt: now,
    );
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
    final world = widget.world ?? (_selectedCards.isNotEmpty ? _selectedCards.first.world : '') ?? '';
    final charName = _playerCharacter?.name ?? (_selectedCharacters.isNotEmpty ? _selectedCharacters.first.name : 'None');
    final cardTitle = _selectedCards.isNotEmpty ? _selectedCards.map((c) => c.title).join(', ') : 'None';
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
                      // multi-select characters via dialog
                      final chosenIds = await showDialog<List<int>>(context: context, builder: (c) {
                        final selected = Set<int>.from(_selectedCharacters.map((e) => e.id!).whereType<int>());
                        return StatefulBuilder(builder: (ctx, setSt) {
                          return AlertDialog(
                            title: const Text('Select Characters'),
                            content: SizedBox(
                              width: double.maxFinite,
                              child: ListView(
                                shrinkWrap: true,
                                children: _allCharacters.map((ch) {
                                  final id = ch.id!;
                                  return CheckboxListTile(
                                    value: selected.contains(id),
                                    title: Text(ch.name),
                                    subtitle: Text(ch.description, maxLines: 1, overflow: TextOverflow.ellipsis),
                                    onChanged: (v) => setSt(() => v == true ? selected.add(id) : selected.remove(id)),
                                  );
                                }).toList(),
                              ),
                            ),
                            actions: [
                              TextButton(onPressed: () => Navigator.of(c).pop(null), child: const Text('Cancel')),
                              ElevatedButton(onPressed: () => Navigator.of(c).pop(selected.toList()), child: const Text('OK')),
                            ],
                          );
                        });
                      });
                      if (chosenIds != null) {
                        setState(() {
                          _selectedCharacters = _allCharacters.where((c) => c.id != null && chosenIds.contains(c.id)).toList();
                          if (_playerCharacter == null || !_selectedCharacters.any((c) => c.id == _playerCharacter?.id)) {
                            _playerCharacter = _selectedCharacters.isNotEmpty ? _selectedCharacters.first : null;
                          }
                        });
                      }
                    },
                    icon: const Icon(Icons.person),
                    label: const Text('Pick Characters'),
                  ),
                ),
                const SizedBox(width: 8),
                Expanded(
                  child: ElevatedButton.icon(
                    onPressed: () async {
                      final chosenIds = await showDialog<List<int>>(context: context, builder: (c) {
                        final filtered = widget.world != null && widget.world!.isNotEmpty ? _allCards.where((a) => a.world == widget.world).toList() : _allCards;
                        final selected = Set<int>.from(_selectedCards.map((e) => e.id!).whereType<int>());
                        return StatefulBuilder(builder: (ctx, setSt) {
                          return AlertDialog(
                            title: const Text('Select Story Cards'),
                            content: SizedBox(
                              width: double.maxFinite,
                              child: ListView(
                                shrinkWrap: true,
                                children: filtered.map((card) {
                                  final id = card.id!;
                                  return CheckboxListTile(
                                    value: selected.contains(id),
                                    title: Text(card.title),
                                    subtitle: Text(card.world.isNotEmpty ? card.world : 'No world'),
                                    onChanged: (v) => setSt(() => v == true ? selected.add(id) : selected.remove(id)),
                                  );
                                }).toList(),
                              ),
                            ),
                            actions: [
                              TextButton(onPressed: () => Navigator.of(c).pop(null), child: const Text('Cancel')),
                              ElevatedButton(onPressed: () => Navigator.of(c).pop(selected.toList()), child: const Text('OK')),
                            ],
                          );
                        });
                      });
                      if (chosenIds != null) {
                        setState(() => _selectedCards = _allCards.where((c) => c.id != null && chosenIds.contains(c.id)).toList());
                      }
                    },
                    icon: const Icon(Icons.menu_book),
                    label: const Text('Pick Story Cards'),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),
            if (_selectedCharacters.isNotEmpty)
              Column(
                children: [
                  ListTile(title: Text('Characters (${_selectedCharacters.length})')),
                  for (final ch in _selectedCharacters)
                    ListTile(title: Text(ch.name), subtitle: Text(ch.description, maxLines: 1, overflow: TextOverflow.ellipsis)),
                  if (_selectedCharacters.length > 1)
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 8.0),
                      child: DropdownButton<RPCharacter>(
                        isExpanded: true,
                        value: _playerCharacter ?? _selectedCharacters.first,
                        items: _selectedCharacters.map((c) => DropdownMenuItem(value: c, child: Text('Player: ${c.name}'))).toList(),
                        onChanged: (v) => setState(() => _playerCharacter = v),
                      ),
                    ),
                ],
              ),
            if (_selectedCards.isNotEmpty)
              Column(
                children: [
                  ListTile(title: Text('Story Cards (${_selectedCards.length})')),
                  for (final card in _selectedCards)
                    ListTile(title: Text(card.title), subtitle: Text(card.world.isNotEmpty ? card.world : 'No world')),
                ],
              ),
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
                      _selectedCharacters.clear();
                      _selectedCards.clear();
                      _playerCharacter = null;
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