import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'roleplay_models.dart';
import 'roleplay_repository.dart';
import 'StoryCardEditor.dart';

class StoryCardList extends StatefulWidget {
  final String? world;
  const StoryCardList({super.key, this.world});

  @override
  State<StoryCardList> createState() => _StoryCardListState();
}

class _StoryCardListState extends State<StoryCardList> {
  final RolePlayRepository _repo = RolePlayRepository();
  List<RPStoryCard> _cards = [];

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    final all = await _repo.getStoryCards();
    final filtered = widget.world == null || widget.world!.isEmpty
        ? all
        : all.where((c) => c.world == widget.world).toList();
    setState(() => _cards = filtered);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(widget.world?.isNotEmpty == true ? 'Cards: ${widget.world}' : 'Story Cards')),
      body: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          children: [
            ElevatedButton.icon(
              onPressed: () async {
                final res = await Navigator.push(context, MaterialPageRoute(builder: (_) => const StoryCardEditor()));
                if (res != null) _load();
              },
              icon: const Icon(Icons.add),
              label: const Text('New Story Card'),
            ),
            const SizedBox(height: 12),
            Expanded(
              child: ListView.builder(
                itemCount: _cards.length,
                itemBuilder: (c, i) {
                  final card = _cards[i];
                  return ListTile(
                    title: Text(card.title),
                    subtitle: Text(card.world.isNotEmpty ? card.world : 'No world'),
                    onTap: () => Navigator.pop(context, card),
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
