// spectral.js
// -----------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the Mozilla Public License,
// v. 2.0. If a copy of the MPL was not distributed with this file, you can obtain
// one at https://mozilla.org/MPL/2.0/.
//
// Author: ENDERZOMBI102 <enderzombi102.end@gmail.com> 2024
// Description: Spectral configuration that redirects all output to a file (api.txt)
//              in addition to printing it to the terminal.
// -----------------------------------------------------------------------------

import fs from 'fs';
import path from 'path';

// Instead of __dirname (which is not defined in ESM), we use process.cwd()
const directory = process.cwd();
const logFilePath = path.join(directory, 'api.txt');

// Clear (or create) the log file at startup
fs.writeFileSync(logFilePath, '', { flag: 'w' });

// Override process.stdout.write so that all stdout is also saved in api.txt
const originalStdoutWrite = process.stdout.write.bind(process.stdout);
process.stdout.write = (chunk, encoding, callback) => {
  fs.appendFileSync(logFilePath, chunk);
  return originalStdoutWrite(chunk, encoding, callback);
};

// Override process.stderr.write so that all stderr is also saved in api.txt
const originalStderrWrite = process.stderr.write.bind(process.stderr);
process.stderr.write = (chunk, encoding, callback) => {
  fs.appendFileSync(logFilePath, chunk);
  return originalStderrWrite(chunk, encoding, callback);
};

// Also override common console methods (optional)
console.log = (...args) => {
  const message = args.join(' ') + "\n";
  fs.appendFileSync(logFilePath, message);
};
console.error = (...args) => {
  const message = args.join(' ') + "\n";
  fs.appendFileSync(logFilePath, message);
};
console.warn = (...args) => {
  const message = args.join(' ') + "\n";
  fs.appendFileSync(logFilePath, message);
};

import ibmCloudValidationRules from '@ibm-cloud/openapi-ruleset';
import { allowedKeywords, propertyCasingConvention } from '@ibm-cloud/openapi-ruleset/src/functions';
import { schemas } from '@ibm-cloud/openapi-ruleset-utilities/src/collections';

console.log('Loaded config from spectral.js');

export default {
  extends: ibmCloudValidationRules,
  rules: {
    // Enforce allowed keywords (using "example" instead of deprecated "examples")
    'ibm-schema-keywords': {
      description: 'Disallows the use of certain keywords',
      message: '{{error}}',
      resolved: true,
      given: schemas,
      severity: 'error',
      then: {
        function: allowedKeywords,
        functionOptions: {
          keywordAllowList: [
            '$ref',
            'additionalProperties',
            'allOf',
            'anyOf',
            'default',
            'description',
            'discriminator',
            'enum',
            'example',
            'exclusiveMaximum',
            'exclusiveMinimum',
            'format',
            'items',
            'maximum',
            'maxItems',
            'maxLength',
            'maxProperties',
            'minimum',
            'minItems',
            'minLength',
            'minProperties',
            'multipleOf',
            'not',
            'oneOf',
            'pattern',
            'patternProperties',
            'properties',
            'readOnly',
            'required',
            'title',
            'type',
            'uniqueItems',
            'unevaluatedProperties',
            'writeOnly'
          ]
        }
      }
    },
    // Enforce camelCase for operation IDs
    'ibm-operationid-casing-convention': {
      description: 'Operation IDs must follow camel case',
      message: '{{error}}',
      resolved: true,
      given: schemas,
      severity: 'warn',
      then: {
        function: propertyCasingConvention,
        functionOptions: { type: 'camel' }
      }
    },
    // Enforce camelCase for property names
    'ibm-property-casing-convention': {
      description: 'Property names must follow camel case',
      message: '{{error}}',
      resolved: true,
      given: schemas,
      severity: 'warn',
      then: {
        function: propertyCasingConvention,
        functionOptions: { type: 'camel' }
      }
    },
    'ibm-property-consistent-name-and-type': 'warn',
    'ibm-request-and-response-content': 'error',
    'ibm-avoid-repeating-path-parameters': 'error',
    // Rules turned off as they are not required for this project
    'ibm-integer-attributes': 'off',
    'ibm-schema-type-format': 'off',
    'ibm-no-array-responses': 'off',
    'ibm-parameter-casing-convention': 'off',
    'ibm-collection-array-property': 'off',
    'ibm-anchored-patterns': 'off',
    'ibm-parameter-description': 'off',
    'operation-tag-defined': 'off',
    'oas3-api-servers': 'off',
    'ibm-major-version-in-path': 'off',
    'ibm-success-response-example': 'off',
    'ibm-operationid-naming-convention': 'off',
    'ibm-pagination-style': 'off',
    'ibm-avoid-inline-schemas': 'off',
    'ibm-requestbody-name': 'off',
    'ibm-error-response-schemas': 'off',
    'ibm-schema-casing-convention': 'off',
    'ibm-no-ambiguous-paths': 'off',
    'ibm-path-segment-casing-convention': 'off',
    'ibm-enum-casing-convention': 'off'
  }
};
