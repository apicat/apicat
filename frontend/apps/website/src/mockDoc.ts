export const doc: any = {
  type: 'doc',
  content: [
    {
      type: 'heading',
      attrs: {
        level: 1,
      },
      content: [
        {
          type: 'text',
          text: 'heading 1',
        },
      ],
    },
    {
      type: 'heading',
      attrs: {
        level: 2,
      },
      content: [
        {
          type: 'text',
          text: 'heading 2',
        },
      ],
    },
    {
      type: 'heading',
      attrs: {
        level: 3,
      },
      content: [
        {
          type: 'text',
          text: 'heading 3',
        },
      ],
    },
    {
      type: 'heading',
      attrs: {
        level: 4,
      },
      content: [
        {
          type: 'text',
          text: 'heading 4',
        },
      ],
    },
    {
      type: 'heading',
      attrs: {
        level: 5,
      },
      content: [
        {
          type: 'text',
          text: 'heading 5',
        },
      ],
    },
    {
      type: 'heading',
      attrs: {
        level: 6,
      },
      content: [
        {
          type: 'text',
          text: 'heading 6',
        },
      ],
    },
    {
      type: 'bulletList',
      content: [
        {
          type: 'listItem',
          content: [
            {
              type: 'paragraph',
              content: [
                {
                  type: 'text',
                  text: 'list item 1',
                },
              ],
            },
          ],
        },
        {
          type: 'listItem',
          content: [
            {
              type: 'paragraph',
              content: [
                {
                  type: 'text',
                  text: 'list item 2',
                },
              ],
            },
          ],
        },
        {
          type: 'listItem',
          content: [
            {
              type: 'paragraph',
              content: [
                {
                  type: 'text',
                  text: 'list item 3',
                },
              ],
            },
          ],
        },
      ],
    },
    {
      type: 'orderedList',
      attrs: {
        start: 1,
      },
      content: [
        {
          type: 'listItem',
          content: [
            {
              type: 'paragraph',
              content: [
                {
                  type: 'text',
                  text: 'list item 1',
                },
              ],
            },
          ],
        },
        {
          type: 'listItem',
          content: [
            {
              type: 'paragraph',
              content: [
                {
                  type: 'text',
                  text: 'list item 2',
                },
              ],
            },
          ],
        },
        {
          type: 'listItem',
          content: [
            {
              type: 'paragraph',
              content: [
                {
                  type: 'text',
                  text: 'list item 3',
                },
              ],
            },
          ],
        },
      ],
    },
    {
      type: 'paragraph',
      content: [
        {
          type: 'text',
          marks: [
            {
              type: 'code',
            },
          ],
          text: 'code',
        },
        {
          type: 'text',
          text: ' ',
        },
        {
          type: 'text',
          marks: [
            {
              type: 'italic',
            },
          ],
          text: '123123',
        },
        {
          type: 'text',
          text: ' ',
        },
        {
          type: 'text',
          marks: [
            {
              type: 'bold',
            },
          ],
          text: 'strong',
        },
        {
          type: 'text',
          text: ' ',
        },
        {
          type: 'text',
          marks: [
            {
              type: 'bold',
            },
          ],
          text: 'bold',
        },
        {
          type: 'text',
          text: ' ',
        },
        {
          type: 'text',
          marks: [
            {
              type: 'link',
              attrs: {
                href: 'https://en.wikipedia.org/wiki/World_Wide_Web',
                target: '_blank',
                rel: 'noopener noreferrer nofollow',
                class: null,
              },
            },
          ],
          text: 'Link',
        },
      ],
    },
    {
      type: 'horizontalRule',
    },
    {
      type: 'paragraph',
      content: [
        {
          type: 'text',
          text: 'text',
        },
      ],
    },
    {
      type: 'blockquote',
      content: [
        {
          type: 'paragraph',
          content: [
            {
              type: 'text',
              text: 'blockquote',
            },
          ],
        },
      ],
    },
    {
      type: 'image',
      attrs: {
        src: 'https://source.unsplash.com/8xznAGy4HcY/800x400',
        alt: null,
        title: null,
        width: 200,
        height: 100,
        alignment: 'left',
      },
    },
    {
      type: 'table',
      content: [
        {
          type: 'tableRow',
          content: [
            {
              type: 'tableHeader',
              attrs: {
                colspan: 1,
                rowspan: 1,
                alignment: null,
              },
              content: [
                {
                  type: 'paragraph',
                  content: [
                    {
                      type: 'text',
                      text: 'Name',
                    },
                  ],
                },
              ],
            },
            {
              type: 'tableHeader',
              attrs: {
                colspan: 1,
                rowspan: 1,
                alignment: null,
              },
              content: [
                {
                  type: 'paragraph',
                  content: [
                    {
                      type: 'text',
                      text: 'Name',
                    },
                  ],
                },
              ],
            },
            {
              type: 'tableHeader',
              attrs: {
                colspan: 1,
                rowspan: 1,
                alignment: null,
              },
              content: [
                {
                  type: 'paragraph',
                  content: [
                    {
                      type: 'text',
                      text: 'Name',
                    },
                  ],
                },
              ],
            },
            {
              type: 'tableHeader',
              attrs: {
                colspan: 1,
                rowspan: 1,
                alignment: null,
              },
              content: [
                {
                  type: 'paragraph',
                  content: [
                    {
                      type: 'text',
                      text: 'Name',
                    },
                  ],
                },
              ],
            },
          ],
        },
        {
          type: 'tableRow',
          content: [
            {
              type: 'tableCell',
              attrs: {
                colspan: 1,
                rowspan: 1,
                alignment: null,
              },
              content: [
                {
                  type: 'paragraph',
                  content: [
                    {
                      type: 'text',
                      text: 'Cyndi Lauper',
                    },
                  ],
                },
              ],
            },
            {
              type: 'tableCell',
              attrs: {
                colspan: 1,
                rowspan: 1,
                alignment: null,
              },
              content: [
                {
                  type: 'paragraph',
                  content: [
                    {
                      type: 'text',
                      text: 'singer',
                    },
                  ],
                },
              ],
            },
            {
              type: 'tableCell',
              attrs: {
                colspan: 1,
                rowspan: 1,
                alignment: null,
              },
              content: [
                {
                  type: 'paragraph',
                  content: [
                    {
                      type: 'text',
                      text: 'songwriter',
                    },
                  ],
                },
              ],
            },
            {
              type: 'tableCell',
              attrs: {
                colspan: 1,
                rowspan: 1,
                alignment: null,
              },
              content: [
                {
                  type: 'paragraph',
                  content: [
                    {
                      type: 'text',
                      text: 'actress',
                    },
                  ],
                },
              ],
            },
          ],
        },
      ],
    },
    {
      type: 'codeBlock',
      attrs: {
        language: 'json',
      },
      content: [
        {
          type: 'text',
          text: 'const a = 1;\\n',
        },
      ],
    },
    {
      type: 'apicat-http-url',
      attrs: {
        path: '/duhan/{duhan}/{duhan}/{id}/{id3}',
        method: 'get',
      },
    },
    {
      type: 'apicat-http-request',
      attrs: {
        globalExcepts: {
          header: [],
          cookie: [],
          query: [],
          path: [],
        },
        parameters: {
          header: [],
          cookie: [],
          query: [],
          path: [
            {
              name: 'duhan',
              required: false,
              schema: {},
            },
            {
              name: 'duhan',
              required: false,
              schema: {},
            },
            {
              name: 'id',
              required: false,
              schema: {},
            },
            {
              name: 'id3',
              required: false,
              schema: {},
            },
          ],
        },
        content: {},
      },
    },
    {
      type: 'paragraph',
    },
  ],
}

export const doc2: any = {
  type: 'doc',
  content: [
    {
      type: 'apicat-http-url',
      attrs: {
        path: '/duhan/{duhan}',
        method: 'get',
      },
    },
    {
      type: 'apicat-http-request',
      attrs: {
        globalExcepts: {
          header: [],
          cookie: [],
          query: [],
          path: [],
        },
        parameters: {
          header: [],
          cookie: [],
          query: [],
          path: [
            {
              name: 'duhan',
              required: true,
              schema: {
                type: 'string',
              },
            },
          ],
        },
        content: {},
      },
    },
    {
      type: 'paragraph',
    },
  ],
}
