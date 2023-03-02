-- +goose Up
-- +goose StatementBegin
create table if not exists role (
    id bigserial,
    name varchar(255) not null,
    type integer default 1,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id)
);
create unique index if not exists role_type_index on role (type);
create table if not exists organization (
    id bigserial,
    name varchar(255) not null,
    slug varchar(255) not null,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id)
);
create unique index if not exists organization_slug_index on organization (slug);
create table if not exists "user" (
    id bigserial,
    name varchar(255) not null,
    email varchar(255) not null,
    picture_url varchar(255),
    password varchar(255),
    payment_customer_id varchar(255),
    provider_uid varchar(255),
    provider integer default 1,
    last_login_at timestamp(0) default now(),
    role_id bigint,
    organization_id bigint,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id),
    foreign key (role_id) references role,
    foreign key (organization_id) references organization on delete cascade
);
create unique index if not exists user_email_index on "user" (email);
create unique index if not exists user_provider_uid_index on "user" (provider_uid);
create unique index if not exists user_payment_customer_id_index on "user" (payment_customer_id);
create index if not exists user_role_id_index on "user" (role_id);
create index if not exists user_organization_id_index on "user" (organization_id);
create index if not exists user_last_login_at_index on "user" (last_login_at);
create table if not exists region (
    id bigserial,
    name varchar(255) not null,
    key varchar(255) not null,
    ip_address varchar(255),
    "default" boolean default false,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id)
);
create unique index if not exists region_key_index on region (key);
create table if not exists monitor_group (
    id bigserial,
    name varchar(255) not null,
    organization_id bigint,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id),
    foreign key (organization_id) references organization on delete cascade
);
create unique index if not exists monitor_group_name_organization_id_index on monitor_group (name, organization_id);
create index if not exists monitor_group_organization_id_index on monitor_group (organization_id);
create table if not exists monitor (
    id bigserial,
    name varchar(255) not null,
    url varchar(255) not null,
    method integer default 1,
    interval integer default 300,
    timeout integer default 5,
    type integer default 1,
    body text,
    body_format integer default 1,
    headers jsonb default '{}'::jsonb,
    username text,
    password text,
    "on" boolean default true,
    muted boolean default false,
    status integer default 1,
    check_ssl boolean default false,
    follow_redirects boolean default false,
    next_check_at timestamp(0),
    last_checked_at timestamp(0),
    last_failed_at timestamp(0),
    user_id bigint,
    monitor_group_id bigint,
    prev_id bigint,
    organization_id bigint,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id),
    constraint monitor_unique_previous_id unique (prev_id, organization_id) deferrable initially deferred,
    foreign key (user_id) references "user",
    foreign key (monitor_group_id) references monitor_group,
    foreign key (prev_id) references monitor on delete cascade,
    foreign key (organization_id) references organization on delete cascade
);
create unique index if not exists monitor_url_organization_id_index on monitor (url, organization_id);
create index if not exists monitor_on_index on monitor ("on");
create index if not exists monitor_status_index on monitor (status);
create index if not exists monitor_user_id_index on monitor (user_id);
create index if not exists monitor_organization_id_index on monitor (organization_id);
create index if not exists monitor_monitor_group_id_index on monitor (monitor_group_id);
create index if not exists monitor_next_check_at_index on monitor (next_check_at);
create index if not exists monitor_last_checked_at_next_check_at_index on monitor (last_checked_at, next_check_at);
create table if not exists assertion (
    id bigserial,
    source integer default 1,
    property varchar(255),
    comparison integer default 1,
    value varchar(255),
    monitor_id bigint,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id),
    foreign key (monitor_id) references monitor on delete cascade
);
create index if not exists assertion_monitor_id_index on assertion (monitor_id);
create unique index if not exists assertion_source_value_monitor_id_index on assertion (source, value, monitor_id);
create table if not exists monitor_region (
    id bigserial,
    down boolean default false,
    last_checked_at timestamp(0),
    monitor_id bigint,
    region_id bigint,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id),
    foreign key (monitor_id) references monitor on delete cascade,
    foreign key (region_id) references region on delete cascade
);
create index if not exists monitor_region_monitor_id_index on monitor_region (monitor_id);
create index if not exists monitor_region_region_id_index on monitor_region (region_id);
create index if not exists monitor_region_last_checked_at_index on monitor_region (last_checked_at);
create unique index if not exists monitor_region_region_id_monitor_id_index on monitor_region (region_id, monitor_id);
create table if not exists "check" (
    id bigserial,
    status_code integer,
    duration integer default 0,
    success boolean default false not null,
    region_id bigint,
    monitor_id bigint,
    organization_id bigint,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id),
    foreign key (region_id) references region,
    foreign key (monitor_id) references monitor on delete cascade,
    foreign key (organization_id) references organization on delete cascade
);
create index if not exists check_region_id_index on "check" (region_id);
create index if not exists check_monitor_id_index on "check" (monitor_id);
create index if not exists check_organization_id_index on "check" (organization_id);
create table if not exists monitor_integration (
    id bigserial,
    name varchar(255),
    type integer,
    config jsonb,
    organization_id bigint,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id),
    foreign key (organization_id) references organization on delete cascade
);
create index if not exists monitor_integration_organization_id_index on monitor_integration (organization_id);
create unique index if not exists monitor_integration_type_organization_id_index on monitor_integration (type, organization_id);
create table if not exists error_log (
    id bigserial,
    text text,
    type integer,
    screenshot_url varchar(255),
    check_id bigint,
    monitor_id bigint,
    assertion_id bigint,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id),
    foreign key (check_id) references "check" on delete cascade,
    foreign key (monitor_id) references monitor on delete cascade,
    foreign key (assertion_id) references assertion
);
create index if not exists error_log_check_id_index on error_log (check_id);
create index if not exists error_log_monitor_id_index on error_log (monitor_id);
create index if not exists error_log_assertion_id_index on error_log (assertion_id);
create table if not exists monitor_notification_policy (
    id bigserial,
    user_id bigint,
    monitor_id bigint,
    organization_id bigint,
    integration_id bigint,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id),
    foreign key (user_id) references "user" on delete cascade,
    foreign key (monitor_id) references monitor on delete cascade,
    foreign key (organization_id) references organization on delete cascade,
    foreign key (integration_id) references monitor_integration on delete cascade
);
create index if not exists monitor_notification_policy_user_id_index on monitor_notification_policy (user_id);
create index if not exists monitor_notification_policy_monitor_id_index on monitor_notification_policy (monitor_id);
create index if not exists monitor_notification_policy_integration_id_index on monitor_notification_policy (integration_id);
create index if not exists monitor_notification_policy_organization_id_index on monitor_notification_policy (organization_id);
create unique index if not exists monitor_notification_policy_user_id_monitor_id_integration_id_o on monitor_notification_policy (
    user_id,
    monitor_id,
    integration_id,
    organization_id
);
create table if not exists alarm (
    id bigserial,
    ongoing boolean,
    resolved_at timestamp(0),
    triggered_by_check_id bigint,
    resolved_by_check_id bigint,
    monitor_id bigint,
    organization_id bigint,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id),
    foreign key (triggered_by_check_id) references "check",
    foreign key (resolved_by_check_id) references "check",
    foreign key (monitor_id) references monitor on delete cascade,
    foreign key (organization_id) references organization on delete cascade
);
create index if not exists alarm_monitor_id_index on alarm (monitor_id);
create index if not exists alarm_organization_id_index on alarm (organization_id);
create unique index if not exists alarm_triggered_by_check_id_index on alarm (triggered_by_check_id);
create unique index if not exists uq_monitor_on_alarm on alarm (monitor_id, ongoing)
where (ongoing = true);
create table if not exists monitor_alarm_policy (
    id bigserial,
    reason varchar(255),
    threshold integer default 0,
    monitor_id bigint,
    organization_id bigint,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id),
    foreign key (monitor_id) references monitor on delete cascade,
    foreign key (organization_id) references organization on delete cascade
);
create index if not exists monitor_alarm_policy_monitor_id_index on monitor_alarm_policy (monitor_id);
create index if not exists monitor_alarm_policy_organization_id_index on monitor_alarm_policy (organization_id);
create unique index if not exists monitor_alarm_policy_reason_monitor_id_organization_id_index on monitor_alarm_policy (reason, monitor_id, organization_id);
create table if not exists daily_report (
    id bigserial,
    successful_checks integer default 0,
    error_checks integer default 0,
    downtime integer default 0,
    date date,
    monitor_id bigint,
    organization_id bigint,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id),
    foreign key (monitor_id) references monitor on delete cascade,
    foreign key (organization_id) references organization on delete cascade
);
create index if not exists daily_report_monitor_id_index on daily_report (monitor_id);
create index if not exists daily_report_organization_id_index on daily_report (organization_id);
create unique index if not exists daily_report_date_monitor_id_index on daily_report (date, monitor_id);
create table if not exists user_contact (
    id bigserial,
    email varchar(255),
    number varchar(255),
    mode integer,
    device_id varchar(255),
    verification_code varchar(255),
    verification_code_expires_at timestamp(0),
    verified boolean default false not null,
    subscribed boolean default true not null,
    bounce_count integer default 0,
    user_id bigint,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id),
    foreign key (user_id) references "user" on delete cascade
);
create index if not exists user_contact_user_id_index on user_contact (user_id);
create unique index if not exists user_contact_email_verified_index on user_contact (email, verified);
create unique index if not exists user_contact_number_verified_index on user_contact (number, verified);
create unique index if not exists user_contact_device_id_index on user_contact (device_id);
create table if not exists monitor_notification (
    id bigserial,
    type integer,
    successful boolean default true not null,
    alarm_id bigint,
    monitor_id bigint,
    user_contact_id bigint,
    organization_id bigint,
    integration_id bigint,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id),
    foreign key (alarm_id) references alarm on delete cascade,
    foreign key (monitor_id) references monitor on delete cascade,
    foreign key (user_contact_id) references user_contact on delete cascade,
    foreign key (organization_id) references organization on delete cascade,
    foreign key (integration_id) references monitor_integration on delete cascade
);
create index if not exists monitor_notification_alarm_id_index on monitor_notification (alarm_id);
create index if not exists monitor_notification_monitor_id_index on monitor_notification (monitor_id);
create index if not exists monitor_notification_user_contact_id_index on monitor_notification (user_contact_id);
create index if not exists monitor_notification_organization_id_index on monitor_notification (organization_id);
create index if not exists monitor_notification_integration_id_index on monitor_notification (integration_id);
create unique index if not exists monitor_notification_alarm_id_type_user_contact_id_integration_ on monitor_notification (alarm_id, type, user_contact_id, integration_id);
create table if not exists guest_user (
    id bigserial,
    email varchar(255) not null,
    code varchar(255) not null,
    expires_at timestamp(0) not null,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id)
);
create unique index if not exists guest_user_code_index on guest_user (code);
create index if not exists guest_user_expires_at_index on guest_user (expires_at);
create table if not exists invitation (
    id bigserial,
    email varchar(255) not null,
    code varchar(255) not null,
    expires_at timestamp(0) not null,
    notification_count integer default 1,
    invited_by_user_id bigint,
    role_id bigint,
    organization_id bigint,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id),
    foreign key (invited_by_user_id) references "user",
    foreign key (role_id) references role,
    foreign key (organization_id) references organization on delete cascade
);
create unique index if not exists invitation_code_index on invitation (code);
create unique index if not exists invitation_email_organization_id_index on invitation (email, organization_id);
create index if not exists invitation_invited_by_user_id_index on invitation (invited_by_user_id);
create index if not exists invitation_email_index on invitation (email);
create index if not exists invitation_expires_at_index on invitation (expires_at);
create index if not exists invitation_role_id_index on invitation (role_id);
create index if not exists invitation_organization_id_index on invitation (organization_id);
create table if not exists organization_user (
    id bigserial,
    status integer default 1,
    role_id bigint,
    user_id bigint,
    organization_id bigint,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id),
    foreign key (role_id) references role,
    foreign key (user_id) references "user" on delete cascade,
    foreign key (organization_id) references organization on delete cascade
);
create unique index if not exists organization_user_user_id_organization_id_index on organization_user (user_id, organization_id);
create index if not exists organization_user_role_id_index on organization_user (role_id);
create index if not exists organization_user_user_id_index on organization_user (user_id);
create index if not exists organization_user_organization_id_index on organization_user (organization_id);
create unique index if not exists uq_superadmin_on_org_user on organization_user (role_id, user_id)
where (role_id = 1);
create table if not exists product (
    id bigserial,
    name varchar(255) not null,
    description varchar(255),
    external_id varchar(255),
    tier integer default 1,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id)
);
create unique index if not exists product_name_index on product (name);
create unique index if not exists product_tier_index on product (tier);
create unique index if not exists product_external_id_index on product (external_id);
create table if not exists plan (
    id bigserial,
    price double precision not null,
    type integer default 1,
    external_id varchar(255),
    product_id bigint,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id),
    foreign key (product_id) references product on delete cascade
);
create index if not exists plan_product_id_index on plan (product_id);
create unique index if not exists plan_external_id_index on plan (external_id);
create unique index if not exists plan_price_type_index on plan (price, type);
create table if not exists subscription (
    id bigserial,
    status integer,
    starts_at timestamp(0),
    expires_at timestamp(0),
    canceled_at timestamp(0),
    is_trial boolean default false,
    external_id varchar(255),
    external_customer_id varchar(255),
    plan_id bigint,
    product_id bigint,
    organization_id bigint,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id),
    foreign key (plan_id) references plan,
    foreign key (product_id) references product,
    foreign key (organization_id) references organization on delete cascade
);
create index if not exists subscription_status_index on subscription (status);
create index if not exists subscription_expires_at_index on subscription (expires_at);
create unique index if not exists subscription_external_id_index on subscription (external_id);
create index if not exists subscription_plan_id_index on subscription (plan_id);
create index if not exists subscription_product_id_index on subscription (product_id);
create index if not exists subscription_organization_id_index on subscription (organization_id);
create table if not exists receipt (
    id bigserial,
    price double precision not null,
    currency varchar(255) default 'usd'::character varying,
    external_id varchar(255),
    external_customer_id varchar(255),
    url varchar(255),
    status integer,
    paid boolean default false,
    paid_at timestamp(0),
    "from" date,
    "to" date,
    is_trial boolean default false,
    plan_id bigint,
    product_id bigint,
    subscription_id bigint,
    organization_id bigint,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id),
    foreign key (plan_id) references plan,
    foreign key (product_id) references product,
    foreign key (subscription_id) references subscription on delete cascade,
    foreign key (organization_id) references organization on delete cascade
);
create unique index if not exists receipt_external_id_index on receipt (external_id);
create index if not exists receipt_plan_id_index on receipt (plan_id);
create index if not exists receipt_product_id_index on receipt (product_id);
create index if not exists receipt_subscription_id_index on receipt (subscription_id);
create index if not exists receipt_organization_id_index on receipt (organization_id);
create table if not exists role_claim (
    id bigserial,
    claim varchar(255) not null,
    role_id bigint,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id),
    foreign key (role_id) references role on delete cascade
);
create index if not exists role_claim_role_id_index on role_claim (role_id);
create unique index if not exists role_claim_claim_role_id_index on role_claim (claim, role_id);
create table if not exists feature (
    id bigserial,
    name varchar(255) not null,
    type integer default 1,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id)
);
create unique index if not exists feature_name_type_index on feature (name, type);
create table if not exists product_feature (
    id bigserial,
    count integer default 1,
    product_id bigint,
    feature_id bigint,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id),
    foreign key (product_id) references product on delete cascade,
    foreign key (feature_id) references feature on delete cascade
);
create index if not exists product_feature_product_id_index on product_feature (product_id);
create index if not exists product_feature_feature_id_index on product_feature (feature_id);
create unique index if not exists product_feature_product_id_feature_id_index on product_feature (product_id, feature_id);
create table if not exists monitor_status_change (
    id bigserial,
    status integer default 1,
    changed_at timestamp(0),
    monitor_id bigint,
    inserted_at timestamp(0) not null default now(),
    updated_at timestamp(0) not null default now(),
    primary key (id),
    foreign key (monitor_id) references monitor on delete cascade
);
create index if not exists monitor_status_change_monitor_id_index on monitor_status_change (monitor_id);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop table monitor_region;
drop table error_log;
drop table assertion;
drop table monitor_notification_policy;
drop table monitor_alarm_policy;
drop table daily_report;
drop table monitor_notification;
drop table monitor_integration;
drop table alarm;
drop table "check";
drop table region;
drop table user_contact;
drop table guest_user;
drop table invitation;
drop table organization_user;
drop table receipt;
drop table subscription;
drop table plan;
drop table role_claim;
drop table product_feature;
drop table product;
drop table feature;
drop table monitor_status_change;
drop table monitor;
drop table "user";
drop table role;
drop table monitor_group;
drop table organization;
-- +goose StatementEnd