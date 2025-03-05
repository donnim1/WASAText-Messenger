import fs from 'fs';
import path from 'path';
const EOL = process.platform === 'win32' ? '\r\n' : '\n';
const directory = process.cwd();
const logFilePath = path.join(directory, 'api.txt');
fs.writeFileSync(logFilePath, '', { flag: 'w' });
const stripAnsi = str => str.replace(/\x1B\[[0-9;]*m/g, '');
const origStdoutWrite = process.stdout.write.bind(process.stdout);
process.stdout.write = (chunk, encoding, callback) => {
  const text = typeof chunk === 'string' ? chunk : chunk.toString();
  fs.appendFileSync(logFilePath, stripAnsi(text));
  return origStdoutWrite(chunk, encoding, callback);
};
const origStderrWrite = process.stderr.write.bind(process.stderr);
process.stderr.write = (chunk, encoding, callback) => {
  const text = typeof chunk === 'string' ? chunk : chunk.toString();
  fs.appendFileSync(logFilePath, stripAnsi(text));
  return origStderrWrite(chunk, encoding, callback);
};
console.log = (...args) => fs.appendFileSync(logFilePath, stripAnsi(args.join(' ') + EOL));
console.error = (...args) => fs.appendFileSync(logFilePath, stripAnsi(args.join(' ') + EOL));
console.warn = (...args) => fs.appendFileSync(logFilePath, stripAnsi(args.join(' ') + EOL));
import ibmCloudValidationRules from '@ibm-cloud/openapi-ruleset';
import { allowedKeywords, propertyCasingConvention } from '@ibm-cloud/openapi-ruleset/src/functions';
import { schemas } from '@ibm-cloud/openapi-ruleset-utilities/src/collections';
console.log('Loaded config from spectral.js');
export default {
  extends: ibmCloudValidationRules,
  rules: {
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
process.on('exit', () => {
  const content = fs.readFileSync(logFilePath, 'utf8').trim();
  if (process.exitCode && process.exitCode !== 0) {
    origStderrWrite("\x1b[31m" + content + "\x1b[0m" + EOL);
  } else {
    const input = process.argv.find(a => /\.(ya?ml|json)$/.test(a)) || 'unknown file';
    origStdoutWrite("\x1b[32mValidation results for " + input + " can be found " + logFilePath + "\x1b[0m" + EOL);
  }
});
