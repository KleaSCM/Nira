import 'package:flutter/material.dart';
import 'roleplay_models.dart';
import 'roleplay_repository.dart';

class SessionManager extends StatefulWidget {
  const SessionManager({super.key});

  @override
  State<SessionManager> createState() => _SessionManagerState();
}

class _SessionManagerState extends State<SessionManager> {
  final List<RPSession> _sessions = [];
  final RolePlayRepository _repo = RolePlayRepository();

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    final s = await _repo.getSessions();
    setState(() => _sessions
      ..clear()
      ..addAll(s));
  }

  void _createSession() async {
    final now = DateTime.now().millisecondsSinceEpoch;
    final name = 'Session ${_sessions.length + 1}';
    final s = RPSession(name: name, metadata: '', createdAt: now);
    await _repo.insertSession(s);
    await _load();
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
            ElevatedButton.icon(onPressed: _createSession, icon: const Icon(Icons.add), label: const Text('New Session')),
            const SizedBox(height: 12),
            Expanded(
              child: ListView.builder(
                itemCount: _sessions.length,
                itemBuilder: (context, i) => ListTile(title: Text(_sessions[i].name), trailing: const Icon(Icons.play_arrow)),
              ),
            ),
          ],
        ),
      ),
    );
  }
}