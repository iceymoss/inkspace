declare const _default: import("vue").DefineComponent<{
    readonly colors: {
        readonly type: import("vue").PropType<string[]>;
        readonly required: true;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    readonly color: {
        readonly type: import("vue").PropType<import("../utils/color").default>;
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
}, {}, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, Record<string, any>, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<import("vue").ExtractPropTypes<{
    readonly colors: {
        readonly type: import("vue").PropType<string[]>;
        readonly required: true;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    readonly color: {
        readonly type: import("vue").PropType<import("../utils/color").default>;
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
}>>, {
    readonly disabled: boolean;
}>;
export default _default;
