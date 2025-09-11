import {saveFileTableCols} from "@/pages/components/save-file-columns.tsx";
import {SaveFileTable} from "@/pages/components/save-file-table.tsx";
import {useSaveTableStore} from "@/store/save-collection.ts";
import SaveMigrationDialog from "@/pages/components/save-migration-dialog.tsx";

export default function PackDataSidebar() {
  const data = useSaveTableStore(state => state.packDataSavefiles)
  const packDataDir = useSaveTableStore(state => state.packDataDir)

  return (
    <div className="container mx-auto py-10">
      <h1>Pack Data Location: {packDataDir}</h1>
      <SaveFileTable columns={saveFileTableCols} data={data} migrationTo={"SAVE_DIRECTORY"}/>
      <SaveMigrationDialog/>
    </div>
  )
}
