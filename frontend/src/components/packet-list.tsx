import React from "react";
import { TableBody, TableCell, TableRow } from "@/components/ui/table";
import { useVirtualizer } from '@tanstack/react-virtual';
import { types } from "../../wailsjs/go/models"

type CapturedPacket = types.CapturedPacket;

interface PacketListProps {
    packets: CapturedPacket[];
    onPacketClick: (packet: CapturedPacket) => void;
}

export function VirtualizedPacketList({ packets, onPacketClick }: PacketListProps) {
    const parentRef = React.useRef<HTMLDivElement>(null);

    const virtualizer = useVirtualizer({
        count: packets.length,
        getScrollElement: () => parentRef.current,
        estimateSize: () => 40, // Estimated row height
        overscan: 5,
    });

    return (
        <div
            ref={parentRef}
            className="max-h-[calc(100vh-200px)] overflow-auto"
        >
            {/* The large inner element to hold all items */}
            <div
                style={{
                    height: `${virtualizer.getTotalSize()}px`,
                    width: '100%',
                    position: 'relative',
                }}
            >
                {/* Only the visible items in the viewport */}
                {virtualizer.getVirtualItems().map((virtualItem) => {
                    const packet = packets[virtualItem.index];
                    return (
                        <div
                            key={virtualItem.key}
                            data-index={virtualItem.index}
                            ref={virtualizer.measureElement}
                            className={`absolute top-0 left-0 w-full flex hover:bg-accent cursor-pointer border-b`}
                            style={{
                                transform: `translateY(${virtualItem.start}px)`,
                            }}
                            onClick={() => onPacketClick(packet)}
                        >
                            <div style={{ width: '5%' }} className="p-2 text-right">{virtualItem.index + 1}</div>
                            <div style={{ width: '10%' }} className="p-2 truncate">{packet.meta.Timestamp}</div>
                            <div style={{ width: '20%' }} className="p-2 truncate">{packet.meta.SrcIP}</div>
                            <div style={{ width: '20%' }} className="p-2 truncate">{packet.meta.DstIP}</div>
                            <div style={{ width: '10%' }} className="p-2 truncate">{packet.meta.Protocol}</div>
                            <div style={{ width: '10%' }} className="p-2 text-right">{packet.meta.Length}</div>
                            <div style={{ width: '25%' }} className="p-2 truncate">{packet.meta.Info || "-"}</div>
                        </div>
                    );
                })}
            </div>
        </div>
    );
}