import 'dart:async';
import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:path/path.dart';
import 'package:sqflite/sqflite.dart';
import 'package:nira_frontend/WebSocketService.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'roleplay_models.dart';

/// RolePlayRepository
///
/// Frontend-only adapter with a thin client -> server boundary.
///
/// - Primary path: uses the backend RP tools over WebSocket (rp_* tools)
///   via WebSocketService.callToolJson(). This is not “backend code” living
///   in the frontend — it’s a remote call to the Go backend, similar to a
///   REST/RPC client.
/// - Fallback path: if the backend is unavailable (e.g., running the web
///   build without the Go server), we use local persistence to keep the UI
///   usable. On web, that’s SharedPreferences; on desktop/mobile, a small
///   sqflite database.
///
/// This keeps the Flutter layer as a GUI while centralizing RP logic and
/// persistence on the backend when it’s running. No RP business logic lives
/// here; it’s just transport + local fallback.

class RolePlayRepository {
  static final RolePlayRepository _instance = RolePlayRepository._internal();
  factory RolePlayRepository() => _instance;
  RolePlayRepository._internal();

  Database? _db;
  Future<Database> get db async {
    if (_db != null) return _db!;
    _db = await _initDb();
    return _db!;
  }

  Future<Database> _initDb() async {
    final databasesPath = await getDatabasesPath();
    final path = join(databasesPath, 'nira_roleplay.db');
    // bump DB version to 6 to add worlds table and session linkage
      return await openDatabase(path, version: 6, onCreate: (db, ver) async {
      await db.execute('''
        CREATE TABLE characters (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          name TEXT NOT NULL,
          description TEXT NOT NULL,
          world TEXT DEFAULT ''
        )
      ''');
      await db.execute('''
        CREATE TABLE story_cards (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          title TEXT NOT NULL,
          content TEXT NOT NULL,
          world TEXT DEFAULT ''
        )
      ''');
      await db.execute('''
        CREATE TABLE sessions (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          name TEXT NOT NULL,
          world TEXT,
          player_character_id INTEGER,
          character_ids TEXT,
          story_card_ids TEXT,
          rules TEXT,
          created_at INTEGER NOT NULL
        )
      ''');
      await db.execute('''
        CREATE TABLE worlds (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          name TEXT NOT NULL UNIQUE,
          description TEXT DEFAULT ''
        )
      ''');
    }, onUpgrade: (db, oldVer, newVer) async {
      if (oldVer < 2) {
        // add world column to story_cards if upgrading from v1
        try {
          await db.execute("ALTER TABLE story_cards ADD COLUMN world TEXT DEFAULT ''");
        } catch (_) {}
      }
      if (oldVer < 3) {
        // previous upgrade steps (ensure world column exists)
        try {
          await db.execute("ALTER TABLE sessions ADD COLUMN world TEXT");
        } catch (_) {}
      }
      if (oldVer < 4) {
        // add session linkage columns for multi-ids and player character
        try {
          await db.execute("ALTER TABLE sessions ADD COLUMN player_character_id INTEGER");
          await db.execute("ALTER TABLE sessions ADD COLUMN character_ids TEXT");
          await db.execute("ALTER TABLE sessions ADD COLUMN story_card_ids TEXT");
          await db.execute("ALTER TABLE sessions ADD COLUMN rules TEXT");
        } catch (_) {}
      }
      if (oldVer < 5) {
        try {
          await db.execute('''
            CREATE TABLE IF NOT EXISTS worlds (
              id INTEGER PRIMARY KEY AUTOINCREMENT,
              name TEXT NOT NULL UNIQUE,
              description TEXT DEFAULT ''
            )
          ''');
        } catch (_) {}
      }
        if (oldVer < 6) {
          // add world column to characters
          try {
            await db.execute("ALTER TABLE characters ADD COLUMN world TEXT DEFAULT ''");
          } catch (_) {}
        }
    });
  }

  // Web fallback using SharedPreferences because sqflite is not supported on web.
  Future<SharedPreferences> _prefs() async => await SharedPreferences.getInstance();

  // Characters
  Future<int> insertCharacter(RPCharacter c) async {
    // Prefer backend RP tool; fallback to local storage if it fails
    try {
      final ws = WebSocketService.instance;
      ws.connect('ws://localhost:8080/ws');
      final res = await ws.callToolJson('rp_character_save', {
        'name': c.name,
        // Map description -> summary; keep background empty for now
        'summary': c.description,
        'background': '',
        'traits': <String>[],
        'goals': <String>[],
        'tags': <String>[],
        'notes': '',
      });
      // Backend returns an object with string id; we ignore local int id
      return 0;
    } catch (_) {
      if (kIsWeb) {
        final prefs = await _prefs();
        final raw = prefs.getString('rp_chars') ?? '[]';
        final list = (jsonDecode(raw) as List).cast<Map<String, dynamic>>();
        final next = prefs.getInt('rp_char_next_id') ?? 1;
        final map = c.toMap();
        map['id'] = next;
        list.insert(0, map);
        await prefs.setString('rp_chars', jsonEncode(list));
        await prefs.setInt('rp_char_next_id', next + 1);
        return next;
      } else {
        final d = await db;
        return await d.insert('characters', c.toMap());
      }
    }
  }
  Future<List<RPCharacter>> getCharacters() async {
    // Prefer backend list tool; fallback to local storage
    try {
      final ws = WebSocketService.instance;
      ws.connect('ws://localhost:8080/ws');
      final res = await ws.callToolJson('rp_character_list', {
        'limit': 200,
        'offset': 0,
      });
      if (res is List) {
        return res.map<RPCharacter>((e) {
          final m = (e as Map).cast<String, dynamic>();
          return RPCharacter(
            id: null,
            name: (m['name'] ?? '') as String,
            description: (m['summary'] ?? m['notes'] ?? '') as String,
            world: '',
          );
        }).toList();
      }
      return <RPCharacter>[];
    } catch (_) {
      if (kIsWeb) {
        final prefs = await _prefs();
        final raw = prefs.getString('rp_chars') ?? '[]';
        final rows = (jsonDecode(raw) as List).cast<Map<String, dynamic>>();
        return rows.map((r) => RPCharacter.fromMap(r)).toList();
      } else {
        final d = await db;
        final rows = await d.query('characters', orderBy: 'id DESC');
        return rows.map((r) => RPCharacter.fromMap(r)).toList();
      }
    }
  }

  // Story cards
  Future<int> insertStoryCard(RPStoryCard s) async {
    // Prefer backend RP tool; fallback to local storage
    try {
      final ws = WebSocketService.instance;
      ws.connect('ws://localhost:8080/ws');
      final _ = await ws.callToolJson('rp_storycard_save', {
        'title': s.title,
        'kind': 'lore',
        'content': s.content,
        'tags': <String>[],
        'links': <String>[],
      });
      return 0;
    } catch (_) {
      if (kIsWeb) {
        final prefs = await _prefs();
        final raw = prefs.getString('rp_cards') ?? '[]';
        final list = (jsonDecode(raw) as List).cast<Map<String, dynamic>>();
        final next = prefs.getInt('rp_card_next_id') ?? 1;
        final map = s.toMap();
        map['id'] = next;
        list.insert(0, map);
        await prefs.setString('rp_cards', jsonEncode(list));
        await prefs.setInt('rp_card_next_id', next + 1);
        return next;
      } else {
        final d = await db;
        return await d.insert('story_cards', s.toMap());
      }
    }
  }
  Future<List<RPStoryCard>> getStoryCards() async {
    try {
      final ws = WebSocketService.instance;
      ws.connect('ws://localhost:8080/ws');
      final res = await ws.callToolJson('rp_storycard_list', {
        'limit': 200,
        'offset': 0,
      });
      if (res is List) {
        return res.map<RPStoryCard>((e) {
          final m = (e as Map).cast<String, dynamic>();
          return RPStoryCard(
            id: null,
            title: (m['title'] ?? '') as String,
            content: (m['content'] ?? '') as String,
            world: '',
          );
        }).toList();
      }
      return <RPStoryCard>[];
    } catch (_) {
      if (kIsWeb) {
        final prefs = await _prefs();
        final raw = prefs.getString('rp_cards') ?? '[]';
        final rows = (jsonDecode(raw) as List).cast<Map<String, dynamic>>();
        return rows.map((r) => RPStoryCard.fromMap(r)).toList();
      } else {
        final d = await db;
        final rows = await d.query('story_cards', orderBy: 'id DESC');
        return rows.map((r) => RPStoryCard.fromMap(r)).toList();
      }
    }
  }



  // Sessions
  Future<int> insertSession(RPSession s) async {
    if (kIsWeb) {
      final prefs = await _prefs();
      final raw = prefs.getString('rp_sessions') ?? '[]';
      final list = (jsonDecode(raw) as List).cast<Map<String, dynamic>>();
      final next = prefs.getInt('rp_session_next_id') ?? 1;
      final map = s.toMap();
      map['id'] = next;
      list.insert(0, map);
      await prefs.setString('rp_sessions', jsonEncode(list));
      await prefs.setInt('rp_session_next_id', next + 1);
      return next;
    } else {
      final d = await db;
      return await d.insert('sessions', s.toMap());
    }
  }
  Future<List<RPSession>> getSessions() async {
    if (kIsWeb) {
      final prefs = await _prefs();
      final raw = prefs.getString('rp_sessions') ?? '[]';
      final rows = (jsonDecode(raw) as List).cast<Map<String, dynamic>>();
      return rows.map((r) => RPSession.fromMap(r)).toList();
    } else {
      final d = await db;
      final rows = await d.query('sessions', orderBy: 'created_at DESC');
      return rows.map((r) => RPSession.fromMap(r)).toList();
    }
  }

  Future<List<String>> getWorlds() async {
    // Prefer worlds table; fall back to distinct world values on story_cards
    if (kIsWeb) {
      final prefs = await _prefs();
      final rawWorlds = prefs.getString('rp_worlds') ?? '[]';
      final list = (jsonDecode(rawWorlds) as List).cast<Map<String, dynamic>>();
      final names = list.map((m) => (m['name'] as String)).toSet();

      // also include any worlds referenced by story cards (web fallback)
      final raw = prefs.getString('rp_cards') ?? '[]';
      final rows = (jsonDecode(raw) as List).cast<Map<String, dynamic>>();
      for (final r in rows) {
        final w = (r['world'] as String?) ?? '';
        if (w.isNotEmpty) names.add(w);
      }
      return names.toList();
    } else {
      final d = await db;
      // check worlds table first
      try {
        final wRows = await d.query('worlds', orderBy: 'name ASC');
        if (wRows.isNotEmpty) return wRows.map((r) => (r['name'] as String)).toList();
      } catch (_) {}

      // fallback: distinct worlds from story_cards
      final rows = await d.rawQuery("SELECT DISTINCT world FROM story_cards WHERE world IS NOT NULL AND world != ''");
      return rows.map((r) => (r['world'] as String)).toList();
    }
  }

  // Worlds
  Future<int> insertWorld(RPWorld w) async {
    if (kIsWeb) {
      final prefs = await _prefs();
      final raw = prefs.getString('rp_worlds') ?? '[]';
      final list = (jsonDecode(raw) as List).cast<Map<String, dynamic>>();
      // prevent duplicates by name
      final exists = list.any((m) => (m['name'] as String).toLowerCase() == w.name.toLowerCase());
      if (exists) return -1;
      final next = prefs.getInt('rp_world_next_id') ?? 1;
      final map = w.toMap();
      map['id'] = next;
      list.insert(0, map);
      await prefs.setString('rp_worlds', jsonEncode(list));
      await prefs.setInt('rp_world_next_id', next + 1);
      return next;
    } else {
      final d = await db;
      try {
        return await d.insert('worlds', w.toMap());
      } catch (e) {
        // unique constraint failed
        return -1;
      }
    }
  }
}