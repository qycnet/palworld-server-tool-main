import { useState } from "react"
import { useTranslation, Trans } from 'react-i18next';
import { ChevronDown } from "lucide-react"

import { DeathPenaltyLabels, AllowConnectPlatformLabels, LogFormatTypeLabels } from "@/consts/dropdownLabels"
import { Button } from "@/components/ui/button"
import { Label } from "@/components/ui/label"
import {
    Command,
    CommandEmpty,
    CommandGroup,
    CommandItem,
    CommandList,
} from "@/components/ui/command"
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuGroup,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
} from "@/components/ui/tooltip"
import { I18nStr } from "@/i18n";

type Labels = typeof DeathPenaltyLabels | typeof AllowConnectPlatformLabels | typeof LogFormatTypeLabels;
export type LabelValue = Labels[number]['name'];
type Key =  'DeathPenalty' | 'AllowConnectPlatform' | 'LogFormatType';

/**
 * 从字典中获取指定键的值，如果键不存在则返回默认值
 *
 * @param dict 字典对象，类型为 Record<string, T>
 * @param key 要获取的键
 * @param defaultValue 如果键不存在时返回的默认值
 * @returns 返回指定键的值或默认值
 */
function get<T>(dict: Record<string, T>, key: string, defaultValue: T): T {
    return Object.prototype.hasOwnProperty.call(dict, key) ? dict[key] : defaultValue;
}

/**
 * 下拉选择组件
 *
 * @param props 组件属性
 * @param props.dKey 键值，用于选择标签的集合
 * @param props.label 当前选中的标签
 * @param props.onLabelChange 标签改变时的回调函数
 * @returns 下拉选择组件
 */
export function DropDown(props: {
    dKey: Key;
    label: LabelValue;
    onLabelChange: (label: string) => void;
}) {
    const { dKey, label, onLabelChange } = props;
    const labels = {
      DeathPenalty: DeathPenaltyLabels,
      AllowConnectPlatform: AllowConnectPlatformLabels,
      LogFormatType: LogFormatTypeLabels
    }[dKey] as Labels;
    const { t } = useTranslation();
    const [open, setOpen] = useState(false);
    const i18nLabelDesc = get(I18nStr.entry.description[dKey] as Record<string, string>, label, "");

    const labelDesc = t(i18nLabelDesc, {
        defaultValue: labels.find((l) => l.name === label)?.desc ?? "",
    });


    return (
        <div className="space-y-1">
            <Label>
                <TooltipProvider>
                    <Tooltip>
                        <TooltipTrigger className="cursor-default">
                            <Trans i18nKey={get(I18nStr.entry.name, dKey, "")} />
                            <TooltipContent>
                                <p>{dKey}</p>
                            </TooltipContent>
                        </TooltipTrigger>
                    </Tooltip>
                </TooltipProvider>
            </Label>
            <div className="flex w-full items-start justify-between rounded-md border px-4 py-1 sm:flex-row sm:items-center">
                <p className="text-sm font-medium leading-none">
                    <span className="mr-2 rounded-lg bg-primary px-2 py-1 text-xs text-primary-foreground normal-case">
                        {label}
                    </span>
                    <span className="text-muted-foreground leading-7">{labelDesc}</span>
                </p>
                <DropdownMenu open={open} onOpenChange={setOpen}>
                    <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="sm" className="ml-auto">
                            <ChevronDown />
                        </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end" className="w-[200px]">
                        <DropdownMenuGroup>
                            <Command>
                                <CommandList>
                                    <CommandEmpty>No label found.</CommandEmpty>
                                    <CommandGroup>
                                        {labels.map((label) => (
                                            <CommandItem
                                                key={label.name}
                                                value={label.name}
                                                onSelect={() => {
                                                    setOpen(false);
                                                    onLabelChange(label.name);
                                                }}
                                            >
                                                {label.name}
                                            </CommandItem>
                                        ))}
                                    </CommandGroup>
                                </CommandList>
                            </Command>
                        </DropdownMenuGroup>
                    </DropdownMenuContent>
                </DropdownMenu>
            </div>
        </div>
    );
}
