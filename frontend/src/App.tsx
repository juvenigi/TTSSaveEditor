import './App.css'
import {Button} from "@/components/ui/button.tsx";
import {Greet} from "../wailsjs/go/main/App";

function App() {
    return (
        <>
            <div className="flex flex-col items-center justify-center">
                <Button onClick={() => Greet("somebody").then(str => console.log(str))}>Click me, baby</Button>
            </div>
        </>
    )
}

export default App
