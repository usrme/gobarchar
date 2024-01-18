import { getOptions, getScenario, testUrlWithParams } from "./helpers.js";

export const options = Object.assign({}, getOptions("local-static-data"), {
  scenarios: {
    local: getScenario("local"),
  },
});

export function local() {
  testUrlWithParams("http://localhost:8080/");
}
