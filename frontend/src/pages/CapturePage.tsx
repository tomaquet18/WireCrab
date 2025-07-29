import { useEffect, useRef, useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from "@/components/ui/table"
import { Play, StopCircle } from "lucide-react"
import { GetCapturedPackets, StartCapture, GetPacketCount } from "../../wailsjs/go/main/App"
import { types } from "../../wailsjs/go/models"
import { useInterfaceStore } from "@/stores/interfaces"

type CapturedPacket = types.CapturedPacket

export default function CapturePage() {
  const [isCapturing, setIsCapturing] = useState(false)
  const [packets, setPackets] = useState<CapturedPacket[]>([])
  const [totalPackets, setTotalPackets] = useState(0)
  const [loading, setLoading] = useState(false)
  const pageSize = 100
  const loadMoreRef = useRef(null)
  const interfaceName = useInterfaceStore((s) => s.selected)

  const loadMorePackets = async () => {
    if (loading) return
    setLoading(true)

    try {
      const newPackets = await GetCapturedPackets(packets.length, pageSize)
      if (newPackets && newPackets.length > 0) {
        setPackets(prev => [...prev, ...newPackets])
      }
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    const observer = new IntersectionObserver(
      entries => {
        if (entries[0].isIntersecting) {
          loadMorePackets()
        }
      },
      { threshold: 0.5 }
    )

    if (loadMoreRef.current) {
      observer.observe(loadMoreRef.current)
    }

    return () => observer.disconnect()
  }, [packets])

  useEffect(() => {
    let interval: ReturnType<typeof setInterval>

    if (isCapturing && interfaceName) {
      StartCapture(interfaceName)

      interval = setInterval(async () => {
        const count = await GetPacketCount()
        setTotalPackets(count)

        // Load new packets only if we're near the end
        if (packets.length >= totalPackets - pageSize) {
          loadMorePackets()
        }
      }, 1000)
    }

    return () => {
      clearInterval(interval)
    }
  }, [isCapturing, interfaceName])

  const handleStart = () => setIsCapturing(true)
  const handleStop = () => setIsCapturing(false)

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

      {/* Scrollable table area with sticky header */}
      <div className="flex-1 overflow-hidden">
        <div className="h-full overflow-y-auto relative">
          <Table className="w-full table-auto">
            <TableHeader className="sticky top-0 z-10 bg-background shadow-sm">
              <TableRow>
                <TableHead className="w-12 text-right">No.</TableHead>
                <TableHead className="min-w-[80px] w-[100px]">Time</TableHead>
                <TableHead className="min-w-[140px]">Source</TableHead>
                <TableHead className="min-w-[140px]">Destination</TableHead>
                <TableHead className="w-[90px]">Protocol</TableHead>
                <TableHead className="w-[70px] text-right">Length</TableHead>
                <TableHead className="min-w-[200px]">Info</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {packets.map((pkt, idx) => (
                <TableRow key={idx} className="hover:bg-accent cursor-pointer">
                  <TableCell className="text-right">{idx + 1}</TableCell>
                  <TableCell>{pkt.meta.Timestamp}</TableCell>
                  <TableCell>{pkt.meta.SrcIP}</TableCell>
                  <TableCell>{pkt.meta.DstIP}</TableCell>
                  <TableCell>{pkt.meta.Protocol}</TableCell>
                  <TableCell className="text-right">{pkt.meta.Length}</TableCell>
                  <TableCell>{pkt.meta.Info || "-"}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>

        

        {/* Packet count indicator */}
        <div className="px-4 py-2 border-t">
          Total Packets: {totalPackets}
        </div>
      </div>
    </div>
  )
}
