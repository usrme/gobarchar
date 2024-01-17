import { dRb, vus, localScenario, testLocalStatic } from './helpers.js'
import http from 'k6/http';

export const options = {
  ext: {
    loadimpact: {
      projectID: `${__ENV.PROJECT_ID}`,
      name: "local-static-data"
    }
  },

  discardResponseBodies: dRb,
  vus: vus,

  scenarios: {
    local: localScenario,
  },
};

export function local() {
  testLocalStatic();
}
