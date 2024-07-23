CREATE TABLE users (
    id UUID,
    fullName String,
    email String,
    password String,
    status Enum('active' = 1, 'invited' = 2),
    resetPasswordToken Nullable(String),
    resetPasswordTokenSentAt Nullable(DateTime),
    invitationToken Nullable(String),
    invitationTokenSentAt Nullable(DateTime),
    trialExpiryDate Date,
    roleId UUID,
    deletedAt Nullable(DateTime),
    createdAt DateTime,
    updatedAt DateTime,
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id);

CREATE TABLE roles (
    id UUID,
    name String,
    created_at DateTime,
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id);

CREATE TABLE permissions (
    id UUID,
    role_id UUID,
    subject String,
    action String,
    conditions Nullable(String),
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id);

CREATE TABLE api_keys (
    id UUID,
    user_id UUID,
    key_hash String,
    name String,
    expires_at DateTime,
    created_at DateTime,
    last_used_at Nullable(DateTime),
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id);

CREATE TABLE workflows (
    id UUID,
    user_id UUID,
    name String,
    description String,
    is_active UInt8,
    created_at DateTime,
    updated_at DateTime,
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id);

CREATE TABLE workflow_stats (
    date Date,
    total_workflows UInt32,
    successful_workflows UInt32,
    failed_workflows UInt32,
    pending_workflows UInt32,
    avg_duration_seconds Float64,
    PRIMARY KEY (date)
) ENGINE = SummingMergeTree() ORDER BY (date);

CREATE TABLE workflow_steps (
    id UUID,
    workflow_id UUID,
    step_order UInt32,
    app_integration_id UUID,
    action_name String,
    input_data String,
    output_data String,
    created_at DateTime,
    updated_at DateTime,
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id, workflow_id);

CREATE TABLE triggers (
    id UUID,
    workflow_id UUID,
    app_integration_id UUID,
    trigger_type String,
    trigger_config String,
    created_at DateTime,
    updated_at DateTime,
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id);

CREATE TABLE executions (
    id UUID,
    workflow_id UUID,
    workflow_name String,
    trigger_id UUID,
    status Enum('pending' = 1, 'processing' = 2, 'completed' = 3, 'failed' = 4),
    start_time DateTime,
    end_time Nullable(DateTime),
    duration_seconds Int32,
    error_message Nullable(String),
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id, start_time);

CREATE TABLE execution_steps (
    id UUID,
    execution_id UUID,
    workflow_step_id UUID,
    status Enum('pending' = 1, 'processing' = 2, 'completed' = 3, 'failed' = 4),
    input_data String,
    output_data String,
    error_details Nullable(String),
    start_time DateTime,
    end_time Nullable(DateTime),
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id, execution_id);

CREATE TABLE app_integrations (
    id UUID,
    name String,
    description String,
    auth_type Enum('oauth' = 1, 'api_key' = 2, 'basic' = 3),
    auth_config String,
    icon_url String,
    created_at DateTime,
    updated_at DateTime,
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id);

CREATE TABLE app_actions (
    id UUID,
    app_integration_id UUID,
    name String,
    description String,
    input_schema String,
    output_schema String,
    created_at DateTime,
    updated_at DateTime,
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id);

CREATE TABLE user_app_integrations (
    id UUID,
    user_id UUID,
    app_integration_id UUID,
    auth_data String,
    created_at DateTime,
    updated_at DateTime,
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id);

CREATE TABLE webhooks (
    id UUID,
    user_id UUID,
    workflow_id UUID,
    url String,
    http_method Enum('GET' = 1, 'POST' = 2, 'PUT' = 3, 'PATCH' = 4, 'DELETE' = 5),
    headers String,
    body_template String,
    created_at DateTime,
    updated_at DateTime,
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id);

CREATE TABLE audit_logs (
    id UUID,
    user_id UUID,
    action String,
    details String,
    ip_address String,
    user_agent String,
    timestamp DateTime,
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id);

CREATE TABLE system_metrics (
    id UUID,
    metric_name String,
    metric_value Float64,
    timestamp DateTime,
    PRIMARY KEY (id, timestamp)
) ENGINE = MergeTree() ORDER BY (id, timestamp);

CREATE TABLE rate_limits (
    id UUID,
    user_id UUID,
    resource_type String,
    max_requests UInt32,
    time_window UInt32,
    current_usage UInt32,
    last_reset_at DateTime,
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id, user_id);

CREATE TABLE notifications (
    id UUID,
    user_id UUID,
    type Enum('email' = 1, 'in_app' = 2, 'sms' = 3),
    subject String,
    content String,
    sent_at DateTime,
    read_at Nullable(DateTime),
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id);

CREATE TABLE alerts (
    id UUID,
    workflow_id UUID,
    type Enum('error' = 1, 'warning' = 2, 'info' = 3),
    message String,
    created_at DateTime,
    resolved_at Nullable(DateTime),
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id);

CREATE TABLE subscriptions (
    id UUID,
    user_id UUID,
    plan_id String,
    status Enum('active' = 1, 'canceled' = 2, 'expired' = 3),
    start_date Date,
    end_date Date,
    created_at DateTime,
    updated_at DateTime,
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id);

CREATE TABLE usage_data (
    id UUID,
    user_id UUID,
    metric_name String,
    metric_value Float64,
    timestamp DateTime,
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id, user_id, timestamp);


CREATE TABLE workflow_templates (
    id UUID,
    name String,
    description String,
    category String,
    template_data String,
    created_by UUID,
    created_at DateTime,
    updated_at DateTime,
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id);

CREATE TABLE workflow_tags (
    id UUID,
    workflow_id UUID,
    tag String,
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id, workflow_id);

CREATE TABLE workflow_versions (
    id UUID,
    workflow_id UUID,
    version_number UInt32,
    workflow_data String,
    created_at DateTime,
    created_by UUID,
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id, workflow_id);

CREATE TABLE custom_events (
    id UUID,
    user_id UUID,
    event_name String,
    event_data String,
    created_at DateTime,
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id, user_id);

CREATE TABLE scheduled_tasks (
    id UUID,
    workflow_id UUID,
    schedule_type Enum('once' = 1, 'recurring' = 2),
    cron_expression String,
    next_run_at DateTime,
    last_run_at Nullable(DateTime),
    created_at DateTime,
    updated_at DateTime,
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id, next_run_at);

CREATE TABLE debug_logs (
    id UUID,
    execution_id UUID,
    workflow_step_id UUID,
    log_level Enum('debug' = 1, 'info' = 2, 'warning' = 3, 'error' = 4),
    message String,
    timestamp DateTime,
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id, execution_id, timestamp);

CREATE TABLE app_connections (
    id UUID,
    user_id UUID,
    source_app_id UUID,
    target_app_id UUID,
    connection_data String,
    created_at DateTime,
    updated_at DateTime,
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id, user_id);

CREATE TABLE user_feedback (
    id UUID,
    user_id UUID,
    feedback_type Enum('bug' = 1, 'feature_request' = 2, 'general' = 3),
    content String,
    status Enum('open' = 1, 'in_progress' = 2, 'closed' = 3),
    created_at DateTime,
    updated_at DateTime,
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id, user_id);

CREATE TABLE global_config (
    id UUID,
    config_key String,
    config_value String,
    created_at DateTime,
    updated_at DateTime,
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (id);

CREATE TABLE recent_activity (
    id UUID,
    user_id UUID,
    user_name String,
    activity_type Enum('created_workflow' = 1, 'uploaded_file' = 2, 'downloaded_report' = 3, 'executed_workflow' = 4),
    activity_description String,
    related_workflow_id Nullable(UUID),
    timestamp DateTime,
    PRIMARY KEY (id)
) ENGINE = MergeTree() ORDER BY (timestamp DESC);

