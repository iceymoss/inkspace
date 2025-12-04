import { defineComponent, inject, computed, openBlock, createElementBlock, normalizeClass, normalizeStyle, withModifiers, renderSlot } from 'vue';
import { DROPDOWN_INJECTION_KEY } from './tokens.mjs';
import { dropdownMenuProps } from './dropdown.mjs';
import { useDropdown } from './useDropdown.mjs';
import _export_sfc from '../../../_virtual/plugin-vue_export-helper.mjs';
import { ROVING_FOCUS_GROUP_INJECTION_KEY } from '../../roving-focus-group/src/tokens.mjs';
import { ROVING_FOCUS_COLLECTION_INJECTION_KEY as COLLECTION_INJECTION_KEY } from '../../roving-focus-group/src/roving-focus-group.mjs';
import { useNamespace } from '../../../hooks/use-namespace/index.mjs';
import { composeRefs } from '../../../utils/vue/refs.mjs';
import { composeEventHandlers, getEventCode } from '../../../utils/dom/event.mjs';
import { EVENT_CODE } from '../../../constants/aria.mjs';

const _sfc_main = defineComponent({
  name: "ElDropdownMenu",
  props: dropdownMenuProps,
  setup(props) {
    const ns = useNamespace("dropdown");
    const { _elDropdownSize } = useDropdown();
    const size = _elDropdownSize.value;
    const { contentRef, role, triggerId, isUsingKeyboard, handleClose } = inject(DROPDOWN_INJECTION_KEY, void 0);
    const {
      rovingFocusGroupRef,
      rovingFocusGroupRootStyle,
      onBlur,
      onFocus,
      onKeydown,
      onMousedown
    } = inject(ROVING_FOCUS_GROUP_INJECTION_KEY, void 0);
    const { collectionRef: rovingFocusGroupCollectionRef } = inject(COLLECTION_INJECTION_KEY, void 0);
    const dropdownKls = computed(() => {
      return [ns.b("menu"), ns.bm("menu", size == null ? void 0 : size.value)];
    });
    const dropdownListWrapperRef = composeRefs(contentRef, rovingFocusGroupRef, rovingFocusGroupCollectionRef);
    const handleKeydown = composeEventHandlers((e) => {
      var _a;
      (_a = props.onKeydown) == null ? void 0 : _a.call(props, e);
    }, (e) => {
      const { currentTarget, target } = e;
      const code = getEventCode(e);
      currentTarget.contains(target);
      if (EVENT_CODE.tab === code) {
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
  return openBlock(), createElementBlock("ul", {
    ref: _ctx.dropdownListWrapperRef,
    class: normalizeClass(_ctx.dropdownKls),
    style: normalizeStyle(_ctx.rovingFocusGroupRootStyle),
    tabindex: -1,
    role: _ctx.role,
    "aria-labelledby": _ctx.triggerId,
    onFocusin: _ctx.handleFocus,
    onFocusout: _ctx.onBlur,
    onKeydown: withModifiers(_ctx.handleKeydown, ["self"]),
    onMousedown: withModifiers(_ctx.onMousedown, ["self"])
  }, [
    renderSlot(_ctx.$slots, "default")
  ], 46, ["role", "aria-labelledby", "onFocusin", "onFocusout", "onKeydown", "onMousedown"]);
}
var DropdownMenu = /* @__PURE__ */ _export_sfc(_sfc_main, [["render", _sfc_render], ["__file", "dropdown-menu.vue"]]);

export { DropdownMenu as default };
//# sourceMappingURL=dropdown-menu.mjs.map
