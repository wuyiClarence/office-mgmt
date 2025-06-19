CREATE TABLE `user`
(
    `id`                  BIGINT       NOT NULL AUTO_INCREMENT,
    `user_name`           VARCHAR(255) NOT NULL DEFAULT '',
    `user_display_name`   VARCHAR(255) NOT NULL DEFAULT '',
    `status`              TINYINT      NOT NULL DEFAULT 0 COMMENT '1:正常状态 2:已删除',
    `email`               VARCHAR(255) NOT NULL DEFAULT '',
    `phone_number`        VARCHAR(255) NOT NULL DEFAULT '',
    `password`            VARCHAR(255) NOT NULL DEFAULT '',
    `permissions`         BIGINT       NOT NULL DEFAULT 0 COMMENT '',
    `password_updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `deleted_at`          DATETIME              DEFAULT NULL,
    `logout_at`           DATETIME              DEFAULT NULL,
    `create_user`         BIGINT       NOT NULL DEFAULT 0,
    `created_at`          DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`          DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_user_name` (`user_name`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT '用户信息表';

CREATE TABLE `role`
(
    `id`            BIGINT       NOT NULL AUTO_INCREMENT,
    `role_name`     VARCHAR(128) NOT NULL DEFAULT '',
    `description`   VARCHAR(255) NOT NULL DEFAULT '',
    `create_user`   BIGINT       NOT NULL DEFAULT 0,
    `system_create` TINYINT      NOT NULL DEFAULT 0 COMMENT '1:系统内置 2:用户创建',
    `created_at`    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_role_name` (`role_name`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT '用户角色表';

CREATE TABLE `device`
(
    `id`          BIGINT       NOT NULL AUTO_INCREMENT,
    `unique_id`   VARCHAR(128) NOT NULL DEFAULT '',
    `device_name` VARCHAR(512) NOT NULL DEFAULT '',
    `alias_name` VARCHAR(512) NOT NULL DEFAULT '',
    `mac`         VARCHAR(128) NOT NULL DEFAULT '',
    `ip`          VARCHAR(50)  NOT NULL DEFAULT '',
    `os_type`     VARCHAR(50)  NOT NULL DEFAULT '',
    `device_type` TINYINT      NOT NULL DEFAULT 0 COMMENT '1:实体服务器 2:虚拟机',
    `status`      TINYINT      NOT NULL DEFAULT 0 COMMENT '1:运行中',
    `op_status`   INT      NOT NULL DEFAULT 0 COMMENT '',
    `wake_device_id`  BIGINT      NOT NULL DEFAULT 0,
    `host_device_id`  BIGINT   NOT NULL DEFAULT 0,
    `online_time`  DATETIME    DEFAULT NULL,
    `created_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_unique_id` (`unique_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT '设备表';

CREATE TABLE `device_group`
(
    `id`          BIGINT       NOT NULL AUTO_INCREMENT,
    `group_name`  VARCHAR(100) NOT NULL DEFAULT '',
    `create_user` BIGINT       NOT NULL DEFAULT 0,
    `created_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_group_name` (`group_name`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT '服务器组表';

CREATE TABLE `device_device_group_rel`
(
    `id`         BIGINT   NOT NULL AUTO_INCREMENT,
    `device_id`  BIGINT   NOT NULL DEFAULT 0,
    `group_id`   BIGINT   NOT NULL DEFAULT 0,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_group_device` (`group_id`, `device_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT '服务器与组关联表';

CREATE TABLE `policy`
(
    `id`                  BIGINT       NOT NULL AUTO_INCREMENT,
    `policy_name`         VARCHAR(100) NOT NULL DEFAULT '',
    `action_type`         VARCHAR(100) NOT NULL DEFAULT '' COMMENT '执行动作',
    `associate_type`      TINYINT      NOT NULL DEFAULT 0 COMMENT '1:单个设备 2:服务器组',
    `status`              TINYINT      NOT NULL DEFAULT 0 COMMENT '0:未启用 1:启用中',
    `execute_type`        INT      NOT NULL DEFAULT 0,
    `start_date`          DATETIME              DEFAULT NULL,
    `end_date`            DATETIME              DEFAULT NULL,
    `execute_time`        VARCHAR(50) NOT NULL DEFAULT '',
    `day_of_week`         INT      NOT NULL DEFAULT 0,
    `created_at`          DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`          DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `create_user`         BIGINT       NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT '策略表';

CREATE TABLE `policy_device_rel`
(
    `id`         BIGINT   NOT NULL AUTO_INCREMENT,
    `policy_id`  BIGINT   NOT NULL DEFAULT 0,
    `device_id`   BIGINT   NOT NULL DEFAULT 0,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_policy_device` (`policy_id`, `device_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT '策略和设备关联表';

CREATE TABLE `policy_device_group_rel`
(
    `id`         BIGINT   NOT NULL AUTO_INCREMENT,
    `policy_id`  BIGINT   NOT NULL DEFAULT 0,
    `device_group_id`   BIGINT   NOT NULL DEFAULT 0,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_policy_device_group` (`policy_id`, `device_group_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT '策略和设备组关联表';

CREATE TABLE `user_role`
(
    `id`         BIGINT   NOT NULL AUTO_INCREMENT,
    `user_id`    BIGINT   NOT NULL DEFAULT 0,
    `role_id`    BIGINT   NOT NULL DEFAULT 0 COMMENT '用户身上绑定的角色id',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_user_role` (`user_id`, `role_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT '用户角色表';



CREATE TABLE `menu_permission`
(
    `id`              BIGINT       NOT NULL AUTO_INCREMENT,
    `permission_name` VARCHAR(128) NOT NULL,
    `permission_key`  VARCHAR(128) NOT NULL,
    `used`            TINYINT      NOT NULL DEFAULT 0,
    `parent_id`       BIGINT       NOT NULL DEFAULT 0,
    `order`           INT          NOT NULL DEFAULT 0,
    `created_at`      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_permission_key` (`permission_key`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT '菜单权限表';


CREATE TABLE `menu_permission_api`
(
    `id`             BIGINT       NOT NULL AUTO_INCREMENT,
    `permission_key` VARCHAR(128) NOT NULL,
    `api_method`     VARCHAR(16)  NOT NULL,
    `api_path`       VARCHAR(256) NOT NULL,
    `used`           TINYINT      NOT NULL DEFAULT 0,
    `created_at`     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_permission_api` (`permission_key`,`api_method`,`api_path`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT '菜单API权限表';

DROP TABLE IF EXISTS `role_menu_permission`;
CREATE TABLE `role_menu_permission`
(
    `id`             BIGINT       NOT NULL AUTO_INCREMENT,
    `role_id`        BIGINT       NOT NULL DEFAULT 0,
    `permission_key` VARCHAR(128) NOT NULL,
    `created_at`     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `create_user`    BIGINT       NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_role_permission` (`role_id`,`permission_key`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT '角色菜单权限表';

DROP TABLE IF EXISTS `role_resource_permission`;
CREATE TABLE `role_resource_permission`
(
    `id`              BIGINT       NOT NULL AUTO_INCREMENT,
    `role_id`         BIGINT       NOT NULL DEFAULT 0,
    `permission_key`  VARCHAR(128) NOT NULL,
    `permission_name` VARCHAR(128) NOT NULL,
    `resource_type`   VARCHAR(128) NOT NULL,
    `resource_id`     BIGINT       NOT NULL,
    `created_at`      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `create_user`     BIGINT       NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_role_key_resource` (`role_id`,`resource_id`,`permission_key`,`resource_type`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT '角色菜单权限表';


DROP TABLE IF EXISTS `user_resource_permission`;
CREATE TABLE `user_resource_permission`
(
    `id`              BIGINT       NOT NULL AUTO_INCREMENT,
    `user_id`         BIGINT       NOT NULL DEFAULT 0,
    `permission_key`  VARCHAR(128) NOT NULL,
    `permission_name` VARCHAR(128) NOT NULL,
    `resource_type`   VARCHAR(128) NOT NULL,
    `resource_id`     BIGINT       NOT NULL,
    `created_at`      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `create_user`     BIGINT       NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_user_key_resource` (`user_id`,`resource_id`,`permission_key`,`resource_type`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT '用户资源权限表';

/*
    初始化数据
 */
INSERT INTO `role` (`role_name`, `description`, `created_at`, `updated_at`, `create_user`, `system_create`)
VALUES ('管理员', '系统内置角色', NOW(), NOW(), 1, 1),
       ('普通用户', '系统内置角色', NOW(), NOW(), 1, 1);

INSERT INTO `user_resource_permission` (`user_id`, `permission_key`, `permission_name`, `resource_type`, `resource_id`,
                                        `created_at`, `updated_at`, `create_user`)
VALUES ('0', 'view', '查看', 'role', 1, NOW(), NOW(), 1),
       ('0', 'view', '查看', 'role', 2, NOW(), NOW(), 1);