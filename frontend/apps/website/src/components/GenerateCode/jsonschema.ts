export default {
  type: 'object',
  properties: {
    res: {
      $ref: '#/definitions/duhan',
      description: '哇哈哈哈',
    },
  },
  title: 'Category',
  definitions: {
    duhan: {
      type: 'object',
      properties: {
        code: {
          type: 'string',
          title: 'code',
        },
      },
      required: ['code'],
      title: 'Result',
      name: 'Result',
      description: 'Result',
    },
  },
  description: 'Category',
}
