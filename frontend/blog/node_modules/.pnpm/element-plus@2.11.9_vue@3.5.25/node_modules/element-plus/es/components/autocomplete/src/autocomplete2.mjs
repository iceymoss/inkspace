import { defineComponent, computed, useAttrs, ref, onBeforeUnmount, onMounted, openBlock, createBlock, unref, withCtx, createElementVNode, normalizeClass, normalizeStyle, createElementBlock, withModifiers, renderSlot, createCommentVNode, createVNode, Fragment, renderList, createTextVNode, toDisplayString, mergeProps, createSlots } from 'vue';
import { pick } from 'lodash-unified';
import { useDebounceFn, onClickOutside } from '@vueuse/core';
import { Loading } from '@element-plus/icons-vue';
import { ElInput } from '../../input/index.mjs';
import { ElScrollbar } from '../../scrollbar/index.mjs';
import { ElTooltip } from '../../tooltip/index.mjs';
import { ElIcon } from '../../icon/index.mjs';
import { autocompleteProps, autocompleteEmits } from './autocomplete.mjs';
import _export_sfc from '../../../_virtual/plugin-vue_export-helper.mjs';
import { inputProps } from '../../input/src/input.mjs';
import { useFormDisabled } from '../../form/src/hooks/use-form-common-props.mjs';
import { useNamespace } from '../../../hooks/use-namespace/index.mjs';
import { useId } from '../../../hooks/use-id/index.mjs';
import { isArray } from '@vue/shared';
import { INPUT_EVENT, UPDATE_MODEL_EVENT, CHANGE_EVENT } from '../../../constants/event.mjs';
import { getEventCode } from '../../../utils/dom/event.mjs';
import { EVENT_CODE } from '../../../constants/aria.mjs';
import { throwError } from '../../../utils/error.mjs';

const COMPONENT_NAME = "ElAutocomplete";
const __default__ = defineComponent({
  name: COMPONENT_NAME,
  inheritAttrs: false
});
const _sfc_main = /* @__PURE__ */ defineComponent({
  ...__default__,
  props: autocompleteProps,
  emits: autocompleteEmits,
  setup(__props, { expose, emit }) {
    const props = __props;
    const passInputProps = computed(() => pick(props, Object.keys(inputProps)));
    const rawAttrs = useAttrs();
    const disabled = useFormDisabled();
    const ns = useNamespace("autocomplete");
    const inputRef = ref();
    const regionRef = ref();
    const popperRef = ref();
    const listboxRef = ref();
    let readonly = false;
    let ignoreFocusEvent = false;
    const suggestions = ref([]);
    const highlightedIndex = ref(-1);
    const dropdownWidth = ref("");
    const activated = ref(false);
    const suggestionDisabled = ref(false);
    const loading = ref(false);
    const listboxId = useId();
    const styles = computed(() => rawAttrs.style);
    const suggestionVisible = computed(() => {
      const isValidData = suggestions.value.length > 0;
      return (isValidData || loading.value) && activated.value;
    });
    const suggestionLoading = computed(() => !props.hideLoading && loading.value);
    const refInput = computed(() => {
      if (inputRef.value) {
        return Array.from(inputRef.value.$el.querySelectorAll("input"));
      }
      return [];
    });
    const onSuggestionShow = () => {
      if (suggestionVisible.value) {
        dropdownWidth.value = `${inputRef.value.$el.offsetWidth}px`;
      }
    };
    const onHide = () => {
      highlightedIndex.value = -1;
    };
    const getData = async (queryString) => {
      if (suggestionDisabled.value)
        return;
      const cb = (suggestionList) => {
        loading.value = false;
        if (suggestionDisabled.value)
          return;
        if (isArray(suggestionList)) {
          suggestions.value = suggestionList;
          highlightedIndex.value = props.highlightFirstItem ? 0 : -1;
        } else {
          throwError(COMPONENT_NAME, "autocomplete suggestions must be an array");
        }
      };
      loading.value = true;
      if (isArray(props.fetchSuggestions)) {
        cb(props.fetchSuggestions);
      } else {
        const result = await props.fetchSuggestions(queryString, cb);
        if (isArray(result))
          cb(result);
      }
    };
    const debounce = computed(() => props.debounce);
    const debouncedGetData = useDebounceFn(getData, debounce);
    const handleInput = (value) => {
      const valuePresented = !!value;
      emit(INPUT_EVENT, value);
      emit(UPDATE_MODEL_EVENT, value);
      suggestionDisabled.value = false;
      activated.value || (activated.value = valuePresented);
      if (!props.triggerOnFocus && !value) {
        suggestionDisabled.value = true;
        suggestions.value = [];
        return;
      }
      debouncedGetData(value);
    };
    const handleMouseDown = (event) => {
      var _a;
      if (disabled.value)
        return;
      if (((_a = event.target) == null ? void 0 : _a.tagName) !== "INPUT" || refInput.value.includes(document.activeElement)) {
        activated.value = true;
      }
    };
    const handleChange = (value) => {
      emit(CHANGE_EVENT, value);
    };
    const handleFocus = (evt) => {
      var _a;
      if (!ignoreFocusEvent) {
        activated.value = true;
        emit("focus", evt);
        const queryString = (_a = props.modelValue) != null ? _a : "";
        if (props.triggerOnFocus && !readonly) {
          debouncedGetData(String(queryString));
        }
      } else {
        ignoreFocusEvent = false;
      }
    };
    const handleBlur = (evt) => {
      setTimeout(() => {
        var _a;
        if ((_a = popperRef.value) == null ? void 0 : _a.isFocusInsideContent()) {
          ignoreFocusEvent = true;
          return;
        }
        activated.value && close();
        emit("blur", evt);
      });
    };
    const handleClear = () => {
      activated.value = false;
      emit(UPDATE_MODEL_EVENT, "");
      emit("clear");
    };
    const handleKeyEnter = async () => {
      var _a;
      if ((_a = inputRef.value) == null ? void 0 : _a.isComposing) {
        return;
      }
      if (suggestionVisible.value && highlightedIndex.value >= 0 && highlightedIndex.value < suggestions.value.length) {
        handleSelect(suggestions.value[highlightedIndex.value]);
      } else {
        if (props.selectWhenUnmatched) {
          emit("select", { value: props.modelValue });
          suggestions.value = [];
          highlightedIndex.value = -1;
        }
        activated.value = true;
        debouncedGetData(String(props.modelValue));
      }
    };
    const handleKeyEscape = (evt) => {
      if (suggestionVisible.value) {
        evt.preventDefault();
        evt.stopPropagation();
        close();
      }
    };
    const close = () => {
      activated.value = false;
    };
    const focus = () => {
      var _a;
      (_a = inputRef.value) == null ? void 0 : _a.focus();
    };
    const blur = () => {
      var _a;
      (_a = inputRef.value) == null ? void 0 : _a.blur();
    };
    const handleSelect = async (item) => {
      emit(INPUT_EVENT, item[props.valueKey]);
      emit(UPDATE_MODEL_EVENT, item[props.valueKey]);
      emit("select", item);
      suggestions.value = [];
      highlightedIndex.value = -1;
    };
    const highlight = (index) => {
      var _a, _b;
      if (!suggestionVisible.value || loading.value)
        return;
      if (index < 0) {
        if (!props.loopNavigation) {
          highlightedIndex.value = -1;
          return;
        }
        index = suggestions.value.length - 1;
      }
      if (index >= suggestions.value.length) {
        index = props.loopNavigation ? 0 : suggestions.value.length - 1;
      }
      const [suggestion, suggestionList] = getSuggestionContext();
      const highlightItem = suggestionList[index];
      const scrollTop = suggestion.scrollTop;
      const { offsetTop, scrollHeight } = highlightItem;
      if (offsetTop + scrollHeight > scrollTop + suggestion.clientHeight) {
        suggestion.scrollTop = offsetTop + scrollHeight - suggestion.clientHeight;
      }
      if (offsetTop < scrollTop) {
        suggestion.scrollTop = offsetTop;
      }
      highlightedIndex.value = index;
      (_b = (_a = inputRef.value) == null ? void 0 : _a.ref) == null ? void 0 : _b.setAttribute("aria-activedescendant", `${listboxId.value}-item-${highlightedIndex.value}`);
    };
    const getSuggestionContext = () => {
      const suggestion = regionRef.value.querySelector(`.${ns.be("suggestion", "wrap")}`);
      const suggestionList = suggestion.querySelectorAll(`.${ns.be("suggestion", "list")} li`);
      return [suggestion, suggestionList];
    };
    const stopHandle = onClickOutside(listboxRef, () => {
      var _a;
      if ((_a = popperRef.value) == null ? void 0 : _a.isFocusInsideContent())
        return;
      suggestionVisible.value && close();
    });
    const handleKeydown = (e) => {
      const code = getEventCode(e);
      switch (code) {
        case EVENT_CODE.up:
          e.preventDefault();
          highlight(highlightedIndex.value - 1);
          break;
        case EVENT_CODE.down:
          e.preventDefault();
          highlight(highlightedIndex.value + 1);
          break;
        case EVENT_CODE.enter:
        case EVENT_CODE.numpadEnter:
          e.preventDefault();
          handleKeyEnter();
          break;
        case EVENT_CODE.tab:
          close();
          break;
        case EVENT_CODE.esc:
          handleKeyEscape(e);
          break;
        case EVENT_CODE.home:
          e.preventDefault();
          highlight(0);
          break;
        case EVENT_CODE.end:
          e.preventDefault();
          highlight(suggestions.value.length - 1);
          break;
        case EVENT_CODE.pageUp:
          e.preventDefault();
          highlight(Math.max(0, highlightedIndex.value - 10));
          break;
        case EVENT_CODE.pageDown:
          e.preventDefault();
          highlight(Math.min(suggestions.value.length - 1, highlightedIndex.value + 10));
          break;
      }
    };
    onBeforeUnmount(() => {
      stopHandle == null ? void 0 : stopHandle();
    });
    onMounted(() => {
      var _a;
      const inputElement = (_a = inputRef.value) == null ? void 0 : _a.ref;
      if (!inputElement)
        return;
      [
        { key: "role", value: "textbox" },
        { key: "aria-autocomplete", value: "list" },
        { key: "aria-controls", value: "id" },
        {
          key: "aria-activedescendant",
          value: `${listboxId.value}-item-${highlightedIndex.value}`
        }
      ].forEach(({ key, value }) => inputElement.setAttribute(key, value));
      readonly = inputElement.hasAttribute("readonly");
    });
    expose({
      highlightedIndex,
      activated,
      loading,
      inputRef,
      popperRef,
      suggestions,
      handleSelect,
      handleKeyEnter,
      focus,
      blur,
      close,
      highlight,
      getData
    });
    return (_ctx, _cache) => {
      return openBlock(), createBlock(unref(ElTooltip), {
        ref_key: "popperRef",
        ref: popperRef,
        visible: unref(suggestionVisible),
        placement: _ctx.placement,
        "fallback-placements": ["bottom-start", "top-start"],
        "popper-class": [unref(ns).e("popper"), _ctx.popperClass],
        "popper-style": _ctx.popperStyle,
        teleported: _ctx.teleported,
        "append-to": _ctx.appendTo,
        "gpu-acceleration": false,
        pure: "",
        "manual-mode": "",
        effect: "light",
        trigger: "click",
        transition: `${unref(ns).namespace.value}-zoom-in-top`,
        persistent: "",
        role: "listbox",
        onBeforeShow: onSuggestionShow,
        onHide
      }, {
        content: withCtx(() => [
          createElementVNode("div", {
            ref_key: "regionRef",
            ref: regionRef,
            class: normalizeClass([unref(ns).b("suggestion"), unref(ns).is("loading", unref(suggestionLoading))]),
            style: normalizeStyle({
              [_ctx.fitInputWidth ? "width" : "minWidth"]: dropdownWidth.value,
              outline: "none"
            }),
            role: "region"
          }, [
            _ctx.$slots.header ? (openBlock(), createElementBlock("div", {
              key: 0,
              class: normalizeClass(unref(ns).be("suggestion", "header")),
              onClick: withModifiers(() => {
              }, ["stop"])
            }, [
              renderSlot(_ctx.$slots, "header")
            ], 10, ["onClick"])) : createCommentVNode("v-if", true),
            createVNode(unref(ElScrollbar), {
              id: unref(listboxId),
              tag: "ul",
              "wrap-class": unref(ns).be("suggestion", "wrap"),
              "view-class": unref(ns).be("suggestion", "list"),
              role: "listbox"
            }, {
              default: withCtx(() => [
                unref(suggestionLoading) ? (openBlock(), createElementBlock("li", { key: 0 }, [
                  renderSlot(_ctx.$slots, "loading", {}, () => [
                    createVNode(unref(ElIcon), {
                      class: normalizeClass(unref(ns).is("loading"))
                    }, {
                      default: withCtx(() => [
                        createVNode(unref(Loading))
                      ]),
                      _: 1
                    }, 8, ["class"])
                  ])
                ])) : (openBlock(true), createElementBlock(Fragment, { key: 1 }, renderList(suggestions.value, (item, index) => {
                  return openBlock(), createElementBlock("li", {
                    id: `${unref(listboxId)}-item-${index}`,
                    key: index,
                    class: normalizeClass({ highlighted: highlightedIndex.value === index }),
                    role: "option",
                    "aria-selected": highlightedIndex.value === index,
                    onClick: ($event) => handleSelect(item)
                  }, [
                    renderSlot(_ctx.$slots, "default", { item }, () => [
                      createTextVNode(toDisplayString(item[_ctx.valueKey]), 1)
                    ])
                  ], 10, ["id", "aria-selected", "onClick"]);
                }), 128))
              ]),
              _: 3
            }, 8, ["id", "wrap-class", "view-class"]),
            _ctx.$slots.footer ? (openBlock(), createElementBlock("div", {
              key: 1,
              class: normalizeClass(unref(ns).be("suggestion", "footer")),
              onClick: withModifiers(() => {
              }, ["stop"])
            }, [
              renderSlot(_ctx.$slots, "footer")
            ], 10, ["onClick"])) : createCommentVNode("v-if", true)
          ], 6)
        ]),
        default: withCtx(() => [
          createElementVNode("div", {
            ref_key: "listboxRef",
            ref: listboxRef,
            class: normalizeClass([unref(ns).b(), _ctx.$attrs.class]),
            style: normalizeStyle(unref(styles)),
            role: "combobox",
            "aria-haspopup": "listbox",
            "aria-expanded": unref(suggestionVisible),
            "aria-owns": unref(listboxId)
          }, [
            createVNode(unref(ElInput), mergeProps({
              ref_key: "inputRef",
              ref: inputRef
            }, mergeProps(unref(passInputProps), _ctx.$attrs), {
              "model-value": _ctx.modelValue,
              disabled: unref(disabled),
              onInput: handleInput,
              onChange: handleChange,
              onFocus: handleFocus,
              onBlur: handleBlur,
              onClear: handleClear,
              onKeydown: handleKeydown,
              onMousedown: handleMouseDown
            }), createSlots({
              _: 2
            }, [
              _ctx.$slots.prepend ? {
                name: "prepend",
                fn: withCtx(() => [
                  renderSlot(_ctx.$slots, "prepend")
                ])
              } : void 0,
              _ctx.$slots.append ? {
                name: "append",
                fn: withCtx(() => [
                  renderSlot(_ctx.$slots, "append")
                ])
              } : void 0,
              _ctx.$slots.prefix ? {
                name: "prefix",
                fn: withCtx(() => [
                  renderSlot(_ctx.$slots, "prefix")
                ])
              } : void 0,
              _ctx.$slots.suffix ? {
                name: "suffix",
                fn: withCtx(() => [
                  renderSlot(_ctx.$slots, "suffix")
                ])
              } : void 0
            ]), 1040, ["model-value", "disabled"])
          ], 14, ["aria-expanded", "aria-owns"])
        ]),
        _: 3
      }, 8, ["visible", "placement", "popper-class", "popper-style", "teleported", "append-to", "transition"]);
    };
  }
});
var Autocomplete = /* @__PURE__ */ _export_sfc(_sfc_main, [["__file", "autocomplete.vue"]]);

export { Autocomplete as default };
//# sourceMappingURL=autocomplete2.mjs.map
