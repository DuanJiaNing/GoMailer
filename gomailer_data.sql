-- --------------------------------------------------------
-- 主机:                           127.0.0.1
-- 服务器版本:                        8.0.19 - MySQL Community Server - GPL
-- 服务器操作系统:                      Win64
-- HeidiSQL 版本:                  10.3.0.5771
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;


-- 导出 gomailer 的数据库结构
CREATE DATABASE IF NOT EXISTS `gomailer` /*!40100 DEFAULT CHARACTER SET utf8 */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `gomailer`;

-- 导出  表 gomailer.dialer 结构
CREATE TABLE IF NOT EXISTS `dialer` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `insert_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `host` varchar(1000) NOT NULL,
  `port` int NOT NULL,
  `user_id` int unsigned NOT NULL,
  `auth_username` varchar(1000) NOT NULL,
  `auth_password` varchar(1000) NOT NULL,
  `name` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`user_id`,`name`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

-- 正在导出表  gomailer.dialer 的数据：~3 rows (大约)
/*!40000 ALTER TABLE `dialer` DISABLE KEYS */;
INSERT INTO `dialer` (`id`, `insert_time`, `host`, `port`, `user_id`, `auth_username`, `auth_password`, `name`) VALUES
	(1, '2020-03-12 12:03:48', 'smtp.qq.com', 465, 3, '666@qq.com', '666aaa111', 'XX公司'),
	(2, '2020-03-12 13:05:21', 'smtp.qq.com', 465, 1, '666@qq.com', '666aaa', 'XX公司'),
	(3, '2020-03-18 17:02:38', 'smtp.qq.com', 465, 1, '2213994603@qq.com', 'athupcbmeyvvdjif', 'DuanJiaNing公司');
/*!40000 ALTER TABLE `dialer` ENABLE KEYS */;

-- 导出  表 gomailer.endpoint 结构
CREATE TABLE IF NOT EXISTS `endpoint` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `insert_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `user_app_id` int unsigned NOT NULL,
  `dialer_id` int unsigned NOT NULL,
  `template_id` int unsigned DEFAULT NULL,
  `user_id` int unsigned NOT NULL,
  `name` varchar(500) NOT NULL,
  `key` varchar(1000) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_app_id` (`user_app_id`,`name`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

-- 正在导出表  gomailer.endpoint 的数据：~2 rows (大约)
/*!40000 ALTER TABLE `endpoint` DISABLE KEYS */;
INSERT INTO `endpoint` (`id`, `insert_time`, `user_app_id`, `dialer_id`, `template_id`, `user_id`, `name`, `key`) VALUES
	(1, '2020-03-12 13:36:26', 1, 3, 38, 1, '发送反馈', 'Wiov3aJpmu'),
	(2, '2020-03-18 17:02:38', 1, 3, 39, 1, '用户反馈', 'P9jwNrzCpz');
/*!40000 ALTER TABLE `endpoint` ENABLE KEYS */;

-- 导出  表 gomailer.endpoint_preference 结构
CREATE TABLE IF NOT EXISTS `endpoint_preference` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `insert_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `endpoint_id` int unsigned NOT NULL,
  `deliver_strategy` varchar(500) DEFAULT NULL,
  `enable_re_captcha` tinyint DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `end_point_id` (`endpoint_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;

-- 正在导出表  gomailer.endpoint_preference 的数据：~1 rows (大约)
/*!40000 ALTER TABLE `endpoint_preference` DISABLE KEYS */;
INSERT INTO `endpoint_preference` (`id`, `insert_time`, `endpoint_id`, `deliver_strategy`, `enable_re_captcha`) VALUES
	(4, '2020-03-12 14:08:41', 1, 'DELIVER_IMMEDIATELY', 2),
	(5, '2020-03-18 17:02:38', 2, 'DELIVER_IMMEDIATELY', 2);
/*!40000 ALTER TABLE `endpoint_preference` ENABLE KEYS */;

-- 导出  表 gomailer.mail 结构
CREATE TABLE IF NOT EXISTS `mail` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `insert_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `endpoint_id` int NOT NULL,
  `state` varchar(100) NOT NULL,
  `delivery_time` timestamp NULL DEFAULT NULL,
  `content` longtext NOT NULL,
  `raw` longtext NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 正在导出表  gomailer.mail 的数据：~0 rows (大约)
/*!40000 ALTER TABLE `mail` DISABLE KEYS */;
/*!40000 ALTER TABLE `mail` ENABLE KEYS */;

-- 导出  表 gomailer.receiver 结构
CREATE TABLE IF NOT EXISTS `receiver` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `insert_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `endpoint_id` int unsigned NOT NULL,
  `user_id` int unsigned NOT NULL,
  `user_app_id` int unsigned NOT NULL,
  `address` varchar(1000) NOT NULL,
  `receiver_type` varchar(500) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=69 DEFAULT CHARSET=utf8;

-- 正在导出表  gomailer.receiver 的数据：~6 rows (大约)
/*!40000 ALTER TABLE `receiver` DISABLE KEYS */;
INSERT INTO `receiver` (`id`, `insert_time`, `endpoint_id`, `user_id`, `user_app_id`, `address`, `receiver_type`) VALUES
	(63, '2020-03-19 11:40:46', 1, 1, 1, 'TO_djn<duan_jia_ning@163.com>', 'To'),
	(64, '2020-03-19 11:40:46', 1, 1, 1, 'CC_djn<duanjianing0@gmail.com>', 'Cc'),
	(65, '2020-03-19 11:40:46', 1, 1, 1, 'BCC_djn<aimeimeits@gmail.com>', 'Bcc'),
	(66, '2020-03-19 11:41:07', 2, 1, 1, 'TO_djn<duan_jia_ning@163.com>', 'To'),
	(67, '2020-03-19 11:41:07', 2, 1, 1, 'CC_djn<duanjianing0@gmail.com>', 'Cc'),
	(68, '2020-03-19 11:41:07', 2, 1, 1, 'BCC_djn<aimeimeits@gmail.com>', 'Bcc');
/*!40000 ALTER TABLE `receiver` ENABLE KEYS */;

-- 导出  表 gomailer.template 结构
CREATE TABLE IF NOT EXISTS `template` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `insert_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `user_id` int unsigned NOT NULL,
  `template` longtext NOT NULL,
  `content_type` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=40 DEFAULT CHARSET=utf8;

-- 正在导出表  gomailer.template 的数据：~39 rows (大约)
/*!40000 ALTER TABLE `template` DISABLE KEYS */;
INSERT INTO `template` (`id`, `insert_time`, `user_id`, `template`, `content_type`) VALUES
	(1, '2020-03-12 13:06:23', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(2, '2020-03-12 13:30:09', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(3, '2020-03-12 13:34:01', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(4, '2020-03-12 13:37:31', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(5, '2020-03-12 13:49:18', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(6, '2020-03-12 13:50:38', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(7, '2020-03-12 13:51:54', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(8, '2020-03-12 13:54:01', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(9, '2020-03-12 14:06:06', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(10, '2020-03-12 14:06:34', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(11, '2020-03-12 14:07:45', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(12, '2020-03-12 14:08:41', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(13, '2020-03-12 14:09:02', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(14, '2020-03-12 16:19:09', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(15, '2020-03-12 16:19:17', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(16, '2020-03-12 16:19:19', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(17, '2020-03-12 16:19:20', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(18, '2020-03-12 16:19:20', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(19, '2020-03-12 16:19:21', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(20, '2020-03-12 16:19:22', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(21, '2020-03-12 16:19:23', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(22, '2020-03-12 16:19:24', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(23, '2020-03-12 16:19:24', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(24, '2020-03-12 16:20:38', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(25, '2020-03-12 16:20:40', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(26, '2020-03-12 16:20:44', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(27, '2020-03-12 16:21:04', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(28, '2020-03-12 16:25:06', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(29, '2020-03-12 16:25:09', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(30, '2020-03-12 16:39:31', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(31, '2020-03-12 16:41:08', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(32, '2020-03-12 17:07:14', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(33, '2020-03-12 17:08:00', 1, '<div><hr><h1>Test email{{msg}}</h1><div/>', 'text/html'),
	(34, '2020-03-18 17:02:38', 1, '<div>来自用户[{{name}}]的反馈, 用户电话号码: {{phone}}, 反馈内容如下:<hr><p>{{content}}</p><div/>', 'text/html'),
	(35, '2020-03-18 17:04:19', 1, '<div>来自用户[{{name}}]的反馈, 用户电话号码: {{phone}}, 反馈内容如下:<hr><p>{{content}}</p><div/>', 'text/html'),
	(36, '2020-03-18 17:58:18', 1, '<div>来自用户[{{name}}]的反馈, 用户电话号码: {{phone}}, 反馈内容如下:<hr><p>{{content}}</p><div/>', 'text/html'),
	(37, '2020-03-19 11:39:00', 1, '<div>来自用户[{{name}}]的反馈, 用户电话号码: {{phone}}, 反馈内容如下:<hr><p>{{content}}</p><div/>', 'text/html'),
	(38, '2020-03-19 11:40:46', 1, '<div>来自用户[{{name}}]的反馈, 用户电话号码: {{phone}}, 反馈内容如下:<hr><p>{{content}}</p><div/>', 'text/html'),
	(39, '2020-03-19 11:41:07', 1, '<div>来自用户[{{name}}]的反馈, 用户电话号码: {{phone}}, 反馈内容如下:<hr><p>{{content}}</p><div/>', 'text/html');
/*!40000 ALTER TABLE `template` ENABLE KEYS */;

-- 导出  表 gomailer.user 结构
CREATE TABLE IF NOT EXISTS `user` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `insert_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `username` varchar(500) NOT NULL,
  `password` varchar(1000) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

-- 正在导出表  gomailer.user 的数据：~3 rows (大约)
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` (`id`, `insert_time`, `username`, `password`) VALUES
	(1, '2020-03-12 11:21:43', 'djn', '123456'),
	(2, '2020-03-12 11:23:26', 'dj1n', '123456'),
	(3, '2020-03-12 11:45:08', 'djn1', '123456');
/*!40000 ALTER TABLE `user` ENABLE KEYS */;

-- 导出  表 gomailer.user_app 结构
CREATE TABLE IF NOT EXISTS `user_app` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `insert_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `user_id` int unsigned NOT NULL,
  `name` varchar(500) NOT NULL,
  `host` varchar(1000) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`user_id`,`name`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

-- 正在导出表  gomailer.user_app 的数据：~2 rows (大约)
/*!40000 ALTER TABLE `user_app` DISABLE KEYS */;
INSERT INTO `user_app` (`id`, `insert_time`, `user_id`, `name`, `host`) VALUES
	(1, '2020-03-12 11:42:40', 1, 'demo', 'demo.com'),
	(2, '2020-03-12 11:45:08', 3, 'demo', 'demo.com'),
	(3, '2020-03-12 12:02:38', 3, 'demo1', 'demo.com');
/*!40000 ALTER TABLE `user_app` ENABLE KEYS */;

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
