import type {ColumnDef} from "@tanstack/react-table";
import type {GameSaveFile} from "@/store/save-collection.ts";
import {HoverCard, HoverCardTrigger} from "@/components/ui/hover-card.tsx";
import {HoverCardContent} from "@radix-ui/react-hover-card";
import {Card, CardContent} from "@/components/ui/card.tsx";
import {getOsPathSeparator} from "@/events/event-registrar.ts";

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

          <HoverCardContent>
            <Card>
              <CardContent style={{textAlign: "left"}}>
                <p><b>{file.directory + getOsPathSeparator() + file.filename}</b></p>
              </CardContent>
            </Card>

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