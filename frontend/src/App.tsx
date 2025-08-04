import './App.css'
import SaveFilePage from "@/pages/save-file-page.tsx";
import {Button} from "@/components/ui/button.tsx";
import {GetTabletopSaves, PanicButton} from "../wailsjs/go/app/CacheManagerApi";

export default function App() {
  return (
    <>
      <SaveFilePage></SaveFilePage>

      <Button onMouseDown={() => GetTabletopSaves().finally(() => console.log("fetch complete"))}>Click Me!</Button>
      <Button onMouseDown={() => PanicButton()}>Crash me!</Button>
    </>
  )
}
