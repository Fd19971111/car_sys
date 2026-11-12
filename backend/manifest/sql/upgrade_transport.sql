-- ============================================================
-- 2026-11-12 运力管理上线 / 库存查询下线
-- 1. 删除「库存查询」一级菜单及其全部二级菜单、部门授权
-- 2. 新增「运力管理」一级菜单，下挂「司机管理」「车辆管理」
-- 3. 新增 driver、vehicle 表（库存共用表保留，避免项目资料/付款/利润报错）
-- ============================================================
SET NAMES utf8mb4;

-- ---------- 1. 删除库存查询菜单 ----------
DELETE dm
FROM sys_department_menu dm
INNER JOIN sys_menu m ON m.id = dm.menu_id
WHERE m.id = 5 OR m.parent_id = 5;

DELETE FROM sys_menu WHERE parent_id = 5;
DELETE FROM sys_menu WHERE id = 5;

-- ---------- 2. 新增运力管理菜单 ----------
INSERT INTO sys_menu
  (id, parent_id, title, route_name, route_path, component, icon, sort_order, enabled, hide_in_menu)
VALUES
  (32, 0,  '运力管理', 'transport',         '/transport',         'layout.base',            'mdi:steering',            4, 1, 0),
  (33, 32, '司机管理', 'transport_driver',  '/transport/driver',  'view.transport_driver',  'mdi:account-tie-outline', 1, 1, 0),
  (34, 32, '车辆管理', 'transport_vehicle', '/transport/vehicle', 'view.transport_vehicle', 'mdi:car-outline',         2, 1, 0)
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

ALTER TABLE sys_menu AUTO_INCREMENT = 35;

-- ---------- 3. 司机表 ----------
CREATE TABLE IF NOT EXISTS `driver` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `company_id` int(10) unsigned NOT NULL DEFAULT '1' COMMENT '所属公司ID',
  `name` varchar(50) NOT NULL COMMENT '司机姓名',
  `phone` varchar(20) NOT NULL COMMENT '手机号',
  `license_number` varchar(50) NOT NULL DEFAULT '' COMMENT '驾驶证号',
  `license_type` varchar(20) NOT NULL DEFAULT '' COMMENT '准驾车型',
  `service_status` varchar(20) NOT NULL DEFAULT 'available' COMMENT '接单状态：available可接单 busy服务中 offline已离线 suspended暂停服务',
  `rating` decimal(3,2) NOT NULL DEFAULT '5.00' COMMENT '服务评分（0-5）',
  `completed_trips` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '完成订单数',
  `joined_at` date DEFAULT NULL COMMENT '加入日期',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_driver_company_status` (`company_id`, `service_status`) USING BTREE,
  KEY `idx_driver_company_phone` (`company_id`, `phone`) USING BTREE,
  KEY `idx_driver_deleted` (`company_id`, `deleted_at`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='司机档案';

-- ---------- 4. 车辆表 ----------
CREATE TABLE IF NOT EXISTS `vehicle` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `company_id` int(10) unsigned NOT NULL DEFAULT '1' COMMENT '所属公司ID',
  `plate_number` varchar(20) NOT NULL COMMENT '车牌号',
  `brand` varchar(100) NOT NULL DEFAULT '' COMMENT '品牌车型',
  `category` varchar(20) NOT NULL DEFAULT 'comfort' COMMENT '服务车型：economy经济型 comfort舒适型 business商务型 six-seat六座商务',
  `seat_count` tinyint(3) unsigned NOT NULL DEFAULT '5' COMMENT '可乘人数',
  `status` varchar(20) NOT NULL DEFAULT 'active' COMMENT '车辆状态：active运营中 maintenance保养中 inactive已停用',
  `driver_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '当前司机ID，0 表示未绑定',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_vehicle_company_status` (`company_id`, `status`) USING BTREE,
  KEY `idx_vehicle_company_plate` (`company_id`, `plate_number`) USING BTREE,
  KEY `idx_vehicle_driver` (`driver_id`) USING BTREE,
  KEY `idx_vehicle_deleted` (`company_id`, `deleted_at`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='车辆档案';
