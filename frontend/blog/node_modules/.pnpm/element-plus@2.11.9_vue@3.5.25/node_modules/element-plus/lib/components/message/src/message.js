'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var runtime = require('../../../utils/vue/props/runtime.js');
var icon = require('../../../utils/vue/icon.js');
var typescript = require('../../../utils/typescript.js');
var core = require('@vueuse/core');

const messageTypes = [
  "primary",
  "success",
  "info",
  "warning",
  "error"
];
const messagePlacement = [
  "top",
  "top-left",
  "top-right",
  "bottom",
  "bottom-left",
  "bottom-right"
];
const MESSAGE_DEFAULT_PLACEMENT = "top";
const messageDefaults = typescript.mutable({
  customClass: "",
  dangerouslyUseHTMLString: false,
  duration: 3e3,
  icon: void 0,
  id: "",
  message: "",
  onClose: void 0,
  showClose: false,
  type: "info",
  plain: false,
  offset: 16,
  placement: void 0,
  zIndex: 0,
  grouping: false,
  repeatNum: 1,
  appendTo: core.isClient ? document.body : void 0
});
const messageProps = runtime.buildProps({
  customClass: {
    type: String,
    default: messageDefaults.customClass
  },
  dangerouslyUseHTMLString: {
    type: Boolean,
    default: messageDefaults.dangerouslyUseHTMLString
  },
  duration: {
    type: Number,
    default: messageDefaults.duration
  },
  icon: {
    type: icon.iconPropType,
    default: messageDefaults.icon
  },
  id: {
    type: String,
    default: messageDefaults.id
  },
  message: {
    type: runtime.definePropType([
      String,
      Object,
      Function
    ]),
    default: messageDefaults.message
  },
  onClose: {
    type: runtime.definePropType(Function),
    default: messageDefaults.onClose
  },
  showClose: {
    type: Boolean,
    default: messageDefaults.showClose
  },
  type: {
    type: String,
    values: messageTypes,
    default: messageDefaults.type
  },
  plain: {
    type: Boolean,
    default: messageDefaults.plain
  },
  offset: {
    type: Number,
    default: messageDefaults.offset
  },
  placement: {
    type: String,
    values: messagePlacement,
    default: messageDefaults.placement
  },
  zIndex: {
    type: Number,
    default: messageDefaults.zIndex
  },
  grouping: {
    type: Boolean,
    default: messageDefaults.grouping
  },
  repeatNum: {
    type: Number,
    default: messageDefaults.repeatNum
  }
});
const messageEmits = {
  destroy: () => true
};

exports.MESSAGE_DEFAULT_PLACEMENT = MESSAGE_DEFAULT_PLACEMENT;
exports.messageDefaults = messageDefaults;
exports.messageEmits = messageEmits;
exports.messagePlacement = messagePlacement;
exports.messageProps = messageProps;
exports.messageTypes = messageTypes;
//# sourceMappingURL=message.js.map
