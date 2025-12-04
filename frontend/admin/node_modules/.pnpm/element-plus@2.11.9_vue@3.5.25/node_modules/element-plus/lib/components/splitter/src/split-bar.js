'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var vue = require('vue');
var iconsVue = require('@element-plus/icons-vue');
var pluginVue_exportHelper = require('../../../_virtual/plugin-vue_export-helper.js');
var index = require('../../../hooks/use-namespace/index.js');

const __default__ = vue.defineComponent({
  name: "ElSplitterBar"
});
const _sfc_main = /* @__PURE__ */ vue.defineComponent({
  ...__default__,
  props: {
    index: {
      type: Number,
      required: true
    },
    layout: {
      type: String,
      values: ["horizontal", "vertical"],
      default: "horizontal"
    },
    resizable: {
      type: Boolean,
      default: true
    },
    lazy: Boolean,
    startCollapsible: Boolean,
    endCollapsible: Boolean
  },
  emits: ["moveStart", "moving", "moveEnd", "collapse"],
  setup(__props, { emit }) {
    const props = __props;
    const ns = index.useNamespace("splitter-bar");
    const isHorizontal = vue.computed(() => props.layout === "horizontal");
    const barWrapStyles = vue.computed(() => {
      if (isHorizontal.value) {
        return { width: 0 };
      }
      return { height: 0 };
    });
    const draggerStyles = vue.computed(() => {
      return {
        width: isHorizontal.value ? "16px" : "100%",
        height: isHorizontal.value ? "100%" : "16px",
        cursor: !props.resizable ? "auto" : isHorizontal.value ? "ew-resize" : "ns-resize",
        touchAction: "none"
      };
    });
    const draggerPseudoClass = vue.computed(() => {
      const prefix = ns.e("dragger");
      return {
        [`${prefix}-horizontal`]: isHorizontal.value,
        [`${prefix}-vertical`]: !isHorizontal.value,
        [`${prefix}-active`]: !!startPos.value
      };
    });
    const startPos = vue.ref(null);
    const onMousedown = (e) => {
      if (!props.resizable)
        return;
      startPos.value = [e.pageX, e.pageY];
      emit("moveStart", props.index);
      window.addEventListener("mouseup", onMouseUp);
      window.addEventListener("mousemove", onMouseMove);
    };
    const onTouchStart = (e) => {
      if (props.resizable && e.touches.length === 1) {
        e.preventDefault();
        const touch = e.touches[0];
        startPos.value = [touch.pageX, touch.pageY];
        emit("moveStart", props.index);
        window.addEventListener("touchend", onTouchEnd);
        window.addEventListener("touchmove", onTouchMove);
      }
    };
    const onMouseMove = (e) => {
      const { pageX, pageY } = e;
      const offsetX = pageX - startPos.value[0];
      const offsetY = pageY - startPos.value[1];
      const offset = isHorizontal.value ? offsetX : offsetY;
      emit("moving", props.index, offset);
    };
    const onTouchMove = (e) => {
      if (e.touches.length === 1) {
        e.preventDefault();
        const touch = e.touches[0];
        const offsetX = touch.pageX - startPos.value[0];
        const offsetY = touch.pageY - startPos.value[1];
        const offset = isHorizontal.value ? offsetX : offsetY;
        emit("moving", props.index, offset);
      }
    };
    const onMouseUp = () => {
      startPos.value = null;
      window.removeEventListener("mouseup", onMouseUp);
      window.removeEventListener("mousemove", onMouseMove);
      emit("moveEnd", props.index);
    };
    const onTouchEnd = () => {
      startPos.value = null;
      window.removeEventListener("touchend", onTouchEnd);
      window.removeEventListener("touchmove", onTouchMove);
      emit("moveEnd", props.index);
    };
    const StartIcon = vue.computed(() => isHorizontal.value ? iconsVue.ArrowLeft : iconsVue.ArrowUp);
    const EndIcon = vue.computed(() => isHorizontal.value ? iconsVue.ArrowRight : iconsVue.ArrowDown);
    return (_ctx, _cache) => {
      return vue.openBlock(), vue.createElementBlock("div", {
        class: vue.normalizeClass([vue.unref(ns).b()]),
        style: vue.normalizeStyle(vue.unref(barWrapStyles))
      }, [
        __props.startCollapsible ? (vue.openBlock(), vue.createElementBlock("div", {
          key: 0,
          class: vue.normalizeClass([vue.unref(ns).e("collapse-icon"), vue.unref(ns).e(`${__props.layout}-collapse-icon-start`)]),
          onClick: ($event) => emit("collapse", __props.index, "start")
        }, [
          vue.renderSlot(_ctx.$slots, "start-collapsible", {}, () => [
            (vue.openBlock(), vue.createBlock(vue.resolveDynamicComponent(vue.unref(StartIcon)), { style: { "width": "12px", "height": "12px" } }))
          ])
        ], 10, ["onClick"])) : vue.createCommentVNode("v-if", true),
        vue.createElementVNode("div", {
          class: vue.normalizeClass([
            vue.unref(ns).e("dragger"),
            vue.unref(draggerPseudoClass),
            vue.unref(ns).is("disabled", !__props.resizable),
            vue.unref(ns).is("lazy", __props.resizable && __props.lazy)
          ]),
          style: vue.normalizeStyle(vue.unref(draggerStyles)),
          onMousedown,
          onTouchstart: onTouchStart
        }, null, 38),
        __props.endCollapsible ? (vue.openBlock(), vue.createElementBlock("div", {
          key: 1,
          class: vue.normalizeClass([vue.unref(ns).e("collapse-icon"), vue.unref(ns).e(`${__props.layout}-collapse-icon-end`)]),
          onClick: ($event) => emit("collapse", __props.index, "end")
        }, [
          vue.renderSlot(_ctx.$slots, "end-collapsible", {}, () => [
            (vue.openBlock(), vue.createBlock(vue.resolveDynamicComponent(vue.unref(EndIcon)), { style: { "width": "12px", "height": "12px" } }))
          ])
        ], 10, ["onClick"])) : vue.createCommentVNode("v-if", true)
      ], 6);
    };
  }
});
var SplitBar = /* @__PURE__ */ pluginVue_exportHelper["default"](_sfc_main, [["__file", "split-bar.vue"]]);

exports["default"] = SplitBar;
//# sourceMappingURL=split-bar.js.map
