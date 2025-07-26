import { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"
import { GetDevices } from "../../wailsjs/go/main/App"
import { Button } from "@/components/ui/button"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Skeleton } from "@/components/ui/skeleton"
import { useInterfaceStore } from "@/stores/interfaces"

interface Device {
  name: string
  description: string
  type: string
}

export default function WelcomePage() {
  const [interfaces, setInterfaces] = useState<Device[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selected, setSelected] = useState<Device | null>(null)
  const navigate = useNavigate()

  const setSelectedInterface = useInterfaceStore((s) => s.setSelected)

  useEffect(() => {
    GetDevices()
      .then(setInterfaces)
      .catch((err) => setError(String(err)))
      .finally(() => setLoading(false))
  }, [])

  const handleContinue = () => {
    if (selected) {
      setSelectedInterface(selected.name)
      navigate("/capture")
    }
  }

  return (
    <div className="min-h-screen flex flex-col items-center justify-center px-6 py-10 space-y-6">
      <div className="text-center space-y-1">
        <h1 className="text-3xl font-bold">Welcome to WireCrab</h1>
        <p className="text-muted-foreground text-sm">Select an interface to start capturing packets</p>
      </div>

      {loading ? (
        <Skeleton className="h-10 w-64" />
      ) : (
        <Select
          onValueChange={(value) => {
            const found = interfaces.find((i) => i.name === value)
            setSelected(found ?? null)
          }}
        >
          <SelectTrigger className="w-64">
            <SelectValue placeholder="Choose interface..." />
          </SelectTrigger>
          <SelectContent>
            {interfaces.map((iface) => (
              <SelectItem key={iface.name} value={iface.name}>
                {iface.description || iface.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      )}

      {selected && (
        <div className="text-sm bg-muted/50 border rounded-lg p-4 w-64 space-y-1 text-left">
          <p><strong>Name:</strong> {selected.name}</p>
          <p><strong>Description:</strong> {selected.description || "-"}</p>
          <p><strong>Type:</strong> {selected.type}</p>
        </div>
      )}

      <Button onClick={handleContinue} disabled={!selected}>
        Start scan
      </Button>

      {error && <p className="text-red-500 text-sm">{error}</p>}
    </div>
  )
}