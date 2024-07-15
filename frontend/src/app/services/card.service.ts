import {Injectable} from "@angular/core";
import {GameCardFormData, setFontSizeOfCustomCard} from "../store/savefile/savefile.state";
import {ObjectState} from "../types/ttstypes";
import {generateJSONPatch, Operation} from "generate-json-patch";


@Injectable({
  providedIn: 'root'
})
export default class CardService {

  public createEditPatch(data: GameCardFormData, original: ObjectState) {
    const objPath = data.path;
    if(!objPath) return;
    const edit = {...original,
      LuaScriptState: JSON.stringify([data.cardText]),
      LuaScript: setFontSizeOfCustomCard(data.fontSize)
    } satisfies ObjectState
    return generateJSONPatch(original, edit)
      .map(operation => ({...operation, path: `/ObjectState/${objPath}${operation.path}`} satisfies Operation)
      );
  }

}
