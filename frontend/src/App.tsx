import './App.css'
import {Button} from "@/components/ui/button.tsx";
import {Greet} from "../wailsjs/go/app/App";
import {useFileStore} from "@/store/save-collection.ts";

export default function App() {
  const files = useFileStore(state => state.files)
  const incFiles = useFileStore(state => state.fetchFiles)

  return (
    <>
      <div className="flex flex-col items-center justify-center">
        <Button onClick={() => Greet("somebody").then(str => console.log(str))}>Click me, baby</Button>
      </div>

      <div className="flex flex-col items-center justify-center">
        <Button onClick={incFiles}>ZustandButton {files}</Button>
      </div>
    </>
  )
}
