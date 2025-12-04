'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var vue = require('vue');
var constants = require('./constants.js');
var columnAlignment = require('../../../constants/column-alignment.js');
var runtime = require('../../../utils/vue/props/runtime.js');

const descriptionItemProps = runtime.buildProps({
  label: {
    type: String,
    default: ""
  },
  span: {
    type: Number,
    default: 1
  },
  rowspan: {
    type: Number,
    default: 1
  },
  width: {
    type: [String, Number],
    default: ""
  },
  minWidth: {
    type: [String, Number],
    default: ""
  },
  labelWidth: {
    type: [String, Number]
  },
  align: {
    type: String,
    values: columnAlignment.columnAlignment,
    default: "left"
  },
  labelAlign: {
    type: String,
    values: columnAlignment.columnAlignment
  },
  className: {
    type: String,
    default: ""
  },
  labelClassName: {
    type: String,
    default: ""
  }
});
const DescriptionItem = vue.defineComponent({
  name: constants.COMPONENT_NAME,
  props: descriptionItemProps
});

exports["default"] = DescriptionItem;
exports.descriptionItemProps = descriptionItemProps;
//# sourceMappingURL=description-item.js.map
