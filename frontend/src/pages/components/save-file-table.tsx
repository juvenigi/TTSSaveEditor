import {type ColumnDef, flexRender, getCoreRowModel, useReactTable} from "@tanstack/react-table";
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "@/components/ui/table.tsx";
import {clsx} from "clsx";
import {getColumnMetaClass} from "@/pages/components/save-file-columns.tsx";
import {type MigrationDestination, useMigrationDialogStore} from "@/store/migration-dialog.ts";
import {useSaveTableStore} from "@/store/save-collection.ts";
import {getOsPathSeparator} from "@/events/event-registrar.ts";

interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[]
  data: TData[]
  migrationTo: MigrationDestination
}

export function SaveFileTable<TData, TValue>({
                                               columns,
                                               data,
                                               migrationTo,
                                             }: DataTableProps<TData, TValue>) {
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    defaultColumn: {
      size: undefined, // Don't define size as it gets defined using column definition
    },
  });
  const saveFiles = useSaveTableStore(state => state.activeSavefiles)
  const packDataFiles = useSaveTableStore(state => state.packDataSavefiles)
  const selectFile = useMigrationDialogStore(state => state.openWith)
  const selectById = (rowIdx: string) => {
    const idx = Number.parseInt(rowIdx, 10)
    const save = migrationTo !== 'PACK_DATA' ? packDataFiles[idx] : saveFiles[idx];
    const saveLoc = save.directory + getOsPathSeparator() + save.filename
    selectFile(saveLoc, migrationTo)
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
