import { generateOptions, testUrlWithParams } from "./helpers.js";

export const options = generateOptions("fly");

export function fly() {
  testUrlWithParams("https://gobarchar.fly.dev/");
}
