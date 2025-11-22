import 'package:flutter/material.dart';
import 'roleplay_models.dart';
import 'roleplay_repository.dart';
import 'RPChatScreen.dart';

class SessionManager extends StatefulWidget {
  final String? world;
  const SessionManager({super.key, this.world});

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
    final metadata = widget.world != null && widget.world!.isNotEmpty ? '{"world":"${widget.world}"}' : '';
    final s = RPSession(name: name, metadata: metadata, createdAt: now);
    final id = await _repo.insertSession(s);
    // create session object with returned id and navigate straight to chat
    final created = RPSession(id: id, name: name, metadata: metadata, createdAt: now);
    await _load();
    // open RP chat for the new session
    if (!mounted) return;
    final res = await Navigator.push(context, MaterialPageRoute(builder: (_) => RPChatScreen(session: created)));
    // show a small confirmation after returning from chat (or immediately)
    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Session created and opened')));
    }
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
            ElevatedButton.icon(onPressed: _createSession, icon: const Icon(Icons.add), label: const Text('New Session')),
            const SizedBox(height: 12),
            Expanded(
              child: ListView.builder(
                itemCount: _sessions.length,
                itemBuilder: (context, i) => ListTile(
                  title: Text(_sessions[i].name),
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