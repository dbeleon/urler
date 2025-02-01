import { check } from "k6";
import http from 'k6/http';
import faker from 'https://cdnjs.cloudflare.com/ajax/libs/Faker/3.1.0/faker.min.js'

export default function() {
    var url = "http://localhost:8000/v1/user";
    const data = {
        name: faker.internet.userName(),
        email: faker.internet.email()
      };

    let res = http.post(url, JSON.stringify(data), {
        headers: { 'Content-Type': 'application/json' },
    });
    check(res, {
        "is status 200": (r) => r.status === 200
    }, { my_tag: res.body },);

    return res.json();
}