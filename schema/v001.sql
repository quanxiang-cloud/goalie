CREATE TABLE `role`(
    `id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(30) NOT NULL,
    `tag` VARCHAR(24) NOT NULL,
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    created_by varchar(64) null,
    updated_by varchar(64) null,
    tenant_id  varchar(64) null,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_role_name` (`name`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- init role -----
BEGIN;
INSERT INTO `role` (`id`, `name`,`tag`,`created_at`,`updated_at`)
    VALUES ("1","超级管理员","super",0,0);
INSERT INTO `role` (`id`, `name`,`tag`,`created_at`,`updated_at`)
    VALUES ("2","管理员","warden",0,0);
INSERT INTO `role` (`id`, `name`,`tag`,`created_at`,`updated_at`)
    VALUES ("3","查看者","viewer",0,0);
COMMIT;


CREATE TABLE `role_owner`(
    `id` VARCHAR(36) NOT NULL,
    `role_id` VARCHAR(36) NOT NULL,
    `type` TINYINT(1) NOT NULL,
    `owner_id` VARCHAR(36) NOT NULL,
    `created_at` BIGINT(20) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_role_owner`(`role_id`,`owner_id`),
    INDEX `nk_role`(`role_id`),
    INDEX `nk_owner`(`owner_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `func`(
    `id` VARCHAR(36) NOT NULL,
    `p_id` VARCHAR(36),
    `func_tag`  VARCHAR(255) NOT NULL,
    `name` VARCHAR(60) NOT NULL,
    `describe` TEXT,
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `nk_func_pid` (`p_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

BEGIN;
-- 应用管理 --
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("1","","application","应用管理","",0,0);
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("101","1","application/read","查看应用","",0,0);
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("102","1","application/create","新建应用","",0,0);
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("103","1","application/update","修改应用","",0,0);
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("104","1","application/delete","删除应用","",0,0);
-- 访问控制 --
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("2","","accessControl","访问控制","",0,0);
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("201","2","accessControl/mailList/read","查看企业通讯录","",0,0);
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("202","2","accessControl/mailList/manage","管理企业通讯录","",0,0);
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("203","2","accessControl/role/read","查看角色","",0,0);
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("204","2","accessControl/role/manage","管理角色","",0,0);
-- 平台设置 --
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("3","","platform","平台设置","",0,0);
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("301","3","platform/read","查看平台设置","",0,0);
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("302","3","platform/mangage","管理平台设置","",0,0);
COMMIT;


CREATE TABLE `role_func`(
    `id` VARCHAR(36) NOT NULL,
    `role_id` VARCHAR(36)  NOT NULL,
    `func_id`  VARCHAR(36) NOT NULL,
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_role_func` (`role_id`,`func_id`),
    INDEX `nk_role` (`role_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 赋予超级管理员全部权限 --
BEGIN;
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("1","1","1",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("2","1","101",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("3","1","102",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("4","1","103",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("5","1","104",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("6","1","2",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("7","1","201",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("8","1","202",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("9","1","203",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("10","1","204",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("11","1","3",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("12","1","301",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("13","1","302",0,0);
COMMIT;

-- 赋予管理员全部权限 --
BEGIN;
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("14","2","1",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("15","2","101",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("16","2","102",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("17","2","103",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("18","2","104",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("19","2","2",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("20","2","201",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("21","2","202",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("22","2","203",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("23","2","204",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("24","2","3",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("25","2","301",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("26","2","302",0,0);
COMMIT;

-- 赋予观察这权限 --
BEGIN;
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("27","3","1",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("28","3","101",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("29","3","2",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("30","3","201",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("31","3","203",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("32","3","3",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("33","3","301",0,0);
COMMIT;
