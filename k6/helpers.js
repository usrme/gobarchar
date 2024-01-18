import { URLSearchParams } from "https://jslib.k6.io/url/1.0.0/index.js";
import http from "k6/http";

const searchParams = new URLSearchParams([
  ["December", "19"],
  ["February", "62"],
  ["March", "93"],
  ["May", "59"],
  ["October", "65"],
  ["September", "70"],
]);

const duration = "10s";
const executor = "constant-vus";
const dRb = true;
const vus = 1;

const getOptions = (name) => {
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

const getScenario = (exec) => {
  return {
    duration: duration,
    exec: exec,
    executor: executor,
  };
};

export const generateOptions = (key) =>
  Object.assign({}, getOptions(`${key}-static-data`), {
    scenarios: {
      [key]: getScenario(key),
    },
  });

export const testUrlWithParams = (url) => {
  http.get(`${url}?${searchParams.toString()}`);
};
