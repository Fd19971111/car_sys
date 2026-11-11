-- Run once on databases created before purchase tax support.
ALTER TABLE material_catalog
  ADD COLUMN tax_mode varchar(20) NOT NULL DEFAULT 'NO_INVOICE' COMMENT '默认税务方式' AFTER unit_price,
  ADD COLUMN tax_rate decimal(5,2) NOT NULL DEFAULT 0.00 COMMENT '默认税率（百分比）' AFTER tax_mode;

ALTER TABLE project_material
  ADD COLUMN tax_mode varchar(20) NOT NULL DEFAULT 'NO_INVOICE' COMMENT '税务方式' AFTER unit_price,
  ADD COLUMN tax_rate decimal(5,2) NOT NULL DEFAULT 0.00 COMMENT '税率（百分比）' AFTER tax_mode;

ALTER TABLE project_purchase_order
  ADD COLUMN tax_mode varchar(20) NOT NULL DEFAULT 'NO_INVOICE' COMMENT '税务方式快照' AFTER unit_price,
  ADD COLUMN tax_rate decimal(5,2) NOT NULL DEFAULT 0.00 COMMENT '税率快照（百分比）' AFTER tax_mode,
  ADD COLUMN net_amount decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '未税金额' AFTER tax_rate,
  ADD COLUMN tax_amount decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '税额' AFTER net_amount;

UPDATE project_purchase_order
SET net_amount=amount, tax_amount=0
WHERE tax_mode='NO_INVOICE';
