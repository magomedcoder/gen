import 'dart:typed_data';

import 'package:gen/core/web_blob_download.dart';

Future<bool> saveDocxToFileImpl(Uint8List bytes, String fileName) => downloadUint8List(
    bytes,
    fileName,
    mimeType: 'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
);
