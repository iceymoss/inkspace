'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var vue = require('vue');
var predefine = require('../props/predefine.js');
var usePredefine = require('../composables/use-predefine.js');
var index = require('../../../../hooks/use-locale/index.js');
var pluginVue_exportHelper = require('../../../../_virtual/plugin-vue_export-helper.js');

const __default__ = vue.defineComponent({
  name: "ElColorPredefine"
});
const _sfc_main = /* @__PURE__ */ vue.defineComponent({
  ...__default__,
  props: predefine.predefineProps,
  setup(__props) {
    const props = __props;
    const { rgbaColors, handleSelect } = usePredefine.usePredefine(props);
    const { rootKls, colorsKls, colorSelectorKls } = usePredefine.usePredefineDOM(props);
    const { t } = index.useLocale();
    const ariaLabel = (value) => {
      return t("el.colorpicker.predefineDescription", { value });
    };
    return (_ctx, _cache) => {
      return vue.openBlock(), vue.createElementBlock("div", {
        class: vue.normalizeClass(vue.unref(rootKls))
      }, [
        vue.createElementVNode("div", {
          class: vue.normalizeClass(vue.unref(colorsKls))
        }, [
          (vue.openBlock(true), vue.createElementBlock(vue.Fragment, null, vue.renderList(vue.unref(rgbaColors), (item, index) => {
            return vue.openBlock(), vue.createElementBlock("button", {
              key: _ctx.colors[index],
              type: "button",
              disabled: _ctx.disabled,
              "aria-label": ariaLabel(item.value),
              class: vue.normalizeClass(vue.unref(colorSelectorKls)(item)),
              onClick: ($event) => vue.unref(handleSelect)(index)
            }, [
              vue.createElementVNode("div", {
                style: vue.normalizeStyle({ backgroundColor: item.value })
              }, null, 4)
            ], 10, ["disabled", "aria-label", "onClick"]);
          }), 128))
        ], 2)
      ], 2);
    };
  }
});
var Predefine = /* @__PURE__ */ pluginVue_exportHelper["default"](_sfc_main, [["__file", "predefine.vue"]]);

exports["default"] = Predefine;
//# sourceMappingURL=predefine.js.map
