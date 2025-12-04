'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var vue = require('vue');
var treeStore = require('./model/tree-store.js');
var util = require('./model/util.js');
var treeNode = require('./tree-node.js');
var useNodeExpandEventBroadcast = require('./model/useNodeExpandEventBroadcast.js');
var useDragNode = require('./model/useDragNode.js');
var useKeydown = require('./model/useKeydown.js');
var tokens = require('./tokens.js');
var lodashUnified = require('lodash-unified');
var tree = require('./tree.js');
var pluginVue_exportHelper = require('../../../_virtual/plugin-vue_export-helper.js');
var runtime = require('../../../utils/vue/props/runtime.js');
var icon = require('../../../utils/vue/icon.js');
var index = require('../../../hooks/use-locale/index.js');
var index$1 = require('../../../hooks/use-namespace/index.js');
var constants = require('../../form/src/constants.js');

const _sfc_main = vue.defineComponent({
  name: "ElTree",
  components: { ElTreeNode: treeNode["default"] },
  props: {
    data: {
      type: runtime.definePropType(Array),
      default: () => []
    },
    emptyText: {
      type: String
    },
    renderAfterExpand: {
      type: Boolean,
      default: true
    },
    nodeKey: String,
    checkStrictly: Boolean,
    defaultExpandAll: Boolean,
    expandOnClickNode: {
      type: Boolean,
      default: true
    },
    checkOnClickNode: Boolean,
    checkOnClickLeaf: {
      type: Boolean,
      default: true
    },
    checkDescendants: Boolean,
    autoExpandParent: {
      type: Boolean,
      default: true
    },
    defaultCheckedKeys: Array,
    defaultExpandedKeys: Array,
    currentNodeKey: [String, Number],
    renderContent: {
      type: runtime.definePropType(Function)
    },
    showCheckbox: Boolean,
    draggable: Boolean,
    allowDrag: {
      type: runtime.definePropType(Function)
    },
    allowDrop: {
      type: runtime.definePropType(Function)
    },
    props: {
      type: Object,
      default: () => ({
        children: "children",
        label: "label",
        disabled: "disabled"
      })
    },
    lazy: Boolean,
    highlightCurrent: Boolean,
    load: Function,
    filterNodeMethod: Function,
    accordion: Boolean,
    indent: {
      type: Number,
      default: 18
    },
    icon: {
      type: icon.iconPropType
    }
  },
  emits: tree.treeEmits,
  setup(props, ctx) {
    const { t } = index.useLocale();
    const ns = index$1.useNamespace("tree");
    const store = vue.ref(new treeStore["default"]({
      key: props.nodeKey,
      data: props.data,
      lazy: props.lazy,
      props: props.props,
      load: props.load,
      currentNodeKey: props.currentNodeKey,
      checkStrictly: props.checkStrictly,
      checkDescendants: props.checkDescendants,
      defaultCheckedKeys: props.defaultCheckedKeys,
      defaultExpandedKeys: props.defaultExpandedKeys,
      autoExpandParent: props.autoExpandParent,
      defaultExpandAll: props.defaultExpandAll,
      filterNodeMethod: props.filterNodeMethod
    }));
    store.value.initialize();
    const root = vue.ref(store.value.root);
    const currentNode = vue.ref(null);
    const el$ = vue.ref(null);
    const dropIndicator$ = vue.ref(null);
    const { broadcastExpanded } = useNodeExpandEventBroadcast.useNodeExpandEventBroadcast(props);
    const { dragState } = useDragNode.useDragNodeHandler({
      props,
      ctx,
      el$,
      dropIndicator$,
      store
    });
    useKeydown.useKeydown({ el$ }, store);
    const instance = vue.getCurrentInstance();
    const isSelectTree = vue.computed(() => {
      let parent = instance == null ? void 0 : instance.parent;
      while (parent) {
        if (parent.type.name === "ElTreeSelect") {
          return true;
        }
        parent = parent.parent;
      }
      return false;
    });
    const isEmpty = vue.computed(() => {
      const { childNodes } = root.value;
      return (!childNodes || childNodes.length === 0 || childNodes.every(({ visible }) => !visible)) && !isSelectTree.value;
    });
    vue.watch(() => props.currentNodeKey, (newVal) => {
      store.value.setCurrentNodeKey(newVal != null ? newVal : null);
    });
    vue.watch(() => props.defaultCheckedKeys, (newVal, oldVal) => {
      if (lodashUnified.isEqual(newVal, oldVal))
        return;
      store.value.setDefaultCheckedKey(newVal != null ? newVal : []);
    });
    vue.watch(() => props.defaultExpandedKeys, (newVal) => {
      store.value.setDefaultExpandedKeys(newVal != null ? newVal : []);
    });
    vue.watch(() => props.data, (newVal) => {
      store.value.setData(newVal);
    }, { deep: true });
    vue.watch(() => props.checkStrictly, (newVal) => {
      store.value.checkStrictly = newVal;
    });
    const filter = (value) => {
      if (!props.filterNodeMethod)
        throw new Error("[Tree] filterNodeMethod is required when filter");
      store.value.filter(value);
    };
    const getNodeKey = (node) => {
      return util.getNodeKey(props.nodeKey, node.data);
    };
    const requireNodeKey = (methodName) => {
      if (!props.nodeKey) {
        throw new Error(`[Tree] nodeKey is required in ${methodName}`);
      }
    };
    const getNodePath = (data) => {
      requireNodeKey("getNodePath");
      const node = store.value.getNode(data);
      if (!node)
        return [];
      const path = [node.data];
      let parent = node.parent;
      while (parent && parent !== root.value) {
        path.push(parent.data);
        parent = parent.parent;
      }
      return path.reverse();
    };
    const getCheckedNodes = (leafOnly, includeHalfChecked) => {
      return store.value.getCheckedNodes(leafOnly, includeHalfChecked);
    };
    const getCheckedKeys = (leafOnly) => {
      return store.value.getCheckedKeys(leafOnly);
    };
    const getCurrentNode = () => {
      const currentNode2 = store.value.getCurrentNode();
      return currentNode2 ? currentNode2.data : null;
    };
    const getCurrentKey = () => {
      requireNodeKey("getCurrentKey");
      const currentNode2 = getCurrentNode();
      return currentNode2 ? currentNode2[props.nodeKey] : null;
    };
    const setCheckedNodes = (nodes, leafOnly) => {
      requireNodeKey("setCheckedNodes");
      store.value.setCheckedNodes(nodes, leafOnly);
    };
    const setCheckedKeys = (keys, leafOnly) => {
      requireNodeKey("setCheckedKeys");
      store.value.setCheckedKeys(keys, leafOnly);
    };
    const setChecked = (data, checked, deep) => {
      store.value.setChecked(data, checked, deep);
    };
    const getHalfCheckedNodes = () => {
      return store.value.getHalfCheckedNodes();
    };
    const getHalfCheckedKeys = () => {
      return store.value.getHalfCheckedKeys();
    };
    const setCurrentNode = (node, shouldAutoExpandParent = true) => {
      requireNodeKey("setCurrentNode");
      util.handleCurrentChange(store, ctx.emit, () => {
        broadcastExpanded(node);
        store.value.setUserCurrentNode(node, shouldAutoExpandParent);
      });
    };
    const setCurrentKey = (key = null, shouldAutoExpandParent = true) => {
      requireNodeKey("setCurrentKey");
      util.handleCurrentChange(store, ctx.emit, () => {
        broadcastExpanded();
        store.value.setCurrentNodeKey(key, shouldAutoExpandParent);
      });
    };
    const getNode = (data) => {
      return store.value.getNode(data);
    };
    const remove = (data) => {
      store.value.remove(data);
    };
    const append = (data, parentNode) => {
      store.value.append(data, parentNode);
    };
    const insertBefore = (data, refNode) => {
      store.value.insertBefore(data, refNode);
    };
    const insertAfter = (data, refNode) => {
      store.value.insertAfter(data, refNode);
    };
    const handleNodeExpand = (nodeData, node, instance2) => {
      broadcastExpanded(node);
      ctx.emit("node-expand", nodeData, node, instance2);
    };
    const updateKeyChildren = (key, data) => {
      requireNodeKey("updateKeyChild");
      store.value.updateChildren(key, data);
    };
    vue.provide(tokens.ROOT_TREE_INJECTION_KEY, {
      ctx,
      props,
      store,
      root,
      currentNode,
      instance
    });
    vue.provide(constants.formItemContextKey, void 0);
    return {
      ns,
      store,
      root,
      currentNode,
      dragState,
      el$,
      dropIndicator$,
      isEmpty,
      filter,
      getNodeKey,
      getNodePath,
      getCheckedNodes,
      getCheckedKeys,
      getCurrentNode,
      getCurrentKey,
      setCheckedNodes,
      setCheckedKeys,
      setChecked,
      getHalfCheckedNodes,
      getHalfCheckedKeys,
      setCurrentNode,
      setCurrentKey,
      t,
      getNode,
      remove,
      append,
      insertBefore,
      insertAfter,
      handleNodeExpand,
      updateKeyChildren
    };
  }
});
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  const _component_el_tree_node = vue.resolveComponent("el-tree-node");
  return vue.openBlock(), vue.createElementBlock("div", {
    ref: "el$",
    class: vue.normalizeClass([
      _ctx.ns.b(),
      _ctx.ns.is("dragging", !!_ctx.dragState.draggingNode),
      _ctx.ns.is("drop-not-allow", !_ctx.dragState.allowDrop),
      _ctx.ns.is("drop-inner", _ctx.dragState.dropType === "inner"),
      { [_ctx.ns.m("highlight-current")]: _ctx.highlightCurrent }
    ]),
    role: "tree"
  }, [
    (vue.openBlock(true), vue.createElementBlock(vue.Fragment, null, vue.renderList(_ctx.root.childNodes, (child) => {
      return vue.openBlock(), vue.createBlock(_component_el_tree_node, {
        key: _ctx.getNodeKey(child),
        node: child,
        props: _ctx.props,
        accordion: _ctx.accordion,
        "render-after-expand": _ctx.renderAfterExpand,
        "show-checkbox": _ctx.showCheckbox,
        "render-content": _ctx.renderContent,
        onNodeExpand: _ctx.handleNodeExpand
      }, null, 8, ["node", "props", "accordion", "render-after-expand", "show-checkbox", "render-content", "onNodeExpand"]);
    }), 128)),
    _ctx.isEmpty ? (vue.openBlock(), vue.createElementBlock("div", {
      key: 0,
      class: vue.normalizeClass(_ctx.ns.e("empty-block"))
    }, [
      vue.renderSlot(_ctx.$slots, "empty", {}, () => {
        var _a;
        return [
          vue.createElementVNode("span", {
            class: vue.normalizeClass(_ctx.ns.e("empty-text"))
          }, vue.toDisplayString((_a = _ctx.emptyText) != null ? _a : _ctx.t("el.tree.emptyText")), 3)
        ];
      })
    ], 2)) : vue.createCommentVNode("v-if", true),
    vue.withDirectives(vue.createElementVNode("div", {
      ref: "dropIndicator$",
      class: vue.normalizeClass(_ctx.ns.e("drop-indicator"))
    }, null, 2), [
      [vue.vShow, _ctx.dragState.showDropIndicator]
    ])
  ], 2);
}
var Tree = /* @__PURE__ */ pluginVue_exportHelper["default"](_sfc_main, [["render", _sfc_render], ["__file", "tree.vue"]]);

exports["default"] = Tree;
//# sourceMappingURL=tree2.js.map
