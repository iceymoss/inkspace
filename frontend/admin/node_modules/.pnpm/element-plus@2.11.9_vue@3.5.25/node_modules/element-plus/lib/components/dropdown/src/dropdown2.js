'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var vue = require('vue');
var index = require('../../button/index.js');
var index$2 = require('../../tooltip/index.js');
var index$1 = require('../../scrollbar/index.js');
var index$3 = require('../../icon/index.js');
var rovingFocusGroup = require('../../roving-focus-group/src/roving-focus-group2.js');
var iconsVue = require('@element-plus/icons-vue');
var dropdown = require('./dropdown.js');
var tokens = require('./tokens.js');
var pluginVue_exportHelper = require('../../../_virtual/plugin-vue_export-helper.js');
var onlyChild = require('../../slot/src/only-child.js');
var index$4 = require('../../../hooks/use-namespace/index.js');
var index$5 = require('../../../hooks/use-locale/index.js');
var style = require('../../../utils/dom/style.js');
var lodashUnified = require('lodash-unified');
var index$6 = require('../../../hooks/use-id/index.js');
var useFormCommonProps = require('../../form/src/hooks/use-form-common-props.js');

const { ButtonGroup: ElButtonGroup } = index.ElButton;
const _sfc_main = vue.defineComponent({
  name: "ElDropdown",
  components: {
    ElButton: index.ElButton,
    ElButtonGroup,
    ElScrollbar: index$1.ElScrollbar,
    ElTooltip: index$2.ElTooltip,
    ElRovingFocusGroup: rovingFocusGroup["default"],
    ElOnlyChild: onlyChild.OnlyChild,
    ElIcon: index$3.ElIcon,
    ArrowDown: iconsVue.ArrowDown
  },
  props: dropdown.dropdownProps,
  emits: ["visible-change", "click", "command"],
  setup(props, { emit }) {
    const _instance = vue.getCurrentInstance();
    const ns = index$4.useNamespace("dropdown");
    const { t } = index$5.useLocale();
    const triggeringElementRef = vue.ref();
    const referenceElementRef = vue.ref();
    const popperRef = vue.ref();
    const contentRef = vue.ref();
    const scrollbar = vue.ref(null);
    const currentTabId = vue.ref(null);
    const isUsingKeyboard = vue.ref(false);
    const wrapStyle = vue.computed(() => ({
      maxHeight: style.addUnit(props.maxHeight)
    }));
    const dropdownTriggerKls = vue.computed(() => [ns.m(dropdownSize.value)]);
    const trigger = vue.computed(() => lodashUnified.castArray(props.trigger));
    const defaultTriggerId = index$6.useId().value;
    const triggerId = vue.computed(() => props.id || defaultTriggerId);
    function handleClick() {
      var _a;
      (_a = popperRef.value) == null ? void 0 : _a.onClose(void 0, 0);
    }
    function handleClose() {
      var _a;
      (_a = popperRef.value) == null ? void 0 : _a.onClose();
    }
    function handleOpen() {
      var _a;
      (_a = popperRef.value) == null ? void 0 : _a.onOpen();
    }
    const dropdownSize = useFormCommonProps.useFormSize();
    function commandHandler(...args) {
      emit("command", ...args);
    }
    function onItemEnter() {
    }
    function onItemLeave() {
      const contentEl = vue.unref(contentRef);
      trigger.value.includes("hover") && (contentEl == null ? void 0 : contentEl.focus({
        preventScroll: true
      }));
      currentTabId.value = null;
    }
    function handleCurrentTabIdChange(id) {
      currentTabId.value = id;
    }
    function handleBeforeShowTooltip() {
      emit("visible-change", true);
    }
    function handleShowTooltip(event) {
      var _a;
      isUsingKeyboard.value = (event == null ? void 0 : event.type) === "keydown";
      (_a = contentRef.value) == null ? void 0 : _a.focus();
    }
    function handleBeforeHideTooltip() {
      emit("visible-change", false);
    }
    vue.provide(tokens.DROPDOWN_INJECTION_KEY, {
      contentRef,
      role: vue.computed(() => props.role),
      triggerId,
      isUsingKeyboard,
      onItemEnter,
      onItemLeave,
      handleClose
    });
    vue.provide(tokens.DROPDOWN_INSTANCE_INJECTION_KEY, {
      instance: _instance,
      dropdownSize,
      handleClick,
      commandHandler,
      trigger: vue.toRef(props, "trigger"),
      hideOnClick: vue.toRef(props, "hideOnClick")
    });
    const handlerMainButtonClick = (event) => {
      emit("click", event);
    };
    return {
      t,
      ns,
      scrollbar,
      wrapStyle,
      dropdownTriggerKls,
      dropdownSize,
      triggerId,
      currentTabId,
      handleCurrentTabIdChange,
      handlerMainButtonClick,
      handleClose,
      handleOpen,
      handleBeforeShowTooltip,
      handleShowTooltip,
      handleBeforeHideTooltip,
      popperRef,
      contentRef,
      triggeringElementRef,
      referenceElementRef
    };
  }
});
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  var _a;
  const _component_el_roving_focus_group = vue.resolveComponent("el-roving-focus-group");
  const _component_el_scrollbar = vue.resolveComponent("el-scrollbar");
  const _component_el_only_child = vue.resolveComponent("el-only-child");
  const _component_el_tooltip = vue.resolveComponent("el-tooltip");
  const _component_el_button = vue.resolveComponent("el-button");
  const _component_arrow_down = vue.resolveComponent("arrow-down");
  const _component_el_icon = vue.resolveComponent("el-icon");
  const _component_el_button_group = vue.resolveComponent("el-button-group");
  return vue.openBlock(), vue.createElementBlock("div", {
    class: vue.normalizeClass([_ctx.ns.b(), _ctx.ns.is("disabled", _ctx.disabled)])
  }, [
    vue.createVNode(_component_el_tooltip, {
      ref: "popperRef",
      role: _ctx.role,
      effect: _ctx.effect,
      "fallback-placements": ["bottom", "top"],
      "popper-options": _ctx.popperOptions,
      "gpu-acceleration": false,
      placement: _ctx.placement,
      "popper-class": [_ctx.ns.e("popper"), _ctx.popperClass],
      "popper-style": _ctx.popperStyle,
      trigger: _ctx.trigger,
      "trigger-keys": _ctx.triggerKeys,
      "trigger-target-el": _ctx.contentRef,
      "show-arrow": _ctx.showArrow,
      "show-after": _ctx.trigger === "hover" ? _ctx.showTimeout : 0,
      "hide-after": _ctx.trigger === "hover" ? _ctx.hideTimeout : 0,
      "virtual-ref": (_a = _ctx.virtualRef) != null ? _a : _ctx.triggeringElementRef,
      "virtual-triggering": _ctx.virtualTriggering || _ctx.splitButton,
      disabled: _ctx.disabled,
      transition: `${_ctx.ns.namespace.value}-zoom-in-top`,
      teleported: _ctx.teleported,
      pure: "",
      "focus-on-target": "",
      persistent: _ctx.persistent,
      onBeforeShow: _ctx.handleBeforeShowTooltip,
      onShow: _ctx.handleShowTooltip,
      onBeforeHide: _ctx.handleBeforeHideTooltip
    }, vue.createSlots({
      content: vue.withCtx(() => [
        vue.createVNode(_component_el_scrollbar, {
          ref: "scrollbar",
          "wrap-style": _ctx.wrapStyle,
          tag: "div",
          "view-class": _ctx.ns.e("list")
        }, {
          default: vue.withCtx(() => [
            vue.createVNode(_component_el_roving_focus_group, {
              loop: _ctx.loop,
              "current-tab-id": _ctx.currentTabId,
              orientation: "horizontal",
              onCurrentTabIdChange: _ctx.handleCurrentTabIdChange
            }, {
              default: vue.withCtx(() => [
                vue.renderSlot(_ctx.$slots, "dropdown")
              ]),
              _: 3
            }, 8, ["loop", "current-tab-id", "onCurrentTabIdChange"])
          ]),
          _: 3
        }, 8, ["wrap-style", "view-class"])
      ]),
      _: 2
    }, [
      !_ctx.splitButton ? {
        name: "default",
        fn: vue.withCtx(() => [
          vue.createVNode(_component_el_only_child, {
            id: _ctx.triggerId,
            ref: "triggeringElementRef",
            role: "button",
            tabindex: _ctx.tabindex
          }, {
            default: vue.withCtx(() => [
              vue.renderSlot(_ctx.$slots, "default")
            ]),
            _: 3
          }, 8, ["id", "tabindex"])
        ])
      } : void 0
    ]), 1032, ["role", "effect", "popper-options", "placement", "popper-class", "popper-style", "trigger", "trigger-keys", "trigger-target-el", "show-arrow", "show-after", "hide-after", "virtual-ref", "virtual-triggering", "disabled", "transition", "teleported", "persistent", "onBeforeShow", "onShow", "onBeforeHide"]),
    _ctx.splitButton ? (vue.openBlock(), vue.createBlock(_component_el_button_group, { key: 0 }, {
      default: vue.withCtx(() => [
        vue.createVNode(_component_el_button, vue.mergeProps({ ref: "referenceElementRef" }, _ctx.buttonProps, {
          size: _ctx.dropdownSize,
          type: _ctx.type,
          disabled: _ctx.disabled,
          tabindex: _ctx.tabindex,
          onClick: _ctx.handlerMainButtonClick
        }), {
          default: vue.withCtx(() => [
            vue.renderSlot(_ctx.$slots, "default")
          ]),
          _: 3
        }, 16, ["size", "type", "disabled", "tabindex", "onClick"]),
        vue.createVNode(_component_el_button, vue.mergeProps({
          id: _ctx.triggerId,
          ref: "triggeringElementRef"
        }, _ctx.buttonProps, {
          role: "button",
          size: _ctx.dropdownSize,
          type: _ctx.type,
          class: _ctx.ns.e("caret-button"),
          disabled: _ctx.disabled,
          tabindex: _ctx.tabindex,
          "aria-label": _ctx.t("el.dropdown.toggleDropdown")
        }), {
          default: vue.withCtx(() => [
            vue.createVNode(_component_el_icon, {
              class: vue.normalizeClass(_ctx.ns.e("icon"))
            }, {
              default: vue.withCtx(() => [
                vue.createVNode(_component_arrow_down)
              ]),
              _: 1
            }, 8, ["class"])
          ]),
          _: 1
        }, 16, ["id", "size", "type", "class", "disabled", "tabindex", "aria-label"])
      ]),
      _: 3
    })) : vue.createCommentVNode("v-if", true)
  ], 2);
}
var Dropdown = /* @__PURE__ */ pluginVue_exportHelper["default"](_sfc_main, [["render", _sfc_render], ["__file", "dropdown.vue"]]);

exports["default"] = Dropdown;
//# sourceMappingURL=dropdown2.js.map
