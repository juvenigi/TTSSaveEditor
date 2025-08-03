import type {ColumnDef} from "@tanstack/react-table";
import type {GameSaveFile} from "@/store/save-collection.ts";

export const saveFileTableCols: ColumnDef<GameSaveFile>[] = [
  {
    accessorKey: "derivedFilename",
    header: "File",
    cell: ({row}) => {
      const file = row.original as GameSaveFile;
      return `${file.filename} (${file.objectCount} objects)`;
    },
    meta: {
      className: "w-[60%] min-w-[200px] max-w-[60vw]",
    }
  },
  {
    accessorKey: "uncachedResources",
    header: "☁️"
  },
  {
    accessorKey: "cachedRemoteResources",
    header: "📥",
  },
  {
    accessorKey: "localResources",
    header: "📁",
  },
  {
    accessorKey: "packedResources",
    header: "📦",
  }
]

export function getColumnMetaClass(input: { columnDef: { meta?: { className?: string } } }): string {
  if ('columnDef' in input) {
    return input.columnDef.meta?.className ?? '';
  }

  return '';
}