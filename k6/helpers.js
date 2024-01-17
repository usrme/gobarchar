import { URLSearchParams } from 'https://jslib.k6.io/url/1.0.0/index.js';
import http from 'k6/http';

export const searchParams = new URLSearchParams([
  ['December', '19'],
  ['February', '62'],
  ['March',    '93'],
  ['May',      '59'],
  ['October',  '65'],
  ['September','70'],
]);

export const duration = '10s';
export const executor = 'constant-vus'
export const dRb = true;
export const vus = 1;

export const localScenario = {
  duration: duration,
  exec: 'local',
  executor: executor,
};

export const flyScenario = {
  duration: duration,
  exec: 'fly',
  executor: executor,
};

export function testLocalStatic() {
  http.get(`${'http://localhost:8080/'}?${searchParams.toString()}`);
}

export function testFlyStatic() {
  http.get(`${'https://gobarchar.fly.dev/'}?${searchParams.toString()}`);
}
