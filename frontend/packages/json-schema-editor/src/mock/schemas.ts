export const definitionSchemas = [
  {
    description: 'User',
    id: 14,
    name: 'User',
    parent_id: 0,
    schema: {
      properties: {
        email: {
          type: 'string',
        },
        firstName: {
          type: 'string',
        },
        id: {
          format: 'int64',
          type: 'integer',
        },
        lastName: {
          type: 'string',
        },
        password: {
          type: 'string',
        },
        phone: {
          type: 'string',
        },
        userStatus: {
          description: 'User Status',
          format: 'int32',
          type: 'integer',
        },
        username: {
          type: 'string',
        },
      },
      type: 'object',
      'x-apicat-orders': ['username', 'firstName', 'lastName', 'email', 'password', 'phone', 'userStatus', 'id'],
    },
    type: 'schema',
  },
  {
    description: 'Tag',
    id: 15,
    name: 'Tag',
    parent_id: 0,
    schema: {
      properties: {
        id: {
          format: 'int64',
          type: 'integer',
        },
        name: {
          type: 'string',
        },
      },
      type: 'object',
      'x-apicat-orders': ['id', 'name'],
    },
    type: 'schema',
  },
  {
    description: 'Pet',
    id: 16,
    name: 'Pet',
    parent_id: 0,
    schema: {
      properties: {
        category: {
          $ref: '#/definitions/schemas/18',
        },
        id: {
          format: 'int64',
          type: 'integer',
        },
        name: {
          example: 'doggie',
          type: 'string',
        },
        photoUrls: {
          items: {
            type: 'string',
          },
          type: 'array',
        },
        status: {
          description: 'pet status in the store',
          type: 'string',
        },
        tags: {
          items: {
            $ref: '#/definitions/schemas/15',
          },
          type: 'array',
        },
      },
      required: ['name', 'photoUrls'],
      type: 'object',
      'x-apicat-orders': ['id', 'category', 'name', 'photoUrls', 'tags', 'status'],
    },
    type: 'schema',
  },
  {
    description: 'Order',
    id: 17,
    name: 'Order',
    parent_id: 0,
    schema: {
      properties: {
        complete: {
          default: false,
          type: 'boolean',
        },
        id: {
          format: 'int64',
          type: 'integer',
        },
        petId: {
          format: 'int64',
          type: 'integer',
        },
        quantity: {
          format: 'int32',
          type: 'integer',
        },
        shipDate: {
          format: 'date-time',
          type: 'string',
        },
        status: {
          description: 'Order Status',
          type: 'string',
        },
      },
      type: 'object',
      'x-apicat-orders': ['status', 'complete', 'id', 'petId', 'quantity', 'shipDate'],
    },
    type: 'schema',
  },
  {
    description: 'Category',
    id: 18,
    name: 'Category',
    parent_id: 0,
    schema: {
      properties: {
        id: {
          format: 'int64',
          type: 'integer',
        },
        name: {
          type: 'string',
        },
      },
      type: 'object',
      'x-apicat-orders': ['name', 'id'],
    },
    type: 'schema',
  },
  {
    description: 'ApiResponse',
    id: 19,
    name: 'ApiResponse',
    parent_id: 0,
    schema: {
      properties: {
        code: {
          format: 'int32',
          type: 'integer',
        },
        message: {
          type: 'string',
        },
        type: {
          type: 'string',
        },
      },
      type: 'object',
      'x-apicat-orders': ['code', 'type', 'message'],
    },
    type: 'schema',
  },
]

export const exampleSchema = {
  type: 'object',
  title: 'title',
  properties: {
    field_0: {
      type: 'array',
      items: {
        type: 'string',
      },
    },
    field_1: {
      $ref: '#/definitions/schemas/18',
    },
    field_2: {
      type: 'object',
      properties: {
        field_3: {
          type: 'object',
          properties: {
            field_5: {
              type: 'string',
            },
          },
          required: ['field_5'],
        },
      },
      required: ['field_3'],
    },
    field_4: {
      type: 'array',
      items: {
        type: 'object',
        properties: {
          field_7: {
            type: 'array',
            items: {
              $ref: '#/definitions/schemas/16',
            },
          },
          field_9: {
            type: 'string',
          },
          field_8: {
            type: 'string',
          },
        },
        required: ['field_7', 'field_8', 'field_9'],
      },
    },
    field_6: {
      type: 'string',
    },
  },
  required: ['field_1', 'field_2', 'field_4', 'field_6'],
}

export const exampleSchema2 = {
  type: 'object',
  properties: {
    field_5: {
      type: 'string',
    },
  },
  required: ['field_5'],
}

export default definitionSchemas
