'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var iconsVue = require('@element-plus/icons-vue');
var button = require('../../button/src/button.js');
var runtime = require('../../../utils/vue/props/runtime.js');
var icon = require('../../../utils/vue/icon.js');
var content = require('../../tooltip/src/content.js');
var trigger = require('../../tooltip/src/trigger.js');

const popconfirmProps = runtime.buildProps({
  title: String,
  confirmButtonText: String,
  cancelButtonText: String,
  confirmButtonType: {
    type: String,
    values: button.buttonTypes,
    default: "primary"
  },
  cancelButtonType: {
    type: String,
    values: button.buttonTypes,
    default: "text"
  },
  icon: {
    type: icon.iconPropType,
    default: () => iconsVue.QuestionFilled
  },
  iconColor: {
    type: String,
    default: "#f90"
  },
  hideIcon: Boolean,
  hideAfter: {
    type: Number,
    default: 200
  },
  effect: {
    ...content.useTooltipContentProps.effect,
    default: "light"
  },
  teleported: content.useTooltipContentProps.teleported,
  persistent: content.useTooltipContentProps.persistent,
  width: {
    type: [String, Number],
    default: 150
  },
  virtualTriggering: trigger.useTooltipTriggerProps.virtualTriggering,
  virtualRef: trigger.useTooltipTriggerProps.virtualRef
});
const popconfirmEmits = {
  confirm: (e) => e instanceof MouseEvent,
  cancel: (e) => e instanceof MouseEvent
};

exports.popconfirmEmits = popconfirmEmits;
exports.popconfirmProps = popconfirmProps;
//# sourceMappingURL=popconfirm.js.map
