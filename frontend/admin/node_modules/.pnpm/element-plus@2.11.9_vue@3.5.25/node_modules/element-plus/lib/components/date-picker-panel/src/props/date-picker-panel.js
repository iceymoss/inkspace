'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var shared = require('../../../time-picker/src/props/shared.js');
var runtime = require('../../../../utils/vue/props/runtime.js');

const datePickerPanelProps = runtime.buildProps({
  valueFormat: String,
  dateFormat: String,
  timeFormat: String,
  disabled: Boolean,
  modelValue: {
    type: runtime.definePropType([Date, Array, String, Number]),
    default: ""
  },
  defaultValue: {
    type: runtime.definePropType([Date, Array])
  },
  defaultTime: {
    type: runtime.definePropType([Date, Array])
  },
  isRange: Boolean,
  ...shared.disabledTimeListsProps,
  disabledDate: {
    type: Function
  },
  cellClassName: {
    type: Function
  },
  shortcuts: {
    type: Array,
    default: () => []
  },
  arrowControl: Boolean,
  unlinkPanels: Boolean,
  showNow: {
    type: Boolean,
    default: true
  },
  showConfirm: Boolean,
  showFooter: Boolean,
  showWeekNumber: Boolean,
  type: {
    type: runtime.definePropType(String),
    default: "date"
  },
  clearable: {
    type: Boolean,
    default: true
  },
  border: {
    type: Boolean,
    default: true
  }
});

exports.datePickerPanelProps = datePickerPanelProps;
//# sourceMappingURL=date-picker-panel.js.map
