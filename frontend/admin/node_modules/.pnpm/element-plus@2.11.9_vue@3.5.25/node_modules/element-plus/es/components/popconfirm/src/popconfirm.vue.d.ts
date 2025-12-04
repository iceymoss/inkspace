declare function __VLS_template(): {
    actions?(_: {
        confirm: (e: MouseEvent) => void;
        cancel: (e: MouseEvent) => void;
    }): any;
    reference?(_: {}): any;
};
declare const __VLS_component: import("vue").DefineComponent<{
    readonly title: StringConstructor;
    readonly confirmButtonText: StringConstructor;
    readonly cancelButtonText: StringConstructor;
    readonly confirmButtonType: import("element-plus/es/utils").EpPropFinalized<StringConstructor, "" | "text" | "default" | "primary" | "success" | "warning" | "info" | "danger", unknown, "primary", boolean>;
    readonly cancelButtonType: import("element-plus/es/utils").EpPropFinalized<StringConstructor, "" | "text" | "default" | "primary" | "success" | "warning" | "info" | "danger", unknown, "text", boolean>;
    readonly icon: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => (string | import("vue").Component) & {}) | (() => string | import("vue").Component) | ((new (...args: any[]) => (string | import("vue").Component) & {}) | (() => string | import("vue").Component))[], unknown, unknown, () => any, boolean>;
    readonly iconColor: import("element-plus/es/utils").EpPropFinalized<StringConstructor, unknown, unknown, "#f90", boolean>;
    readonly hideIcon: BooleanConstructor;
    readonly hideAfter: import("element-plus/es/utils").EpPropFinalized<NumberConstructor, unknown, unknown, 200, boolean>;
    readonly effect: {
        readonly default: "light";
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => string) | (() => import("element-plus").PopperEffect) | ((new (...args: any[]) => string) | (() => import("element-plus").PopperEffect))[], unknown, unknown>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        readonly __epPropKey: true;
    };
    readonly teleported: import("element-plus/es/utils").EpPropFinalized<BooleanConstructor, unknown, unknown, true, boolean>;
    readonly persistent: BooleanConstructor;
    readonly width: import("element-plus/es/utils").EpPropFinalized<readonly [StringConstructor, NumberConstructor], unknown, unknown, 150, boolean>;
    readonly virtualTriggering: BooleanConstructor;
    readonly virtualRef: {
        readonly type: import("vue").PropType<import("element-plus").Measurable>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
}, {
    popperRef: import("vue").ComputedRef<import("element-plus").PopperInstance | undefined>;
    hide: () => void;
}, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {
    cancel: (e: MouseEvent) => void;
    confirm: (e: MouseEvent) => void;
}, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<import("vue").ExtractPropTypes<{
    readonly title: StringConstructor;
    readonly confirmButtonText: StringConstructor;
    readonly cancelButtonText: StringConstructor;
    readonly confirmButtonType: import("element-plus/es/utils").EpPropFinalized<StringConstructor, "" | "text" | "default" | "primary" | "success" | "warning" | "info" | "danger", unknown, "primary", boolean>;
    readonly cancelButtonType: import("element-plus/es/utils").EpPropFinalized<StringConstructor, "" | "text" | "default" | "primary" | "success" | "warning" | "info" | "danger", unknown, "text", boolean>;
    readonly icon: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => (string | import("vue").Component) & {}) | (() => string | import("vue").Component) | ((new (...args: any[]) => (string | import("vue").Component) & {}) | (() => string | import("vue").Component))[], unknown, unknown, () => any, boolean>;
    readonly iconColor: import("element-plus/es/utils").EpPropFinalized<StringConstructor, unknown, unknown, "#f90", boolean>;
    readonly hideIcon: BooleanConstructor;
    readonly hideAfter: import("element-plus/es/utils").EpPropFinalized<NumberConstructor, unknown, unknown, 200, boolean>;
    readonly effect: {
        readonly default: "light";
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => string) | (() => import("element-plus").PopperEffect) | ((new (...args: any[]) => string) | (() => import("element-plus").PopperEffect))[], unknown, unknown>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        readonly __epPropKey: true;
    };
    readonly teleported: import("element-plus/es/utils").EpPropFinalized<BooleanConstructor, unknown, unknown, true, boolean>;
    readonly persistent: BooleanConstructor;
    readonly width: import("element-plus/es/utils").EpPropFinalized<readonly [StringConstructor, NumberConstructor], unknown, unknown, 150, boolean>;
    readonly virtualTriggering: BooleanConstructor;
    readonly virtualRef: {
        readonly type: import("vue").PropType<import("element-plus").Measurable>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
}>> & {
    onConfirm?: ((e: MouseEvent) => any) | undefined;
    onCancel?: ((e: MouseEvent) => any) | undefined;
}, {
    readonly width: import("element-plus/es/utils").EpPropMergeType<readonly [StringConstructor, NumberConstructor], unknown, unknown>;
    readonly icon: import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => (string | import("vue").Component) & {}) | (() => string | import("vue").Component) | ((new (...args: any[]) => (string | import("vue").Component) & {}) | (() => string | import("vue").Component))[], unknown, unknown>;
    readonly effect: import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => string) | (() => import("element-plus").PopperEffect) | ((new (...args: any[]) => string) | (() => import("element-plus").PopperEffect))[], unknown, unknown>;
    readonly hideAfter: number;
    readonly virtualTriggering: boolean;
    readonly teleported: import("element-plus/es/utils").EpPropMergeType<BooleanConstructor, unknown, unknown>;
    readonly persistent: boolean;
    readonly confirmButtonType: import("element-plus/es/utils").EpPropMergeType<StringConstructor, "" | "text" | "default" | "primary" | "success" | "warning" | "info" | "danger", unknown>;
    readonly cancelButtonType: import("element-plus/es/utils").EpPropMergeType<StringConstructor, "" | "text" | "default" | "primary" | "success" | "warning" | "info" | "danger", unknown>;
    readonly iconColor: string;
    readonly hideIcon: boolean;
}>;
declare const _default: __VLS_WithTemplateSlots<typeof __VLS_component, ReturnType<typeof __VLS_template>>;
export default _default;
type __VLS_WithTemplateSlots<T, S> = T & {
    new (): {
        $slots: S;
    };
};
