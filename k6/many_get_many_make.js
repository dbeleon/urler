import { check } from "k6";
import http from 'k6/http';
import faker from 'https://cdnjs.cloudflare.com/ajax/libs/Faker/3.1.0/faker.min.js'
import { Trend, Counter } from 'k6/metrics';

const getUrlCnt = new Counter('zcounter_get_url');
const getFakeUrlCnt = new Counter('zcounter_get_fake_url');
const makeUrlCnt = new Counter('zcounter_make_url');

const urler = "http://localhost:8000"

const parseShorts = () => {
    return JSON.parse(open("shorts.json", "r")).shorts;
};

const shorts = parseShorts();

export const options = {
    scenarios: {
        constant_request_rate: {
            executor: 'constant-arrival-rate',
            rate: 5000,
            timeUnit: '1s',
            duration: '10s',
            preAllocatedVUs: 100, // how large the initial pool of VUs would be
            maxVUs: 200, // if the preAllocatedVUs are not enough, we can initialize more
        },
    },
};

const getReq = (invalidShort) => {
    var url = urler + "/"
    if (invalidShort == true) {
        url += Math.floor(Math.random() * 999999);
    } else {
        url += shorts[Math.floor(Math.random() * shorts.length)];
    }
    let res = http.get(url);
    
    let needStatus = 200;
    if (invalidShort == true) {
        needStatus = 404;
        getFakeUrlCnt.add(1);
    } else {
        getUrlCnt.add(1);
    }

    check(res, {
        "is status ok": (r) => r.status === needStatus
    });

    return res;
}

const makeUrl = () => {
    var url = urler+"/v1/url";
    const data = {
        user: Math.floor(1 + 50000 * Math.random()),
        url: "http://"+faker.internet.userName()+Math.floor(99999999 * Math.random())+".ru",
        };

    let res = http.post(url, JSON.stringify(data), {
        headers: { 'Content-Type': 'application/json' },
    });

    check(res, {
        "is status ok": (r) => r.status === 200
    }, { my_tag: res.body },);

    makeUrlCnt.add(1);

    return res.json();
}

export default function() {
    if (Math.random() < 0.5) {
        makeUrl();
    } else {
        getReq(Math.random() < 0.01);
    }
}