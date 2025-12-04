import type { CheckboxGroupValueType } from './checkbox-group';
declare function __VLS_template(): {
    default?(_: {}): any;
};
declare const __VLS_component: import("vue").DefineComponent<{
    readonly ariaLabel: StringConstructor;
    readonly modelValue: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => CheckboxGroupValueType) | (() => CheckboxGroupValueType) | ((new (...args: any[]) => CheckboxGroupValueType) | (() => CheckboxGroupValueType))[], unknown, unknown, () => never[], boolean>;
    readonly disabled: BooleanConstructor;
    readonly min: NumberConstructor;
    readonly max: NumberConstructor;
    readonly size: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<StringConstructor, "" | "small" | "default" | "large", never>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    readonly fill: StringConstructor;
    readonly textColor: StringConstructor;
    readonly tag: import("element-plus/es/utils").EpPropFinalized<StringConstructor, unknown, unknown, "div", boolean>;
    readonly validateEvent: import("element-plus/es/utils").EpPropFinalized<BooleanConstructor, unknown, unknown, true, boolean>;
    readonly options: {
        readonly type: import("vue").PropType<import("./checkbox-group").CheckboxOption[]>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    readonly props: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => {
        value?: string;
        label?: string;
        disabled?: string;
    }) | (() => {
        value?: string;
        label?: string;
        disabled?: string;
    }) | ((new (...args: any[]) => {
        value?: string;
        label?: string;
        disabled?: string;
    }) | (() => {
        value?: string;
        label?: string;
        disabled?: string;
    }))[], unknown, unknown, () => Required<{
        value?: string;
        label?: string;
        disabled?: string;
    }>, boolean>;
    readonly type: import("element-plus/es/utils").EpPropFinalized<StringConstructor, "button" | "checkbox", unknown, "checkbox", boolean>;
}, {}, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {
    "update:modelValue": (val: CheckboxGroupValueType) => void;
    change: (val: import("element-plus").CheckboxValueType[]) => void;
}, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<import("vue").ExtractPropTypes<{
    readonly ariaLabel: StringConstructor;
    readonly modelValue: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => CheckboxGroupValueType) | (() => CheckboxGroupValueType) | ((new (...args: any[]) => CheckboxGroupValueType) | (() => CheckboxGroupValueType))[], unknown, unknown, () => never[], boolean>;
    readonly disabled: BooleanConstructor;
    readonly min: NumberConstructor;
    readonly max: NumberConstructor;
    readonly size: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<StringConstructor, "" | "small" | "default" | "large", never>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    readonly fill: StringConstructor;
    readonly textColor: StringConstructor;
    readonly tag: import("element-plus/es/utils").EpPropFinalized<StringConstructor, unknown, unknown, "div", boolean>;
    readonly validateEvent: import("element-plus/es/utils").EpPropFinalized<BooleanConstructor, unknown, unknown, true, boolean>;
    readonly options: {
        readonly type: import("vue").PropType<import("./checkbox-group").CheckboxOption[]>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    readonly props: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => {
        value?: string;
        label?: string;
        disabled?: string;
    }) | (() => {
        value?: string;
        label?: string;
        disabled?: string;
    }) | ((new (...args: any[]) => {
        value?: string;
        label?: string;
        disabled?: string;
    }) | (() => {
        value?: string;
        label?: string;
        disabled?: string;
    }))[], unknown, unknown, () => Required<{
        value?: string;
        label?: string;
        disabled?: string;
    }>, boolean>;
    readonly type: import("element-plus/es/utils").EpPropFinalized<StringConstructor, "button" | "checkbox", unknown, "checkbox", boolean>;
}>> & {
    "onUpdate:modelValue"?: ((val: CheckboxGroupValueType) => any) | undefined;
    onChange?: ((val: import("element-plus").CheckboxValueType[]) => any) | undefined;
}, {
    readonly disabled: boolean;
    readonly type: import("element-plus/es/utils").EpPropMergeType<StringConstructor, "button" | "checkbox", unknown>;
    readonly props: {
        value?: string;
        label?: string;
        disabled?: string;
    };
    readonly modelValue: CheckboxGroupValueType;
    readonly validateEvent: import("element-plus/es/utils").EpPropMergeType<BooleanConstructor, unknown, unknown>;
    readonly tag: string;
}>;
declare const _default: __VLS_WithTemplateSlots<typeof __VLS_component, ReturnType<typeof __VLS_template>>;
export default _default;
type __VLS_WithTemplateSlots<T, S> = T & {
    new (): {
        $slots: S;
    };
};
