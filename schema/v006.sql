-- 默认超管 --

insert into role_owner (`id`, `role_id`, `type`, `owner_id`, `created_at`)values ('1', '1', '1', '1',0);

-- role 新增字段 --
alter table role
    add created_by varchar(64) null;

alter table role
    add updated_by varchar(64) null;

alter table role
    add tenant_id varchar(64) null;
