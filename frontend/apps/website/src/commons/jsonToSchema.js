// ref https://github.com/mohsen1/json-to-json-schema/blob/master/src/index.js
'use strict'

/*
 * Converts a JSON object to a JSON Schema
 * @param {any} json
 * @param {object} options
 * @returns {object} a json schema
 */
function convert(json, options = {}) {
  if (typeof json === 'function') {
    throw new TypeError('Can not convert a function')
  }

  if (json === undefined) {
    return {}
  }

  // primitives
  if (typeof json === 'string') {
    return { type: 'string' }
  }

  if (typeof json === 'boolean') {
    return { type: 'boolean' }
  }

  if (typeof json === 'number') {
    if (Number.isInteger(json)) {
      return { type: 'integer' }
    } else {
      return { type: 'number' }
    }
  }

  if (json === null) {
    return { type: 'null' }
  }

  if (Array.isArray(json)) {
    let schema = { type: 'array' }

    if (!json.length) {
      // default to empty array
      schema.items = {
        type: 'string',
      }
      return schema
    }

    let schemas = json.map(convert)

    // set default to first schema
    schema.items = schemas[0]

    // // if all schemas are the same use that schema for items
    // if (schemas.every((s) => isEqual(s, schemas[0]))) {
    //   schema.items = schemas[0]

    //   // if there are multiple schemas use oneOf
    // } else {
    //   schema.items = { oneOf: unique(schemas) }
    // }

    return schema
  }

  let schema = { type: 'object' }

  if (!Object.keys(json).length) {
    return schema
  }

  schema.properties = Object.keys(json).reduce((properties, key) => {
    properties[key] = convert(json[key])
    return properties
  }, {})

  return schema
}

/*
 * Removes duplicates from array using isEqual comparator
 * @param {array}
 * @return {array}
 */
function unique(arr = []) {
  return arr.reduce((result, item) => {
    // result does not contain item
    if (!result.some((i) => isEqual(i, item))) {
      result.push(item)
    }
    return result
  }, [])
}

export { convert }
