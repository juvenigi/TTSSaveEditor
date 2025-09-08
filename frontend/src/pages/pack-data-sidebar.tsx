import {saveFileTableCols} from "@/pages/components/save-file-columns.tsx";
import {SaveFileTable} from "@/pages/components/save-file-table.tsx";
import {useSaveTableStore} from "@/store/save-collection.ts";

export default function PackDataSidebar() {
  const data = useSaveTableStore(state => state.packDataSavefiles)
  const packDataDir = useSaveTableStore(state => state.packDataDir)

  return (
    <div className="container mx-auto py-10">
      <h1>Pack Data Location: {packDataDir}</h1>
      <SaveFileTable columns={saveFileTableCols} data={data}/>
    </div>
  )
}
