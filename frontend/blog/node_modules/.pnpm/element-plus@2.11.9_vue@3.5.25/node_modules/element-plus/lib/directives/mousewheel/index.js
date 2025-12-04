'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var normalizeWheel = require('normalize-wheel-es');

function _interopDefaultLegacy (e) { return e && typeof e === 'object' && 'default' in e ? e : { 'default': e }; }

var normalizeWheel__default = /*#__PURE__*/_interopDefaultLegacy(normalizeWheel);

const SCOPE = "_Mousewheel";
const mousewheel = function(element, callback) {
  if (element && element.addEventListener) {
    removeWheelHandler(element);
    const fn = function(event) {
      const normalized = normalizeWheel__default["default"](event);
      callback && Reflect.apply(callback, this, [event, normalized]);
    };
    element[SCOPE] = { wheelHandler: fn };
    element.addEventListener("wheel", fn, { passive: true });
  }
};
const removeWheelHandler = (element) => {
  var _a;
  if ((_a = element[SCOPE]) == null ? void 0 : _a.wheelHandler) {
    element.removeEventListener("wheel", element[SCOPE].wheelHandler);
    element[SCOPE] = null;
  }
};
const Mousewheel = {
  beforeMount(el, binding) {
    mousewheel(el, binding.value);
  },
  unmounted(el) {
    removeWheelHandler(el);
  },
  updated(el, binding) {
    if (binding.value !== binding.oldValue) {
      mousewheel(el, binding.value);
    }
  }
};

exports.SCOPE = SCOPE;
exports["default"] = Mousewheel;
//# sourceMappingURL=index.js.map
