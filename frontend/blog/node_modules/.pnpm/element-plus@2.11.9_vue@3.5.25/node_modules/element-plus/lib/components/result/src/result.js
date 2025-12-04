'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var iconsVue = require('@element-plus/icons-vue');
var runtime = require('../../../utils/vue/props/runtime.js');

const IconMap = {
  primary: "icon-primary",
  success: "icon-success",
  warning: "icon-warning",
  error: "icon-error",
  info: "icon-info"
};
const IconComponentMap = {
  [IconMap.primary]: iconsVue.InfoFilled,
  [IconMap.success]: iconsVue.CircleCheckFilled,
  [IconMap.warning]: iconsVue.WarningFilled,
  [IconMap.error]: iconsVue.CircleCloseFilled,
  [IconMap.info]: iconsVue.InfoFilled
};
const resultProps = runtime.buildProps({
  title: {
    type: String,
    default: ""
  },
  subTitle: {
    type: String,
    default: ""
  },
  icon: {
    type: String,
    values: ["primary", "success", "warning", "info", "error"],
    default: "info"
  }
});

exports.IconComponentMap = IconComponentMap;
exports.IconMap = IconMap;
exports.resultProps = resultProps;
//# sourceMappingURL=result.js.map
