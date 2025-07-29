import React, { useState } from "react";
import { ScrollArea } from "@/components/ui/scroll-area";
import {
    ResizableHandle,
    ResizablePanel,
    ResizablePanelGroup
} from "@/components/ui/resizable"
import { tshark } from "wailsjs/go/models";

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

export const HexDump: React.FC<{ data: string; highlightStart?: number; highlightSize?: number }> = ({
    data,
    highlightStart,
    highlightSize
}) => {
    return (
        <pre className="font-mono text-sm whitespace-pre-wrap">
            {data.split('\n').map((line, idx) => (
                <div key={idx} className="hover:bg-accent">
                    {line}
                </div>
            ))}
        </pre>
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