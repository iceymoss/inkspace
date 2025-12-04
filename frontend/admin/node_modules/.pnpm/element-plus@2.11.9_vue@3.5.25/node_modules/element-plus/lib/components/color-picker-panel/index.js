'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var colorPickerPanel$1 = require('./src/color-picker-panel2.js');
var colorPickerPanel = require('./src/color-picker-panel.js');
var install = require('../../utils/vue/install.js');

const ElColorPickerPanel = install.withInstall(colorPickerPanel$1["default"]);

exports.ROOT_COMMON_COLOR_INJECTION_KEY = colorPickerPanel.ROOT_COMMON_COLOR_INJECTION_KEY;
exports.colorPickerPanelContextKey = colorPickerPanel.colorPickerPanelContextKey;
exports.colorPickerPanelEmits = colorPickerPanel.colorPickerPanelEmits;
exports.colorPickerPanelProps = colorPickerPanel.colorPickerPanelProps;
exports.ElColorPickerPanel = ElColorPickerPanel;
exports["default"] = ElColorPickerPanel;
//# sourceMappingURL=index.js.map
