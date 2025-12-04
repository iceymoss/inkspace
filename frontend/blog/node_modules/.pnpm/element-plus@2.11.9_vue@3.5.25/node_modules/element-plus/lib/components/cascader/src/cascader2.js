'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var vue = require('vue');
var lodashUnified = require('lodash-unified');
var core = require('@vueuse/core');
var index$b = require('../../cascader-panel/index.js');
var index$6 = require('../../input/index.js');
var index$5 = require('../../tooltip/index.js');
var index$9 = require('../../scrollbar/index.js');
var index$8 = require('../../tag/index.js');
var index$7 = require('../../icon/index.js');
var iconsVue = require('@element-plus/icons-vue');
var cascader = require('./cascader.js');
var pluginVue_exportHelper = require('../../../_virtual/plugin-vue_export-helper.js');
var index$a = require('../../../directives/click-outside/index.js');
var index = require('../../../hooks/use-namespace/index.js');
var index$1 = require('../../../hooks/use-locale/index.js');
var useFormItem = require('../../form/src/hooks/use-form-item.js');
var useFormCommonProps = require('../../form/src/hooks/use-form-common-props.js');
var index$2 = require('../../../hooks/use-empty-values/index.js');
var index$3 = require('../../../hooks/use-composition/index.js');
var index$4 = require('../../../hooks/use-focus-controller/index.js');
var event = require('../../../constants/event.js');
var error = require('../../../utils/error.js');
var shared = require('@vue/shared');
var event$1 = require('../../../utils/dom/event.js');
var aria = require('../../../constants/aria.js');
var aria$1 = require('../../../utils/dom/aria.js');

const __default__ = vue.defineComponent({
  name: "ElCascader"
});
const _sfc_main = /* @__PURE__ */ vue.defineComponent({
  ...__default__,
  props: cascader.cascaderProps,
  emits: cascader.cascaderEmits,
  setup(__props, { expose, emit }) {
    const props = __props;
    const popperOptions = {
      modifiers: [
        {
          name: "arrowPosition",
          enabled: true,
          phase: "main",
          fn: ({ state }) => {
            const { modifiersData, placement } = state;
            if (["right", "left", "bottom", "top"].includes(placement))
              return;
            if (modifiersData.arrow) {
              modifiersData.arrow.x = 35;
            }
          },
          requires: ["arrow"]
        }
      ]
    };
    const attrs = vue.useAttrs();
    let inputInitialHeight = 0;
    let pressDeleteCount = 0;
    const nsCascader = index.useNamespace("cascader");
    const nsInput = index.useNamespace("input");
    const { t } = index$1.useLocale();
    const { formItem } = useFormItem.useFormItem();
    const isDisabled = useFormCommonProps.useFormDisabled();
    const { valueOnClear } = index$2.useEmptyValues(props);
    const { isComposing, handleComposition } = index$3.useComposition({
      afterComposition(event) {
        var _a;
        const text = (_a = event.target) == null ? void 0 : _a.value;
        handleInput(text);
      }
    });
    const tooltipRef = vue.ref();
    const tagTooltipRef = vue.ref();
    const inputRef = vue.ref();
    const tagWrapper = vue.ref();
    const cascaderPanelRef = vue.ref();
    const suggestionPanel = vue.ref();
    const popperVisible = vue.ref(false);
    const inputHover = vue.ref(false);
    const filtering = vue.ref(false);
    const inputValue = vue.ref("");
    const searchInputValue = vue.ref("");
    const tags = vue.ref([]);
    const suggestions = vue.ref([]);
    const showTagList = vue.computed(() => {
      if (!props.props.multiple) {
        return [];
      }
      return props.collapseTags ? tags.value.slice(0, props.maxCollapseTags) : tags.value;
    });
    const collapseTagList = vue.computed(() => {
      if (!props.props.multiple) {
        return [];
      }
      return props.collapseTags ? tags.value.slice(props.maxCollapseTags) : [];
    });
    const cascaderStyle = vue.computed(() => {
      return attrs.style;
    });
    const inputPlaceholder = vue.computed(() => {
      var _a;
      return (_a = props.placeholder) != null ? _a : t("el.cascader.placeholder");
    });
    const currentPlaceholder = vue.computed(() => searchInputValue.value || tags.value.length > 0 || isComposing.value ? "" : inputPlaceholder.value);
    const realSize = useFormCommonProps.useFormSize();
    const tagSize = vue.computed(() => realSize.value === "small" ? "small" : "default");
    const multiple = vue.computed(() => !!props.props.multiple);
    const readonly = vue.computed(() => !props.filterable || multiple.value);
    const searchKeyword = vue.computed(() => multiple.value ? searchInputValue.value : inputValue.value);
    const checkedNodes = vue.computed(() => {
      var _a;
      return ((_a = cascaderPanelRef.value) == null ? void 0 : _a.checkedNodes) || [];
    });
    const { wrapperRef, isFocused, handleBlur } = index$4.useFocusController(inputRef, {
      disabled: isDisabled,
      beforeBlur(event) {
        var _a, _b;
        return ((_a = tooltipRef.value) == null ? void 0 : _a.isFocusInsideContent(event)) || ((_b = tagTooltipRef.value) == null ? void 0 : _b.isFocusInsideContent(event));
      },
      afterBlur() {
        var _a;
        if (props.validateEvent) {
          (_a = formItem == null ? void 0 : formItem.validate) == null ? void 0 : _a.call(formItem, "blur").catch((err) => error.debugWarn(err));
        }
      }
    });
    const clearBtnVisible = vue.computed(() => {
      if (!props.clearable || isDisabled.value || filtering.value || !inputHover.value && !isFocused.value)
        return false;
      return !!checkedNodes.value.length;
    });
    const presentText = vue.computed(() => {
      const { showAllLevels, separator } = props;
      const nodes = checkedNodes.value;
      return nodes.length ? multiple.value ? "" : nodes[0].calcText(showAllLevels, separator) : "";
    });
    const validateState = vue.computed(() => (formItem == null ? void 0 : formItem.validateState) || "");
    const checkedValue = vue.computed({
      get() {
        return lodashUnified.cloneDeep(props.modelValue);
      },
      set(val) {
        const value = val != null ? val : valueOnClear.value;
        emit(event.UPDATE_MODEL_EVENT, value);
        emit(event.CHANGE_EVENT, value);
        if (props.validateEvent) {
          formItem == null ? void 0 : formItem.validate("change").catch((err) => error.debugWarn(err));
        }
      }
    });
    const cascaderKls = vue.computed(() => {
      return [
        nsCascader.b(),
        nsCascader.m(realSize.value),
        nsCascader.is("disabled", isDisabled.value),
        attrs.class
      ];
    });
    const cascaderIconKls = vue.computed(() => {
      return [
        nsInput.e("icon"),
        "icon-arrow-down",
        nsCascader.is("reverse", popperVisible.value)
      ];
    });
    const inputClass = vue.computed(() => nsCascader.is("focus", isFocused.value));
    const contentRef = vue.computed(() => {
      var _a, _b;
      return (_b = (_a = tooltipRef.value) == null ? void 0 : _a.popperRef) == null ? void 0 : _b.contentRef;
    });
    const handleClickOutside = (event) => {
      if (isFocused.value) {
        const _event = new FocusEvent("blur", event);
        handleBlur(_event);
      }
      togglePopperVisible(false);
    };
    const togglePopperVisible = (visible) => {
      var _a, _b, _c;
      if (isDisabled.value)
        return;
      visible = visible != null ? visible : !popperVisible.value;
      if (visible !== popperVisible.value) {
        popperVisible.value = visible;
        (_b = (_a = inputRef.value) == null ? void 0 : _a.input) == null ? void 0 : _b.setAttribute("aria-expanded", `${visible}`);
        if (visible) {
          updatePopperPosition();
          vue.nextTick((_c = cascaderPanelRef.value) == null ? void 0 : _c.scrollToExpandingNode);
        } else if (props.filterable) {
          syncPresentTextValue();
        }
        emit("visibleChange", visible);
      }
    };
    const updatePopperPosition = () => {
      vue.nextTick(() => {
        var _a;
        (_a = tooltipRef.value) == null ? void 0 : _a.updatePopper();
      });
    };
    const hideSuggestionPanel = () => {
      filtering.value = false;
    };
    const genTag = (node) => {
      const { showAllLevels, separator } = props;
      return {
        node,
        key: node.uid,
        text: node.calcText(showAllLevels, separator),
        hitState: false,
        closable: !isDisabled.value && !node.isDisabled
      };
    };
    const deleteTag = (tag) => {
      var _a;
      const node = tag.node;
      node.doCheck(false);
      (_a = cascaderPanelRef.value) == null ? void 0 : _a.calculateCheckedValue();
      emit("removeTag", node.valueByOption);
    };
    const getStrategyCheckedNodes = () => {
      switch (props.showCheckedStrategy) {
        case "child":
          return checkedNodes.value;
        case "parent": {
          const clickedNodes = getCheckedNodes(false);
          const clickedNodesValue = clickedNodes.map((o) => o.value);
          const parentNodes = clickedNodes.filter((o) => !o.parent || !clickedNodesValue.includes(o.parent.value));
          return parentNodes;
        }
        default:
          return [];
      }
    };
    const calculatePresentTags = () => {
      if (!multiple.value)
        return;
      const nodes = getStrategyCheckedNodes();
      const allTags = [];
      nodes.forEach((node) => allTags.push(genTag(node)));
      tags.value = allTags;
    };
    const calculateSuggestions = () => {
      var _a, _b;
      const { filterMethod, showAllLevels, separator } = props;
      const res = (_b = (_a = cascaderPanelRef.value) == null ? void 0 : _a.getFlattedNodes(!props.props.checkStrictly)) == null ? void 0 : _b.filter((node) => {
        if (node.isDisabled)
          return false;
        node.calcText(showAllLevels, separator);
        return filterMethod(node, searchKeyword.value);
      });
      if (multiple.value) {
        tags.value.forEach((tag) => {
          tag.hitState = false;
        });
      }
      filtering.value = true;
      suggestions.value = res;
      updatePopperPosition();
    };
    const focusFirstNode = () => {
      var _a;
      let firstNode;
      if (filtering.value && suggestionPanel.value) {
        firstNode = suggestionPanel.value.$el.querySelector(`.${nsCascader.e("suggestion-item")}`);
      } else {
        firstNode = (_a = cascaderPanelRef.value) == null ? void 0 : _a.$el.querySelector(`.${nsCascader.b("node")}[tabindex="-1"]`);
      }
      if (firstNode) {
        firstNode.focus();
        !filtering.value && firstNode.click();
      }
    };
    const updateStyle = () => {
      var _a, _b;
      const inputInner = (_a = inputRef.value) == null ? void 0 : _a.input;
      const tagWrapperEl = tagWrapper.value;
      const suggestionPanelEl = (_b = suggestionPanel.value) == null ? void 0 : _b.$el;
      if (!core.isClient || !inputInner)
        return;
      if (suggestionPanelEl) {
        const suggestionList = suggestionPanelEl.querySelector(`.${nsCascader.e("suggestion-list")}`);
        suggestionList.style.minWidth = `${inputInner.offsetWidth}px`;
      }
      if (tagWrapperEl) {
        const { offsetHeight } = tagWrapperEl;
        const height = tags.value.length > 0 ? `${Math.max(offsetHeight, inputInitialHeight) - 2}px` : `${inputInitialHeight}px`;
        inputInner.style.height = height;
        updatePopperPosition();
      }
    };
    const getCheckedNodes = (leafOnly) => {
      var _a;
      return (_a = cascaderPanelRef.value) == null ? void 0 : _a.getCheckedNodes(leafOnly);
    };
    const handleExpandChange = (value) => {
      updatePopperPosition();
      emit("expandChange", value);
    };
    const handleKeyDown = (e) => {
      if (isComposing.value)
        return;
      const code = event$1.getEventCode(e);
      switch (code) {
        case aria.EVENT_CODE.enter:
        case aria.EVENT_CODE.numpadEnter:
          togglePopperVisible();
          break;
        case aria.EVENT_CODE.down:
          togglePopperVisible(true);
          vue.nextTick(focusFirstNode);
          e.preventDefault();
          break;
        case aria.EVENT_CODE.esc:
          if (popperVisible.value === true) {
            e.preventDefault();
            e.stopPropagation();
            togglePopperVisible(false);
          }
          break;
        case aria.EVENT_CODE.tab:
          togglePopperVisible(false);
          break;
      }
    };
    const handleClear = () => {
      var _a;
      (_a = cascaderPanelRef.value) == null ? void 0 : _a.clearCheckedNodes();
      if (!popperVisible.value && props.filterable) {
        syncPresentTextValue();
      }
      togglePopperVisible(false);
      emit("clear");
    };
    const syncPresentTextValue = () => {
      const { value } = presentText;
      inputValue.value = value;
      searchInputValue.value = value;
    };
    const handleSuggestionClick = (node) => {
      var _a, _b;
      const { checked } = node;
      if (multiple.value) {
        (_a = cascaderPanelRef.value) == null ? void 0 : _a.handleCheckChange(node, !checked, false);
      } else {
        !checked && ((_b = cascaderPanelRef.value) == null ? void 0 : _b.handleCheckChange(node, true, false));
        togglePopperVisible(false);
      }
    };
    const handleSuggestionKeyDown = (e) => {
      const target = e.target;
      const code = event$1.getEventCode(e);
      switch (code) {
        case aria.EVENT_CODE.up:
        case aria.EVENT_CODE.down: {
          e.preventDefault();
          const distance = code === aria.EVENT_CODE.up ? -1 : 1;
          aria$1.focusNode(aria$1.getSibling(target, distance, `.${nsCascader.e("suggestion-item")}[tabindex="-1"]`));
          break;
        }
        case aria.EVENT_CODE.enter:
        case aria.EVENT_CODE.numpadEnter:
          target.click();
          break;
      }
    };
    const handleDelete = () => {
      const lastTag = tags.value[tags.value.length - 1];
      pressDeleteCount = searchInputValue.value ? 0 : pressDeleteCount + 1;
      if (!lastTag || !pressDeleteCount || props.collapseTags && tags.value.length > 1)
        return;
      if (lastTag.hitState) {
        deleteTag(lastTag);
      } else {
        lastTag.hitState = true;
      }
    };
    const debounce = vue.computed(() => props.debounce);
    const handleFilter = core.useDebounceFn(() => {
      const { value } = searchKeyword;
      if (!value)
        return;
      const passed = props.beforeFilter(value);
      if (shared.isPromise(passed)) {
        passed.then(calculateSuggestions).catch(() => {
        });
      } else if (passed !== false) {
        calculateSuggestions();
      } else {
        hideSuggestionPanel();
      }
    }, debounce);
    const handleInput = (val, e) => {
      !popperVisible.value && togglePopperVisible(true);
      if (e == null ? void 0 : e.isComposing)
        return;
      val ? handleFilter() : hideSuggestionPanel();
    };
    const getInputInnerHeight = (inputInner) => Number.parseFloat(core.useCssVar(nsInput.cssVarName("input-height"), inputInner).value) - 2;
    const focus = () => {
      var _a;
      (_a = inputRef.value) == null ? void 0 : _a.focus();
    };
    const blur = () => {
      var _a;
      (_a = inputRef.value) == null ? void 0 : _a.blur();
    };
    vue.watch(filtering, updatePopperPosition);
    vue.watch([
      checkedNodes,
      isDisabled,
      () => props.collapseTags,
      () => props.maxCollapseTags
    ], calculatePresentTags);
    vue.watch(tags, () => {
      vue.nextTick(() => updateStyle());
    });
    vue.watch(realSize, async () => {
      await vue.nextTick();
      const inputInner = inputRef.value.input;
      inputInitialHeight = getInputInnerHeight(inputInner) || inputInitialHeight;
      updateStyle();
    });
    vue.watch(presentText, syncPresentTextValue, { immediate: true });
    vue.watch(() => popperVisible.value, (val) => {
      var _a;
      if (val && props.props.lazy && props.props.lazyLoad) {
        (_a = cascaderPanelRef.value) == null ? void 0 : _a.loadLazyRootNodes();
      }
    });
    vue.onMounted(() => {
      const inputInner = inputRef.value.input;
      const inputInnerHeight = getInputInnerHeight(inputInner);
      inputInitialHeight = inputInner.offsetHeight || inputInnerHeight;
      core.useResizeObserver(inputInner, updateStyle);
    });
    expose({
      getCheckedNodes,
      cascaderPanelRef,
      togglePopperVisible,
      contentRef,
      presentText,
      focus,
      blur
    });
    return (_ctx, _cache) => {
      return vue.openBlock(), vue.createBlock(vue.unref(index$5.ElTooltip), {
        ref_key: "tooltipRef",
        ref: tooltipRef,
        visible: popperVisible.value,
        teleported: _ctx.teleported,
        "popper-class": [vue.unref(nsCascader).e("dropdown"), _ctx.popperClass],
        "popper-style": _ctx.popperStyle,
        "popper-options": popperOptions,
        "fallback-placements": _ctx.fallbackPlacements,
        "stop-popper-mouse-event": false,
        "gpu-acceleration": false,
        placement: _ctx.placement,
        transition: `${vue.unref(nsCascader).namespace.value}-zoom-in-top`,
        effect: _ctx.effect,
        pure: "",
        persistent: _ctx.persistent,
        onHide: hideSuggestionPanel
      }, {
        default: vue.withCtx(() => [
          vue.withDirectives((vue.openBlock(), vue.createElementBlock("div", {
            ref_key: "wrapperRef",
            ref: wrapperRef,
            class: vue.normalizeClass(vue.unref(cascaderKls)),
            style: vue.normalizeStyle(vue.unref(cascaderStyle)),
            onClick: () => togglePopperVisible(vue.unref(readonly) ? void 0 : true),
            onKeydown: handleKeyDown,
            onMouseenter: ($event) => inputHover.value = true,
            onMouseleave: ($event) => inputHover.value = false
          }, [
            vue.createVNode(vue.unref(index$6.ElInput), {
              ref_key: "inputRef",
              ref: inputRef,
              modelValue: inputValue.value,
              "onUpdate:modelValue": ($event) => inputValue.value = $event,
              placeholder: vue.unref(currentPlaceholder),
              readonly: vue.unref(readonly),
              disabled: vue.unref(isDisabled),
              "validate-event": false,
              size: vue.unref(realSize),
              class: vue.normalizeClass(vue.unref(inputClass)),
              tabindex: vue.unref(multiple) && _ctx.filterable && !vue.unref(isDisabled) ? -1 : void 0,
              onCompositionstart: vue.unref(handleComposition),
              onCompositionupdate: vue.unref(handleComposition),
              onCompositionend: vue.unref(handleComposition),
              onInput: handleInput
            }, vue.createSlots({
              suffix: vue.withCtx(() => [
                vue.unref(clearBtnVisible) ? (vue.openBlock(), vue.createBlock(vue.unref(index$7.ElIcon), {
                  key: "clear",
                  class: vue.normalizeClass([vue.unref(nsInput).e("icon"), "icon-circle-close"]),
                  onClick: vue.withModifiers(handleClear, ["stop"])
                }, {
                  default: vue.withCtx(() => [
                    (vue.openBlock(), vue.createBlock(vue.resolveDynamicComponent(_ctx.clearIcon)))
                  ]),
                  _: 1
                }, 8, ["class", "onClick"])) : (vue.openBlock(), vue.createBlock(vue.unref(index$7.ElIcon), {
                  key: "arrow-down",
                  class: vue.normalizeClass(vue.unref(cascaderIconKls)),
                  onClick: vue.withModifiers(($event) => togglePopperVisible(), ["stop"])
                }, {
                  default: vue.withCtx(() => [
                    vue.createVNode(vue.unref(iconsVue.ArrowDown))
                  ]),
                  _: 1
                }, 8, ["class", "onClick"]))
              ]),
              _: 2
            }, [
              _ctx.$slots.prefix ? {
                name: "prefix",
                fn: vue.withCtx(() => [
                  vue.renderSlot(_ctx.$slots, "prefix")
                ])
              } : void 0
            ]), 1032, ["modelValue", "onUpdate:modelValue", "placeholder", "readonly", "disabled", "size", "class", "tabindex", "onCompositionstart", "onCompositionupdate", "onCompositionend"]),
            vue.unref(multiple) ? (vue.openBlock(), vue.createElementBlock("div", {
              key: 0,
              ref_key: "tagWrapper",
              ref: tagWrapper,
              class: vue.normalizeClass([
                vue.unref(nsCascader).e("tags"),
                vue.unref(nsCascader).is("validate", Boolean(vue.unref(validateState)))
              ])
            }, [
              vue.renderSlot(_ctx.$slots, "tag", {
                data: tags.value,
                deleteTag
              }, () => [
                (vue.openBlock(true), vue.createElementBlock(vue.Fragment, null, vue.renderList(vue.unref(showTagList), (tag) => {
                  return vue.openBlock(), vue.createBlock(vue.unref(index$8.ElTag), {
                    key: tag.key,
                    type: _ctx.tagType,
                    size: vue.unref(tagSize),
                    effect: _ctx.tagEffect,
                    hit: tag.hitState,
                    closable: tag.closable,
                    "disable-transitions": "",
                    onClose: ($event) => deleteTag(tag)
                  }, {
                    default: vue.withCtx(() => [
                      vue.createElementVNode("span", null, vue.toDisplayString(tag.text), 1)
                    ]),
                    _: 2
                  }, 1032, ["type", "size", "effect", "hit", "closable", "onClose"]);
                }), 128))
              ]),
              _ctx.collapseTags && tags.value.length > _ctx.maxCollapseTags ? (vue.openBlock(), vue.createBlock(vue.unref(index$5.ElTooltip), {
                key: 0,
                ref_key: "tagTooltipRef",
                ref: tagTooltipRef,
                disabled: popperVisible.value || !_ctx.collapseTagsTooltip,
                "fallback-placements": ["bottom", "top", "right", "left"],
                placement: "bottom",
                "popper-class": _ctx.popperClass,
                "popper-style": _ctx.popperStyle,
                effect: _ctx.effect,
                persistent: _ctx.persistent
              }, {
                default: vue.withCtx(() => [
                  vue.createVNode(vue.unref(index$8.ElTag), {
                    closable: false,
                    size: vue.unref(tagSize),
                    type: _ctx.tagType,
                    effect: _ctx.tagEffect,
                    "disable-transitions": ""
                  }, {
                    default: vue.withCtx(() => [
                      vue.createElementVNode("span", {
                        class: vue.normalizeClass(vue.unref(nsCascader).e("tags-text"))
                      }, " + " + vue.toDisplayString(tags.value.length - _ctx.maxCollapseTags), 3)
                    ]),
                    _: 1
                  }, 8, ["size", "type", "effect"])
                ]),
                content: vue.withCtx(() => [
                  vue.createVNode(vue.unref(index$9.ElScrollbar), { "max-height": _ctx.maxCollapseTagsTooltipHeight }, {
                    default: vue.withCtx(() => [
                      vue.createElementVNode("div", {
                        class: vue.normalizeClass(vue.unref(nsCascader).e("collapse-tags"))
                      }, [
                        (vue.openBlock(true), vue.createElementBlock(vue.Fragment, null, vue.renderList(vue.unref(collapseTagList), (tag, idx) => {
                          return vue.openBlock(), vue.createElementBlock("div", {
                            key: idx,
                            class: vue.normalizeClass(vue.unref(nsCascader).e("collapse-tag"))
                          }, [
                            (vue.openBlock(), vue.createBlock(vue.unref(index$8.ElTag), {
                              key: tag.key,
                              class: "in-tooltip",
                              type: _ctx.tagType,
                              size: vue.unref(tagSize),
                              effect: _ctx.tagEffect,
                              hit: tag.hitState,
                              closable: tag.closable,
                              "disable-transitions": "",
                              onClose: ($event) => deleteTag(tag)
                            }, {
                              default: vue.withCtx(() => [
                                vue.createElementVNode("span", null, vue.toDisplayString(tag.text), 1)
                              ]),
                              _: 2
                            }, 1032, ["type", "size", "effect", "hit", "closable", "onClose"]))
                          ], 2);
                        }), 128))
                      ], 2)
                    ]),
                    _: 1
                  }, 8, ["max-height"])
                ]),
                _: 1
              }, 8, ["disabled", "popper-class", "popper-style", "effect", "persistent"])) : vue.createCommentVNode("v-if", true),
              _ctx.filterable && !vue.unref(isDisabled) ? vue.withDirectives((vue.openBlock(), vue.createElementBlock("input", {
                key: 1,
                "onUpdate:modelValue": ($event) => searchInputValue.value = $event,
                type: "text",
                class: vue.normalizeClass(vue.unref(nsCascader).e("search-input")),
                placeholder: vue.unref(presentText) ? "" : vue.unref(inputPlaceholder),
                onInput: (e) => handleInput(searchInputValue.value, e),
                onClick: vue.withModifiers(($event) => togglePopperVisible(true), ["stop"]),
                onKeydown: vue.withKeys(handleDelete, ["delete"]),
                onCompositionstart: vue.unref(handleComposition),
                onCompositionupdate: vue.unref(handleComposition),
                onCompositionend: vue.unref(handleComposition)
              }, null, 42, ["onUpdate:modelValue", "placeholder", "onInput", "onClick", "onKeydown", "onCompositionstart", "onCompositionupdate", "onCompositionend"])), [
                [vue.vModelText, searchInputValue.value]
              ]) : vue.createCommentVNode("v-if", true)
            ], 2)) : vue.createCommentVNode("v-if", true)
          ], 46, ["onClick", "onMouseenter", "onMouseleave"])), [
            [vue.unref(index$a["default"]), handleClickOutside, vue.unref(contentRef)]
          ])
        ]),
        content: vue.withCtx(() => [
          _ctx.$slots.header ? (vue.openBlock(), vue.createElementBlock("div", {
            key: 0,
            class: vue.normalizeClass(vue.unref(nsCascader).e("header")),
            onClick: vue.withModifiers(() => {
            }, ["stop"])
          }, [
            vue.renderSlot(_ctx.$slots, "header")
          ], 10, ["onClick"])) : vue.createCommentVNode("v-if", true),
          vue.withDirectives(vue.createVNode(vue.unref(index$b.ElCascaderPanel), {
            ref_key: "cascaderPanelRef",
            ref: cascaderPanelRef,
            modelValue: vue.unref(checkedValue),
            "onUpdate:modelValue": ($event) => vue.isRef(checkedValue) ? checkedValue.value = $event : null,
            options: _ctx.options,
            props: props.props,
            border: false,
            "render-label": _ctx.$slots.default,
            onExpandChange: handleExpandChange,
            onClose: ($event) => _ctx.$nextTick(() => togglePopperVisible(false))
          }, {
            empty: vue.withCtx(() => [
              vue.renderSlot(_ctx.$slots, "empty")
            ]),
            _: 3
          }, 8, ["modelValue", "onUpdate:modelValue", "options", "props", "render-label", "onClose"]), [
            [vue.vShow, !filtering.value]
          ]),
          _ctx.filterable ? vue.withDirectives((vue.openBlock(), vue.createBlock(vue.unref(index$9.ElScrollbar), {
            key: 1,
            ref_key: "suggestionPanel",
            ref: suggestionPanel,
            tag: "ul",
            class: vue.normalizeClass(vue.unref(nsCascader).e("suggestion-panel")),
            "view-class": vue.unref(nsCascader).e("suggestion-list"),
            onKeydown: handleSuggestionKeyDown
          }, {
            default: vue.withCtx(() => [
              suggestions.value.length ? (vue.openBlock(true), vue.createElementBlock(vue.Fragment, { key: 0 }, vue.renderList(suggestions.value, (item) => {
                return vue.openBlock(), vue.createElementBlock("li", {
                  key: item.uid,
                  class: vue.normalizeClass([
                    vue.unref(nsCascader).e("suggestion-item"),
                    vue.unref(nsCascader).is("checked", item.checked)
                  ]),
                  tabindex: -1,
                  onClick: ($event) => handleSuggestionClick(item)
                }, [
                  vue.renderSlot(_ctx.$slots, "suggestion-item", { item }, () => [
                    vue.createElementVNode("span", null, vue.toDisplayString(item.text), 1),
                    item.checked ? (vue.openBlock(), vue.createBlock(vue.unref(index$7.ElIcon), { key: 0 }, {
                      default: vue.withCtx(() => [
                        vue.createVNode(vue.unref(iconsVue.Check))
                      ]),
                      _: 1
                    })) : vue.createCommentVNode("v-if", true)
                  ])
                ], 10, ["onClick"]);
              }), 128)) : vue.renderSlot(_ctx.$slots, "empty", { key: 1 }, () => [
                vue.createElementVNode("li", {
                  class: vue.normalizeClass(vue.unref(nsCascader).e("empty-text"))
                }, vue.toDisplayString(vue.unref(t)("el.cascader.noMatch")), 3)
              ])
            ]),
            _: 3
          }, 8, ["class", "view-class"])), [
            [vue.vShow, filtering.value]
          ]) : vue.createCommentVNode("v-if", true),
          _ctx.$slots.footer ? (vue.openBlock(), vue.createElementBlock("div", {
            key: 2,
            class: vue.normalizeClass(vue.unref(nsCascader).e("footer")),
            onClick: vue.withModifiers(() => {
            }, ["stop"])
          }, [
            vue.renderSlot(_ctx.$slots, "footer")
          ], 10, ["onClick"])) : vue.createCommentVNode("v-if", true)
        ]),
        _: 3
      }, 8, ["visible", "teleported", "popper-class", "popper-style", "fallback-placements", "placement", "transition", "effect", "persistent"]);
    };
  }
});
var Cascader = /* @__PURE__ */ pluginVue_exportHelper["default"](_sfc_main, [["__file", "cascader.vue"]]);

exports["default"] = Cascader;
//# sourceMappingURL=cascader2.js.map
