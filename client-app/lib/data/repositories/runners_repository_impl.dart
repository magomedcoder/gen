import 'package:gen/domain/entities/runner_info.dart';
import 'package:gen/domain/repositories/runners_repository.dart';
import 'package:gen/data/data_sources/remote/runners_remote_datasource.dart';

class RunnersRepositoryImpl implements RunnersRepository {
  final IRunnersRemoteDataSource _remote;

  RunnersRepositoryImpl(this._remote);

  @override
  Future<List<RunnerInfo>> getRunners() => _remote.getRunners();

  @override
  Future<void> setRunnerEnabled(String address, bool enabled) => _remote.setRunnerEnabled(address, enabled);

  @override
  Future<bool> getRunnersStatus() => _remote.getRunnersStatus();
}
