import Ajv from "ajv";

const ajv = new Ajv({allErrors: true});
// FIXME: I cannot guarantee that the generated SaveState JSON Schema corresponds to the typescript type
// export const validateSaveData = () => ajv.compile<SaveState>(schema);
