import 'package:gen/domain/entities/spreadsheet_apply_result.dart';
import 'package:gen/domain/repositories/chat_repository.dart';

class ApplySpreadsheetUseCase {
  final ChatRepository _repository;

  ApplySpreadsheetUseCase(this._repository);

  Future<SpreadsheetApplyResult> call({
    List<int>? workbookXlsx,
    required String operationsJson,
    String previewSheet = '',
    String previewRange = '',
  }) {
    return _repository.applySpreadsheet(
      workbookXlsx: workbookXlsx,
      operationsJson: operationsJson,
      previewSheet: previewSheet,
      previewRange: previewRange,
    );
  }
}
