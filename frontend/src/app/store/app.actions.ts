import {createActionGroup, props} from "@ngrx/store";

export const AppActions = createActionGroup({
  source: 'AppState',
  events: {
    'Request Savefile': props<{ fsPath: string }>()
  }
})
export const SaveFileApiActions = createActionGroup({
  source: 'SaveFile API',
  events: {
    'Request Savefiles': props<{ fsPath: string }>(),
    'Request Savefiles Success': props<{ paths: string[], fsPath: string }>(),
    'Savefile Retrieve Success': props<{ data: string, fsPath: string }>(),
    'Savefile Retrieve Failure': props<{ message: string, statusCode: number, fsPath: string }>()
  }
})
