import {createActionGroup, props} from "@ngrx/store";
import {Directory} from "./directory.models";
import {SaveFile} from "../savefile/savefile.models";


export const DirectoryApiActions = createActionGroup({
  source: 'Directory API',
  events: {
    'Request Directory': props<{ rootPath: string }>(),
    'Delete Savefile': props<{ path: string }>(),
    'Copy Savefile': props<{ path: string, name: string }>(),
    'Write Savefile': props<{ file: SaveFile }>(),
    'Write Savefile Success': props<{ file: SaveFile }>(),
    'Directory Success Response': props<{directory: Directory}>(),
    'Directory Failure Response': props<{ request: string, message: any }>()
  }
});
