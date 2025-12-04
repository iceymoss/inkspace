import { buttonProps } from './button.mjs';
import { definePropType } from '../../../utils/vue/props/runtime.mjs';

const buttonGroupProps = {
  size: buttonProps.size,
  type: buttonProps.type,
  direction: {
    type: definePropType(String),
    values: ["horizontal", "vertical"],
    default: "horizontal"
  }
};

export { buttonGroupProps };
//# sourceMappingURL=button-group.mjs.map
