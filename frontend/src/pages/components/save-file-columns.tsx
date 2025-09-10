import type {ColumnDef} from "@tanstack/react-table";
import type {GameSaveFile} from "@/store/save-collection.ts";
import {HoverCard, HoverCardContent, HoverCardTrigger} from "@/components/ui/hover-card.tsx";


export const saveFileTableCols: ColumnDef<GameSaveFile>[] = [
  {
    accessorKey: "derivedFilename",
    header: "File",
    cell: ({row}) => {
      const file = row.original as GameSaveFile;
      const object = file.objectCount === 1 ? "object" : "objects";
      const isLongFilename = file.filename.length > 30;

      const finalFilename = `${isLongFilename ? "..." + file.filename.substring(file.filename.length - 30) : file.filename}`;

      // todo: the hoverCard->Card is a hack
      // todo: need to handle click event differently (depending on packdata vs non-packdata)
      return (
        <HoverCard>
          {file.savename.length === 0
            ? <HoverCardTrigger><i>{finalFilename}</i>{` (${file.objectCount} ${object})`}</HoverCardTrigger>
            : <HoverCardTrigger>{`${file.savename} (${file.objectCount} ${object})`}</HoverCardTrigger>}
          <HoverCardContent className="text-left">
            {file.savename.length === 0 ? <i>save name is blank</i> : <p>save name: <b>{file.savename}</b></p>}
            <p>file name: <b>{finalFilename}</b></p>
          </HoverCardContent>
        </HoverCard>
      );
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