export default (type) => (state) => {
  const { from, $from, to, empty } = state.selection;
  return empty
    ? type.isInSet(state.storedMarks || $from.marks())
    : state.doc.rangeHasMark(from, to, type);
};
