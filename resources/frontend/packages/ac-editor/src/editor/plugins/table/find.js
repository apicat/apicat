import {findParentNode } from 'prosemirror-utils';

// Iterates over parent nodes, returning the closest table node.
export const findTable = (selection)=>
  findParentNode(
    node => node.type.spec.tableRole && node.type.spec.tableRole === 'table',
  )(selection);
