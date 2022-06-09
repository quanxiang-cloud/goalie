-- 操作日志 ---
BEGIN;
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("5","","audit","操作日志","",0,0);
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("501","5","audit/read","查看","",0,0);
COMMIT;

-- 赋予超级管理员权限 --
BEGIN;
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("42","1","5",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("43","1","501",0,0);
COMMIT;

-- 赋予管理员权限 --
BEGIN;
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("44","2","5",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("45","2","501",0,0);
COMMIT;

-- 赋予观察这权限 --
BEGIN;
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("46","3","5",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("47","3","501",0,0);
COMMIT;

-- 数据集 ---
BEGIN;
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("6","","system","数据集","",0,0);
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("601","6","dataset/read","查看","",0,0);
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("602","6","dataset/create","新建","",0,0);
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("603","6","dataset/update","修改","",0,0);
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("604","6","dataset/delete","删除","",0,0);
COMMIT;

-- 赋予超级管理员权限 --
BEGIN;
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("48","1","6",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("49","1","601",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("50","1","602",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("51","1","603",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("52","1","604",0,0);
COMMIT;

-- 赋予管理员权限 --
BEGIN;
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("53","2","6",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("54","2","601",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("55","2","602",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("56","2","603",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("57","2","604",0,0);
COMMIT;

-- 赋予观察这权限 --
BEGIN;
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("58","3","6",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("59","3","601",0,0);
COMMIT;


-- 异常流程 ---
BEGIN;
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("7","","abnormalFlow","异常流程","",0,0);
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("701","7","abnormalFlow/read","查看","",0,0);
INSERT INTO `func` (`id`, `p_id`,`func_tag`,`name`,`describe`,`created_at`,`updated_at`)
    VALUES ("702","7","abnormalFlow/dispose","处理","",0,0);
COMMIT;

-- 赋予超级管理员权限 --
BEGIN;
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("60","1","7",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("61","1","701",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("62","1","702",0,0);
COMMIT;

-- 赋予管理员权限 --
BEGIN;
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("63","2","7",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("64","2","701",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("65","2","702",0,0);
COMMIT;

-- 赋予观察这权限 --
BEGIN;
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("66","3","7",0,0);
INSERT INTO `role_func` (`id`, `role_id`,`func_id`,`created_at`,`updated_at`)
    VALUES ("67","3","701",0,0);
COMMIT;