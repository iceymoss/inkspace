'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var button = require('./button.js');
var runtime = require('../../../utils/vue/props/runtime.js');

const buttonGroupProps = {
  size: button.buttonProps.size,
  type: button.buttonProps.type,
  direction: {
    type: runtime.definePropType(String),
    values: ["horizontal", "vertical"],
    default: "horizontal"
  }
};

exports.buttonGroupProps = buttonGroupProps;
//# sourceMappingURL=button-group.js.map
