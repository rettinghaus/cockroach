// Copyright 2019 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

const fs                 = require('fs');
const assert             = require('assert');
const rejectsWithPGError = require('./rejects-with-pg-error');
const client             = require('./client');

// We orchestrate a failure here to ensure that a failing test actually causes
// the docker build to fail.
if (process.env.SHOULD_FAIL) {
  describe('failure smoke test', () => {
    it('causes the docker build to fail on a test failure', () => {
      assert.fail();
    });
  });
}

describe('select', () => {
  it('lets you select values', () => {
    return client.query("SELECT 1 as first, 2+$1 as second, ARRAY['\"','',''] as third", [3])
      .then(results => {
        assert.deepEqual(results.rows, [{
          first: 1,
          second: 5,
          third: ['"', '', '']
        }]);
      });
  });
});

describe('error cases', () => {
  const cases = [{
    name: 'not enough params',
    query: { text: 'SELECT 3', values: ['foo'] },
    msg: "bind message supplies 1 parameters, but requires 0",
    code: '08P01',
  }, {
    name: 'invalid utf8',
    query: { text: 'SELECT $1::STRING', values: [new Buffer([167])] },
    msg: "invalid UTF-8 sequence",
    code: '22021',
  }];

  cases.forEach(({ name, query, msg, code }) => {
    it(`${name} # ${query.text}`, () => {
      return rejectsWithPGError(client.query(query), { msg, code });
    });
  });
});

const NUMERIC_TYPES = ['INT', 'FLOAT', 'DECIMAL'];

describe('arrays', () => {
  it('can be selected', () => {
    return client.query('SELECT ARRAY[1, 2, 3] a')
      .then(results => {
        assert.deepEqual([1, 2, 3], results.rows[0].a);
      });
  });

  NUMERIC_TYPES.forEach(t => {
    it(`can be passed as a placeholder for a ${t}[]`, () => {
      return client.query(`SELECT $1:::${t}[] a`, [[1, 2, 3]])
        .then(results => {
          assert.deepEqual([1, 2, 3], results.rows[0].a);
        });
    });
  });
});

// Temporarily disabled until https://github.com/brianc/node-postgres/issues/3487
// gets resolved. The binary encoding in node-postgres was regressed at some
// point leading to protocol violation errors.
/*
describe('regression tests', () => {
  it('allows you to switch between format modes for arrays', () => {
    return client.query({
            text: 'SELECT $1:::int[] as b',
            values: [[1, 2, 8]],
            binary: false,
          }).then(r => {
            return client.query({
              text: 'SELECT $1:::int[] a',
              values: [[4, 5, 6]],
              binary: true,
            });
          }).then(results => {
            assert.deepEqual([4, 5, 6], results.rows[0].a);
          });
  });
})
*/
