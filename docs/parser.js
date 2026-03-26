const fs = require('fs');
const path = require('path');

class Parser {
  constructor(config = {}) {
    this.config = {
      encoding: 'utf8',
      ...config
    };
  }

  parseFile(filePath) {
    try {
      const absolutePath = path.resolve(filePath);
      if (!fs.existsSync(absolutePath)) {
        throw new Error(`File not found: ${absolutePath}`);
      }

      const fileContent = fs.readFileSync(absolutePath, this.config.encoding);
      return this._parseContent(fileContent);
    } catch (error) {
      throw new Error(`Failed to parse file: ${error.message}`);
    }
  }

  _parseContent(content) {
    try {
      return JSON.parse(content);
    } catch (error) {
      throw new Error(`Invalid JSON content: ${error.message}`);
    }
  }

  static validateSchema(data, schema) {
    if (!data || typeof data !== 'object') {
      throw new Error('Invalid data: must be an object');
    }

    for (const key in schema) {
      if (!data.hasOwnProperty(key)) {
        throw new Error(`Missing required field: ${key}`);
      }
      if (typeof data[key] !== schema[key]) {
        throw new Error(`Invalid type for field ${key}: expected ${schema[key]}`);
      }
    }

    return true;
  }
}

module.exports = Parser;