import { Label } from "@/components/ui/label"
import { Switch } from "@/components/ui/switch"
import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
} from "@/components/ui/tooltip"

/**
 * SwitchInput 组件用于创建一个带有开关的输入控件。
 *
 * @param props - 组件的属性对象
 * @param props.id - 开关的唯一标识符。
 * @param props.name - 开关的名称。
 * @param props.checked - 开关的当前状态（true 表示选中，false 表示未选中）。
 * @param props.onCheckedChange - 当开关状态改变时触发的回调函数，参数为新的状态值（true 或 false）。
 * @param props.disabled - 可选参数，指定开关是否禁用（true 表示禁用，false 表示启用，默认为 false）。
 * @returns 返回 JSX 元素，包含 Tooltip 和 Switch 组件。
 */
export function SwitchInput(props: {
    id: string;
    name: string;
    checked: boolean;
    onCheckedChange: (checked: boolean) => void;
    disabled?: boolean;
}) {
    const {
        id,
        name,
        checked,
        onCheckedChange,
        disabled,
    } = props;
    return (
        <div className="flex">
            <TooltipProvider>
                <Tooltip>
                    <TooltipTrigger className="cursor-default leading-6">
                        <Label>{name}</Label>
                        <TooltipContent>
                            <p>{id}</p>
                        </TooltipContent>
                    </TooltipTrigger>
                </Tooltip>
            </TooltipProvider>
            <Switch className="ml-auto" checked={checked} id={id} onCheckedChange={onCheckedChange} disabled={disabled} />
        </div>
    );
}