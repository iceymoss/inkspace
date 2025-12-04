import { defineComponent, openBlock, createElementBlock, normalizeClass, unref, createElementVNode, Fragment, renderList, normalizeStyle } from 'vue';
import { predefineProps } from '../props/predefine.mjs';
import { usePredefine, usePredefineDOM } from '../composables/use-predefine.mjs';
import { useLocale } from '../../../../hooks/use-locale/index.mjs';
import _export_sfc from '../../../../_virtual/plugin-vue_export-helper.mjs';

const __default__ = defineComponent({
  name: "ElColorPredefine"
});
const _sfc_main = /* @__PURE__ */ defineComponent({
  ...__default__,
  props: predefineProps,
  setup(__props) {
    const props = __props;
    const { rgbaColors, handleSelect } = usePredefine(props);
    const { rootKls, colorsKls, colorSelectorKls } = usePredefineDOM(props);
    const { t } = useLocale();
    const ariaLabel = (value) => {
      return t("el.colorpicker.predefineDescription", { value });
    };
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", {
        class: normalizeClass(unref(rootKls))
      }, [
        createElementVNode("div", {
          class: normalizeClass(unref(colorsKls))
        }, [
          (openBlock(true), createElementBlock(Fragment, null, renderList(unref(rgbaColors), (item, index) => {
            return openBlock(), createElementBlock("button", {
              key: _ctx.colors[index],
              type: "button",
              disabled: _ctx.disabled,
              "aria-label": ariaLabel(item.value),
              class: normalizeClass(unref(colorSelectorKls)(item)),
              onClick: ($event) => unref(handleSelect)(index)
            }, [
              createElementVNode("div", {
                style: normalizeStyle({ backgroundColor: item.value })
              }, null, 4)
            ], 10, ["disabled", "aria-label", "onClick"]);
          }), 128))
        ], 2)
      ], 2);
    };
  }
});
var Predefine = /* @__PURE__ */ _export_sfc(_sfc_main, [["__file", "predefine.vue"]]);

export { Predefine as default };
//# sourceMappingURL=predefine.mjs.map
