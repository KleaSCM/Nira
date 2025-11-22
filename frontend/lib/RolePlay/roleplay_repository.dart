import 'dart:async';
import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:path/path.dart';
import 'package:sqflite/sqflite.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'roleplay_models.dart';

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
    // bump DB version to 2 to add 'world' column for story_cards
    return await openDatabase(path, version: 2, onCreate: (db, ver) async {
      await db.execute('''
        CREATE TABLE characters (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          name TEXT NOT NULL,
          description TEXT NOT NULL
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
          metadata TEXT,
          created_at INTEGER NOT NULL
        )
      ''');
    }, onUpgrade: (db, oldVer, newVer) async {
      if (oldVer < 2) {
        // add world column to story_cards if upgrading from v1
        try {
          await db.execute("ALTER TABLE story_cards ADD COLUMN world TEXT DEFAULT ''");
        } catch (_) {
          // ignore if column already exists
        }
      }
    });
  }

  // Web fallback using SharedPreferences because sqflite is not supported on web.
  Future<SharedPreferences> _prefs() async => await SharedPreferences.getInstance();

  // Characters
  Future<int> insertCharacter(RPCharacter c) async {
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
  Future<List<RPCharacter>> getCharacters() async {
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

  // Story cards
  Future<int> insertStoryCard(RPStoryCard s) async {
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
  Future<List<RPStoryCard>> getStoryCards() async {
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
    // Return distinct worlds from story cards
    if (kIsWeb) {
      final prefs = await _prefs();
      final raw = prefs.getString('rp_cards') ?? '[]';
      final rows = (jsonDecode(raw) as List).cast<Map<String, dynamic>>();
      final worlds = <String>{};
      for (final r in rows) {
        final w = (r['world'] as String?) ?? '';
        if (w.isNotEmpty) worlds.add(w);
      }
      return worlds.toList();
    } else {
      final d = await db;
      final rows = await d.rawQuery("SELECT DISTINCT world FROM story_cards WHERE world IS NOT NULL AND world != ''");
      return rows.map((r) => (r['world'] as String)).toList();
    }
  }
}