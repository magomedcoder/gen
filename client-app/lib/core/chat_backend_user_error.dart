import 'package:gen/core/failures.dart';
import 'package:gen/core/grpc_unavailable.dart';
import 'package:grpc/grpc.dart';

String chatBackendErrorMessage(Object error) {
  if (error is Failure) {
    return error.message;
  }
  
  if (error is GrpcError) {
    return _grpcRunnerOrGenMessage(error);
  }

  final raw = error.toString().trim();
  if (raw.isEmpty) {
    return 'Неизвестная error';
  }

  return raw;
}

String? chatHeadlineIfBackendReachable(Object error, String headline) {
  if (isGrpcUnavailable(error)) {
    return null;
  }

  return '$headline: ${chatBackendErrorMessage(error)}';
}

String _grpcRunnerOrGenMessage(GrpcError e) {
  final detail = e.message?.trim();
  switch (e.code) {
    case StatusCode.deadlineExceeded:
      return 'превышено время ожидания gen или локального раннера. Проверьте timeout в настройках сессии и сеть.';
    case StatusCode.cancelled:
      return 'request отменён.';
    case StatusCode.unavailable:
      return 'gen или раннер недоступен по сети.';
    case StatusCode.unauthenticated:
      return 'требуется повторный вход (gen отклонил request).';
    case StatusCode.permissionDenied:
      return 'доступ запрещён на стороне gen.';
    case StatusCode.notFound:
      return 'ресурс не найден на стороне gen.';
    case StatusCode.invalidArgument:
      if (detail != null && detail.isNotEmpty) {
        return 'invalid request: $detail';
      }

      return 'invalid request к gen или раннеру.';
    case StatusCode.resourceExhausted:
      return 'лимит или размер requestа (gen).';
    case StatusCode.unimplemented:
      return 'операция не поддержана раннером или конфигурацией gen.';
    case StatusCode.internal:
      if (detail != null && detail.isNotEmpty) {
        return 'внутренняя error gen/раннера: $detail';
      }

      return 'внутренняя error gen или раннера (см. логи сервера).';
    default:
      if (detail != null && detail.isNotEmpty) {
        return 'error gen/раннера (${e.codeName}): $detail';
      }
      return 'error gen/раннера (${e.codeName}).';
  }
}
