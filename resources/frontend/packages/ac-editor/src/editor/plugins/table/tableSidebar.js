import {Decoration, DecorationSet} from "prosemirror-view"
import {Plugin, PluginKey} from "prosemirror-state"
import {CellSelection, TableMap,isInTable } from "prosemirror-tables"
import {findTableDepth, selectRow, selectCol} from './commands'
// import {selectColumn} from "@/editor/plugins/table/select-nodes";

export const tableSidebarKey = new PluginKey("tableSidebar")

let selectedAnchorRow = -1
let selectedAnchorCol = -1

export function tableSidebar() {
  return new Plugin({
    key: tableSidebarKey,

    state: {
      init() { return null },
      apply() { return null }
    },

    props: {
      decorations: (state) => {
        if (!isInTable(state)) {
          return null
        }

        return new SidebarDecoration(state)
      },
    },
  })
}



class SidebarDecoration {
  constructor(state) {
    const tableDepth = findTableDepth(state)
    const sel = state.selection

    this.cells = []
    this.state = state
    this.selection = state.selection
    this.table = sel.$anchor.node(tableDepth)
    this.tableStart = sel.$anchor.start(tableDepth)
    this.tableEnd = this.tableStart - 1 + this.table.nodeSize - 1
    this.tableMap = TableMap.get(this.table)

    this.rect = sel instanceof CellSelection ? this.tableMap.rectBetween(sel.$anchorCell.pos - this.tableStart, sel.$headCell.pos - this.tableStart) : {}
    this.renderWidget()
    return DecorationSet.create(state.doc, this.cells)
  }

  renderWidget () {
    const {tableStart} = this

    // table widget
    this.cells.push(Decoration.widget(tableStart - 1, (view) => {
      const widget = document.createElement('div')
      widget.className = 'ProseMirror-tableSidebar'
      this.widget = widget
      this.view = view
      this.renderRow()
      this.renderColumn()

      return widget
    }, {
      // stop decoration selection
      ignoreSelection: true
    }))
  }

  renderRow () {
    const {view, state, widget, table, tableMap} = this

    // select row
    const sidebarRowContainer = widget.appendChild(document.createElement('div'))
    sidebarRowContainer.className = 'ProseMirror-tableSidebar-row-container'

    for (let row = 0; row < tableMap.height; row++) {
      // calculate row height from cell left (first col)
      const pos = tableMap.map[row * tableMap.width]
      const cell = table.nodeAt(pos)

      // skip merged cells
      row += cell.attrs.rowspan - 1

      const sidebarRow = sidebarRowContainer.appendChild(document.createElement('div'))
      sidebarRow.className = 'ProseMirror-tableSidebar-row'

      this.updateSidebarRowStatus(sidebarRow, row)

      const isLast = row === tableMap.height - 1
      // this.updateRowHeight(sidebarRow, pos, isLast)

      // calc height after prosemirror doc update
      requestAnimationFrame(() => {
        this.updateRowHeight(sidebarRow, pos, isLast)
      })

      sidebarRow.addEventListener('mousedown', (event) => {
        event.preventDefault()
        if (view.dom.classList.contains('resize-row-cursor')) return
        selectedAnchorRow = row
        selectRow(state, view.dispatch, selectedAnchorRow)
      })
    }
  }

  renderColumn () {
    const {view, state, widget, table, tableMap} = this

      const sidebarColContainer = widget.appendChild(document.createElement('div'))
      sidebarColContainer.className = 'ProseMirror-tableSidebar-col-container'

      for (let col = 0; col < tableMap.width; col++) {
        // calculate col width from cell top (first row)
        const pos = tableMap.map[col]
        const cell = table.nodeAt(pos)

        // skip merged cells
        col += cell.attrs.colspan - 1

        const sidebarCol = sidebarColContainer.appendChild(document.createElement('div'))
        sidebarCol.className = 'ProseMirror-tableSidebar-col'

        this.updateSidebarColStatus(sidebarCol, col)

        const isLast = col === tableMap.width - 1
        // this.updateColWidth(sidebarCol, pos, isLast)

        requestAnimationFrame(() => {
          this.updateColWidth(sidebarCol, pos, isLast)
        })

        sidebarCol.addEventListener('mousedown', (event) => {
          event.preventDefault()
          selectedAnchorCol = col
          selectCol(state, view.dispatch, selectedAnchorCol)
          // console.log(selectCol,col)
          // view.dispatch(selectColumn(selectedAnchorCol)(state.tr))
        })
      }
  }

  updateRowHeight (sidebarRow, pos, isLast) {
    const {view, tableStart} = this
    const dom = view.nodeDOM(tableStart + pos)
    // check dom when redo delete table
    const height = dom && dom.getBoundingClientRect ? dom.getBoundingClientRect().height : 0
    // if last cell, add table border right 1px
    if (height > 0) { sidebarRow.style.height = `${height + (isLast ? 1 : 0)}px` }
  }

  updateColWidth (sidebarCol, pos, isLast) {
    const {view, tableStart} = this
    const dom = view.nodeDOM(tableStart + pos)
    // check dom when redo delete table
    const width = dom && dom.getBoundingClientRect ? dom.getBoundingClientRect().width : 0
    // if last cell, add table border bottom 1px
    if (width > 0) { sidebarCol.style.width = `${width + (isLast ? 1 : 0)}px` }
  }

  updateSidebarRowStatus (sidebarRow, row) {
    const {tableMap, rect} = this
    const {left, right, top, bottom} = rect

    // when selected Cell (cellSelection), sidebar highlight
    if (row >= top && row < bottom) {
      sidebarRow.setAttribute('selected', '')

      // when selected row，sidebar color
      if (left === 0 && right === tableMap.width) {
        sidebarRow.setAttribute('selected-row', '')
      } else {
        sidebarRow.removeAttribute('selected')
      }
    } else {
      sidebarRow.removeAttribute('selected')
    }
  }

  updateSidebarColStatus (sidebarCol, col) {
    const {tableMap, rect} = this
    const {left, right, top, bottom} = rect

    // when selected Cell (cellSelection), sidebar highlight
    if (col >= left && col < right) {
      sidebarCol.setAttribute('selected', '')

      // when selected col，sidebar color
      if (top === 0 && bottom === tableMap.height) {
        sidebarCol.setAttribute('selected-col', '')
      } else {
        sidebarCol.removeAttribute('selected')
      }
    } else {
      sidebarCol.removeAttribute('selected')
    }
  }
}
