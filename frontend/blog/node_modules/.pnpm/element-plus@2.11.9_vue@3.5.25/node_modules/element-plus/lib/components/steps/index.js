'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var steps$1 = require('./src/steps2.js');
var item$1 = require('./src/item2.js');
var item = require('./src/item.js');
var steps = require('./src/steps.js');
var tokens = require('./src/tokens.js');
var install = require('../../utils/vue/install.js');

const ElSteps = install.withInstall(steps$1["default"], {
  Step: item$1["default"]
});
const ElStep = install.withNoopInstall(item$1["default"]);

exports.stepProps = item.stepProps;
exports.stepsEmits = steps.stepsEmits;
exports.stepsProps = steps.stepsProps;
exports.STEPS_INJECTION_KEY = tokens.STEPS_INJECTION_KEY;
exports.ElStep = ElStep;
exports.ElSteps = ElSteps;
exports["default"] = ElSteps;
//# sourceMappingURL=index.js.map
