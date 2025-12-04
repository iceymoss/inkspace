import { defineComponent, getCurrentInstance, inject, computed, resolveComponent, openBlock, createBlock, withCtx, createVNode, mergeProps, renderSlot } from 'vue';
import ElRovingFocusItem from '../../roving-focus-group/src/roving-focus-item.mjs';
import ElDropdownItemImpl from './dropdown-item-impl.mjs';
import { useDropdown } from './useDropdown.mjs';
import { dropdownItemProps } from './dropdown.mjs';
import { DROPDOWN_INJECTION_KEY } from './tokens.mjs';
import _export_sfc from '../../../_virtual/plugin-vue_export-helper.mjs';
import { composeEventHandlers, whenMouse } from '../../../utils/dom/event.mjs';

const _sfc_main = defineComponent({
  name: "ElDropdownItem",
  components: {
    ElRovingFocusItem,
    ElDropdownItemImpl
  },
  inheritAttrs: false,
  props: dropdownItemProps,
  emits: ["pointermove", "pointerleave", "click"],
  setup(props, { emit, attrs }) {
    const { elDropdown } = useDropdown();
    const _instance = getCurrentInstance();
    const { onItemEnter, onItemLeave } = inject(DROPDOWN_INJECTION_KEY, void 0);
    const handlePointerMove = composeEventHandlers((e) => {
      emit("pointermove", e);
      return e.defaultPrevented;
    }, whenMouse((e) => {
      if (props.disabled) {
        onItemLeave(e);
        return;
      }
      const target = e.currentTarget;
      if (target === document.activeElement || target.contains(document.activeElement)) {
        return;
      }
      onItemEnter(e);
      if (!e.defaultPrevented) {
        target == null ? void 0 : target.focus({
          preventScroll: true
        });
      }
    }));
    const handlePointerLeave = composeEventHandlers((e) => {
      emit("pointerleave", e);
      return e.defaultPrevented;
    }, whenMouse(onItemLeave));
    const handleClick = composeEventHandlers((e) => {
      if (props.disabled) {
        return;
      }
      emit("click", e);
      return e.type !== "keydown" && e.defaultPrevented;
    }, (e) => {
      var _a, _b, _c;
      if (props.disabled) {
        e.stopImmediatePropagation();
        return;
      }
      if ((_a = elDropdown == null ? void 0 : elDropdown.hideOnClick) == null ? void 0 : _a.value) {
        (_b = elDropdown.handleClick) == null ? void 0 : _b.call(elDropdown);
      }
      (_c = elDropdown.commandHandler) == null ? void 0 : _c.call(elDropdown, props.command, _instance, e);
    });
    const propsAndAttrs = computed(() => ({ ...props, ...attrs }));
    return {
      handleClick,
      handlePointerMove,
      handlePointerLeave,
      propsAndAttrs
    };
  }
});
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  const _component_el_dropdown_item_impl = resolveComponent("el-dropdown-item-impl");
  const _component_el_roving_focus_item = resolveComponent("el-roving-focus-item");
  return openBlock(), createBlock(_component_el_roving_focus_item, {
    focusable: !_ctx.disabled
  }, {
    default: withCtx(() => [
      createVNode(_component_el_dropdown_item_impl, mergeProps(_ctx.propsAndAttrs, {
        onPointerleave: _ctx.handlePointerLeave,
        onPointermove: _ctx.handlePointerMove,
        onClickimpl: _ctx.handleClick
      }), {
        default: withCtx(() => [
          renderSlot(_ctx.$slots, "default")
        ]),
        _: 3
      }, 16, ["onPointerleave", "onPointermove", "onClickimpl"])
    ]),
    _: 3
  }, 8, ["focusable"]);
}
var DropdownItem = /* @__PURE__ */ _export_sfc(_sfc_main, [["render", _sfc_render], ["__file", "dropdown-item.vue"]]);

export { DropdownItem as default };
//# sourceMappingURL=dropdown-item.mjs.map
