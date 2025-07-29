import React from "react";
import { ScrollArea } from "@/components/ui/scroll-area";
import {
    Sheet,
    SheetContent,
    SheetHeader,
    SheetTitle,
} from "@/components/ui/sheet"
import { tshark } from "wailsjs/go/models";

interface PacketDetailsDrawerProps {
    isOpen: boolean;
    onClose: () => void;
    protocolInfo?: tshark.ProtocolInfo;
}

interface ProtocolNodeProps {
    protocolInfo: tshark.ProtocolInfo;
}

const ProtocolNode: React.FC<ProtocolNodeProps> = ({ protocolInfo }) => {
    return (
        <div className="mb-4">
            <h3 className="text-lg font-semibold mb-2">{protocolInfo.Name}</h3>
            {protocolInfo.Detail && (
                <div className="pl-4">
                    {Object.entries(protocolInfo.Detail).map(([key, detail]) => {
                        const value = (detail as any)?.value;
                        if (!value) return null;
                        return (
                            <div key={key} className="text-sm mb-1">
                                <span className="font-medium">{key}</span>
                                <span className="text-muted-foreground">{value}</span>
                            </div>
                        );
                    })}
                </div>
            )}
            {protocolInfo.Child && (
                <div className="pl-4 mt-2 border-l">
                    <ProtocolNode protocolInfo={protocolInfo.Child} />
                </div>
            )}
        </div>
    );
};

export const PacketDetailsDrawer: React.FC<PacketDetailsDrawerProps> = ({
    isOpen,
    onClose,
    protocolInfo
}) => {
    return (
        <Sheet open={isOpen} onOpenChange={onClose}>
            <SheetContent side="right" className="w-[400px] sm:w-[540px]">
                <SheetHeader>
                    <SheetTitle>Packet Details</SheetTitle>
                </SheetHeader>
                <ScrollArea className="h-[calc(100vh-80px)] mt-4">
                    {protocolInfo ? (
                        <>
                            <div className="mb-2 text-sm text-muted-foreground">
                                Protocol: {protocolInfo.Name}
                            </div>
                            <ProtocolNode protocolInfo={protocolInfo} />
                        </>
                    ) : (
                        <div className="text-center text-muted-foreground">
                            No packet selected or failed to load details
                        </div>
                    )}
                </ScrollArea>
            </SheetContent>
        </Sheet>
    );
};