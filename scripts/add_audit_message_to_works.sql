-- 添加 audit_message 字段到 works 表
ALTER TABLE `works`
ADD COLUMN IF NOT EXISTS `audit_message` TEXT NULL COMMENT '审核消息（审核通过或拒绝的原因）' AFTER `is_recommend`;

