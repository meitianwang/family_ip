-- 087_add_setting_group_label.sql
-- 给 settings 表增加 setting_group 和 label 字段，用于分类和可读描述

ALTER TABLE settings ADD COLUMN IF NOT EXISTS setting_group VARCHAR(50) NOT NULL DEFAULT 'general';
ALTER TABLE settings ADD COLUMN IF NOT EXISTS label VARCHAR(200);

CREATE INDEX IF NOT EXISTS idx_settings_setting_group ON settings(setting_group);

-- 将已有的支付配置项归类到 payment 分组
UPDATE settings SET setting_group = 'payment', label = '订单超时时间（分钟）'       WHERE key = 'pay_order_timeout_minutes';
UPDATE settings SET setting_group = 'payment', label = '最小充值金额'               WHERE key = 'pay_min_recharge_amount';
UPDATE settings SET setting_group = 'payment', label = '最大充值金额'               WHERE key = 'pay_max_recharge_amount';
UPDATE settings SET setting_group = 'payment', label = '每日最大充值金额（每用户）'  WHERE key = 'pay_max_daily_recharge_amount';
UPDATE settings SET setting_group = 'payment', label = '商品名称'                   WHERE key = 'pay_product_name';
UPDATE settings SET setting_group = 'payment', label = '启用的支付渠道'             WHERE key = 'pay_providers';
UPDATE settings SET setting_group = 'payment', label = '支付帮助图片 URL'           WHERE key = 'pay_help_image_url';
UPDATE settings SET setting_group = 'payment', label = '支付帮助文本'               WHERE key = 'pay_help_text';
UPDATE settings SET setting_group = 'payment', label = '支付宝每日限额'             WHERE key = 'pay_max_daily_amount_alipay';
UPDATE settings SET setting_group = 'payment', label = '微信支付每日限额'           WHERE key = 'pay_max_daily_amount_wxpay';
UPDATE settings SET setting_group = 'payment', label = 'Stripe 每日限额'            WHERE key = 'pay_max_daily_amount_stripe';
