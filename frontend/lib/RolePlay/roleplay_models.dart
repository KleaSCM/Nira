import 'dart:convert';

class RPCharacter {
  int? id;
  String name;
  String description;

  RPCharacter({this.id, required this.name, required this.description});

  Map<String, dynamic> toMap() => {
        if (id != null) 'id': id,
        'name': name,
        'description': description,
      };

  factory RPCharacter.fromMap(Map<String, dynamic> m) => RPCharacter(
        id: m['id'] as int?,
        name: m['name'] as String,
        description: m['description'] as String,
      );
}

class RPStoryCard {
  int? id;
  String title;
  String content;
  String world;

  RPStoryCard({this.id, required this.title, required this.content, this.world = ''});

  Map<String, dynamic> toMap() => {
        if (id != null) 'id': id,
        'title': title,
        'content': content,
        'world': world,
      };

  factory RPStoryCard.fromMap(Map<String, dynamic> m) => RPStoryCard(
        id: m['id'] as int?,
        title: m['title'] as String,
        content: m['content'] as String,
        world: (m['world'] as String?) ?? '',
      );
}

class RPSession {
  int? id;
  String name;
  String world;
  int? playerCharacterId;
  List<int> characterIds;
  List<int> storyCardIds;
  String rules;
  int createdAt;

  RPSession({
    this.id,
    required this.name,
    this.world = '',
    this.playerCharacterId,
    List<int>? characterIds,
    List<int>? storyCardIds,
    this.rules = '',
    required this.createdAt,
  })  : characterIds = characterIds ?? [],
        storyCardIds = storyCardIds ?? [];

  Map<String, dynamic> toMap() => {
        if (id != null) 'id': id,
        'name': name,
        'world': world,
        'player_character_id': playerCharacterId,
        'character_ids': jsonEncode(characterIds),
        'story_card_ids': jsonEncode(storyCardIds),
        'rules': rules,
        'created_at': createdAt,
      };

  factory RPSession.fromMap(Map<String, dynamic> m) => RPSession(
        id: m['id'] as int?,
        name: m['name'] as String,
        world: (m['world'] as String?) ?? '',
        playerCharacterId: m['player_character_id'] as int?,
        characterIds: m['character_ids'] != null && (m['character_ids'] is String)
            ? (jsonDecode(m['character_ids'] as String) as List).cast<int>()
            : <int>[],
        storyCardIds: m['story_card_ids'] != null && (m['story_card_ids'] is String)
            ? (jsonDecode(m['story_card_ids'] as String) as List).cast<int>()
            : <int>[],
        rules: (m['rules'] as String?) ?? '',
        createdAt: m['created_at'] as int,
      );
}

class RPWorld {
  int? id;
  String name;
  String description;

  RPWorld({this.id, required this.name, this.description = ''});

  Map<String, dynamic> toMap() => {
        if (id != null) 'id': id,
        'name': name,
        'description': description,
      };

  factory RPWorld.fromMap(Map<String, dynamic> m) => RPWorld(
        id: m['id'] as int?,
        name: m['name'] as String,
        description: (m['description'] as String?) ?? '',
      );
}