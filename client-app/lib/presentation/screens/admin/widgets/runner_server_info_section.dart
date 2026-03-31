import 'package:flutter/material.dart';
import 'package:gen/domain/entities/server_info.dart';

class RunnerServerInfoSection extends StatelessWidget {
  final ServerInfo serverInfo;
  final String? defaultModel;
  final ValueChanged<String>? onDefaultModelChanged;

  const RunnerServerInfoSection({
    super.key,
    required this.serverInfo,
    this.defaultModel,
    this.onDefaultModelChanged,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        if (serverInfo.models.isNotEmpty) ...[
          DropdownButtonFormField<String>(
            initialValue: defaultModel != null && serverInfo.models.contains(defaultModel)
                ? defaultModel
                : serverInfo.models.first,
            isExpanded: true,
            decoration: const InputDecoration(
              border: OutlineInputBorder(),
              isDense: true,
              labelText: 'Модель',
            ),
            items: [
              for (final model in serverInfo.models)
                DropdownMenuItem<String>(
                  value: model,
                  child: Text(model),
                ),
            ],
            onChanged: onDefaultModelChanged == null ? null : (String? v) {
              if (v != null) {
                onDefaultModelChanged!(v);
              }
            },
          ),
        ],
        const SizedBox(height: 12),
        const Divider(height: 1),
        const SizedBox(height: 12),
        Text(
          'Сервер',
          style: theme.textTheme.labelMedium?.copyWith(
            color: theme.colorScheme.onSurfaceVariant,
          ),
        ),
        const SizedBox(height: 6),
        Wrap(
          spacing: 12,
          runSpacing: 6,
          children: [
            if (serverInfo.hostname.isNotEmpty)
              RunnerInfoChip(
                icon: Icons.computer,
                label: 'Хост',
                value: serverInfo.hostname,
              ),
            if (serverInfo.os.isNotEmpty)
              RunnerInfoChip(
                icon: Icons.terminal,
                label: 'ОС',
                value: '${serverInfo.os}/${serverInfo.arch}',
              ),
            if (serverInfo.cpuCores > 0)
              RunnerInfoChip(
                icon: Icons.memory,
                label: 'Ядра CPU',
                value: '${serverInfo.cpuCores}',
              ),
            if (serverInfo.memoryTotalMb > 0)
              RunnerInfoChip(
                icon: Icons.storage,
                label: 'ОЗУ',
                value: '${serverInfo.memoryTotalMb} МБ',
              ),
          ],
        ),
      ],
    );
  }
}

class RunnerInfoChip extends StatelessWidget {
  final IconData icon;
  final String label;
  final String value;

  const RunnerInfoChip({
    super.key,
    required this.icon,
    required this.label,
    required this.value,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(icon, size: 16),
        const SizedBox(width: 4),
        Text(
          '$label: ',
          style: theme.textTheme.bodySmall?.copyWith(
            color: theme.colorScheme.onSurfaceVariant,
          ),
        ),
        Text(
          value,
          style: theme.textTheme.bodySmall?.copyWith(
            fontWeight: FontWeight.w600,
          ),
        ),
      ],
    );
  }
}
