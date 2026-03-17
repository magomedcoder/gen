import 'package:flutter/material.dart';
import 'package:gen/domain/entities/session.dart';
import 'package:gen/presentation/screens/chat/bloc/chat_state.dart';

class ChatAppBarTitle extends StatelessWidget {
  const ChatAppBarTitle({
    super.key,
    required this.state,
    required this.useDrawer,
    required this.isSidebarExpanded,
    required this.onToggleSidebar,
  });

  final ChatState state;
  final bool useDrawer;
  final bool isSidebarExpanded;
  final VoidCallback onToggleSidebar;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final currentSession = state.sessions.firstWhere(
      (session) => session.id == state.currentSessionId,
      orElse: () => ChatSession(
        id: 0,
        title: 'Новый чат',
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      ),
    );

    return Row(
      children: [
        if (!useDrawer)
          IconButton(
            icon: Icon(
              isSidebarExpanded ? Icons.menu_open : Icons.menu,
              color: theme.colorScheme.onSurfaceVariant,
            ),
            onPressed: onToggleSidebar,
            tooltip: isSidebarExpanded ? 'Скрыть меню' : 'Показать меню',
          ),
        if (!useDrawer) const SizedBox(width: 8),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(
                currentSession.title,
                style: TextStyle(
                  fontSize: useDrawer ? 18 : 16,
                  fontWeight: FontWeight.w600,
                ),
                overflow: TextOverflow.ellipsis,
              ),
              if (!state.isConnected)
                Row(
                  children: [
                    Icon(
                      Icons.wifi_off,
                      size: 12,
                      color: theme.colorScheme.error,
                    ),
                    const SizedBox(width: 4),
                    Text(
                      'Нет подключения',
                      style: TextStyle(
                        fontSize: 11,
                        color: theme.colorScheme.error,
                      ),
                    ),
                  ],
                ),
            ],
          ),
        ),
      ],
    );
  }
}
