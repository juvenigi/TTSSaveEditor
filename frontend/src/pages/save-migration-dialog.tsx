import {Button} from "@/components/ui/button";
import {Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle} from "@/components/ui/dialog.tsx";
import {useMigrationDialogStore} from "@/store/migration-dialog.ts";
import {useState} from "react";
import {GetProperties, GetTabletopSaves, WriteToPackData} from "../../wailsjs/go/internal/CacheManagerApi";
import {Input} from "@/components/ui/input.tsx";
import {useSaveTableStore} from "@/store/save-collection.ts";


export default function SaveMigrationDialog() {
  const dest = useMigrationDialogStore(state => state.destination)
  const isOpen = useMigrationDialogStore(state => state.isOpen)
  const close = useMigrationDialogStore(state => state.close)
  const filename = useMigrationDialogStore(state => state.fileLocation)

  const [saveName, setSaveName] = useState("Untitled")
  const [processing, setProcessing] = useState(false)

  const handleConfirm = async () => {
    const properties = await GetProperties();
    try {
      await WriteToPackData(filename, saveName)
    } finally {
      close()
      setProcessing(false)
      useSaveTableStore.getState().clearPackDataFiles()
      GetTabletopSaves(properties.PackDataDir).finally()
    }
  }

  return (
    <Dialog open={isOpen} onOpenChange={close}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Are you absolutely sure?</DialogTitle>
          <DialogDescription>{filename}
            <Input value={saveName} onChange={e => setSaveName(e.target.value)}></Input>
          </DialogDescription>
          <Button disabled={processing} onMouseDown={handleConfirm}>Yarr to {dest} we go</Button>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  )
}