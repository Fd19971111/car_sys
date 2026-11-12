-- ============================================================
-- 2026-11-12 订单中心上线 / 项目管理下线
-- 执行顺序：旧库先导入 yiminjia.sql / 已运行库直接执行本脚本
-- 1. 新建 ride_order 订单表
-- 2. 删除项目管理、报价、项目材料/采购/库存业务表
-- 3. 财务收付款表保留，去掉项目关联字段（项目已下线）
-- 4. 菜单：删除项目管理、项目利润，新增订单中心/全部订单/行程监控
-- ============================================================
SET NAMES utf8mb4;

-- ---------- 1. 订单表 ----------
CREATE TABLE IF NOT EXISTS `ride_order` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `company_id` int(10) unsigned NOT NULL DEFAULT '1' COMMENT '所属公司ID',
  `order_number` varchar(32) NOT NULL COMMENT '订单编号',
  `customer_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '客户ID',
  `customer_name` varchar(50) NOT NULL DEFAULT '' COMMENT '客户姓名快照',
  `customer_phone` varchar(20) NOT NULL DEFAULT '' COMMENT '客户手机号快照',
  `status` varchar(30) NOT NULL DEFAULT 'quoting' COMMENT '订单状态：quoting待报价 awaiting-confirmation待确认 dispatching待派单 driver-en-route司机前往中 in-progress行程中 completed已完成 cancelled已取消',
  `passenger_count` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '乘车人数',
  `vehicle_category` varchar(20) NOT NULL DEFAULT 'comfort' COMMENT '车型：economy comfort business six-seat',
  `pickup_address` varchar(255) NOT NULL DEFAULT '' COMMENT '上车地点',
  `destination_address` varchar(255) NOT NULL DEFAULT '' COMMENT '目的地点',
  `estimated_distance_km` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '预计里程（公里）',
  `estimated_duration_min` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '预计时长（分钟）',
  `quoted_fare` decimal(12,2) NOT NULL DEFAULT '0.00' COMMENT '客户报价（元）',
  `final_fare` decimal(12,2) NOT NULL DEFAULT '0.00' COMMENT '最终费用（元）',
  `customer_response` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '客户反馈：pending等待回复 accepted接受报价 rejected拒绝报价',
  `scheduled_at` datetime DEFAULT NULL COMMENT '预约用车时间',
  `driver_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '指派司机ID',
  `vehicle_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '指派车辆ID',
  `completed_at` datetime DEFAULT NULL COMMENT '行程结束时间',
  `internal_notes` varchar(500) NOT NULL DEFAULT '' COMMENT '内部备注',
  `created_by` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '创建人ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_ride_order_company_number` (`company_id`, `order_number`) USING BTREE,
  KEY `idx_ride_order_company_status` (`company_id`, `status`) USING BTREE,
  KEY `idx_ride_order_scheduled` (`company_id`, `scheduled_at`) USING BTREE,
  KEY `idx_ride_order_customer` (`customer_id`) USING BTREE,
  KEY `idx_ride_order_driver` (`driver_id`) USING BTREE,
  KEY `idx_ride_order_vehicle` (`vehicle_id`) USING BTREE,
  KEY `idx_ride_order_deleted` (`company_id`, `deleted_at`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='行程订单';

-- ---------- 2. 财务表去项目化 ----------
RENAME TABLE `project_payment` TO `payment`, `project_payout` TO `payout`;

ALTER TABLE `payment` DROP COLUMN `project_id`;

ALTER TABLE `payout`
  DROP COLUMN `project_id`,
  DROP COLUMN `purchase_order_id`,
  DROP COLUMN `vendor_id`,
  ADD COLUMN `vendor_name` varchar(100) NOT NULL DEFAULT '' COMMENT '收款方名称' AFTER `payout_date`;

-- ---------- 3. 删除项目管理 / 报价 / 材料采购库存表 ----------
DROP TABLE IF EXISTS `project_material_usage`;
DROP TABLE IF EXISTS `project_material`;
DROP TABLE IF EXISTS `project_purchase_order_item`;
DROP TABLE IF EXISTS `project_purchase_order`;
DROP TABLE IF EXISTS `inventory_record`;
DROP TABLE IF EXISTS `project_log`;
DROP TABLE IF EXISTS `project_change`;
DROP TABLE IF EXISTS `project_contract`;
DROP TABLE IF EXISTS `project_quote`;
DROP TABLE IF EXISTS `warehouse`;
DROP TABLE IF EXISTS `material_catalog`;
DROP TABLE IF EXISTS `project`;

-- ---------- 4. 菜单替换 ----------
DELETE dm
FROM sys_department_menu dm
INNER JOIN sys_menu m ON m.id = dm.menu_id
WHERE m.id IN (2, 18) OR m.parent_id IN (2, 18);

DELETE FROM sys_menu WHERE parent_id IN (2, 18);
DELETE FROM sys_menu WHERE id IN (2, 18);

INSERT INTO sys_menu
  (id, parent_id, title, route_name, route_path, component, icon, sort_order, enabled, hide_in_menu)
VALUES
  (35, 0,  '订单中心', 'order',         '/order',         'layout.base',        'mdi:clipboard-text-clock-outline', 2, 1, 0),
  (36, 35, '全部订单', 'order_all',     '/order/all',     'view.order_all',     'mdi:format-list-bulleted',         1, 1, 0),
  (37, 35, '行程监控', 'order_monitor', '/order/monitor', 'view.order_monitor', 'mdi:map-marker-path',              2, 1, 0)
ON DUPLICATE KEY UPDATE
  parent_id = VALUES(parent_id),
  title = VALUES(title),
  route_name = VALUES(route_name),
  route_path = VALUES(route_path),
  component = VALUES(component),
  icon = VALUES(icon),
  sort_order = VALUES(sort_order),
  enabled = VALUES(enabled),
  hide_in_menu = VALUES(hide_in_menu);

ALTER TABLE sys_menu AUTO_INCREMENT = 38;
