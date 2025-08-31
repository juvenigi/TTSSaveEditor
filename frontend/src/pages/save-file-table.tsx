import {type ColumnDef, flexRender, getCoreRowModel, useReactTable} from "@tanstack/react-table";
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "@/components/ui/table.tsx";
import {clsx} from "clsx";
import {getColumnMetaClass} from "@/pages/save-file-columns.tsx";
import {useMigrationDialogStore} from "@/store/migration-dialog.ts";
import {useSaveTableStore} from "@/store/save-collection.ts";
import {getOsPathSeparator} from "@/events/event-registrar.ts";


interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[]
  data: TData[]
}

export function SaveFileTable<TData, TValue>({
                                               columns,
                                               data,
                                             }: DataTableProps<TData, TValue>) {
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    defaultColumn: {
      size: undefined, // Don't define size
    },
  });

  // todo: decouple from useSaveTableStore

  const saveFiles = useSaveTableStore(state => state.activeSavefiles)
  const selectFile = useMigrationDialogStore(state => state.openWith)
  const selectById = (rowIdx: string) => {
    const idx = Number.parseInt(rowIdx, 10)
    if (idx === saveFiles.length) {
      return
    }
    const save = saveFiles[idx];

    const saveLoc = save.directory + getOsPathSeparator() + save.filename
    selectFile(saveLoc, "PACK_DATA")
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
                onMouseDown={() => selectById(row.id)}
              >
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id} className={getColumnMetaClass(cell.column)}>
                    {flexRender(
                      cell.column.columnDef.cell,
                      cell.getContext()
                    )}
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
