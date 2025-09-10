import {type ColumnDef, flexRender, getCoreRowModel, useReactTable} from "@tanstack/react-table";
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "@/components/ui/table.tsx";
import {clsx} from "clsx";
import {getColumnMetaClass} from "@/pages/components/save-file-columns.tsx";
import {useMigrationDialogStore} from "@/store/migration-dialog.ts";
import {useSaveTableStore} from "@/store/save-collection.ts";
import {getOsPathSeparator} from "@/events/event-registrar.ts";
import * as React from "react";
import {MouseButton} from "@/lib/utils.ts";
import {ContextMenu, ContextMenuContent, ContextMenuItem, ContextMenuTrigger} from "@/components/ui/context-menu.tsx";


interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[]
  data: TData[]
}

export function SaveFileTable<TData, TValue>({
                                               columns,
                                               data,
                                             }: DataTableProps<TData, TValue>, packData = false) {
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    defaultColumn: {
      size: undefined, // Don't define size as it gets defined using column definition
    },
  });

  // todo: decouple from useSaveTableStore

  const saveFiles = useSaveTableStore(state => state.activeSavefiles)
  const packDataFiles = useSaveTableStore(state => state.packDataSavefiles)
  const selectFile = useMigrationDialogStore(state => state.openWith)
  const selectById = (rowIdx: string) => {
    const idx = Number.parseInt(rowIdx, 10)
    if (idx === saveFiles.length) {
      return
    }
    const save = packData ? packDataFiles[idx] : saveFiles[idx];
    return save.directory + getOsPathSeparator() + save.filename
  }

  const handleOnMouseDown = (ev: React.MouseEvent<HTMLTableRowElement>, rowId: string) => {
    ev.stopPropagation();
    ev.preventDefault();
    const saveLoc = selectById(rowId);
    if (!saveLoc) {
      return;
    }
    if (ev.button === MouseButton.Left) {
      selectFile(saveLoc, "PACK_DATA")
    }
  }

  return (
    <div className="overflow-hidden rounded-md border">
      <Table className="table-fixed w-full">
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <TableHead
                  key={header.id}
                  className={clsx("text-center", getColumnMetaClass(header.column))}
                >
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                      header.column.columnDef.header,
                      header.getContext()
                    )}
                </TableHead>
              ))}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {table.getRowModel().rows.length ? (
            table.getRowModel().rows.map((row) => (
              <TableRow
                key={row.id}
                data-state={row.getIsSelected() && "selected"}
                onMouseDown={(ev: React.MouseEvent<HTMLTableRowElement>) => handleOnMouseDown(ev, row.id)}
              >
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id} className={getColumnMetaClass(cell.column)}>
                    <ContextMenu>
                      <ContextMenuTrigger>
                        {flexRender(
                          cell.column.columnDef.cell,
                          cell.getContext()
                        )}
                      </ContextMenuTrigger>
                      <ContextMenuContent>
                        <ContextMenuItem>
                          Delete...
                        </ContextMenuItem>
                      </ContextMenuContent>

                    </ContextMenu>
                  </TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell
                colSpan={columns.length}
                className="h-24 text-center"
              >
                No results.
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
}
