import 'package:desktop_drop/desktop_drop.dart';
import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:gen/core/layout/responsive.dart';
import 'package:gen/domain/entities/session.dart';
import 'package:gen/presentation/screens/chat/bloc/chat_bloc.dart';
import 'package:gen/presentation/screens/chat/bloc/chat_event.dart';
import 'package:gen/presentation/screens/chat/bloc/chat_state.dart';
import 'package:gen/presentation/screens/chat/widgets/chat_app_bar_title.dart';
import 'package:gen/presentation/screens/chat/widgets/chat_dialogs.dart';
import 'package:gen/presentation/screens/chat/widgets/chat_input_bar.dart';
import 'package:gen/presentation/screens/chat/widgets/chat_messages_panel.dart';
import 'package:gen/presentation/screens/chat/widgets/sessions_sidebar.dart';

class ChatScreen extends StatefulWidget {
  const ChatScreen({super.key});

  @override
  State<ChatScreen> createState() => _ChatScreenState();
}

class _ChatScreenState extends State<ChatScreen> {
  final _scrollController = ScrollController();
  final _scaffoldKey = GlobalKey<ScaffoldState>();
  final _inputBarKey = GlobalKey<ChatInputBarState>();
  final TextEditingController _sessionTitleController = TextEditingController();
  bool _isSidebarExpanded = true;
  bool _isDraggingFile = false;
  double get _sidebarWidth => Breakpoints.sidebarDefaultWidth;

  static const double _scrollThreshold = 80.0;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<ChatBloc>().add(ChatStarted());
    });
  }

  void _scrollToBottom() {
    if (!mounted) {
      return;
    }

    if (!_scrollController.hasClients) {
      return;
    }

    final pos = _scrollController.position;
    if (pos.maxScrollExtent - pos.pixels <= _scrollThreshold) {
      _scrollController.animateTo(
        pos.maxScrollExtent,
        duration: const Duration(milliseconds: 300),
        curve: Curves.easeOut,
      );
    }
  }

  void _toggleSidebar() {
    setState(() {
      _isSidebarExpanded = !_isSidebarExpanded;
    });
  }

  Future<void> _createNewSession() async {
    final result = await showNewSessionDialog(context, _sessionTitleController);
    if (!mounted) {
      return;
    }

    if (result != null) {
      context.read<ChatBloc>().add(ChatCreateSession(title: result));
      _sessionTitleController.clear();
    }
  }

  void _selectSession(ChatSession session) {
    context.read<ChatBloc>().add(ChatSelectSession(session.id));
  }

  void _selectSessionAndCloseDrawer(ChatSession session) {
    _selectSession(session);
    if (Breakpoints.useDrawerForSessions(context)) {
      Navigator.of(context).pop();
    }
  }

  void _deleteSession(int sessionId, String sessionTitle) {
    showDeleteSessionDialog(
      context,
      sessionId: sessionId,
      sessionTitle: sessionTitle,
      chatBloc: context.read<ChatBloc>(),
    );
  }

  Future<void> _onFilesDropped(DropDoneDetails details) async {
    setState(() => _isDraggingFile = false);
    if (details.files.isEmpty) {
      return;
    }

    final item = details.files.first;
    if (item is! DropItemFile) {
      return;
    }

    try {
      final bytes = await item.readAsBytes();
      final name = item.name.isNotEmpty
        ? item.name
        : item.path.split(RegExp(r'[/\\]')).last;
      if (!mounted) return;
      _inputBarKey.currentState?.setDroppedFile(
        PlatformFile(
          name: name,
          size: bytes.length,
          bytes: bytes,
        ),
      );
    } catch (_) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Не удалось прочитать файл')),
        );
      }
    }
  }

  @override
  void dispose() {
    _scrollController.dispose();
    _sessionTitleController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return BlocListener<ChatBloc, ChatState>(
      listener: (context, state) {
        WidgetsBinding.instance.addPostFrameCallback((_) {
          if (state.messages.isNotEmpty) {
            _scrollToBottom();
          }
        });

        if (state.error != null) {
          WidgetsBinding.instance.addPostFrameCallback((_) {
            if (!mounted) return;
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.error!),
                backgroundColor: Theme.of(context).colorScheme.error,
                behavior: SnackBarBehavior.floating,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(8),
                ),
              ),
            );
          });
        }
      },
      child: Builder(
        builder: (context) {
          final useDrawer = Breakpoints.useDrawerForSessions(context);
          return Scaffold(
            key: _scaffoldKey,
            drawer: useDrawer
              ? Drawer(
                child: SafeArea(
                  child: SessionsSidebar(
                    isInDrawer: true,
                    onCreateNewSession: _createNewSession,
                    onSelectSession: _selectSessionAndCloseDrawer,
                    onDeleteSession: _deleteSession,
                  ),
                ),
              )
              : null,
            appBar: AppBar(
              leading: useDrawer
                ? IconButton(
                  icon: const Icon(Icons.menu),
                  onPressed: () => _scaffoldKey.currentState?.openDrawer(),
                  tooltip: 'Меню сессий',
                )
                : null,
              title: BlocBuilder<ChatBloc, ChatState>(
                builder: (context, state) => ChatAppBarTitle(
                  state: state,
                  useDrawer: useDrawer,
                  isSidebarExpanded: _isSidebarExpanded,
                  onToggleSidebar: _toggleSidebar,
                ),
              ),
              actions: [
                BlocBuilder<ChatBloc, ChatState>(
                  builder: (context, state) {
                    if (state.isLoading && !state.isStreaming) {
                      return const Padding(
                        padding: EdgeInsets.only(right: 16),
                        child: SizedBox(
                          width: 16,
                          height: 16,
                          child: CircularProgressIndicator(strokeWidth: 2),
                        ),
                      );
                    }
                    return const SizedBox();
                  },
                ),
              ],
            ),
            body: SafeArea(
              top: false,
              bottom: true,
              left: false,
              right: false,
              child: Row(
                children: [
                  if (!useDrawer)
                    AnimatedContainer(
                      duration: const Duration(milliseconds: 300),
                      width: _isSidebarExpanded ? _sidebarWidth : 0,
                      curve: Curves.easeInOut,
                      decoration: BoxDecoration(
                        border: Border(
                          right: BorderSide(
                            color: Theme.of(context)
                              .dividerColor
                              .withValues(alpha: 0.1),
                            width: 1,
                          ),
                        ),
                      ),
                      child: _isSidebarExpanded
                        ? SessionsSidebar(
                          onCreateNewSession: _createNewSession,
                          onSelectSession: _selectSession,
                          onDeleteSession: _deleteSession,
                        )
                        : const SizedBox.shrink(),
                    ),
                  Expanded(
                    child: BlocBuilder<ChatBloc, ChatState>(
                      builder: (context, state) {
                        final canDropFile = state.isConnected && !state.isLoading && (state.hasActiveRunners != false);
                        return ChatMessagesPanel(
                          state: state,
                          scrollController: _scrollController,
                          inputBarKey: _inputBarKey,
                          isDraggingFile: _isDraggingFile,
                          canDropFile: canDropFile,
                          onDragEntered: (_) => setState(() => _isDraggingFile = true),
                          onDragExited: (_) => setState(() => _isDraggingFile = false),
                          onDragDone: _onFilesDropped,
                        );
                      },
                    ),
                  ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }
}
