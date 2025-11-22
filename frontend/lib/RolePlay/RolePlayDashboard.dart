import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'CharacterEditor.dart';
import 'StoryCardEditor.dart';
import 'SessionManager.dart';
import 'RolePlaySettings.dart';
import 'CharacterList.dart';
import 'StoryCardList.dart';
import 'roleplay_models.dart';
import 'roleplay_repository.dart';
import 'package:shared_preferences/shared_preferences.dart';

const Color rpBackground = Color(0xFFFFF0F5);
const Color rpAccent = Color(0xFFFF69B4);
const Color rpText = Color(0xFF8B4789);

class RolePlayDashboard extends StatefulWidget {
  const RolePlayDashboard({super.key});

  @override
  State<RolePlayDashboard> createState() => _RolePlayDashboardState();
}

class _RolePlayDashboardState extends State<RolePlayDashboard> {
  final RolePlayRepository repo = RolePlayRepository();

  void _refresh() => setState(() {});

  String selectedWorld = ''; // empty = all
  Future<List<String>>? _worldsFuture;

  @override
  void initState() {
    super.initState();
    _worldsFuture = repo.getWorlds();
    _loadPrefs();
  }

  Future<void> _loadPrefs() async {
    final prefs = await SharedPreferences.getInstance();
    setState(() {
      selectedWorld = prefs.getString('rp_selected_world') ?? '';
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: rpBackground,
      appBar: AppBar(
        backgroundColor: rpAccent,
        title: Text('Roleplay Dashboard', style: GoogleFonts.pacifico()),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Roleplay (RP) Dashboard', style: GoogleFonts.pacifico(fontSize: 26, color: rpText)),
            const SizedBox(height: 8),
            Text('Create characters, manage story cards, and start RP sessions.', style: GoogleFonts.quicksand(color: rpText)),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: FutureBuilder<List<String>>(
                    future: _worldsFuture,
                    builder: (context, snap) {
                      final worlds = snap.data ?? [];
                      final items = ['All'] + worlds;
                      return DropdownButtonFormField<String>(
                        value: selectedWorld.isEmpty ? 'All' : selectedWorld,
                        items: items.map((w) => DropdownMenuItem(value: w, child: Text(w))).toList(),
                        onChanged: (v) async {
                          setState(() {
                            if (v == null || v == 'All') selectedWorld = '';
                            else selectedWorld = v;
                          });
                          final prefs = await SharedPreferences.getInstance();
                          await prefs.setString('rp_selected_world', selectedWorld);
                        },
                        decoration: const InputDecoration(labelText: 'Active World'),
                      );
                    },
                  ),
                ),
                const SizedBox(width: 12),
              ],
            ),
            const SizedBox(height: 12),
            Wrap(
              spacing: 12,
              runSpacing: 12,
              children: [
                ElevatedButton.icon(
                  style: ElevatedButton.styleFrom(backgroundColor: rpAccent),
                  onPressed: () async {
                    final result = await Navigator.push(context, MaterialPageRoute(builder: (_) => const CharacterEditor()));
                    if (result != null) {
                      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Character saved')));
                      _refresh();
                    }
                  },
                  icon: const Icon(Icons.person_add),
                  label: const Text('New Character'),
                ),
                ElevatedButton.icon(
                  style: ElevatedButton.styleFrom(backgroundColor: Color(0xFFFFB3D9)),
                  onPressed: () async {
                    final res = await Navigator.push(context, MaterialPageRoute(builder: (_) => CharacterList()));
                    if (res != null) {
                      // user picked a character from list
                      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Character selected')));
                    }
                    _refresh();
                  },
                  icon: const Icon(Icons.list),
                  label: const Text('Characters'),
                ),
                ElevatedButton.icon(
                  style: ElevatedButton.styleFrom(backgroundColor: Color(0xFFE6B3FF)),
                  onPressed: () async {
                    final result = await Navigator.push(context, MaterialPageRoute(builder: (_) => const StoryCardEditor()));
                    if (result != null) {
                      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Story card saved')));
                      _refresh();
                    }
                  },
                  icon: const Icon(Icons.menu_book),
                  label: const Text('Story Card'),
                ),
                ElevatedButton.icon(
                  style: ElevatedButton.styleFrom(backgroundColor: Color(0xFFE6B3FF).withOpacity(0.9)),
                  onPressed: () async {
                    final res = await Navigator.push(context, MaterialPageRoute(builder: (_) => StoryCardList(world: selectedWorld.isEmpty ? null : selectedWorld)));
                    if (res != null) {
                      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Story card selected')));
                    }
                    _refresh();
                  },
                  icon: const Icon(Icons.menu_book_outlined),
                  label: const Text('Story Cards'),
                ),
                ElevatedButton.icon(
                  style: ElevatedButton.styleFrom(backgroundColor: Color(0xFFFFD6E8)),
                  onPressed: () => Navigator.push(context, MaterialPageRoute(builder: (_) => SessionManager(world: selectedWorld))),
                  icon: const Icon(Icons.play_circle_outline),
                  label: const Text('Start Session'),
                ),
                ElevatedButton.icon(
                  style: ElevatedButton.styleFrom(backgroundColor: Color(0xFFFFF8DC)),
                  onPressed: () async {
                    final result = await Navigator.push(context, MaterialPageRoute(builder: (_) => const RolePlaySettings()));
                    if (result != null) {
                      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('RP settings saved')));
                    }
                  },
                  icon: const Icon(Icons.settings),
                  label: const Text('RP Settings'),
                ),
              ],
            ),
            const SizedBox(height: 24),
            FutureBuilder<List<dynamic>>(
              // Let type inference handle the Future.wait generic parameter
              future: Future.wait([
                repo.getCharacters(),
                repo.getStoryCards(),
                repo.getSessions(),
              ]),
              builder: (context, snapshot) {
                if (!snapshot.hasData) {
                  return const SizedBox(
                    height: 64,
                    child: Center(child: CircularProgressIndicator()),
                  );
                }
                final chars = (snapshot.data![0] as List).length;
                final cards = (snapshot.data![1] as List).length;
                final sessions = (snapshot.data![2] as List).length;
                return Card(
                  elevation: 2,
                  child: Padding(
                    padding: const EdgeInsets.all(12),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text('Library', style: GoogleFonts.quicksand(fontWeight: FontWeight.w700)),
                        const SizedBox(height: 8),
                        Text('Characters: $chars   Story Cards: $cards   Sessions: $sessions', style: GoogleFonts.quicksand(color: Colors.grey[700])),
                      ],
                    ),
                  ),
                );
              },
            ),
          ],
        ),
      ),
    );
  }
}