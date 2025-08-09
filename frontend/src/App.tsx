import './App.css'
import SaveFilePage from "@/pages/save-file-page.tsx";
import {Button} from "@/components/ui/button.tsx";
import {useSaveTableStore} from "@/store/save-collection.ts";
import {GetTabletopSaves} from "../wailsjs/go/internal/CacheManagerApi";

export default function App() {
  const setSegments = useSaveTableStore(store => store.setChildSegments);
  return (
    <>
      <SaveFilePage></SaveFilePage>
      <Button
        onMouseDown={() => GetTabletopSaves().then(() => setSegments())}>Click Me!</Button>
    </>
  )
}
