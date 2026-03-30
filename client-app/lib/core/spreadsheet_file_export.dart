import 'dart:typed_data';

import 'package:gen/core/spreadsheet_file_export_io.dart'
    if (dart.library.html) 'package:gen/core/spreadsheet_file_export_stub.dart' as impl;

Future<bool> saveSpreadsheetToFile(Uint8List bytes, String fileName) => impl.saveSpreadsheetToFileImpl(bytes, fileName);

Future<bool> saveCsvToFile(String utf8Text, String fileName) => impl.saveCsvToFileImpl(utf8Text, fileName);
