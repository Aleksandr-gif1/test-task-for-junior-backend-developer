-- Добавляем поля периодичности в таблицу задач
ALTER TABLE tasks
  ADD COLUMN recurrence_type TEXT NOT NULL DEFAULT '',
  ADD COLUMN recurrence_config JSONB,
  ADD COLUMN recurrence_active BOOLEAN NOT NULL DEFAULT false;

-- Обновляем существующие записи: убираем пустой тип в NULL для консистентности
UPDATE tasks SET recurrence_type = NULL WHERE recurrence_type = '';

-- Создаём таблицу экземпляров периодических задач
CREATE TABLE task_instances (
  id BIGSERIAL PRIMARY KEY,
  parent_task_id BIGINT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
  scheduled_date DATE NOT NULL,
  title TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'new',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT task_instances_unique_task_date UNIQUE (parent_task_id, scheduled_date)
);

CREATE INDEX idx_task_instances_parent ON task_instances(parent_task_id);
CREATE INDEX idx_task_instances_scheduled ON task_instances(scheduled_date);
