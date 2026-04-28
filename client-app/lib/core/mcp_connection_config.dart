import 'dart:convert';

import 'package:gen/domain/entities/mcp_server_entity.dart';

class McpConnectionConfig {
  const McpConnectionConfig({
    required this.transport,
    required this.command,
    required this.args,
    required this.env,
    required this.url,
    required this.headers,
    required this.timeoutSeconds,
  });

  final String transport;
  final String command;
  final List<String> args;
  final Map<String, String> env;
  final String url;
  final Map<String, String> headers;
  final int timeoutSeconds;

  static const Set<String> _transports = {'sse', 'streamable'};

  static String defaultJsonPretty() {
    return const JsonEncoder.withIndent('  ').convert(_defaultMap());
  }

  static Map<String, Object> _defaultMap() => {
    'transport': 'sse',
    'command': '',
    'args': <String>[],
    'env': <String, String>{},
    'url': '',
    'headers': <String, String>{},
    'timeoutSeconds': 120,
  };

  static String _pretty(Map<String, Object> m) => const JsonEncoder.withIndent('  ').convert(m);

  static String exampleJsonSse() {
    return _pretty({
      'transport': 'sse',
      'command': '',
      'args': <String>[],
      'env': <String, String>{},
      'url': 'https://сайт/mcp/sse',
      'headers': <String, String>{
        'Authorization': 'Bearer токен',
      },
      'timeoutSeconds': 120,
    });
  }

  static String exampleJsonStreamable() {
    return _pretty({
      'transport': 'streamable',
      'command': '',
      'args': <String>[],
      'env': <String, String>{},
      'url': 'https://сайт/mcp',
      'headers': <String, String>{},
      'timeoutSeconds': 180,
    });
  }

  static String exampleJsonFullRemote() {
    return _pretty({
      'transport': 'streamable',
      'command': '',
      'args': <String>[],
      'env': <String, String>{},
      'url': 'https://сайт/mcp/v1/stream?region=eu&trace=1',
      'headers': <String, String>{
        'Authorization': 'Bearer токен',
        'Content-Type': 'application/json',
        'Accept': 'application/json, text/event-stream',
        'X-Request-ID': 'example-correlation-id',
        'X-Client-Name': 'gen-mcp',
        'X-Api-Version': '',
      },
      'timeoutSeconds': 600,
    });
  }

  static const String documentation = ''
      'Поля:\n'
      ' transport - способ подключения:\n'
      '   - sse - сервер по HTTP, обычно с Server-Sent Events\n'
      '   - streamable - HTTP-транспорт с потоковой передачей\n'
      ' command - для HTTP не используется, оставьте пустую строку "".\n'
      ' args - для HTTP не используется, оставьте empty массив [].\n'
      ' env - для HTTP не используется, оставьте empty объект {}.\n'
      ' url - полный URL для sse/streamable (схема https:// или http://, путь к endpoint).\n'
      ' headers - заголовки HTTP для sse/streamable (например Authorization, X-Api-Key). Значения - строки.\n'
      ' timeoutSeconds - timeout вызовов инструментов, целое число от 1 до 600 (секунды).\n'
      '\n'
      'Режимы:\n'
      '  Удалённый сервер: transport (sse или streamable) и url, при необходимости headers.\n'
      '\n'
      'Кнопки sse и streamable - короткие примеры. Полный HTTP - развёрнутый образец с несколькими заголовками и timeoutом 600 с. '
      'Замените пути, хосты и секреты на свои.\n'
      '\n'
      'При редактировании сохранённого сервера значения секретов в env и headers на сервере могут отображаться как *** - не меняйте эти строки, если не хотите перезаписать секрет.';

  static String prettyFromEntity(McpServerEntity e) {
    return const JsonEncoder.withIndent('  ').convert({
      'transport': e.transport.trim().isEmpty ? 'sse' : e.transport.trim(),
      'command': e.command,
      'args': e.args,
      'env': e.env,
      'url': e.url,
      'headers': e.headers,
      'timeoutSeconds': e.timeoutSeconds > 0 ? e.timeoutSeconds : 120,
    });
  }

  static McpConnectionConfig parse(String raw) {
    final trimmed = raw.trim();
    if (trimmed.isEmpty) {
      throw const FormatException('Укажите JSON с настройками подключения.');
    }

    final dynamic decoded = jsonDecode(trimmed);
    if (decoded is! Map) {
      throw const FormatException('JSON должен быть объектом { ... }.');
    }

    final map = decoded.cast<String, dynamic>();

    final transport = (map['transport'] as String?)?.trim() ?? 'sse';
    if (!_transports.contains(transport)) {
      throw FormatException('transport должен быть одним из: ${_transports.join(", ")}.');
    }

    final command = (map['command'] as String?) ?? '';
    final args = _stringList(map['args']);
    final env = _stringMap(map['env']);
    final url = (map['url'] as String?) ?? '';
    final headers = _stringMap(map['headers']);

    final ts = map['timeoutSeconds'];
    int timeout = 120;
    if (ts is int) {
      timeout = ts;
    } else if (ts is num) {
      timeout = ts.round();
    } else if (ts != null) {
      throw const FormatException('timeoutSeconds должно быть числом.');
    }

    if (timeout <= 0 || timeout > 600) {
      throw const FormatException('timeoutSeconds: укажите число от 1 до 600.');
    }

    return McpConnectionConfig(
      transport: transport,
      command: command,
      args: args,
      env: env,
      url: url,
      headers: headers,
      timeoutSeconds: timeout,
    );
  }

  static List<String> _stringList(dynamic v) {
    if (v == null) {
      return [];
    }

    if (v is List) {
      return v.map((e) => '$e').toList();
    }

    throw const FormatException('args должен быть массивом строк.');
  }

  static Map<String, String> _stringMap(dynamic v) {
    if (v == null) {
      return {};
    }

    if (v is! Map) {
      throw const FormatException('env и headers должны быть объектами строк.');
    }

    final out = <String, String>{};
    for (final e in v.entries) {
      out['${e.key}'] = '${e.value}';
    }

    return out;
  }

  McpServerEntity toEntity({
    required int id,
    required String name,
    required bool enabled,
    int ownerUserId = 0,
  }) {
    return McpServerEntity(
      id: id,
      name: name,
      enabled: enabled,
      transport: transport,
      command: command,
      args: args,
      env: env,
      url: url,
      headers: headers,
      timeoutSeconds: timeoutSeconds,
      ownerUserId: ownerUserId,
    );
  }
}
