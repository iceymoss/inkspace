'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var runtime = require('../../../utils/vue/props/runtime.js');
var typescript = require('../../../utils/typescript.js');
var types = require('../../../utils/types.js');

const imageViewerProps = runtime.buildProps({
  urlList: {
    type: runtime.definePropType(Array),
    default: () => typescript.mutable([])
  },
  zIndex: {
    type: Number
  },
  initialIndex: {
    type: Number,
    default: 0
  },
  infinite: {
    type: Boolean,
    default: true
  },
  hideOnClickModal: Boolean,
  teleported: Boolean,
  closeOnPressEscape: {
    type: Boolean,
    default: true
  },
  zoomRate: {
    type: Number,
    default: 1.2
  },
  scale: {
    type: Number,
    default: 1
  },
  minScale: {
    type: Number,
    default: 0.2
  },
  maxScale: {
    type: Number,
    default: 7
  },
  showProgress: Boolean,
  crossorigin: {
    type: runtime.definePropType(String)
  }
});
const imageViewerEmits = {
  close: () => true,
  error: (evt) => evt instanceof Event,
  switch: (index) => types.isNumber(index),
  rotate: (deg) => types.isNumber(deg)
};

exports.imageViewerEmits = imageViewerEmits;
exports.imageViewerProps = imageViewerProps;
//# sourceMappingURL=image-viewer.js.map
