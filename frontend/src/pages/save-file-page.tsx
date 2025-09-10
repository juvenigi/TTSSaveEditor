import {useSaveTableStore} from "@/store/save-collection.ts";
import {saveFileTableCols} from "@/pages/components/save-file-columns.tsx";
import {SaveFileTable} from "@/pages/components/save-file-table.tsx";
import SaveMigrationDialog from "@/pages/components/save-migration-dialog.tsx";
import {Button} from "@/components/ui/button.tsx";
import {GetTabletopSaves} from "../../wailsjs/go/internal/CacheManagerApi";
import OpenInExplorerBtn from "@/pages/components/open-in-explorer-btn.tsx";
import {getOsPathSeparator} from "@/events/event-registrar.ts";


export default function SaveFilePage() {
  const data = useSaveTableStore(store => store.activeSavefiles);
  const activeDir = useSaveTableStore(store => store.activeDir);
  const subDirs = useSaveTableStore(store => store.childSegments);
  const setActiveDir = useSaveTableStore(store => store.setActiveSegment);
  const retToParent = useSaveTableStore(store => store.retToParent);
  const saveDir = useSaveTableStore(store => store.rootSaveDir)
  return (
    <>
      <div className="container mx-auto py-10">
        {folderNavigator(retToParent, saveDir, activeDir, subDirs, setActiveDir)}
        <SaveFileTable columns={saveFileTableCols} data={data}/>
      </div>
      <SaveMigrationDialog></SaveMigrationDialog>
    </>
  )
}

function folderNavigator(retToParent: (times: number) => void, rootSaveDir: string, activeDir: string[], subdirs: string[], setActiveDir: (segment: (string | "..")) => void) {
  return <div className="flex items-center gap-2 mb-6">
    <OpenInExplorerBtn path={rootSaveDir + "/" + activeDir.join(getOsPathSeparator())}>
      Open Directory
    </OpenInExplorerBtn>
    <RefreshFilesBtn path={rootSaveDir}>Refresh</RefreshFilesBtn>
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

interface RefreshFilesBtnProps {
  path: string;
  children?: React.ReactNode
}

function RefreshFilesBtn({path, children}: RefreshFilesBtnProps) {
  const clear = useSaveTableStore(state => state.clearSaveFilesList)
  const update = useSaveTableStore(state => state.setChildSegments)
  const press = async () => {
    clear()
    try {
      await GetTabletopSaves(path)
    } finally {
      update()
    }
  }

  return (
    <Button onMouseDown={press}>{children}</Button>
  )
}