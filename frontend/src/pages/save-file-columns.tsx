import type {ColumnDef} from "@tanstack/react-table";
import type {GameSaveFile} from "@/store/save-collection.ts";

export const saveFileTableCols: ColumnDef<GameSaveFile>[] = [
  {
    accessorKey: "derivedFilename",
    header: "File",
    cell: ({row}) => {
      const file = row.original as GameSaveFile;
      const object = file.objectCount === 1 ? "object" : "objects";
      const isLongFilename = file.filename.length > 30;

      return `${isLongFilename ? "..." + file.filename.substring(file.filename.length - 30) : file.filename} (${file.objectCount} ${object})`;
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