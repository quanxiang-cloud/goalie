-- 系统管理 --
BEGIN;
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("4","","system","系统设置","",0,0);
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("401","4","system/read","消息查看","",0,0);
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("402","4","system/mangage","消息管理","",0,0);
COMMIT;


-- 赋予超级管理员权限 --
BEGIN;
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("34","1","4",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("35","1","401",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("36","1","402",0,0);
COMMIT;

-- 赋予管理员权限 --
BEGIN;
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("37","2","4",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("38","2","401",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("39","2","402",0,0);
COMMIT;

-- 赋予观察这权限 --
BEGIN;
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("40","3","4",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("41","3","401",0,0);
COMMIT;