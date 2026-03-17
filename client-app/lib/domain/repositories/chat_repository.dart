import 'dart:async';

import 'package:gen/domain/entities/message.dart';
import 'package:gen/domain/entities/session.dart';

abstract interface class ChatRepository {
  Future<bool> checkConnection();

  Future<List<String>> getModels();

  Stream<String> sendMessage(
    int sessionId,
    List<Message> messages, {
    String? model,
  });

  Future<ChatSession> createSession(String title, {String? model});

  Future<ChatSession> getSession(int sessionId);

  Future<List<ChatSession>> listSessions(int page, int pageSize);

  Future<List<Message>> getSessionMessages(
    int sessionId,
    int page,
    int pageSize,
  );

  Future<void> deleteSession(int sessionId);

  Future<ChatSession> updateSessionTitle(int sessionId, String title);

  Future<ChatSession> updateSessionModel(int sessionId, String model);

  Future<String?> getSessionModel(int sessionId);

  Future<void> setSessionModel(int sessionId, String model);
}
