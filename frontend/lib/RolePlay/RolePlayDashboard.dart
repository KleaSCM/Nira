import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'CharacterEditor.dart';
import 'StoryCardEditor.dart';
import 'SessionManager.dart';
import 'RolePlaySettings.dart';

const Color rpBackground = Color(0xFFFFF0F5);
const Color rpAccent = Color(0xFFFF69B4);
const Color rpText = Color(0xFF8B4789);

class RolePlayDashboard extends StatelessWidget {
  const RolePlayDashboard({super.key});

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
            Wrap(
              spacing: 12,
              runSpacing: 12,
              children: [
                ElevatedButton.icon(
                  style: ElevatedButton.styleFrom(backgroundColor: rpAccent),
                  onPressed: () => Navigator.push(context, MaterialPageRoute(builder: (_) => const CharacterEditor())),
                  icon: const Icon(Icons.person_add),
                  label: const Text('New Character'),
                ),
                ElevatedButton.icon(
                  style: ElevatedButton.styleFrom(backgroundColor: Color(0xFFE6B3FF)),
                  onPressed: () => Navigator.push(context, MaterialPageRoute(builder: (_) => const StoryCardEditor())),
                  icon: const Icon(Icons.menu_book),
                  label: const Text('Story Card'),
                ),
                ElevatedButton.icon(
                  style: ElevatedButton.styleFrom(backgroundColor: Color(0xFFFFD6E8)),
                  onPressed: () => Navigator.push(context, MaterialPageRoute(builder: (_) => const SessionManager())),
                  icon: const Icon(Icons.play_circle_outline),
                  label: const Text('Start Session'),
                ),
                ElevatedButton.icon(
                  style: ElevatedButton.styleFrom(backgroundColor: Color(0xFFFFF8DC)),
                  onPressed: () async {
                    // Open RP Settings screen
                    final result = await Navigator.push(context, MaterialPageRoute(builder: (_) => const RolePlaySettings()));
                    // result can be a Map<String, dynamic> with saved settings (persist as needed)
                    if (result != null) {
                      // TODO: persist or apply the returned settings
                    }
                  },
                  icon: const Icon(Icons.settings),
                  label: const Text('RP Settings'),
                ),
              ],
            ),
            const SizedBox(height: 24),
            Card(
              elevation: 2,
              child: Padding(
                padding: const EdgeInsets.all(12),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text('Recent RP Sessions', style: GoogleFonts.quicksand(fontWeight: FontWeight.w700)),
                    const SizedBox(height: 8),
                    Text('No sessions yet. Create one using "Start Session".', style: GoogleFonts.quicksand(color: Colors.grey[700])),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}