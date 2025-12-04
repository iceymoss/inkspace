'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var vue = require('vue');
var lodashUnified = require('lodash-unified');
var index$3 = require('../../input/index.js');
var index$4 = require('../../tooltip/index.js');
var mention = require('./mention.js');
var helper = require('./helper.js');
var mentionDropdown = require('./mention-dropdown2.js');
var pluginVue_exportHelper = require('../../../_virtual/plugin-vue_export-helper.js');
var input = require('../../input/src/input.js');
var index = require('../../../hooks/use-namespace/index.js');
var useFormCommonProps = require('../../form/src/hooks/use-form-common-props.js');
var index$1 = require('../../../hooks/use-id/index.js');
var index$2 = require('../../../hooks/use-focus-controller/index.js');
var event = require('../../../constants/event.js');
var event$1 = require('../../../utils/dom/event.js');
var aria = require('../../../constants/aria.js');
var shared = require('@vue/shared');

const __default__ = vue.defineComponent({
  name: "ElMention",
  inheritAttrs: false
});
const _sfc_main = /* @__PURE__ */ vue.defineComponent({
  ...__default__,
  props: mention.mentionProps,
  emits: mention.mentionEmits,
  setup(__props, { expose, emit }) {
    const props = __props;
    const passInputProps = vue.computed(() => lodashUnified.pick(props, Object.keys(input.inputProps)));
    const ns = index.useNamespace("mention");
    const disabled = useFormCommonProps.useFormDisabled();
    const contentId = index$1.useId();
    const elInputRef = vue.ref();
    const tooltipRef = vue.ref();
    const dropdownRef = vue.ref();
    const visible = vue.ref(false);
    const cursorStyle = vue.ref();
    const mentionCtx = vue.ref();
    const computedPlacement = vue.computed(() => props.showArrow ? props.placement : `${props.placement}-start`);
    const computedFallbackPlacements = vue.computed(() => props.showArrow ? ["bottom", "top"] : ["bottom-start", "top-start"]);
    const aliasProps = vue.computed(() => ({
      ...mention.mentionDefaultProps,
      ...props.props
    }));
    const mapOption = (option) => {
      const base = {
        label: option[aliasProps.value.label],
        value: option[aliasProps.value.value],
        disabled: option[aliasProps.value.disabled]
      };
      return { ...option, ...base };
    };
    const options = vue.computed(() => props.options.map(mapOption));
    const filteredOptions = vue.computed(() => {
      const { filterOption } = props;
      if (!mentionCtx.value || !filterOption)
        return options.value;
      return options.value.filter((option) => filterOption(mentionCtx.value.pattern, option));
    });
    const dropdownVisible = vue.computed(() => {
      return visible.value && (!!filteredOptions.value.length || props.loading);
    });
    const hoveringId = vue.computed(() => {
      var _a;
      return `${contentId.value}-${(_a = dropdownRef.value) == null ? void 0 : _a.hoveringIndex}`;
    });
    const handleInputChange = (value) => {
      emit(event.UPDATE_MODEL_EVENT, value);
      emit(event.INPUT_EVENT, value);
      syncAfterCursorMove();
    };
    const handleInputKeyDown = (event$2) => {
      var _a, _b, _c, _d;
      if ((_a = elInputRef.value) == null ? void 0 : _a.isComposing)
        return;
      const code = event$1.getEventCode(event$2);
      switch (code) {
        case aria.EVENT_CODE.left:
        case aria.EVENT_CODE.right:
          syncAfterCursorMove();
          break;
        case aria.EVENT_CODE.up:
        case aria.EVENT_CODE.down:
          if (!visible.value)
            return;
          event$2.preventDefault();
          (_b = dropdownRef.value) == null ? void 0 : _b.navigateOptions(code === aria.EVENT_CODE.up ? "prev" : "next");
          break;
        case aria.EVENT_CODE.enter:
        case aria.EVENT_CODE.numpadEnter:
          if (!visible.value) {
            props.type !== "textarea" && syncAfterCursorMove();
            return;
          }
          event$2.preventDefault();
          if ((_c = dropdownRef.value) == null ? void 0 : _c.hoverOption) {
            (_d = dropdownRef.value) == null ? void 0 : _d.selectHoverOption();
          } else {
            visible.value = false;
          }
          break;
        case aria.EVENT_CODE.esc:
          if (!visible.value)
            return;
          event$2.preventDefault();
          visible.value = false;
          break;
        case aria.EVENT_CODE.backspace:
          if (props.whole && mentionCtx.value) {
            const { splitIndex, selectionEnd, pattern, prefixIndex, prefix } = mentionCtx.value;
            const inputEl = getInputEl();
            if (!inputEl)
              return;
            const inputValue = inputEl.value;
            const matchOption = options.value.find((item) => item.value === pattern);
            const isWhole = shared.isFunction(props.checkIsWhole) ? props.checkIsWhole(pattern, prefix) : matchOption;
            if (isWhole && splitIndex !== -1 && splitIndex + 1 === selectionEnd) {
              event$2.preventDefault();
              const newValue = inputValue.slice(0, prefixIndex) + inputValue.slice(splitIndex + 1);
              emit(event.UPDATE_MODEL_EVENT, newValue);
              emit(event.INPUT_EVENT, newValue);
              emit("whole-remove", pattern, prefix);
              const newSelectionEnd = prefixIndex;
              vue.nextTick(() => {
                inputEl.selectionStart = newSelectionEnd;
                inputEl.selectionEnd = newSelectionEnd;
                syncDropdownVisible();
              });
            }
          }
      }
    };
    const { wrapperRef } = index$2.useFocusController(elInputRef, {
      disabled,
      afterFocus() {
        syncAfterCursorMove();
      },
      beforeBlur(event) {
        var _a;
        return (_a = tooltipRef.value) == null ? void 0 : _a.isFocusInsideContent(event);
      },
      afterBlur() {
        visible.value = false;
      }
    });
    const handleInputMouseDown = () => {
      syncAfterCursorMove();
    };
    const getOriginalOption = (mentionOption) => {
      return props.options.find((option) => {
        return mentionOption.value === option[aliasProps.value.value];
      });
    };
    const handleSelect = (item) => {
      if (!mentionCtx.value)
        return;
      const inputEl = getInputEl();
      if (!inputEl)
        return;
      const inputValue = inputEl.value;
      const { split } = props;
      const newEndPart = inputValue.slice(mentionCtx.value.end);
      const alreadySeparated = newEndPart.startsWith(split);
      const newMiddlePart = `${item.value}${alreadySeparated ? "" : split}`;
      const newValue = inputValue.slice(0, mentionCtx.value.start) + newMiddlePart + newEndPart;
      emit(event.UPDATE_MODEL_EVENT, newValue);
      emit(event.INPUT_EVENT, newValue);
      emit("select", getOriginalOption(item), mentionCtx.value.prefix);
      const newSelectionEnd = mentionCtx.value.start + newMiddlePart.length + (alreadySeparated ? 1 : 0);
      vue.nextTick(() => {
        inputEl.selectionStart = newSelectionEnd;
        inputEl.selectionEnd = newSelectionEnd;
        inputEl.focus();
        syncDropdownVisible();
      });
    };
    const getInputEl = () => {
      var _a, _b;
      return props.type === "textarea" ? (_a = elInputRef.value) == null ? void 0 : _a.textarea : (_b = elInputRef.value) == null ? void 0 : _b.input;
    };
    const syncAfterCursorMove = () => {
      setTimeout(() => {
        syncCursor();
        syncDropdownVisible();
        vue.nextTick(() => {
          var _a;
          return (_a = tooltipRef.value) == null ? void 0 : _a.updatePopper();
        });
      }, 0);
    };
    const syncCursor = () => {
      const inputEl = getInputEl();
      if (!inputEl)
        return;
      const caretPosition = helper.getCursorPosition(inputEl);
      const inputRect = inputEl.getBoundingClientRect();
      const wrapperRect = wrapperRef.value.getBoundingClientRect();
      cursorStyle.value = {
        position: "absolute",
        width: 0,
        height: `${caretPosition.height}px`,
        left: `${caretPosition.left + inputRect.left - wrapperRect.left}px`,
        top: `${caretPosition.top + inputRect.top - wrapperRect.top}px`
      };
    };
    const syncDropdownVisible = () => {
      const inputEl = getInputEl();
      if (document.activeElement !== inputEl) {
        visible.value = false;
        return;
      }
      const { prefix, split } = props;
      mentionCtx.value = helper.getMentionCtx(inputEl, prefix, split);
      if (mentionCtx.value && mentionCtx.value.splitIndex === -1) {
        visible.value = true;
        emit("search", mentionCtx.value.pattern, mentionCtx.value.prefix);
        return;
      }
      visible.value = false;
    };
    expose({
      input: elInputRef,
      tooltip: tooltipRef,
      dropdownVisible
    });
    return (_ctx, _cache) => {
      return vue.openBlock(), vue.createElementBlock("div", {
        ref_key: "wrapperRef",
        ref: wrapperRef,
        class: vue.normalizeClass(vue.unref(ns).b())
      }, [
        vue.createVNode(vue.unref(index$3.ElInput), vue.mergeProps(vue.mergeProps(vue.unref(passInputProps), _ctx.$attrs), {
          ref_key: "elInputRef",
          ref: elInputRef,
          "model-value": _ctx.modelValue,
          disabled: vue.unref(disabled),
          role: vue.unref(dropdownVisible) ? "combobox" : void 0,
          "aria-activedescendant": vue.unref(dropdownVisible) ? vue.unref(hoveringId) || "" : void 0,
          "aria-controls": vue.unref(dropdownVisible) ? vue.unref(contentId) : void 0,
          "aria-expanded": vue.unref(dropdownVisible) || void 0,
          "aria-label": _ctx.ariaLabel,
          "aria-autocomplete": vue.unref(dropdownVisible) ? "none" : void 0,
          "aria-haspopup": vue.unref(dropdownVisible) ? "listbox" : void 0,
          onInput: handleInputChange,
          onKeydown: handleInputKeyDown,
          onMousedown: handleInputMouseDown
        }), vue.createSlots({
          _: 2
        }, [
          vue.renderList(_ctx.$slots, (_, name) => {
            return {
              name,
              fn: vue.withCtx((slotProps) => [
                vue.renderSlot(_ctx.$slots, name, vue.normalizeProps(vue.guardReactiveProps(slotProps)))
              ])
            };
          })
        ]), 1040, ["model-value", "disabled", "role", "aria-activedescendant", "aria-controls", "aria-expanded", "aria-label", "aria-autocomplete", "aria-haspopup"]),
        vue.createVNode(vue.unref(index$4.ElTooltip), {
          ref_key: "tooltipRef",
          ref: tooltipRef,
          visible: vue.unref(dropdownVisible),
          "popper-class": [vue.unref(ns).e("popper"), _ctx.popperClass],
          "popper-style": _ctx.popperStyle,
          "popper-options": _ctx.popperOptions,
          placement: vue.unref(computedPlacement),
          "fallback-placements": vue.unref(computedFallbackPlacements),
          effect: "light",
          pure: "",
          offset: _ctx.offset,
          "show-arrow": _ctx.showArrow
        }, {
          default: vue.withCtx(() => [
            vue.createElementVNode("div", {
              style: vue.normalizeStyle(cursorStyle.value)
            }, null, 4)
          ]),
          content: vue.withCtx(() => {
            var _a;
            return [
              vue.createVNode(mentionDropdown["default"], {
                ref_key: "dropdownRef",
                ref: dropdownRef,
                options: vue.unref(filteredOptions),
                disabled: vue.unref(disabled),
                loading: _ctx.loading,
                "content-id": vue.unref(contentId),
                "aria-label": _ctx.ariaLabel,
                onSelect: handleSelect,
                onClick: vue.withModifiers((_a = elInputRef.value) == null ? void 0 : _a.focus, ["stop"])
              }, vue.createSlots({
                _: 2
              }, [
                vue.renderList(_ctx.$slots, (_, name) => {
                  return {
                    name,
                    fn: vue.withCtx((slotProps) => [
                      vue.renderSlot(_ctx.$slots, name, vue.normalizeProps(vue.guardReactiveProps(slotProps)))
                    ])
                  };
                })
              ]), 1032, ["options", "disabled", "loading", "content-id", "aria-label", "onClick"])
            ];
          }),
          _: 3
        }, 8, ["visible", "popper-class", "popper-style", "popper-options", "placement", "fallback-placements", "offset", "show-arrow"])
      ], 2);
    };
  }
});
var Mention = /* @__PURE__ */ pluginVue_exportHelper["default"](_sfc_main, [["__file", "mention.vue"]]);

exports["default"] = Mention;
//# sourceMappingURL=mention2.js.map
