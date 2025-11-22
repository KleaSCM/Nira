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
  Future<List<RPSession>>? _sessionsFuture;

  @override
  void initState() {
    super.initState();
    _worldsFuture = repo.getWorlds();
    _loadPrefs();
    _loadSessionsFuture();
  }

  Future<void> _loadPrefs() async {
    final prefs = await SharedPreferences.getInstance();
    setState(() {
      selectedWorld = prefs.getString('rp_selected_world') ?? '';
    });
  }

  void _loadSessionsFuture() {
    _sessionsFuture = repo.getSessions().then((list) {
      if (selectedWorld.isEmpty) return list;
      return list.where((s) => (s.world ?? '') == selectedWorld).toList();
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
                          // refresh session list for the selected world
                          setState(() => _loadSessionsFuture());
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
                  style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFFFFEFD5)),
                  onPressed: () async {
                    final nameCtrl = TextEditingController();
                    final descCtrl = TextEditingController();
                    final res = await showDialog<bool>(
                      context: context,
                      builder: (c) => AlertDialog(
                        title: const Text('Create World'),
                        content: Column(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            TextField(controller: nameCtrl, decoration: const InputDecoration(labelText: 'World name')),
                            TextField(controller: descCtrl, decoration: const InputDecoration(labelText: 'Description (optional)')),
                          ],
                        ),
                        actions: [
                          TextButton(onPressed: () => Navigator.of(c).pop(false), child: const Text('Cancel')),
                          ElevatedButton(onPressed: () => Navigator.of(c).pop(true), child: const Text('Create')),
                        ],
                      ),
                    );
                    if (res == true && nameCtrl.text.trim().isNotEmpty) {
                      final id = await repo.insertWorld(RPWorld(name: nameCtrl.text.trim(), description: descCtrl.text.trim()));
                      if (id == -1) {
                        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('World already exists')));
                      } else {
                        // refresh worlds and select the new one
                        setState(() {
                          _worldsFuture = repo.getWorlds();
                          selectedWorld = nameCtrl.text.trim();
                          _loadSessionsFuture();
                        });
                        final prefs = await SharedPreferences.getInstance();
                        await prefs.setString('rp_selected_world', selectedWorld);
                        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('World created')));
                      }
                    }
                  },
                  icon: const Icon(Icons.public),
                  label: const Text('New World'),
                ),
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
            const SizedBox(height: 16),
            Text('World Sessions', style: GoogleFonts.quicksand(fontWeight: FontWeight.w700)),
            const SizedBox(height: 8),
            FutureBuilder<List<RPSession>>(
              future: _sessionsFuture,
              builder: (context, snap) {
                if (!snap.hasData) return const SizedBox(height: 80, child: Center(child: CircularProgressIndicator()));
                final sessions = snap.data!;
                if (sessions.isEmpty) return Text(selectedWorld.isEmpty ? 'No sessions yet.' : 'No sessions for "$selectedWorld"');
                return Card(
                  elevation: 1,
                  child: Column(
                    children: sessions.map((s) {
                      return ListTile(
                        title: Text(s.name),
                        subtitle: Text(s.world ?? ''),
                        trailing: const Icon(Icons.play_arrow),
                        onTap: () => Navigator.push(context, MaterialPageRoute(builder: (_) => SessionManager(world: s.world))),
                      );
                    }).toList(),
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