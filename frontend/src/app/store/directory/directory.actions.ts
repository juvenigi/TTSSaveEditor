import {createActionGroup, props} from "@ngrx/store";
import {Directory} from "./directory.state";
import {SaveFile} from "../savefile/savefile.state";


export const DirectoryApiActions = createActionGroup({
  source: 'Directory API',
  events: {
    'Request Directory': props<{ rootPath: string }>(),
    'Directory Success Response': props<{ directory: Directory }>(),
    'Directory Failure Response': props<{ request: string, message: any }>(),
    'Navigate To Folder': props<{ folder: string }>()
  }
});

export const IOActions = createActionGroup({
  source: 'IO',
  events: {
    'Delete Savefile': props<{ path: string }>(),
    'Copy Savefile': props<{ path: string, name: string }>(),
    'Write Savefile': props<{ file: SaveFile }>(),
    'Write Savefile Success': props<{ file: SaveFile }>(),
  }
});
