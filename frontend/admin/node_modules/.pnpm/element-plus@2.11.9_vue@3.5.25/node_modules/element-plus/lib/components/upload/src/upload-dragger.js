'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var runtime = require('../../../utils/vue/props/runtime.js');
var shared = require('@vue/shared');

const uploadDraggerProps = runtime.buildProps({
  disabled: Boolean
});
const uploadDraggerEmits = {
  file: (file) => shared.isArray(file)
};

exports.uploadDraggerEmits = uploadDraggerEmits;
exports.uploadDraggerProps = uploadDraggerProps;
//# sourceMappingURL=upload-dragger.js.map
