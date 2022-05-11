
// Creates a new transaction object from a given transaction
export const cloneTr = (tr) => Object.assign(Object.create(tr), tr).setTime(Date.now());
