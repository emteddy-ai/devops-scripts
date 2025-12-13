const fs = require('fs');
const path = require('path');

function parseXml(xmlFile) {
    const parser = new DOMParser();
    const xmlDoc = parser.parseFromString(fs.readFileSync(xmlFile, 'utf8'), 'text/xml');
    return xmlDoc;
}

function parseYaml(yamlFile) {
    const yaml = require('js-yaml');
    return yaml.safeLoad(fs.readFileSync(yamlFile, 'utf8'));
}

function parseJson(jsonFile) {
    const json = require('json5');
    return json.parse(fs.readFileSync(jsonFile, 'utf8'));
}

module.exports = {
    parseXml,
    parseYaml,
    parseJson
};