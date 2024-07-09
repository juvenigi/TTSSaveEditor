import {createActionGroup, emptyProps, props} from "@ngrx/store";
import {SaveFile} from "./savefile.models";

export const SaveFileApiActions = createActionGroup({
  source: 'SaveFile API',
  events: {
    'Request Savefile By Path': props<{ path: string }>(),
    'Cache Active Savefile': emptyProps(),
    'Savefile Cache Success': props<{data: SaveFile}>(),
    'Savefile Cache Failure': props<{ request: string, message: any }>(),
    'Savefile Success': props<{ data: SaveFile }>(),
    'Savefile Failure': props<{ request: string, message: any }>()
  }
})

export const SaveFileEditActions = createActionGroup({
  source: 'SaveFile Edits',
  events: {
    'Push New Object State': props<{ guid: string, edit: any, undoDepth: number }>(),
    'Squash Object Edits': props<{ guid: string }>()
  }
})
