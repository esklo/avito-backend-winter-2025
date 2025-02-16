import http from 'k6/http';
import {check, sleep} from 'k6';
import {randomItem} from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

const baseUrl = 'http://host.docker.internal:8080';

const getStages = (targets) => {
    let stages = [];
    for (const targetsKey in targets) {
        stages.push({duration: '1s', target: targets[targetsKey]})
        stages.push({duration: '30s', target: targets[targetsKey]})
    }
    return stages
}

const scenarios = {
    auth: {
        executor: 'ramping-vus',
        startVUs: 0,
        stages: getStages([150, 200, 250]),
        exec: 'testAuth',
    },

    info: {
        executor: 'ramping-vus',
        startVUs: 0,
        stages: getStages([700, 1000, 1100]),
        exec: 'testInfo',
    },

    sendCoin: {
        executor: 'ramping-vus',
        startVUs: 0,
        stages: getStages([750, 1000, 1100]),
        exec: 'testSendCoin',
    },
    buy: {
        executor: 'ramping-vus',
        startVUs: 0,
        stages: getStages([1100]),
        exec: 'testBuy',
    },
};

const {SCENARIO} = __ENV;
export const options = {
    thresholds: {
        http_req_failed: [{threshold: 'rate<0.0001'}],
        http_req_duration: ['p(95)<50'],
    },
    rps: 1100,
    scenarios: SCENARIO ? {
        [SCENARIO]: scenarios[SCENARIO]
    } : {},
};

export function setup() {
    let tokens = {}

    for (let i = 0; i < 100; i++) {
        const username = `user${i}`
        const payload = JSON.stringify({
            username: username,
            password: 'password123'
        });

        const res = http.post(`${baseUrl}/api/auth`, payload, {
            headers: {'Content-Type': 'application/json'},
        });
        check("setup tokens", {'auth success': (r) => r.status === 200});

        tokens[username] = res.json('token')
    }
    return {tokens: tokens};
}

export function testInfo(data) {
    const randomToken = randomItem(Object.values(data.tokens));
    const res = http.get(`${baseUrl}/api/info`, {
        headers: {Authorization: `Bearer ${randomToken}`},
    });

    check(res, {
        'info status is 200': (r) => r.status === 200,
    });
    sleep(1);
}

export function testSendCoin(data) {
    var tokens = {...data.tokens};
    const randomUser = randomItem(Object.keys(tokens));
    delete tokens[randomUser]

    const token = randomItem(Object.values(tokens));
    const payload = JSON.stringify({
        toUser: randomUser,
        amount: 1,
    });

    const res = http.post(`${baseUrl}/api/sendCoin`, payload, {
        headers: {Authorization: `Bearer ${token}`},
    });

    check(res, {
        'sendCoin status is 200': (r) => r.status === 200,
    });
    sleep(1);
}

export function testBuy(data) {
    const randomToken = randomItem(Object.values(data.tokens));
    const items = ['pen', 'socks', 'cup'];
    const item = items[Math.floor(Math.random() * items.length)];
    const res = http.get(`${baseUrl}/api/buy/${item}`, {
        headers: {Authorization: `Bearer ${randomToken}`},
    });

    check(res, {
        'buy status is 200': (r) => r.status === 200,
    });
    sleep(1);
}


export function testAuth() {
    const payload = JSON.stringify({
        username: `user1`,
        password: 'password123'
    });

    const res = http.post(`${baseUrl}/api/auth`, payload, {
        headers: {'Content-Type': 'application/json'},
    });

    check(res, {
        'auth status is 200': (r) => r.status === 200,
    });
    sleep(1);
}