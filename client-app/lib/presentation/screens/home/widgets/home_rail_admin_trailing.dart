import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:gen/presentation/screens/auth/bloc/auth_bloc.dart';
import 'package:gen/presentation/screens/auth/bloc/auth_state.dart';

class HomeRailAdminTrailing extends StatelessWidget {
  const HomeRailAdminTrailing({
    super.key,
    required this.onOpenUsersAdmin,
    required this.onOpenRunnersAdmin,
  });

  final VoidCallback onOpenUsersAdmin;
  final VoidCallback onOpenRunnersAdmin;

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<AuthBloc, AuthState>(
      builder: (context, authState) {
        if (!(authState.user?.isAdmin ?? false)) {
          return const SizedBox.shrink();
        }
        return Padding(
          padding: const EdgeInsets.only(bottom: 8),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Divider(
                height: 1,
                color: Theme.of(context).dividerColor.withValues(alpha: 0.2),
              ),
              const SizedBox(height: 8),
              Tooltip(
                message: 'Пользователи',
                child: IconButton(
                  icon: const Icon(Icons.supervisor_account_outlined),
                  onPressed: onOpenUsersAdmin,
                ),
              ),
              Tooltip(
                message: 'Раннеры',
                child: IconButton(
                  icon: const Icon(Icons.dns_outlined),
                  onPressed: onOpenRunnersAdmin,
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}
