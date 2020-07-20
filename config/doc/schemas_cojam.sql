
-- 사용자 셋팅 : grit
-- create user 'grit'@'%' identified by 'grit';
-- grant all privileges on *.* to 'grit'@'%';

-- START 모든 테이블 삭제
SET @tables = NULL;
SELECT GROUP_CONCAT(table_schema, '.', table_name) INTO @tables
  FROM information_schema.tables 
  WHERE table_schema = 'cojam'; -- specify DB name here.

SET @tables = CONCAT('DROP TABLE ', @tables);
PREPARE stmt FROM @tables;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
-- END 모든 테이블 삭제




CREATE TABLE `user` (
  `user_seq` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` varchar(40) NOT NULL,
  `user_name` varchar(20) NOT NULL,
  `nick_name` varchar(20) DEFAULT NULL,
  `profile_img_url` varchar(300) DEFAULT NULL,
  `use_yn` char(1) NOT NULL DEFAULT 'N',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`user_seq`),
  UNIQUE KEY `user_id_UNIQUE` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `cojam`.`user_profile_img` (
  `user_seq` INT UNSIGNED NOT NULL,
  `user_profile_img_seq` VARCHAR(36) NULL,
  `domain` VARCHAR(100) NOT NULL,
  `path` VARCHAR(300) NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  PRIMARY KEY (`user_seq`),
  UNIQUE INDEX `user_profile_img_seq_UNIQUE` (`user_profile_img_seq` ASC),
  CONSTRAINT `fk_user_profile_img_seq_user1`
    FOREIGN KEY (`user_seq`)
    REFERENCES `cojam`.`user` (`user_seq`)
    ON DELETE CASCADE
    ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `user_auth` (
  `user_auth_seq` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_seq` int(10) unsigned NOT NULL,
  `auth_code` char(2) NOT NULL,
  `auth_name` varchar(10) NOT NULL,
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`user_auth_seq`,`user_seq`),
  KEY `fk_user_auth_user1_idx` (`user_seq`),
  CONSTRAINT `fk_user_auth_user1` FOREIGN KEY (`user_seq`) REFERENCES `user` (`user_seq`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='auth_name : auth_code = [admin: 01, core: 02, jammer: 03]';

CREATE TABLE `user_join_info` (
  `provider` varchar(20) NOT NULL,
  `oauth_code` varchar(40) NOT NULL,
  `user_seq` int(10) unsigned NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`provider`,`oauth_code`),
  KEY `fk_user_join_info_user` (`user_seq`),
  CONSTRAINT `fk_user_join_info_user` FOREIGN KEY (`user_seq`) REFERENCES `user` (`user_seq`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `cojam`.`user_notice_agree` (
  `user_seq` INT UNSIGNED NOT NULL,
  `event_yn` CHAR(1) NOT NULL DEFAULT 'N',
  `subscribe_yn` CHAR(1) NOT NULL DEFAULT 'N',
  `channel_yn` CHAR(1) NOT NULL DEFAULT 'N',
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  PRIMARY KEY (`user_seq`),
  CONSTRAINT `fk_user_notice_agree_user1`
    FOREIGN KEY (`user_seq`)
    REFERENCES `cojam`.`user` (`user_seq`)
    ON DELETE CASCADE
    ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `cojam`.`user_blacklist_detail` (
  `user_blacklist_detail_seq` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_seq` INT UNSIGNED NOT NULL,
  `status` VARCHAR(2) NOT NULL,
  `description` VARCHAR(500) NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  PRIMARY KEY (`user_blacklist_detail_seq`, `user_seq`),
  INDEX `fk_user_blacklist_detail_user1_idx` (`user_seq` ASC),
  CONSTRAINT `fk_user_blacklist_detail_user1`
    FOREIGN KEY (`user_seq`)
    REFERENCES `cojam`.`user` (`user_seq`)
    ON DELETE CASCADE
    ON UPDATE CASCADE)
ENGINE = InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `channel` (
  `channel_seq` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_seq` int(10) unsigned NOT NULL,
  `title` varchar(100) NOT NULL,
  `description` varchar(500) DEFAULT NULL,
  `subscription_cnt` int(10) unsigned DEFAULT 0,
  `live_join_yn` char(1) NOT NULL DEFAULT 'Y',
  `vod_join_yn` char(1) NOT NULL DEFAULT 'Y',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`channel_seq`,`user_seq`),
  KEY `fk_channel_user1_idx` (`user_seq`),
  CONSTRAINT `fk_channel_user1` FOREIGN KEY (`user_seq`) REFERENCES `user` (`user_seq`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- CREATE TABLE `channel_desc` (
--   `channel_seq` int(10) unsigned NOT NULL,
--   `description` varchar(500) DEFAULT NULL,
--   `created_at` datetime NOT NULL,
--   `updated_at` datetime NOT NULL,
--   PRIMARY KEY (`channel_seq`),
--   CONSTRAINT `fk_channel_desc_channel1` FOREIGN KEY (`channel_seq`) REFERENCES `channel` (`channel_seq`) ON DELETE CASCADE ON UPDATE CASCADE
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `channel_subscription` (
  `channel_seq` int(10) unsigned NOT NULL,
  `user_seq` int(10) unsigned NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`channel_seq`,`user_seq`),
  INDEX `fk_channel_subscription_user1_idx` (`user_seq` ASC),
  CONSTRAINT `fk_channel_subscription_channel1` FOREIGN KEY (`channel_seq`) REFERENCES `channel` (`channel_seq`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_channel_subscription_user1`
    FOREIGN KEY (`user_seq`)
    REFERENCES `cojam`.`user` (`user_seq`)
    ON DELETE CASCADE
    ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `channel_thumbnail` (
  `channel_seq` int(10) unsigned NOT NULL,
  `domain` varchar(100) DEFAULT NULL,
  `path` varchar(300) DEFAULT NULL,
  `file_name` varchar(100) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`channel_seq`),
  CONSTRAINT `fk_channel_thumbnail_channel1` FOREIGN KEY (`channel_seq`) REFERENCES `channel` (`channel_seq`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



CREATE TABLE `common_code_group` (
  `group_seq` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `group_name` varchar(20) NOT NULL,
  `category` varchar(20) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`group_seq`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `common_code` (
  `group_seq` int(10) unsigned NOT NULL,
  `code` varchar(4) NOT NULL,
  `code_name` varchar(50) NOT NULL,
  `description` varchar(100) DEFAULT '',
  `order_num` tinyint(4) DEFAULT NULL,
  `use_yn` char(1) NOT NULL DEFAULT 'Y',
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  `updated_at` datetime NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`group_seq`,`code`),
  CONSTRAINT `fk_common_code_common_code_group1` FOREIGN KEY (`group_seq`) REFERENCES `common_code_group` (`group_seq`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `faq` (
  `faq_seq` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_seq` int(10) unsigned NOT NULL,
  `faq_code` varchar(4) NOT NULL,
  `status` varchar(4) NOT NULL,
  `question` mediumtext NOT NULL,
  `answer` mediumtext DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`faq_seq`,`user_seq`),
  KEY `fk_faq_user1_idx` (`user_seq`),
  CONSTRAINT `fk_faq_user1` FOREIGN KEY (`user_seq`) REFERENCES `user` (`user_seq`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `live` (
  `live_seq` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_seq` int(10) unsigned NOT NULL,
  `room_token` text DEFAULT NULL,
  `status` char(4) NOT NULL,
  `start_dt` datetime DEFAULT NULL,
  `end_dt` datetime DEFAULT NULL,
  `title` varchar(30) DEFAULT NULL,
  `description` varchar(500) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`live_seq`,`user_seq`),
  KEY `fk_live_user1_idx` (`user_seq`),
  CONSTRAINT `fk_live_user1` FOREIGN KEY (`user_seq`) REFERENCES `user` (`user_seq`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `cojam`.`live_join_history` (
  `live_seq` INT UNSIGNED NOT NULL,
  `user_seq` INT UNSIGNED NOT NULL,
  `created_at` DATETIME NOT NULL COMMENT '참여일',
  `deleted_at` DATETIME NOT NULL COMMENT '삭제일',
  PRIMARY KEY (`live_seq`, `user_seq`)
) ENGINE = InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `live_cnt_info` (
  `live_seq` int(10) unsigned NOT NULL,
  `live_cnt` int(10) unsigned DEFAULT NULL,
  `view_cnt` int(10) unsigned DEFAULT NULL,
  `like_cnt` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`live_seq`),
  CONSTRAINT `fk_live_cnt_info_live1` FOREIGN KEY (`live_seq`) REFERENCES `live` (`live_seq`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `live_like` (
  `live_seq` int(10) unsigned NOT NULL,
  `user_seq` int(10) unsigned NOT NULL,
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`live_seq`,`user_seq`),
  CONSTRAINT `fk_live_like_live1` FOREIGN KEY (`live_seq`) REFERENCES `live` (`live_seq`) ON DELETE CASCADE ON UPDATE CASCADE,
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `live_thumbnail` (
  `live_thumbnail_seq` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `live_seq` int(10) unsigned NOT NULL,
  `order_num` int(10) unsigned NOT NULL,
  `domain` varchar(100) NOT NULL,
  `path` varchar(300) NOT NULL,
  `file_name` varchar(100) NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`live_thumbnail_seq`),
  KEY `fk_live_thumbnail_live1_idx` (`live_seq`),
  CONSTRAINT `fk_live_thumbnail_live1` FOREIGN KEY (`live_seq`) REFERENCES `live` (`live_seq`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



CREATE TABLE `vod` (
  `vod_seq` varchar(36) NOT NULL,
  `channel_seq` int(10) unsigned NOT NULL,
  `live_seq` int(10) unsigned NOT NULL,
  `user_seq` int(10) unsigned NOT NULL,
  `title` varchar(30) DEFAULT NULL,
  `description` varchar(500) DEFAULT NULL,
  `status` varchar(4) NOT NULL DEFAULT '1001',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`vod_seq`),
  KEY `fk_vod_user1_idx` (`user_seq`),
  KEY `fk_vod_live1_idx` (`live_seq`),
  KEY `fk_vod_channel1_idx` (`channel_seq`),
  CONSTRAINT `fk_vod_join_user1` FOREIGN KEY (`user_seq`) REFERENCES `user` (`user_seq`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `cojam`.`vod_like` (
  `vod_seq` VARCHAR(36) NOT NULL,
  `user_seq` INT UNSIGNED NOT NULL,
  `created_at` DATETIME NULL DEFAULT now(),
  PRIMARY KEY (`vod_seq`, `user_seq`),
  INDEX `fk_vod_like_vod1_idx` (`vod_seq` ASC),
  CONSTRAINT `fk_vod_like_vod1`
    FOREIGN KEY (`vod_seq`)
    REFERENCES `cojam`.`vod` (`vod_seq`)
    ON DELETE CASCADE
    ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `vod_cnt_info` (
  `vod_seq` varchar(36) NOT NULL,
  `view_cnt` int(10) unsigned NOT NULL DEFAULT 0,
  `like_cnt` int(10) unsigned NOT NULL DEFAULT 0,
  `join_cnt` int(10) unsigned NOT NULL DEFAULT 0,
  `comment_cnt` int(10) unsigned NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`vod_seq`),
  CONSTRAINT `fk_vod_cnt_info_vod1` FOREIGN KEY (`vod_seq`) REFERENCES `vod` (`vod_seq`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `vod_join` (
  `channel_seq` int(10) unsigned NOT NULL,
  `vod_seq` varchar(36) NOT NULL,
  `user_seq` int(10) unsigned NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`channel_seq`,`vod_seq`),
  KEY `fk_vod_join_vod1_idx` (`vod_seq`),
  KEY `fk_vod_join_user2_idx` (`user_seq`),
  CONSTRAINT `fk_vod_join_vod1` FOREIGN KEY (`vod_seq`) REFERENCES `vod` (`vod_seq`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `cojam`.`vod_join_history` (
  `channel_seq` INT UNSIGNED NOT NULL,
  `vod_seq` VARCHAR(36) NOT NULL,
  `user_seq` INT NOT NULL,
  `created_at` DATETIME NOT NULL COMMENT '참여일',
  `updated_at` DATETIME NOT NULL COMMENT '참여 업데이트일',
  `delete_at` DATETIME NOT NULL COMMENT '삭제일'
) ENGINE = InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `vod_path` (
  `vod_seq` varchar(36) NOT NULL,
  `vod_type` varchar(4) NOT NULL,
  `domain` varchar(100) NOT NULL,
  `path` varchar(300) NOT NULL,
  `file_name` varchar(100) NOT NULL,
  `duration` int(10) unsigned NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`vod_seq`,`vod_type`),
  CONSTRAINT `fk_vod_path_vod1` FOREIGN KEY (`vod_seq`) REFERENCES `vod` (`vod_seq`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `vod_thumbnail` (
  `vod_thumbnail_seq` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `vod_seq` varchar(36) NOT NULL,
  `order_num` int(11) NOT NULL,
  `domain` varchar(100) NOT NULL,
  `path` varchar(300) NOT NULL,
  `file_name` varchar(100) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`vod_thumbnail_seq`),
  KEY `fk_vod_thumbnail_vod1_idx` (`vod_seq`),
  CONSTRAINT `fk_vod_thumbnail_vod1` FOREIGN KEY (`vod_seq`) REFERENCES `vod` (`vod_seq`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


CREATE TABLE IF NOT EXISTS `cojam`.`withdrawal_user` (
  `user_seq` INT UNSIGNED NOT NULL,
  `user_id` varchar(40) NOT NULL,
  `description` VARCHAR(500) NULL,
  `created_at` DATETIME NOT NULL,
  PRIMARY KEY (`user_seq`),
  CONSTRAINT `fk_withdrawal_user_user1`
    FOREIGN KEY (`user_seq`)
    REFERENCES `cojam`.`user` (`user_seq`)
    ON DELETE CASCADE
    ON UPDATE CASCADE)
ENGINE = InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE IF NOT EXISTS `cojam`.`vod_comment` (
  `vod_comment_seq` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `vod_seq` VARCHAR(36) NOT NULL,
  `user_seq` INT UNSIGNED NOT NULL,
  `comment` VARCHAR(300) NOT NULL,
  `like_cnt` INT UNSIGNED NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  PRIMARY KEY (`vod_comment_seq`),
  INDEX `fk_vod_comment_vod1_idx` (`vod_seq` ASC),
  INDEX `fk_vod_comment_user1_idx` (`user_seq` ASC),
  CONSTRAINT `fk_vod_comment_vod1`
    FOREIGN KEY (`vod_seq`)
    REFERENCES `cojam`.`vod` (`vod_seq`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `fk_vod_comment_user1`
    FOREIGN KEY (`user_seq`)
    REFERENCES `cojam`.`user` (`user_seq`)
    ON DELETE CASCADE
    ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `cojam`.`report_illegality` (
  `report_illegality_seq` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `contents_type` VARCHAR(20) NOT NULL COMMENT '생방송, 동영상, 채팅, 댓글',
  `contents_seq` INT UNSIGNED NOT NULL COMMENT 'live_seq, vod_seq ...',
  `reporter` INT UNSIGNED NOT NULL COMMENT '신고자(user_seq)',
  `status` VARCHAR(4) NOT NULL,
  `report_code` VARCHAR(4) NOT NULL,
  `report_detail` VARCHAR(500) NULL,
  `created_at` DATETIME NOT NULL DEFAULT now(),
  `updated_at` DATETIME NOT NULL DEFAULT now(),
  PRIMARY KEY (`report_illegality_seq`)
) ENGINE = InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;



