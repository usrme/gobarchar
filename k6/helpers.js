import { URLSearchParams } from "https://jslib.k6.io/url/1.0.0/index.js";
import http from "k6/http";

export const searchParams = new URLSearchParams([
  ["December", "19"],
  ["February", "62"],
  ["March", "93"],
  ["May", "59"],
  ["October", "65"],
  ["September", "70"],
]);

export const duration = "10s";
export const executor = "constant-vus";
export const dRb = true;
export const vus = 1;

export const getOptions = (name) => {
  return Object.assign(
    {},
    {
      ext: {
        loadimpact: {
          projectID: `${__ENV.PROJECT_ID}`,
          name: name,
        },
      },
      discardResponseBodies: dRb,
      vus: vus,
    },
  );
};

export const getScenario = (exec) => {
  return {
    duration: duration,
    exec: exec,
    executor: executor,
  };
};

export const testUrlWithParams = (url) => {
  http.get(`${url}?${searchParams.toString()}`);
};
