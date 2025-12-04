'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var splitter$1 = require('./src/splitter2.js');
var splitPanel$1 = require('./src/split-panel2.js');
var splitter = require('./src/splitter.js');
var splitPanel = require('./src/split-panel.js');
var install = require('../../utils/vue/install.js');

const ElSplitter = install.withInstall(splitter$1["default"], {
  SplitPanel: splitPanel$1["default"]
});
const ElSplitterPanel = install.withNoopInstall(splitPanel$1["default"]);

exports.splitterEmits = splitter.splitterEmits;
exports.splitterProps = splitter.splitterProps;
exports.splitterPanelEmits = splitPanel.splitterPanelEmits;
exports.splitterPanelProps = splitPanel.splitterPanelProps;
exports.ElSplitter = ElSplitter;
exports.ElSplitterPanel = ElSplitterPanel;
exports["default"] = ElSplitter;
//# sourceMappingURL=index.js.map
