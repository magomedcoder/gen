import 'dart:typed_data';

import 'package:gen/core/user_file_save_io.dart'
    if (dart.library.html) 'package:gen/core/user_file_save_stub.dart' as impl;

Future<bool> saveUserPickedFile(Uint8List bytes, String fileName) => impl.saveUserPickedFileImpl(bytes, fileName);
