import {CellSelection, isInTable, selectionCell, TableMap} from "prosemirror-tables";
import {NodeSelection} from "prosemirror-state";

export function findTableDepth (state) {
  let $head = state.selection.$head
  for (let d = $head.depth; d > 0; d--) if ($head.node(d).type.spec.tableRole === "table") return d
  return false
}

export function selectTable (state, dispatch) {
  if (!isInTable(state)) return false

  const { tr, doc } = state
  const tableDepth = findTableDepth(state)
  const tableStart = state.selection.$anchor.start(tableDepth)

  const selection = NodeSelection.create(doc, tableStart - 1)
  tr.setSelection(selection)
  dispatch(tr)
}

export function selectRow (state, dispatch, anchorRow, headRow = anchorRow) {
  if (!isInTable(state)) return false

  const {tr, doc} = state
  let $anchorCell, $headCell

  // when pararm choice a row num
  if (anchorRow !== undefined) {
    const $pos = selectionCell(state)
    const table = $pos.node(-1)
    const tableStart = $pos.start(-1)
    const map = TableMap.get(table)

    // check anchorRow and headRow in table ranges
    if (!(
      Math.min(anchorRow, headRow) >= 0 &&
      Math.max(anchorRow, headRow) < map.height
    )) return false

    $anchorCell = doc.resolve(tableStart + map.positionAt(anchorRow, 0, table))
    $headCell = anchorRow === headRow ? $anchorCell : doc.resolve(tableStart + map.positionAt(headRow, 0, table))

  // when selected cell
  } else if (state.selection instanceof CellSelection) {
    $anchorCell = state.selection.$anchorCell
    $headCell =  state.selection.$headCell

  // when selected text
  } else {
    $headCell = $anchorCell = selectionCell(state)
  }

  const selection = CellSelection.rowSelection($anchorCell, $headCell)
  tr.setSelection(selection)
  dispatch(tr)
}

export function selectCol (state, dispatch, anchorCol, headCol = anchorCol) {
  if (!isInTable(state)) return false

  const {tr, doc} = state
  let $anchorCell, $headCell

  // when param choice a column
  if (anchorCol !== undefined) {
    const $pos = selectionCell(state)
    const table = $pos.node(-1)
    const tableStart = $pos.start(-1)
    const map = TableMap.get(table)

    // check anchorCol and headCol in table ranges
    if (anchorCol >= map.width || headCol >= map.width) return false
    if (!(
      Math.min(anchorCol, headCol) >= 0 &&
      Math.max(anchorCol, headCol) < map.width
    )) return false

    $anchorCell = doc.resolve(tableStart + map.positionAt(0, anchorCol, table))
    $headCell = headCol === anchorCol ? $anchorCell : doc.resolve(tableStart + map.positionAt(0, headCol, table))

  // when selected cell
  } else if (state.selection instanceof CellSelection) {
      $anchorCell = state.selection
    $headCell = state.selection

  // when selected text
  } else {
    $headCell = $anchorCell = selectionCell(state)
  }

  const selection = CellSelection.colSelection($anchorCell, $headCell)
  tr.setSelection(selection)
  dispatch(tr)
}
