import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'roleplay_models.dart';
import 'roleplay_repository.dart';
import 'CharacterEditor.dart';

class CharacterList extends StatefulWidget {
  final String? world;
  const CharacterList({super.key, this.world});

  @override
  State<CharacterList> createState() => _CharacterListState();
}

class _CharacterListState extends State<CharacterList> {
  final RolePlayRepository _repo = RolePlayRepository();
  List<RPCharacter> _chars = [];

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    final list = await _repo.getCharacters();
    final filtered = widget.world == null || widget.world!.isEmpty ? list : list.where((c) => c.world == widget.world).toList();
    setState(() => _chars = filtered);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(widget.world?.isNotEmpty == true ? 'Characters: ${widget.world}' : 'Characters')),
      body: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          children: [
            ElevatedButton.icon(
              onPressed: () async {
                // navigate to editor to create new
                final res = await Navigator.push(context, MaterialPageRoute(builder: (_) => const CharacterEditor()));
                if (res != null) _load();
              },
              icon: const Icon(Icons.add),
              label: const Text('New Character'),
            ),
            const SizedBox(height: 12),
            Expanded(
              child: ListView.builder(
                itemCount: _chars.length,
                itemBuilder: (c, i) {
                  final ch = _chars[i];
                  return ListTile(
                    title: Text(ch.name),
                    subtitle: Text(ch.description, maxLines: 2, overflow: TextOverflow.ellipsis),
                    onTap: () => Navigator.pop(context, ch),
                  );
                },
              ),
            ),
          ],
        ),
      ),
    );
  }
}


