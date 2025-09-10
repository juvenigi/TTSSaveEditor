import {OpenInExplorer} from "../../../wailsjs/go/internal/CacheManagerApi";
import {Button} from "@/components/ui/button.tsx";
import {toast} from "sonner";


// Explicit props interface
interface OpenInExplorerBtnProps {
  path: string
  children?: React.ReactNode
}

export default function OpenInExplorerBtn({path, children}: OpenInExplorerBtnProps) {
  const callGoCode = async () => {
    try {
      console.log(path)
      await OpenInExplorer(path);
    } catch (e) {
      if (e instanceof Error) {
        toast(e.message);
      } else {
        toast("Error opening in-explorer");
      }
    }
  }

  return (
    <Button onClick={callGoCode}>{children}</Button>
  )
}