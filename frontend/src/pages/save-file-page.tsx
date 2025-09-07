import {useSaveTableStore} from "@/store/save-collection.ts";
import {saveFileTableCols} from "@/pages/save-file-columns.tsx";
import {SaveFileTable} from "@/pages/save-file-table.tsx";
import PackDataSidebar from "@/pages/pack-data-sidebar.tsx";
import SaveMigrationDialog from "@/pages/save-migration-dialog.tsx";
import {Button} from "@/components/ui/button.tsx";
import {PanicButton} from "../../wailsjs/go/internal/CacheManagerApi";


export default function SaveFilePage() {
  const data = useSaveTableStore(store => store.activeSavefiles);
  const activeDir = useSaveTableStore(store => store.activeDir);
  const subdirs = useSaveTableStore(store => store.childSegments);
  const setActiveDir = useSaveTableStore(store => store.setActiveSegment);
  const retToParent = useSaveTableStore(store => store.retToParent);
  return (
    <>
      <h1>Save Files</h1>
      <Button onMouseDown={PanicButton}>Panic Button</Button>
      <PackDataSidebar></PackDataSidebar>
      <SaveMigrationDialog></SaveMigrationDialog>
      {folderNavigator(retToParent, activeDir, subdirs, setActiveDir)}
      <div className="container mx-auto py-10">
        <SaveFileTable columns={saveFileTableCols} data={data}/>
      </div>
    </>
  )
}

function folderNavigator(retToParent: (times: number) => void, activeDir: string[], subdirs: string[], setActiveDir: (segment: (string | "..")) => void) {
  return <div className="flex items-center gap-2 mb-6">
    directory:
    <div
      onMouseDown={() => retToParent(activeDir.length)}
      className="px-3 py-1 bg-muted text-muted-foreground rounded-md border border-input text-sm cursor-pointer hover:bg-accent hover:text-accent-foreground transition-colors"
    >🏠
    </div>
    {activeDir.map((item, idx) => (
      <div
        key={idx}
        onMouseDown={() => retToParent(activeDir.length - 1 - idx)}
        className="px-3 py-1 bg-muted text-muted-foreground rounded-md border border-input text-sm cursor-pointer hover:bg-accent hover:text-accent-foreground transition-colors"
      >
        {item}
      </div>
    ))}/
    {subdirs.length > 0 && subdirs.map((item, idx) => (
      <div
        key={idx}
        onMouseDown={() => setActiveDir(item)}
        className="px-3 py-1 bg-muted text-muted-foreground rounded-md border border-input text-sm cursor-pointer hover:bg-accent hover:text-accent-foreground transition-colors"
      >
        {item}
      </div>
    ))}
    {subdirs.length === 0 && <div
      className="px-3 py-1 bg-muted text-muted-foreground rounded-md border border-input text-sm cursor-pointer hover:bg-accent hover:text-accent-foreground transition-colors">
      <i>no further subfolders</i></div>}
  </div>;
}