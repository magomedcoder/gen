import 'package:gen/domain/repositories/chat_repository.dart';

class ApplyMarkdownPatchUseCase {
  final ChatRepository _repository;

  ApplyMarkdownPatchUseCase(this._repository);

  Future<String> call({
    required String baseText,
    required String patchJson,
  }) {
    return _repository.applyMarkdownPatch(
      baseText: baseText,
      patchJson: patchJson,
    );
  }
}
