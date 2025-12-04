'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var vue = require('vue');
var core = require('@vueuse/core');
var constants = require('./constants.js');
var tabPane = require('./tab-pane2.js');
var pluginVue_exportHelper = require('../../../_virtual/plugin-vue_export-helper.js');
var error = require('../../../utils/error.js');
var index = require('../../../hooks/use-namespace/index.js');

const COMPONENT_NAME = "ElTabPane";
const __default__ = vue.defineComponent({
  name: COMPONENT_NAME
});
const _sfc_main = /* @__PURE__ */ vue.defineComponent({
  ...__default__,
  props: tabPane.tabPaneProps,
  setup(__props) {
    const props = __props;
    const instance = vue.getCurrentInstance();
    const slots = vue.useSlots();
    const tabsRoot = vue.inject(constants.tabsRootContextKey);
    if (!tabsRoot)
      error.throwError(COMPONENT_NAME, "usage: <el-tabs><el-tab-pane /></el-tabs/>");
    const ns = index.useNamespace("tab-pane");
    const paneRef = vue.ref();
    const index$1 = vue.ref();
    const isClosable = vue.computed(() => {
      var _a;
      return (_a = props.closable) != null ? _a : tabsRoot.props.closable;
    });
    const active = core.eagerComputed(() => {
      var _a;
      return tabsRoot.currentName.value === ((_a = props.name) != null ? _a : index$1.value);
    });
    const loaded = vue.ref(active.value);
    const paneName = vue.computed(() => {
      var _a;
      return (_a = props.name) != null ? _a : index$1.value;
    });
    const shouldBeRender = core.eagerComputed(() => !props.lazy || loaded.value || active.value);
    const isFocusInsidePane = () => {
      var _a;
      return (_a = paneRef.value) == null ? void 0 : _a.contains(document.activeElement);
    };
    vue.watch(active, (val) => {
      if (val)
        loaded.value = true;
    });
    const pane = vue.reactive({
      uid: instance.uid,
      getVnode: () => instance.vnode,
      slots,
      props,
      paneName,
      active,
      index: index$1,
      isClosable,
      isFocusInsidePane
    });
    tabsRoot.registerPane(pane);
    vue.onBeforeUnmount(() => {
      tabsRoot.unregisterPane(pane);
    });
    vue.onBeforeUpdate(() => {
      var _a;
      if (slots.label)
        (_a = tabsRoot.nav$.value) == null ? void 0 : _a.scheduleRender();
    });
    return (_ctx, _cache) => {
      return vue.unref(shouldBeRender) ? vue.withDirectives((vue.openBlock(), vue.createElementBlock("div", {
        key: 0,
        id: `pane-${vue.unref(paneName)}`,
        ref_key: "paneRef",
        ref: paneRef,
        class: vue.normalizeClass(vue.unref(ns).b()),
        role: "tabpanel",
        "aria-hidden": !vue.unref(active),
        "aria-labelledby": `tab-${vue.unref(paneName)}`
      }, [
        vue.renderSlot(_ctx.$slots, "default")
      ], 10, ["id", "aria-hidden", "aria-labelledby"])), [
        [vue.vShow, vue.unref(active)]
      ]) : vue.createCommentVNode("v-if", true);
    };
  }
});
var TabPane = /* @__PURE__ */ pluginVue_exportHelper["default"](_sfc_main, [["__file", "tab-pane.vue"]]);

exports["default"] = TabPane;
//# sourceMappingURL=tab-pane.js.map
