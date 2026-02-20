import 'package:equatable/equatable.dart';

class RunnerInfo extends Equatable {
  final String address;
  final bool enabled;

  const RunnerInfo({
    required this.address,
    required this.enabled,
  });

  @override
  List<Object?> get props => [address, enabled];
}
