import { dRb, vus, flyScenario, testFlyStatic } from './helpers.js'

export const options = {
  ext: {
    loadimpact: {
      projectID: `${__ENV.PROJECT_ID}`,
      name: "fly-static-data"
    }
  },

  discardResponseBodies: dRb,
  vus: vus,

  scenarios: {
    fly: flyScenario,
  },
};

export function fly() {
  testFlyStatic();
}
