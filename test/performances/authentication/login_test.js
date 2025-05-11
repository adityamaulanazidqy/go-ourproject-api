import http from 'k6/http';
import { check, group, sleep } from 'k6';
import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";
import { textSummary } from "https://jslib.k6.io/k6-summary/0.0.1/index.js";

export const options = {
    stages: [
        { duration: '10s', target: 1 },
    ],
    thresholds: {
        http_req_duration: ['p(95)<1000'],
        http_req_failed: ['rate<0.05'],
    },
};

const BASE_URL = 'http://localhost:8673/auth';

const randomNumber = Math.floor(Math.random() * 1000000);
const testUser = {
    username: `user${randomNumber}`,
    email: `user${randomNumber}@siswa.smktiannajiyah.sch.id`,
    password: 'ValidPass123!',
    role_id: 1,
    major_id: 1,
    batch: 2023,
    photo: 'default.jpg'
};

export function setup() {
    console.log('=== DEBUG INFO ===');
    console.log('Data registrasi yang dikirim:', JSON.stringify(testUser, null, 2));
    console.log('Endpoint:', `${BASE_URL}/register`);

    const registerRes = http.post(`${BASE_URL}/register`, JSON.stringify(testUser), {
        headers: { 'Content-Type': 'application/json' },
        timeout: '30s'
    });

    console.log('Response registrasi:', {
        status: registerRes.status,
        body: registerRes.body,
        headers: registerRes.headers
    });

    if (registerRes.status !== 201) {
        console.log('=== RESPONSE ERROR DETAILS ===');
        console.log('Status:', registerRes.status);
        console.log('Body:', registerRes.body);
        console.log('Headers:', registerRes.headers);
        throw new Error(`Gagal registrasi user: ${registerRes.status} - ${registerRes.body}`);
    }

    return { user: testUser };
}

export default function (data) {
    group('Test Login', function () {
        const params = {
            headers: { 'Content-Type': 'application/json' },
            timeout: '30s'
        };

        console.log('Data login yang dikirim:', {
            email: data.user.email,
            password: data.user.password
        });

        const loginRes = http.post(`${BASE_URL}/login`, JSON.stringify({
            email: data.user.email,
            password: data.user.password
        }), params);

        console.log('Response login:', {
            status: loginRes.status,
            body: loginRes.body,
            headers: loginRes.headers
        });

        check(loginRes, {
            'Status 200': (r) => r.status === 200,
            'Token diterima': (r) => JSON.parse(r.body).token !== undefined,
        });
    });

    sleep(1);
}

export function handleSummary(data) {
    return {
        "login_test_report.html": htmlReport(data),
        stdout: textSummary(data, { indent: ' ', enableColors: true })
    };
}