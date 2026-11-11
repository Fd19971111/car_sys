/*
 Navicat Premium Dump SQL

 Source Server         : 127.0.0.1
 Source Server Type    : MySQL
 Source Server Version : 50726 (5.7.26)
 Source Host           : localhost:3306
 Source Schema         : yiminjia

 Target Server Type    : MySQL
 Target Server Version : 50726 (5.7.26)
 File Encoding         : 65001

 Date: 06/11/2026 16:41:34
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for company
-- ----------------------------
DROP TABLE IF EXISTS `company`;
CREATE TABLE `company`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `code` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `status` tinyint(3) UNSIGNED NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_company_code`(`code`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '租户公司' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of company
-- ----------------------------
INSERT INTO `company` VALUES (1, '亿民家', 'YMJ', 1, '2026-12-05 15:52:30', '2026-11-03 15:52:35');

-- ----------------------------
-- Table structure for customer
-- ----------------------------
DROP TABLE IF EXISTS `customer`;
CREATE TABLE `customer`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `owner_user_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '归属员工ID（site_user.id，0=未分配）',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '客户姓名',
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机号',
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '房屋地址',
  `source` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '客户来源：转介绍/线上/门店等',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态：1=跟进中 2=已成交 3=已流失',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_phone`(`phone`) USING BTREE,
  INDEX `idx_status`(`status`) USING BTREE,
  INDEX `idx_owner`(`company_id`, `owner_user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 16 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '客户表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of customer
-- ----------------------------
INSERT INTO `customer` VALUES (1, 1, 2, '张伟明', '13800138001', '徐汇区锦绣路88弄12号201室', '转介绍', 2, '全屋翻新，朋友介绍，沟通顺畅', '2025-12-20 10:20:00', '2025-12-20 10:20:00', NULL);
INSERT INTO `customer` VALUES (2, 1, 3, '王丽华', '13912345678', '徐汇区漕溪北路518弄6号1801室', '老客户推荐', 2, '新房装修，注重环保材料', '2025-12-28 14:30:00', '2025-12-28 14:30:00', NULL);
INSERT INTO `customer` VALUES (3, 1, 2, '陈志强', '13611112222', '静安区仲益大厦A座706室', '线上推广', 2, '办公室改造，需夜间施工', '2026-07-18 09:10:00', '2026-07-18 09:10:00', NULL);
INSERT INTO `customer` VALUES (4, 1, 5, '黄建军', '13722223333', '杨浦区国定路335号2号楼1102室', '门店到访', 2, '厨卫局部改造，预算 7 万左右', '2026-09-26 16:40:00', '2026-09-26 16:40:00', NULL);
INSERT INTO `customer` VALUES (5, 1, 4, '吴敏', '13544445555', '静安区恒丰路600号22层', '转介绍', 2, '老房翻新，家里有老人，注意防滑', '2026-06-15 11:00:00', '2026-06-15 11:00:00', NULL);
INSERT INTO `customer` VALUES (6, 1, 2, '赵磊', '13866667777', '浦东新区金桥路939号1层商铺', '电话咨询', 2, '商铺装修，工期要求紧', '2026-04-22 15:20:00', '2026-04-22 15:20:00', NULL);
INSERT INTO `customer` VALUES (7, 1, 9, '刘婷婷', '15900001111', '闵行区都市路2899弄5号', '线上推广', 1, 'loft 改造，已量房待报价', '2026-10-12 10:05:00', '2026-10-12 10:05:00', NULL);
INSERT INTO `customer` VALUES (8, 1, 9, '孙浩', '13633334444', '普陀区金沙江路1518弄8号', '门店到访', 1, '婚房装修，计划明年开工', '2026-10-26 13:40:00', '2026-10-26 13:40:00', NULL);
INSERT INTO `customer` VALUES (9, 1, 9, '李静', '13755556666', '长宁区天山路989号', '其他', 3, '预算分歧，暂时搁置', '2026-03-08 09:30:00', '2026-03-08 09:30:00', NULL);

-- ----------------------------
-- Table structure for customer_token
-- ----------------------------
DROP TABLE IF EXISTS `customer_token`;
CREATE TABLE `customer_token`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `customer_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '业主ID，对应 customer.id',
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `token` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '访问令牌',
  `client` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '客户端类型：mp=小程序',
  `ip` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '登录IP',
  `expire_at` datetime NOT NULL COMMENT '过期时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_token`(`token`) USING BTREE,
  INDEX `idx_customer`(`customer_id`) USING BTREE,
  INDEX `idx_expire`(`expire_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '业主登录令牌表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of customer_token
-- ----------------------------

-- ----------------------------
-- Table structure for fin_account
-- ----------------------------
DROP TABLE IF EXISTS `fin_account`;
CREATE TABLE `fin_account`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '科目ID',
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `account_code` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '科目编码',
  `account_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '科目名称',
  `category_code` varchar(24) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '科目类别代码：ASSET/LIABILITY/EQUITY/REVENUE/EXPENSE',
  `direction_code` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '余额方向代码，由类别决定：DEBIT/CREDIT',
  `enabled` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否启用：0否 1是',
  `is_system` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否系统预置科目：0否 1是',
  `sort_order` int(11) NOT NULL DEFAULT 0 COMMENT '排序号',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_fin_account_company_code`(`company_id`, `account_code`) USING BTREE,
  INDEX `idx_fin_account_category`(`category_code`, `enabled`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 187 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '会计科目表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of fin_account
-- ----------------------------
INSERT INTO `fin_account` VALUES (167, 1, '5001', '主营业务收入', 'REVENUE', 'CREDIT', 1, 0, 0, '', '2026-11-06 15:52:48', '2026-11-06 16:07:03');
INSERT INTO `fin_account` VALUES (168, 1, '5401', '主营业务成本', 'EXPENSE', 'DEBIT', 1, 0, 0, '', '2026-11-06 15:52:48', '2026-11-06 16:07:03');
INSERT INTO `fin_account` VALUES (169, 1, '5601', '销售费用', 'EXPENSE', 'DEBIT', 1, 0, 0, '', '2026-11-06 15:52:48', '2026-11-06 16:07:03');
INSERT INTO `fin_account` VALUES (170, 1, '5602', '管理费用', 'EXPENSE', 'DEBIT', 1, 0, 0, '', '2026-11-06 15:52:48', '2026-11-06 16:07:03');
INSERT INTO `fin_account` VALUES (171, 1, '1001', '库存现金', 'ASSET', 'DEBIT', 1, 0, 0, '', '2026-11-06 15:52:48', '2026-11-06 16:07:03');
INSERT INTO `fin_account` VALUES (172, 1, '1002', '银行存款', 'ASSET', 'DEBIT', 1, 0, 0, '', '2026-11-06 15:52:48', '2026-11-06 16:07:03');
INSERT INTO `fin_account` VALUES (173, 1, '1122', '应收账款', 'ASSET', 'DEBIT', 1, 0, 0, '', '2026-11-06 15:52:48', '2026-11-06 16:07:03');
INSERT INTO `fin_account` VALUES (174, 1, '2202', '应付账款', 'LIABILITY', 'CREDIT', 1, 0, 0, '', '2026-11-06 15:52:48', '2026-11-06 16:07:03');
INSERT INTO `fin_account` VALUES (175, 1, '4001', '实收资本', 'EQUITY', 'CREDIT', 1, 0, 0, '', '2026-11-06 15:52:48', '2026-11-06 16:07:03');
INSERT INTO `fin_account` VALUES (176, 1, '4103', '本年利润', 'EQUITY', 'CREDIT', 1, 0, 0, '', '2026-11-06 15:52:48', '2026-11-06 16:07:03');

-- ----------------------------
-- Table structure for fin_account_balance
-- ----------------------------
DROP TABLE IF EXISTS `fin_account_balance`;
CREATE TABLE `fin_account_balance`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '余额记录ID',
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `account_id` bigint(20) UNSIGNED NOT NULL COMMENT '会计科目ID',
  `period_code` char(7) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '会计期间YYYY-MM',
  `opening_debit` decimal(18, 2) NOT NULL DEFAULT 0.00 COMMENT '期初借方余额',
  `opening_credit` decimal(18, 2) NOT NULL DEFAULT 0.00 COMMENT '期初贷方余额',
  `current_debit` decimal(18, 2) NOT NULL DEFAULT 0.00 COMMENT '本期借方发生额',
  `current_credit` decimal(18, 2) NOT NULL DEFAULT 0.00 COMMENT '本期贷方发生额',
  `ending_debit` decimal(18, 2) NOT NULL DEFAULT 0.00 COMMENT '期末借方余额',
  `ending_credit` decimal(18, 2) NOT NULL DEFAULT 0.00 COMMENT '期末贷方余额',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_fin_balance_company_account_period`(`company_id`, `account_id`, `period_code`) USING BTREE,
  INDEX `idx_fin_balance_period`(`period_code`) USING BTREE,
  INDEX `fk_fin_balance_account`(`account_id`) USING BTREE,
  CONSTRAINT `fk_fin_balance_account` FOREIGN KEY (`account_id`) REFERENCES `fin_account` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 3269 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '科目期间余额表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of fin_account_balance
-- ----------------------------
INSERT INTO `fin_account_balance` VALUES (2493, 1, 167, '2026-01', 0.00, 0.00, 0.00, 37200.00, 0.00, 37200.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2494, 1, 168, '2026-01', 0.00, 0.00, 29696.00, 0.00, 29696.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2496, 1, 170, '2026-01', 0.00, 0.00, 8000.00, 0.00, 8000.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2497, 1, 171, '2026-01', 100000.00, 0.00, 0.00, 27696.00, 72304.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2498, 1, 172, '2026-01', 0.00, 0.00, 37200.00, 10000.00, 27200.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2502, 1, 176, '2026-01', 0.00, 100000.00, 0.00, 0.00, 0.00, 100000.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2503, 1, 167, '2026-02', 0.00, 37200.00, 0.00, 146600.00, 0.00, 183800.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2504, 1, 168, '2026-02', 29696.00, 0.00, 63040.00, 0.00, 92736.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2506, 1, 170, '2026-02', 8000.00, 0.00, 8000.00, 0.00, 16000.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2507, 1, 171, '2026-02', 72304.00, 0.00, 0.00, 0.00, 72304.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2508, 1, 172, '2026-02', 27200.00, 0.00, 146600.00, 71040.00, 102760.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2512, 1, 176, '2026-02', 0.00, 100000.00, 0.00, 0.00, 0.00, 100000.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2513, 1, 167, '2026-03', 0.00, 183800.00, 0.00, 55800.00, 0.00, 239600.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2514, 1, 168, '2026-03', 92736.00, 0.00, 98440.00, 0.00, 191176.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2515, 1, 169, '2026-03', 0.00, 0.00, 5000.00, 0.00, 5000.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2516, 1, 170, '2026-03', 16000.00, 0.00, 6000.00, 0.00, 22000.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2517, 1, 171, '2026-03', 72304.00, 0.00, 0.00, 0.00, 72304.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2518, 1, 172, '2026-03', 102760.00, 0.00, 55800.00, 109440.00, 49120.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2522, 1, 176, '2026-03', 0.00, 100000.00, 0.00, 0.00, 0.00, 100000.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2523, 1, 167, '2026-04', 0.00, 239600.00, 0.00, 134000.00, 0.00, 373600.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2524, 1, 168, '2026-04', 191176.00, 0.00, 23840.00, 0.00, 215016.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2525, 1, 169, '2026-04', 5000.00, 0.00, 0.00, 0.00, 5000.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2526, 1, 170, '2026-04', 22000.00, 0.00, 0.00, 0.00, 22000.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2527, 1, 171, '2026-04', 72304.00, 0.00, 0.00, 1200.00, 71104.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2528, 1, 172, '2026-04', 49120.00, 0.00, 134000.00, 22640.00, 160480.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2532, 1, 176, '2026-04', 0.00, 100000.00, 0.00, 0.00, 0.00, 100000.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2533, 1, 167, '2026-05', 0.00, 373600.00, 0.00, 0.00, 0.00, 373600.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2534, 1, 168, '2026-05', 215016.00, 0.00, 36000.00, 0.00, 251016.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2535, 1, 169, '2026-05', 5000.00, 0.00, 0.00, 0.00, 5000.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2536, 1, 170, '2026-05', 22000.00, 0.00, 12000.00, 0.00, 34000.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2537, 1, 171, '2026-05', 71104.00, 0.00, 0.00, 0.00, 71104.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2538, 1, 172, '2026-05', 160480.00, 0.00, 0.00, 48000.00, 112480.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2542, 1, 176, '2026-05', 0.00, 100000.00, 0.00, 0.00, 0.00, 100000.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2543, 1, 167, '2026-06', 0.00, 373600.00, 0.00, 76400.00, 0.00, 450000.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2544, 1, 168, '2026-06', 251016.00, 0.00, 42430.00, 0.00, 293446.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2545, 1, 169, '2026-06', 5000.00, 0.00, 14000.00, 0.00, 19000.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2546, 1, 170, '2026-06', 34000.00, 0.00, 8000.00, 0.00, 42000.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2547, 1, 171, '2026-06', 71104.00, 0.00, 0.00, 0.00, 71104.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2548, 1, 172, '2026-06', 112480.00, 0.00, 76400.00, 64430.00, 124450.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2552, 1, 176, '2026-06', 0.00, 100000.00, 0.00, 0.00, 0.00, 100000.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2553, 1, 167, '2026-07', 0.00, 450000.00, 0.00, 0.00, 0.00, 450000.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2554, 1, 168, '2026-07', 293446.00, 0.00, 14000.00, 0.00, 307446.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2555, 1, 169, '2026-07', 19000.00, 0.00, 0.00, 0.00, 19000.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2556, 1, 170, '2026-07', 42000.00, 0.00, 10000.00, 0.00, 52000.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2557, 1, 171, '2026-07', 71104.00, 0.00, 0.00, 0.00, 71104.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2558, 1, 172, '2026-07', 124450.00, 0.00, 0.00, 24000.00, 100450.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2562, 1, 176, '2026-07', 0.00, 100000.00, 0.00, 0.00, 0.00, 100000.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2563, 1, 167, '2026-08', 0.00, 450000.00, 0.00, 54000.00, 0.00, 504000.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2564, 1, 168, '2026-08', 307446.00, 0.00, 42760.00, 0.00, 350206.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2565, 1, 169, '2026-08', 19000.00, 0.00, 6000.00, 0.00, 25000.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2566, 1, 170, '2026-08', 52000.00, 0.00, 0.00, 0.00, 52000.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2567, 1, 171, '2026-08', 71104.00, 0.00, 0.00, 0.00, 71104.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2568, 1, 172, '2026-08', 100450.00, 0.00, 54000.00, 48760.00, 105690.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2572, 1, 176, '2026-08', 0.00, 100000.00, 0.00, 0.00, 0.00, 100000.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2573, 1, 167, '2026-09', 0.00, 504000.00, 0.00, 76000.00, 0.00, 580000.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2574, 1, 168, '2026-09', 350206.00, 0.00, 53840.00, 0.00, 404046.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2575, 1, 169, '2026-09', 25000.00, 0.00, 0.00, 0.00, 25000.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2576, 1, 170, '2026-09', 52000.00, 0.00, 12000.00, 0.00, 64000.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2577, 1, 171, '2026-09', 71104.00, 0.00, 0.00, 10000.00, 61104.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2578, 1, 172, '2026-09', 105690.00, 0.00, 86000.00, 65840.00, 125850.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2582, 1, 176, '2026-09', 0.00, 100000.00, 0.00, 0.00, 0.00, 100000.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2583, 1, 167, '2026-10', 0.00, 580000.00, 0.00, 73800.00, 0.00, 653800.00, '2026-11-06 16:18:01');
INSERT INTO `fin_account_balance` VALUES (2584, 1, 168, '2026-10', 404046.00, 0.00, 20000.00, 0.00, 424046.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2585, 1, 169, '2026-10', 25000.00, 0.00, 4000.00, 0.00, 29000.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2586, 1, 170, '2026-10', 64000.00, 0.00, 0.00, 0.00, 64000.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2587, 1, 171, '2026-10', 61104.00, 0.00, 35400.00, 0.00, 96504.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2588, 1, 172, '2026-10', 125850.00, 0.00, 38400.00, 24000.00, 140250.00, 0.00, '2026-11-06 16:18:01');
INSERT INTO `fin_account_balance` VALUES (2592, 1, 176, '2026-10', 0.00, 100000.00, 0.00, 0.00, 0.00, 100000.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2593, 1, 167, '2026-11', 0.00, 653800.00, 0.00, 20400.00, 0.00, 674200.00, '2026-11-06 16:37:00');
INSERT INTO `fin_account_balance` VALUES (2594, 1, 168, '2026-11', 424046.00, 0.00, 10640.00, 0.00, 434686.00, 0.00, '2026-11-06 16:37:19');
INSERT INTO `fin_account_balance` VALUES (2595, 1, 169, '2026-11', 29000.00, 0.00, 0.00, 0.00, 29000.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2596, 1, 170, '2026-11', 64000.00, 0.00, 0.00, 0.00, 64000.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2597, 1, 171, '2026-11', 96504.00, 0.00, 20400.00, 4600.00, 112304.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2598, 1, 172, '2026-11', 140250.00, 0.00, 0.00, 6040.00, 134210.00, 0.00, '2026-11-06 16:37:19');
INSERT INTO `fin_account_balance` VALUES (2602, 1, 176, '2026-11', 0.00, 100000.00, 0.00, 0.00, 0.00, 100000.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2603, 1, 167, '2026-12', 0.00, 674200.00, 0.00, 0.00, 0.00, 674200.00, '2026-11-06 16:37:00');
INSERT INTO `fin_account_balance` VALUES (2604, 1, 168, '2026-12', 434686.00, 0.00, 0.00, 0.00, 434686.00, 0.00, '2026-11-06 16:37:19');
INSERT INTO `fin_account_balance` VALUES (2605, 1, 169, '2026-12', 29000.00, 0.00, 0.00, 0.00, 29000.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2606, 1, 170, '2026-12', 64000.00, 0.00, 0.00, 0.00, 64000.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2607, 1, 171, '2026-12', 112304.00, 0.00, 0.00, 0.00, 112304.00, 0.00, '2026-11-05 12:00:00');
INSERT INTO `fin_account_balance` VALUES (2608, 1, 172, '2026-12', 134210.00, 0.00, 0.00, 0.00, 134210.00, 0.00, '2026-11-06 16:37:19');
INSERT INTO `fin_account_balance` VALUES (2612, 1, 176, '2026-12', 0.00, 100000.00, 0.00, 0.00, 0.00, 100000.00, '2026-11-05 12:00:00');

-- ----------------------------
-- Table structure for fin_accounting_period
-- ----------------------------
DROP TABLE IF EXISTS `fin_accounting_period`;
CREATE TABLE `fin_accounting_period`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '期间ID',
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `period_code` char(7) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '会计期间，格式YYYY-MM',
  `period_year` smallint(5) UNSIGNED NOT NULL COMMENT '会计年度',
  `period_month` tinyint(3) UNSIGNED NOT NULL COMMENT '会计月份，1-12',
  `status_code` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'FUTURE' COMMENT '期间状态代码：OPEN/CLOSED/FUTURE',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_fin_period_company_code`(`company_id`, `period_code`) USING BTREE,
  INDEX `idx_fin_period_year_month`(`period_year`, `period_month`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 313 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '会计期间表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of fin_accounting_period
-- ----------------------------
INSERT INTO `fin_accounting_period` VALUES (253, 1, '2026-01', 2026, 1, 'OPEN', '2026-11-06 15:52:48', '2026-11-06 15:52:48');
INSERT INTO `fin_accounting_period` VALUES (254, 1, '2026-02', 2026, 2, 'OPEN', '2026-11-06 15:52:48', '2026-11-06 16:07:44');
INSERT INTO `fin_accounting_period` VALUES (255, 1, '2026-03', 2026, 3, 'OPEN', '2026-11-06 15:52:48', '2026-11-06 16:07:44');
INSERT INTO `fin_accounting_period` VALUES (256, 1, '2026-04', 2026, 4, 'OPEN', '2026-11-06 15:52:48', '2026-11-06 16:07:44');
INSERT INTO `fin_accounting_period` VALUES (257, 1, '2026-05', 2026, 5, 'OPEN', '2026-11-06 15:52:48', '2026-11-06 16:07:44');
INSERT INTO `fin_accounting_period` VALUES (258, 1, '2026-06', 2026, 6, 'OPEN', '2026-11-06 15:52:48', '2026-11-06 16:07:44');
INSERT INTO `fin_accounting_period` VALUES (259, 1, '2026-07', 2026, 7, 'OPEN', '2026-11-06 15:52:48', '2026-11-06 16:07:44');
INSERT INTO `fin_accounting_period` VALUES (260, 1, '2026-08', 2026, 8, 'OPEN', '2026-11-06 15:52:48', '2026-11-06 16:07:44');
INSERT INTO `fin_accounting_period` VALUES (261, 1, '2026-09', 2026, 9, 'OPEN', '2026-11-06 15:52:48', '2026-11-06 16:07:44');
INSERT INTO `fin_accounting_period` VALUES (262, 1, '2026-10', 2026, 10, 'OPEN', '2026-11-06 15:52:48', '2026-11-06 16:07:44');
INSERT INTO `fin_accounting_period` VALUES (263, 1, '2026-11', 2026, 11, 'OPEN', '2026-11-06 15:52:48', '2026-11-06 16:07:44');
INSERT INTO `fin_accounting_period` VALUES (264, 1, '2026-12', 2026, 12, 'FUTURE', '2026-11-06 15:52:48', '2026-11-06 15:52:48');

-- ----------------------------
-- Table structure for fin_dictionary
-- ----------------------------
DROP TABLE IF EXISTS `fin_dictionary`;
CREATE TABLE `fin_dictionary`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `dict_type` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '字典类别代码',
  `dict_code` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '选项代码（前后端共用）',
  `label` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '显示名称',
  `sort_order` int(11) NOT NULL DEFAULT 0 COMMENT '排序号',
  `enabled` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否启用：0否 1是',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_fin_dictionary_type_code`(`dict_type`, `dict_code`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '财务枚举字典' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of fin_dictionary
-- ----------------------------

-- ----------------------------
-- Table structure for fin_fund_account
-- ----------------------------
DROP TABLE IF EXISTS `fin_fund_account`;
CREATE TABLE `fin_fund_account`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '资金账户ID',
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `account_type` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '账户类型代码：CASH/BANK',
  `account_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '账户显示名称',
  `bank_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '开户银行名称，现金账户为空',
  `account_last4` char(4) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '银行卡末四位，不保存完整卡号',
  `ledger_account_id` bigint(20) UNSIGNED NOT NULL COMMENT '关联会计科目ID',
  `enabled` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否启用：0否 1是',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_fin_fund_account_type_enabled`(`account_type`, `enabled`) USING BTREE,
  INDEX `idx_fin_fund_account_ledger`(`ledger_account_id`) USING BTREE,
  CONSTRAINT `fk_fin_fund_account_ledger` FOREIGN KEY (`ledger_account_id`) REFERENCES `fin_account` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 28 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '现金及银行资金账户' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of fin_fund_account
-- ----------------------------
INSERT INTO `fin_fund_account` VALUES (24, 1, 'CASH', '库存现金', '', '', 171, 1, '2026-11-06 15:52:48', '2026-11-06 16:07:03');
INSERT INTO `fin_fund_account` VALUES (25, 1, 'BANK', '公司基本户', '中国银行', '2578', 172, 1, '2026-11-06 15:52:48', '2026-11-06 16:07:03');

-- ----------------------------
-- Table structure for fin_setup_config
-- ----------------------------
DROP TABLE IF EXISTS `fin_setup_config`;
CREATE TABLE `fin_setup_config`  (
  `id` int(10) UNSIGNED NOT NULL COMMENT '主键，等于公司ID',
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `is_initialized` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否完成首次初始化：0否 1是',
  `period_year` smallint(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT '当前初始化会计年度',
  `start_month` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT '启用起始月份，1-12',
  `initialized_by` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '最后配置人site_user.id',
  `initialized_at` datetime NULL DEFAULT NULL COMMENT '首次完成时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_fin_setup_company`(`company_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '财务初始化状态及配置' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of fin_setup_config
-- ----------------------------
INSERT INTO `fin_setup_config` VALUES (1, 1, 1, 2026, 1, 1, '2026-11-06 15:52:48', '2026-11-06 15:52:48');

-- ----------------------------
-- Table structure for fin_transaction
-- ----------------------------
DROP TABLE IF EXISTS `fin_transaction`;
CREATE TABLE `fin_transaction`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '资金流水ID',
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `source_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'MANUAL' COMMENT '来源类型：MANUAL手工，预留业务来源扩展',
  `source_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '来源业务记录ID',
  `account_id` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '关联会计科目ID',
  `occurred_at` datetime NOT NULL COMMENT '发生时间',
  `amount` decimal(18, 2) NOT NULL COMMENT '流水金额（正数）',
  `direction_code` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '资金方向代码：INCOME/EXPENSE',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_fin_transaction_source`(`source_type`, `source_id`) USING BTREE,
  INDEX `idx_fin_transaction_occurred_at`(`occurred_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 336 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '资金流水表（本期仅建表，不写入初始化流水）' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of fin_transaction
-- ----------------------------
INSERT INTO `fin_transaction` VALUES (247, 1, 'PROJECT_PAYOUT', 1, 171, '2026-01-05 10:00:00', 8000.00, 'EXPENSE', '资金流出', '2026-01-05 10:00:00');
INSERT INTO `fin_transaction` VALUES (248, 1, 'PROJECT_PAYMENT', 1, 172, '2026-01-12 10:00:00', 37200.00, 'INCOME', '收到张伟明 · 锦绣路全屋翻新定金', '2026-01-12 10:00:00');
INSERT INTO `fin_transaction` VALUES (249, 1, 'PROJECT_PAYOUT', 2, 171, '2026-01-13 10:00:00', 5120.00, 'EXPENSE', '资金流出', '2026-01-13 10:00:00');
INSERT INTO `fin_transaction` VALUES (250, 1, 'PROJECT_PAYOUT', 3, 171, '2026-01-14 10:00:00', 2520.00, 'EXPENSE', '资金流出', '2026-01-14 10:00:00');
INSERT INTO `fin_transaction` VALUES (251, 1, 'PROJECT_PAYOUT', 4, 171, '2026-01-16 10:00:00', 4416.00, 'EXPENSE', '资金流出', '2026-01-16 10:00:00');
INSERT INTO `fin_transaction` VALUES (252, 1, 'PROJECT_PAYOUT', 5, 171, '2026-01-18 10:00:00', 1360.00, 'EXPENSE', '资金流出', '2026-01-18 10:00:00');
INSERT INTO `fin_transaction` VALUES (253, 1, 'PROJECT_PAYOUT', 6, 171, '2026-01-18 10:00:00', 1280.00, 'EXPENSE', '资金流出', '2026-01-18 10:00:00');
INSERT INTO `fin_transaction` VALUES (254, 1, 'PROJECT_PAYOUT', 7, 172, '2026-01-20 10:00:00', 10000.00, 'EXPENSE', '资金流出', '2026-01-20 10:00:00');
INSERT INTO `fin_transaction` VALUES (255, 1, 'PROJECT_PAYOUT', 8, 171, '2026-01-25 10:00:00', 5000.00, 'EXPENSE', '资金流出', '2026-01-25 10:00:00');
INSERT INTO `fin_transaction` VALUES (256, 1, 'PROJECT_PAYOUT', 9, 172, '2026-02-05 10:00:00', 8000.00, 'EXPENSE', '资金流出', '2026-02-05 10:00:00');
INSERT INTO `fin_transaction` VALUES (257, 1, 'PROJECT_PAYOUT', 10, 172, '2026-02-10 10:00:00', 15000.00, 'EXPENSE', '资金流出', '2026-02-10 10:00:00');
INSERT INTO `fin_transaction` VALUES (258, 1, 'PROJECT_PAYMENT', 2, 172, '2026-02-15 10:00:00', 93000.00, 'INCOME', '收到张伟明 · 锦绣路全屋翻新中期款', '2026-02-15 10:00:00');
INSERT INTO `fin_transaction` VALUES (259, 1, 'PROJECT_PAYMENT', 3, 172, '2026-02-18 10:00:00', 53600.00, 'INCOME', '收到王丽华 · 漕溪北路新房装修定金', '2026-02-18 10:00:00');
INSERT INTO `fin_transaction` VALUES (260, 1, 'PROJECT_PAYOUT', 11, 172, '2026-02-22 10:00:00', 9440.00, 'EXPENSE', '资金流出', '2026-02-22 10:00:00');
INSERT INTO `fin_transaction` VALUES (261, 1, 'PROJECT_PAYOUT', 12, 172, '2026-02-25 10:00:00', 16000.00, 'EXPENSE', '资金流出', '2026-02-25 10:00:00');
INSERT INTO `fin_transaction` VALUES (262, 1, 'PROJECT_PAYOUT', 13, 172, '2026-02-28 10:00:00', 12000.00, 'EXPENSE', '资金流出', '2026-02-28 10:00:00');
INSERT INTO `fin_transaction` VALUES (263, 1, 'PROJECT_PAYOUT', 14, 172, '2026-02-28 10:00:00', 10600.00, 'EXPENSE', '资金流出', '2026-02-28 10:00:00');
INSERT INTO `fin_transaction` VALUES (264, 1, 'PROJECT_PAYOUT', 15, 172, '2026-03-03 10:00:00', 5200.00, 'EXPENSE', '资金流出', '2026-03-03 10:00:00');
INSERT INTO `fin_transaction` VALUES (265, 1, 'PROJECT_PAYOUT', 16, 172, '2026-03-05 10:00:00', 10000.00, 'EXPENSE', '资金流出', '2026-03-05 10:00:00');
INSERT INTO `fin_transaction` VALUES (266, 1, 'PROJECT_PAYOUT', 17, 172, '2026-03-07 10:00:00', 1800.00, 'EXPENSE', '资金流出', '2026-03-07 10:00:00');
INSERT INTO `fin_transaction` VALUES (267, 1, 'PROJECT_PAYOUT', 18, 172, '2026-03-08 10:00:00', 30000.00, 'EXPENSE', '资金流出', '2026-03-08 10:00:00');
INSERT INTO `fin_transaction` VALUES (268, 1, 'PROJECT_PAYOUT', 19, 172, '2026-03-10 10:00:00', 5000.00, 'EXPENSE', '资金流出', '2026-03-10 10:00:00');
INSERT INTO `fin_transaction` VALUES (269, 1, 'PROJECT_PAYOUT', 20, 172, '2026-03-12 10:00:00', 12400.00, 'EXPENSE', '资金流出', '2026-03-12 10:00:00');
INSERT INTO `fin_transaction` VALUES (270, 1, 'PROJECT_PAYOUT', 21, 172, '2026-03-15 10:00:00', 6000.00, 'EXPENSE', '资金流出', '2026-03-15 10:00:00');
INSERT INTO `fin_transaction` VALUES (271, 1, 'PROJECT_PAYOUT', 22, 172, '2026-03-17 10:00:00', 8040.00, 'EXPENSE', '资金流出', '2026-03-17 10:00:00');
INSERT INTO `fin_transaction` VALUES (272, 1, 'PROJECT_PAYMENT', 4, 172, '2026-03-18 10:00:00', 55800.00, 'INCOME', '收到张伟明 · 锦绣路全屋翻新尾款', '2026-03-18 10:00:00');
INSERT INTO `fin_transaction` VALUES (273, 1, 'PROJECT_PAYOUT', 23, 172, '2026-03-18 10:00:00', 9000.00, 'EXPENSE', '资金流出', '2026-03-18 10:00:00');
INSERT INTO `fin_transaction` VALUES (274, 1, 'PROJECT_PAYOUT', 24, 172, '2026-03-28 10:00:00', 22000.00, 'EXPENSE', '资金流出', '2026-03-28 10:00:00');
INSERT INTO `fin_transaction` VALUES (275, 1, 'PROJECT_PAYOUT', 25, 172, '2026-04-04 10:00:00', 640.00, 'EXPENSE', '资金流出', '2026-04-04 10:00:00');
INSERT INTO `fin_transaction` VALUES (276, 1, 'PROJECT_PAYOUT', 26, 171, '2026-04-15 10:00:00', 1200.00, 'EXPENSE', '资金流出', '2026-04-15 10:00:00');
INSERT INTO `fin_transaction` VALUES (277, 1, 'PROJECT_PAYOUT', 27, 172, '2026-04-18 10:00:00', 22000.00, 'EXPENSE', '资金流出', '2026-04-18 10:00:00');
INSERT INTO `fin_transaction` VALUES (278, 1, 'PROJECT_PAYMENT', 5, 172, '2026-04-20 10:00:00', 134000.00, 'INCOME', '收到王丽华 · 漕溪北路新房装修中期款', '2026-04-20 10:00:00');
INSERT INTO `fin_transaction` VALUES (279, 1, 'PROJECT_PAYOUT', 28, 172, '2026-05-10 10:00:00', 12000.00, 'EXPENSE', '资金流出', '2026-05-10 10:00:00');
INSERT INTO `fin_transaction` VALUES (280, 1, 'PROJECT_PAYOUT', 29, 172, '2026-05-20 10:00:00', 18000.00, 'EXPENSE', '资金流出', '2026-05-20 10:00:00');
INSERT INTO `fin_transaction` VALUES (281, 1, 'PROJECT_PAYOUT', 30, 172, '2026-05-22 10:00:00', 18000.00, 'EXPENSE', '资金流出', '2026-05-22 10:00:00');
INSERT INTO `fin_transaction` VALUES (282, 1, 'PROJECT_PAYMENT', 6, 172, '2026-06-05 10:00:00', 22800.00, 'INCOME', '收到赵磊 · 金桥商铺装修定金', '2026-06-05 10:00:00');
INSERT INTO `fin_transaction` VALUES (283, 1, 'PROJECT_PAYOUT', 31, 172, '2026-06-08 10:00:00', 4000.00, 'EXPENSE', '资金流出', '2026-06-08 10:00:00');
INSERT INTO `fin_transaction` VALUES (284, 1, 'PROJECT_PAYOUT', 32, 172, '2026-06-08 10:00:00', 8000.00, 'EXPENSE', '资金流出', '2026-06-08 10:00:00');
INSERT INTO `fin_transaction` VALUES (285, 1, 'PROJECT_PAYOUT', 33, 172, '2026-06-10 10:00:00', 12000.00, 'EXPENSE', '资金流出', '2026-06-10 10:00:00');
INSERT INTO `fin_transaction` VALUES (286, 1, 'PROJECT_PAYOUT', 34, 172, '2026-06-14 10:00:00', 4050.00, 'EXPENSE', '资金流出', '2026-06-14 10:00:00');
INSERT INTO `fin_transaction` VALUES (287, 1, 'PROJECT_PAYOUT', 35, 172, '2026-06-17 10:00:00', 5100.00, 'EXPENSE', '资金流出', '2026-06-17 10:00:00');
INSERT INTO `fin_transaction` VALUES (288, 1, 'PROJECT_PAYOUT', 36, 172, '2026-06-18 10:00:00', 6000.00, 'EXPENSE', '资金流出', '2026-06-18 10:00:00');
INSERT INTO `fin_transaction` VALUES (289, 1, 'PROJECT_PAYOUT', 37, 172, '2026-06-20 10:00:00', 7280.00, 'EXPENSE', '资金流出', '2026-06-20 10:00:00');
INSERT INTO `fin_transaction` VALUES (290, 1, 'PROJECT_PAYMENT', 7, 172, '2026-06-25 10:00:00', 53600.00, 'INCOME', '收到王丽华 · 漕溪北路新房装修尾款', '2026-06-25 10:00:00');
INSERT INTO `fin_transaction` VALUES (291, 1, 'PROJECT_PAYOUT', 38, 172, '2026-06-25 10:00:00', 10000.00, 'EXPENSE', '资金流出', '2026-06-25 10:00:00');
INSERT INTO `fin_transaction` VALUES (292, 1, 'PROJECT_PAYOUT', 39, 172, '2026-06-28 10:00:00', 8000.00, 'EXPENSE', '资金流出', '2026-06-28 10:00:00');
INSERT INTO `fin_transaction` VALUES (293, 1, 'PROJECT_PAYOUT', 40, 172, '2026-07-05 10:00:00', 8000.00, 'EXPENSE', '资金流出', '2026-07-05 10:00:00');
INSERT INTO `fin_transaction` VALUES (294, 1, 'PROJECT_PAYOUT', 41, 172, '2026-07-15 10:00:00', 6000.00, 'EXPENSE', '资金流出', '2026-07-15 10:00:00');
INSERT INTO `fin_transaction` VALUES (295, 1, 'PROJECT_PAYOUT', 42, 172, '2026-07-20 10:00:00', 8000.00, 'EXPENSE', '资金流出', '2026-07-20 10:00:00');
INSERT INTO `fin_transaction` VALUES (296, 1, 'PROJECT_PAYOUT', 43, 172, '2026-07-28 10:00:00', 2000.00, 'EXPENSE', '资金流出', '2026-07-28 10:00:00');
INSERT INTO `fin_transaction` VALUES (297, 1, 'PROJECT_PAYOUT', 44, 172, '2026-08-01 10:00:00', 4000.00, 'EXPENSE', '资金流出', '2026-08-01 10:00:00');
INSERT INTO `fin_transaction` VALUES (298, 1, 'PROJECT_PAYMENT', 8, 172, '2026-08-05 10:00:00', 23600.00, 'INCOME', '收到吴敏 · 恒丰路老房翻新定金', '2026-08-05 10:00:00');
INSERT INTO `fin_transaction` VALUES (299, 1, 'PROJECT_PAYOUT', 45, 172, '2026-08-06 10:00:00', 6000.00, 'EXPENSE', '资金流出', '2026-08-06 10:00:00');
INSERT INTO `fin_transaction` VALUES (300, 1, 'PROJECT_PAYOUT', 46, 172, '2026-08-08 10:00:00', 8000.00, 'EXPENSE', '资金流出', '2026-08-08 10:00:00');
INSERT INTO `fin_transaction` VALUES (301, 1, 'PROJECT_PAYOUT', 47, 172, '2026-08-12 10:00:00', 3800.00, 'EXPENSE', '资金流出', '2026-08-12 10:00:00');
INSERT INTO `fin_transaction` VALUES (302, 1, 'PROJECT_PAYOUT', 48, 172, '2026-08-14 10:00:00', 8160.00, 'EXPENSE', '资金流出', '2026-08-14 10:00:00');
INSERT INTO `fin_transaction` VALUES (303, 1, 'PROJECT_PAYOUT', 49, 172, '2026-08-17 10:00:00', 5120.00, 'EXPENSE', '资金流出', '2026-08-17 10:00:00');
INSERT INTO `fin_transaction` VALUES (304, 1, 'PROJECT_PAYOUT', 50, 172, '2026-08-18 10:00:00', 1680.00, 'EXPENSE', '资金流出', '2026-08-18 10:00:00');
INSERT INTO `fin_transaction` VALUES (305, 1, 'PROJECT_PAYMENT', 9, 172, '2026-08-20 10:00:00', 30400.00, 'INCOME', '收到赵磊 · 金桥商铺装修尾款', '2026-08-20 10:00:00');
INSERT INTO `fin_transaction` VALUES (306, 1, 'PROJECT_PAYOUT', 51, 172, '2026-08-28 10:00:00', 12000.00, 'EXPENSE', '资金流出', '2026-08-28 10:00:00');
INSERT INTO `fin_transaction` VALUES (307, 1, 'PROJECT_PAYMENT', 10, 172, '2026-09-03 10:00:00', 28800.00, 'INCOME', '收到陈志强 · 仲益大厦办公室改造定金', '2026-09-03 10:00:00');
INSERT INTO `fin_transaction` VALUES (308, 1, 'PROJECT_PAYOUT', 52, 172, '2026-09-05 10:00:00', 8000.00, 'EXPENSE', '资金流出', '2026-09-05 10:00:00');
INSERT INTO `fin_transaction` VALUES (309, 1, 'PROJECT_PAYOUT', 53, 172, '2026-09-07 10:00:00', 2400.00, 'EXPENSE', '资金流出', '2026-09-07 10:00:00');
INSERT INTO `fin_transaction` VALUES (310, 1, 'PROJECT_PAYOUT', 54, 172, '2026-09-10 10:00:00', 5400.00, 'EXPENSE', '资金流出', '2026-09-10 10:00:00');
INSERT INTO `fin_transaction` VALUES (311, 1, 'PROJECT_PAYOUT', 55, 172, '2026-09-10 10:00:00', 7000.00, 'EXPENSE', '资金流出', '2026-09-10 10:00:00');
INSERT INTO `fin_transaction` VALUES (312, 1, 'PROJECT_PAYOUT', 56, 172, '2026-09-15 10:00:00', 5000.00, 'EXPENSE', '资金流出', '2026-09-15 10:00:00');
INSERT INTO `fin_transaction` VALUES (313, 1, 'PROJECT_PAYMENT', 11, 172, '2026-09-18 10:00:00', 47200.00, 'INCOME', '收到吴敏 · 恒丰路老房翻新中期款', '2026-09-18 10:00:00');
INSERT INTO `fin_transaction` VALUES (314, 1, 'PROJECT_PAYOUT', 57, 172, '2026-09-25 10:00:00', 8000.00, 'EXPENSE', '资金流出', '2026-09-25 10:00:00');
INSERT INTO `fin_transaction` VALUES (315, 1, 'PROJECT_PAYOUT', 58, 172, '2026-09-28 10:00:00', 19800.00, 'EXPENSE', '资金流出', '2026-09-28 10:00:00');
INSERT INTO `fin_transaction` VALUES (316, 1, 'PROJECT_PAYOUT', 59, 172, '2026-09-30 10:00:00', 3120.00, 'EXPENSE', '资金流出', '2026-09-30 10:00:00');
INSERT INTO `fin_transaction` VALUES (317, 1, 'PROJECT_PAYOUT', 60, 172, '2026-09-30 10:00:00', 3120.00, 'EXPENSE', '资金流出', '2026-09-30 10:00:00');
INSERT INTO `fin_transaction` VALUES (318, 1, 'PROJECT_PAYOUT', 61, 172, '2026-09-30 10:00:00', 4000.00, 'EXPENSE', '资金流出', '2026-09-30 10:00:00');
INSERT INTO `fin_transaction` VALUES (319, 1, 'MANUAL', 0, 172, '2026-09-30 10:00:00', 10000.00, 'INCOME', '现金存入银行（内部转账）', '2026-09-30 10:00:00');
INSERT INTO `fin_transaction` VALUES (320, 1, 'MANUAL', 0, 171, '2026-09-30 10:00:00', 10000.00, 'EXPENSE', '资金流出', '2026-09-30 10:00:00');
INSERT INTO `fin_transaction` VALUES (321, 1, 'PROJECT_PAYOUT', 62, 172, '2026-10-06 10:00:00', 12000.00, 'EXPENSE', '资金流出', '2026-10-06 10:00:00');
INSERT INTO `fin_transaction` VALUES (322, 1, 'PROJECT_PAYOUT', 63, 172, '2026-10-09 10:00:00', 4000.00, 'EXPENSE', '资金流出', '2026-10-09 10:00:00');
INSERT INTO `fin_transaction` VALUES (323, 1, 'PROJECT_PAYMENT', 12, 172, '2026-10-15 10:00:00', 38400.00, 'INCOME', '收到陈志强 · 仲益大厦办公室改造中期款', '2026-10-15 10:00:00');
INSERT INTO `fin_transaction` VALUES (325, 1, 'PROJECT_PAYMENT', 14, 171, '2026-10-25 10:00:00', 35400.00, 'INCOME', '收到吴敏 · 恒丰路老房翻新尾款', '2026-10-25 10:00:00');
INSERT INTO `fin_transaction` VALUES (326, 1, 'PROJECT_PAYOUT', 64, 172, '2026-10-28 10:00:00', 8000.00, 'EXPENSE', '资金流出', '2026-10-28 10:00:00');
INSERT INTO `fin_transaction` VALUES (327, 1, 'PROJECT_PAYOUT', 65, 172, '2026-11-01 10:00:00', 3040.00, 'EXPENSE', '资金流出', '2026-11-01 10:00:00');
INSERT INTO `fin_transaction` VALUES (328, 1, 'PROJECT_PAYMENT', 15, 171, '2026-11-02 10:00:00', 20400.00, 'INCOME', '收到黄建军 · 国定路厨卫改造定金', '2026-11-02 10:00:00');
INSERT INTO `fin_transaction` VALUES (329, 1, 'PROJECT_PAYOUT', 66, 171, '2026-11-03 10:00:00', 2000.00, 'EXPENSE', '资金流出', '2026-11-03 10:00:00');
INSERT INTO `fin_transaction` VALUES (330, 1, 'PROJECT_PAYOUT', 67, 172, '2026-11-03 10:00:00', 3000.00, 'EXPENSE', '资金流出', '2026-11-03 10:00:00');
INSERT INTO `fin_transaction` VALUES (331, 1, 'PROJECT_PAYOUT', 68, 171, '2026-11-03 10:00:00', 2600.00, 'EXPENSE', '资金流出', '2026-11-03 10:00:00');

-- ----------------------------
-- Table structure for fin_voucher
-- ----------------------------
DROP TABLE IF EXISTS `fin_voucher`;
CREATE TABLE `fin_voucher`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '凭证ID',
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `voucher_no` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '凭证编号',
  `voucher_date` date NOT NULL COMMENT '凭证日期',
  `period_code` char(7) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '所属会计期间YYYY-MM',
  `voucher_type_code` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '凭证类型代码',
  `summary` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '凭证摘要',
  `source_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'MANUAL' COMMENT '来源类型：MANUAL手工，预留自动来源扩展',
  `source_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '来源业务记录ID，手工凭证为0',
  `status_code` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'DRAFT' COMMENT '凭证状态代码：DRAFT/POSTED/VOID',
  `created_by` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '制单人site_user.id',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_fin_voucher_company_no`(`company_id`, `voucher_no`) USING BTREE,
  INDEX `idx_fin_voucher_period_date`(`period_code`, `voucher_date`) USING BTREE,
  INDEX `idx_fin_voucher_source`(`source_type`, `source_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 91 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '会计凭证主表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of fin_voucher
-- ----------------------------
INSERT INTO `fin_voucher` VALUES (1, 1, '付-2026-01-0001', '2026-01-05', '2026-01', 'PAYMENT', '支付上海汇智物业管理有限公司管理费', 'PROJECT_PAYOUT', 1, 'POSTED', 4, '2026-01-05 10:00:00', '2026-01-05 10:00:00');
INSERT INTO `fin_voucher` VALUES (2, 1, '收-2026-01-0002', '2026-01-12', '2026-01', 'RECEIPT', '收到张伟明 · 锦绣路全屋翻新定金', 'PROJECT_PAYMENT', 1, 'POSTED', 4, '2026-01-12 10:00:00', '2026-01-12 10:00:00');
INSERT INTO `fin_voucher` VALUES (3, 1, '付-2026-01-0003', '2026-01-13', '2026-01', 'PAYMENT', '支付红星美凯龙建材（上海供应处）材料款', 'PROJECT_PAYOUT', 2, 'POSTED', 4, '2026-01-13 10:00:00', '2026-01-13 10:00:00');
INSERT INTO `fin_voucher` VALUES (4, 1, '付-2026-01-0004', '2026-01-14', '2026-01', 'PAYMENT', '支付红星美凯龙建材（上海供应处）材料款', 'PROJECT_PAYOUT', 3, 'POSTED', 4, '2026-01-14 10:00:00', '2026-01-14 10:00:00');
INSERT INTO `fin_voucher` VALUES (5, 1, '付-2026-01-0005', '2026-01-16', '2026-01', 'PAYMENT', '支付立邦涂料上海总代理材料款', 'PROJECT_PAYOUT', 4, 'POSTED', 4, '2026-01-16 10:00:00', '2026-01-16 10:00:00');
INSERT INTO `fin_voucher` VALUES (6, 1, '付-2026-01-0006', '2026-01-18', '2026-01', 'PAYMENT', '支付红星美凯龙建材（上海供应处）材料款', 'PROJECT_PAYOUT', 5, 'POSTED', 4, '2026-01-18 10:00:00', '2026-01-18 10:00:00');
INSERT INTO `fin_voucher` VALUES (7, 1, '付-2026-01-0007', '2026-01-18', '2026-01', 'PAYMENT', '支付红星美凯龙建材（上海供应处）材料款', 'PROJECT_PAYOUT', 6, 'POSTED', 4, '2026-01-18 10:00:00', '2026-01-18 10:00:00');
INSERT INTO `fin_voucher` VALUES (8, 1, '付-2026-01-0008', '2026-01-20', '2026-01', 'PAYMENT', '支付红星美凯龙建材（上海供应处）材料款', 'PROJECT_PAYOUT', 7, 'POSTED', 4, '2026-01-20 10:00:00', '2026-01-20 10:00:00');
INSERT INTO `fin_voucher` VALUES (9, 1, '付-2026-01-0009', '2026-01-25', '2026-01', 'PAYMENT', '支付上海安居物流配送有限公司其他', 'PROJECT_PAYOUT', 8, 'POSTED', 4, '2026-01-25 10:00:00', '2026-01-25 10:00:00');
INSERT INTO `fin_voucher` VALUES (10, 1, '付-2026-02-0001', '2026-02-05', '2026-02', 'PAYMENT', '支付上海汇智物业管理有限公司管理费', 'PROJECT_PAYOUT', 9, 'POSTED', 4, '2026-02-05 10:00:00', '2026-02-05 10:00:00');
INSERT INTO `fin_voucher` VALUES (11, 1, '付-2026-02-0002', '2026-02-10', '2026-02', 'PAYMENT', '支付沪上名家水电安装队人工费', 'PROJECT_PAYOUT', 10, 'POSTED', 4, '2026-02-10 10:00:00', '2026-02-10 10:00:00');
INSERT INTO `fin_voucher` VALUES (12, 1, '收-2026-02-0003', '2026-02-15', '2026-02', 'RECEIPT', '收到张伟明 · 锦绣路全屋翻新中期款', 'PROJECT_PAYMENT', 2, 'POSTED', 4, '2026-02-15 10:00:00', '2026-02-15 10:00:00');
INSERT INTO `fin_voucher` VALUES (13, 1, '收-2026-02-0004', '2026-02-18', '2026-02', 'RECEIPT', '收到王丽华 · 漕溪北路新房装修定金', 'PROJECT_PAYMENT', 3, 'POSTED', 4, '2026-02-18 10:00:00', '2026-02-18 10:00:00');
INSERT INTO `fin_voucher` VALUES (14, 1, '付-2026-02-0005', '2026-02-22', '2026-02', 'PAYMENT', '支付红星美凯龙建材（上海供应处）材料款', 'PROJECT_PAYOUT', 11, 'POSTED', 4, '2026-02-22 10:00:00', '2026-02-22 10:00:00');
INSERT INTO `fin_voucher` VALUES (15, 1, '付-2026-02-0006', '2026-02-25', '2026-02', 'PAYMENT', '支付上海素合空间设计工作室设计费', 'PROJECT_PAYOUT', 12, 'POSTED', 4, '2026-02-25 10:00:00', '2026-02-25 10:00:00');
INSERT INTO `fin_voucher` VALUES (16, 1, '付-2026-02-0007', '2026-02-28', '2026-02', 'PAYMENT', '支付欧派家居上海经销商材料款', 'PROJECT_PAYOUT', 13, 'POSTED', 4, '2026-02-28 10:00:00', '2026-02-28 10:00:00');
INSERT INTO `fin_voucher` VALUES (17, 1, '付-2026-02-0008', '2026-02-28', '2026-02', 'PAYMENT', '支付沪上名家水电安装队人工费', 'PROJECT_PAYOUT', 14, 'POSTED', 4, '2026-02-28 10:00:00', '2026-02-28 10:00:00');
INSERT INTO `fin_voucher` VALUES (18, 1, '付-2026-03-0001', '2026-03-03', '2026-03', 'PAYMENT', '支付立邦涂料上海总代理材料款', 'PROJECT_PAYOUT', 15, 'POSTED', 4, '2026-03-03 10:00:00', '2026-03-03 10:00:00');
INSERT INTO `fin_voucher` VALUES (19, 1, '付-2026-03-0002', '2026-03-05', '2026-03', 'PAYMENT', '支付红星美凯龙建材（上海供应处）材料款', 'PROJECT_PAYOUT', 16, 'POSTED', 4, '2026-03-05 10:00:00', '2026-03-05 10:00:00');
INSERT INTO `fin_voucher` VALUES (20, 1, '付-2026-03-0003', '2026-03-07', '2026-03', 'PAYMENT', '支付红星美凯龙建材（上海供应处）材料款', 'PROJECT_PAYOUT', 17, 'POSTED', 4, '2026-03-07 10:00:00', '2026-03-07 10:00:00');
INSERT INTO `fin_voucher` VALUES (21, 1, '付-2026-03-0004', '2026-03-08', '2026-03', 'PAYMENT', '支付红星美凯龙建材（上海供应处）材料款', 'PROJECT_PAYOUT', 18, 'POSTED', 4, '2026-03-08 10:00:00', '2026-03-08 10:00:00');
INSERT INTO `fin_voucher` VALUES (22, 1, '付-2026-03-0005', '2026-03-10', '2026-03', 'PAYMENT', '支付上海巨量引擎广告服务有限公司销售费', 'PROJECT_PAYOUT', 19, 'POSTED', 4, '2026-03-10 10:00:00', '2026-03-10 10:00:00');
INSERT INTO `fin_voucher` VALUES (23, 1, '付-2026-03-0006', '2026-03-12', '2026-03', 'PAYMENT', '支付沪上名家水电安装队施工费', 'PROJECT_PAYOUT', 20, 'POSTED', 4, '2026-03-12 10:00:00', '2026-03-12 10:00:00');
INSERT INTO `fin_voucher` VALUES (24, 1, '付-2026-03-0007', '2026-03-15', '2026-03', 'PAYMENT', '支付上海汇智物业管理有限公司管理费', 'PROJECT_PAYOUT', 21, 'POSTED', 4, '2026-03-15 10:00:00', '2026-03-15 10:00:00');
INSERT INTO `fin_voucher` VALUES (25, 1, '付-2026-03-0008', '2026-03-17', '2026-03', 'PAYMENT', '支付欧派家居上海经销商材料款', 'PROJECT_PAYOUT', 22, 'POSTED', 4, '2026-03-17 10:00:00', '2026-03-17 10:00:00');
INSERT INTO `fin_voucher` VALUES (26, 1, '收-2026-03-0009', '2026-03-18', '2026-03', 'RECEIPT', '收到张伟明 · 锦绣路全屋翻新尾款', 'PROJECT_PAYMENT', 4, 'POSTED', 4, '2026-03-18 10:00:00', '2026-03-18 10:00:00');
INSERT INTO `fin_voucher` VALUES (27, 1, '付-2026-03-0010', '2026-03-18', '2026-03', 'PAYMENT', '支付上海素合空间设计工作室设计费', 'PROJECT_PAYOUT', 23, 'POSTED', 4, '2026-03-18 10:00:00', '2026-03-18 10:00:00');
INSERT INTO `fin_voucher` VALUES (28, 1, '付-2026-03-0011', '2026-03-28', '2026-03', 'PAYMENT', '支付沪上名家水电安装队人工费', 'PROJECT_PAYOUT', 24, 'POSTED', 4, '2026-03-28 10:00:00', '2026-03-28 10:00:00');
INSERT INTO `fin_voucher` VALUES (29, 1, '付-2026-04-0001', '2026-04-04', '2026-04', 'PAYMENT', '支付红星美凯龙建材（上海供应处）材料款', 'PROJECT_PAYOUT', 25, 'POSTED', 4, '2026-04-04 10:00:00', '2026-04-04 10:00:00');
INSERT INTO `fin_voucher` VALUES (30, 1, '付-2026-04-0002', '2026-04-15', '2026-04', 'PAYMENT', '支付上海安居物流配送有限公司其他', 'PROJECT_PAYOUT', 26, 'POSTED', 4, '2026-04-15 10:00:00', '2026-04-15 10:00:00');
INSERT INTO `fin_voucher` VALUES (31, 1, '付-2026-04-0003', '2026-04-18', '2026-04', 'PAYMENT', '支付红星美凯龙建材（上海供应处）材料款', 'PROJECT_PAYOUT', 27, 'POSTED', 4, '2026-04-18 10:00:00', '2026-04-18 10:00:00');
INSERT INTO `fin_voucher` VALUES (32, 1, '收-2026-04-0004', '2026-04-20', '2026-04', 'RECEIPT', '收到王丽华 · 漕溪北路新房装修中期款', 'PROJECT_PAYMENT', 5, 'POSTED', 4, '2026-04-20 10:00:00', '2026-04-20 10:00:00');
INSERT INTO `fin_voucher` VALUES (33, 1, '付-2026-05-0001', '2026-05-10', '2026-05', 'PAYMENT', '支付上海汇智物业管理有限公司管理费', 'PROJECT_PAYOUT', 28, 'POSTED', 4, '2026-05-10 10:00:00', '2026-05-10 10:00:00');
INSERT INTO `fin_voucher` VALUES (34, 1, '付-2026-05-0002', '2026-05-20', '2026-05', 'PAYMENT', '支付沪上名家水电安装队人工费', 'PROJECT_PAYOUT', 29, 'POSTED', 4, '2026-05-20 10:00:00', '2026-05-20 10:00:00');
INSERT INTO `fin_voucher` VALUES (35, 1, '付-2026-05-0003', '2026-05-22', '2026-05', 'PAYMENT', '支付红星美凯龙建材（上海供应处）材料款', 'PROJECT_PAYOUT', 30, 'POSTED', 4, '2026-05-22 10:00:00', '2026-05-22 10:00:00');
INSERT INTO `fin_voucher` VALUES (36, 1, '收-2026-06-0001', '2026-06-05', '2026-06', 'RECEIPT', '收到赵磊 · 金桥商铺装修定金', 'PROJECT_PAYMENT', 6, 'POSTED', 4, '2026-06-05 10:00:00', '2026-06-05 10:00:00');
INSERT INTO `fin_voucher` VALUES (37, 1, '付-2026-06-0002', '2026-06-08', '2026-06', 'PAYMENT', '支付上海素合空间设计工作室设计费', 'PROJECT_PAYOUT', 31, 'POSTED', 4, '2026-06-08 10:00:00', '2026-06-08 10:00:00');
INSERT INTO `fin_voucher` VALUES (38, 1, '付-2026-06-0003', '2026-06-08', '2026-06', 'PAYMENT', '支付上海巨量引擎广告服务有限公司销售费', 'PROJECT_PAYOUT', 32, 'POSTED', 4, '2026-06-08 10:00:00', '2026-06-08 10:00:00');
INSERT INTO `fin_voucher` VALUES (39, 1, '付-2026-06-0004', '2026-06-10', '2026-06', 'PAYMENT', '支付沪上名家水电安装队人工费', 'PROJECT_PAYOUT', 33, 'POSTED', 4, '2026-06-10 10:00:00', '2026-06-10 10:00:00');
INSERT INTO `fin_voucher` VALUES (40, 1, '付-2026-06-0005', '2026-06-14', '2026-06', 'PAYMENT', '支付红星美凯龙建材（上海供应处）材料款', 'PROJECT_PAYOUT', 34, 'POSTED', 4, '2026-06-14 10:00:00', '2026-06-14 10:00:00');
INSERT INTO `fin_voucher` VALUES (41, 1, '付-2026-06-0006', '2026-06-17', '2026-06', 'PAYMENT', '支付红星美凯龙建材（上海供应处）材料款', 'PROJECT_PAYOUT', 35, 'POSTED', 4, '2026-06-17 10:00:00', '2026-06-17 10:00:00');
INSERT INTO `fin_voucher` VALUES (42, 1, '付-2026-06-0007', '2026-06-18', '2026-06', 'PAYMENT', '支付上海巨量引擎广告服务有限公司销售费', 'PROJECT_PAYOUT', 36, 'POSTED', 4, '2026-06-18 10:00:00', '2026-06-18 10:00:00');
INSERT INTO `fin_voucher` VALUES (43, 1, '付-2026-06-0008', '2026-06-20', '2026-06', 'PAYMENT', '支付红星美凯龙建材（上海供应处）材料款', 'PROJECT_PAYOUT', 37, 'POSTED', 4, '2026-06-20 10:00:00', '2026-06-20 10:00:00');
INSERT INTO `fin_voucher` VALUES (44, 1, '收-2026-06-0009', '2026-06-25', '2026-06', 'RECEIPT', '收到王丽华 · 漕溪北路新房装修尾款', 'PROJECT_PAYMENT', 7, 'POSTED', 4, '2026-06-25 10:00:00', '2026-06-25 10:00:00');
INSERT INTO `fin_voucher` VALUES (45, 1, '付-2026-06-0010', '2026-06-25', '2026-06', 'PAYMENT', '支付沪上名家水电安装队人工费', 'PROJECT_PAYOUT', 38, 'POSTED', 4, '2026-06-25 10:00:00', '2026-06-25 10:00:00');
INSERT INTO `fin_voucher` VALUES (46, 1, '付-2026-06-0011', '2026-06-28', '2026-06', 'PAYMENT', '支付上海汇智物业管理有限公司管理费', 'PROJECT_PAYOUT', 39, 'POSTED', 4, '2026-06-28 10:00:00', '2026-06-28 10:00:00');
INSERT INTO `fin_voucher` VALUES (47, 1, '付-2026-07-0001', '2026-07-05', '2026-07', 'PAYMENT', '支付上海汇智物业管理有限公司管理费', 'PROJECT_PAYOUT', 40, 'POSTED', 4, '2026-07-05 10:00:00', '2026-07-05 10:00:00');
INSERT INTO `fin_voucher` VALUES (48, 1, '付-2026-07-0002', '2026-07-15', '2026-07', 'PAYMENT', '支付红星美凯龙建材（上海供应处）材料款', 'PROJECT_PAYOUT', 41, 'POSTED', 4, '2026-07-15 10:00:00', '2026-07-15 10:00:00');
INSERT INTO `fin_voucher` VALUES (49, 1, '付-2026-07-0003', '2026-07-20', '2026-07', 'PAYMENT', '支付沪上名家水电安装队人工费', 'PROJECT_PAYOUT', 42, 'POSTED', 4, '2026-07-20 10:00:00', '2026-07-20 10:00:00');
INSERT INTO `fin_voucher` VALUES (50, 1, '付-2026-07-0004', '2026-07-28', '2026-07', 'PAYMENT', '支付上海汇智物业管理有限公司管理费', 'PROJECT_PAYOUT', 43, 'POSTED', 4, '2026-07-28 10:00:00', '2026-07-28 10:00:00');
INSERT INTO `fin_voucher` VALUES (51, 1, '付-2026-08-0001', '2026-08-01', '2026-08', 'PAYMENT', '支付沪上名家水电安装队施工费', 'PROJECT_PAYOUT', 44, 'POSTED', 4, '2026-08-01 10:00:00', '2026-08-01 10:00:00');
INSERT INTO `fin_voucher` VALUES (52, 1, '收-2026-08-0002', '2026-08-05', '2026-08', 'RECEIPT', '收到吴敏 · 恒丰路老房翻新定金', 'PROJECT_PAYMENT', 8, 'POSTED', 4, '2026-08-05 10:00:00', '2026-08-05 10:00:00');
INSERT INTO `fin_voucher` VALUES (53, 1, '付-2026-08-0003', '2026-08-06', '2026-08', 'PAYMENT', '支付上海巨量引擎广告服务有限公司销售费', 'PROJECT_PAYOUT', 45, 'POSTED', 4, '2026-08-06 10:00:00', '2026-08-06 10:00:00');
INSERT INTO `fin_voucher` VALUES (54, 1, '付-2026-08-0004', '2026-08-08', '2026-08', 'PAYMENT', '支付上海素合空间设计工作室设计费', 'PROJECT_PAYOUT', 46, 'POSTED', 4, '2026-08-08 10:00:00', '2026-08-08 10:00:00');
INSERT INTO `fin_voucher` VALUES (55, 1, '付-2026-08-0005', '2026-08-12', '2026-08', 'PAYMENT', '支付东方雨虹防水工程部材料款', 'PROJECT_PAYOUT', 47, 'POSTED', 4, '2026-08-12 10:00:00', '2026-08-12 10:00:00');
INSERT INTO `fin_voucher` VALUES (56, 1, '付-2026-08-0006', '2026-08-14', '2026-08', 'PAYMENT', '支付红星美凯龙建材（上海供应处）材料款', 'PROJECT_PAYOUT', 48, 'POSTED', 4, '2026-08-14 10:00:00', '2026-08-14 10:00:00');
INSERT INTO `fin_voucher` VALUES (57, 1, '付-2026-08-0007', '2026-08-17', '2026-08', 'PAYMENT', '支付欧派家居上海经销商材料款', 'PROJECT_PAYOUT', 49, 'POSTED', 4, '2026-08-17 10:00:00', '2026-08-17 10:00:00');
INSERT INTO `fin_voucher` VALUES (58, 1, '付-2026-08-0008', '2026-08-18', '2026-08', 'PAYMENT', '支付红星美凯龙建材（上海供应处）材料款', 'PROJECT_PAYOUT', 50, 'POSTED', 4, '2026-08-18 10:00:00', '2026-08-18 10:00:00');
INSERT INTO `fin_voucher` VALUES (59, 1, '收-2026-08-0009', '2026-08-20', '2026-08', 'RECEIPT', '收到赵磊 · 金桥商铺装修尾款', 'PROJECT_PAYMENT', 9, 'POSTED', 4, '2026-08-20 10:00:00', '2026-08-20 10:00:00');
INSERT INTO `fin_voucher` VALUES (60, 1, '付-2026-08-0010', '2026-08-28', '2026-08', 'PAYMENT', '支付沪上名家水电安装队人工费', 'PROJECT_PAYOUT', 51, 'POSTED', 4, '2026-08-28 10:00:00', '2026-08-28 10:00:00');
INSERT INTO `fin_voucher` VALUES (61, 1, '收-2026-09-0001', '2026-09-03', '2026-09', 'RECEIPT', '收到陈志强 · 仲益大厦办公室改造定金', 'PROJECT_PAYMENT', 10, 'POSTED', 4, '2026-09-03 10:00:00', '2026-09-03 10:00:00');
INSERT INTO `fin_voucher` VALUES (62, 1, '付-2026-09-0002', '2026-09-05', '2026-09', 'PAYMENT', '支付上海汇智物业管理有限公司管理费', 'PROJECT_PAYOUT', 52, 'POSTED', 4, '2026-09-05 10:00:00', '2026-09-05 10:00:00');
INSERT INTO `fin_voucher` VALUES (63, 1, '付-2026-09-0003', '2026-09-07', '2026-09', 'PAYMENT', '支付红星美凯龙建材（上海供应处）材料款', 'PROJECT_PAYOUT', 53, 'POSTED', 4, '2026-09-07 10:00:00', '2026-09-07 10:00:00');
INSERT INTO `fin_voucher` VALUES (64, 1, '付-2026-09-0004', '2026-09-10', '2026-09', 'PAYMENT', '支付红星美凯龙建材（上海供应处）材料款', 'PROJECT_PAYOUT', 54, 'POSTED', 4, '2026-09-10 10:00:00', '2026-09-10 10:00:00');
INSERT INTO `fin_voucher` VALUES (65, 1, '付-2026-09-0005', '2026-09-10', '2026-09', 'PAYMENT', '支付红星美凯龙建材（上海供应处）材料款', 'PROJECT_PAYOUT', 55, 'POSTED', 4, '2026-09-10 10:00:00', '2026-09-10 10:00:00');
INSERT INTO `fin_voucher` VALUES (66, 1, '付-2026-09-0006', '2026-09-15', '2026-09', 'PAYMENT', '支付上海素合空间设计工作室设计费', 'PROJECT_PAYOUT', 56, 'POSTED', 4, '2026-09-15 10:00:00', '2026-09-15 10:00:00');
INSERT INTO `fin_voucher` VALUES (67, 1, '收-2026-09-0007', '2026-09-18', '2026-09', 'RECEIPT', '收到吴敏 · 恒丰路老房翻新中期款', 'PROJECT_PAYMENT', 11, 'POSTED', 4, '2026-09-18 10:00:00', '2026-09-18 10:00:00');
INSERT INTO `fin_voucher` VALUES (68, 1, '付-2026-09-0008', '2026-09-25', '2026-09', 'PAYMENT', '支付沪上名家水电安装队人工费', 'PROJECT_PAYOUT', 57, 'POSTED', 4, '2026-09-25 10:00:00', '2026-09-25 10:00:00');
INSERT INTO `fin_voucher` VALUES (69, 1, '付-2026-09-0009', '2026-09-28', '2026-09', 'PAYMENT', '支付红星美凯龙建材（上海供应处）材料款', 'PROJECT_PAYOUT', 58, 'POSTED', 4, '2026-09-28 10:00:00', '2026-09-28 10:00:00');
INSERT INTO `fin_voucher` VALUES (70, 1, '付-2026-09-0010', '2026-09-30', '2026-09', 'PAYMENT', '支付立邦涂料上海总代理材料款', 'PROJECT_PAYOUT', 59, 'POSTED', 4, '2026-09-30 10:00:00', '2026-09-30 10:00:00');
INSERT INTO `fin_voucher` VALUES (71, 1, '付-2026-09-0011', '2026-09-30', '2026-09', 'PAYMENT', '支付立邦涂料上海总代理材料款', 'PROJECT_PAYOUT', 60, 'POSTED', 4, '2026-09-30 10:00:00', '2026-09-30 10:00:00');
INSERT INTO `fin_voucher` VALUES (72, 1, '付-2026-09-0012', '2026-09-30', '2026-09', 'PAYMENT', '支付上海汇智物业管理有限公司管理费', 'PROJECT_PAYOUT', 61, 'POSTED', 4, '2026-09-30 10:00:00', '2026-09-30 10:00:00');
INSERT INTO `fin_voucher` VALUES (73, 1, '转-2026-09-0013', '2026-09-30', '2026-09', 'TRANSFER', '现金存入银行（内部转账）', 'MANUAL', 0, 'POSTED', 4, '2026-09-30 10:00:00', '2026-09-30 10:00:00');
INSERT INTO `fin_voucher` VALUES (74, 1, '付-2026-10-0001', '2026-10-06', '2026-10', 'PAYMENT', '支付沪上名家水电安装队人工费', 'PROJECT_PAYOUT', 62, 'POSTED', 4, '2026-10-06 10:00:00', '2026-10-06 10:00:00');
INSERT INTO `fin_voucher` VALUES (75, 1, '付-2026-10-0002', '2026-10-09', '2026-10', 'PAYMENT', '支付上海巨量引擎广告服务有限公司销售费', 'PROJECT_PAYOUT', 63, 'POSTED', 4, '2026-10-09 10:00:00', '2026-10-09 10:00:00');
INSERT INTO `fin_voucher` VALUES (76, 1, '收-2026-10-0003', '2026-10-15', '2026-10', 'RECEIPT', '收到陈志强 · 仲益大厦办公室改造中期款', 'PROJECT_PAYMENT', 12, 'POSTED', 4, '2026-10-15 10:00:00', '2026-10-15 10:00:00');
INSERT INTO `fin_voucher` VALUES (78, 1, '收-2026-10-0005', '2026-10-25', '2026-10', 'RECEIPT', '收到吴敏 · 恒丰路老房翻新尾款', 'PROJECT_PAYMENT', 14, 'POSTED', 4, '2026-10-25 10:00:00', '2026-10-25 10:00:00');
INSERT INTO `fin_voucher` VALUES (79, 1, '付-2026-10-0006', '2026-10-28', '2026-10', 'PAYMENT', '支付沪上名家水电安装队人工费', 'PROJECT_PAYOUT', 64, 'POSTED', 4, '2026-10-28 10:00:00', '2026-10-28 10:00:00');
INSERT INTO `fin_voucher` VALUES (80, 1, '付-2026-11-0001', '2026-11-01', '2026-11', 'PAYMENT', '支付东方雨虹防水工程部材料款', 'PROJECT_PAYOUT', 65, 'POSTED', 4, '2026-11-01 10:00:00', '2026-11-01 10:00:00');
INSERT INTO `fin_voucher` VALUES (81, 1, '收-2026-11-0002', '2026-11-02', '2026-11', 'RECEIPT', '收到黄建军 · 国定路厨卫改造定金', 'PROJECT_PAYMENT', 15, 'POSTED', 4, '2026-11-02 10:00:00', '2026-11-02 10:00:00');
INSERT INTO `fin_voucher` VALUES (82, 1, '付-2026-11-0003', '2026-11-03', '2026-11', 'PAYMENT', '支付上海安居物流配送有限公司其他', 'PROJECT_PAYOUT', 66, 'POSTED', 4, '2026-11-03 10:00:00', '2026-11-03 10:00:00');
INSERT INTO `fin_voucher` VALUES (83, 1, '付-2026-11-0004', '2026-11-03', '2026-11', 'PAYMENT', '支付上海素合空间设计工作室设计费', 'PROJECT_PAYOUT', 67, 'POSTED', 4, '2026-11-03 10:00:00', '2026-11-03 10:00:00');
INSERT INTO `fin_voucher` VALUES (84, 1, '付-2026-11-0005', '2026-11-03', '2026-11', 'PAYMENT', '支付上海安居物流配送有限公司其他', 'PROJECT_PAYOUT', 68, 'POSTED', 4, '2026-11-03 10:00:00', '2026-11-03 10:00:00');
INSERT INTO `fin_voucher` VALUES (87, 1, '收-202611-0006', '2026-11-06', '2026-11', 'RECEIPT', '收到黄建军 · 国定路厨卫改造???', 'PROJECT_PAYMENT', 17, 'VOID', 1, '2026-11-06 16:35:53', '2026-11-06 16:36:37');
INSERT INTO `fin_voucher` VALUES (88, 1, '收-202611-0007', '2026-11-06', '2026-11', 'RECEIPT', '收到黄建军 · 国定路厨卫改造中期款', 'PROJECT_PAYMENT', 18, 'VOID', 1, '2026-11-06 16:36:46', '2026-11-06 16:36:59');
INSERT INTO `fin_voucher` VALUES (89, 1, '收-202611-0008', '2026-11-06', '2026-11', 'RECEIPT', '收到黄建军 · 国定路厨卫改造中期款', 'PROJECT_PAYMENT', 18, 'VOID', 1, '2026-11-06 16:36:59', '2026-11-06 16:37:00');
INSERT INTO `fin_voucher` VALUES (90, 1, '付-202611-0009', '2026-11-06', '2026-11', 'PAYMENT', '支付红星美凯龙建材（上海供应处）材料款', 'PROJECT_PAYOUT', 69, 'VOID', 1, '2026-11-06 16:37:19', '2026-11-06 16:37:19');

-- ----------------------------
-- Table structure for fin_voucher_entry
-- ----------------------------
DROP TABLE IF EXISTS `fin_voucher_entry`;
CREATE TABLE `fin_voucher_entry`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '分录ID',
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `voucher_id` bigint(20) UNSIGNED NOT NULL COMMENT '凭证ID',
  `line_no` smallint(5) UNSIGNED NOT NULL COMMENT '凭证内行号',
  `account_id` bigint(20) UNSIGNED NOT NULL COMMENT '会计科目ID',
  `account_code` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '科目编码快照',
  `account_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '科目名称快照',
  `summary` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '分录摘要',
  `direction_code` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '借贷方向代码：DEBIT/CREDIT',
  `amount` decimal(18, 2) NOT NULL COMMENT '分录金额（正数）',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_fin_voucher_entry_line`(`voucher_id`, `line_no`) USING BTREE,
  INDEX `idx_fin_voucher_entry_account`(`account_id`) USING BTREE,
  CONSTRAINT `fk_fin_voucher_entry_account` FOREIGN KEY (`account_id`) REFERENCES `fin_account` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_fin_voucher_entry_voucher` FOREIGN KEY (`voucher_id`) REFERENCES `fin_voucher` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 665 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '会计凭证借贷分录明细表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of fin_voucher_entry
-- ----------------------------
INSERT INTO `fin_voucher_entry` VALUES (489, 1, 1, 1, 170, '5602', '管理费用', '支付上海汇智物业管理有限公司管理费', 'DEBIT', 8000.00, '2026-01-05 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (490, 1, 1, 2, 171, '1001', '库存现金', '资金流出', 'CREDIT', 8000.00, '2026-01-05 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (491, 1, 2, 1, 172, '1002', '银行存款', '收到张伟明 · 锦绣路全屋翻新定金', 'DEBIT', 37200.00, '2026-01-12 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (492, 1, 2, 2, 167, '5001', '主营业务收入', '确认主营业务收入', 'CREDIT', 37200.00, '2026-01-12 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (493, 1, 3, 1, 168, '5401', '主营业务成本', '支付红星美凯龙建材（上海供应处）材料款', 'DEBIT', 5120.00, '2026-01-13 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (494, 1, 3, 2, 171, '1001', '库存现金', '资金流出', 'CREDIT', 5120.00, '2026-01-13 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (495, 1, 4, 1, 168, '5401', '主营业务成本', '支付红星美凯龙建材（上海供应处）材料款', 'DEBIT', 2520.00, '2026-01-14 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (496, 1, 4, 2, 171, '1001', '库存现金', '资金流出', 'CREDIT', 2520.00, '2026-01-14 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (497, 1, 5, 1, 168, '5401', '主营业务成本', '支付立邦涂料上海总代理材料款', 'DEBIT', 4416.00, '2026-01-16 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (498, 1, 5, 2, 171, '1001', '库存现金', '资金流出', 'CREDIT', 4416.00, '2026-01-16 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (499, 1, 6, 1, 168, '5401', '主营业务成本', '支付红星美凯龙建材（上海供应处）材料款', 'DEBIT', 1360.00, '2026-01-18 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (500, 1, 6, 2, 171, '1001', '库存现金', '资金流出', 'CREDIT', 1360.00, '2026-01-18 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (501, 1, 7, 1, 168, '5401', '主营业务成本', '支付红星美凯龙建材（上海供应处）材料款', 'DEBIT', 1280.00, '2026-01-18 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (502, 1, 7, 2, 171, '1001', '库存现金', '资金流出', 'CREDIT', 1280.00, '2026-01-18 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (503, 1, 8, 1, 168, '5401', '主营业务成本', '支付红星美凯龙建材（上海供应处）材料款', 'DEBIT', 10000.00, '2026-01-20 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (504, 1, 8, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 10000.00, '2026-01-20 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (505, 1, 9, 1, 168, '5401', '主营业务成本', '支付上海安居物流配送有限公司其他', 'DEBIT', 5000.00, '2026-01-25 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (506, 1, 9, 2, 171, '1001', '库存现金', '资金流出', 'CREDIT', 5000.00, '2026-01-25 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (507, 1, 10, 1, 170, '5602', '管理费用', '支付上海汇智物业管理有限公司管理费', 'DEBIT', 8000.00, '2026-02-05 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (508, 1, 10, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 8000.00, '2026-02-05 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (509, 1, 11, 1, 168, '5401', '主营业务成本', '支付沪上名家水电安装队人工费', 'DEBIT', 15000.00, '2026-02-10 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (510, 1, 11, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 15000.00, '2026-02-10 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (511, 1, 12, 1, 172, '1002', '银行存款', '收到张伟明 · 锦绣路全屋翻新中期款', 'DEBIT', 93000.00, '2026-02-15 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (512, 1, 12, 2, 167, '5001', '主营业务收入', '确认主营业务收入', 'CREDIT', 93000.00, '2026-02-15 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (513, 1, 13, 1, 172, '1002', '银行存款', '收到王丽华 · 漕溪北路新房装修定金', 'DEBIT', 53600.00, '2026-02-18 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (514, 1, 13, 2, 167, '5001', '主营业务收入', '确认主营业务收入', 'CREDIT', 53600.00, '2026-02-18 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (515, 1, 14, 1, 168, '5401', '主营业务成本', '支付红星美凯龙建材（上海供应处）材料款', 'DEBIT', 9440.00, '2026-02-22 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (516, 1, 14, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 9440.00, '2026-02-22 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (517, 1, 15, 1, 168, '5401', '主营业务成本', '支付上海素合空间设计工作室设计费', 'DEBIT', 16000.00, '2026-02-25 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (518, 1, 15, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 16000.00, '2026-02-25 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (519, 1, 16, 1, 168, '5401', '主营业务成本', '支付欧派家居上海经销商材料款', 'DEBIT', 12000.00, '2026-02-28 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (520, 1, 16, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 12000.00, '2026-02-28 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (521, 1, 17, 1, 168, '5401', '主营业务成本', '支付沪上名家水电安装队人工费', 'DEBIT', 10600.00, '2026-02-28 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (522, 1, 17, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 10600.00, '2026-02-28 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (523, 1, 18, 1, 168, '5401', '主营业务成本', '支付立邦涂料上海总代理材料款', 'DEBIT', 5200.00, '2026-03-03 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (524, 1, 18, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 5200.00, '2026-03-03 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (525, 1, 19, 1, 168, '5401', '主营业务成本', '支付红星美凯龙建材（上海供应处）材料款', 'DEBIT', 10000.00, '2026-03-05 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (526, 1, 19, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 10000.00, '2026-03-05 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (527, 1, 20, 1, 168, '5401', '主营业务成本', '支付红星美凯龙建材（上海供应处）材料款', 'DEBIT', 1800.00, '2026-03-07 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (528, 1, 20, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 1800.00, '2026-03-07 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (529, 1, 21, 1, 168, '5401', '主营业务成本', '支付红星美凯龙建材（上海供应处）材料款', 'DEBIT', 30000.00, '2026-03-08 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (530, 1, 21, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 30000.00, '2026-03-08 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (531, 1, 22, 1, 169, '5601', '销售费用', '支付上海巨量引擎广告服务有限公司销售费', 'DEBIT', 5000.00, '2026-03-10 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (532, 1, 22, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 5000.00, '2026-03-10 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (533, 1, 23, 1, 168, '5401', '主营业务成本', '支付沪上名家水电安装队施工费', 'DEBIT', 12400.00, '2026-03-12 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (534, 1, 23, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 12400.00, '2026-03-12 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (535, 1, 24, 1, 170, '5602', '管理费用', '支付上海汇智物业管理有限公司管理费', 'DEBIT', 6000.00, '2026-03-15 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (536, 1, 24, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 6000.00, '2026-03-15 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (537, 1, 25, 1, 168, '5401', '主营业务成本', '支付欧派家居上海经销商材料款', 'DEBIT', 8040.00, '2026-03-17 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (538, 1, 25, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 8040.00, '2026-03-17 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (539, 1, 26, 1, 172, '1002', '银行存款', '收到张伟明 · 锦绣路全屋翻新尾款', 'DEBIT', 55800.00, '2026-03-18 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (540, 1, 26, 2, 167, '5001', '主营业务收入', '确认主营业务收入', 'CREDIT', 55800.00, '2026-03-18 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (541, 1, 27, 1, 168, '5401', '主营业务成本', '支付上海素合空间设计工作室设计费', 'DEBIT', 9000.00, '2026-03-18 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (542, 1, 27, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 9000.00, '2026-03-18 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (543, 1, 28, 1, 168, '5401', '主营业务成本', '支付沪上名家水电安装队人工费', 'DEBIT', 22000.00, '2026-03-28 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (544, 1, 28, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 22000.00, '2026-03-28 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (545, 1, 29, 1, 168, '5401', '主营业务成本', '支付红星美凯龙建材（上海供应处）材料款', 'DEBIT', 640.00, '2026-04-04 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (546, 1, 29, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 640.00, '2026-04-04 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (547, 1, 30, 1, 168, '5401', '主营业务成本', '支付上海安居物流配送有限公司其他', 'DEBIT', 1200.00, '2026-04-15 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (548, 1, 30, 2, 171, '1001', '库存现金', '资金流出', 'CREDIT', 1200.00, '2026-04-15 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (549, 1, 31, 1, 168, '5401', '主营业务成本', '支付红星美凯龙建材（上海供应处）材料款', 'DEBIT', 22000.00, '2026-04-18 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (550, 1, 31, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 22000.00, '2026-04-18 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (551, 1, 32, 1, 172, '1002', '银行存款', '收到王丽华 · 漕溪北路新房装修中期款', 'DEBIT', 134000.00, '2026-04-20 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (552, 1, 32, 2, 167, '5001', '主营业务收入', '确认主营业务收入', 'CREDIT', 134000.00, '2026-04-20 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (553, 1, 33, 1, 170, '5602', '管理费用', '支付上海汇智物业管理有限公司管理费', 'DEBIT', 12000.00, '2026-05-10 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (554, 1, 33, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 12000.00, '2026-05-10 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (555, 1, 34, 1, 168, '5401', '主营业务成本', '支付沪上名家水电安装队人工费', 'DEBIT', 18000.00, '2026-05-20 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (556, 1, 34, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 18000.00, '2026-05-20 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (557, 1, 35, 1, 168, '5401', '主营业务成本', '支付红星美凯龙建材（上海供应处）材料款', 'DEBIT', 18000.00, '2026-05-22 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (558, 1, 35, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 18000.00, '2026-05-22 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (559, 1, 36, 1, 172, '1002', '银行存款', '收到赵磊 · 金桥商铺装修定金', 'DEBIT', 22800.00, '2026-06-05 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (560, 1, 36, 2, 167, '5001', '主营业务收入', '确认主营业务收入', 'CREDIT', 22800.00, '2026-06-05 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (561, 1, 37, 1, 168, '5401', '主营业务成本', '支付上海素合空间设计工作室设计费', 'DEBIT', 4000.00, '2026-06-08 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (562, 1, 37, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 4000.00, '2026-06-08 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (563, 1, 38, 1, 169, '5601', '销售费用', '支付上海巨量引擎广告服务有限公司销售费', 'DEBIT', 8000.00, '2026-06-08 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (564, 1, 38, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 8000.00, '2026-06-08 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (565, 1, 39, 1, 168, '5401', '主营业务成本', '支付沪上名家水电安装队人工费', 'DEBIT', 12000.00, '2026-06-10 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (566, 1, 39, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 12000.00, '2026-06-10 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (567, 1, 40, 1, 168, '5401', '主营业务成本', '支付红星美凯龙建材（上海供应处）材料款', 'DEBIT', 4050.00, '2026-06-14 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (568, 1, 40, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 4050.00, '2026-06-14 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (569, 1, 41, 1, 168, '5401', '主营业务成本', '支付红星美凯龙建材（上海供应处）材料款', 'DEBIT', 5100.00, '2026-06-17 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (570, 1, 41, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 5100.00, '2026-06-17 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (571, 1, 42, 1, 169, '5601', '销售费用', '支付上海巨量引擎广告服务有限公司销售费', 'DEBIT', 6000.00, '2026-06-18 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (572, 1, 42, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 6000.00, '2026-06-18 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (573, 1, 43, 1, 168, '5401', '主营业务成本', '支付红星美凯龙建材（上海供应处）材料款', 'DEBIT', 7280.00, '2026-06-20 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (574, 1, 43, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 7280.00, '2026-06-20 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (575, 1, 44, 1, 172, '1002', '银行存款', '收到王丽华 · 漕溪北路新房装修尾款', 'DEBIT', 53600.00, '2026-06-25 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (576, 1, 44, 2, 167, '5001', '主营业务收入', '确认主营业务收入', 'CREDIT', 53600.00, '2026-06-25 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (577, 1, 45, 1, 168, '5401', '主营业务成本', '支付沪上名家水电安装队人工费', 'DEBIT', 10000.00, '2026-06-25 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (578, 1, 45, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 10000.00, '2026-06-25 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (579, 1, 46, 1, 170, '5602', '管理费用', '支付上海汇智物业管理有限公司管理费', 'DEBIT', 8000.00, '2026-06-28 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (580, 1, 46, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 8000.00, '2026-06-28 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (581, 1, 47, 1, 170, '5602', '管理费用', '支付上海汇智物业管理有限公司管理费', 'DEBIT', 8000.00, '2026-07-05 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (582, 1, 47, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 8000.00, '2026-07-05 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (583, 1, 48, 1, 168, '5401', '主营业务成本', '支付红星美凯龙建材（上海供应处）材料款', 'DEBIT', 6000.00, '2026-07-15 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (584, 1, 48, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 6000.00, '2026-07-15 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (585, 1, 49, 1, 168, '5401', '主营业务成本', '支付沪上名家水电安装队人工费', 'DEBIT', 8000.00, '2026-07-20 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (586, 1, 49, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 8000.00, '2026-07-20 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (587, 1, 50, 1, 170, '5602', '管理费用', '支付上海汇智物业管理有限公司管理费', 'DEBIT', 2000.00, '2026-07-28 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (588, 1, 50, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 2000.00, '2026-07-28 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (589, 1, 51, 1, 168, '5401', '主营业务成本', '支付沪上名家水电安装队施工费', 'DEBIT', 4000.00, '2026-08-01 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (590, 1, 51, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 4000.00, '2026-08-01 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (591, 1, 52, 1, 172, '1002', '银行存款', '收到吴敏 · 恒丰路老房翻新定金', 'DEBIT', 23600.00, '2026-08-05 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (592, 1, 52, 2, 167, '5001', '主营业务收入', '确认主营业务收入', 'CREDIT', 23600.00, '2026-08-05 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (593, 1, 53, 1, 169, '5601', '销售费用', '支付上海巨量引擎广告服务有限公司销售费', 'DEBIT', 6000.00, '2026-08-06 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (594, 1, 53, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 6000.00, '2026-08-06 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (595, 1, 54, 1, 168, '5401', '主营业务成本', '支付上海素合空间设计工作室设计费', 'DEBIT', 8000.00, '2026-08-08 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (596, 1, 54, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 8000.00, '2026-08-08 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (597, 1, 55, 1, 168, '5401', '主营业务成本', '支付东方雨虹防水工程部材料款', 'DEBIT', 3800.00, '2026-08-12 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (598, 1, 55, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 3800.00, '2026-08-12 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (599, 1, 56, 1, 168, '5401', '主营业务成本', '支付红星美凯龙建材（上海供应处）材料款', 'DEBIT', 8160.00, '2026-08-14 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (600, 1, 56, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 8160.00, '2026-08-14 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (601, 1, 57, 1, 168, '5401', '主营业务成本', '支付欧派家居上海经销商材料款', 'DEBIT', 5120.00, '2026-08-17 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (602, 1, 57, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 5120.00, '2026-08-17 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (603, 1, 58, 1, 168, '5401', '主营业务成本', '支付红星美凯龙建材（上海供应处）材料款', 'DEBIT', 1680.00, '2026-08-18 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (604, 1, 58, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 1680.00, '2026-08-18 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (605, 1, 59, 1, 172, '1002', '银行存款', '收到赵磊 · 金桥商铺装修尾款', 'DEBIT', 30400.00, '2026-08-20 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (606, 1, 59, 2, 167, '5001', '主营业务收入', '确认主营业务收入', 'CREDIT', 30400.00, '2026-08-20 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (607, 1, 60, 1, 168, '5401', '主营业务成本', '支付沪上名家水电安装队人工费', 'DEBIT', 12000.00, '2026-08-28 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (608, 1, 60, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 12000.00, '2026-08-28 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (609, 1, 61, 1, 172, '1002', '银行存款', '收到陈志强 · 仲益大厦办公室改造定金', 'DEBIT', 28800.00, '2026-09-03 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (610, 1, 61, 2, 167, '5001', '主营业务收入', '确认主营业务收入', 'CREDIT', 28800.00, '2026-09-03 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (611, 1, 62, 1, 170, '5602', '管理费用', '支付上海汇智物业管理有限公司管理费', 'DEBIT', 8000.00, '2026-09-05 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (612, 1, 62, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 8000.00, '2026-09-05 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (613, 1, 63, 1, 168, '5401', '主营业务成本', '支付红星美凯龙建材（上海供应处）材料款', 'DEBIT', 2400.00, '2026-09-07 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (614, 1, 63, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 2400.00, '2026-09-07 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (615, 1, 64, 1, 168, '5401', '主营业务成本', '支付红星美凯龙建材（上海供应处）材料款', 'DEBIT', 5400.00, '2026-09-10 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (616, 1, 64, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 5400.00, '2026-09-10 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (617, 1, 65, 1, 168, '5401', '主营业务成本', '支付红星美凯龙建材（上海供应处）材料款', 'DEBIT', 7000.00, '2026-09-10 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (618, 1, 65, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 7000.00, '2026-09-10 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (619, 1, 66, 1, 168, '5401', '主营业务成本', '支付上海素合空间设计工作室设计费', 'DEBIT', 5000.00, '2026-09-15 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (620, 1, 66, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 5000.00, '2026-09-15 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (621, 1, 67, 1, 172, '1002', '银行存款', '收到吴敏 · 恒丰路老房翻新中期款', 'DEBIT', 47200.00, '2026-09-18 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (622, 1, 67, 2, 167, '5001', '主营业务收入', '确认主营业务收入', 'CREDIT', 47200.00, '2026-09-18 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (623, 1, 68, 1, 168, '5401', '主营业务成本', '支付沪上名家水电安装队人工费', 'DEBIT', 8000.00, '2026-09-25 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (624, 1, 68, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 8000.00, '2026-09-25 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (625, 1, 69, 1, 168, '5401', '主营业务成本', '支付红星美凯龙建材（上海供应处）材料款', 'DEBIT', 19800.00, '2026-09-28 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (626, 1, 69, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 19800.00, '2026-09-28 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (627, 1, 70, 1, 168, '5401', '主营业务成本', '支付立邦涂料上海总代理材料款', 'DEBIT', 3120.00, '2026-09-30 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (628, 1, 70, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 3120.00, '2026-09-30 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (629, 1, 71, 1, 168, '5401', '主营业务成本', '支付立邦涂料上海总代理材料款', 'DEBIT', 3120.00, '2026-09-30 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (630, 1, 71, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 3120.00, '2026-09-30 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (631, 1, 72, 1, 170, '5602', '管理费用', '支付上海汇智物业管理有限公司管理费', 'DEBIT', 4000.00, '2026-09-30 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (632, 1, 72, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 4000.00, '2026-09-30 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (633, 1, 73, 1, 172, '1002', '银行存款', '现金存入银行（内部转账）', 'DEBIT', 10000.00, '2026-09-30 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (634, 1, 73, 2, 171, '1001', '库存现金', '资金流出', 'CREDIT', 10000.00, '2026-09-30 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (635, 1, 74, 1, 168, '5401', '主营业务成本', '支付沪上名家水电安装队人工费', 'DEBIT', 12000.00, '2026-10-06 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (636, 1, 74, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 12000.00, '2026-10-06 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (637, 1, 75, 1, 169, '5601', '销售费用', '支付上海巨量引擎广告服务有限公司销售费', 'DEBIT', 4000.00, '2026-10-09 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (638, 1, 75, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 4000.00, '2026-10-09 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (639, 1, 76, 1, 172, '1002', '银行存款', '收到陈志强 · 仲益大厦办公室改造中期款', 'DEBIT', 38400.00, '2026-10-15 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (640, 1, 76, 2, 167, '5001', '主营业务收入', '确认主营业务收入', 'CREDIT', 38400.00, '2026-10-15 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (643, 1, 78, 1, 171, '1001', '库存现金', '收到吴敏 · 恒丰路老房翻新尾款', 'DEBIT', 35400.00, '2026-10-25 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (644, 1, 78, 2, 167, '5001', '主营业务收入', '确认主营业务收入', 'CREDIT', 35400.00, '2026-10-25 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (645, 1, 79, 1, 168, '5401', '主营业务成本', '支付沪上名家水电安装队人工费', 'DEBIT', 8000.00, '2026-10-28 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (646, 1, 79, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 8000.00, '2026-10-28 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (647, 1, 80, 1, 168, '5401', '主营业务成本', '支付东方雨虹防水工程部材料款', 'DEBIT', 3040.00, '2026-11-01 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (648, 1, 80, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 3040.00, '2026-11-01 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (649, 1, 81, 1, 171, '1001', '库存现金', '收到黄建军 · 国定路厨卫改造定金', 'DEBIT', 20400.00, '2026-11-02 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (650, 1, 81, 2, 167, '5001', '主营业务收入', '确认主营业务收入', 'CREDIT', 20400.00, '2026-11-02 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (651, 1, 82, 1, 168, '5401', '主营业务成本', '支付上海安居物流配送有限公司其他', 'DEBIT', 2000.00, '2026-11-03 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (652, 1, 82, 2, 171, '1001', '库存现金', '资金流出', 'CREDIT', 2000.00, '2026-11-03 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (653, 1, 83, 1, 168, '5401', '主营业务成本', '支付上海素合空间设计工作室设计费', 'DEBIT', 3000.00, '2026-11-03 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (654, 1, 83, 2, 172, '1002', '银行存款', '资金流出', 'CREDIT', 3000.00, '2026-11-03 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (655, 1, 84, 1, 168, '5401', '主营业务成本', '支付上海安居物流配送有限公司其他', 'DEBIT', 2600.00, '2026-11-03 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (656, 1, 84, 2, 171, '1001', '库存现金', '资金流出', 'CREDIT', 2600.00, '2026-11-03 10:00:00');
INSERT INTO `fin_voucher_entry` VALUES (657, 1, 87, 1, 172, '1002', '银行存款', '收到黄建军 · 国定路厨卫改造???', 'DEBIT', 3000.00, '2026-11-06 16:35:53');
INSERT INTO `fin_voucher_entry` VALUES (658, 1, 87, 2, 167, '5001', '主营业务收入', '确认项目收入', 'CREDIT', 3000.00, '2026-11-06 16:35:53');
INSERT INTO `fin_voucher_entry` VALUES (659, 1, 88, 1, 172, '1002', '银行存款', '收到黄建军 · 国定路厨卫改造中期款', 'DEBIT', 3000.00, '2026-11-06 16:36:46');
INSERT INTO `fin_voucher_entry` VALUES (660, 1, 88, 2, 167, '5001', '主营业务收入', '确认项目收入', 'CREDIT', 3000.00, '2026-11-06 16:36:46');
INSERT INTO `fin_voucher_entry` VALUES (661, 1, 89, 1, 172, '1002', '银行存款', '收到黄建军 · 国定路厨卫改造中期款', 'DEBIT', 3800.00, '2026-11-06 16:36:59');
INSERT INTO `fin_voucher_entry` VALUES (662, 1, 89, 2, 167, '5001', '主营业务收入', '确认项目收入', 'CREDIT', 3800.00, '2026-11-06 16:36:59');
INSERT INTO `fin_voucher_entry` VALUES (663, 1, 90, 1, 168, '5401', '主营业务成本', '支付红星美凯龙建材（上海供应处）材料款', 'DEBIT', 500.00, '2026-11-06 16:37:19');
INSERT INTO `fin_voucher_entry` VALUES (664, 1, 90, 2, 172, '1002', '银行存款', '资金流出：公司基本户', 'CREDIT', 500.00, '2026-11-06 16:37:19');

-- ----------------------------
-- Table structure for inventory_record
-- ----------------------------
DROP TABLE IF EXISTS `inventory_record`;
CREATE TABLE `inventory_record`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `project_id` bigint(20) NOT NULL COMMENT '关联项目ID',
  `material_id` bigint(20) NOT NULL COMMENT '关联材料ID',
  `warehouse_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '仓库ID',
  `record_type` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'IN' COMMENT 'IN=入库 OUT=出库',
  `biz_type` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'NORMAL' COMMENT 'NORMAL=业务单据 ADJUST=补录',
  `quantity` decimal(12, 2) NOT NULL DEFAULT 0.00 COMMENT '数量',
  `unit_price` decimal(12, 2) NOT NULL DEFAULT 0.00 COMMENT '单价',
  `amount` decimal(12, 2) NOT NULL DEFAULT 0.00 COMMENT '金额',
  `vendor_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '供应商ID（入库生成付款单时记录）',
  `fund_account_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '资金账户ID',
  `payout_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '关联付款单ID',
  `payment_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '关联收款单ID',
  `source_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '来源单据类型 PURCHASE=采购单 USAGE=领用单',
  `source_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '来源单据ID',
  `biz_date` date NOT NULL COMMENT '业务日期',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  `created_by` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建人ID',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间，NULL表示未删除',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_record_company`(`company_id`, `record_type`) USING BTREE,
  INDEX `idx_record_project`(`project_id`) USING BTREE,
  INDEX `idx_record_material`(`material_id`) USING BTREE,
  INDEX `idx_record_source`(`source_type`, `source_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '出入库流水记录' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of inventory_record
-- ----------------------------
INSERT INTO `inventory_record` VALUES (1, 1, 1, 1, 1, 'IN', 'NORMAL', 10.00, 128.00, 1280.00, 105, 25, 6, 0, '', 0, '2026-01-18', '补货入库，已生成付款单', 3, '2026-01-18 11:00:00', '2026-01-18 11:00:00', NULL);
INSERT INTO `inventory_record` VALUES (2, 1, 5, 24, 1, 'IN', 'NORMAL', 6.00, 520.00, 3120.00, 109, 25, 60, 0, '', 0, '2026-09-30', '乳胶漆入库，已生成付款单', 4, '2026-09-30 11:00:00', '2026-09-30 11:00:00', NULL);
INSERT INTO `inventory_record` VALUES (3, 1, 3, 12, 1, 'OUT', 'NORMAL', 10.00, 165.00, 1650.00, 0, 25, 0, 0, '', 0, '2026-10-20', '调拨出库给业主（项目领用，不生成财务单据）', 2, '2026-10-20 11:00:00', '2026-11-06 16:18:01', NULL);
INSERT INTO `inventory_record` VALUES (4, 1, 1, 3, 1, 'IN', 'ADJUST', 20.00, 8.50, 170.00, 0, 0, 0, 0, '', 0, '2026-02-02', '盘点补录：供应商赠送合页', 3, '2026-02-02 11:00:00', '2026-02-02 11:00:00', NULL);
INSERT INTO `inventory_record` VALUES (5, 1, 1, 4, 1, 'OUT', 'ADJUST', 5.00, 42.00, 210.00, 0, 0, 0, 0, '', 0, '2026-02-05', '盘亏调整：现场破损 5 张', 2, '2026-02-05 11:00:00', '2026-02-05 11:00:00', NULL);
INSERT INTO `inventory_record` VALUES (6, 1, 2, 8, 1, 'OUT', 'ADJUST', 4.00, 45.00, 180.00, 0, 0, 0, 0, '', 0, '2026-05-08', '盘亏调整：受潮报废 4 袋', 3, '2026-05-08 11:00:00', '2026-05-08 11:00:00', NULL);
INSERT INTO `inventory_record` VALUES (7, 1, 2, 9, 1, 'IN', 'ADJUST', 50.00, 32.00, 1600.00, 0, 0, 0, 0, '', 0, '2026-04-01', '盘点补录：前期漏录采购', 4, '2026-04-01 11:00:00', '2026-04-01 11:00:00', NULL);
INSERT INTO `inventory_record` VALUES (8, 1, 5, 22, 1, 'IN', 'ADJUST', 2.00, 380.00, 760.00, 0, 0, 0, 0, '', 0, '2026-09-02', '盘点补录：供应商补发 2 桶', 4, '2026-09-02 11:00:00', '2026-09-02 11:00:00', NULL);
INSERT INTO `inventory_record` VALUES (9, 1, 6, 26, 1, 'OUT', 'ADJUST', 2.00, 260.00, 520.00, 0, 0, 0, 0, '', 0, '2026-08-10', '盘亏调整：运输破损 2㎡', 2, '2026-08-10 11:00:00', '2026-08-10 11:00:00', NULL);
INSERT INTO `inventory_record` VALUES (10, 1, 6, 28, 1, 'IN', 'ADJUST', 40.00, 45.00, 1800.00, 0, 0, 0, 0, '', 0, '2026-07-01', '盘点补录：灯具到货未录入', 2, '2026-07-01 11:00:00', '2026-07-01 11:00:00', NULL);

-- ----------------------------
-- Table structure for material_catalog
-- ----------------------------
DROP TABLE IF EXISTS `material_catalog`;
CREATE TABLE `material_catalog`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '材料名称',
  `spec` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '规格 / 型号',
  `unit` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '单位',
  `unit_price` decimal(12, 2) NOT NULL DEFAULT 0.00 COMMENT '默认单价',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_material_catalog_company_name_spec_unit`(`company_id`, `name`, `spec`, `unit`) USING BTREE,
  INDEX `idx_material_catalog_deleted`(`company_id`, `deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 63 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '公司级材料档案' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of material_catalog
-- ----------------------------
INSERT INTO `material_catalog` VALUES (32, 1, 'LED 筒灯', '7W 中性光', '个', 45.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (33, 1, '不锈钢合页', '4寸 缓冲', '个', 8.50, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (34, 1, '乳胶漆', '18L', '桶', 520.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (35, 1, '乳胶漆', '立邦净味120 18L', '桶', 520.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (36, 1, '前台灯箱', '定制 2.4m×0.8m', '个', 3600.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (37, 1, '办公地毯', '方块拼接地毯', '㎡', 85.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (38, 1, '多层实木板', '18mm E1 级', '张', 128.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (39, 1, '实木复合门', '标准 樘', '樘', 1280.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (40, 1, '开关插座', '西门子 五孔', '个', 32.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (41, 1, '强化复合地板', '1210×165×12mm', '㎡', 118.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (42, 1, '整体橱柜', '石英石台面', '延米', 1800.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (43, 1, '木纹铝方通', '30×80mm', '米', 45.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (44, 1, '橱柜板材', '爱格板 18mm', '张', 268.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (45, 1, '水性木器漆', '哑光白 5L', '桶', 368.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (46, 1, '水龙头', '摩恩 单孔', '个', 460.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (47, 1, '瓷砖', '600×600', '片', 68.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (48, 1, '瓷砖', '600×600 防滑', '片', 68.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (49, 1, '矿棉板', '600×600', '块', 18.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (50, 1, '纸面石膏板', '9mm', '张', 42.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (51, 1, '耐水腻子粉', '20kg', '袋', 45.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (52, 1, '轻钢龙骨', '75 型', '米', 12.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (53, 1, '钢化玻璃隔断', '10mm 钢化', '㎡', 260.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (54, 1, '铝扣板', '300×300', '㎡', 85.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (55, 1, '门套线', '实木 6cm', '米', 28.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (56, 1, '防水涂料', '20kg', '桶', 380.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (57, 1, '防水涂料', '东方雨虹 20kg', '桶', 380.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);
INSERT INTO `material_catalog` VALUES (58, 1, '防静电地板', '600×600', '块', 165.00, '', '2026-11-06 16:08:18', '2026-11-06 16:08:18', NULL);

-- ----------------------------
-- Table structure for print_setting
-- ----------------------------
DROP TABLE IF EXISTS `print_setting`;
CREATE TABLE `print_setting`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `company_id` int(10) UNSIGNED NOT NULL COMMENT '所属公司ID',
  `title` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '抬头标题',
  `sub_title` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '副标题',
  `footer` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '页脚文字',
  `show_date` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否显示打印日期',
  `show_seal` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否加盖电子章',
  `seal_image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '印章图片相对路径',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_print_setting_company`(`company_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '打印设置' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of print_setting
-- ----------------------------
INSERT INTO `print_setting` VALUES (1, 1, '亿民家装饰工程有限公司', '工程报价单 / 施工合同 / 收款凭证', '联系电话：021-6288 6688　地址：上海市静安区江宁路 200 号', 1, 1, '', '2026-01-02 09:00:00', '2026-01-02 09:00:00');

-- ----------------------------
-- Table structure for project
-- ----------------------------
DROP TABLE IF EXISTS `project`;
CREATE TABLE `project`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `customer_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '客户ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '项目名称',
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '施工地址',
  `manager_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '负责人ID（site_user.id）',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态：1=待开工 2=施工中 3=已完工 4=已结算',
  `amount` decimal(10, 2) NOT NULL DEFAULT 0.00 COMMENT '合同金额',
  `start_date` date NULL DEFAULT NULL COMMENT '开工日期',
  `end_date` date NULL DEFAULT NULL COMMENT '预计完工日期',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_customer`(`customer_id`) USING BTREE,
  INDEX `idx_manager`(`manager_id`) USING BTREE,
  INDEX `idx_status`(`status`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 9 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '装修项目表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of project
-- ----------------------------
INSERT INTO `project` VALUES (1, 1, 1, '张伟明 · 锦绣路全屋翻新', '徐汇区锦绣路88弄12号201室', 2, 4, 186000.00, '2026-01-10', '2026-03-20', '含水电改造、柜体定制与全屋油漆', '2026-01-05 09:00:00', '2026-03-20 17:30:00', NULL);
INSERT INTO `project` VALUES (2, 1, 2, '王丽华 · 漕溪北路新房装修', '徐汇区漕溪北路518弄6号1801室', 3, 3, 268000.00, '2026-02-15', '2026-06-30', '三房两厅，环保材料，含定制柜体', '2026-02-08 10:30:00', '2026-06-30 18:00:00', NULL);
INSERT INTO `project` VALUES (3, 1, 3, '陈志强 · 仲益大厦办公室改造', '静安区仲益大厦A座706室', 2, 2, 96000.00, '2026-09-01', '2026-12-15', '办公区 320㎡，夜间施工，含前台形象墙', '2026-08-22 14:00:00', '2026-11-03 09:20:00', NULL);
INSERT INTO `project` VALUES (4, 1, 4, '黄建军 · 国定路厨卫改造', '杨浦区国定路335号2号楼1102室', 5, 1, 68000.00, '2026-11-15', '2027-01-10', '厨房整体翻新 + 卫生间防水改造', '2026-10-28 16:10:00', '2026-11-04 10:00:00', NULL);
INSERT INTO `project` VALUES (5, 1, 5, '吴敏 · 恒丰路老房翻新', '静安区恒丰路600号22层', 4, 2, 118000.00, '2026-08-01', '2026-10-30', '老房全屋翻新，注意成品保护与防滑', '2026-07-24 09:40:00', '2026-11-02 17:00:00', NULL);
INSERT INTO `project` VALUES (6, 1, 6, '赵磊 · 金桥商铺装修', '浦东新区金桥路939号1层商铺', 2, 3, 76000.00, '2026-06-01', '2026-08-15', '临街商铺，含门头招牌与玻璃隔断', '2026-05-22 11:20:00', '2026-08-15 16:30:00', NULL);

-- ----------------------------
-- Table structure for project_change
-- ----------------------------
DROP TABLE IF EXISTS `project_change`;
CREATE TABLE `project_change`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `project_id` bigint(20) NOT NULL COMMENT '关联项目ID',
  `change_type` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '变更类型：ADD 增项 / REDUCE 减项',
  `title` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '变更事项',
  `amount` decimal(12, 2) NOT NULL DEFAULT 0.00 COMMENT '变更金额（正数）',
  `change_date` date NOT NULL COMMENT '变更日期',
  `reason` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '变更原因',
  `images` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '附件图片相对路径，逗号分隔',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间，NULL表示未删除',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_change_project`(`project_id`) USING BTREE,
  INDEX `idx_change_date`(`change_date`) USING BTREE,
  INDEX `idx_change_deleted`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '项目变更记录（增项/减项）' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of project_change
-- ----------------------------
INSERT INTO `project_change` VALUES (1, 1, 1, 'ADD', '厨房增加嵌入式洗碗机及水电点位', 4800.00, '2026-02-06', '业主中途增加电器，需调整橱柜与水电', '/uploads/demo/kitchen-counter.jpg', '已计入结算金额', '2026-02-06 10:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_change` VALUES (2, 1, 2, 'ADD', '阳台增加储物柜与晾晒区', 6800.00, '2026-04-02', '业主希望提升阳台收纳', '/uploads/demo/site-construction.jpg', '定制柜体，工期顺延 3 天', '2026-04-02 10:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_change` VALUES (3, 1, 2, 'REDUCE', '取消书客房榻榻米定制', 4200.00, '2026-05-10', '业主改为购买成品家具', '/uploads/demo/ceiling.jpg', '已退还对应款项', '2026-05-10 10:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_change` VALUES (4, 1, 3, 'ADD', '增加前台形象墙及灯箱', 5600.00, '2026-10-09', '甲方品牌升级需求', '/uploads/demo/storefront-sign.jpg', '灯箱由启明广告制作', '2026-10-09 10:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_change` VALUES (5, 1, 6, 'ADD', '增加后厨不锈钢操作台', 3800.00, '2026-07-12', '现场实际尺寸与原方案不符', '/uploads/demo/kitchen-counter.jpg', '含运输与安装', '2026-07-12 10:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_change` VALUES (6, 1, 6, 'REDUCE', '取消原吊顶造型，改为平顶', 2600.00, '2026-08-05', '商铺层高受限，简化吊顶', '/uploads/demo/ceiling.jpg', '已与业主确认', '2026-08-05 10:00:00', '2026-11-06 16:08:18', NULL);

-- ----------------------------
-- Table structure for project_contract
-- ----------------------------
DROP TABLE IF EXISTS `project_contract`;
CREATE TABLE `project_contract`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `project_id` bigint(20) NOT NULL COMMENT '关联项目ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '合同名称',
  `contract_no` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '合同编号',
  `amount` decimal(12, 2) NOT NULL DEFAULT 0.00 COMMENT '合同金额',
  `signed_at` date NOT NULL COMMENT '签订日期',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  `attachments` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '附件JSON（路径和原始文件名）',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_contract_company_project`(`company_id`, `project_id`, `deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '项目合同' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of project_contract
-- ----------------------------
INSERT INTO `project_contract` VALUES (1, 1, 1, '锦绣路全屋翻新施工合同', 'HT-2026-001', 186000.00, '2026-01-08', '含主材代购，工期 80 天，质保两年', '[{\"path\":\"/uploads/demo/site-construction.jpg\",\"name\":\"锦绣路全屋翻新施工合同-盖章页.jpg\"}]', '2026-01-08 10:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_contract` VALUES (2, 1, 2, '漕溪北路新房装修施工合同', 'HT-2026-002', 268000.00, '2026-02-10', '环保材料，分段验收，尾款 10% 质保金', '[{\"path\":\"/uploads/demo/water-electric.jpg\",\"name\":\"漕溪北路新房装修施工合同-盖章页.jpg\"}]', '2026-02-10 10:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_contract` VALUES (3, 1, 3, '仲益大厦办公室改造施工合同', 'HT-2026-003', 96000.00, '2026-08-25', '夜间施工，含消防配合与保洁', '[{\"path\":\"/uploads/demo/glass-partition.jpg\",\"name\":\"仲益大厦办公室改造施工合同-盖章页.jpg\"}]', '2026-08-25 10:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_contract` VALUES (4, 1, 4, '国定路厨卫改造施工合同', 'HT-2026-004', 68000.00, '2026-11-01', '工期 60 天，防水质保五年', '[{\"path\":\"/uploads/demo/kitchen-counter.jpg\",\"name\":\"国定路厨卫改造施工合同-盖章页.jpg\"}]', '2026-11-01 10:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_contract` VALUES (5, 1, 5, '恒丰路老房翻新施工合同', 'HT-2026-005', 118000.00, '2026-07-26', '含拆除与垃圾清运，老人友好改造', '[{\"path\":\"/uploads/demo/floor-tile.jpg\",\"name\":\"恒丰路老房翻新施工合同-盖章页.jpg\"}]', '2026-07-26 10:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_contract` VALUES (6, 1, 6, '金桥商铺装修施工合同', 'HT-2026-006', 76000.00, '2026-05-26', '夜间施工，含门头审批配合', '[{\"path\":\"/uploads/demo/storefront-sign.jpg\",\"name\":\"金桥商铺装修施工合同-盖章页.jpg\"}]', '2026-05-26 10:00:00', '2026-11-06 16:08:18', NULL);

-- ----------------------------
-- Table structure for project_log
-- ----------------------------
DROP TABLE IF EXISTS `project_log`;
CREATE TABLE `project_log`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `project_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '项目ID',
  `user_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '记录人ID（site_user.id）',
  `content` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '进度内容',
  `log_date` date NULL DEFAULT NULL COMMENT '进度业务日期(YYYY-MM-DD)',
  `images` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '图片URL，多个用逗号分隔',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_project`(`project_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 196 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '项目进度记录表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of project_log
-- ----------------------------
INSERT INTO `project_log` VALUES (170, 1, 1, 2, '开工交底完成，材料进场：多层实木板40张、石膏板60张已验收。', '2026-01-11', '/uploads/demo/site-construction.jpg', '2026-01-11 18:00:00');
INSERT INTO `project_log` VALUES (171, 1, 1, 1, '水电开槽完成，打压测试合格，业主已确认点位。', '2026-01-24', '/uploads/demo/water-electric.jpg', '2026-01-24 18:00:00');
INSERT INTO `project_log` VALUES (172, 1, 1, 3, '木工吊顶完成 70%，柜体尺寸复尺完成。', '2026-02-12', '/uploads/demo/ceiling.jpg', '2026-02-12 18:00:00');
INSERT INTO `project_log` VALUES (173, 1, 1, 2, '墙面腻子第二遍完成，等待打磨。', '2026-02-26', '[]', '2026-02-26 18:00:00');
INSERT INTO `project_log` VALUES (174, 1, 1, 4, '乳胶漆面漆完成，灯具安装调试完成。', '2026-03-12', '/uploads/demo/ceiling.jpg,/uploads/demo/finished-store.jpg', '2026-03-12 18:00:00');
INSERT INTO `project_log` VALUES (175, 1, 1, 1, '全屋保洁完成，业主验收通过，项目结算。', '2026-03-20', '/uploads/demo/finished-store.jpg', '2026-03-20 18:00:00');
INSERT INTO `project_log` VALUES (176, 1, 2, 3, '开工交底，成品保护完成。', '2026-02-16', '/uploads/demo/site-construction.jpg', '2026-02-16 18:00:00');
INSERT INTO `project_log` VALUES (177, 1, 2, 2, '水电改造完成并验收，隐蔽工程拍照存档。', '2026-03-05', '/uploads/demo/water-electric.jpg', '2026-03-05 18:00:00');
INSERT INTO `project_log` VALUES (178, 1, 2, 7, '瓦工进场，厨卫墙砖铺贴完成。', '2026-04-10', '/uploads/demo/floor-tile.jpg', '2026-04-10 18:00:00');
INSERT INTO `project_log` VALUES (179, 1, 2, 3, '定制柜体安装完成，五金调试。', '2026-05-18', '/uploads/demo/kitchen-counter.jpg', '2026-05-18 18:00:00');
INSERT INTO `project_log` VALUES (180, 1, 2, 2, '软装进场，整体验收通过，待收尾款。', '2026-06-28', '/uploads/demo/finished-store.jpg', '2026-06-28 18:00:00');
INSERT INTO `project_log` VALUES (181, 1, 3, 2, '办理物业开工手续，围挡搭建完成。', '2026-09-02', '/uploads/demo/site-construction.jpg', '2026-09-02 18:00:00');
INSERT INTO `project_log` VALUES (182, 1, 3, 8, '拆除与隔墙完成，消防点位确认。', '2026-09-20', '/uploads/demo/glass-partition.jpg', '2026-09-20 18:00:00');
INSERT INTO `project_log` VALUES (183, 1, 3, 2, '吊顶龙骨完成，强弱电布线完成。', '2026-10-12', '/uploads/demo/ceiling.jpg', '2026-10-12 18:00:00');
INSERT INTO `project_log` VALUES (184, 1, 3, 5, '前台形象墙基层完成，等待灯箱安装。', '2026-11-02', '/uploads/demo/storefront-sign.jpg', '2026-11-02 18:00:00');
INSERT INTO `project_log` VALUES (185, 1, 4, 5, '合同签订，主材清单确认。', '2026-11-01', '/uploads/demo/kitchen-counter.jpg', '2026-11-01 18:00:00');
INSERT INTO `project_log` VALUES (186, 1, 4, 2, '开工前现场复尺，确定 11 月 15 日进场。', '2026-11-04', '/uploads/demo/site-construction.jpg', '2026-11-04 18:00:00');
INSERT INTO `project_log` VALUES (187, 1, 5, 4, '开工，旧房拆除完成。', '2026-08-02', '/uploads/demo/site-construction.jpg', '2026-08-02 18:00:00');
INSERT INTO `project_log` VALUES (188, 1, 5, 2, '水电改造完成，防水闭水试验通过。', '2026-08-20', '/uploads/demo/water-electric.jpg', '2026-08-20 18:00:00');
INSERT INTO `project_log` VALUES (189, 1, 5, 7, '泥木阶段完成，门窗安装。', '2026-09-15', '/uploads/demo/floor-tile.jpg', '2026-09-15 18:00:00');
INSERT INTO `project_log` VALUES (190, 1, 5, 4, '墙面乳胶漆完成。', '2026-10-08', '[]', '2026-10-08 18:00:00');
INSERT INTO `project_log` VALUES (191, 1, 5, 8, '灯具与五金安装，准备竣工验收。', '2026-11-02', '/uploads/demo/finished-store.jpg', '2026-11-02 18:00:00');
INSERT INTO `project_log` VALUES (192, 1, 6, 2, '商铺夜间施工交底，办理施工许可。', '2026-06-02', '/uploads/demo/site-construction.jpg', '2026-06-02 18:00:00');
INSERT INTO `project_log` VALUES (193, 1, 6, 3, '水电与地面找平完成。', '2026-06-20', '/uploads/demo/water-electric.jpg,/uploads/demo/floor-tile.jpg', '2026-06-20 18:00:00');
INSERT INTO `project_log` VALUES (194, 1, 6, 7, '门头招牌安装完成，玻璃隔断进场。', '2026-07-15', '/uploads/demo/storefront-sign.jpg,/uploads/demo/glass-partition.jpg', '2026-07-15 18:00:00');
INSERT INTO `project_log` VALUES (195, 1, 6, 1, '竣工验收，交付业主。', '2026-08-15', '/uploads/demo/finished-store.jpg,/uploads/demo/glass-partition.jpg', '2026-08-15 18:00:00');

-- ----------------------------
-- Table structure for project_material
-- ----------------------------
DROP TABLE IF EXISTS `project_material`;
CREATE TABLE `project_material`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `project_id` bigint(20) NOT NULL COMMENT '关联项目ID',
  `catalog_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '公司级材料档案ID',
  `warehouse_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属仓库ID，关联 warehouse.id',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '材料名称',
  `spec` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '规格 / 型号',
  `unit` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '单位，如：米、平米、个',
  `quantity` decimal(12, 2) NOT NULL DEFAULT 0.00 COMMENT '清单数量',
  `unit_price` decimal(12, 2) NOT NULL DEFAULT 0.00 COMMENT '单价',
  `amount` decimal(12, 2) NOT NULL DEFAULT 0.00 COMMENT '金额 = 数量 × 单价',
  `supplier` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '供应商',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间，NULL表示未删除',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_material_project`(`project_id`) USING BTREE,
  INDEX `idx_material_deleted`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 30 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '项目材料清单' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of project_material
-- ----------------------------
INSERT INTO `project_material` VALUES (1, 1, 1, 38, 1, '多层实木板', '18mm E1 级', '张', 50.00, 128.00, 6400.00, '红星美凯龙建材（上海供应处）', '', '2026-01-10 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (2, 1, 1, 45, 1, '水性木器漆', '哑光白 5L', '桶', 12.00, 368.00, 4416.00, '立邦涂料上海总代理', '', '2026-01-10 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (3, 1, 1, 33, 1, '不锈钢合页', '4寸 缓冲', '个', 180.00, 8.50, 1530.00, '红星美凯龙建材（上海供应处）', '', '2026-01-10 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (4, 1, 1, 50, 1, '纸面石膏板', '9mm', '张', 60.00, 42.00, 2520.00, '红星美凯龙建材（上海供应处）', '', '2026-01-10 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (5, 1, 2, 41, 1, '强化复合地板', '1210×165×12mm', '㎡', 80.00, 118.00, 9440.00, '红星美凯龙建材（上海供应处）', '', '2026-02-15 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (6, 1, 2, 35, 1, '乳胶漆', '立邦净味120 18L', '桶', 10.00, 520.00, 5200.00, '立邦涂料上海总代理', '', '2026-02-15 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (7, 1, 2, 44, 1, '橱柜板材', '爱格板 18mm', '张', 30.00, 268.00, 8040.00, '欧派家居上海经销商', '', '2026-02-15 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (8, 1, 2, 51, 1, '耐水腻子粉', '20kg', '袋', 40.00, 45.00, 1800.00, '红星美凯龙建材（上海供应处）', '', '2026-02-15 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (9, 1, 2, 40, 1, '开关插座', '西门子 五孔', '个', 70.00, 32.00, 2240.00, '红星美凯龙建材（上海供应处）', '', '2026-02-15 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (10, 1, 3, 52, 1, '轻钢龙骨', '75 型', '米', 200.00, 12.00, 2400.00, '红星美凯龙建材（上海供应处）', '', '2026-09-01 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (11, 1, 3, 49, 1, '矿棉板', '600×600', '块', 300.00, 18.00, 5400.00, '红星美凯龙建材（上海供应处）', '', '2026-09-01 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (12, 1, 3, 58, 1, '防静电地板', '600×600', '块', 120.00, 165.00, 19800.00, '红星美凯龙建材（上海供应处）', '', '2026-09-01 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (13, 1, 3, 34, 1, '乳胶漆', '18L', '桶', 6.00, 520.00, 3120.00, '立邦涂料上海总代理', '', '2026-09-01 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (14, 1, 3, 54, 1, '铝扣板', '300×300', '㎡', 0.00, 85.00, 0.00, '红星美凯龙建材（上海供应处）', '', '2026-09-01 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (15, 1, 3, 36, 1, '前台灯箱', '定制 2.4m×0.8m', '个', 0.00, 3600.00, 0.00, '上海启明广告制作有限公司', '', '2026-09-01 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (16, 1, 4, 57, 1, '防水涂料', '东方雨虹 20kg', '桶', 8.00, 380.00, 3040.00, '东方雨虹防水工程部', '', '2026-11-15 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (17, 1, 4, 48, 1, '瓷砖', '600×600 防滑', '片', 0.00, 68.00, 0.00, '红星美凯龙建材（上海供应处）', '', '2026-11-15 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (18, 1, 4, 46, 1, '水龙头', '摩恩 单孔', '个', 0.00, 460.00, 0.00, '欧派家居上海经销商', '', '2026-11-15 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (19, 1, 4, 42, 1, '整体橱柜', '石英石台面', '延米', 0.00, 1800.00, 0.00, '欧派家居上海经销商', '', '2026-11-15 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (20, 1, 5, 39, 1, '实木复合门', '标准 樘', '樘', 4.00, 1280.00, 5120.00, '欧派家居上海经销商', '', '2026-08-01 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (21, 1, 5, 55, 1, '门套线', '实木 6cm', '米', 60.00, 28.00, 1680.00, '红星美凯龙建材（上海供应处）', '', '2026-08-01 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (22, 1, 5, 56, 1, '防水涂料', '20kg', '桶', 12.00, 380.00, 4560.00, '东方雨虹防水工程部', '', '2026-08-01 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (23, 1, 5, 47, 1, '瓷砖', '600×600', '片', 120.00, 68.00, 8160.00, '红星美凯龙建材（上海供应处）', '', '2026-08-01 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (24, 1, 5, 34, 1, '乳胶漆', '18L', '桶', 6.00, 520.00, 3120.00, '立邦涂料上海总代理', '', '2026-08-01 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (25, 1, 6, 43, 1, '木纹铝方通', '30×80mm', '米', 90.00, 45.00, 4050.00, '红星美凯龙建材（上海供应处）', '', '2026-06-01 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (26, 1, 6, 53, 1, '钢化玻璃隔断', '10mm 钢化', '㎡', 28.00, 260.00, 7280.00, '红星美凯龙建材（上海供应处）', '', '2026-06-01 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (27, 1, 6, 37, 1, '办公地毯', '方块拼接地毯', '㎡', 60.00, 85.00, 5100.00, '红星美凯龙建材（上海供应处）', '', '2026-06-01 09:00:00', '2026-11-06 16:08:18', NULL);
INSERT INTO `project_material` VALUES (28, 1, 6, 32, 1, 'LED 筒灯', '7W 中性光', '个', 40.00, 45.00, 1800.00, '红星美凯龙建材（上海供应处）', '', '2026-06-01 09:00:00', '2026-11-06 16:08:18', NULL);

-- ----------------------------
-- Table structure for project_material_usage
-- ----------------------------
DROP TABLE IF EXISTS `project_material_usage`;
CREATE TABLE `project_material_usage`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `project_id` bigint(20) NOT NULL COMMENT '关联项目ID',
  `material_id` bigint(20) NOT NULL COMMENT '关联材料ID',
  `warehouse_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '出库仓库ID，0 表示沿用材料所属仓库',
  `quantity` decimal(12, 2) NOT NULL DEFAULT 0.00 COMMENT '本次领用数量',
  `used_at` date NOT NULL COMMENT '领用日期',
  `receiver` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '领用人',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  `created_by` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建人ID，0 表示历史数据',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `stocked_at` datetime NULL DEFAULT NULL COMMENT '出库流水生成时间，NULL 表示未出库',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间，NULL表示未删除',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_usage_project`(`project_id`) USING BTREE,
  INDEX `idx_usage_material`(`material_id`) USING BTREE,
  INDEX `idx_usage_deleted`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 82 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '项目材料领用记录' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of project_material_usage
-- ----------------------------
INSERT INTO `project_material_usage` VALUES (55, 1, 1, 1, 0, 30.00, '2026-01-15', '赵明', '施工领用', 0, '2026-01-15 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (56, 1, 1, 1, 0, 18.00, '2026-02-20', '赵明', '施工领用', 0, '2026-02-20 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (57, 1, 1, 2, 0, 12.00, '2026-03-08', '赵明', '施工领用', 0, '2026-03-08 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (58, 1, 1, 3, 0, 120.00, '2026-02-22', '赵明', '施工领用', 0, '2026-02-22 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (59, 1, 1, 4, 0, 55.00, '2026-01-20', '赵明', '施工领用', 0, '2026-01-20 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (60, 1, 1, 4, 0, 5.00, '2026-02-05', '赵明', '盘亏出库', 0, '2026-02-05 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (61, 1, 2, 5, 0, 76.00, '2026-03-10', '陈晓东', '施工领用', 0, '2026-03-10 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (62, 1, 2, 6, 0, 9.00, '2026-03-25', '陈晓东', '施工领用', 0, '2026-03-25 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (63, 1, 2, 7, 0, 28.00, '2026-04-20', '陈晓东', '施工领用', 0, '2026-04-20 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (64, 1, 2, 8, 0, 36.00, '2026-04-05', '陈晓东', '施工领用', 0, '2026-04-05 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (65, 1, 2, 8, 0, 4.00, '2026-05-08', '陈晓东', '盘亏出库', 0, '2026-05-08 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (66, 1, 2, 9, 0, 60.00, '2026-05-10', '陈晓东', '施工领用', 0, '2026-05-10 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (67, 1, 3, 10, 0, 150.00, '2026-09-18', '郑雨欣', '施工领用', 0, '2026-09-18 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (68, 1, 3, 11, 0, 240.00, '2026-09-22', '郑雨欣', '施工领用', 0, '2026-09-22 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (69, 1, 3, 12, 0, 100.00, '2026-10-10', '郑雨欣', '施工领用', 0, '2026-10-10 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (70, 1, 3, 12, 0, 10.00, '2026-10-20', '陈志强', '调拨出库（项目领用）', 0, '2026-10-20 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (71, 1, 3, 13, 0, 5.00, '2026-10-20', '郑雨欣', '施工领用', 0, '2026-10-20 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (72, 1, 5, 20, 0, 4.00, '2026-09-25', '王建国', '施工领用', 0, '2026-09-25 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (73, 1, 5, 21, 0, 52.00, '2026-09-26', '王建国', '施工领用', 0, '2026-09-26 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (74, 1, 5, 22, 0, 8.00, '2026-08-18', '王建国', '施工领用', 0, '2026-08-18 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (75, 1, 5, 23, 0, 115.00, '2026-09-05', '王建国', '施工领用', 0, '2026-09-05 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (76, 1, 5, 24, 0, 5.00, '2026-10-12', '王建国', '施工领用', 0, '2026-10-12 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (77, 1, 6, 25, 0, 85.00, '2026-06-24', '赵明', '施工领用', 0, '2026-06-24 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (78, 1, 6, 26, 0, 26.00, '2026-07-08', '赵明', '施工领用', 0, '2026-07-08 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (79, 1, 6, 26, 0, 2.00, '2026-08-10', '赵明', '盘亏出库', 0, '2026-08-10 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (80, 1, 6, 27, 0, 58.00, '2026-07-22', '赵明', '施工领用', 0, '2026-07-22 09:00:00', NULL, NULL);
INSERT INTO `project_material_usage` VALUES (81, 1, 6, 28, 0, 36.00, '2026-07-28', '赵明', '施工领用', 0, '2026-07-28 09:00:00', NULL, NULL);

-- ----------------------------
-- Table structure for project_payment
-- ----------------------------
DROP TABLE IF EXISTS `project_payment`;
CREATE TABLE `project_payment`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `project_id` bigint(20) NOT NULL COMMENT '关联项目ID',
  `amount` decimal(10, 2) NOT NULL DEFAULT 0.00 COMMENT '收款金额',
  `payment_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '款项类型：定金/首期款/中期款/尾款',
  `payment_date` date NOT NULL COMMENT '收款日期',
  `fund_account_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '收款资金账户ID，对应 fin_fund_account.id',
  `voucher_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '自动生成的记账凭证ID，0 表示未记账',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '备注，如：微信转账',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间，NULL表示未删除',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_project_id`(`project_id`) USING BTREE,
  INDEX `idx_payment_date`(`payment_date`) USING BTREE,
  INDEX `idx_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 19 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '项目收款记录表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of project_payment
-- ----------------------------
INSERT INTO `project_payment` VALUES (1, 1, 1, 37200.00, '定金', '2026-01-12', 25, 2, '合同签订，首期定金 20%', '2026-01-12 10:00:00', NULL);
INSERT INTO `project_payment` VALUES (2, 1, 1, 93000.00, '中期款', '2026-02-15', 25, 12, '水电验收后进度款 50%', '2026-02-15 10:00:00', NULL);
INSERT INTO `project_payment` VALUES (3, 1, 2, 53600.00, '定金', '2026-02-18', 25, 13, '合同签订，首期定金 20%', '2026-02-18 10:00:00', NULL);
INSERT INTO `project_payment` VALUES (4, 1, 1, 55800.00, '尾款', '2026-03-18', 25, 26, '竣工验收尾款 30%', '2026-03-18 10:00:00', NULL);
INSERT INTO `project_payment` VALUES (5, 1, 2, 134000.00, '中期款', '2026-04-20', 25, 32, '泥木验收后进度款 50%', '2026-04-20 10:00:00', NULL);
INSERT INTO `project_payment` VALUES (6, 1, 6, 22800.00, '定金', '2026-06-05', 25, 36, '合同签订，首期定金 30%', '2026-06-05 10:00:00', NULL);
INSERT INTO `project_payment` VALUES (7, 1, 2, 53600.00, '尾款', '2026-06-25', 25, 44, '已收 90%，余 26,800 待收', '2026-06-25 10:00:00', NULL);
INSERT INTO `project_payment` VALUES (8, 1, 5, 23600.00, '定金', '2026-08-05', 25, 52, '合同签订，首期定金 20%', '2026-08-05 10:00:00', NULL);
INSERT INTO `project_payment` VALUES (9, 1, 6, 30400.00, '尾款', '2026-08-20', 25, 59, '尾款部分到账，余 22,800 待收', '2026-08-20 10:00:00', NULL);
INSERT INTO `project_payment` VALUES (10, 1, 3, 28800.00, '定金', '2026-09-03', 25, 61, '开工前定金 30%', '2026-09-03 10:00:00', NULL);
INSERT INTO `project_payment` VALUES (11, 1, 5, 47200.00, '中期款', '2026-09-18', 25, 67, '水电验收后进度款 40%', '2026-09-18 10:00:00', NULL);
INSERT INTO `project_payment` VALUES (12, 1, 3, 38400.00, '中期款', '2026-10-15', 25, 76, '吊顶完成后进度款 40%', '2026-10-15 10:00:00', NULL);
INSERT INTO `project_payment` VALUES (14, 1, 5, 35400.00, '尾款', '2026-10-25', 24, 78, '部分尾款，余 11,800 待收', '2026-10-25 10:00:00', NULL);
INSERT INTO `project_payment` VALUES (15, 1, 4, 20400.00, '定金', '2026-11-02', 24, 81, '开工前定金 30%（现金）', '2026-11-02 10:00:00', NULL);
INSERT INTO `project_payment` VALUES (17, 1, 4, 3000.00, '???', '2026-11-06', 25, 87, '????????', '2026-11-06 08:35:53', '2026-11-06 08:36:38');
INSERT INTO `project_payment` VALUES (18, 1, 4, 3800.00, '中期款', '2026-11-06', 25, 89, '现场测量增项预收（调整）', '2026-11-06 08:36:47', '2026-11-06 08:37:00');

-- ----------------------------
-- Table structure for project_payout
-- ----------------------------
DROP TABLE IF EXISTS `project_payout`;
CREATE TABLE `project_payout`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `project_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '关联项目ID，0 表示公司公共支出',
  `amount` decimal(10, 2) NOT NULL DEFAULT 0.00 COMMENT '付款金额',
  `payout_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '费用类型：材料款/人工费/施工费/设计费/管理费/销售费/其他',
  `payout_date` date NOT NULL COMMENT '付款日期',
  `vendor_id` int(11) NOT NULL COMMENT '收款方：供应商 / 工人 / 公司',
  `purchase_order_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '关联采购单ID，0表示手工付款',
  `fund_account_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '付款资金账户ID，对应 fin_fund_account.id',
  `voucher_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '自动生成的记账凭证ID，0 表示未记账',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '备注，如：转账、结算尾款',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间，NULL表示未删除',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_payout_project`(`project_id`) USING BTREE,
  INDEX `idx_payout_date`(`payout_date`) USING BTREE,
  INDEX `idx_payout_deleted`(`deleted_at`) USING BTREE,
  INDEX `idx_payout_purchase_order`(`purchase_order_id`) USING BTREE,
  INDEX `idx_payout_vendor`(`vendor_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 70 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '项目付款记录表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of project_payout
-- ----------------------------
INSERT INTO `project_payout` VALUES (1, 1, 0, 8000.00, '管理费', '2026-01-05', 110, 0, 24, 1, '办公室一月物业与水电', '2026-01-05 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (2, 1, 1, 5120.00, '材料款', '2026-01-13', 105, 1, 24, 3, '采购单 PO-2026-0001 入库结算', '2026-01-13 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (3, 1, 1, 2520.00, '材料款', '2026-01-14', 105, 4, 24, 4, '采购单 PO-2026-0004 入库结算', '2026-01-14 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (4, 1, 1, 4416.00, '材料款', '2026-01-16', 109, 2, 24, 5, '采购单 PO-2026-0002 入库结算', '2026-01-16 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (5, 1, 1, 1360.00, '材料款', '2026-01-18', 105, 3, 24, 6, '采购单 PO-2026-0003 入库结算', '2026-01-18 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (6, 1, 1, 1280.00, '材料款', '2026-01-18', 105, 0, 24, 7, '多层实木板补货入库（现结）', '2026-01-18 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (7, 1, 1, 10000.00, '材料款', '2026-01-20', 105, 0, 25, 8, '瓷砖、水泥辅料现结', '2026-01-20 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (8, 1, 1, 5000.00, '其他', '2026-01-25', 113, 0, 24, 9, '建筑垃圾清运两次', '2026-01-25 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (9, 1, 0, 8000.00, '管理费', '2026-02-05', 110, 0, 25, 10, '办公室二月物业与水电', '2026-02-05 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (10, 1, 1, 15000.00, '人工费', '2026-02-10', 107, 0, 25, 11, '木工与瓦工阶段人工费', '2026-02-10 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (11, 1, 2, 9440.00, '材料款', '2026-02-22', 105, 5, 25, 14, '采购单 PO-2026-0005 入库结算', '2026-02-22 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (12, 1, 2, 16000.00, '设计费', '2026-02-25', 112, 0, 25, 15, '全案设计与出图费', '2026-02-25 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (13, 1, 1, 12000.00, '材料款', '2026-02-28', 108, 0, 25, 16, '定制柜体中期料款', '2026-02-28 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (14, 1, 1, 10600.00, '人工费', '2026-02-28', 107, 0, 25, 17, '油漆班组人工费', '2026-02-28 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (15, 1, 2, 5200.00, '材料款', '2026-03-03', 109, 6, 25, 18, '采购单 PO-2026-0006 入库结算', '2026-03-03 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (16, 1, 1, 10000.00, '材料款', '2026-03-05', 105, 0, 25, 19, '五金辅料补货', '2026-03-05 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (17, 1, 2, 1800.00, '材料款', '2026-03-07', 105, 8, 25, 20, '采购单 PO-2026-0008 入库结算', '2026-03-07 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (18, 1, 2, 30000.00, '材料款', '2026-03-08', 105, 0, 25, 21, '瓷砖、地板首付款', '2026-03-08 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (19, 1, 0, 5000.00, '销售费', '2026-03-10', 111, 0, 25, 22, '抖音本地推三月投放', '2026-03-10 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (20, 1, 1, 12400.00, '施工费', '2026-03-12', 107, 0, 25, 23, '安装与收尾施工费', '2026-03-12 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (21, 1, 1, 6000.00, '管理费', '2026-03-15', 110, 0, 25, 24, '项目现场管理费', '2026-03-15 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (22, 1, 2, 8040.00, '材料款', '2026-03-17', 108, 7, 25, 25, '采购单 PO-2026-0007 入库结算', '2026-03-17 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (23, 1, 1, 9000.00, '设计费', '2026-03-18', 112, 0, 25, 27, '全案设计尾款', '2026-03-18 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (24, 1, 2, 22000.00, '人工费', '2026-03-28', 107, 0, 25, 28, '水电与泥工人工费', '2026-03-28 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (25, 1, 2, 640.00, '材料款', '2026-04-04', 105, 9, 25, 29, '采购单 PO-2026-0009 入库结算', '2026-04-04 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (26, 1, 0, 1200.00, '其他', '2026-04-15', 113, 0, 24, 30, '办公耗材与快递配送', '2026-04-15 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (27, 1, 2, 22000.00, '材料款', '2026-04-18', 105, 0, 25, 31, '卫浴与五金到货款', '2026-04-18 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (28, 1, 0, 12000.00, '管理费', '2026-05-10', 110, 0, 25, 33, '办公室二季度租金分摊', '2026-05-10 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (29, 1, 2, 18000.00, '人工费', '2026-05-20', 107, 0, 25, 34, '木工与油漆人工费', '2026-05-20 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (30, 1, 2, 18000.00, '材料款', '2026-05-22', 105, 0, 25, 35, '软装与辅材结算', '2026-05-22 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (31, 1, 6, 4000.00, '设计费', '2026-06-08', 112, 0, 25, 37, '商铺门头与空间设计费', '2026-06-08 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (32, 1, 0, 8000.00, '销售费', '2026-06-08', 111, 0, 25, 38, '六月线上推广费', '2026-06-08 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (33, 1, 2, 12000.00, '人工费', '2026-06-10', 107, 0, 25, 39, '安装收尾人工费', '2026-06-10 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (34, 1, 6, 4050.00, '材料款', '2026-06-14', 105, 24, 25, 40, '采购单 PO-2026-0024 入库结算', '2026-06-14 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (35, 1, 6, 5100.00, '材料款', '2026-06-17', 105, 26, 25, 41, '采购单 PO-2026-0026 入库结算', '2026-06-17 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (36, 1, 6, 6000.00, '销售费', '2026-06-18', 111, 0, 25, 42, '本地生活平台推广费', '2026-06-18 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (37, 1, 6, 7280.00, '材料款', '2026-06-20', 105, 25, 25, 43, '采购单 PO-2026-0025 入库结算', '2026-06-20 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (38, 1, 6, 10000.00, '人工费', '2026-06-25', 107, 0, 25, 45, '水电与地面施工人工费', '2026-06-25 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (39, 1, 2, 8000.00, '管理费', '2026-06-28', 110, 0, 25, 46, '项目现场管理费', '2026-06-28 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (40, 1, 0, 8000.00, '管理费', '2026-07-05', 110, 0, 25, 47, '办公室七月物业与水电', '2026-07-05 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (41, 1, 6, 6000.00, '材料款', '2026-07-15', 105, 0, 25, 48, '铝方通与灯具补货', '2026-07-15 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (42, 1, 6, 8000.00, '人工费', '2026-07-20', 107, 0, 25, 49, '门头与玻璃安装人工费', '2026-07-20 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (43, 1, 6, 2000.00, '管理费', '2026-07-28', 110, 0, 25, 50, '项目现场管理费', '2026-07-28 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (44, 1, 6, 4000.00, '施工费', '2026-08-01', 107, 0, 25, 51, '灯具与招牌安装收尾', '2026-08-01 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (45, 1, 0, 6000.00, '销售费', '2026-08-06', 111, 0, 25, 53, '八月线上推广费', '2026-08-06 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (46, 1, 5, 8000.00, '设计费', '2026-08-08', 112, 0, 25, 54, '老房翻新全案设计费', '2026-08-08 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (47, 1, 5, 3800.00, '材料款', '2026-08-12', 106, 22, 25, 55, '采购单 PO-2026-0022 入库结算', '2026-08-12 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (48, 1, 5, 8160.00, '材料款', '2026-08-14', 105, 23, 25, 56, '采购单 PO-2026-0023 入库结算', '2026-08-14 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (49, 1, 5, 5120.00, '材料款', '2026-08-17', 108, 20, 25, 57, '采购单 PO-2026-0020 入库结算', '2026-08-17 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (50, 1, 5, 1680.00, '材料款', '2026-08-18', 105, 21, 25, 58, '采购单 PO-2026-0021 入库结算', '2026-08-18 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (51, 1, 5, 12000.00, '人工费', '2026-08-28', 107, 0, 25, 60, '拆除与水电人工费', '2026-08-28 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (52, 1, 0, 8000.00, '管理费', '2026-09-05', 110, 0, 25, 62, '办公室九月物业与水电', '2026-09-05 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (53, 1, 3, 2400.00, '材料款', '2026-09-07', 105, 10, 25, 63, '采购单 PO-2026-0010 入库结算', '2026-09-07 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (54, 1, 3, 5400.00, '材料款', '2026-09-10', 105, 11, 25, 64, '采购单 PO-2026-0011 入库结算', '2026-09-10 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (55, 1, 5, 7000.00, '材料款', '2026-09-10', 105, 0, 25, 65, '门窗五金补货', '2026-09-10 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (56, 1, 3, 5000.00, '设计费', '2026-09-15', 112, 0, 25, 66, '办公空间方案设计费', '2026-09-15 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (57, 1, 3, 8000.00, '人工费', '2026-09-25', 107, 0, 25, 68, '拆除与隔墙人工费', '2026-09-25 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (58, 1, 3, 19800.00, '材料款', '2026-09-28', 105, 12, 25, 69, '采购单 PO-2026-0012 入库结算', '2026-09-28 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (59, 1, 3, 3120.00, '材料款', '2026-09-30', 109, 13, 25, 70, '采购单 PO-2026-0013 入库结算', '2026-09-30 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (60, 1, 5, 3120.00, '材料款', '2026-09-30', 109, 0, 25, 71, '乳胶漆补货入库（现结）', '2026-09-30 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (61, 1, 5, 4000.00, '管理费', '2026-09-30', 110, 0, 25, 72, '项目现场管理费', '2026-09-30 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (62, 1, 5, 12000.00, '人工费', '2026-10-06', 107, 0, 25, 74, '泥木与油漆人工费', '2026-10-06 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (63, 1, 0, 4000.00, '销售费', '2026-10-09', 111, 0, 25, 75, '十月线上推广费', '2026-10-09 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (64, 1, 3, 8000.00, '人工费', '2026-10-28', 107, 0, 25, 79, '吊顶与布线人工费', '2026-10-28 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (65, 1, 4, 3040.00, '材料款', '2026-11-01', 106, 16, 25, 80, '采购单 PO-2026-0016 入库结算', '2026-11-01 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (66, 1, 3, 2000.00, '其他', '2026-11-03', 113, 0, 24, 82, '夜间施工垃圾清运', '2026-11-03 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (67, 1, 4, 3000.00, '设计费', '2026-11-03', 112, 0, 25, 83, '厨卫改造方案设计费', '2026-11-03 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (68, 1, 0, 2600.00, '其他', '2026-11-03', 113, 0, 24, 84, '门店物料搬运与配送', '2026-11-03 10:00:00', NULL);
INSERT INTO `project_payout` VALUES (69, 1, 4, 500.00, '材料款', '2026-11-06', 105, 0, 25, 90, '联调测试付款', '2026-11-06 08:37:19', '2026-11-06 08:37:19');

-- ----------------------------
-- Table structure for project_purchase_order
-- ----------------------------
DROP TABLE IF EXISTS `project_purchase_order`;
CREATE TABLE `project_purchase_order`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `order_no` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '采购单号',
  `project_id` bigint(20) NOT NULL COMMENT '关联项目ID',
  `material_id` bigint(20) NOT NULL COMMENT '关联材料ID',
  `vendor_id` int(10) UNSIGNED NOT NULL COMMENT '供应商ID',
  `warehouse_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '出入库仓库ID，0 表示沿用材料所属仓库',
  `purchase_date` date NOT NULL COMMENT '采购日期',
  `material_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '材料名称快照',
  `spec` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '规格快照',
  `unit` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '单位快照',
  `quantity` decimal(12, 2) NOT NULL COMMENT '采购数量',
  `unit_price` decimal(12, 2) NOT NULL COMMENT '采购单价',
  `amount` decimal(12, 2) NOT NULL COMMENT '采购金额',
  `supplier` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '供应商名称快照',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  `created_by` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建人ID，0 表示历史数据',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'DRAFT' COMMENT '状态：DRAFT草稿 / CONFIRMED待付款 / PAID已付款',
  `payout_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '关联付款记录ID',
  `confirmed_at` datetime NULL DEFAULT NULL COMMENT '确认时间',
  `paid_at` datetime NULL DEFAULT NULL COMMENT '付款时间',
  `stocked_at` datetime NULL DEFAULT NULL COMMENT '入库时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_purchase_order_no`(`order_no`) USING BTREE,
  INDEX `idx_purchase_order_no`(`order_no`) USING BTREE,
  INDEX `idx_purchase_project`(`project_id`) USING BTREE,
  INDEX `idx_purchase_vendor`(`vendor_id`) USING BTREE,
  INDEX `idx_purchase_status`(`status`) USING BTREE,
  INDEX `idx_purchase_payout`(`payout_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 28 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '材料采购单' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of project_purchase_order
-- ----------------------------
INSERT INTO `project_purchase_order` VALUES (1, 1, 'PO-2026-0001', 1, 1, 105, 0, '2026-01-11', '多层实木板', '18mm E1 级', '张', 40.00, 128.00, 5120.00, '红星美凯龙建材（上海供应处）', '款已付，材料已入库', 0, 'PAID', 2, '2026-01-11 10:00:00', '2026-01-13 10:00:00', '2026-01-13 15:00:00', '2026-01-11 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (2, 1, 'PO-2026-0002', 1, 2, 109, 0, '2026-01-14', '水性木器漆', '哑光白 5L', '桶', 12.00, 368.00, 4416.00, '立邦涂料上海总代理', '款已付，材料已入库', 0, 'PAID', 4, '2026-01-14 10:00:00', '2026-01-16 10:00:00', '2026-01-16 15:00:00', '2026-01-14 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (3, 1, 'PO-2026-0003', 1, 3, 105, 0, '2026-01-16', '不锈钢合页', '4寸 缓冲', '个', 160.00, 8.50, 1360.00, '红星美凯龙建材（上海供应处）', '款已付，材料已入库', 0, 'PAID', 5, '2026-01-16 10:00:00', '2026-01-18 10:00:00', '2026-01-18 15:00:00', '2026-01-16 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (4, 1, 'PO-2026-0004', 1, 4, 105, 0, '2026-01-12', '纸面石膏板', '9mm', '张', 60.00, 42.00, 2520.00, '红星美凯龙建材（上海供应处）', '款已付，材料已入库', 0, 'PAID', 3, '2026-01-12 10:00:00', '2026-01-14 10:00:00', '2026-01-14 15:00:00', '2026-01-12 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (5, 1, 'PO-2026-0005', 2, 5, 105, 0, '2026-02-20', '强化复合地板', '1210×165×12mm', '㎡', 80.00, 118.00, 9440.00, '红星美凯龙建材（上海供应处）', '款已付，材料已入库', 0, 'PAID', 11, '2026-02-20 10:00:00', '2026-02-22 10:00:00', '2026-02-22 15:00:00', '2026-02-20 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (6, 1, 'PO-2026-0006', 2, 6, 109, 0, '2026-03-01', '乳胶漆', '立邦净味120 18L', '桶', 10.00, 520.00, 5200.00, '立邦涂料上海总代理', '款已付，材料已入库', 0, 'PAID', 15, '2026-03-01 10:00:00', '2026-03-03 10:00:00', '2026-03-03 15:00:00', '2026-03-01 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (7, 1, 'PO-2026-0007', 2, 7, 108, 0, '2026-03-15', '橱柜板材', '爱格板 18mm', '张', 30.00, 268.00, 8040.00, '欧派家居上海经销商', '款已付，材料已入库', 0, 'PAID', 22, '2026-03-15 10:00:00', '2026-03-17 10:00:00', '2026-03-17 15:00:00', '2026-03-15 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (8, 1, 'PO-2026-0008', 2, 8, 105, 0, '2026-03-05', '耐水腻子粉', '20kg', '袋', 40.00, 45.00, 1800.00, '红星美凯龙建材（上海供应处）', '款已付，材料已入库', 0, 'PAID', 17, '2026-03-05 10:00:00', '2026-03-07 10:00:00', '2026-03-07 15:00:00', '2026-03-05 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (9, 1, 'PO-2026-0009', 2, 9, 105, 0, '2026-04-02', '开关插座', '西门子 五孔', '个', 20.00, 32.00, 640.00, '红星美凯龙建材（上海供应处）', '款已付，材料已入库', 0, 'PAID', 25, '2026-04-02 10:00:00', '2026-04-04 10:00:00', '2026-04-04 15:00:00', '2026-04-02 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (10, 1, 'PO-2026-0010', 3, 10, 105, 0, '2026-09-05', '轻钢龙骨', '75 型', '米', 200.00, 12.00, 2400.00, '红星美凯龙建材（上海供应处）', '款已付，材料已入库', 0, 'PAID', 53, '2026-09-05 10:00:00', '2026-09-07 10:00:00', '2026-09-07 15:00:00', '2026-09-05 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (11, 1, 'PO-2026-0011', 3, 11, 105, 0, '2026-09-08', '矿棉板', '600×600', '块', 300.00, 18.00, 5400.00, '红星美凯龙建材（上海供应处）', '款已付，材料已入库', 0, 'PAID', 54, '2026-09-08 10:00:00', '2026-09-10 10:00:00', '2026-09-10 15:00:00', '2026-09-08 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (12, 1, 'PO-2026-0012', 3, 12, 105, 0, '2026-09-26', '防静电地板', '600×600', '块', 120.00, 165.00, 19800.00, '红星美凯龙建材（上海供应处）', '款已付，材料已入库', 0, 'PAID', 58, '2026-09-26 10:00:00', '2026-09-28 10:00:00', '2026-09-28 15:00:00', '2026-09-26 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (13, 1, 'PO-2026-0013', 3, 13, 109, 0, '2026-09-28', '乳胶漆', '18L', '桶', 6.00, 520.00, 3120.00, '立邦涂料上海总代理', '款已付，材料已入库', 0, 'PAID', 59, '2026-09-28 10:00:00', '2026-09-30 10:00:00', '2026-09-30 15:00:00', '2026-09-28 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (14, 1, 'PO-2026-0014', 3, 14, 105, 0, '2026-10-30', '铝扣板', '300×300', '㎡', 60.00, 85.00, 5100.00, '红星美凯龙建材（上海供应处）', '已确认，等待付款', 0, 'CONFIRMED', 0, '2026-10-30 10:00:00', NULL, NULL, '2026-10-30 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (15, 1, 'PO-2026-0015', 3, 15, 114, 0, '2026-10-25', '前台灯箱', '定制 2.4m×0.8m', '个', 1.00, 3600.00, 3600.00, '上海启明广告制作有限公司', '已确认，等待付款', 0, 'CONFIRMED', 0, '2026-10-25 10:00:00', NULL, NULL, '2026-10-25 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (16, 1, 'PO-2026-0016', 4, 16, 106, 0, '2026-10-30', '防水涂料', '东方雨虹 20kg', '桶', 8.00, 380.00, 3040.00, '东方雨虹防水工程部', '款已付，材料已入库', 0, 'PAID', 65, '2026-10-30 10:00:00', '2026-11-01 10:00:00', '2026-11-01 15:00:00', '2026-10-30 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (17, 1, 'PO-2026-0017', 4, 17, 105, 0, '2026-11-01', '瓷砖', '600×600 防滑', '片', 100.00, 68.00, 6800.00, '红星美凯龙建材（上海供应处）', '已确认，等待付款', 0, 'CONFIRMED', 0, '2026-11-01 10:00:00', NULL, NULL, '2026-11-01 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (18, 1, 'PO-2026-0018', 4, 18, 108, 0, '2026-11-03', '水龙头', '摩恩 单孔', '个', 6.00, 460.00, 2760.00, '欧派家居上海经销商', '草稿，待确认', 0, 'DRAFT', 0, NULL, NULL, NULL, '2026-11-03 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (19, 1, 'PO-2026-0019', 4, 19, 108, 0, '2026-11-02', '整体橱柜', '石英石台面', '延米', 4.00, 1800.00, 7200.00, '欧派家居上海经销商', '已确认，等待付款', 0, 'CONFIRMED', 0, '2026-11-02 10:00:00', NULL, NULL, '2026-11-02 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (20, 1, 'PO-2026-0020', 5, 20, 108, 0, '2026-08-15', '实木复合门', '标准 樘', '樘', 4.00, 1280.00, 5120.00, '欧派家居上海经销商', '款已付，材料已入库', 0, 'PAID', 49, '2026-08-15 10:00:00', '2026-08-17 10:00:00', '2026-08-17 15:00:00', '2026-08-15 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (21, 1, 'PO-2026-0021', 5, 21, 105, 0, '2026-08-16', '门套线', '实木 6cm', '米', 60.00, 28.00, 1680.00, '红星美凯龙建材（上海供应处）', '款已付，材料已入库', 0, 'PAID', 50, '2026-08-16 10:00:00', '2026-08-18 10:00:00', '2026-08-18 15:00:00', '2026-08-16 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (22, 1, 'PO-2026-0022', 5, 22, 106, 0, '2026-08-10', '防水涂料', '20kg', '桶', 10.00, 380.00, 3800.00, '东方雨虹防水工程部', '款已付，材料已入库', 0, 'PAID', 47, '2026-08-10 10:00:00', '2026-08-12 10:00:00', '2026-08-12 15:00:00', '2026-08-10 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (23, 1, 'PO-2026-0023', 5, 23, 105, 0, '2026-08-12', '瓷砖', '600×600', '片', 120.00, 68.00, 8160.00, '红星美凯龙建材（上海供应处）', '款已付，材料已入库', 0, 'PAID', 48, '2026-08-12 10:00:00', '2026-08-14 10:00:00', '2026-08-14 15:00:00', '2026-08-12 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (24, 1, 'PO-2026-0024', 6, 25, 105, 0, '2026-06-12', '木纹铝方通', '30×80mm', '米', 90.00, 45.00, 4050.00, '红星美凯龙建材（上海供应处）', '款已付，材料已入库', 0, 'PAID', 34, '2026-06-12 10:00:00', '2026-06-14 10:00:00', '2026-06-14 15:00:00', '2026-06-12 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (25, 1, 'PO-2026-0025', 6, 26, 105, 0, '2026-06-18', '钢化玻璃隔断', '10mm 钢化', '㎡', 28.00, 260.00, 7280.00, '红星美凯龙建材（上海供应处）', '款已付，材料已入库', 0, 'PAID', 37, '2026-06-18 10:00:00', '2026-06-20 10:00:00', '2026-06-20 15:00:00', '2026-06-18 09:30:00', NULL);
INSERT INTO `project_purchase_order` VALUES (26, 1, 'PO-2026-0026', 6, 27, 105, 0, '2026-06-15', '办公地毯', '方块拼接地毯', '㎡', 60.00, 85.00, 5100.00, '红星美凯龙建材（上海供应处）', '款已付，材料已入库', 0, 'PAID', 35, '2026-06-15 10:00:00', '2026-06-17 10:00:00', '2026-06-17 15:00:00', '2026-06-15 09:30:00', NULL);

-- ----------------------------
-- Table structure for project_purchase_order_item
-- ----------------------------
DROP TABLE IF EXISTS `project_purchase_order_item`;
CREATE TABLE `project_purchase_order_item`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `purchase_order_id` bigint(20) UNSIGNED NOT NULL COMMENT '采购单ID',
  `material_id` bigint(20) NOT NULL COMMENT '材料清单ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '材料名称快照',
  `spec` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '规格快照',
  `unit` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '单位快照',
  `quantity` decimal(12, 2) NOT NULL DEFAULT 0.00 COMMENT '采购数量',
  `unit_price` decimal(12, 2) NOT NULL DEFAULT 0.00 COMMENT '采购单价',
  `amount` decimal(12, 2) NOT NULL DEFAULT 0.00 COMMENT '采购金额',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_purchase_item_order`(`purchase_order_id`) USING BTREE,
  INDEX `idx_purchase_item_material`(`material_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '项目采购单明细' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of project_purchase_order_item
-- ----------------------------

-- ----------------------------
-- Table structure for project_quote
-- ----------------------------
DROP TABLE IF EXISTS `project_quote`;
CREATE TABLE `project_quote`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `project_id` bigint(20) NOT NULL COMMENT '关联项目ID',
  `item_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '项目名称，如：拆墙、水电改造',
  `unit` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '单位，如：平米、米、项',
  `price` decimal(10, 2) NOT NULL DEFAULT 0.00 COMMENT '单价',
  `quantity` decimal(10, 2) NOT NULL DEFAULT 0.00 COMMENT '数量',
  `total` decimal(10, 2) NOT NULL DEFAULT 0.00 COMMENT '小计 = 单价 * 数量',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '备注，如：不含瓷砖',
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间，NULL表示未删除',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_project_id`(`project_id`) USING BTREE,
  INDEX `idx_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 254 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '项目报价明细表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of project_quote
-- ----------------------------
INSERT INTO `project_quote` VALUES (216, 1, 1, '墙体拆除与垃圾清运', '项', 4800.00, 1.00, 4800.00, '', '2026-01-05 10:00:00', '2026-01-05 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (217, 1, 1, '水电改造（含开槽、布线、打压）', '项', 26000.00, 1.00, 26000.00, '', '2026-01-05 10:00:00', '2026-01-05 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (218, 1, 1, '木工制作（吊顶、柜体、门套）', '项', 58000.00, 1.00, 58000.00, '', '2026-01-05 10:00:00', '2026-01-05 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (219, 1, 1, '墙面油漆（腻子、底漆、面漆）', '项', 32000.00, 1.00, 32000.00, '', '2026-01-05 10:00:00', '2026-01-05 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (220, 1, 1, '主材采购与安装', '项', 52000.00, 1.00, 52000.00, '', '2026-01-05 10:00:00', '2026-01-05 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (221, 1, 1, '全案设计费', '项', 13200.00, 1.00, 13200.00, '', '2026-01-05 10:00:00', '2026-01-05 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (222, 1, 2, '全案设计费', '项', 16800.00, 1.00, 16800.00, '', '2026-02-08 10:00:00', '2026-02-08 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (223, 1, 2, '基础施工（水电、泥木、油漆）', '项', 92000.00, 1.00, 92000.00, '', '2026-02-08 10:00:00', '2026-02-08 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (224, 1, 2, '主材（地板、瓷砖、卫浴）', '项', 108000.00, 1.00, 108000.00, '', '2026-02-08 10:00:00', '2026-02-08 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (225, 1, 2, '定制柜体（橱柜、衣柜、阳台柜）', '项', 38000.00, 1.00, 38000.00, '', '2026-02-08 10:00:00', '2026-02-08 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (226, 1, 2, '软装搭配（窗帘、灯具）', '项', 13200.00, 1.00, 13200.00, '', '2026-02-08 10:00:00', '2026-02-08 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (227, 1, 3, '拆除与垃圾清运', '项', 6000.00, 1.00, 6000.00, '', '2026-08-22 10:00:00', '2026-08-22 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (228, 1, 3, '水电改造（办公区强弱电）', '项', 18000.00, 1.00, 18000.00, '', '2026-08-22 10:00:00', '2026-08-22 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (229, 1, 3, '吊顶与轻钢隔墙', '项', 26000.00, 1.00, 26000.00, '', '2026-08-22 10:00:00', '2026-08-22 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (230, 1, 3, '墙面油漆与矿棉板', '项', 14000.00, 1.00, 14000.00, '', '2026-08-22 10:00:00', '2026-08-22 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (231, 1, 3, '前台形象墙与灯箱', '项', 20000.00, 1.00, 20000.00, '', '2026-08-22 10:00:00', '2026-08-22 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (232, 1, 3, '开荒保洁', '项', 4000.00, 1.00, 4000.00, '', '2026-08-22 10:00:00', '2026-08-22 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (233, 1, 3, '设计费', '项', 8000.00, 1.00, 8000.00, '', '2026-08-22 10:00:00', '2026-08-22 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (234, 1, 4, '拆除与垃圾清运', '项', 3000.00, 1.00, 3000.00, '', '2026-10-28 10:00:00', '2026-10-28 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (235, 1, 4, '防水施工（东方雨虹）', '项', 4800.00, 1.00, 4800.00, '', '2026-10-28 10:00:00', '2026-10-28 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (236, 1, 4, '水电改造', '项', 12000.00, 1.00, 12000.00, '', '2026-10-28 10:00:00', '2026-10-28 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (237, 1, 4, '墙地砖铺贴', '项', 18000.00, 1.00, 18000.00, '', '2026-10-28 10:00:00', '2026-10-28 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (238, 1, 4, '整体橱柜', '项', 22000.00, 1.00, 22000.00, '', '2026-10-28 10:00:00', '2026-10-28 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (239, 1, 4, '洁具与五金', '项', 6200.00, 1.00, 6200.00, '', '2026-10-28 10:00:00', '2026-10-28 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (240, 1, 4, '设计费', '项', 2000.00, 1.00, 2000.00, '', '2026-10-28 10:00:00', '2026-10-28 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (241, 1, 5, '全案设计费', '项', 8000.00, 1.00, 8000.00, '', '2026-07-24 10:00:00', '2026-07-24 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (242, 1, 5, '拆除与垃圾清运', '项', 5000.00, 1.00, 5000.00, '', '2026-07-24 10:00:00', '2026-07-24 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (243, 1, 5, '水电改造', '项', 21000.00, 1.00, 21000.00, '', '2026-07-24 10:00:00', '2026-07-24 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (244, 1, 5, '泥木施工', '项', 38000.00, 1.00, 38000.00, '', '2026-07-24 10:00:00', '2026-07-24 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (245, 1, 5, '墙面油漆', '项', 16000.00, 1.00, 16000.00, '', '2026-07-24 10:00:00', '2026-07-24 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (246, 1, 5, '主材采购与安装', '项', 30000.00, 1.00, 30000.00, '', '2026-07-24 10:00:00', '2026-07-24 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (247, 1, 6, '设计费', '项', 4000.00, 1.00, 4000.00, '', '2026-05-22 10:00:00', '2026-05-22 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (248, 1, 6, '拆除与垃圾清运', '项', 4000.00, 1.00, 4000.00, '', '2026-05-22 10:00:00', '2026-05-22 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (249, 1, 6, '水电改造', '项', 12000.00, 1.00, 12000.00, '', '2026-05-22 10:00:00', '2026-05-22 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (250, 1, 6, '墙面装修', '项', 14000.00, 1.00, 14000.00, '', '2026-05-22 10:00:00', '2026-05-22 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (251, 1, 6, '地面与玻璃隔断', '项', 16000.00, 1.00, 16000.00, '', '2026-05-22 10:00:00', '2026-05-22 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (252, 1, 6, '门头招牌与灯箱', '项', 18000.00, 1.00, 18000.00, '', '2026-05-22 10:00:00', '2026-05-22 10:00:00', NULL);
INSERT INTO `project_quote` VALUES (253, 1, 6, '灯具与安装', '项', 8000.00, 1.00, 8000.00, '', '2026-05-22 10:00:00', '2026-05-22 10:00:00', NULL);

-- ----------------------------
-- Table structure for site_user
-- ----------------------------
DROP TABLE IF EXISTS `site_user`;
CREATE TABLE `site_user`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '登录账号',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '密码（bcrypt加密）',
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '昵称/姓名',
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机号',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '头像URL',
  `role` tinyint(3) UNSIGNED NOT NULL DEFAULT 3 COMMENT '角色：1=管理员 2=员工 3=商家',
  `status` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态：0=禁用 1=正常',
  `merchant_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '所属商家ID（商家自身或员工归属）',
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `account_created` tinyint(3) UNSIGNED NOT NULL DEFAULT 1 COMMENT '是否已开通独立登录账号',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  `last_login_at` datetime NULL DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '最后登录IP',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间',
  `department_id` int(11) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_username`(`username`) USING BTREE,
  INDEX `idx_role_status`(`role`, `status`) USING BTREE,
  INDEX `idx_merchant`(`merchant_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 10 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '系统用户表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of site_user
-- ----------------------------
INSERT INTO `site_user` VALUES (1, 'admin', '$2a$10$EBJtupiv1TELywksMaWAFOdC7/gI6k7FbZo3GHTy94J/KwhUolyfi', '超级管理员', '', '', 1, 1, 0, 1, 1, '', '2026-11-06 08:35:46', '127.0.0.1', '2026-09-22 09:58:11', '2026-11-06 08:35:46', NULL, NULL);
INSERT INTO `site_user` VALUES (2, 'wangjianguo', '$2a$10$EBJtupiv1TELywksMaWAFOdC7/gI6k7FbZo3GHTy94J/KwhUolyfi', '王建国', '13900001111', '', 2, 1, 0, 1, 1, '演示数据：项目经理', NULL, '', '2026-09-30 09:30:05', '2026-11-04 17:40:14', NULL, 6);
INSERT INTO `site_user` VALUES (3, 'liqiang', '$2a$10$EBJtupiv1TELywksMaWAFOdC7/gI6k7FbZo3GHTy94J/KwhUolyfi', '李强', '13900002222', '', 2, 1, 0, 1, 1, '演示数据：施工队长', '2026-11-04 07:32:22', '127.0.0.1', '2026-09-30 09:30:05', '2026-11-04 17:40:14', NULL, 6);
INSERT INTO `site_user` VALUES (4, 'zhaoming', '$2a$10$EBJtupiv1TELywksMaWAFOdC7/gI6k7FbZo3GHTy94J/KwhUolyfi', '赵明', '13900003333', '', 2, 1, 0, 1, 1, '演示数据：项目经理', '2026-11-04 08:17:47', '127.0.0.1', '2026-09-30 09:30:06', '2026-11-04 17:40:14', NULL, 7);
INSERT INTO `site_user` VALUES (5, 'chenxiaodong', '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92', '陈晓东', '13900004444', '', 2, 1, 0, 1, 1, '演示数据：工程监理', NULL, '', '2026-09-30 09:30:06', '2026-11-04 17:40:14', NULL, 8);
INSERT INTO `site_user` VALUES (7, 'lisi', '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92', '李四', '18810589665', '', 2, 1, 0, 1, 1, '', NULL, '', '2026-11-03 16:06:19', '2026-11-04 17:40:14', NULL, 6);
INSERT INTO `site_user` VALUES (8, 'wuxiaofeng', '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92', '吴晓峰', '13766668888', '', 2, 1, 0, 1, 1, '负责水电与施工统筹', NULL, '', '2026-11-04 17:39:04', '2026-11-04 17:39:04', NULL, 6);
INSERT INTO `site_user` VALUES (9, '13777779999', '779999', '郑雨欣', '13777779999', '', 2, 1, 0, 1, 0, '线上渠道获客', NULL, '', '2026-11-04 17:39:04', '2026-11-04 17:39:04', NULL, 7);

-- ----------------------------
-- Table structure for site_user_token
-- ----------------------------
DROP TABLE IF EXISTS `site_user_token`;
CREATE TABLE `site_user_token`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` int(10) UNSIGNED NOT NULL COMMENT '用户ID，对应 site_user.id',
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `token` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'Token（JWT的jti或随机串）',
  `client` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '客户端类型：web/app',
  `ip` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '登录IP',
  `expire_at` datetime NOT NULL COMMENT '过期时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_token`(`token`) USING BTREE,
  INDEX `idx_user`(`user_id`) USING BTREE,
  INDEX `idx_expire`(`expire_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 132 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户登录令牌表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of site_user_token
-- ----------------------------
INSERT INTO `site_user_token` VALUES (93, 3, 1, '12b0ba09d70abf6e63b60809d1aa934aa57f45ee7f94e19c4e87fbee8013de7f', 'web', '127.0.0.1', '2026-11-11 07:32:22', '2026-11-04 07:32:22', '2026-11-04 07:32:22');
INSERT INTO `site_user_token` VALUES (100, 4, 1, '5bb78b5a54a8d776e715d23134caf1015801723d98791cb2eb51ba2aad1eacc3', 'web', '127.0.0.1', '2026-11-11 08:17:47', '2026-11-04 08:17:47', '2026-11-04 08:17:47');
INSERT INTO `site_user_token` VALUES (131, 1, 1, 'e59fbf2729871c21ba9960b322abb7d32a87bd7e26ec3b1df8ccd6a6a4f0b1b2', 'web', '127.0.0.1', '2026-11-13 08:35:46', '2026-11-06 08:35:46', '2026-11-06 08:35:46');

-- ----------------------------
-- Table structure for sys_department
-- ----------------------------
DROP TABLE IF EXISTS `sys_department`;
CREATE TABLE `sys_department`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `name` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `code` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `sort_order` int(11) NOT NULL DEFAULT 0,
  `status` tinyint(3) UNSIGNED NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_department_company_code`(`company_id`, `code`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 10 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '公司部门' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_department
-- ----------------------------
INSERT INTO `sys_department` VALUES (6, 1, '工程部', 'GC', 1, 1, '2026-11-04 17:39:04', '2026-11-04 17:39:04');
INSERT INTO `sys_department` VALUES (7, 1, '市场部', 'SC', 2, 1, '2026-11-04 17:39:04', '2026-11-04 17:39:04');
INSERT INTO `sys_department` VALUES (8, 1, '设计部', 'SJ', 3, 1, '2026-11-04 17:39:04', '2026-11-04 17:39:04');
INSERT INTO `sys_department` VALUES (9, 1, '财务部', 'CW', 4, 1, '2026-11-04 17:39:04', '2026-11-04 17:39:04');

-- ----------------------------
-- Table structure for sys_department_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_department_menu`;
CREATE TABLE `sys_department_menu`  (
  `department_id` int(10) UNSIGNED NOT NULL,
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `menu_id` int(10) UNSIGNED NOT NULL,
  PRIMARY KEY (`department_id`, `menu_id`) USING BTREE,
  INDEX `idx_department_menu_menu`(`menu_id`) USING BTREE,
  CONSTRAINT `fk_department_menu_department` FOREIGN KEY (`department_id`) REFERENCES `sys_department` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `fk_department_menu_menu` FOREIGN KEY (`menu_id`) REFERENCES `sys_menu` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '部门菜单授权' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_department_menu
-- ----------------------------

-- ----------------------------
-- Table structure for sys_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_menu`;
CREATE TABLE `sys_menu`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `parent_id` int(10) UNSIGNED NOT NULL DEFAULT 0,
  `title` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `route_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `route_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `component` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `icon` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `sort_order` int(11) NOT NULL DEFAULT 0,
  `enabled` tinyint(1) NOT NULL DEFAULT 1,
  `hide_in_menu` tinyint(1) NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_menu_route_name`(`route_name`) USING BTREE,
  INDEX `idx_menu_parent_order`(`parent_id`, `sort_order`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 32 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '后台动态菜单' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_menu
-- ----------------------------
INSERT INTO `sys_menu` VALUES (1, 0, '首页', 'home', '/home', 'layout.base$view.home', 'mdi:home-outline', 1, 1, 0, '2026-11-04 16:10:07', '2026-11-04 16:10:07');
INSERT INTO `sys_menu` VALUES (2, 0, '项目管理', 'project', '/project', 'layout.base', 'mdi:folder-multiple-outline', 2, 1, 0, '2026-11-04 16:10:07', '2026-11-04 16:10:07');
INSERT INTO `sys_menu` VALUES (3, 0, '财务管理', 'finance', '/finance', 'layout.base', 'mdi:cash-multiple', 3, 1, 0, '2026-11-04 16:10:07', '2026-11-04 16:10:07');
INSERT INTO `sys_menu` VALUES (4, 0, '数据报表', 'report', '/report', 'layout.base', 'mdi:chart-box-outline', 6, 1, 0, '2026-11-04 16:10:07', '2026-11-06 09:14:19');
INSERT INTO `sys_menu` VALUES (5, 0, '库存查询', 'inventory', '/inventory', 'layout.base', 'mdi:package-variant-closed', 4, 1, 0, '2026-11-04 16:10:07', '2026-11-06 09:14:11');
INSERT INTO `sys_menu` VALUES (6, 0, '企业管理', 'enterprise', '/enterprise', 'layout.base', 'mdi:office-building-outline', 5, 1, 0, '2026-11-04 16:10:07', '2026-11-06 09:14:15');
INSERT INTO `sys_menu` VALUES (7, 0, '基础设置', 'settings', '/settings', 'layout.base', 'mdi:cog-outline', 7, 1, 0, '2026-11-04 16:10:07', '2026-11-04 16:13:43');
INSERT INTO `sys_menu` VALUES (8, 2, '项目列表', 'project_project-list', '/project/project-list', 'view.project_project-list', 'mdi:format-list-bulleted', 1, 1, 0, '2026-11-04 16:10:07', '2026-11-04 16:10:07');
INSERT INTO `sys_menu` VALUES (9, 2, '项目资料', 'project_data', '/project/data', 'view.project_data', 'mdi:file-document-outline', 2, 1, 0, '2026-11-04 16:10:07', '2026-11-04 16:10:07');
INSERT INTO `sys_menu` VALUES (10, 2, '报价管理', 'project_quote', '/project/quote', 'view.project_quote', 'mdi:file-sign', 3, 1, 0, '2026-11-04 16:10:07', '2026-11-04 16:10:07');
INSERT INTO `sys_menu` VALUES (11, 3, '收款管理', 'finance_payment', '/finance/payment', 'view.finance_payment', 'mdi:cash-check', 1, 1, 0, '2026-11-04 16:10:07', '2026-11-04 16:10:07');
INSERT INTO `sys_menu` VALUES (12, 3, '付款管理', 'finance_payout', '/finance/payout', 'view.finance_payout', 'mdi:cash-minus', 2, 1, 0, '2026-11-04 16:10:07', '2026-11-04 16:10:07');
INSERT INTO `sys_menu` VALUES (13, 3, '记账凭证', 'finance_voucher', '/finance/voucher', 'view.finance_voucher', 'mdi:receipt-text-outline', 3, 1, 0, '2026-11-04 16:10:07', '2026-11-04 16:10:07');
INSERT INTO `sys_menu` VALUES (14, 3, '财务基础资料', 'finance_basic', '/finance/basic', 'view.finance_basic', 'mdi:database-cog-outline', 4, 1, 0, '2026-11-04 16:10:07', '2026-11-04 16:10:07');
INSERT INTO `sys_menu` VALUES (15, 4, '经营总览', 'report_overview', '/report/overview', 'view.report_overview', 'mdi:view-dashboard-outline', 1, 1, 0, '2026-11-04 16:10:07', '2026-11-04 16:10:07');
INSERT INTO `sys_menu` VALUES (16, 4, '资金流水', 'report_flows', '/report/flows', 'view.report_flows', 'mdi:swap-horizontal', 2, 1, 0, '2026-11-04 16:10:07', '2026-11-04 16:10:07');
INSERT INTO `sys_menu` VALUES (17, 4, '科目余额', 'report_balance', '/report/balance', 'view.report_balance', 'mdi:scale-balance', 3, 1, 0, '2026-11-04 16:10:07', '2026-11-04 16:10:07');
INSERT INTO `sys_menu` VALUES (18, 4, '项目利润', 'report_profit', '/report/profit', 'view.report_profit', 'mdi:trending-up', 4, 1, 0, '2026-11-04 16:10:07', '2026-11-04 16:10:07');
INSERT INTO `sys_menu` VALUES (19, 5, '材料清单', 'inventory_materials', '/inventory/materials', 'view.inventory_materials', 'mdi:package-variant', 1, 1, 0, '2026-11-04 16:10:07', '2026-11-04 16:10:07');
INSERT INTO `sys_menu` VALUES (20, 5, '领取记录', 'inventory_usages', '/inventory/usages', 'view.inventory_usages', 'mdi:clipboard-text-outline', 2, 1, 0, '2026-11-04 16:10:07', '2026-11-04 16:10:07');
INSERT INTO `sys_menu` VALUES (21, 5, '采购管理', 'inventory_procurement', '/inventory/procurement', 'view.inventory_procurement', 'mdi:cart-outline', 3, 1, 0, '2026-11-04 16:10:07', '2026-11-04 16:10:07');
INSERT INTO `sys_menu` VALUES (22, 6, '员工管理', 'enterprise_staff', '/enterprise/staff', 'view.enterprise_staff', 'mdi:account-group-outline', 1, 1, 0, '2026-11-04 16:10:07', '2026-11-04 16:10:07');
INSERT INTO `sys_menu` VALUES (23, 6, '部门管理', 'enterprise_department', '/enterprise/department', 'view.enterprise_department', 'mdi:sitemap-outline', 2, 1, 0, '2026-11-04 16:10:07', '2026-11-04 16:10:07');
INSERT INTO `sys_menu` VALUES (24, 6, '账号管理', 'enterprise_user', '/enterprise/user', 'view.enterprise_user', 'mdi:account-key-outline', 3, 1, 0, '2026-11-04 16:10:07', '2026-11-04 16:10:07');
INSERT INTO `sys_menu` VALUES (25, 7, '客户管理', 'settings_customer', '/settings/customer', 'view.settings_customer', 'mdi:account-tie-outline', 1, 1, 0, '2026-11-04 16:10:07', '2026-11-04 16:10:07');
INSERT INTO `sys_menu` VALUES (26, 7, '供应商', 'settings_supplier', '/settings/supplier', 'view.settings_supplier', 'mdi:truck-outline', 2, 1, 0, '2026-11-04 16:10:07', '2026-11-04 16:10:07');
INSERT INTO `sys_menu` VALUES (27, 7, '菜单管理', 'settings_menu', '/settings/menu', 'view.settings_menu', 'mdi:menu', 3, 1, 0, '2026-11-04 16:10:07', '2026-11-04 16:10:07');
INSERT INTO `sys_menu` VALUES (28, 5, '库存管理', 'inventory_stock', '/inventory/stock', 'view.inventory_stock', 'mdi:warehouse', 4, 1, 0, '2026-11-05 09:33:30', '2026-11-05 09:33:30');
INSERT INTO `sys_menu` VALUES (29, 5, '仓库管理', 'inventory_warehouse', '/inventory/warehouse', 'view.inventory_warehouse', 'mdi:garage', 5, 1, 0, '2026-11-05 10:27:44', '2026-11-05 13:53:15');
INSERT INTO `sys_menu` VALUES (30, 7, '打印设置', 'settings_print', '/settings/print', 'view.settings_print', 'mdi:printer-settings', 4, 1, 0, '2026-11-05 11:31:54', '2026-11-06 09:15:26');
INSERT INTO `sys_menu` VALUES (31, 4, '操作日志', 'report_operation-log', '/report/operation-log', 'view.report_operation-log', 'mdi:history', 5, 1, 0, '2026-11-06 16:28:23', '2026-11-06 16:34:35');

-- ----------------------------
-- Table structure for sys_operation_log
-- ----------------------------
DROP TABLE IF EXISTS `sys_operation_log`;
CREATE TABLE `sys_operation_log`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `user_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '操作人ID，0 表示系统',
  `user_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '操作人姓名（冗余保存，人员删除后仍可追溯）',
  `module` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '业务模块：收款/付款/项目/库存/客户/供应商',
  `action` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '动作：CREATE/UPDATE/DELETE',
  `target_type` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '对象类型，如 payment / payout',
  `target_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '对象ID',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '操作摘要',
  `amount` decimal(14, 2) NOT NULL DEFAULT 0.00 COMMENT '涉及金额，无金额时为 0',
  `detail` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '变更详情（JSON）',
  `ip` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '操作来源IP',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_oplog_company_created`(`company_id`, `created_at`) USING BTREE,
  INDEX `idx_oplog_module_action`(`company_id`, `module`, `action`) USING BTREE,
  INDEX `idx_oplog_target`(`target_type`, `target_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 8 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '操作日志表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sys_operation_log
-- ----------------------------
INSERT INTO `sys_operation_log` VALUES (3, 1, 1, '超级管理员', '收款', 'CREATE', 'payment', 18, '新增收款：黄建军 · 国定路厨卫改造（中期款 ￥3000.00）', 3000.00, '{\"amount\":3000,\"paymentDate\":\"2026-11-06\",\"paymentId\":18,\"paymentType\":\"中期款\",\"projectId\":4,\"projectName\":\"黄建军 · 国定路厨卫改造\",\"remark\":\"现场测量增项预收\",\"voucherNo\":\"收-202611-0007\"}', '127.0.0.1', '2026-11-06 08:36:47');
INSERT INTO `sys_operation_log` VALUES (4, 1, 1, '超级管理员', '收款', 'UPDATE', 'payment', 18, '修改收款：黄建军 · 国定路厨卫改造（中期款 ￥3800.00）', 3800.00, '{\"amount\":3800,\"paymentDate\":\"2026-11-06\",\"paymentId\":18,\"paymentType\":\"中期款\",\"projectId\":4,\"projectName\":\"黄建军 · 国定路厨卫改造\",\"remark\":\"现场测量增项预收（调整）\",\"voucherNo\":\"收-202611-0008\"}', '127.0.0.1', '2026-11-06 08:37:00');
INSERT INTO `sys_operation_log` VALUES (5, 1, 1, '超级管理员', '收款', 'DELETE', 'payment', 18, '删除收款：黄建军 · 国定路厨卫改造（中期款 ￥3800.00）', 3800.00, '{\"amount\":3800,\"paymentDate\":\"2026-11-06\",\"paymentId\":18,\"paymentType\":\"中期款\",\"projectId\":4,\"projectName\":\"黄建军 · 国定路厨卫改造\",\"remark\":\"现场测量增项预收（调整）\",\"voucherNo\":\"收-202611-0008\"}', '127.0.0.1', '2026-11-06 08:37:00');
INSERT INTO `sys_operation_log` VALUES (6, 1, 1, '超级管理员', '付款', 'CREATE', 'payout', 69, '新增付款：黄建军 · 国定路厨卫改造（材料款 付 红星美凯龙建材（上海供应处） ￥500.00）', 500.00, '{\"amount\":500,\"payoutDate\":\"2026-11-06\",\"payoutId\":69,\"payoutType\":\"材料款\",\"projectId\":4,\"projectName\":\"黄建军 · 国定路厨卫改造\",\"remark\":\"联调测试付款\",\"vendorId\":105,\"vendorName\":\"红星美凯龙建材（上海供应处）\",\"voucherNo\":\"付-202611-0009\"}', '127.0.0.1', '2026-11-06 08:37:19');
INSERT INTO `sys_operation_log` VALUES (7, 1, 1, '超级管理员', '付款', 'DELETE', 'payout', 69, '删除付款：黄建军 · 国定路厨卫改造（材料款 付 红星美凯龙建材（上海供应处） ￥500.00）', 500.00, '{\"amount\":500,\"payoutDate\":\"2026-11-06\",\"payoutId\":69,\"payoutType\":\"材料款\",\"projectId\":4,\"projectName\":\"黄建军 · 国定路厨卫改造\",\"remark\":\"联调测试付款\",\"vendorId\":105,\"vendorName\":\"红星美凯龙建材（上海供应处）\",\"voucherNo\":\"付-202611-0009\"}', '127.0.0.1', '2026-11-06 08:37:19');

-- ----------------------------
-- Table structure for vendor
-- ----------------------------
DROP TABLE IF EXISTS `vendor`;
CREATE TABLE `vendor`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '供应商名称',
  `contact_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '联系人',
  `phone` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '联系电话',
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '地址',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_vendor_name`(`name`) USING BTREE,
  INDEX `idx_vendor_phone`(`phone`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 118 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '供应商表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of vendor
-- ----------------------------
INSERT INTO `vendor` VALUES (105, 1, '红星美凯龙建材（上海供应处）', '马经理', '021-54886632', '上海市普陀区真北路1108号', '瓷砖/水泥辅料主供应商', '2026-01-03 09:00:00', '2026-01-03 09:00:00', NULL);
INSERT INTO `vendor` VALUES (106, 1, '东方雨虹防水工程部', '孙工', '13901234567', '上海市闵行区莲花南路288号', '防水施工专业分包', '2026-01-03 09:00:00', '2026-01-03 09:00:00', NULL);
INSERT INTO `vendor` VALUES (107, 1, '沪上名家水电安装队', '周队长', '13700137001', '上海市宝山区沪太路4365号', '水电/瓦工/油漆劳务', '2026-01-03 09:00:00', '2026-01-03 09:00:00', NULL);
INSERT INTO `vendor` VALUES (108, 1, '欧派家居上海经销商', '林经理', '021-62985577', '上海市静安区汶水路401号', '橱柜/定制柜体', '2026-01-03 09:00:00', '2026-01-03 09:00:00', NULL);
INSERT INTO `vendor` VALUES (109, 1, '立邦涂料上海总代理', '何经理', '021-56082233', '上海市嘉定区曹安公路4800号', '乳胶漆/墙面材料', '2026-01-03 09:00:00', '2026-01-03 09:00:00', NULL);
INSERT INTO `vendor` VALUES (110, 1, '上海汇智物业管理有限公司', '前台', '021-62550099', '上海市静安区江宁路200号', '办公室租赁物业', '2026-01-03 09:00:00', '2026-01-03 09:00:00', NULL);
INSERT INTO `vendor` VALUES (111, 1, '上海巨量引擎广告服务有限公司', '客户经理', '021-61350066', '上海市杨浦区淞沪路433号', '抖音本地推投放', '2026-01-03 09:00:00', '2026-01-03 09:00:00', NULL);
INSERT INTO `vendor` VALUES (112, 1, '上海素合空间设计工作室', '顾设计', '021-62889102', '上海市静安区愚园路546号', '全案设计/效果图外包', '2026-01-03 09:00:00', '2026-01-03 09:00:00', NULL);
INSERT INTO `vendor` VALUES (113, 1, '上海安居物流配送有限公司', '刘师傅', '13818886677', '上海市松江区九亭镇沪松公路1288号', '建材配送/垃圾清运', '2026-01-03 09:00:00', '2026-01-03 09:00:00', NULL);
INSERT INTO `vendor` VALUES (114, 1, '上海启明广告制作有限公司', '张经理', '021-64370055', '上海市徐汇区宜山路700号', '灯箱/门头/广告物料制作', '2026-01-03 09:00:00', '2026-01-03 09:00:00', NULL);

-- ----------------------------
-- Table structure for warehouse
-- ----------------------------
DROP TABLE IF EXISTS `warehouse`;
CREATE TABLE `warehouse`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `company_id` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '仓库名称',
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '仓库编号',
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '仓库地址',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  `is_default` tinyint(1) UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否默认仓库 1-是 0-否',
  `sort` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '排序（越小越靠前）',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_warehouse_company`(`company_id`) USING BTREE,
  INDEX `idx_warehouse_name`(`name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '仓库表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of warehouse
-- ----------------------------
INSERT INTO `warehouse` VALUES (1, 1, '默认仓库', 'WH-001', '上海市松江区九亭镇沪松公路1288号', '系统初始仓库', 1, 0, '2026-01-01 09:00:00', '2026-01-01 09:00:00', NULL);

ALTER TABLE `material_catalog`
  ADD COLUMN `tax_mode` varchar(20) NOT NULL DEFAULT 'NO_INVOICE' COMMENT '默认税务方式' AFTER `unit_price`,
  ADD COLUMN `tax_rate` decimal(5,2) NOT NULL DEFAULT 0.00 COMMENT '默认税率（百分比）' AFTER `tax_mode`;

ALTER TABLE `project_material`
  ADD COLUMN `tax_mode` varchar(20) NOT NULL DEFAULT 'NO_INVOICE' COMMENT '税务方式' AFTER `unit_price`,
  ADD COLUMN `tax_rate` decimal(5,2) NOT NULL DEFAULT 0.00 COMMENT '税率（百分比）' AFTER `tax_mode`;

ALTER TABLE `project_purchase_order`
  ADD COLUMN `tax_mode` varchar(20) NOT NULL DEFAULT 'NO_INVOICE' COMMENT '税务方式快照' AFTER `unit_price`,
  ADD COLUMN `tax_rate` decimal(5,2) NOT NULL DEFAULT 0.00 COMMENT '税率快照（百分比）' AFTER `tax_mode`,
  ADD COLUMN `net_amount` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '未税金额' AFTER `tax_rate`,
  ADD COLUMN `tax_amount` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '税额' AFTER `net_amount`;

UPDATE `project_purchase_order` SET `net_amount`=`amount`, `tax_amount`=0 WHERE `tax_mode`='NO_INVOICE';

SET FOREIGN_KEY_CHECKS = 1;
