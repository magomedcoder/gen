import 'dart:typed_data';

import 'package:gen/core/web_blob_download.dart';

Future<bool> saveSpreadsheetToFileImpl(Uint8List bytes, String fileName) => downloadUint8List(
  bytes,
  fileName,
  mimeType: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
);

Future<bool> saveCsvToFileImpl(String utf8Text, String fileName) => downloadUtf8Text(
  utf8Text,
  fileName,
  mimeType: 'text/csv;charset=utf-8',
);
