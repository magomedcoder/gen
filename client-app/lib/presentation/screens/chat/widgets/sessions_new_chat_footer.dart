import 'package:flutter/material.dart';
import 'package:gen/core/layout/responsive.dart';

class SessionsNewChatFooter extends StatelessWidget {
  const SessionsNewChatFooter({
    super.key,
    required this.onPressed,
    required this.isInDrawer,
  });

  final VoidCallback onPressed;
  final bool isInDrawer;

  @override
  Widget build(BuildContext context) {
    final padding = isInDrawer && Breakpoints.isMobile(context)
        ? const EdgeInsets.symmetric(horizontal: 12, vertical: 12)
        : const EdgeInsets.all(16);

    return Container(
      padding: padding,
      decoration: BoxDecoration(
        border: Border(
          top: BorderSide(
            color: Theme.of(context).dividerColor.withValues(alpha: 0.1),
          ),
        ),
      ),
      child: ElevatedButton.icon(
        icon: const Icon(Icons.add, size: 18),
        label: const Text('Новый чат'),
        onPressed: onPressed,
        style: ElevatedButton.styleFrom(
          minimumSize: const Size(double.infinity, 48),
          backgroundColor: Theme.of(context).colorScheme.primary,
          foregroundColor: Theme.of(context).colorScheme.onPrimary,
        ),
      ),
    );
  }
}
