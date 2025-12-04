'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var lodashUnified = require('lodash-unified');
var runtime = require('../../../utils/vue/props/runtime.js');
var index = require('../../../hooks/use-size/index.js');
var content = require('../../tooltip/src/content.js');
var index$1 = require('../../../hooks/use-empty-values/index.js');
var index$2 = require('../../../hooks/use-aria/index.js');
var event = require('../../../constants/event.js');
var shared = require('@vue/shared');

const colorPickerProps = runtime.buildProps({
  persistent: {
    type: Boolean,
    default: true
  },
  modelValue: {
    type: runtime.definePropType(String),
    default: void 0
  },
  id: String,
  showAlpha: Boolean,
  colorFormat: String,
  disabled: Boolean,
  size: index.useSizeProp,
  popperClass: content.useTooltipContentProps.popperClass,
  popperStyle: content.useTooltipContentProps.popperStyle,
  tabindex: {
    type: [String, Number],
    default: 0
  },
  teleported: content.useTooltipContentProps.teleported,
  appendTo: content.useTooltipContentProps.appendTo,
  predefine: {
    type: runtime.definePropType(Array)
  },
  validateEvent: {
    type: Boolean,
    default: true
  },
  ...index$1.useEmptyValuesProps,
  ...index$2.useAriaProps(["ariaLabel"])
});
const colorPickerEmits = {
  [event.UPDATE_MODEL_EVENT]: (val) => shared.isString(val) || lodashUnified.isNil(val),
  [event.CHANGE_EVENT]: (val) => shared.isString(val) || lodashUnified.isNil(val),
  activeChange: (val) => shared.isString(val) || lodashUnified.isNil(val),
  focus: (evt) => evt instanceof FocusEvent,
  blur: (evt) => evt instanceof FocusEvent
};

exports.colorPickerEmits = colorPickerEmits;
exports.colorPickerProps = colorPickerProps;
//# sourceMappingURL=color-picker.js.map
