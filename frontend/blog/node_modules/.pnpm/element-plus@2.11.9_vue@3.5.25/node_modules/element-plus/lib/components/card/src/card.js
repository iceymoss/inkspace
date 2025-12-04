'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var runtime = require('../../../utils/vue/props/runtime.js');

const cardProps = runtime.buildProps({
  header: {
    type: String,
    default: ""
  },
  footer: {
    type: String,
    default: ""
  },
  bodyStyle: {
    type: runtime.definePropType([String, Object, Array]),
    default: ""
  },
  headerClass: String,
  bodyClass: String,
  footerClass: String,
  shadow: {
    type: String,
    values: ["always", "hover", "never"],
    default: void 0
  }
});
const cardContextKey = Symbol("cardContextKey");

exports.cardContextKey = cardContextKey;
exports.cardProps = cardProps;
//# sourceMappingURL=card.js.map
