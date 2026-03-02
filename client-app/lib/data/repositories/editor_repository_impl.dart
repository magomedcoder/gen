import 'package:gen/core/failures.dart';
import 'package:gen/core/log/logs.dart';
import 'package:gen/data/data_sources/remote/editor_remote_datasource.dart';
import 'package:gen/domain/repositories/editor_repository.dart';

class EditorRepositoryImpl implements EditorRepository {
  final IEditorRemoteDataSource dataSource;

  EditorRepositoryImpl(this.dataSource);

  @override
  Future<String> transform({
    required String text,
    String? model,
  }) async {
    try {
      return await dataSource.transform(
        text: text,
        model: model,
      );
    } catch (e) {
      if (e is Failure) rethrow;
      Logs().e('EditorRepository: неожиданная ошибка transform',
          exception: e);
      throw ApiFailure('Ошибка обработки текста');
    }
  }
}
