import 'package:gen/core/grpc_channel_manager.dart';
import 'package:gen/domain/entities/runner_info.dart' as domain;
import 'package:gen/generated/grpc_pb/runner.pb.dart' as pb;

abstract class IRunnersRemoteDataSource {
  Future<List<domain.RunnerInfo>> getRunners();

  Future<void> setRunnerEnabled(String address, bool enabled);

  Future<bool> getRunnersStatus();
}

class RunnersRemoteDataSource implements IRunnersRemoteDataSource {
  final GrpcChannelManager _channelManager;

  RunnersRemoteDataSource(this._channelManager);

  @override
  Future<List<domain.RunnerInfo>> getRunners() async {
    final resp = await _channelManager.runnerAdminClient.getRunners(pb.Empty());
    return resp.runners
        .map((r) => domain.RunnerInfo(address: r.address, enabled: r.enabled))
        .toList();
  }

  @override
  Future<void> setRunnerEnabled(String address, bool enabled) async {
    await _channelManager.runnerAdminClient.setRunnerEnabled(
      pb.SetRunnerEnabledRequest(address: address, enabled: enabled),
    );
  }

  @override
  Future<bool> getRunnersStatus() async {
    final resp = await _channelManager.runnerAdminClient.getRunnersStatus(pb.Empty());
    return resp.hasActiveRunners;
  }
}
