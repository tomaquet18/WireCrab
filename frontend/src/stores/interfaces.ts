import { create } from "zustand"

interface InterfaceStore {
  selected: string
  setSelected: (name: string) => void
}

export const useInterfaceStore = create<InterfaceStore>((set) => ({
  selected: "",
  setSelected: (name: string) => set({ selected: name }),
}))
