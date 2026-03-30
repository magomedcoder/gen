import 'dart:typed_data';

import 'package:gen/core/web_blob_download.dart';

Future<bool> saveUserPickedFileImpl(Uint8List bytes, String fileName) => downloadUint8List(
    bytes,
    fileName,
    mimeType: 'application/octet-stream'
);
