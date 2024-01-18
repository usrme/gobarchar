import { getOptions, getScenario, testUrlWithParams } from "./helpers.js";

export const options = Object.assign({}, getOptions("fly-static-data"), {
  scenarios: {
    fly: getScenario("fly"),
  },
});

export function fly() {
  testUrlWithParams("https://gobarchar.fly.dev/");
}
