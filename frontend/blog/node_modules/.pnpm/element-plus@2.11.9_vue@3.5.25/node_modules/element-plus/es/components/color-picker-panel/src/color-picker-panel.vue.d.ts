import type { InputInstance } from 'element-plus/es/components/input';
declare function update(): void;
declare function __VLS_template(): {
    footer?(_: {}): any;
};
declare const __VLS_component: import("vue").DefineComponent<{
    readonly modelValue: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => string) | (() => string | null) | ((new (...args: any[]) => string) | (() => string | null))[], unknown, unknown, undefined, boolean>;
    readonly border: import("element-plus/es/utils").EpPropFinalized<BooleanConstructor, unknown, unknown, true, boolean>;
    readonly showAlpha: BooleanConstructor;
    readonly colorFormat: StringConstructor;
    readonly disabled: BooleanConstructor;
    readonly predefine: {
        readonly type: import("vue").PropType<string[]>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    readonly validateEvent: import("element-plus/es/utils").EpPropFinalized<BooleanConstructor, unknown, unknown, true, boolean>;
}, {
    /**
     * @description current color object
     */
    color: import("./utils/color.js").default;
    /**
     * @description custom input ref
     */
    inputRef: import("vue").Ref<InputInstance | undefined>;
    /**
     * @description update sub components
     */
    update: typeof update;
}, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {
    "update:modelValue": (val: string | null) => void;
}, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<import("vue").ExtractPropTypes<{
    readonly modelValue: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => string) | (() => string | null) | ((new (...args: any[]) => string) | (() => string | null))[], unknown, unknown, undefined, boolean>;
    readonly border: import("element-plus/es/utils").EpPropFinalized<BooleanConstructor, unknown, unknown, true, boolean>;
    readonly showAlpha: BooleanConstructor;
    readonly colorFormat: StringConstructor;
    readonly disabled: BooleanConstructor;
    readonly predefine: {
        readonly type: import("vue").PropType<string[]>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    readonly validateEvent: import("element-plus/es/utils").EpPropFinalized<BooleanConstructor, unknown, unknown, true, boolean>;
}>> & {
    "onUpdate:modelValue"?: ((val: string | null) => any) | undefined;
}, {
    readonly disabled: boolean;
    readonly border: import("element-plus/es/utils").EpPropMergeType<BooleanConstructor, unknown, unknown>;
    readonly modelValue: import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => string) | (() => string | null) | ((new (...args: any[]) => string) | (() => string | null))[], unknown, unknown>;
    readonly validateEvent: import("element-plus/es/utils").EpPropMergeType<BooleanConstructor, unknown, unknown>;
    readonly showAlpha: boolean;
}>;
declare const _default: __VLS_WithTemplateSlots<typeof __VLS_component, ReturnType<typeof __VLS_template>>;
export default _default;
type __VLS_WithTemplateSlots<T, S> = T & {
    new (): {
        $slots: S;
    };
};
