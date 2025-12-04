'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var lodashUnified = require('lodash-unified');
var runtime = require('../../../utils/vue/props/runtime.js');
var event = require('../../../constants/event.js');
var shared = require('@vue/shared');

const colorPickerPanelProps = runtime.buildProps({
  modelValue: {
    type: runtime.definePropType(String),
    default: void 0
  },
  border: {
    type: Boolean,
    default: true
  },
  showAlpha: Boolean,
  colorFormat: String,
  disabled: Boolean,
  predefine: {
    type: runtime.definePropType(Array)
  },
  validateEvent: {
    type: Boolean,
    default: true
  }
});
const colorPickerPanelEmits = {
  [event.UPDATE_MODEL_EVENT]: (val) => shared.isString(val) || lodashUnified.isNil(val)
};
const ROOT_COMMON_COLOR_INJECTION_KEY = Symbol("colorCommonPickerKey");
const colorPickerPanelContextKey = Symbol("colorPickerPanelContextKey");

exports.ROOT_COMMON_COLOR_INJECTION_KEY = ROOT_COMMON_COLOR_INJECTION_KEY;
exports.colorPickerPanelContextKey = colorPickerPanelContextKey;
exports.colorPickerPanelEmits = colorPickerPanelEmits;
exports.colorPickerPanelProps = colorPickerPanelProps;
//# sourceMappingURL=color-picker-panel.js.map
