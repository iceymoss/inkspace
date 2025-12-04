import Steps from './src/steps2.mjs';
import Step from './src/item2.mjs';
export { stepProps } from './src/item.mjs';
export { stepsEmits, stepsProps } from './src/steps.mjs';
export { STEPS_INJECTION_KEY } from './src/tokens.mjs';
import { withInstall, withNoopInstall } from '../../utils/vue/install.mjs';

const ElSteps = withInstall(Steps, {
  Step
});
const ElStep = withNoopInstall(Step);

export { ElStep, ElSteps, ElSteps as default };
//# sourceMappingURL=index.mjs.map
