import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'roleplay_models.dart';
import 'roleplay_repository.dart';

class CharacterEditor extends StatefulWidget {
  final String? world;
  const CharacterEditor({super.key, this.world});

  @override
  State<CharacterEditor> createState() => _CharacterEditorState();
}

class _CharacterEditorState extends State<CharacterEditor> {
  final TextEditingController _name = TextEditingController();
  final TextEditingController _desc = TextEditingController();
  final TextEditingController _world = TextEditingController();

  @override
  void initState() {
    super.initState();
    _world.text = widget.world ?? '';
  }

  @override
  void dispose() {
    _name.dispose();
    _desc.dispose();
    super.dispose();
  }

  void _save() async {
    final name = _name.text.trim();
    if (name.isEmpty) return;
    final desc = _desc.text.trim();
    final world = _world.text.trim();
    final c = RPCharacter(name: name, description: desc, world: world);
    final id = await RolePlayRepository().insertCharacter(c);
    final created = RPCharacter(id: id, name: name, description: desc, world: world);
    // return created object so callers can refresh and show feedback
    Navigator.pop(context, created);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Character Editor'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            TextField(controller: _name, decoration: const InputDecoration(labelText: 'Name')),
            const SizedBox(height: 12),
            TextField(controller: _desc, decoration: const InputDecoration(labelText: 'Description'), maxLines: 6),
            const SizedBox(height: 12),
            TextField(controller: _world, decoration: const InputDecoration(labelText: 'World (optional)'), maxLines: 1),
            const SizedBox(height: 16),
            ElevatedButton(onPressed: _save, child: const Text('Save')),
            const SizedBox(height: 16),
            ElevatedButton(onPressed: _save, child: const Text('Save')),
          ],
        ),
      ),
    );
  }
}