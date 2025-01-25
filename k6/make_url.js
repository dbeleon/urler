import { check } from "k6";
import http from 'k6/http';
import faker from 'https://cdnjs.cloudflare.com/ajax/libs/Faker/3.1.0/faker.min.js'

export default function() {
    var url = "http://localhost:8000/v1/url";
    const data = {
        user: Math.floor(1 + 50000 * Math.random()),
        url: faker.internet.url(),//"http://"+Math.floor(99999999 * Math.random())+".ru",
      };

    let res = http.post(url, JSON.stringify(data), {
        headers: { 'Content-Type': 'application/json' },
    });

    check(res, {
        "is status 200": (r) => r.status === 200
    }, { my_tag: res.body },);

    return res.json();
}