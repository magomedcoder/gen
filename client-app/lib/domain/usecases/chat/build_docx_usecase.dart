import 'dart:typed_data';

import 'package:gen/domain/repositories/chat_repository.dart';

class BuildDocxUseCase {
  final ChatRepository _repository;

  BuildDocxUseCase(this._repository);

  Future<Uint8List> call({required String specJson}) {
    return _repository.buildDocx(specJson: specJson);
  }
}
