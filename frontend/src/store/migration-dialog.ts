import {create} from "zustand/react";

export type MigrationDestination = "PACK_DATA" | "SAVE_DIRECTORY"

export type MigrationDialogState = {
  isOpen: boolean
  fileLocation: string
  destination: MigrationDestination
}

export type MigrationDialogActions = {
  openWith: (fileLocation: string, destination: MigrationDestination) => void
  close: () => void
}

type MigrationDialogStore = MigrationDialogState & MigrationDialogActions

export const useMigrationDialogStore = create<MigrationDialogStore>((set) => ({
  isOpen: false,
  fileLocation: "",
  destination: "SAVE_DIRECTORY",
  openWith: (fileLocation: string, destination: MigrationDestination) => {
    set(() => ({
      isOpen: true,
      fileLocation: fileLocation,
      destination
    }))
  },
  close: () => set(() => ({isOpen: false, fileLocation: ""})),
}))