/* eslint-disable no-undef */
import { TestLogin } from './testCases/loginTest.js';
import { TestRegistration } from './testCases/registerTest.js'
import config from './config.js';
import { TestDeleteManageCat, TestGetManageCat, TestPutManageCat, TestPostManageCat } from './testCases/manageCat.js';
import { TestDeleteManageCatMatch, TestGetManageCatMatch, TestPostManageCatApprove, TestPostManageCatMatch, TestPostManageCatReject } from './testCases/manageCatMatch.js';
import { generateRandomNumber } from './helper.js';

export const options = {
    stages: [],
    summaryTrendStats: ['avg', 'min', 'med', 'max', 'p(95)', 'p(99)'],
};

// Ramping up from 50 VUs to 300 VUs if __ENV.LOAD_TEST is true
// eslint-disable-next-line no-undef
if (__ENV.LOAD_TEST) {
    options.stages.push(
        { target: 50, iterations: 1, duration: "15s" },
        { target: 100, iterations: 1, duration: "15s" },
        { target: 150, iterations: 1, duration: "30s" },
        { target: 200, iterations: 1, duration: "50s" },
        { target: 250, iterations: 1, duration: "1m" },
        { target: 300, iterations: 1, duration: "1m" },
        { target: 600, iterations: 1, duration: "1m" }
    );
} else {
    options.stages.push({
        target: 1,
        iterations: 1
    });
}
const positiveCaseConfig = Object.assign(config, {
    POSITIVE_CASE: true
})

const usersKv = {
    getRandomUser() {
        const keys = Object.keys(usersKv)
        const i = generateRandomNumber(0, keys.length - 1)
        if (!this.usedKeys.includes(keys[i])) {
            this.usedKeys.push(keys[i])
            return this[keys[i]]
        }
        return this.getRandomUser()
    },
    clearUsedKeys() {
        this.usedKeys = []
    },
    usedKeys: []
}

export default function () {
    let user;
    const currentTarget = options.stages[0].target;
    const currentStage = options.stages[0]; // Get the current stage
    const totalVUs = currentStage.target; // Total VUs for the current stage
    const percentageVUs10 = (__VU - 1) % Math.ceil(totalVUs / Math.round(totalVUs * 0.1)) === 0; // Calculate 10% of total VUs
    const percentageVUs20 = (__VU - 1) % Math.ceil(totalVUs / Math.round(totalVUs * 0.2)) === 0; // Calculate 20% of total VUs
    const percentageVUs30 = (__VU - 1) % Math.ceil(totalVUs / Math.round(totalVUs * 0.3)) === 0; // Calculate 30% of total VUs
    const percentageVUs40 = (__VU - 1) % Math.ceil(totalVUs / Math.round(totalVUs * 0.4)) === 0; // Calculate 40% of total VUs
    const percentageVUs50 = (__VU - 1) % Math.ceil(totalVUs / Math.round(totalVUs * 0.5)) === 0; // Calculate 50% of total VUs
    const percentageVUs60 = (__VU - 1) % Math.ceil(totalVUs / Math.round(totalVUs * 0.6)) === 0; // Calculate 60% of total VUs
    const percentageVUs70 = (__VU - 1) % Math.ceil(totalVUs / Math.round(totalVUs * 0.7)) === 0; // Calculate 70% of total VUs
    const percentageVUs80 = (__VU - 1) % Math.ceil(totalVUs / Math.round(totalVUs * 0.8)) === 0; // Calculate 80% of total VUs
    const percentageVUs90 = (__VU - 1) % Math.ceil(totalVUs / Math.round(totalVUs * 0.9)) === 0; // Calculate 90% of total VUs

    if (currentTarget === 50) {
        // auth
        user = TestRegistration(positiveCaseConfig);
        usersKv[user.email] = user
        if (percentageVUs20) {
            user = TestLogin(positiveCaseConfig, user);
        }
        // manage cat
        if (percentageVUs80) {
            let cat = TestPostManageCat(positiveCaseConfig, user);
            TestGetManageCat(positiveCaseConfig, user, cat);
            if (percentageVUs60) {
                TestPutManageCat(positiveCaseConfig, user);
                if (percentageVUs10) {
                    TestDeleteManageCat(positiveCaseConfig, user)
                }
            }
            if (percentageVUs90) {
                TestPostManageCatMatch(positiveCaseConfig, user);
            }
        }
        // match cat
        if (percentageVUs90) {
            TestPostManageCatMatch(positiveCaseConfig, user);
            TestGetManageCatMatch(positiveCaseConfig, user);
            if (percentageVUs30) {
                TestDeleteManageCatMatch(positiveCaseConfig, user, usersKv.getRandomUser())
            }
            if (percentageVUs80) {
                TestPostManageCatApprove(positiveCaseConfig, user, usersKv.getRandomUser())
            } else {
                TestPostManageCatReject(positiveCaseConfig, user, usersKv.getRandomUser())
            }
        }
    } else if (currentTarget === 100) {
        // auth
        user = TestLogin(positiveCaseConfig, usersKv.getRandomUser());
        if (percentageVUs50) {
            user = TestRegistration(positiveCaseConfig, user);
            usersKv[user.email] = user
        }
        // manage cat
        if (percentageVUs80) {
            let cat = TestPostManageCat(config, user);
            TestGetManageCat(config, user, cat);
            if (percentageVUs50) {
                TestPutManageCat(config, user);
                if (percentageVUs10) {
                    TestDeleteManageCat(config, user)
                }
            }
            if (percentageVUs90) {
                TestPostManageCatMatch(positiveCaseConfig, user);
                TestPostManageCatApprove(positiveCaseConfig, user, usersKv.getRandomUser())
            }
        }
        // match cat
        if (percentageVUs90) {
            TestPostManageCatMatch(positiveCaseConfig, user);
            TestGetManageCatMatch(positiveCaseConfig, user);
            if (percentageVUs30) {
                TestDeleteManageCatMatch(positiveCaseConfig, user, usersKv.getRandomUser())
            }
            if (percentageVUs80) {
                TestPostManageCatApprove(positiveCaseConfig, user, usersKv.getRandomUser())
            } else {
                TestPostManageCatReject(positiveCaseConfig, user, usersKv.getRandomUser())
            }
        }
    } else if (currentTarget === 200) {
        // auth
        user = TestLogin(positiveCaseConfig, usersKv.getRandomUser());
        if (percentageVUs50) {
            user = TestRegistration(positiveCaseConfig);
            usersKv[user.email] = user
        }
        if (percentageVUs80) {
            let cat = TestPostManageCat(config, user);
            TestGetManageCat(config, user, cat);
            if (percentageVUs50) {
                TestPutManageCat(config, user);
                if (percentageVUs10) {
                    TestDeleteManageCat(config, user)
                }
            }
        }
        // manage cat
        if (percentageVUs90) {
            TestPostManageCatMatch(positiveCaseConfig, user);
            TestGetManageCatMatch(positiveCaseConfig, user)
            if (percentageVUs30) {
                TestDeleteManageCatMatch(positiveCaseConfig, user, usersKv.getRandomUser())
            }
            if (percentageVUs80) {
                TestPostManageCatApprove(positiveCaseConfig, user, usersKv.getRandomUser())
            } else {
                TestPostManageCatReject(positiveCaseConfig, user, usersKv.getRandomUser())
            }
        }

        // match cat
        if (percentageVUs90) {
            if (percentageVUs20) {
                TestPostManageCatMatch(config, user);
                TestGetManageCatMatch(config, user);
            } else {
                TestPostManageCatMatch(positiveCaseConfig, user);
                TestGetManageCatMatch(positiveCaseConfig, user);
            }
            if (percentageVUs30) {
                if (percentageVUs10) {
                    TestDeleteManageCatMatch(config, user, usersKv.getRandomUser())
                } else {
                    TestDeleteManageCatMatch(positiveCaseConfig, user, usersKv.getRandomUser())
                }
            }
            if (percentageVUs80) {
                if (percentageVUs20) {
                    TestPostManageCatApprove(config, user, usersKv.getRandomUser())
                } else {
                    TestPostManageCatApprove(positiveCaseConfig, user, usersKv.getRandomUser())
                }
            } else {
                if (percentageVUs10) {
                    TestPostManageCatReject(config, user, usersKv.getRandomUser())
                } else {
                    TestPostManageCatReject(positiveCaseConfig, user, usersKv.getRandomUser())
                }
            }
        } else if (currentTarget === 300) {

            // auth
            if (percentageVUs40) {
                user = TestLogin(positiveCaseConfig, usersKv.getRandomUser());
            }
            else {
                user = TestLogin(positiveCaseConfig, usersKv.getRandomUser());
            }
            if (percentageVUs10) {
                user = TestRegistration(positiveCaseConfig);
                usersKv[user.email] = user
            }

            // manage cat
            if (percentageVUs70) {
                let cat
                if (percentageVUs30) {
                    cat = TestPostManageCat(config, user);
                    TestGetManageCat(config, user, cat);
                } else {
                    cat = TestPostManageCat(positiveCaseConfig, user);
                    TestGetManageCat(positiveCaseConfig, user, cat);
                }
                if (percentageVUs50) {
                    if (percentageVUs20) {
                        TestPutManageCat(config, user);
                    } else
                        TestPutManageCat(positiveCaseConfig, user);
                }
                if (percentageVUs10) {
                    if (percentageVUs10) {
                        TestDeleteManageCat(config, user)
                    }
                    else {
                        TestDeleteManageCat(positiveCaseConfig, user)
                    }
                }
            }

            // match cat
            if (percentageVUs90) {
                if (percentageVUs20) {
                    TestPostManageCatMatch(config, user);
                    TestGetManageCatMatch(config, user);
                } else {
                    TestPostManageCatMatch(positiveCaseConfig, user);
                    TestGetManageCatMatch(positiveCaseConfig, user);
                }
                if (percentageVUs30) {
                    if (percentageVUs10) {
                        TestDeleteManageCatMatch(config, user, usersKv.getRandomUser())
                    } else {
                        TestDeleteManageCatMatch(positiveCaseConfig, user, usersKv.getRandomUser())
                    }
                }
                if (percentageVUs80) {
                    if (percentageVUs20) {
                        TestPostManageCatApprove(config, user, usersKv.getRandomUser())
                    } else {
                        TestPostManageCatApprove(positiveCaseConfig, user, usersKv.getRandomUser())
                    }
                } else {
                    if (percentageVUs10) {
                        TestPostManageCatReject(config, user, usersKv.getRandomUser())
                    } else {
                        TestPostManageCatReject(positiveCaseConfig, user, usersKv.getRandomUser())
                    }
                }
            }

        } else if (currentTarget === 600) {
            // auth
            if (percentageVUs40) {
                user = TestLogin(positiveCaseConfig, usersKv.getRandomUser());
            }
            else {
                user = TestLogin(positiveCaseConfig, usersKv.getRandomUser());
            }
            if (percentageVUs10) {
                if (percentageVUs20) {
                    user = TestRegistration(positiveCaseConfig);
                } else {
                    user = TestRegistration(positiveCaseConfig);
                }

                usersKv[user.email] = user
            }

            // manage cat
            if (percentageVUs70) {
                let cat
                if (percentageVUs30) {
                    cat = TestPostManageCat(config, user);
                    TestGetManageCat(config, user, cat);
                } else {
                    cat = TestPostManageCat(positiveCaseConfig, user);
                    TestGetManageCat(positiveCaseConfig, user, cat);
                }
                if (percentageVUs50) {
                    if (percentageVUs20) {
                        TestPutManageCat(config, user);
                    } else
                        TestPutManageCat(positiveCaseConfig, user);
                }
                if (percentageVUs10) {
                    if (percentageVUs10) {
                        TestDeleteManageCat(config, user)
                    }
                    else {
                        TestDeleteManageCat(positiveCaseConfig, user)
                    }
                }
            }

            // match cat
            if (percentageVUs90) {
                if (percentageVUs20) {
                    TestPostManageCatMatch(config, user);
                    TestGetManageCatMatch(config, user);
                } else {
                    TestPostManageCatMatch(positiveCaseConfig, user);
                    TestGetManageCatMatch(positiveCaseConfig, user);
                }
                if (percentageVUs30) {
                    if (percentageVUs10) {
                        TestDeleteManageCatMatch(config, user, usersKv.getRandomUser())
                    } else {
                        TestDeleteManageCatMatch(positiveCaseConfig, user, usersKv.getRandomUser())
                    }
                }
                if (percentageVUs80) {
                    if (percentageVUs20) {
                        TestPostManageCatApprove(config, user, usersKv.getRandomUser())
                    } else {
                        TestPostManageCatApprove(positiveCaseConfig, user, usersKv.getRandomUser())
                    }
                } else {
                    if (percentageVUs10) {
                        TestPostManageCatReject(config, user, usersKv.getRandomUser())
                    } else {
                        TestPostManageCatReject(positiveCaseConfig, user, usersKv.getRandomUser())
                    }
                }
            }
        }


    } else {
        const users = []
        for (let index = 0; index < 3; index++) {
            user = TestRegistration(config);
            user = TestLogin(config, user);
            users.push(user)
            let cat = TestPostManageCat(config, user);
            TestGetManageCat(config, user, cat);
            TestPutManageCat(config, user);
            TestDeleteManageCat(config, user);
        }
        console.log("user credentials:", user)

        TestPostManageCatMatch(config, user);
        TestGetManageCatMatch(config, user);
        TestDeleteManageCatMatch(config, users[2], users[1],);
        TestPostManageCatApprove(config, users[2], users[1], {});
        TestPostManageCatReject(config, users[2], users[1], {});
    }
}
