'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var dialog$1 = require('./src/dialog2.js');
var useDialog = require('./src/use-dialog.js');
var dialog = require('./src/dialog.js');
var constants = require('./src/constants.js');
var install = require('../../utils/vue/install.js');

const ElDialog = install.withInstall(dialog$1["default"]);

exports.useDialog = useDialog.useDialog;
exports.dialogContextKey = dialog.dialogContextKey;
exports.dialogEmits = dialog.dialogEmits;
exports.dialogProps = dialog.dialogProps;
exports.DEFAULT_DIALOG_TRANSITION = constants.DEFAULT_DIALOG_TRANSITION;
exports.dialogInjectionKey = constants.dialogInjectionKey;
exports.ElDialog = ElDialog;
exports["default"] = ElDialog;
//# sourceMappingURL=index.js.map
