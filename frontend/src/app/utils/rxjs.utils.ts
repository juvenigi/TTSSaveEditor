import {filter, OperatorFunction} from "rxjs";


export function requiredFilter<T>(): OperatorFunction<T | undefined, T> {
  return filter((value): value is T => value !== undefined);
}

export function instanceOfFilter<T>(clazz: new (...args: any[]) => T): OperatorFunction<any, T> {
  return filter((value): value is T => value instanceof clazz);
}
