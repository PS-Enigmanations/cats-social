import { TestLogin } from './testCases/loginTest.js';
import { TestRegistration } from './testCases/registerTest.js'
import config from './config.js';
import { TestDeleteManageCat, TestGetManageCat, TestPutManageCat, TestPostManageCat } from './testCases/manageCat.js';
import { TestDeleteManageCatMatch, TestGetManageCatMatch, TestPostManageCatMatch } from './testCases/manageCatMatch.js';

export const options = {
    // A number specifying the number of VUs to run concurrently.
    vus: 1,
    // A string specifying the total duration of the test run.
    iterations: 1,
};

export default function () {
    let user
    for (let index = 0; index < 3; index++) {
        user = TestRegistration(config)
        user = TestLogin(config, user)

        let cat = TestPostManageCat(config, user)
        TestGetManageCat(config, user, cat)
        TestPutManageCat(config, user)
        TestDeleteManageCat(config, user)
    }

    TestPostManageCatMatch(config, user)
    TestGetManageCatMatch(config, user)
    TestDeleteManageCatMatch(config, user)
}
