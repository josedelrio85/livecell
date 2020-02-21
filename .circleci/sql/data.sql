USE `test_db`;

CREATE TABLE `leadlive` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `sou_id` bigint(20) DEFAULT NULL,
  `type_id` bigint(20) DEFAULT NULL,
  `queue_id` bigint(20) DEFAULT NULL,
  `smartcenter_id` bigint(20) DEFAULT NULL,
  `cat_id` bigint(20) DEFAULT NULL,
  `subcat_id` bigint(20) DEFAULT NULL,
  `ord_id` bigint(20) DEFAULT NULL,
  `wsid` bigint(20) DEFAULT NULL,
  `closed` bigint(20) DEFAULT NULL,
  `phone` varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  `url` text COLLATE utf8_spanish_ci,
  `is_client` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_leadlive_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COLLATE=utf8_spanish_ci;


CREATE SCHEMA `crmti`;
USE `crmti`;

CREATE TABLE `que_queues` (
  `que_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `que_description` varchar(255) COLLATE utf8_spanish_ci NOT NULL DEFAULT '',
  `que_group` int(11) unsigned NOT NULL DEFAULT '1',
  `que_group_tel` int(11) unsigned NOT NULL DEFAULT '0',
  `que_type` int(11) unsigned NOT NULL DEFAULT '1',
  `que_source` int(11) unsigned NOT NULL DEFAULT '1',
  `que_lot` int(11) unsigned NOT NULL DEFAULT '1',
  `que_weight` int(2) unsigned NOT NULL DEFAULT '50' COMMENT 'Peso, a mayor numero mas peso',
  `que_stack` int(1) unsigned NOT NULL DEFAULT '0' COMMENT 'Tipo de pila 0 - fifo 1 lifo',
  `que_form` int(11) unsigned NOT NULL DEFAULT '1',
  `que_form_cli` int(11) unsigned NOT NULL DEFAULT '1',
  `que_form_ord` int(11) NOT NULL  DEFAULT '0',
  `que_active` int(1) unsigned NOT NULL DEFAULT '0',
  `que_cost` float(6,4) unsigned NOT NULL DEFAULT '0.0000',
  `que_dnis` int(11) unsigned NOT NULL DEFAULT '0',
  `que_speed` int(3) NOT NULL DEFAULT '0',
  `que_turnlimit` int(3) NOT NULL DEFAULT '25',
  `que_autodial` int(1) NOT NULL DEFAULT '1',
  `que_new_allow` int(1) NOT NULL DEFAULT '1',
  `que_pause_tel` int(1) NOT NULL DEFAULT '0',
  `que_lost` int(11) NOT NULL DEFAULT '0',
  `que_engine` int(1) NOT NULL DEFAULT '0',
  `que_autoclose` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`que_id`),
  UNIQUE KEY `Segmentacion Duplicada` (`que_type`,`que_source`) USING BTREE,
  KEY `que_group` (`que_group`) USING BTREE,
  KEY `que_type` (`que_type`) USING BTREE,
  KEY `que_source` (`que_source`) USING BTREE,
  KEY `que_lot` (`que_lot`) USING BTREE,
  KEY `que_form` (`que_form`) USING BTREE,
  KEY `que_form_cli` (`que_form_cli`) USING BTREE,
  KEY `que_dnis` (`que_dnis`) USING BTREE,
  KEY `que_weight` (`que_weight`) USING BTREE,
  KEY `que_active` (`que_active`) USING BTREE,
  KEY `que_engine` (`que_engine`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COLLATE=utf8_spanish_ci;

insert into crmti.que_queues (que_source, que_type, que_id) values (73,2,244);

CREATE TABLE `sub_subcategories` (
  `sub_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `sub_cat` int(11) unsigned DEFAULT NULL,
  `sub_description` varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  `sub_action` varchar(1024) COLLATE utf8_spanish_ci NOT NULL DEFAULT '{"result":"2-cierre","cierreTipo":"2-negativo"}',
  `sub_system` enum('1','0') COLLATE utf8_spanish_ci NOT NULL DEFAULT '0',
  `sub_active` int(1) NOT NULL DEFAULT '0',
  `sub_closing` int(2) NOT NULL DEFAULT '0',
  `sub_util` int(1) DEFAULT '0',
  `sub_sch_auto` int(1) DEFAULT NULL,
  `sub_aux` varchar(10) COLLATE utf8_spanish_ci DEFAULT NULL,
  `sub_callback` varchar(1024) COLLATE utf8_spanish_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`sub_id`),
  UNIQUE KEY `sub_id` (`sub_id`) USING BTREE,
  KEY `sub_cat` (`sub_cat`) USING BTREE,
  KEY `sub_description` (`sub_description`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COLLATE=utf8_spanish_ci;

insert into crmti.sub_subcategories (sub_id, sub_action) values (341, '{"result":"2-cierre","cierreTipo":"1-positivo"}');