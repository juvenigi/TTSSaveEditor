import {createActionGroup, emptyProps, props} from "@ngrx/store";
import {GameCardFormData, SaveFile} from "./savefile.state";

export const SaveFileApiActions = createActionGroup({
  source: 'SaveFile API',
  events: {
    'Request Savefile By Path': props<{ path: string }>(),
    'Cache Active Savefile': emptyProps(),
    'Savefile Cache Success': props<{ data: SaveFile }>(),
    'Savefile Cache Failure': props<{ request: string, message: any }>(),
    'Savefile Success': props<{ data: SaveFile }>(),
    'Savefile Failure': props<{ request: string, message: any }>()
  }
});

// TODO implement effects / reducer triggers
export const SaveFileEditActions = createActionGroup({
  source: 'SaveFile Edits',
  events: {
    'Push New Object State': props<{ guid: string, edit: any, undoDepth: number }>(),
    'Squash Object Edits': props<{ guid: string }>(),
  }
});

export const CustomCardActions = createActionGroup({
  source: 'Custom Card',
  events: {
    'Push Changes': props<{cardData: GameCardFormData}>(),
    'Push Changes Fail': props<{reason: string}>()
  }
});
