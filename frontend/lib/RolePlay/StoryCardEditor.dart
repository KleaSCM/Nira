import 'package:flutter/material.dart';

class StoryCardEditor extends StatefulWidget {
  const StoryCardEditor({super.key});

  @override
  State<StoryCardEditor> createState() => _StoryCardEditorState();
}

class _StoryCardEditorState extends State<StoryCardEditor> {
  final TextEditingController _title = TextEditingController();
  final TextEditingController _content = TextEditingController();

  @override
  void dispose() {
    _title.dispose();
    _content.dispose();
    super.dispose();
  }

  void _save() {
    // TODO: persist story card
    Navigator.pop(context);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Story Card'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            TextField(controller: _title, decoration: const InputDecoration(labelText: 'Title')),
            const SizedBox(height: 12),
            TextField(controller: _content, decoration: const InputDecoration(labelText: 'Content'), maxLines: 8),
            const SizedBox(height: 16),
            ElevatedButton(onPressed: _save, child: const Text('Save')),
          ],
        ),
      ),
    );
  }
}