import { generateOptions, testUrlWithParams } from "./helpers.js";

export const options = generateOptions("local");

export function local() {
  testUrlWithParams("http://localhost:8080/");
}
