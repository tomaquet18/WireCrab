import { Suspense, useCallback, useEffect, useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Play, StopCircle } from "lucide-react"
import { GetCapturedPackets, StartCapture, GetPacketCount, GetPacketDetails } from "../../wailsjs/go/main/App"
import { types } from "../../wailsjs/go/models"
import { useInterfaceStore } from "@/stores/interfaces"
import { PacketDetailsPanel } from "@/components/packet-details-panel";
import { ResizablePanelGroup, ResizablePanel, ResizableHandle } from "@/components/ui/resizable"
import { VirtualTable, Column } from "@/components/virtual-table"

type CapturedPacket = types.CapturedPacket;

export default function CapturePage() {
  const [isCapturing, setIsCapturing] = useState(false)
  const [packets, setPackets] = useState<CapturedPacket[]>([])
  const [totalPackets, setTotalPackets] = useState(0)
  const [loading, setLoading] = useState(false)
  const [selectedPacketDetails, setSelectedPacketDetails] = useState<any>(null);
  const [hexDump, setHexDump] = useState<string>("");
  const [isStarted, setIsStarted] = useState(false);
  const interfaceName = useInterfaceStore((s) => s.selected)

  const columns: Column<CapturedPacket>[] = [
    {
      header: "No.",
      accessorFn: (row, index) => index + 1,
      width: 60,
      minWidth: 50,
    },
    {
      header: "Time",
      accessorFn: (row) => row.meta.Timestamp,
      width: 120,
      minWidth: 100,
    },
    {
      header: "Source",
      accessorFn: (row) => row.meta.SrcIP || "",
      width: 160,
      minWidth: 120,
    },
    {
      header: "Destination",
      accessorFn: (row) => row.meta.DstIP || "",
      width: 160,
      minWidth: 120,
    },
    {
      header: "Protocol",
      accessorFn: (row) => row.meta.Protocol || "",
      width: 100,
      minWidth: 80,
    },
    {
      header: "Length",
      accessorFn: (row) => row.meta.Length || 0,
      width: 80,
      minWidth: 60,
    },
    {
      header: "Info",
      accessorFn: (row) => row.meta.Info || "-",
      width: 300,
      minWidth: 150,
    },
  ];


  // Memoize packet loading function
  const loadMorePackets = useCallback(async () => {
    if (loading) return;
    setLoading(true);

    try {
      // Get current packet count
      const count = await GetPacketCount();
      setTotalPackets(count);

      // If there are new packets to load
      if (count > packets.length) {
        const newPackets = await GetCapturedPackets(packets.length, count - packets.length);
        if (newPackets && newPackets.length > 0) {
          setPackets(prev => [...prev, ...newPackets]);
        }
      }
    } finally {
      setLoading(false);
    }
  }, [loading, packets.length]);

  useEffect(() => {
    let intervalId: number | undefined;

    const startCapture = async () => {
      if (!isStarted && interfaceName) {
        try {
          // Start the capture
          await StartCapture(interfaceName);
          setIsStarted(true);
        } catch (error) {
          console.error('Error startig capture:', error);
          setIsCapturing(false);
        }
      }
    };
    
    if (isCapturing) {
      startCapture();
      intervalId = window.setInterval(loadMorePackets, 1000);
    } else {
      setIsStarted(false);
    }
    
    return () => {
      if (intervalId) {
        clearInterval(intervalId);
        intervalId = undefined;
      }
    }
  }, [isCapturing, interfaceName, isStarted]);

  const handleStart = useCallback(() => {
    setPackets([]); // Clear existing packets when starting new capture
    setTotalPackets(0);
    setIsStarted(false);
    setIsCapturing(true);
  }, []);

  const handleStop = useCallback(() => {
    setIsCapturing(false);
  }, []);
  
  // Memoize packet click handler
  const handlePacketClick = useCallback(async (packet: CapturedPacket) => {
    const frameNumber = packet.parsed?.Detail?.["frame.number"]?.value;
    if (!frameNumber) return;

    try {
      const details = await GetPacketDetails(parseInt(frameNumber));
      setSelectedPacketDetails(details.Info);
      setHexDump(details.HexDump);
    } catch (error) {
      console.error('Failed to get packet details: ', error);
    }
  }, []);


  return (
    <div className="h-screen w-full flex flex-col">
      {/* Toolbar */}
      <div className="px-6 py-3 border-b flex items-center gap-4 shrink-0">
        <Button onClick={handleStart} disabled={isCapturing}>
          <Play className="w-4 h-4 mr-2" />
          Start
        </Button>
        <Button
          variant="destructive"
          onClick={handleStop}
          disabled={!isCapturing}
        >
          <StopCircle className="w-4 h-4 mr-2" />
          Stop
        </Button>
        <Input
          placeholder="Filter (e.g., tcp.port == 443)"
          className="w-64 ml-auto"
        />
      </div>

      {/* Main content area */}
      <div className="flex-1 overflow-hidden">
        <ResizablePanelGroup direction="vertical" className="h-full" autoSaveId="capture-page-layout">
          {/* Packet List */}
          <ResizablePanel defaultSize={75} minSize={30} maxSize={85} style={{ overflow: 'hidden' }}>
            <div className="h-full overflow-auto relative">
              {/* Table */}
              <div className="w-full h-full relative">
                <VirtualTable
                  data={packets}
                  columns={columns}
                  onRowClick={handlePacketClick}
                />
              </div>
            </div>
          </ ResizablePanel>

          <ResizableHandle className="h-2 bg-border hover:bg-foreground/10 transition-colors" />

          {/* Details Panel */}
          <ResizablePanel defaultSize={25} minSize={15} maxSize={70}>
            <div className="h-full w-full">
              <Suspense fallback={<div>Loading details...</div>}>
                <PacketDetailsPanel
                  protocolInfo={selectedPacketDetails}
                  hexDump={hexDump}
                />
              </Suspense>

              {/* Packet count indicator */}
              <div className="px-4 py-2 border-t">
                Total Packets: {totalPackets}
              </div>
            </div>
          </ResizablePanel>
        </ResizablePanelGroup>  
      </div>
    </div>
  )
}
