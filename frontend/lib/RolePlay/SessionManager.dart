import 'package:flutter/material.dart';

class SessionManager extends StatefulWidget {
  const SessionManager({super.key});

  @override
  State<SessionManager> createState() => _SessionManagerState();
}

class _SessionManagerState extends State<SessionManager> {
  final List<String> _sessions = [];

  void _createSession() {
    // TODO: create and persist RP session
    setState(() => _sessions.add('Session ${_sessions.length + 1}'));
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
                itemBuilder: (context, i) => ListTile(title: Text(_sessions[i]), trailing: const Icon(Icons.play_arrow)),
              ),
            ),
          ],
        ),
      ),
    );
  }
}