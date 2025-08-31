import {Button} from "@/components/ui/button.tsx";
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger
} from "@/components/ui/drawer.tsx";
import {saveFileTableCols} from "@/pages/save-file-columns.tsx";
import {SaveFileTable} from "@/pages/save-file-table.tsx";
import {useSaveTableStore} from "@/store/save-collection.ts";

export default function PackDataSidebar() {

  const data = useSaveTableStore(state => state.packDataSavefiles)

  return (
    <Drawer>
      <DrawerTrigger>Open PackData</DrawerTrigger>
      <DrawerContent>
        <DrawerHeader>
          <DrawerTitle>Pack Data</DrawerTitle>
          <DrawerDescription>ooh yeah</DrawerDescription>
        </DrawerHeader>
        <SaveFileTable columns={saveFileTableCols} data={data}/>
        <DrawerFooter>
          <Button>Submit</Button>
          <DrawerClose>
            Cancel
          </DrawerClose>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  )
}
