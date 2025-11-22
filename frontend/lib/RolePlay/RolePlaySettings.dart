import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

const Color settingsBackground = Color(0xFFFFF7FB);
const Color settingsAccent = Color(0xFFFF69B4);
const Color settingsText = Color(0xFF8B4789);

class RolePlaySettings extends StatefulWidget {
  const RolePlaySettings({super.key});

  @override
  State<RolePlaySettings> createState() => _RolePlaySettingsState();
}

class _RolePlaySettingsState extends State<RolePlaySettings> {
  bool autosave = true;
  bool showHints = true;
  String selectedModel = 'HammerAI/mythomax-l2';
  final TextEditingController _systemPrompt = TextEditingController(text: 'You are a helpful RP assistant.');

  final List<String> availableModels = [
    'HammerAI/mythomax-l2',
    'local-embedder-small',
    'experimental-rp-model',
  ];

  @override
  void dispose() {
    _systemPrompt.dispose();
    super.dispose();
  }

  void _save() {
    final settings = {
      'autosave': autosave,
      'showHints': showHints,
      'model': selectedModel,
      'systemPrompt': _systemPrompt.text.trim(),
    };
    Navigator.pop(context, settings);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: settingsBackground,
      appBar: AppBar(
        backgroundColor: settingsAccent,
        title: Text('RP Settings', style: GoogleFonts.pacifico()),
        actions: [
          TextButton(
            onPressed: _save,
            child: Text('Save', style: GoogleFonts.quicksand(color: Colors.white)),
          ),
        ],
      ),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('General', style: GoogleFonts.quicksand(fontSize: 18, fontWeight: FontWeight.w700, color: settingsText)),
            const SizedBox(height: 8),
            SwitchListTile(
              title: Text('Autosave sessions', style: GoogleFonts.quicksand(color: settingsText)),
              value: autosave,
              onChanged: (v) => setState(() => autosave = v),
            ),
            SwitchListTile(
              title: Text('Show context hints', style: GoogleFonts.quicksand(color: settingsText)),
              value: showHints,
              onChanged: (v) => setState(() => showHints = v),
            ),
            const SizedBox(height: 12),
            Text('Model & Prompt', style: GoogleFonts.quicksand(fontSize: 18, fontWeight: FontWeight.w700, color: settingsText)),
            const SizedBox(height: 8),
            DropdownButtonFormField<String>(
              value: selectedModel,
              items: availableModels.map((m) => DropdownMenuItem(value: m, child: Text(m))).toList(),
              onChanged: (v) => setState(() => selectedModel = v ?? selectedModel),
              decoration: const InputDecoration(border: OutlineInputBorder(), labelText: 'Default RP Model'),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _systemPrompt,
              maxLines: 6,
              decoration: const InputDecoration(
                border: OutlineInputBorder(),
                labelText: 'Default System Prompt',
                alignLabelWithHint: true,
              ),
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                ElevatedButton.icon(
                  onPressed: _save,
                  icon: const Icon(Icons.save),
                  label: Text('Save', style: GoogleFonts.quicksand()),
                  style: ElevatedButton.styleFrom(backgroundColor: settingsAccent),
                ),
                const SizedBox(width: 12),
                OutlinedButton(
                  onPressed: () => Navigator.pop(context),
                  child: Text('Cancel', style: GoogleFonts.quicksand(color: settingsText)),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}