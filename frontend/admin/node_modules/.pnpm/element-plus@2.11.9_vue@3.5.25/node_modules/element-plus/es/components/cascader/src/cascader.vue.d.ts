import type { ComputedRef, StyleValue } from 'vue';
import type { CascaderNode, CascaderPanelInstance, CascaderValue, Tag } from 'element-plus/es/components/cascader-panel';
declare function __VLS_template(): {
    prefix?(_: {}): any;
    tag?(_: {
        data: {
            node?: {
                readonly uid: number;
                readonly level: number;
                readonly value: import("element-plus/es/components/cascader-panel").CascaderNodeValue;
                readonly label: string;
                readonly pathNodes: any[];
                readonly pathValues: import("element-plus/es/components/cascader-panel").CascaderNodeValue[];
                readonly pathLabels: string[];
                childrenData: {
                    [x: string]: unknown;
                    label?: string | undefined;
                    value?: import("element-plus/es/components/cascader-panel").CascaderNodeValue | undefined;
                    children?: any[] | undefined;
                    disabled?: boolean | undefined;
                    leaf?: boolean | undefined;
                }[] | undefined;
                children: any[];
                text: string;
                loaded: boolean;
                checked: boolean;
                indeterminate: boolean;
                loading: boolean;
                readonly data: {
                    [x: string]: unknown;
                    label?: string | undefined;
                    value?: import("element-plus/es/components/cascader-panel").CascaderNodeValue | undefined;
                    children?: any[] | undefined;
                    disabled?: boolean | undefined;
                    leaf?: boolean | undefined;
                };
                readonly config: {
                    expandTrigger: import("element-plus/es/components/cascader-panel").ExpandTrigger;
                    multiple: boolean;
                    checkStrictly: boolean;
                    emitPath: boolean;
                    lazy: boolean;
                    lazyLoad: import("element-plus/es/components/cascader-panel").LazyLoad;
                    value: string;
                    label: string;
                    children: string;
                    disabled: string | import("element-plus/es/components/cascader-panel").isDisabled;
                    leaf: string | import("element-plus/es/components/cascader-panel").isLeaf;
                    hoverThreshold: number;
                    checkOnClickNode: boolean;
                    checkOnClickLeaf: boolean;
                    showPrefix: boolean;
                };
                readonly parent?: any | undefined;
                readonly root: boolean;
                readonly isDisabled: boolean;
                readonly isLeaf: boolean;
                readonly valueByOption: import("element-plus/es/components/cascader-panel").CascaderNodeValue | import("element-plus/es/components/cascader-panel").CascaderNodeValue[];
                appendChild: (childData: import("element-plus/es/components/cascader-panel").CascaderOption) => CascaderNode;
                calcText: (allLevels: boolean, separator: string) => string;
                broadcast: (checked: boolean) => void;
                emit: () => void;
                onParentCheck: (checked: boolean) => void;
                onChildCheck: () => void;
                setCheckState: (checked: boolean) => void;
                doCheck: (checked: boolean) => void;
            } | undefined;
            key: number;
            text: string;
            hitState?: boolean | undefined;
            closable: boolean;
        }[];
        deleteTag: (tag: Tag) => void;
    }): any;
    header?(_: {}): any;
    empty?(_: {}): any;
    empty?(_: {}): any;
    "suggestion-item"?(_: {
        item: {
            readonly uid: number;
            readonly level: number;
            readonly value: import("element-plus/es/components/cascader-panel").CascaderNodeValue;
            readonly label: string;
            readonly pathNodes: any[];
            readonly pathValues: import("element-plus/es/components/cascader-panel").CascaderNodeValue[];
            readonly pathLabels: string[];
            childrenData: {
                [x: string]: unknown;
                label?: string | undefined;
                value?: import("element-plus/es/components/cascader-panel").CascaderNodeValue | undefined;
                children?: any[] | undefined;
                disabled?: boolean | undefined;
                leaf?: boolean | undefined;
            }[] | undefined;
            children: any[];
            text: string;
            loaded: boolean;
            checked: boolean;
            indeterminate: boolean;
            loading: boolean;
            readonly data: {
                [x: string]: unknown;
                label?: string | undefined;
                value?: import("element-plus/es/components/cascader-panel").CascaderNodeValue | undefined;
                children?: any[] | undefined;
                disabled?: boolean | undefined;
                leaf?: boolean | undefined;
            };
            readonly config: {
                expandTrigger: import("element-plus/es/components/cascader-panel").ExpandTrigger;
                multiple: boolean;
                checkStrictly: boolean;
                emitPath: boolean;
                lazy: boolean;
                lazyLoad: import("element-plus/es/components/cascader-panel").LazyLoad;
                value: string;
                label: string;
                children: string;
                disabled: string | import("element-plus/es/components/cascader-panel").isDisabled;
                leaf: string | import("element-plus/es/components/cascader-panel").isLeaf;
                hoverThreshold: number;
                checkOnClickNode: boolean;
                checkOnClickLeaf: boolean;
                showPrefix: boolean;
            };
            readonly parent?: any | undefined;
            readonly root: boolean;
            readonly isDisabled: boolean;
            readonly isLeaf: boolean;
            readonly valueByOption: import("element-plus/es/components/cascader-panel").CascaderNodeValue | import("element-plus/es/components/cascader-panel").CascaderNodeValue[];
            appendChild: (childData: import("element-plus/es/components/cascader-panel").CascaderOption) => CascaderNode;
            calcText: (allLevels: boolean, separator: string) => string;
            broadcast: (checked: boolean) => void;
            emit: () => void;
            onParentCheck: (checked: boolean) => void;
            onChildCheck: () => void;
            setCheckState: (checked: boolean) => void;
            doCheck: (checked: boolean) => void;
        };
    }): any;
    footer?(_: {}): any;
};
declare const __VLS_component: import("vue").DefineComponent<{
    emptyValues: ArrayConstructor;
    valueOnClear: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => string | number | boolean | Function) | (() => string | number | boolean | Function | null) | ((new (...args: any[]) => string | number | boolean | Function) | (() => string | number | boolean | Function | null))[], unknown, unknown, undefined, boolean>;
    size: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<StringConstructor, "" | "small" | "default" | "large", never>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    placeholder: StringConstructor;
    disabled: BooleanConstructor;
    clearable: BooleanConstructor;
    clearIcon: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => (string | import("vue").Component) & {}) | (() => string | import("vue").Component) | ((new (...args: any[]) => (string | import("vue").Component) & {}) | (() => string | import("vue").Component))[], unknown, unknown>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    filterable: BooleanConstructor;
    filterMethod: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => (node: CascaderNode, keyword: string) => boolean) | (() => (node: CascaderNode, keyword: string) => boolean) | {
        (): (node: CascaderNode, keyword: string) => boolean;
        new (): any;
        readonly prototype: any;
    } | ((new (...args: any[]) => (node: CascaderNode, keyword: string) => boolean) | (() => (node: CascaderNode, keyword: string) => boolean) | {
        (): (node: CascaderNode, keyword: string) => boolean;
        new (): any;
        readonly prototype: any;
    })[], unknown, unknown, (node: CascaderNode, keyword: string) => boolean, boolean>;
    separator: import("element-plus/es/utils").EpPropFinalized<StringConstructor, unknown, unknown, string, boolean>;
    showAllLevels: import("element-plus/es/utils").EpPropFinalized<BooleanConstructor, unknown, unknown, boolean, boolean>;
    collapseTags: BooleanConstructor;
    maxCollapseTags: import("element-plus/es/utils").EpPropFinalized<NumberConstructor, unknown, unknown, number, boolean>;
    collapseTagsTooltip: BooleanConstructor;
    maxCollapseTagsTooltipHeight: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<(NumberConstructor | StringConstructor)[], unknown, unknown>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    debounce: import("element-plus/es/utils").EpPropFinalized<NumberConstructor, unknown, unknown, number, boolean>;
    beforeFilter: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => (value: string) => boolean | Promise<any>) | (() => (value: string) => boolean | Promise<any>) | {
        (): (value: string) => boolean | Promise<any>;
        new (): any;
        readonly prototype: any;
    } | ((new (...args: any[]) => (value: string) => boolean | Promise<any>) | (() => (value: string) => boolean | Promise<any>) | {
        (): (value: string) => boolean | Promise<any>;
        new (): any;
        readonly prototype: any;
    })[], unknown, unknown, () => true, boolean>;
    placement: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => "left" | "right" | "top" | "bottom" | "auto" | "auto-start" | "auto-end" | "top-start" | "top-end" | "bottom-start" | "bottom-end" | "right-start" | "right-end" | "left-start" | "left-end") | (() => import("element-plus/es/components/popper").Placement) | ((new (...args: any[]) => "left" | "right" | "top" | "bottom" | "auto" | "auto-start" | "auto-end" | "top-start" | "top-end" | "bottom-start" | "bottom-end" | "right-start" | "right-end" | "left-start" | "left-end") | (() => import("element-plus/es/components/popper").Placement))[], import("element-plus/es/components/popper").Placement, unknown, string, boolean>;
    fallbackPlacements: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => import("element-plus/es/components/popper").Placement[]) | (() => import("element-plus/es/components/popper").Placement[]) | ((new (...args: any[]) => import("element-plus/es/components/popper").Placement[]) | (() => import("element-plus/es/components/popper").Placement[]))[], unknown, unknown, string[], boolean>;
    popperClass: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | any)[])[])[])[])[])[])[])[])[])[])[]) | (() => string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | any)[])[])[])[])[])[])[])[])[])[])[]) | ((new (...args: any[]) => string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | any)[])[])[])[])[])[])[])[])[])[])[]) | (() => string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | any)[])[])[])[])[])[])[])[])[])[])[]))[], unknown, unknown>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    popperStyle: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => string | import("vue").CSSProperties | StyleValue[]) | (() => StyleValue) | ((new (...args: any[]) => string | import("vue").CSSProperties | StyleValue[]) | (() => StyleValue))[], unknown, unknown>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    teleported: import("element-plus/es/utils").EpPropFinalized<BooleanConstructor, unknown, unknown, true, boolean>;
    effect: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => string) | (() => import("element-plus/es/components/popper").PopperEffect) | ((new (...args: any[]) => string) | (() => import("element-plus/es/components/popper").PopperEffect))[], unknown, unknown, string, boolean>;
    tagType: {
        default: string;
        type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<StringConstructor, "primary" | "success" | "warning" | "info" | "danger", unknown>>;
        required: false;
        validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    tagEffect: {
        default: string;
        type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<StringConstructor, "dark" | "light" | "plain", unknown>>;
        required: false;
        validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    validateEvent: import("element-plus/es/utils").EpPropFinalized<BooleanConstructor, unknown, unknown, boolean, boolean>;
    persistent: import("element-plus/es/utils").EpPropFinalized<BooleanConstructor, unknown, unknown, boolean, boolean>;
    showCheckedStrategy: import("element-plus/es/utils").EpPropFinalized<StringConstructor, string, unknown, string, boolean>;
    checkOnClickNode: BooleanConstructor;
    showPrefix: import("element-plus/es/utils").EpPropFinalized<BooleanConstructor, unknown, unknown, boolean, boolean>;
    modelValue: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => string | number | Record<string, any> | import("element-plus/es/components/cascader-panel").CascaderNodePathValue | (import("element-plus/es/components/cascader-panel").CascaderNodeValue | import("element-plus/es/components/cascader-panel").CascaderNodePathValue)[]) | (() => CascaderValue | null) | ((new (...args: any[]) => string | number | Record<string, any> | import("element-plus/es/components/cascader-panel").CascaderNodePathValue | (import("element-plus/es/components/cascader-panel").CascaderNodeValue | import("element-plus/es/components/cascader-panel").CascaderNodePathValue)[]) | (() => CascaderValue | null))[], unknown, unknown>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    options: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => import("element-plus/es/components/cascader-panel").CascaderOption[]) | (() => import("element-plus/es/components/cascader-panel").CascaderOption[]) | ((new (...args: any[]) => import("element-plus/es/components/cascader-panel").CascaderOption[]) | (() => import("element-plus/es/components/cascader-panel").CascaderOption[]))[], unknown, unknown, () => import("element-plus/es/components/cascader-panel").CascaderOption[], boolean>;
    props: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => import("element-plus/es/components/cascader-panel").CascaderProps) | (() => import("element-plus/es/components/cascader-panel").CascaderProps) | ((new (...args: any[]) => import("element-plus/es/components/cascader-panel").CascaderProps) | (() => import("element-plus/es/components/cascader-panel").CascaderProps))[], unknown, unknown, () => import("element-plus/es/components/cascader-panel").CascaderProps, boolean>;
}, {
    /**
     * @description get an array of currently selected node,(leafOnly) whether only return the leaf checked nodes, default is `false`
     */
    getCheckedNodes: (leafOnly: boolean) => CascaderNode[] | undefined;
    /**
     * @description cascader panel ref
     */
    cascaderPanelRef: import("vue").Ref<CascaderPanelInstance | undefined>;
    /**
     * @description toggle the visible of popper
     */
    togglePopperVisible: (visible?: boolean) => void;
    /**
     * @description cascader content ref
     */
    contentRef: ComputedRef<HTMLElement | undefined>;
    /**
     * @description selected content text
     */
    presentText: ComputedRef<string>;
    /** @description focus the input element */
    focus: () => void;
    /** @description blur the input element */
    blur: () => void;
}, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {
    focus: (evt: FocusEvent) => void;
    "update:modelValue": (value: CascaderValue | null | undefined) => void;
    change: (value: CascaderValue | null | undefined) => void;
    clear: () => void;
    blur: (evt: FocusEvent) => void;
    visibleChange: (val: boolean) => void;
    expandChange: (val: CascaderValue) => void;
    removeTag: (val: import("element-plus/es/components/cascader-panel").CascaderNodeValue | import("element-plus/es/components/cascader-panel").CascaderNodePathValue) => void;
}, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<import("vue").ExtractPropTypes<{
    emptyValues: ArrayConstructor;
    valueOnClear: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => string | number | boolean | Function) | (() => string | number | boolean | Function | null) | ((new (...args: any[]) => string | number | boolean | Function) | (() => string | number | boolean | Function | null))[], unknown, unknown, undefined, boolean>;
    size: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<StringConstructor, "" | "small" | "default" | "large", never>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    placeholder: StringConstructor;
    disabled: BooleanConstructor;
    clearable: BooleanConstructor;
    clearIcon: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => (string | import("vue").Component) & {}) | (() => string | import("vue").Component) | ((new (...args: any[]) => (string | import("vue").Component) & {}) | (() => string | import("vue").Component))[], unknown, unknown>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    filterable: BooleanConstructor;
    filterMethod: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => (node: CascaderNode, keyword: string) => boolean) | (() => (node: CascaderNode, keyword: string) => boolean) | {
        (): (node: CascaderNode, keyword: string) => boolean;
        new (): any;
        readonly prototype: any;
    } | ((new (...args: any[]) => (node: CascaderNode, keyword: string) => boolean) | (() => (node: CascaderNode, keyword: string) => boolean) | {
        (): (node: CascaderNode, keyword: string) => boolean;
        new (): any;
        readonly prototype: any;
    })[], unknown, unknown, (node: CascaderNode, keyword: string) => boolean, boolean>;
    separator: import("element-plus/es/utils").EpPropFinalized<StringConstructor, unknown, unknown, string, boolean>;
    showAllLevels: import("element-plus/es/utils").EpPropFinalized<BooleanConstructor, unknown, unknown, boolean, boolean>;
    collapseTags: BooleanConstructor;
    maxCollapseTags: import("element-plus/es/utils").EpPropFinalized<NumberConstructor, unknown, unknown, number, boolean>;
    collapseTagsTooltip: BooleanConstructor;
    maxCollapseTagsTooltipHeight: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<(NumberConstructor | StringConstructor)[], unknown, unknown>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    debounce: import("element-plus/es/utils").EpPropFinalized<NumberConstructor, unknown, unknown, number, boolean>;
    beforeFilter: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => (value: string) => boolean | Promise<any>) | (() => (value: string) => boolean | Promise<any>) | {
        (): (value: string) => boolean | Promise<any>;
        new (): any;
        readonly prototype: any;
    } | ((new (...args: any[]) => (value: string) => boolean | Promise<any>) | (() => (value: string) => boolean | Promise<any>) | {
        (): (value: string) => boolean | Promise<any>;
        new (): any;
        readonly prototype: any;
    })[], unknown, unknown, () => true, boolean>;
    placement: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => "left" | "right" | "top" | "bottom" | "auto" | "auto-start" | "auto-end" | "top-start" | "top-end" | "bottom-start" | "bottom-end" | "right-start" | "right-end" | "left-start" | "left-end") | (() => import("element-plus/es/components/popper").Placement) | ((new (...args: any[]) => "left" | "right" | "top" | "bottom" | "auto" | "auto-start" | "auto-end" | "top-start" | "top-end" | "bottom-start" | "bottom-end" | "right-start" | "right-end" | "left-start" | "left-end") | (() => import("element-plus/es/components/popper").Placement))[], import("element-plus/es/components/popper").Placement, unknown, string, boolean>;
    fallbackPlacements: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => import("element-plus/es/components/popper").Placement[]) | (() => import("element-plus/es/components/popper").Placement[]) | ((new (...args: any[]) => import("element-plus/es/components/popper").Placement[]) | (() => import("element-plus/es/components/popper").Placement[]))[], unknown, unknown, string[], boolean>;
    popperClass: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | any)[])[])[])[])[])[])[])[])[])[])[]) | (() => string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | any)[])[])[])[])[])[])[])[])[])[])[]) | ((new (...args: any[]) => string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | any)[])[])[])[])[])[])[])[])[])[])[]) | (() => string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | (string | {
            [x: string]: boolean;
        } | any)[])[])[])[])[])[])[])[])[])[])[]))[], unknown, unknown>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    popperStyle: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => string | import("vue").CSSProperties | StyleValue[]) | (() => StyleValue) | ((new (...args: any[]) => string | import("vue").CSSProperties | StyleValue[]) | (() => StyleValue))[], unknown, unknown>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    teleported: import("element-plus/es/utils").EpPropFinalized<BooleanConstructor, unknown, unknown, true, boolean>;
    effect: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => string) | (() => import("element-plus/es/components/popper").PopperEffect) | ((new (...args: any[]) => string) | (() => import("element-plus/es/components/popper").PopperEffect))[], unknown, unknown, string, boolean>;
    tagType: {
        default: string;
        type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<StringConstructor, "primary" | "success" | "warning" | "info" | "danger", unknown>>;
        required: false;
        validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    tagEffect: {
        default: string;
        type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<StringConstructor, "dark" | "light" | "plain", unknown>>;
        required: false;
        validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    validateEvent: import("element-plus/es/utils").EpPropFinalized<BooleanConstructor, unknown, unknown, boolean, boolean>;
    persistent: import("element-plus/es/utils").EpPropFinalized<BooleanConstructor, unknown, unknown, boolean, boolean>;
    showCheckedStrategy: import("element-plus/es/utils").EpPropFinalized<StringConstructor, string, unknown, string, boolean>;
    checkOnClickNode: BooleanConstructor;
    showPrefix: import("element-plus/es/utils").EpPropFinalized<BooleanConstructor, unknown, unknown, boolean, boolean>;
    modelValue: {
        readonly type: import("vue").PropType<import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => string | number | Record<string, any> | import("element-plus/es/components/cascader-panel").CascaderNodePathValue | (import("element-plus/es/components/cascader-panel").CascaderNodeValue | import("element-plus/es/components/cascader-panel").CascaderNodePathValue)[]) | (() => CascaderValue | null) | ((new (...args: any[]) => string | number | Record<string, any> | import("element-plus/es/components/cascader-panel").CascaderNodePathValue | (import("element-plus/es/components/cascader-panel").CascaderNodeValue | import("element-plus/es/components/cascader-panel").CascaderNodePathValue)[]) | (() => CascaderValue | null))[], unknown, unknown>>;
        readonly required: false;
        readonly validator: ((val: unknown) => boolean) | undefined;
        __epPropKey: true;
    };
    options: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => import("element-plus/es/components/cascader-panel").CascaderOption[]) | (() => import("element-plus/es/components/cascader-panel").CascaderOption[]) | ((new (...args: any[]) => import("element-plus/es/components/cascader-panel").CascaderOption[]) | (() => import("element-plus/es/components/cascader-panel").CascaderOption[]))[], unknown, unknown, () => import("element-plus/es/components/cascader-panel").CascaderOption[], boolean>;
    props: import("element-plus/es/utils").EpPropFinalized<(new (...args: any[]) => import("element-plus/es/components/cascader-panel").CascaderProps) | (() => import("element-plus/es/components/cascader-panel").CascaderProps) | ((new (...args: any[]) => import("element-plus/es/components/cascader-panel").CascaderProps) | (() => import("element-plus/es/components/cascader-panel").CascaderProps))[], unknown, unknown, () => import("element-plus/es/components/cascader-panel").CascaderProps, boolean>;
}>> & {
    "onUpdate:modelValue"?: ((value: CascaderValue | null | undefined) => any) | undefined;
    onChange?: ((value: CascaderValue | null | undefined) => any) | undefined;
    onFocus?: ((evt: FocusEvent) => any) | undefined;
    onBlur?: ((evt: FocusEvent) => any) | undefined;
    onClear?: (() => any) | undefined;
    onVisibleChange?: ((val: boolean) => any) | undefined;
    onExpandChange?: ((val: CascaderValue) => any) | undefined;
    onRemoveTag?: ((val: import("element-plus/es/components/cascader-panel").CascaderNodeValue | import("element-plus/es/components/cascader-panel").CascaderNodePathValue) => any) | undefined;
}, {
    disabled: boolean;
    separator: string;
    props: import("element-plus/es/components/cascader-panel").CascaderProps;
    placement: import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => "left" | "right" | "top" | "bottom" | "auto" | "auto-start" | "auto-end" | "top-start" | "top-end" | "bottom-start" | "bottom-end" | "right-start" | "right-end" | "left-start" | "left-end") | (() => import("element-plus/es/components/popper").Placement) | ((new (...args: any[]) => "left" | "right" | "top" | "bottom" | "auto" | "auto-start" | "auto-end" | "top-start" | "top-end" | "bottom-start" | "bottom-end" | "right-start" | "right-end" | "left-start" | "left-end") | (() => import("element-plus/es/components/popper").Placement))[], import("element-plus/es/components/popper").Placement, unknown>;
    options: import("element-plus/es/components/cascader-panel").CascaderOption[];
    effect: import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => string) | (() => import("element-plus/es/components/popper").PopperEffect) | ((new (...args: any[]) => string) | (() => import("element-plus/es/components/popper").PopperEffect))[], unknown, unknown>;
    valueOnClear: import("element-plus/es/utils").EpPropMergeType<(new (...args: any[]) => string | number | boolean | Function) | (() => string | number | boolean | Function | null) | ((new (...args: any[]) => string | number | boolean | Function) | (() => string | number | boolean | Function | null))[], unknown, unknown>;
    validateEvent: import("element-plus/es/utils").EpPropMergeType<BooleanConstructor, unknown, unknown>;
    clearable: boolean;
    fallbackPlacements: import("element-plus/es/components/popper").Placement[];
    teleported: import("element-plus/es/utils").EpPropMergeType<BooleanConstructor, unknown, unknown>;
    persistent: import("element-plus/es/utils").EpPropMergeType<BooleanConstructor, unknown, unknown>;
    debounce: number;
    checkOnClickNode: boolean;
    showPrefix: import("element-plus/es/utils").EpPropMergeType<BooleanConstructor, unknown, unknown>;
    filterable: boolean;
    filterMethod: (node: CascaderNode, keyword: string) => boolean;
    showAllLevels: import("element-plus/es/utils").EpPropMergeType<BooleanConstructor, unknown, unknown>;
    collapseTags: boolean;
    maxCollapseTags: number;
    collapseTagsTooltip: boolean;
    beforeFilter: (value: string) => boolean | Promise<any>;
    tagType: import("element-plus/es/utils").EpPropMergeType<StringConstructor, "primary" | "success" | "warning" | "info" | "danger", unknown>;
    tagEffect: import("element-plus/es/utils").EpPropMergeType<StringConstructor, "dark" | "light" | "plain", unknown>;
    showCheckedStrategy: string;
}>;
declare const _default: __VLS_WithTemplateSlots<typeof __VLS_component, ReturnType<typeof __VLS_template>>;
export default _default;
type __VLS_WithTemplateSlots<T, S> = T & {
    new (): {
        $slots: S;
    };
};
