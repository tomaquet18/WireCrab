import * as React from "react";
import { useVirtualizer } from "@tanstack/react-virtual";
import { Table } from "@/components/ui/table";
import { ResizableBox } from "react-resizable"
import "react-resizable/css/styles.css";
import { types } from "../../wailsjs/go/models";

type CapturedPacket = types.CapturedPacket;

export interface Column<T> {
    header: string;
    accessorFn: (row: T, index: number) => string | number;
    width: number;
    minWidth?: number;
}

interface VirtualTableProps<T> {
    data: T[];
    columns: Column<T>[];
    onRowClick?: (row: T) => void;
    selectedPacketId?: number;
}

function getPacketRowColor(packet: CapturedPacket) {
    const protocol = packet.meta.Protocol?.toLowerCase() || '';
    const info = packet.meta.Info?.toLowerCase() || '';

    // Protocol-based coloring
    switch(protocol) {
        case 'tcp':
            return 'bg-purple-100 hover:bg-purple-200';
        case 'udp':
            return 'bg-blue-100 hover:bg-blue-200';
        case 'http':
            return 'bg-green-100 hover:bg-green-200';
        case 'smb':
        case 'netbios':
            return 'bg-yellow-100 hover:bg-yellow-200';
        case 'ospf':
        case 'bgp':
        case 'rip':
            return 'bg-yellow-200 hover:bg-yellow-300';
        default:
            return 'hover:bg-muted/50';
    }
}

export function VirtualTable<T>({ data, columns, onRowClick, selectedPacketId }: VirtualTableProps<T>) {
    const parentRef = React.useRef<HTMLDivElement>(null);
    const [columnWidths, setColumnWidths] = React.useState<number[]>(
        columns.map((col) => col.width)
    );

    const virtualizer = useVirtualizer({
        count: data.length,
        getScrollElement: () => parentRef.current,
        estimateSize: () => 40,
        overscan: 5,
    });

    const onResize = (index: number) => (e: any, { size }: { size: { width: number } }) => {
        const newColumnWidths = [...columnWidths];
        newColumnWidths[index] = size.width;
        setColumnWidths(newColumnWidths);
    };

    return (
        <div className="w-full h-full">
            <Table>
                {/* Sticky Header */}
                <div className="sticky top-0 z-10 bg-background border-b">
                    <div className="flex">
                        {columns.map((column, columnIndex) => (
                            <ResizableBox
                                key={columnIndex}
                                width={columnWidths[columnIndex]}
                                height={40}
                                minConstraints={[column.minWidth || 50, 40]}
                                onResize={onResize(columnIndex)}
                                axis="x"
                                handle={columnIndex < columns.length - 1 ? (
                                    <div className="absolute right-0 top-0 bottom-0 w-px bg-border hover:bg-foreground/50 cursor-col-resize
                                                    transition-colors duration-150 ease-in-out
                                                    after:content-[''] after:absolute after:right-[-4px] after:top-0 after:bottom-0 after:w-[8px]"
                                    />
                                ): <div />}
                            >
                                <div className="px-4 py-2 font-medium">
                                    {column.header}
                                </div>
                            </ResizableBox>
                        ))}
                    </div>
                </div>

                {/* Virtualized Body */}
                <div
                    ref={parentRef}
                    className="overflow-auto"
                    style={{ height: 'calc(100% - 40px)' }}
                >
                    <div
                        style={{
                            height: `${virtualizer.getTotalSize()}px`,
                            width: '100%',
                            position: "relative",
                        }}
                    >
                        {virtualizer.getVirtualItems().map((virtualRow) => {
                            const row = data[virtualRow.index];
                            const rowColorClass = getPacketRowColor(row as CapturedPacket);
                            const isSelected = (row as CapturedPacket).parsed?.Detail?.["frame.number"]?.value === selectedPacketId?.toString();

                            return (
                                <div
                                    key={virtualRow.index}
                                    className={`absolute w-full flex border-b hover:bg-muted/50 cursor-pointer transition-colors
                                        ${rowColorClass}
                                        ${isSelected ? 'ring-2 ring-blue-200 ring-inset bg-blue-100/50' : ''}`}
                                    style={{
                                        height: `${virtualRow.size}px`,
                                        transform: `translateY(${virtualRow.start}px)`,
                                    }}
                                    onClick={() => onRowClick?.(row)}
                                >
                                    {columns.map((column, columnIndex) => (
                                        <div
                                            key={columnIndex}
                                            style={{ 
                                                width: columnWidths[columnIndex],
                                                ...(columnIndex === columns.length - 1 && {
                                                    flex: '1',
                                                    minWidth: column.minWidth || 50,
                                                })
                                            }}
                                            className="px-4 py-2 truncate"
                                        >
                                            {String(column.accessorFn(row, virtualRow.index))}
                                        </div>
                                    ))}
                                </div>
                            );
                        })}
                    </div>
                </div>
            </Table>
        </div>
    );
}