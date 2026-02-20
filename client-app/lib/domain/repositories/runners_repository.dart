import 'package:gen/domain/entities/runner_info.dart';

abstract class RunnersRepository {
  Future<List<RunnerInfo>> getRunners();

  Future<void> setRunnerEnabled(String address, bool enabled);

  Future<bool> getRunnersStatus();
}
