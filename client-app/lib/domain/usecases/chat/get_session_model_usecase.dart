import 'package:gen/domain/repositories/chat_repository.dart';

class GetSessionModelUseCase {
  final ChatRepository repository;

  GetSessionModelUseCase(this.repository);

  Future<String?> call(int sessionId) {
    return repository.getSessionModel(sessionId);
  }
}
