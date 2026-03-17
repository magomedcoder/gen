import 'package:flutter/material.dart';
import 'package:gen/core/layout/responsive.dart';
import 'package:gen/presentation/screens/chat/bloc/chat_state.dart';
import 'package:gen/presentation/screens/chat/widgets/chat_model_selector.dart';
import 'package:gen/presentation/screens/chat/widgets/chat_supported_formats_button.dart';

class ChatToolbar extends StatelessWidget {
  const ChatToolbar({super.key, required this.state});

  final ChatState state;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: EdgeInsets.symmetric(
        horizontal: Breakpoints.isMobile(context) ? 12 : 20,
        vertical: 8,
      ),
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.surface,
        border: Border(
          bottom: BorderSide(
            color: Theme.of(context).dividerColor.withValues(alpha: 0.08),
          ),
        ),
      ),
      child: Row(
        children: [
          ChatModelSelector(state: state),
          const Spacer(),
          const ChatSupportedFormatsButton(),
        ],
      ),
    );
  }
}
