import 'dart:convert';
import 'dart:js_interop';
import 'dart:typed_data';

import 'package:web/web.dart';

String safeDownloadBaseName(String fileName) {
  var name = fileName.trim();
  final slash = name.lastIndexOf('/');
  final bslash = name.lastIndexOf(r'\');
  final idx = slash > bslash ? slash : bslash;
  if (idx >= 0 && idx < name.length - 1) {
    name = name.substring(idx + 1);
  }

  if (name.isEmpty) {
    return 'download';
  }

  return name.replaceAll(RegExp(r'[/\\<>:"|?*\x00-\x1f]'), '_');
}

Future<bool> downloadUint8List(
  Uint8List bytes,
  String fileName, {
  required String mimeType,
}) async {
  final safe = safeDownloadBaseName(fileName);
  final blob = Blob(
    [bytes.toJS].toJS,
    BlobPropertyBag(type: mimeType),
  );
  final url = URL.createObjectURL(blob);
  final anchor = HTMLAnchorElement()
    ..href = url
    ..download = safe
    ..style.display = 'none';
  document.body?.appendChild(anchor);
  anchor.click();
  anchor.remove();
  URL.revokeObjectURL(url);
  return true;
}

Future<bool> downloadUtf8Text(
  String text,
  String fileName, {
  required String mimeType,
}) => downloadUint8List(
  Uint8List.fromList(utf8.encode(text)),
  fileName,
  mimeType: mimeType,
);
