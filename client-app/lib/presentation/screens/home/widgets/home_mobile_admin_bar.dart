import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:gen/presentation/screens/auth/bloc/auth_bloc.dart';
import 'package:gen/presentation/screens/auth/bloc/auth_state.dart';

class HomeMobileAdminBar extends StatelessWidget {
  const HomeMobileAdminBar({
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
        final cs = Theme.of(context).colorScheme;
        return Material(
          elevation: 2,
          color: cs.surfaceContainerHighest,
          child: SafeArea(
            top: false,
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 6),
              child: Row(
                children: [
                  Text(
                    'Админ',
                    style: Theme.of(context).textTheme.labelSmall?.copyWith(
                          color: cs.onSurfaceVariant,
                          fontWeight: FontWeight.w600,
                        ),
                  ),
                  const Spacer(),
                  TextButton.icon(
                    icon: const Icon(Icons.supervisor_account_outlined, size: 18),
                    label: const Text('Пользователи'),
                    onPressed: onOpenUsersAdmin,
                  ),
                  TextButton.icon(
                    icon: const Icon(Icons.dns_outlined, size: 18),
                    label: const Text('Раннеры'),
                    onPressed: onOpenRunnersAdmin,
                  ),
                ],
              ),
            ),
          ),
        );
      },
    );
  }
}
