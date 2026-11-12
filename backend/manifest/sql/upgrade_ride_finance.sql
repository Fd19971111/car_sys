-- 自营打车财务：行程司机成本与司机结算
-- 仅增加结构，不删除历史财务表、凭证或流水。
SET NAMES utf8mb4;

SET @ride_order_has_driver_fare = (
  SELECT COUNT(*)
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'ride_order'
    AND COLUMN_NAME = 'driver_fare'
);
SET @ride_order_driver_fare_sql = IF(
  @ride_order_has_driver_fare = 0,
  'ALTER TABLE `ride_order` ADD COLUMN `driver_fare` decimal(12,2) NOT NULL DEFAULT 0.00 AFTER `final_fare`',
  'SELECT 1'
);
PREPARE ride_order_driver_fare_stmt FROM @ride_order_driver_fare_sql;
EXECUTE ride_order_driver_fare_stmt;
DEALLOCATE PREPARE ride_order_driver_fare_stmt;

SET @payment_has_order_id = (
  SELECT COUNT(*)
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'payment'
    AND COLUMN_NAME = 'ride_order_id'
);
SET @payment_order_id_sql = IF(
  @payment_has_order_id = 0,
  'ALTER TABLE `payment` ADD COLUMN `ride_order_id` bigint(20) unsigned NOT NULL DEFAULT 0 AFTER `company_id`',
  'SELECT 1'
);
PREPARE payment_order_id_stmt FROM @payment_order_id_sql;
EXECUTE payment_order_id_stmt;
DEALLOCATE PREPARE payment_order_id_stmt;

SET @payment_order_index_exists = (
  SELECT COUNT(*)
  FROM information_schema.STATISTICS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'payment'
    AND INDEX_NAME = 'idx_payment_company_order'
);
SET @payment_order_index_sql = IF(
  @payment_order_index_exists = 0,
  'CREATE INDEX `idx_payment_company_order` ON `payment` (`company_id`,`ride_order_id`,`deleted_at`)',
  'SELECT 1'
);
PREPARE payment_order_index_stmt FROM @payment_order_index_sql;
EXECUTE payment_order_index_stmt;
DEALLOCATE PREPARE payment_order_index_stmt;

CREATE TABLE IF NOT EXISTS `ride_driver_earning` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '司机应得明细ID',
  `company_id` int(10) unsigned NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `order_id` bigint(20) unsigned NOT NULL COMMENT '行程订单ID',
  `driver_id` bigint(20) unsigned NOT NULL COMMENT '司机ID',
  `fare_amount` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '订单车费快照（元）',
  `driver_amount` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '司机应得车费（元）',
  `earning_date` date NOT NULL COMMENT '应得日期',
  `status` varchar(16) NOT NULL DEFAULT 'POSTED' COMMENT '状态：POSTED有效 VOID已冲销',
  `voucher_id` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '应计成本凭证ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ride_driver_earning_order` (`company_id`,`order_id`),
  KEY `idx_ride_driver_earning_driver_date` (`company_id`,`driver_id`,`earning_date`),
  KEY `idx_ride_driver_earning_status` (`company_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='行程司机应得车费';

CREATE TABLE IF NOT EXISTS `ride_driver_settlement` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '司机结算单ID',
  `company_id` int(10) unsigned NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `driver_id` bigint(20) unsigned NOT NULL COMMENT '司机ID',
  `settlement_no` varchar(40) NOT NULL COMMENT '结算单号',
  `settlement_date` date NOT NULL COMMENT '付款日期',
  `amount` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '本次结算金额（元）',
  `fund_account_id` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '付款资金账户ID',
  `voucher_id` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '冲减司机应付凭证ID',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `status` varchar(16) NOT NULL DEFAULT 'POSTED' COMMENT '状态：POSTED有效 VOID已作废',
  `created_by` int(10) unsigned NOT NULL DEFAULT 0 COMMENT '制单人',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ride_driver_settlement_no` (`company_id`,`settlement_no`),
  KEY `idx_ride_driver_settlement_driver_date` (`company_id`,`driver_id`,`settlement_date`),
  KEY `idx_ride_driver_settlement_status` (`company_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='司机车费结算单';

CREATE TABLE IF NOT EXISTS `ride_driver_settlement_item` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '结算分摊明细ID',
  `company_id` int(10) unsigned NOT NULL DEFAULT 1 COMMENT '所属公司ID',
  `settlement_id` bigint(20) unsigned NOT NULL COMMENT '司机结算单ID',
  `earning_id` bigint(20) unsigned NOT NULL COMMENT '司机应得明细ID',
  `amount` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '本次结算分摊金额（元）',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_ride_driver_settlement_item_settlement` (`company_id`,`settlement_id`),
  KEY `idx_ride_driver_settlement_item_earning` (`company_id`,`earning_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='司机结算分摊明细';

INSERT INTO `sys_menu`
  (`id`,`parent_id`,`title`,`route_name`,`route_path`,`component`,`icon`,`sort_order`,`enabled`,`hide_in_menu`)
VALUES
  (38,3,'司机结算','finance_driver_settlement','/finance/driver-settlement','view.finance_driver_settlement','mdi:account-cash-outline',4,1,0)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),`title`=VALUES(`title`),`route_name`=VALUES(`route_name`),
  `route_path`=VALUES(`route_path`),`component`=VALUES(`component`),`icon`=VALUES(`icon`),
  `sort_order`=VALUES(`sort_order`),`enabled`=VALUES(`enabled`),`hide_in_menu`=VALUES(`hide_in_menu`);

ALTER TABLE `sys_menu` AUTO_INCREMENT = 39;