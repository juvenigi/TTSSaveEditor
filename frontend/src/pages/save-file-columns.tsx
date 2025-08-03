import type {ColumnDef} from "@tanstack/react-table";
import type {GameSaveFile} from "@/store/save-collection.ts";

export const saveFileTableCols: ColumnDef<GameSaveFile>[] = [
  {
    accessorKey: "filename",
    header: "File",
    // meta: {
    //   className: "w-[40%] min-w-[200px] max-w-[60vw]",
    // }
  },
  {
    accessorKey: "objectCount",
    header: "objects",
  },
  {
    accessorKey: "uncachedResources",
    header: "uncached"
  },
  {
    accessorKey: "cachedRemoteResources",
    header: "cached",
  },
  {
    accessorKey: "localResources",
    header: "local",
  },
  {
    accessorKey: "packedResources",
    header: "packed",
  }
]