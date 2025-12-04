'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var vue = require('vue');
var runtime = require('../../../utils/vue/props/runtime.js');
var shared = require('@vue/shared');
var event = require('../../../constants/event.js');

const CommonProps = runtime.buildProps({
  modelValue: {
    type: runtime.definePropType([Number, String, Array, Object])
  },
  options: {
    type: runtime.definePropType(Array),
    default: () => []
  },
  props: {
    type: runtime.definePropType(Object),
    default: () => ({})
  }
});
const DefaultProps = {
  expandTrigger: "click",
  multiple: false,
  checkStrictly: false,
  emitPath: true,
  lazy: false,
  lazyLoad: shared.NOOP,
  value: "value",
  label: "label",
  children: "children",
  leaf: "leaf",
  disabled: "disabled",
  hoverThreshold: 500,
  checkOnClickNode: false,
  checkOnClickLeaf: true,
  showPrefix: true
};
const cascaderPanelProps = runtime.buildProps({
  ...CommonProps,
  border: {
    type: Boolean,
    default: true
  },
  renderLabel: {
    type: Function
  }
});
const emitChangeFn = (value) => true;
const cascaderPanelEmits = {
  [event.UPDATE_MODEL_EVENT]: emitChangeFn,
  [event.CHANGE_EVENT]: emitChangeFn,
  close: () => true,
  "expand-change": (value) => value
};
const useCascaderConfig = (props) => {
  return vue.computed(() => ({
    ...DefaultProps,
    ...props.props
  }));
};

exports.CommonProps = CommonProps;
exports.DefaultProps = DefaultProps;
exports.cascaderPanelEmits = cascaderPanelEmits;
exports.cascaderPanelProps = cascaderPanelProps;
exports.useCascaderConfig = useCascaderConfig;
//# sourceMappingURL=config.js.map
