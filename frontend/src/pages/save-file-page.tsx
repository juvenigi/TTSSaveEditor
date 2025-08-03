import {useSaveTableStore} from "@/store/save-collection.ts";
import {saveFileTableCols} from "@/pages/save-file-columns.tsx";
import {SaveFileTable} from "@/pages/save-file-table.tsx";

export default function SaveFilePage() {

  const data = useSaveTableStore(store => store.files)

  return (
    <div className="container mx-auto py-10">
      <SaveFileTable columns={saveFileTableCols} data={data}/>
    </div>
  )
}