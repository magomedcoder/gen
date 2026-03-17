import 'package:flutter/material.dart';

class ChatRunnersInactiveBanner extends StatelessWidget {
  const ChatRunnersInactiveBanner({super.key});

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
      color: cs.errorContainer.withValues(alpha: 0.5),
      child: Text(
        'Нет активных раннеров. Чат недоступен.',
        style: TextStyle(color: cs.onErrorContainer),
      ),
    );
  }
}
