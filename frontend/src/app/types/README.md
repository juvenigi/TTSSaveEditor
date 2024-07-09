## generated types
- borrow type definitions from https://github.com/matanlurey/tts-save-format
- generate the schema using typedef 
```bash
npx ts-json-schema-generator --path ttstypes.d.ts --type SaveState -o ttstypes.schema.json
```
- use said schema in `ajv` validator
- `TODO` obtain correct type inference from the schema instead of `.d.ts`

## reasons to do this opposed to using the JSON Schemas provided
- (I don't know how to use) all objects are stored in a separate schema file
- I need a type inference first, validation second
