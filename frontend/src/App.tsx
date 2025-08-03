import './App.css'
import {Button} from "@/components/ui/button.tsx";
import {Greet} from "../wailsjs/go/app/App";
import {useDummyStore} from "@/store/save-collection.ts";
import SaveFilePage from "@/pages/save-file-page.tsx";

export default function App() {
  const files = useDummyStore(state => state.files)
  const incFiles = useDummyStore(state => state.fetchFiles)

  return (
    <>
      <SaveFilePage></SaveFilePage>
      <div className="flex flex-col items-center justify-center">
        <Button onClick={() => Greet("somebody").then(str => console.log(str))}>Click me, baby</Button>
      </div>

      <div className="flex flex-col items-center justify-center">
        <Button onClick={incFiles}>ZustandButton {files}</Button>
      </div>
    </>
  )
}
