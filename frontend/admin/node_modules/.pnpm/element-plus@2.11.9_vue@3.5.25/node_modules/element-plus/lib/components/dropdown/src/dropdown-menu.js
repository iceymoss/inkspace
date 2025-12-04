'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var vue = require('vue');
var tokens = require('./tokens.js');
var dropdown = require('./dropdown.js');
var useDropdown = require('./useDropdown.js');
var pluginVue_exportHelper = require('../../../_virtual/plugin-vue_export-helper.js');
var tokens$1 = require('../../roving-focus-group/src/tokens.js');
var rovingFocusGroup = require('../../roving-focus-group/src/roving-focus-group.js');
var index = require('../../../hooks/use-namespace/index.js');
var refs = require('../../../utils/vue/refs.js');
var event = require('../../../utils/dom/event.js');
var aria = require('../../../constants/aria.js');

const _sfc_main = vue.defineComponent({
  name: "ElDropdownMenu",
  props: dropdown.dropdownMenuProps,
  setup(props) {
    const ns = index.useNamespace("dropdown");
    const { _elDropdownSize } = useDropdown.useDropdown();
    const size = _elDropdownSize.value;
    const { contentRef, role, triggerId, isUsingKeyboard, handleClose } = vue.inject(tokens.DROPDOWN_INJECTION_KEY, void 0);
    const {
      rovingFocusGroupRef,
      rovingFocusGroupRootStyle,
      onBlur,
      onFocus,
      onKeydown,
      onMousedown
    } = vue.inject(tokens$1.ROVING_FOCUS_GROUP_INJECTION_KEY, void 0);
    const { collectionRef: rovingFocusGroupCollectionRef } = vue.inject(rovingFocusGroup.ROVING_FOCUS_COLLECTION_INJECTION_KEY, void 0);
    const dropdownKls = vue.computed(() => {
      return [ns.b("menu"), ns.bm("menu", size == null ? void 0 : size.value)];
    });
    const dropdownListWrapperRef = refs.composeRefs(contentRef, rovingFocusGroupRef, rovingFocusGroupCollectionRef);
    const handleKeydown = event.composeEventHandlers((e) => {
      var _a;
      (_a = props.onKeydown) == null ? void 0 : _a.call(props, e);
    }, (e) => {
      const { currentTarget, target } = e;
      const code = event.getEventCode(e);
      currentTarget.contains(target);
      if (aria.EVENT_CODE.tab === code) {
        return handleClose();
      }
      onKeydown(e);
    });
    function handleFocus(e) {
      isUsingKeyboard.value && onFocus(e);
    }
    return {
      size,
      rovingFocusGroupRootStyle,
      dropdownKls,
      role,
      triggerId,
      dropdownListWrapperRef,
      handleKeydown,
      onBlur,
      handleFocus,
      onMousedown
    };
  }
});
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return vue.openBlock(), vue.createElementBlock("ul", {
    ref: _ctx.dropdownListWrapperRef,
    class: vue.normalizeClass(_ctx.dropdownKls),
    style: vue.normalizeStyle(_ctx.rovingFocusGroupRootStyle),
    tabindex: -1,
    role: _ctx.role,
    "aria-labelledby": _ctx.triggerId,
    onFocusin: _ctx.handleFocus,
    onFocusout: _ctx.onBlur,
    onKeydown: vue.withModifiers(_ctx.handleKeydown, ["self"]),
    onMousedown: vue.withModifiers(_ctx.onMousedown, ["self"])
  }, [
    vue.renderSlot(_ctx.$slots, "default")
  ], 46, ["role", "aria-labelledby", "onFocusin", "onFocusout", "onKeydown", "onMousedown"]);
}
var DropdownMenu = /* @__PURE__ */ pluginVue_exportHelper["default"](_sfc_main, [["render", _sfc_render], ["__file", "dropdown-menu.vue"]]);

exports["default"] = DropdownMenu;
//# sourceMappingURL=dropdown-menu.js.map
