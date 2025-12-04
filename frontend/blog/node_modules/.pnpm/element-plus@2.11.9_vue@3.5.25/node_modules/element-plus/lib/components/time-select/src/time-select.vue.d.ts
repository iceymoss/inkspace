declare const _default: import("vue").DefineComponent<{
    readonly emptyValues: ArrayConstructor;
    readonly valueOnClear: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => string | number | boolean | Function) | (() => string | number | boolean | Function | null) | ((new (...args: any[]) => string | number | boolean | Function) | (() => string | number | boolean | Function | null))[], unknown, unknown, undefined, boolean>;
    readonly format: import("element-plus/es/utils").EpPropFinalized<StringConstructor, unknown, unknown, "HH:mm", boolean>;
    readonly modelValue: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => string) | (() => string | null) | ((new (...args: any[]) => string) | (() => string | null))[], unknown, unknown>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    readonly disabled: BooleanConstructor;
    readonly editable: import("element-plus/es/utils").EpPropFinalized<BooleanConstructor, unknown, unknown, true, boolean>;
    readonly effect: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => string) | (() => import("element-plus").PopperEffect) | ((new (...args: any[]) => string) | (() => import("element-plus").PopperEffect))[], unknown, unknown, "light", boolean>;
    readonly clearable: import("element-plus/es/utils").EpPropFinalized<BooleanConstructor, unknown, unknown, true, boolean>;
    readonly size: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<StringConstructor, "" | "small" | "default" | "large", never>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    readonly placeholder: StringConstructor;
    readonly start: import("element-plus/es/utils").EpPropFinalized<StringConstructor, unknown, unknown, "09:00", boolean>;
    readonly end: import("element-plus/es/utils").EpPropFinalized<StringConstructor, unknown, unknown, "18:00", boolean>;
    readonly step: import("element-plus/es/utils").EpPropFinalized<StringConstructor, unknown, unknown, "00:30", boolean>;
    readonly minTime: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => string) | (() => string | null) | ((new (...args: any[]) => string) | (() => string | null))[], unknown, unknown>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    readonly maxTime: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => string) | (() => string | null) | ((new (...args: any[]) => string) | (() => string | null))[], unknown, unknown>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    readonly includeEndTime: BooleanConstructor;
    readonly name: StringConstructor;
    readonly prefixIcon: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => (string | import("vue").Component) & {}) | (() => string | import("vue").Component) | ((new (...args: any[]) => (string | import("vue").Component) & {}) | (() => string | import("vue").Component))[], unknown, unknown, () => any, boolean>;
    readonly clearIcon: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => (string | import("vue").Component) & {}) | (() => string | import("vue").Component) | ((new (...args: any[]) => (string | import("vue").Component) & {}) | (() => string | import("vue").Component))[], unknown, unknown, () => any, boolean>;
    readonly popperClass: import("element-plus/es/utils").EpPropFinalized<StringConstructor, unknown, unknown, "", boolean>;
    readonly popperStyle: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => string | import("vue").CSSProperties) | (() => string | import("vue").CSSProperties) | ((new (...args: any[]) => string | import("vue").CSSProperties) | (() => string | import("vue").CSSProperties))[], unknown, unknown>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
}, {
    /**
     * @description blur the Input component
     */
    blur: () => void;
    /**
     * @description focus the Input component
     */
    focus: () => void;
}, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {
    focus: (...args: any[]) => void;
    "update:modelValue": (...args: any[]) => void;
    change: (...args: any[]) => void;
    clear: (...args: any[]) => void;
    blur: (...args: any[]) => void;
}, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<import("vue").ExtractPropTypes<{
    readonly emptyValues: ArrayConstructor;
    readonly valueOnClear: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => string | number | boolean | Function) | (() => string | number | boolean | Function | null) | ((new (...args: any[]) => string | number | boolean | Function) | (() => string | number | boolean | Function | null))[], unknown, unknown, undefined, boolean>;
    readonly format: import("element-plus/es/utils").EpPropFinalized<StringConstructor, unknown, unknown, "HH:mm", boolean>;
    readonly modelValue: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => string) | (() => string | null) | ((new (...args: any[]) => string) | (() => string | null))[], unknown, unknown>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    readonly disabled: BooleanConstructor;
    readonly editable: import("element-plus/es/utils").EpPropFinalized<BooleanConstructor, unknown, unknown, true, boolean>;
    readonly effect: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => string) | (() => import("element-plus").PopperEffect) | ((new (...args: any[]) => string) | (() => import("element-plus").PopperEffect))[], unknown, unknown, "light", boolean>;
    readonly clearable: import("element-plus/es/utils").EpPropFinalized<BooleanConstructor, unknown, unknown, true, boolean>;
    readonly size: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<StringConstructor, "" | "small" | "default" | "large", never>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    readonly placeholder: StringConstructor;
    readonly start: import("element-plus/es/utils").EpPropFinalized<StringConstructor, unknown, unknown, "09:00", boolean>;
    readonly end: import("element-plus/es/utils").EpPropFinalized<StringConstructor, unknown, unknown, "18:00", boolean>;
    readonly step: import("element-plus/es/utils").EpPropFinalized<StringConstructor, unknown, unknown, "00:30", boolean>;
    readonly minTime: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => string) | (() => string | null) | ((new (...args: any[]) => string) | (() => string | null))[], unknown, unknown>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    readonly maxTime: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => string) | (() => string | null) | ((new (...args: any[]) => string) | (() => string | null))[], unknown, unknown>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    readonly includeEndTime: BooleanConstructor;
    readonly name: StringConstructor;
    readonly prefixIcon: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => (string | import("vue").Component) & {}) | (() => string | import("vue").Component) | ((new (...args: any[]) => (string | import("vue").Component) & {}) | (() => string | import("vue").Component))[], unknown, unknown, () => any, boolean>;
    readonly clearIcon: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => (string | import("vue").Component) & {}) | (() => string | import("vue").Component) | ((new (...args: any[]) => (string | import("vue").Component) & {}) | (() => string | import("vue").Component))[], unknown, unknown, () => any, boolean>;
    readonly popperClass: import("element-plus/es/utils").EpPropFinalized<StringConstructor, unknown, unknown, "", boolean>;
    readonly popperStyle: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => string | import("vue").CSSProperties) | (() => string | import("vue").CSSProperties) | ((new (...args: any[]) => string | import("vue").CSSProperties) | (() => string | import("vue").CSSProperties))[], unknown, unknown>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
}>> & {
    "onUpdate:modelValue"?: ((...args: any[]) => any) | undefined;
    onChange?: ((...args: any[]) => any) | undefined;
    onFocus?: ((...args: any[]) => any) | undefined;
    onBlur?: ((...args: any[]) => any) | undefined;
    onClear?: ((...args: any[]) => any) | undefined;
}, {
    readonly disabled: boolean;
    readonly end: string;
    readonly start: string;
    readonly effect: import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => string) | (() => import("element-plus").PopperEffect) | ((new (...args: any[]) => string) | (() => import("element-plus").PopperEffect))[], unknown, unknown>;
    readonly valueOnClear: import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => string | number | boolean | Function) | (() => string | number | boolean | Function | null) | ((new (...args: any[]) => string | number | boolean | Function) | (() => string | number | boolean | Function | null))[], unknown, unknown>;
    readonly format: string;
    readonly clearIcon: import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => (string | import("vue").Component) & {}) | (() => string | import("vue").Component) | ((new (...args: any[]) => (string | import("vue").Component) & {}) | (() => string | import("vue").Component))[], unknown, unknown>;
    readonly prefixIcon: import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => (string | import("vue").Component) & {}) | (() => string | import("vue").Component) | ((new (...args: any[]) => (string | import("vue").Component) & {}) | (() => string | import("vue").Component))[], unknown, unknown>;
    readonly clearable: import("element-plus/es/utils").EpPropMergeType<BooleanConstructor, unknown, unknown>;
    readonly popperClass: string;
    readonly editable: import("element-plus/es/utils").EpPropMergeType<BooleanConstructor, unknown, unknown>;
    readonly step: string;
    readonly includeEndTime: boolean;
}>;
export default _default;
