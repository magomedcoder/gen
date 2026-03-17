import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:gen/core/injector.dart' as di;
import 'package:gen/core/layout/responsive.dart';
import 'package:gen/presentation/screens/admin/bloc/users_admin_bloc.dart';
import 'package:gen/presentation/screens/admin/bloc/users_admin_event.dart';
import 'package:gen/presentation/screens/admin/runners_admin_screen.dart';
import 'package:gen/presentation/screens/admin/users_admin_screen.dart';
import 'package:gen/presentation/screens/auth/bloc/auth_bloc.dart';
import 'package:gen/presentation/screens/chat/chat_screen.dart';
import 'package:gen/presentation/screens/editor/bloc/editor_bloc.dart';
import 'package:gen/presentation/screens/editor/bloc/editor_event.dart';
import 'package:gen/presentation/screens/editor/editor_screen.dart';
import 'package:gen/presentation/screens/home/widgets/home_mobile_admin_bar.dart';
import 'package:gen/presentation/screens/home/widgets/home_rail_admin_trailing.dart';

class HomeShell extends StatefulWidget {
  const HomeShell({super.key});

  @override
  State<HomeShell> createState() => _HomeShellState();
}

class _HomeShellState extends State<HomeShell> {
  int _index = 0;
  late final Widget _editorPage = BlocProvider(
    create: (_) => di.sl<EditorBloc>()..add(const EditorStarted()),
    child: const EditorScreen(),
  );

  void _openUsersAdmin() {
    final authBloc = context.read<AuthBloc>();
    Navigator.of(context).push(
      MaterialPageRoute<void>(
        builder: (_) => MultiBlocProvider(
          providers: [
            BlocProvider.value(value: authBloc),
            BlocProvider(
              create: (_) => di.sl<UsersAdminBloc>()
                ..add(const UsersAdminLoadRequested()),
            ),
          ],
          child: const UsersAdminScreen(),
        ),
      ),
    );
  }

  void _openRunnersAdmin() {
    final authBloc = context.read<AuthBloc>();
    Navigator.of(context).push(
      MaterialPageRoute<void>(
        builder: (_) => BlocProvider.value(
          value: authBloc,
          child: const RunnersAdminScreen(),
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final mobile = Breakpoints.isMobile(context);

    final pages = <Widget>[
      const ChatScreen(),
      _editorPage,
    ];

    if (mobile) {
      return Scaffold(
        body: IndexedStack(
          index: _index,
          children: pages,
        ),
        bottomNavigationBar: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            HomeMobileAdminBar(
              onOpenUsersAdmin: _openUsersAdmin,
              onOpenRunnersAdmin: _openRunnersAdmin,
            ),
            NavigationBar(
              selectedIndex: _index,
              onDestinationSelected: (i) => setState(() => _index = i),
              destinations: const [
                NavigationDestination(
                  icon: Icon(Icons.chat_bubble_outline),
                  selectedIcon: Icon(Icons.chat_rounded),
                  label: 'Чат',
                ),
                NavigationDestination(
                  icon: Icon(Icons.edit_note_outlined),
                  selectedIcon: Icon(Icons.edit_note_rounded),
                  label: 'Редактор',
                ),
              ],
            ),
          ],
        ),
      );
    }

    return Scaffold(
      body: Row(
        children: [
          NavigationRail(
            extended: false,
            selectedIndex: _index,
            onDestinationSelected: (i) => setState(() => _index = i),
            trailing: HomeRailAdminTrailing(
              onOpenUsersAdmin: _openUsersAdmin,
              onOpenRunnersAdmin: _openRunnersAdmin,
            ),
            destinations: const [
              NavigationRailDestination(
                icon: Icon(Icons.chat_bubble_outline),
                selectedIcon: Icon(Icons.chat_rounded),
                label: Text('Чат'),
              ),
              NavigationRailDestination(
                icon: Icon(Icons.edit_note_outlined),
                selectedIcon: Icon(Icons.edit_note_rounded),
                label: Text('Редактор'),
              ),
            ],
          ),
          const VerticalDivider(width: 1, thickness: 1),
          Expanded(
            child: IndexedStack(
              index: _index,
              children: pages,
            ),
          ),
        ],
      ),
    );
  }
}
