import type { RadioPropsPublic } from './radio';
import type { ExtractPropTypes, __ExtractPublicPropTypes } from 'vue';
import type RadioGroup from './radio-group.vue';
export declare const radioGroupProps: {
    readonly ariaLabel: StringConstructor;
    readonly id: import("element-plus/es/utils").EpPropFinalized<StringConstructor, unknown, unknown, undefined, boolean>;
    readonly size: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<StringConstructor, "" | "small" | "default" | "large", never>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    readonly disabled: BooleanConstructor;
    readonly modelValue: import("element-plus/es/utils").EpPropFinalized<readonly [StringConstructor, NumberConstructor, BooleanConstructor], unknown, unknown, undefined, boolean>;
    readonly fill: import("element-plus/es/utils").EpPropFinalized<StringConstructor, unknown, unknown, "", boolean>;
    readonly textColor: import("element-plus/es/utils").EpPropFinalized<StringConstructor, unknown, unknown, "", boolean>;
    readonly name: import("element-plus/es/utils").EpPropFinalized<StringConstructor, unknown, unknown, undefined, boolean>;
    readonly validateEvent: import("element-plus/es/utils").EpPropFinalized<BooleanConstructor, unknown, unknown, true, boolean>;
    readonly options: {
        readonly type: import("vue").PropType<radioOption[]>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    readonly props: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => radioOptionProp) | (() => radioOptionProp) | ((new (...args: any[]) => radioOptionProp) | (() => radioOptionProp))[], unknown, unknown, () => Required<radioOptionProp>, boolean>;
    readonly type: import("element-plus/es/utils").EpPropFinalized<StringConstructor, "button" | "radio", unknown, "radio", boolean>;
};
export type RadioGroupProps = ExtractPropTypes<typeof radioGroupProps>;
export type RadioGroupPropsPublic = __ExtractPublicPropTypes<typeof radioGroupProps>;
export declare const radioGroupEmits: {
    "update:modelValue": (val: string | number | boolean | undefined) => val is string | number | boolean;
    change: (val: string | number | boolean | undefined) => val is string | number | boolean;
};
export type RadioGroupEmits = typeof radioGroupEmits;
export type RadioGroupInstance = InstanceType<typeof RadioGroup> & unknown;
export type radioOption = RadioPropsPublic & Record<string, any>;
export declare const radioDefaultProps: Required<radioOptionProp>;
export type radioOptionProp = {
    value?: string;
    label?: string;
    disabled?: string;
};
