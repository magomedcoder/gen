import 'package:flutter/material.dart';
import 'package:gen/presentation/screens/chat/widgets/chat_dialogs.dart';

class ChatSupportedFormatsButton extends StatelessWidget {
  const ChatSupportedFormatsButton({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Tooltip(
      message: 'Поддерживаемые форматы вложений',
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          onTap: () => showSupportedFormatsDialog(context),
          borderRadius: BorderRadius.circular(6),
          child: Container(
            padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 5),
            decoration: BoxDecoration(
              color: theme.colorScheme.surfaceContainerLow,
              borderRadius: BorderRadius.circular(6),
              border: Border.all(
                color: theme.colorScheme.outline.withValues(alpha: 0.2),
                width: 1,
              ),
            ),
            child: Icon(
              Icons.help_outline,
              size: 16,
              color: theme.colorScheme.onSurfaceVariant,
            ),
          ),
        ),
      ),
    );
  }
}
