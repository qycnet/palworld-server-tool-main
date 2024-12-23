import { Label } from "@/components/ui/label"
import { Input } from "@/components/ui/input"
import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
} from "@/components/ui/tooltip"

/**
 * 创建一个文本输入框组件
 *
 * @param props 组件属性
 * @param props.id 输入框的id
 * @param props.name 输入框的名称
 * @param props.value 输入框的初始值
 * @param props.onChange 输入框内容变化时的回调函数
 * @param props.type 输入框的类型（可选），默认为'text'
 * @param props.disabled 是否禁用输入框（可选），默认为false
 * @param props.multiline 是否为多行输入框（可选），默认为false
 * @returns 返回创建的输入框组件
 */
function TextInput(props: {
    id: string;
    name: string;
    value: string;
    onChange: React.ChangeEventHandler<HTMLInputElement | HTMLTextAreaElement>;
    type?: string;
    disabled?: boolean;
    multiline?: boolean;
}) {
    const {
        id,
        name,
        type,
        value,
        onChange,
        disabled,
        multiline,
    } = props;
    const inputElement = multiline ? (
        <textarea
            id={id}
            value={value}
            onChange={onChange}
            className="w-[98%] h-[100px] p-2 border border-gray-300 rounded-md"
            disabled={disabled}
        />
    ): (
        <Input
            value={value}
            id={id}
            onChange={onChange}
            type={type}
            className="w-[98%]"
            disabled={disabled}
        />
    )
    return (
        <div className="flex flex-col items-center space-y-2">
            <TooltipProvider>
                <Tooltip>
                    <TooltipTrigger className="cursor-default mr-auto">
                        <Label>{name}</Label>
                        <TooltipContent>
                            <p>{id}</p>
                        </TooltipContent>
                    </TooltipTrigger>
                </Tooltip>
            </TooltipProvider>
            {inputElement}            
        </div>
    );
}

export { TextInput };