'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var input = require('../../input/src/input.js');
var content = require('../../tooltip/src/content.js');
var runtime = require('../../../utils/vue/props/runtime.js');
var shared = require('@vue/shared');
var event = require('../../../constants/event.js');
var types = require('../../../utils/types.js');

const autocompleteProps = runtime.buildProps({
  ...input.inputProps,
  valueKey: {
    type: String,
    default: "value"
  },
  modelValue: {
    type: [String, Number],
    default: ""
  },
  debounce: {
    type: Number,
    default: 300
  },
  placement: {
    type: runtime.definePropType(String),
    values: [
      "top",
      "top-start",
      "top-end",
      "bottom",
      "bottom-start",
      "bottom-end"
    ],
    default: "bottom-start"
  },
  fetchSuggestions: {
    type: runtime.definePropType([Function, Array]),
    default: shared.NOOP
  },
  popperClass: content.useTooltipContentProps.popperClass,
  popperStyle: content.useTooltipContentProps.popperStyle,
  triggerOnFocus: {
    type: Boolean,
    default: true
  },
  selectWhenUnmatched: Boolean,
  hideLoading: Boolean,
  teleported: content.useTooltipContentProps.teleported,
  appendTo: content.useTooltipContentProps.appendTo,
  highlightFirstItem: Boolean,
  fitInputWidth: Boolean,
  loopNavigation: {
    type: Boolean,
    default: true
  }
});
const autocompleteEmits = {
  [event.UPDATE_MODEL_EVENT]: (value) => shared.isString(value) || types.isNumber(value),
  [event.INPUT_EVENT]: (value) => shared.isString(value) || types.isNumber(value),
  [event.CHANGE_EVENT]: (value) => shared.isString(value) || types.isNumber(value),
  focus: (evt) => evt instanceof FocusEvent,
  blur: (evt) => evt instanceof FocusEvent,
  clear: () => true,
  select: (item) => shared.isObject(item)
};

exports.autocompleteEmits = autocompleteEmits;
exports.autocompleteProps = autocompleteProps;
//# sourceMappingURL=autocomplete.js.map
