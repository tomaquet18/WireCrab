import React, { useEffect, useState } from "react";
import { ScrollArea } from "@/components/ui/scroll-area";
import {
    ResizableHandle,
    ResizablePanel,
    ResizablePanelGroup
} from "@/components/ui/resizable"
import { tshark } from "wailsjs/go/models";
import { parse } from "postcss";

interface PacketDetailsPanelProps {
    protocolInfo?: tshark.ProtocolInfo;
    hexDump?: string;
}

interface ProtocolNodeProps {
    protocolInfo: tshark.ProtocolInfo;
    onFieldSelect: (pos: string, size: string) => void;
}

const ProtocolNode: React.FC<ProtocolNodeProps> = ({ protocolInfo, onFieldSelect }) => {
    return (
        <div className="mb-4">
            <h3 className="text-lg font-semibold mb-2">{protocolInfo.Name}</h3>
            {protocolInfo.Detail && (
                <div className="pl-4">
                    {Object.entries(protocolInfo.Detail).map(([key, detail]: [string, any]) => {
                        if (!detail.value) return null;
                        return (
                            <div
                                key={key}
                                className="text-sm mb-1"
                                onClick={() => onFieldSelect(detail.pos, detail.size)}
                            >
                                <span className="font-medium">{key}: </span>
                                <span className="text-muted-foreground">{detail.value}</span>
                            </div>
                        );
                    })}
                </div>
            )}
            {protocolInfo.Child && (
                <div className="pl-4 mt-2 border-l">
                    <ProtocolNode protocolInfo={protocolInfo.Child} onFieldSelect={onFieldSelect} />
                </div>
            )}
        </div>
    );
};

interface ParsedHexLine {
    offset: string;
    hexBytes: {
        value: string;
        position: number;
    }[];
    ascii: {
        char: string;
        position: number;
    }[];
}

const parseHexDump = (data: string): ParsedHexLine[] => {
    return data.split('\n').map((line, lineIndex) => {
        if (!line.trim()) return null;

        // Extract offset (0000, 0010, etc.)
        const offset = line.substring(0, 4);

        // Extract hex values
        const hexPart = line.substring(6, 54).trim();
        const hexBytes = hexPart.split(' ')
            .filter(x => x)
            .map((byte, idx) => ({
                value: byte,
                position: lineIndex * 16 + idx
            }));

        // Extract ASCII part (after the last two spaces)
        const asciiPart = line.split('  ').pop()?.trim() || '';
        const ascii = Array.from(asciiPart).map((char, idx) => ({
            char,
            position: lineIndex * 16 + idx
        }));

        return { offset, hexBytes, ascii };
    }).filter((line): line is ParsedHexLine => line !== null);
};

export const HexDump: React.FC<{
    data: string;
    highlightStart?: number;
    highlightSize?: number;
    onByteSelect: (position: number, size: number) => void;
}> = ({
    data,
    highlightStart,
    highlightSize,
    onByteSelect
}) => {
    const [selectedStart, setSelectedStart] = useState<number | null>(null);
    const [isSelecting, setIsSelecting] = useState(false);
    const parsedLines = parseHexDump(data);

    const isHighlighted = (position: number) => {
        if (highlightStart !== undefined && highlightSize !== undefined) {
            return position >= highlightStart && position < (highlightStart + highlightSize);
        }
        return false;
    };

    const handleMouseDown = (position: number) => {
        setSelectedStart(position);
        setIsSelecting(true);
        onByteSelect(position, 1);
    };

    const handleMouseMove = (position: number) => {
        if (isSelecting && selectedStart !== null) {
            const start = Math.min(selectedStart, position);
            const size = Math.abs(position - selectedStart) + 1;
            onByteSelect(start, size);
        }
    };

    const handleMouseUp = () => {
        setIsSelecting(false);
    }

    // Global mouse up hundler
    useEffect(() => {
        const handleGlobalMouseUp = () => {
            setIsSelecting(false);
        };

        window.addEventListener('mouseup', handleGlobalMouseUp);
        return () => {
            window.removeEventListener('mouseup', handleGlobalMouseUp);
        };
    }, []);

    return (
        <div className="font-mono text-sm select-none whitespace-nowrap">
            {/* Header row */}
            <div className="flex mb-2">
                <div className="w-16 text-muted-foreground">Offset</div>
                <div className="flex">
                    {Array.from({ length: 16 }, (_, i) => (
                        <div key={i} className="w-5 text-muted-foreground text-center">
                            {i.toString(16).padStart(2, '0').toUpperCase()}
                        </div>
                    ))}
                </div>
                <div className="ml-4 text-muted-foreground">ASCII</div>
            </div>

            {/* Data rows */}
            {parsedLines.map((line, lineIdx) => (
                <div key={lineIdx} className="flex hover:bg-accent/20 py-0.5">
                    {/* Offset column */}
                    <div className="w-16 text-muted-foreground">
                        {line.offset}
                    </div>

                    {/* Hex values */}
                    <div className="flex">
                        {line.hexBytes.map(({ value, position }) => (
                            <div
                                key={position}
                                className={`w-5 text-center cursor-pointer select-none
                                    ${isHighlighted(position) ? 'bg-blue-500/30' : ''}
                                    ${selectedStart !== null ? 'hover:bg-blue-500/20' : ''}`}
                                onMouseDown={() => handleMouseDown(position)}
                                onMouseEnter={(e) => handleMouseMove(position)}
                            >
                                {value}
                            </div>
                        ))}

                        {/* Padding for incomplete rows */}
                        {line.hexBytes.length < 16 && (
                            <div style={{ width: `${(16 - line.hexBytes.length) * 1.25}rem` }} />
                        )}
                    </div>

                    {/* ASCII representation */}
                    <div className="ml-4 font-mono">
                        {line.ascii.map(({ char, position }) => (
                            <span
                                key={position}
                                className={`cursor-pointer select-none
                                    ${isHighlighted(position) ? 'bg-blue-500/30' : ''}
                                    ${selectedStart !== null ? 'hover:bg-blue-500/20' : ''}`}
                                onMouseDown={() => handleMouseDown(position)}
                                onMouseEnter={(e) => handleMouseMove(position)}
                            >
                                {char}
                            </span>
                        ))}
                    </div>
                </div>
            ))}
        </div>
    );
};

export const PacketDetailsPanel: React.FC<PacketDetailsPanelProps> = ({
    protocolInfo,
    hexDump
}) => {
    const [selectedPos, setSelectedPos] = useState<string>();
    const [selectedSize, setSelectedSize] = useState<string>();

    const handleFieldSelect = (pos: string, size: string) => {
        setSelectedPos(pos);
        setSelectedSize(size);
    }

    const handleByteSelect = (position: number, size: number) => {
        setSelectedPos(position.toString());
        setSelectedSize(size.toString());
    }

    return (
        <ResizablePanelGroup direction="horizontal">
            <ResizablePanel defaultSize={50}>
                <ScrollArea className="h-full">
                    {protocolInfo ? (
                        <ProtocolNode
                            protocolInfo={protocolInfo}
                            onFieldSelect={handleFieldSelect}
                        />
                    ) : (
                        <div className="text-center text-muted-foreground">
                            No packet selected
                        </div>
                    )}
                </ScrollArea>
            </ResizablePanel>

            <ResizableHandle />

            <ResizablePanel defaultSize={50}>
                <ScrollArea className="h-full">
                    {hexDump ? (
                        <HexDump
                            data={hexDump}
                            highlightStart={parseInt(selectedPos!)}
                            highlightSize={parseInt(selectedSize!)}
                            onByteSelect={handleByteSelect}
                        />
                    ) : (
                        <div className="text-center text-muted-foreground">
                            No hex dump available
                        </div>
                    )}
                </ScrollArea>
            </ResizablePanel>
        </ResizablePanelGroup>
    );
};