CREATE TABLE roles
(
    id String,
    name String,
    created_at DateTime
)
ENGINE = MergeTree
ORDER BY id;

CREATE TABLE executions
(
    id String,
    workflow_id String,
    workflow_name String,
    trigger_id String,
    status UInt8 DEFAULT 2, -- Enum8('pending' = 1, 'processing' = 2, 'completed' = 3, 'failed' = 4)
    start_time DateTime,
    end_time Nullable(DateTime),
    duration UInt32 MATERIALIZED IF(end_time > start_time, dateDiff('second', start_time, end_time), NULL),
    error_message Nullable(String)
)
ENGINE = MergeTree
PARTITION BY toYYYYMMDD(start_time)
ORDER BY (id, workflow_id, start_time);

CREATE TABLE executions_stats
(
    date Date,
    workflow_id String,
    total_executions UInt32,
    successful_executions UInt32,
    failed_executions UInt32,
    pending_executions UInt32,
    processing_executions UInt32,
    avg_duration_seconds Float64,
    min_duration_seconds Int32,
    max_duration_seconds Int32
)
ENGINE = SummingMergeTree
PARTITION BY toYYYYMM(date)
ORDER BY (date, workflow_id);

CREATE TABLE execution_steps
(
    id String,
    execution_id String,
    workflow_step_id String,
    status Enum8('pending' = 1, 'processing' = 2, 'completed' = 3, 'failed' = 4) DEFAULT 'processing',
    input_data String,
    output_data String,
    error_details Nullable(String),
    start_time DateTime,
    end_time Nullable(DateTime)
)
ENGINE = MergeTree
PARTITION BY toYYYYMMDD(start_time)
ORDER BY (id, execution_id, start_time);

CREATE TABLE permissions
(
    id String,
    role_id String,
    subject String,
    action String,
    conditions Nullable(String)
)
ENGINE = MergeTree
ORDER BY (id, role_id);

CREATE TABLE rate_limits
(
    id String,
    user_id String,
    resource_type String,
    max_requests UInt32,
    time_window UInt32,
    current_usage UInt32,
    last_reset_at DateTime
)
ENGINE = MergeTree
ORDER BY (id, user_id, resource_type);

CREATE TABLE recent_activity
(
    id String,
    user_id String,
    user_name String,
    activity_type Enum8('created_workflow' = 1, 'uploaded_file' = 2, 'downloaded_report' = 3, 'executed_workflow' = 4) DEFAULT 'created_workflow',
    activity_description String,
    related_workflow_id Nullable(String),
    timestamp DateTime
)
ENGINE = MergeTree
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (timestamp, user_id, activity_type);

CREATE TABLE scheduled_tasks
(
    id String,
    workflow_id String,
    schedule_type Enum8('once' = 1, 'recurring' = 2) DEFAULT 'once',
    cron_expression String,
    next_run_at DateTime,
    last_run_at Nullable(DateTime),
    created_at DateTime,
    updated_at DateTime
)
ENGINE = MergeTree
ORDER BY (id, workflow_id, next_run_at);

CREATE TABLE subscriptions
(
    id String,
    user_id String,
    plan_id String,
    status Enum8('active' = 1, 'canceled' = 2, 'expired' = 3) DEFAULT 'active',
    start_date Date,
    end_date Date,
    created_at DateTime,
    updated_at DateTime
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(start_date)
ORDER BY (id, user_id, start_date);

CREATE TABLE triggers
(
    id String,
    workflow_id String,
    app_integration_id String,
    trigger_type String,
    trigger_config String,
    created_at DateTime,
    updated_at DateTime
)
ENGINE = MergeTree
ORDER BY (id, workflow_id);

CREATE TABLE usage_data
(
    id String,
    user_id String,
    metric_name String,
    metric_value Float64,
    timestamp DateTime
)
ENGINE = MergeTree
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (id, user_id, metric_name, timestamp);

CREATE TABLE users
(
    id UUID DEFAULT generateUUIDv4(),
    stub String,
    status UInt8 DEFAULT 1, -- Enum8('active' = 1, 'invited' = 2, 'pending' = 3, 'blocked' = 4)
    resetPasswordToken Nullable(String),
    resetPasswordTokenSentAt Nullable(DateTime),
    invitationToken Nullable(String),
    invitationTokenSentAt Nullable(DateTime),
    trialExpiryDate Date,
    roleId LowCardinality(String) DEFAULT 2,
    deleted_at Nullable(DateTime),
    created_at DateTime DEFAULT now(),
    updated_at DateTime DEFAULT now()
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(created_at)
ORDER BY stub;

CREATE TABLE webhooks
(
    id String,
    user_id String,
    workflow_id String,
    url String,
    http_method Enum8('GET' = 1, 'POST' = 2, 'PUT' = 3, 'PATCH' = 4, 'DELETE' = 5),
    headers String,
    body_template String,
    created_at DateTime,
    updated_at DateTime
)
ENGINE = MergeTree
ORDER BY (id, user_id, workflow_id);

CREATE TABLE workflows
(
    id String,
    user_id String,
    name String,
    description String,
    is_active UInt8, -- Enum8('active' = 1, 'draft' = 2, 'paused' = 3) DEFAULT 'active'
    created_at DateTime,
    updated_at DateTime,
    workflow_init DateTime,
    workflow_completed DateTime,
    status UInt8 DEFAULT 1 -- Enum8('completed' = 1, 'pending' = 2, 'processing' = 3, 'failed' = 4) DEFAULT 'processing'
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(created_at)
ORDER BY (id, user_id);

CREATE TABLE workflow_stats
(
    date Date,
    total_workflows UInt32,
    successful_workflows UInt32,
    failed_workflows UInt32,
    pending_workflows UInt32,
    avg_duration_seconds Float64
)
ENGINE = SummingMergeTree
PARTITION BY toYYYYMM(date)
ORDER BY date;

CREATE TABLE workflow_steps
(
    id String,
    workflow_id String,
    step_order UInt32,
    key Nullable(String),
    appKey Nullable(String),
    type UInt8 DEFAULT 1, -- 'action' = 1, 'trigger' = 2
    connectionId Nullable(String),
    status UInt8 DEFAULT 1, -- Enum8('incomplete' = 1, 'completed' = 2) DEFAULT 'incomplete',
    parameters String,
    webhookPath Nullable(String),
    delete_at Nullable(DateTime64(3)),
    created_at DateTime64(3) DEFAULT now64(),
    updated_at DateTime64(3) DEFAULT now64()
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(created_at)
ORDER BY (workflow_id, step_order);

CREATE TABLE workflow_tags
(
    id String,
    workflow_id String,
    tag String,
    created_at DateTime
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(created_at)
ORDER BY (id, workflow_id, tag);

CREATE TABLE workflow_templates
(
    id String,
    name String,
    description String,
    category String,
    template_data String,
    created_by String,
    created_at DateTime,
    updated_at DateTime
)
ENGINE = MergeTree
ORDER BY (id, category);
