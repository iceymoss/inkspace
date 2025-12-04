'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var vue = require('vue');
var lodashUnified = require('lodash-unified');
var checkboxGroup = require('./checkbox-group.js');
var constants = require('./constants.js');
var checkbox = require('./checkbox2.js');
var checkboxButton = require('./checkbox-button.js');
var pluginVue_exportHelper = require('../../../_virtual/plugin-vue_export-helper.js');
var index = require('../../../hooks/use-namespace/index.js');
var useFormItem = require('../../form/src/hooks/use-form-item.js');
var error = require('../../../utils/error.js');
var event = require('../../../constants/event.js');

const __default__ = vue.defineComponent({
  name: "ElCheckboxGroup"
});
const _sfc_main = /* @__PURE__ */ vue.defineComponent({
  ...__default__,
  props: checkboxGroup.checkboxGroupProps,
  emits: checkboxGroup.checkboxGroupEmits,
  setup(__props, { emit }) {
    const props = __props;
    const ns = index.useNamespace("checkbox");
    const { formItem } = useFormItem.useFormItem();
    const { inputId: groupId, isLabeledByFormItem } = useFormItem.useFormItemInputId(props, {
      formItemContext: formItem
    });
    const changeEvent = async (value) => {
      emit(event.UPDATE_MODEL_EVENT, value);
      await vue.nextTick();
      emit(event.CHANGE_EVENT, value);
    };
    const modelValue = vue.computed({
      get() {
        return props.modelValue;
      },
      set(val) {
        changeEvent(val);
      }
    });
    const aliasProps = vue.computed(() => ({
      ...checkboxGroup.checkboxDefaultProps,
      ...props.props
    }));
    const getOptionProps = (option) => {
      const { label, value, disabled } = aliasProps.value;
      const base = {
        label: option[label],
        value: option[value],
        disabled: option[disabled]
      };
      return { ...lodashUnified.omit(option, [label, value, disabled]), ...base };
    };
    const optionComponent = vue.computed(() => props.type === "button" ? checkboxButton["default"] : checkbox["default"]);
    vue.provide(constants.checkboxGroupContextKey, {
      ...lodashUnified.pick(vue.toRefs(props), [
        "size",
        "min",
        "max",
        "disabled",
        "validateEvent",
        "fill",
        "textColor"
      ]),
      modelValue,
      changeEvent
    });
    vue.watch(() => props.modelValue, (newVal, oldValue) => {
      if (props.validateEvent && !lodashUnified.isEqual(newVal, oldValue)) {
        formItem == null ? void 0 : formItem.validate("change").catch((err) => error.debugWarn(err));
      }
    });
    return (_ctx, _cache) => {
      var _a;
      return vue.openBlock(), vue.createBlock(vue.resolveDynamicComponent(_ctx.tag), {
        id: vue.unref(groupId),
        class: vue.normalizeClass(vue.unref(ns).b("group")),
        role: "group",
        "aria-label": !vue.unref(isLabeledByFormItem) ? _ctx.ariaLabel || "checkbox-group" : void 0,
        "aria-labelledby": vue.unref(isLabeledByFormItem) ? (_a = vue.unref(formItem)) == null ? void 0 : _a.labelId : void 0
      }, {
        default: vue.withCtx(() => [
          vue.renderSlot(_ctx.$slots, "default", {}, () => [
            (vue.openBlock(true), vue.createElementBlock(vue.Fragment, null, vue.renderList(_ctx.options, (item, index) => {
              return vue.openBlock(), vue.createBlock(vue.resolveDynamicComponent(vue.unref(optionComponent)), vue.mergeProps({ key: index }, getOptionProps(item)), null, 16);
            }), 128))
          ])
        ]),
        _: 3
      }, 8, ["id", "class", "aria-label", "aria-labelledby"]);
    };
  }
});
var CheckboxGroup = /* @__PURE__ */ pluginVue_exportHelper["default"](_sfc_main, [["__file", "checkbox-group.vue"]]);

exports["default"] = CheckboxGroup;
//# sourceMappingURL=checkbox-group2.js.map
