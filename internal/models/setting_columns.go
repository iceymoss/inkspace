package models

// TableName explicitly fixes the persisted settings table name.
func (Setting) TableName() string { return "settings" }

// SettingColumns centralizes persisted columns used by the public identity query.
var SettingColumns = struct{ Key, IsPublic string }{Key: "key", IsPublic: "is_public"}
