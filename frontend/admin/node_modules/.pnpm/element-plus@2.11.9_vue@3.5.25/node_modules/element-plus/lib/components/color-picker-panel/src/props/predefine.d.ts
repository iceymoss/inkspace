import type { ExtractPropTypes, __ExtractPublicPropTypes } from 'vue';
import type Color from '../utils/color';
export declare const predefineProps: {
    readonly colors: {
        readonly type: import("vue").PropType<string[]>;
        readonly required: true;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    readonly color: {
        readonly type: import("vue").PropType<Color>;
        readonly required: true;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    readonly enableAlpha: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<BooleanConstructor, unknown, unknown>>;
        readonly required: true;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    readonly disabled: BooleanConstructor;
};
export type PredefineProps = ExtractPropTypes<typeof predefineProps>;
export type PredefinePropsPublic = __ExtractPublicPropTypes<typeof predefineProps>;
