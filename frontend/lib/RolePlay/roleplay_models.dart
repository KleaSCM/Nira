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
  int? characterId;
  int? storyCardId;
  String rules;
  int createdAt;

  RPSession({
    this.id,
    required this.name,
    this.world = '',
    this.characterId,
    this.storyCardId,
    this.rules = '',
    required this.createdAt,
  });

  Map<String, dynamic> toMap() => {
        if (id != null) 'id': id,
        'name': name,
        'world': world,
        'character_id': characterId,
        'story_card_id': storyCardId,
        'rules': rules,
        'created_at': createdAt,
      };

  factory RPSession.fromMap(Map<String, dynamic> m) => RPSession(
        id: m['id'] as int?,
        name: m['name'] as String,
        world: (m['world'] as String?) ?? '',
        characterId: m['character_id'] as int?,
        storyCardId: m['story_card_id'] as int?,
        rules: (m['rules'] as String?) ?? '',
        createdAt: m['created_at'] as int,
      );
}