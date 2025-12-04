'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var vue = require('vue');
var index$1 = require('../../input/index.js');
var alphaSlider = require('./components/alpha-slider.js');
var hueSlider = require('./components/hue-slider.js');
var predefine = require('./components/predefine.js');
var svPanel = require('./components/sv-panel.js');
var colorPickerPanel = require('./color-picker-panel.js');
var useCommonColor = require('./composables/use-common-color.js');
var pluginVue_exportHelper = require('../../../_virtual/plugin-vue_export-helper.js');
var index = require('../../../hooks/use-namespace/index.js');
var useFormItem = require('../../form/src/hooks/use-form-item.js');
var useFormCommonProps = require('../../form/src/hooks/use-form-common-props.js');
var event = require('../../../constants/event.js');
var error = require('../../../utils/error.js');

const __default__ = vue.defineComponent({
  name: "ElColorPickerPanel"
});
const _sfc_main = /* @__PURE__ */ vue.defineComponent({
  ...__default__,
  props: colorPickerPanel.colorPickerPanelProps,
  emits: colorPickerPanel.colorPickerPanelEmits,
  setup(__props, { expose, emit }) {
    const props = __props;
    const ns = index.useNamespace("color-picker-panel");
    const { formItem } = useFormItem.useFormItem();
    const disabled = useFormCommonProps.useFormDisabled();
    const hueRef = vue.ref();
    const svRef = vue.ref();
    const alphaRef = vue.ref();
    const inputRef = vue.ref();
    const customInput = vue.ref("");
    const { color } = vue.inject(colorPickerPanel.ROOT_COMMON_COLOR_INJECTION_KEY, () => useCommonColor.useCommonColor(props, emit), true);
    function handleConfirm() {
      color.fromString(customInput.value);
      if (color.value !== customInput.value) {
        customInput.value = color.value;
      }
    }
    function handleFocusout() {
      var _a;
      if (props.validateEvent) {
        (_a = formItem == null ? void 0 : formItem.validate) == null ? void 0 : _a.call(formItem, "blur").catch((err) => error.debugWarn(err));
      }
    }
    function update() {
      var _a, _b, _c;
      (_a = hueRef.value) == null ? void 0 : _a.update();
      (_b = svRef.value) == null ? void 0 : _b.update();
      (_c = alphaRef.value) == null ? void 0 : _c.update();
    }
    vue.onMounted(() => {
      if (props.modelValue) {
        customInput.value = color.value;
      }
      vue.nextTick(update);
    });
    vue.watch(() => props.modelValue, (newVal) => {
      if (newVal !== color.value) {
        newVal ? color.fromString(newVal) : color.clear();
      }
    });
    vue.watch(() => color.value, (val) => {
      emit(event.UPDATE_MODEL_EVENT, val);
      customInput.value = val;
      if (props.validateEvent) {
        formItem == null ? void 0 : formItem.validate("change").catch((err) => error.debugWarn(err));
      }
    });
    vue.provide(colorPickerPanel.colorPickerPanelContextKey, {
      currentColor: vue.computed(() => color.value)
    });
    expose({
      color,
      inputRef,
      update
    });
    return (_ctx, _cache) => {
      return vue.openBlock(), vue.createElementBlock("div", {
        class: vue.normalizeClass([vue.unref(ns).b(), vue.unref(ns).is("disabled", vue.unref(disabled)), vue.unref(ns).is("border", _ctx.border)]),
        onFocusout: handleFocusout
      }, [
        vue.createElementVNode("div", {
          class: vue.normalizeClass(vue.unref(ns).e("wrapper"))
        }, [
          vue.createVNode(hueSlider["default"], {
            ref_key: "hueRef",
            ref: hueRef,
            class: "hue-slider",
            color: vue.unref(color),
            vertical: "",
            disabled: vue.unref(disabled)
          }, null, 8, ["color", "disabled"]),
          vue.createVNode(svPanel["default"], {
            ref_key: "svRef",
            ref: svRef,
            color: vue.unref(color),
            disabled: vue.unref(disabled)
          }, null, 8, ["color", "disabled"])
        ], 2),
        _ctx.showAlpha ? (vue.openBlock(), vue.createBlock(alphaSlider["default"], {
          key: 0,
          ref_key: "alphaRef",
          ref: alphaRef,
          color: vue.unref(color),
          disabled: vue.unref(disabled)
        }, null, 8, ["color", "disabled"])) : vue.createCommentVNode("v-if", true),
        _ctx.predefine ? (vue.openBlock(), vue.createBlock(predefine["default"], {
          key: 1,
          ref: "predefine",
          "enable-alpha": _ctx.showAlpha,
          color: vue.unref(color),
          colors: _ctx.predefine,
          disabled: vue.unref(disabled)
        }, null, 8, ["enable-alpha", "color", "colors", "disabled"])) : vue.createCommentVNode("v-if", true),
        vue.createElementVNode("div", {
          class: vue.normalizeClass(vue.unref(ns).e("footer"))
        }, [
          vue.createVNode(vue.unref(index$1.ElInput), {
            ref_key: "inputRef",
            ref: inputRef,
            modelValue: customInput.value,
            "onUpdate:modelValue": ($event) => customInput.value = $event,
            "validate-event": false,
            size: "small",
            disabled: vue.unref(disabled),
            onChange: handleConfirm
          }, null, 8, ["modelValue", "onUpdate:modelValue", "disabled"]),
          vue.renderSlot(_ctx.$slots, "footer")
        ], 2)
      ], 34);
    };
  }
});
var ColorPickerPanel = /* @__PURE__ */ pluginVue_exportHelper["default"](_sfc_main, [["__file", "color-picker-panel.vue"]]);

exports["default"] = ColorPickerPanel;
//# sourceMappingURL=color-picker-panel2.js.map
