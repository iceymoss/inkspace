import type { ObjectDirective } from 'vue';
import type { NormalizedWheelEvent } from 'normalize-wheel-es';
export declare const SCOPE = "_Mousewheel";
interface WheelElement extends HTMLElement {
    [SCOPE]: null | {
        wheelHandler?: (event: WheelEvent) => void;
    };
}
type MousewheelCallback = (e: WheelEvent, normalized: NormalizedWheelEvent) => void;
declare const Mousewheel: ObjectDirective<WheelElement, MousewheelCallback>;
export default Mousewheel;
